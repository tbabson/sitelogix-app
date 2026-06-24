package dailylog

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

func (r *Repository) Create(req CreateLogRequest, createdBy string) (*DailyLog, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var l DailyLog
	err = tx.QueryRowx(`
		INSERT INTO daily_logs (project_id, created_by, date, weather, notes, gps_lat, gps_lng)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, project_id, created_by, date, weather, notes, gps_lat, gps_lng, signature_url, status, created_at, updated_at`,
		req.ProjectID, createdBy, req.Date, req.Weather, req.Notes, req.GPSLat, req.GPSLng,
	).StructScan(&l)
	if err != nil {
		return nil, err
	}

	// Record any materials supplied with the log in the same transaction so the
	// log and its materials are saved atomically.
	for _, m := range req.Materials {
		if err = r.insertLogMaterial(tx, l.ID, l.ProjectID, createdBy, m); err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	if len(req.Materials) > 0 {
		l.Materials, _ = r.ListLogMaterials(l.ID)
	}
	return &l, nil
}

// insertLogMaterial records one material against a log within an existing transaction,
// deducting project stock on a best-effort basis (a savepoint isolates the deduction so
// a missing project-inventory row does not abort the surrounding transaction).
func (r *Repository) insertLogMaterial(tx *sqlx.Tx, logID, projectID, createdBy string, req AddLogMaterialRequest) error {
	if _, spErr := tx.Exec("SAVEPOINT sp_deduct"); spErr == nil {
		if deductErr := r.invRepo.DeductProjectStock(tx, req.MaterialID, projectID, req.QuantityUsed, logID, createdBy); deductErr != nil {
			tx.Exec("ROLLBACK TO SAVEPOINT sp_deduct")
		} else {
			tx.Exec("RELEASE SAVEPOINT sp_deduct")
		}
	}

	_, err := tx.Exec(`
		INSERT INTO log_materials (log_id, material_id, quantity_used)
		VALUES ($1, $2, $3)`,
		logID, req.MaterialID, req.QuantityUsed)
	return err
}

func (r *Repository) FindByID(id string) (*DailyLog, error) {
	var l DailyLog
	err := r.db.Get(&l,
		`SELECT dl.id, dl.project_id, dl.created_by, dl.date, dl.weather, dl.notes,
		        dl.gps_lat, dl.gps_lng, dl.signature_url, dl.status, dl.created_at, dl.updated_at,
		        p.name AS project_name, u.name AS creator_name
		 FROM daily_logs dl
		 JOIN projects p ON p.id = dl.project_id
		 JOIN users u ON u.id = dl.created_by
		 WHERE dl.id = $1`, id)
	return &l, err
}

func (r *Repository) List(f LogFilter) ([]DailyLog, int, error) {
	args := []interface{}{}
	conditions := []string{"1=1"}
	idx := 1

	if f.ProjectID != "" {
		conditions = append(conditions, fmt.Sprintf("dl.project_id = $%d", idx))
		args = append(args, f.ProjectID)
		idx++
	}
	if f.Username != "" {
		conditions = append(conditions, fmt.Sprintf("dl.created_by IN (SELECT id FROM users WHERE name ILIKE $%d)", idx))
		args = append(args, "%"+f.Username+"%")
		idx++
	}
	if f.Status != "" {
		conditions = append(conditions, fmt.Sprintf("dl.status = $%d", idx))
		args = append(args, f.Status)
		idx++
	}
	if f.DateFrom != "" {
		conditions = append(conditions, fmt.Sprintf("dl.date >= $%d", idx))
		args = append(args, f.DateFrom)
		idx++
	}
	if f.DateTo != "" {
		conditions = append(conditions, fmt.Sprintf("dl.date <= $%d", idx))
		args = append(args, f.DateTo)
		idx++
	}

	where := strings.Join(conditions, " AND ")
	var total int
	_ = r.db.Get(&total, fmt.Sprintf(`
		SELECT COUNT(*) FROM daily_logs dl
		JOIN projects p ON p.id = dl.project_id
		JOIN users u ON u.id = dl.created_by
		WHERE %s`, where), args...)

	offset := (f.Page - 1) * f.Limit
	args = append(args, f.Limit, offset)
	var logs []DailyLog
	err := r.db.Select(&logs, fmt.Sprintf(`
		SELECT dl.id, dl.project_id, dl.created_by, dl.date, dl.weather, dl.notes,
		       dl.gps_lat, dl.gps_lng, dl.signature_url, dl.status, dl.created_at, dl.updated_at,
		       p.name AS project_name, u.name AS creator_name
		FROM daily_logs dl
		JOIN projects p ON p.id = dl.project_id
		JOIN users u ON u.id = dl.created_by
		WHERE %s ORDER BY dl.date DESC LIMIT $%d OFFSET $%d`,
		where, idx, idx+1), args...)
	return logs, total, err
}

func (r *Repository) Update(id string, req UpdateLogRequest) (*DailyLog, error) {
	var l DailyLog
	err := r.db.QueryRowx(`
		UPDATE daily_logs SET
			weather = COALESCE($1, weather),
			notes = COALESCE($2, notes),
			gps_lat = COALESCE($3, gps_lat),
			gps_lng = COALESCE($4, gps_lng),
			signature_url = COALESCE($5, signature_url),
			updated_at = NOW()
		WHERE id = $6
		RETURNING id, project_id, created_by, date, weather, notes, gps_lat, gps_lng, signature_url, status, created_at, updated_at`,
		req.Weather, req.Notes, req.GPSLat, req.GPSLng, req.SignatureURL, id,
	).StructScan(&l)
	return &l, err
}

func (r *Repository) UpdateStatus(id string, status LogStatus) (*DailyLog, error) {
	var l DailyLog
	err := r.db.QueryRowx(`
		UPDATE daily_logs SET status = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING id, project_id, created_by, date, weather, notes, gps_lat, gps_lng, signature_url, status, created_at, updated_at`,
		status, id,
	).StructScan(&l)
	return &l, err
}

func (r *Repository) Delete(id, createdBy string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var projectID string
	if err = tx.Get(&projectID, `SELECT project_id FROM daily_logs WHERE id = $1`, id); err != nil {
		return err
	}

	type matRow struct {
		MaterialID   string  `db:"material_id"`
		QuantityUsed float64 `db:"quantity_used"`
	}
	var mats []matRow
	if err = tx.Select(&mats, `SELECT material_id, quantity_used FROM log_materials WHERE log_id = $1`, id); err != nil {
		return err
	}
	for _, m := range mats {
		// Best-effort restoration — if no project stock row exists (material was recorded
		// without deduction), the restore is a no-op and we still delete the log.
		_ = r.invRepo.RestoreProjectStock(tx, m.MaterialID, projectID, m.QuantityUsed, id, createdBy)
	}

	_, err = tx.Exec(`DELETE FROM daily_logs WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (r *Repository) ListLogMaterials(logID string) ([]LogMaterial, error) {
	var mats []LogMaterial
	err := r.db.Select(&mats, `
		SELECT lm.id, lm.log_id, lm.material_id, lm.quantity_used, lm.created_at,
		       m.name AS material_name, m.unit AS material_unit
		FROM log_materials lm
		JOIN materials m ON m.id = lm.material_id
		WHERE lm.log_id = $1 ORDER BY m.name`, logID)
	return mats, err
}

func (r *Repository) AddLogMaterial(logID, createdBy string, req AddLogMaterialRequest) (*LogMaterial, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var projectID string
	if err = tx.Get(&projectID, `SELECT project_id FROM daily_logs WHERE id = $1`, logID); err != nil {
		return nil, fmt.Errorf("log not found")
	}

	if err = r.insertLogMaterial(tx, logID, projectID, createdBy, req); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	var lm LogMaterial

	// Load the full row with material name/unit
	results, err := r.ListLogMaterials(logID)
	if err != nil {
		return &lm, nil
	}
	for _, m := range results {
		if m.MaterialID == req.MaterialID {
			return &m, nil
		}
	}
	return &lm, nil
}

func (r *Repository) RemoveLogMaterial(logID, materialID, createdBy string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var entry struct {
		QuantityUsed float64 `db:"quantity_used"`
	}
	if err = tx.Get(&entry,
		`SELECT quantity_used FROM log_materials WHERE log_id = $1 AND material_id = $2`,
		logID, materialID); err != nil {
		return fmt.Errorf("log material not found")
	}

	var projectID string
	if err = tx.Get(&projectID, `SELECT project_id FROM daily_logs WHERE id = $1`, logID); err != nil {
		return err
	}

	// Best-effort restoration — if the original add had no stock to deduct, there is
	// nothing to restore; ignore the error so removal still proceeds.
	if _, spErr := tx.Exec("SAVEPOINT sp_restore"); spErr == nil {
		if restErr := r.invRepo.RestoreProjectStock(tx, materialID, projectID, entry.QuantityUsed, logID, createdBy); restErr != nil {
			tx.Exec("ROLLBACK TO SAVEPOINT sp_restore")
		} else {
			tx.Exec("RELEASE SAVEPOINT sp_restore")
		}
	}

	_, err = tx.Exec(`DELETE FROM log_materials WHERE log_id = $1 AND material_id = $2`, logID, materialID)
	if err != nil {
		return err
	}
	return tx.Commit()
}
