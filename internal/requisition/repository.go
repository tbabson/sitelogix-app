package requisition

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sitelogix/backend/internal/inventory"
)

type Repository struct {
	db      *sqlx.DB
	invRepo *inventory.Repository
}

func NewRepository(db *sqlx.DB, invRepo *inventory.Repository) *Repository {
	return &Repository{db: db, invRepo: invRepo}
}

func (r *Repository) Create(req CreateRequest, requestedBy string) (*Requisition, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var rq Requisition
	err = tx.QueryRowx(`
		INSERT INTO requisitions (project_id, requested_by, notes)
		VALUES ($1, $2, $3)
		RETURNING id, project_id, requested_by, status::text AS status, notes, created_at, updated_at`,
		req.ProjectID, requestedBy, req.Notes,
	).StructScan(&rq)
	if err != nil {
		return nil, err
	}

	for _, item := range req.Items {
		_, err = tx.Exec(`
			INSERT INTO requisition_items (requisition_id, material_id, quantity_requested)
			VALUES ($1, $2, $3)`,
			rq.ID, item.MaterialID, item.QuantityRequested)
		if err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return r.FindByID(rq.ID)
}

func (r *Repository) FindByID(id string) (*Requisition, error) {
	var rq Requisition
	err := r.db.Get(&rq, `
		SELECT r.id, r.project_id, r.requested_by, r.status::text AS status, r.notes,
		       r.created_at, r.updated_at,
		       p.name AS project_name, u.name AS requester_name
		FROM requisitions r
		JOIN projects p ON p.id = r.project_id
		JOIN users u ON u.id = r.requested_by
		WHERE r.id = $1`, id)
	if err != nil {
		return nil, err
	}

	err = r.db.Select(&rq.Items, `
		SELECT ri.id, ri.requisition_id, ri.material_id,
		       ri.quantity_requested, ri.quantity_approved, ri.created_at,
		       m.name AS material_name, m.unit AS material_unit
		FROM requisition_items ri
		JOIN materials m ON m.id = ri.material_id
		WHERE ri.requisition_id = $1 ORDER BY m.name`, id)
	return &rq, err
}

func (r *Repository) List(projectID, status string) ([]Requisition, error) {
	args := []interface{}{}
	conditions := []string{"1=1"}
	idx := 1

	if projectID != "" {
		conditions = append(conditions, fmt.Sprintf("r.project_id = $%d", idx))
		args = append(args, projectID)
		idx++
	}
	if status != "" {
		conditions = append(conditions, fmt.Sprintf("r.status::text = $%d", idx))
		args = append(args, status)
		idx++
	}
	_ = idx
	where := strings.Join(conditions, " AND ")

	var list []Requisition
	err := r.db.Select(&list, fmt.Sprintf(`
		SELECT r.id, r.project_id, r.requested_by, r.status::text AS status, r.notes,
		       r.created_at, r.updated_at,
		       p.name AS project_name, u.name AS requester_name
		FROM requisitions r
		JOIN projects p ON p.id = r.project_id
		JOIN users u ON u.id = r.requested_by
		WHERE %s ORDER BY r.created_at DESC`, where), args...)
	return list, err
}

func (r *Repository) Approve(id string, req ApproveRequest) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, item := range req.Items {
		_, err = tx.Exec(`UPDATE requisition_items SET quantity_approved = $1 WHERE id = $2::uuid`,
			item.QuantityApproved, item.ItemID)
		if err != nil {
			return err
		}
	}

	_, err = tx.Exec(`
		UPDATE requisitions
		SET status = 'approved', notes = COALESCE($1, notes), updated_at = NOW()
		WHERE id = $2`, req.Notes, id)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (r *Repository) Reject(id string, notes *string) error {
	_, err := r.db.Exec(`
		UPDATE requisitions
		SET status = 'rejected', notes = COALESCE($1, notes), updated_at = NOW()
		WHERE id = $2`, notes, id)
	return err
}

func (r *Repository) Fulfill(id, createdBy string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var rq struct {
		ProjectID string `db:"project_id"`
		Status    string `db:"status"`
	}
	if err = tx.Get(&rq, `SELECT project_id, status::text AS status FROM requisitions WHERE id = $1`, id); err != nil {
		return err
	}
	if rq.Status != "approved" {
		return fmt.Errorf("requisition must be approved before fulfillment")
	}

	type itemRow struct {
		MaterialID       string   `db:"material_id"`
		QuantityApproved *float64 `db:"quantity_approved"`
	}
	var items []itemRow
	if err = tx.Select(&items, `SELECT material_id, quantity_approved FROM requisition_items WHERE requisition_id = $1`, id); err != nil {
		return err
	}

	for _, it := range items {
		qty := 0.0
		if it.QuantityApproved != nil {
			qty = *it.QuantityApproved
		}
		if qty <= 0 {
			continue
		}
		if err = r.invRepo.FulfillTransfer(tx, it.MaterialID, rq.ProjectID, qty, id, createdBy); err != nil {
			return err
		}
	}

	_, err = tx.Exec(`UPDATE requisitions SET status = 'fulfilled', updated_at = NOW() WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (r *Repository) Cancel(id, userID, role string) error {
	var reqBy string
	if err := r.db.Get(&reqBy,
		`SELECT requested_by FROM requisitions WHERE id = $1 AND status = 'pending'`, id); err != nil {
		return fmt.Errorf("requisition not found or not in pending status")
	}
	if reqBy != userID && role != "admin" {
		return fmt.Errorf("not authorized to cancel this requisition")
	}
	_, err := r.db.Exec(`DELETE FROM requisitions WHERE id = $1`, id)
	return err
}
