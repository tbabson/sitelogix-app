package attendance

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateWorker(req CreateWorkerRequest) (*Worker, error) {
	var w Worker
	err := r.db.QueryRowx(`
		INSERT INTO workers (name, project_id, phone, trade)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, project_id, phone, trade, user_id, created_at`,
		req.Name, req.ProjectID, req.Phone, req.Trade,
	).StructScan(&w)
	return &w, err
}

func (r *Repository) LinkUser(workerID, userID string) (*Worker, error) {
	var w Worker
	err := r.db.QueryRowx(`
		UPDATE workers SET user_id = $1 WHERE id = $2
		RETURNING id, name, project_id, phone, trade, user_id, created_at`,
		userID, workerID,
	).StructScan(&w)
	return &w, err
}

func (r *Repository) FindByUserID(userID string) ([]Worker, error) {
	var workers []Worker
	err := r.db.Select(&workers,
		`SELECT id, name, project_id, phone, trade, user_id, created_at
		 FROM workers WHERE user_id = $1`, userID)
	return workers, err
}

func (r *Repository) FindByUserAndProject(userID, projectID string) (*Worker, error) {
	var w Worker
	err := r.db.Get(&w,
		`SELECT id, name, project_id, phone, trade, user_id, created_at
		 FROM workers WHERE user_id = $1 AND project_id = $2`,
		userID, projectID)
	return &w, err
}

func (r *Repository) GetOpenAttendance(workerID string) (*Attendance, error) {
	var a Attendance
	err := r.db.Get(&a,
		`SELECT id, project_id, worker_id, check_in, check_out, notes, created_by, created_at
		 FROM attendance WHERE worker_id = $1 AND check_out IS NULL
		 ORDER BY check_in DESC LIMIT 1`, workerID)
	return &a, err
}

func (r *Repository) ListWorkers(projectID string) ([]Worker, error) {
	var workers []Worker
	err := r.db.Select(&workers,
		`SELECT id, name, project_id, phone, trade, user_id, created_at FROM workers WHERE project_id = $1 ORDER BY name`, projectID)
	return workers, err
}

func (r *Repository) CheckIn(req CheckInRequest, createdBy string) (*Attendance, error) {
	var a Attendance
	err := r.db.QueryRowx(`
		INSERT INTO attendance (project_id, worker_id, check_in, notes, created_by)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, project_id, worker_id, check_in, check_out, notes, created_by, created_at`,
		req.ProjectID, req.WorkerID, req.CheckIn, req.Notes, createdBy,
	).StructScan(&a)
	return &a, err
}

func (r *Repository) CheckOut(id string, req CheckOutRequest) (*Attendance, error) {
	var a Attendance
	err := r.db.QueryRowx(`
		UPDATE attendance SET check_out = $1 WHERE id = $2
		RETURNING id, project_id, worker_id, check_in, check_out, notes, created_by, created_at`,
		req.CheckOut, id,
	).StructScan(&a)
	return &a, err
}

func (r *Repository) List(f AttendanceFilter) ([]AttendanceWithWorker, int, error) {
	args := []interface{}{}
	conditions := []string{"1=1"}
	idx := 1

	if f.ProjectID != "" {
		conditions = append(conditions, fmt.Sprintf("a.project_id = $%d", idx))
		args = append(args, f.ProjectID)
		idx++
	}
	if f.WorkerID != "" {
		conditions = append(conditions, fmt.Sprintf("a.worker_id = $%d", idx))
		args = append(args, f.WorkerID)
		idx++
	}
	if f.DateFrom != "" {
		conditions = append(conditions, fmt.Sprintf("a.check_in >= $%d", idx))
		args = append(args, f.DateFrom)
		idx++
	}
	if f.DateTo != "" {
		conditions = append(conditions, fmt.Sprintf("a.check_in <= $%d", idx))
		args = append(args, f.DateTo)
		idx++
	}

	where := strings.Join(conditions, " AND ")
	var total int
	_ = r.db.Get(&total, fmt.Sprintf(`SELECT COUNT(*) FROM attendance a WHERE %s`, where), args...)

	offset := (f.Page - 1) * f.Limit
	args = append(args, f.Limit, offset)
	var records []AttendanceWithWorker
	err := r.db.Select(&records, fmt.Sprintf(`
		SELECT a.id, a.project_id, a.worker_id, a.check_in, a.check_out, a.notes, a.created_by, a.created_at,
		       w.name AS worker_name, w.trade AS worker_trade
		FROM attendance a
		JOIN workers w ON w.id = a.worker_id
		WHERE %s ORDER BY a.check_in DESC LIMIT $%d OFFSET $%d`,
		where, idx, idx+1), args...)
	return records, total, err
}
