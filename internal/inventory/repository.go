package inventory

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Repository struct{ db *sqlx.DB }

func NewRepository(db *sqlx.DB) *Repository { return &Repository{db: db} }

func (r *Repository) GetHeadOffice() ([]Item, error) {
	var items []Item
	err := r.db.Select(&items, `
		SELECT ii.id, ii.material_id, ii.location_type, ii.project_id,
		       ii.quantity, ii.unit_cost, ii.created_at, ii.updated_at,
		       m.name AS material_name, m.unit AS material_unit
		FROM inventory_items ii
		JOIN materials m ON m.id = ii.material_id
		WHERE ii.location_type = 'head_office'
		ORDER BY m.name`)
	return items, err
}

func (r *Repository) GetProject(projectID string) ([]Item, error) {
	var items []Item
	err := r.db.Select(&items, `
		SELECT ii.id, ii.material_id, ii.location_type, ii.project_id,
		       ii.quantity, ii.unit_cost, ii.created_at, ii.updated_at,
		       m.name AS material_name, m.unit AS material_unit
		FROM inventory_items ii
		JOIN materials m ON m.id = ii.material_id
		WHERE ii.location_type = 'project' AND ii.project_id = $1
		ORDER BY m.name`, projectID)
	return items, err
}

// AddStock upserts head-office stock (admin adds physical stock)
func (r *Repository) AddStock(req AddStockRequest, createdBy string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		INSERT INTO inventory_items (material_id, location_type, project_id, quantity, unit_cost)
		VALUES ($1, 'head_office', NULL, $2, $3)
		ON CONFLICT (material_id) WHERE location_type = 'head_office'
		DO UPDATE SET
			quantity   = inventory_items.quantity + EXCLUDED.quantity,
			unit_cost  = COALESCE(EXCLUDED.unit_cost, inventory_items.unit_cost),
			updated_at = NOW()`,
		req.MaterialID, req.Quantity, req.UnitCost)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO inventory_transactions
			(material_id, location_type, project_id, quantity_change, transaction_type, notes, created_by)
		VALUES ($1, 'head_office', NULL, $2, 'stock_in', $3, $4)`,
		req.MaterialID, req.Quantity, req.Notes, createdBy)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Transfer deducts from head-office immediately but does NOT credit the project until
// the supervisor confirms receipt via ConfirmReceipt.
func (r *Repository) Transfer(req TransferRequest, createdBy string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.Exec(`
		UPDATE inventory_items
		SET quantity = quantity - $1, updated_at = NOW()
		WHERE material_id = $2 AND location_type = 'head_office' AND quantity >= $1`,
		req.Quantity, req.MaterialID)
	if err != nil {
		return err
	}
	if n, _ := result.RowsAffected(); n == 0 {
		return fmt.Errorf("insufficient head-office stock")
	}

	_, err = tx.Exec(`
		INSERT INTO inventory_transactions
			(material_id, location_type, project_id, quantity_change, transaction_type, notes, created_by, received)
		VALUES ($1, 'head_office', NULL, $2, 'transfer_out', $3, $4, true)`,
		req.MaterialID, -req.Quantity, req.Notes, createdBy)
	if err != nil {
		return err
	}

	// transfer_in is recorded but received=false — project stock is credited only on confirmation
	_, err = tx.Exec(`
		INSERT INTO inventory_transactions
			(material_id, location_type, project_id, quantity_change, transaction_type, notes, created_by, received)
		VALUES ($1, 'project', $2, $3, 'transfer_in', $4, $5, false)`,
		req.MaterialID, req.ProjectID, req.Quantity, req.Notes, createdBy)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetPendingTransfers returns unconfirmed transfer_in records for a project.
func (r *Repository) GetPendingTransfers(projectID string) ([]Transaction, error) {
	var txs []Transaction
	err := r.db.Select(&txs, `
		SELECT it.id, it.material_id, it.location_type, it.project_id,
		       it.quantity_change, it.transaction_type, it.reference_id, it.reference_type,
		       it.notes, it.created_by, it.created_at,
		       it.received, it.received_by, it.received_at,
		       m.name AS material_name, m.unit AS material_unit,
		       p.name AS project_name, u.name AS creator_name
		FROM inventory_transactions it
		JOIN materials m ON m.id = it.material_id
		LEFT JOIN projects p ON p.id = it.project_id
		LEFT JOIN users u ON u.id = it.created_by
		WHERE it.transaction_type = 'transfer_in'
		  AND it.location_type = 'project'
		  AND it.project_id = $1
		  AND it.received = false
		ORDER BY it.created_at ASC`, projectID)
	return txs, err
}

type pendingRow struct {
	MaterialID string  `db:"material_id"`
	ProjectID  string  `db:"project_id"`
	Quantity   float64 `db:"quantity_change"`
}

// ConfirmReceipt credits project inventory and marks the transfer_in as received.
func (r *Repository) ConfirmReceipt(transactionID, receivedBy string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var row pendingRow
	err = tx.QueryRowx(`
		SELECT material_id, project_id, quantity_change
		FROM inventory_transactions
		WHERE id = $1 AND transaction_type = 'transfer_in' AND location_type = 'project' AND received = false
		FOR UPDATE`, transactionID).StructScan(&row)
	if err != nil {
		return fmt.Errorf("pending transfer not found")
	}

	_, err = tx.Exec(`
		INSERT INTO inventory_items (material_id, location_type, project_id, quantity)
		VALUES ($1, 'project', $2, $3)
		ON CONFLICT (material_id, project_id) WHERE location_type = 'project'
		DO UPDATE SET quantity = inventory_items.quantity + EXCLUDED.quantity, updated_at = NOW()`,
		row.MaterialID, row.ProjectID, row.Quantity)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE inventory_transactions
		SET received = true, received_by = $1, received_at = NOW()
		WHERE id = $2`, receivedBy, transactionID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// FulfillTransfer is called by the requisition package within an existing transaction
func (r *Repository) FulfillTransfer(tx *sqlx.Tx, materialID, projectID string, quantity float64, requisitionID, createdBy string) error {
	result, err := tx.Exec(`
		UPDATE inventory_items
		SET quantity = quantity - $1, updated_at = NOW()
		WHERE material_id = $2 AND location_type = 'head_office' AND quantity >= $1`,
		quantity, materialID)
	if err != nil {
		return err
	}
	if n, _ := result.RowsAffected(); n == 0 {
		return fmt.Errorf("insufficient head-office stock for material %s", materialID)
	}

	_, err = tx.Exec(`
		INSERT INTO inventory_items (material_id, location_type, project_id, quantity)
		VALUES ($1, 'project', $2, $3)
		ON CONFLICT (material_id, project_id) WHERE location_type = 'project'
		DO UPDATE SET quantity = inventory_items.quantity + EXCLUDED.quantity, updated_at = NOW()`,
		materialID, projectID, quantity)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`INSERT INTO inventory_transactions
		(material_id, location_type, project_id, quantity_change, transaction_type, reference_id, reference_type, created_by)
		VALUES ($1,'head_office',NULL,$2,'transfer_out',$3,'requisition',$4)`,
		materialID, -quantity, requisitionID, createdBy)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`INSERT INTO inventory_transactions
		(material_id, location_type, project_id, quantity_change, transaction_type, reference_id, reference_type, created_by)
		VALUES ($1,'project',$2,$3,'transfer_in',$4,'requisition',$5)`,
		materialID, projectID, quantity, requisitionID, createdBy)
	return err
}

func (r *Repository) Beginx() (*sqlx.Tx, error) { return r.db.Beginx() }

// DeductProjectStock deducts from a project's inventory within an existing transaction
func (r *Repository) DeductProjectStock(tx *sqlx.Tx, materialID, projectID string, quantity float64, logID, createdBy string) error {
	result, err := tx.Exec(`
		UPDATE inventory_items
		SET quantity = quantity - $1, updated_at = NOW()
		WHERE material_id = $2 AND location_type = 'project' AND project_id = $3 AND quantity >= $1`,
		quantity, materialID, projectID)
	if err != nil {
		return err
	}
	if n, _ := result.RowsAffected(); n == 0 {
		return fmt.Errorf("insufficient project stock for material %s", materialID)
	}
	_, err = tx.Exec(`
		INSERT INTO inventory_transactions
			(material_id, location_type, project_id, quantity_change, transaction_type, reference_id, reference_type, created_by)
		VALUES ($1, 'project', $2, $3, 'usage', $4, 'daily_log', $5)`,
		materialID, projectID, -quantity, logID, createdBy)
	return err
}

// RestoreProjectStock adds stock back to a project's inventory within an existing transaction
func (r *Repository) RestoreProjectStock(tx *sqlx.Tx, materialID, projectID string, quantity float64, logID, createdBy string) error {
	_, err := tx.Exec(`
		UPDATE inventory_items
		SET quantity = quantity + $1, updated_at = NOW()
		WHERE material_id = $2 AND location_type = 'project' AND project_id = $3`,
		quantity, materialID, projectID)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
		INSERT INTO inventory_transactions
			(material_id, location_type, project_id, quantity_change, transaction_type, reference_id, reference_type, created_by)
		VALUES ($1, 'project', $2, $3, 'usage_reversal', $4, 'daily_log', $5)`,
		materialID, projectID, quantity, logID, createdBy)
	return err
}

func (r *Repository) GetTransactions(materialID, projectID string, limit int) ([]Transaction, error) {
	args := []interface{}{}
	conditions := []string{"1=1"}
	idx := 1

	if materialID != "" {
		conditions = append(conditions, fmt.Sprintf("it.material_id = $%d", idx))
		args = append(args, materialID)
		idx++
	}
	if projectID != "" {
		conditions = append(conditions, fmt.Sprintf("it.project_id = $%d", idx))
		args = append(args, projectID)
		idx++
	}
	if limit <= 0 {
		limit = 50
	}
	where := strings.Join(conditions, " AND ")
	args = append(args, limit)

	var txs []Transaction
	err := r.db.Select(&txs, fmt.Sprintf(`
		SELECT it.id, it.material_id, it.location_type, it.project_id,
		       it.quantity_change, it.transaction_type, it.reference_id, it.reference_type,
		       it.notes, it.created_by, it.created_at,
		       it.received, it.received_by, it.received_at,
		       m.name AS material_name, m.unit AS material_unit,
		       p.name AS project_name, u.name AS creator_name
		FROM inventory_transactions it
		JOIN materials m ON m.id = it.material_id
		LEFT JOIN projects p ON p.id = it.project_id
		LEFT JOIN users u ON u.id = it.created_by
		WHERE %s
		ORDER BY it.created_at DESC LIMIT $%d`, where, idx), args...)
	return txs, err
}
