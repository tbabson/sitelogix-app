package media

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

func (r *Repository) Save(m *Media) (*Media, error) {
	var saved Media
	err := r.db.QueryRowx(`
		INSERT INTO media (project_id, uploaded_by, filename, content_type, url, github_sha, github_key, size, log_id, issue_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, project_id, uploaded_by, filename, content_type, url, github_sha, github_key, size, log_id, issue_id, created_at`,
		m.ProjectID, m.UploadedBy, m.Filename, m.ContentType, m.URL, m.GitHubSHA, m.GitHubKey, m.Size, m.LogID, m.IssueID,
	).StructScan(&saved)
	return &saved, err
}

func (r *Repository) FindByID(id string) (*Media, error) {
	var m Media
	err := r.db.Get(&m,
		`SELECT id, project_id, uploaded_by, filename, content_type, url, github_sha, github_key, size, log_id, issue_id, created_at
		 FROM media WHERE id = $1`, id)
	return &m, err
}

func (r *Repository) List(f MediaFilter) ([]Media, int, error) {
	args := []interface{}{}
	conditions := []string{"1=1"}
	idx := 1

	if f.ProjectID != "" {
		conditions = append(conditions, fmt.Sprintf("project_id = $%d", idx))
		args = append(args, f.ProjectID)
		idx++
	}
	if f.LogID != "" {
		conditions = append(conditions, fmt.Sprintf("log_id = $%d", idx))
		args = append(args, f.LogID)
		idx++
	}
	if f.IssueID != "" {
		conditions = append(conditions, fmt.Sprintf("issue_id = $%d", idx))
		args = append(args, f.IssueID)
		idx++
	}
	if f.DateFrom != "" {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", idx))
		args = append(args, f.DateFrom)
		idx++
	}
	if f.DateTo != "" {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", idx))
		args = append(args, f.DateTo)
		idx++
	}

	where := strings.Join(conditions, " AND ")
	var total int
	_ = r.db.Get(&total, fmt.Sprintf(`SELECT COUNT(*) FROM media WHERE %s`, where), args...)

	offset := (f.Page - 1) * f.Limit
	args = append(args, f.Limit, offset)
	var list []Media
	err := r.db.Select(&list, fmt.Sprintf(`
		SELECT id, project_id, uploaded_by, filename, content_type, url, github_sha, github_key, size, log_id, issue_id, created_at
		FROM media WHERE %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`,
		where, idx, idx+1), args...)
	return list, total, err
}

func (r *Repository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM media WHERE id = $1`, id)
	return err
}

func (r *Repository) AttachToLog(mediaID, logID string) error {
	_, err := r.db.Exec(`UPDATE media SET log_id = $1 WHERE id = $2`, logID, mediaID)
	return err
}

func (r *Repository) AttachToIssue(mediaID, issueID string) error {
	_, err := r.db.Exec(`UPDATE media SET issue_id = $1 WHERE id = $2`, issueID, mediaID)
	return err
}
