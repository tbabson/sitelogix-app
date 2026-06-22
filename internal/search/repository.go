package search

import "github.com/jmoiron/sqlx"

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SearchProjects(q, userID, role string) ([]projectRow, error) {
	like := "%" + q + "%"
	var rows []projectRow

	sql := `
		SELECT p.id, p.name, p.description, p.status
		FROM projects p`

	if role != "admin" {
		sql += ` JOIN project_members pm ON pm.project_id = p.id AND pm.user_id = $2`
	}

	sql += `
		WHERE (p.name ILIKE $1 OR p.description ILIKE $1 OR p.location ILIKE $1)
		ORDER BY p.name LIMIT 10`

	if role != "admin" {
		err := r.db.Select(&rows, sql, like, userID)
		return rows, err
	}
	err := r.db.Select(&rows, sql, like)
	return rows, err
}

func (r *Repository) SearchLogs(q, userID, role string) ([]logRow, error) {
	like := "%" + q + "%"
	var rows []logRow

	sql := `
		SELECT dl.id, p.name AS project_name,
		       TO_CHAR(dl.date, 'DD Mon YYYY') AS date,
		       dl.notes, dl.status
		FROM daily_logs dl
		JOIN projects p ON p.id = dl.project_id`

	if role != "admin" {
		sql += ` JOIN project_members pm ON pm.project_id = dl.project_id AND pm.user_id = $2`
	}

	sql += `
		WHERE (dl.notes ILIKE $1 OR dl.weather ILIKE $1)
		ORDER BY dl.date DESC LIMIT 10`

	if role != "admin" {
		err := r.db.Select(&rows, sql, like, userID)
		return rows, err
	}
	err := r.db.Select(&rows, sql, like)
	return rows, err
}

func (r *Repository) SearchIssues(q, userID, role string) ([]issueRow, error) {
	like := "%" + q + "%"
	var rows []issueRow

	sql := `
		SELECT i.id, p.name AS project_name, i.title, i.description, i.priority, i.status
		FROM issues i
		JOIN projects p ON p.id = i.project_id`

	if role != "admin" {
		sql += ` JOIN project_members pm ON pm.project_id = i.project_id AND pm.user_id = $2`
	}

	sql += `
		WHERE (i.title ILIKE $1 OR i.description ILIKE $1)
		ORDER BY i.created_at DESC LIMIT 10`

	if role != "admin" {
		err := r.db.Select(&rows, sql, like, userID)
		return rows, err
	}
	err := r.db.Select(&rows, sql, like)
	return rows, err
}

func (r *Repository) SearchWorkers(q, userID, role string) ([]workerRow, error) {
	like := "%" + q + "%"
	var rows []workerRow

	sql := `
		SELECT w.id, w.name, p.name AS project_name, w.trade
		FROM workers w
		JOIN projects p ON p.id = w.project_id`

	if role != "admin" {
		sql += ` JOIN project_members pm ON pm.project_id = w.project_id AND pm.user_id = $2`
	}

	sql += `
		WHERE (w.name ILIKE $1 OR w.trade ILIKE $1)
		ORDER BY w.name LIMIT 10`

	if role != "admin" {
		err := r.db.Select(&rows, sql, like, userID)
		return rows, err
	}
	err := r.db.Select(&rows, sql, like)
	return rows, err
}
