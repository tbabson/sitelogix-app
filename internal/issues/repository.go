package issues

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

func (r *Repository) Create(req CreateIssueRequest, reportedBy string) (*Issue, error) {
	var i Issue
	err := r.db.QueryRowx(`
		INSERT INTO issues (project_id, title, description, priority, assigned_to, reported_by, gps_lat, gps_lng)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, project_id, title, description, priority, status, assigned_to, reported_by, gps_lat, gps_lng, created_at, updated_at`,
		req.ProjectID, req.Title, req.Description, req.Priority, req.AssignedTo, reportedBy, req.GPSLat, req.GPSLng,
	).StructScan(&i)
	return &i, err
}

func (r *Repository) FindByID(id string) (*Issue, error) {
	var i Issue
	err := r.db.Get(&i,
		`SELECT id, project_id, title, description, priority, status, assigned_to, reported_by, gps_lat, gps_lng, created_at, updated_at
		 FROM issues WHERE id = $1`, id)
	return &i, err
}

func (r *Repository) List(f IssueFilter) ([]Issue, int, error) {
	args := []interface{}{}
	conditions := []string{"1=1"}
	idx := 1

	if f.ProjectID != "" {
		conditions = append(conditions, fmt.Sprintf("project_id = $%d", idx))
		args = append(args, f.ProjectID)
		idx++
	}
	if f.AssignedTo != "" {
		conditions = append(conditions, fmt.Sprintf("assigned_to = $%d", idx))
		args = append(args, f.AssignedTo)
		idx++
	}
	if f.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", idx))
		args = append(args, f.Status)
		idx++
	}
	if f.Priority != "" {
		conditions = append(conditions, fmt.Sprintf("priority = $%d", idx))
		args = append(args, f.Priority)
		idx++
	}

	where := strings.Join(conditions, " AND ")
	var total int
	_ = r.db.Get(&total, fmt.Sprintf(`SELECT COUNT(*) FROM issues WHERE %s`, where), args...)

	offset := (f.Page - 1) * f.Limit
	args = append(args, f.Limit, offset)
	var list []Issue
	err := r.db.Select(&list, fmt.Sprintf(`
		SELECT id, project_id, title, description, priority, status, assigned_to, reported_by, gps_lat, gps_lng, created_at, updated_at
		FROM issues WHERE %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`,
		where, idx, idx+1), args...)
	return list, total, err
}

func (r *Repository) Update(id string, req UpdateIssueRequest) (*Issue, error) {
	var i Issue
	err := r.db.QueryRowx(`
		UPDATE issues SET
			title       = CASE WHEN $1::text <> '' THEN $1::text ELSE title END,
			description = COALESCE($2::text, description),
			priority    = CASE WHEN $3::text <> '' THEN $3::text::issue_priority ELSE priority END,
			status      = CASE WHEN $4::text <> '' THEN $4::text::issue_status  ELSE status   END,
			assigned_to = CASE WHEN $5::text IS NOT NULL THEN $5::uuid ELSE assigned_to END,
			gps_lat     = COALESCE($6::float8, gps_lat),
			gps_lng     = COALESCE($7::float8, gps_lng),
			updated_at  = NOW()
		WHERE id = $8::uuid
		RETURNING id, project_id, title, description, priority, status, assigned_to, reported_by, gps_lat, gps_lng, created_at, updated_at`,
		req.Title, req.Description, string(req.Priority), string(req.Status),
		req.AssignedTo, req.GPSLat, req.GPSLng, id,
	).StructScan(&i)
	return &i, err
}

func (r *Repository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM issues WHERE id = $1`, id)
	return err
}
