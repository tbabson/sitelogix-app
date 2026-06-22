package reports

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetProject(id string) (*ProjectData, error) {
	var p ProjectData
	err := r.db.QueryRowx(`
		SELECT id, name, description, location, status,
		       TO_CHAR(start_date, 'DD Mon YYYY') AS start_date,
		       TO_CHAR(end_date, 'DD Mon YYYY') AS end_date
		FROM projects WHERE id = $1`, id).StructScan(&p)
	return &p, err
}

func (r *Repository) GetLogs(projectID string, from, to *time.Time) ([]LogData, error) {
	args := []interface{}{projectID}
	q := `
		SELECT dl.id, TO_CHAR(dl.date, 'DD Mon YYYY') AS date,
		       u.name AS created_by_name, dl.weather, dl.notes, dl.status,
		       dl.gps_lat, dl.gps_lng
		FROM daily_logs dl
		JOIN users u ON u.id = dl.created_by
		WHERE dl.project_id = $1`

	if from != nil {
		args = append(args, from)
		q += ` AND dl.date >= $2`
	}
	if to != nil {
		args = append(args, to)
		if from != nil {
			q += ` AND dl.date <= $3`
		} else {
			q += ` AND dl.date <= $2`
		}
	}
	q += ` ORDER BY dl.date ASC`

	var logs []LogData
	err := r.db.Select(&logs, q, args...)
	return logs, err
}

func (r *Repository) GetIssues(projectID string) ([]IssueData, error) {
	var issues []IssueData
	err := r.db.Select(&issues, `
		SELECT i.id, i.title, i.priority, i.status,
		       ur.name AS reported_by_name,
		       ua.name AS assigned_to_name,
		       TO_CHAR(i.created_at, 'DD Mon YYYY') AS created_at
		FROM issues i
		JOIN users ur ON ur.id = i.reported_by
		LEFT JOIN users ua ON ua.id = i.assigned_to
		WHERE i.project_id = $1
		ORDER BY i.created_at ASC`, projectID)
	return issues, err
}

func (r *Repository) GetAttendance(projectID string, from, to *time.Time) ([]AttendanceData, error) {
	args := []interface{}{projectID}
	q := `
		SELECT w.name AS worker_name, w.trade,
		       TO_CHAR(a.check_in, 'DD Mon YYYY HH24:MI') AS check_in,
		       TO_CHAR(a.check_out, 'DD Mon YYYY HH24:MI') AS check_out
		FROM attendance a
		JOIN workers w ON w.id = a.worker_id
		WHERE a.project_id = $1`

	if from != nil {
		args = append(args, from)
		q += ` AND DATE(a.check_in) >= $2`
	}
	if to != nil {
		args = append(args, to)
		if from != nil {
			q += ` AND DATE(a.check_in) <= $3`
		} else {
			q += ` AND DATE(a.check_in) <= $2`
		}
	}
	q += ` ORDER BY a.check_in ASC`

	var records []AttendanceData
	err := r.db.Select(&records, q, args...)
	return records, err
}
