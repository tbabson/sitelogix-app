package project

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

func (r *Repository) Create(p CreateProjectRequest, createdBy string) (*Project, error) {
	var proj Project
	err := r.db.QueryRowx(`
		INSERT INTO projects (name, description, location, start_date, end_date, created_by)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, description, location, start_date, end_date, status, created_by, created_at, updated_at`,
		p.Name, p.Description, p.Location, p.StartDate, p.EndDate, createdBy,
	).StructScan(&proj)
	return &proj, err
}

func (r *Repository) FindByID(id string) (*Project, error) {
	var p Project
	err := r.db.Get(&p,
		`SELECT id, name, description, location, start_date, end_date, status, created_by, created_at, updated_at
		 FROM projects WHERE id = $1`, id)
	return &p, err
}

func (r *Repository) List(limit, offset int) ([]Project, int, error) {
	var projects []Project
	err := r.db.Select(&projects,
		`SELECT id, name, description, location, start_date, end_date, status, created_by, created_at, updated_at
		 FROM projects ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	var total int
	_ = r.db.Get(&total, `SELECT COUNT(*) FROM projects`)
	return projects, total, nil
}

func (r *Repository) ListForUser(userID string, limit, offset int) ([]Project, int, error) {
	var projects []Project
	err := r.db.Select(&projects, `
		SELECT p.id, p.name, p.description, p.location, p.start_date, p.end_date, p.status, p.created_by, p.created_at, p.updated_at
		FROM projects p
		INNER JOIN project_members pm ON pm.project_id = p.id
		WHERE pm.user_id = $1
		ORDER BY p.created_at DESC LIMIT $2 OFFSET $3`, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	var total int
	_ = r.db.Get(&total, `SELECT COUNT(*) FROM project_members WHERE user_id = $1`, userID)
	return projects, total, nil
}

func (r *Repository) Update(id string, p UpdateProjectRequest) (*Project, error) {
	var proj Project
	err := r.db.QueryRowx(`
		UPDATE projects SET
			name = COALESCE(NULLIF($1,''), name),
			description = COALESCE($2, description),
			location = COALESCE($3, location),
			start_date = COALESCE($4, start_date),
			end_date = COALESCE($5, end_date),
			status = COALESCE(NULLIF($6,''), status::text)::project_status,
			updated_at = NOW()
		WHERE id = $7
		RETURNING id, name, description, location, start_date, end_date, status, created_by, created_at, updated_at`,
		p.Name, p.Description, p.Location, p.StartDate, p.EndDate, string(p.Status), id,
	).StructScan(&proj)
	return &proj, err
}

func (r *Repository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM projects WHERE id = $1`, id)
	return err
}

func (r *Repository) AddMember(projectID, userID, role string) (*Member, error) {
	var m Member
	err := r.db.QueryRowx(`
		INSERT INTO project_members (project_id, user_id, role)
		VALUES ($1, $2, $3)
		ON CONFLICT (project_id, user_id) DO UPDATE SET role = EXCLUDED.role
		RETURNING id, project_id, user_id, role, created_at`,
		projectID, userID, role,
	).StructScan(&m)
	return &m, err
}

func (r *Repository) RemoveMember(projectID, userID string) error {
	_, err := r.db.Exec(`DELETE FROM project_members WHERE project_id = $1 AND user_id = $2`, projectID, userID)
	return err
}

type memberRow struct {
	ID        string    `db:"id"`
	ProjectID string    `db:"project_id"`
	UserID    string    `db:"user_id"`
	Role      string    `db:"role"`
	CreatedAt time.Time `db:"created_at"`
	UserName  string    `db:"user_name"`
	UserEmail string    `db:"user_email"`
	UserRole  string    `db:"user_role"`
	AvatarURL *string   `db:"avatar_url"`
}

func (r *Repository) ListMembers(projectID string) ([]MemberDetail, error) {
	var rows []memberRow
	err := r.db.Select(&rows, `
		SELECT pm.id, pm.project_id, pm.user_id, pm.role, pm.created_at,
		       u.name AS user_name, u.email AS user_email, u.role AS user_role, u.avatar_url
		FROM project_members pm
		JOIN users u ON u.id = pm.user_id
		WHERE pm.project_id = $1
		ORDER BY pm.created_at ASC`, projectID)
	if err != nil {
		return nil, err
	}
	members := make([]MemberDetail, len(rows))
	for i, row := range rows {
		members[i] = MemberDetail{
			ID:        row.ID,
			ProjectID: row.ProjectID,
			UserID:    row.UserID,
			Role:      row.Role,
			CreatedAt: row.CreatedAt,
			User: &MemberUser{
				Name:      row.UserName,
				Email:     row.UserEmail,
				Role:      row.UserRole,
				AvatarURL: row.AvatarURL,
			},
		}
	}
	return members, nil
}

func (r *Repository) IsMember(projectID, userID string) (bool, error) {
	var count int
	err := r.db.Get(&count,
		`SELECT COUNT(*) FROM project_members WHERE project_id = $1 AND user_id = $2`, projectID, userID)
	return count > 0, err
}

func (r *Repository) UpdateMemberRole(projectID, userID, role string) error {
	_, err := r.db.Exec(
		`UPDATE project_members SET role = $1 WHERE project_id = $2 AND user_id = $3`,
		role, projectID, userID,
	)
	return err
}

func (r *Repository) SetStatus(id string, status Status) (*Project, error) {
	var p Project
	err := r.db.QueryRowx(`
		UPDATE projects SET status = $1, updated_at = NOW() WHERE id = $2
		RETURNING id, name, description, location, start_date, end_date, status, created_by, created_at, updated_at`,
		status, id,
	).StructScan(&p)
	return &p, err
}

func (r *Repository) Duplicate(id, newName, createdBy string) (*Project, error) {
	var p Project
	err := r.db.QueryRowx(`
		INSERT INTO projects (name, description, location, start_date, end_date, created_by)
		SELECT $1, description, location, NULL, NULL, $2
		FROM projects WHERE id = $3
		RETURNING id, name, description, location, start_date, end_date, status, created_by, created_at, updated_at`,
		newName, createdBy, id,
	).StructScan(&p)
	return &p, err
}

func (r *Repository) GetStats(projectID string) (*ProjectStats, error) {
	var s ProjectStats
	err := r.db.QueryRowx(`
		SELECT
			(SELECT COUNT(*) FROM daily_logs     WHERE project_id = $1)                           AS total_logs,
			(SELECT COUNT(*) FROM daily_logs     WHERE project_id = $1 AND status = 'submitted')  AS submitted_logs,
			(SELECT COUNT(*) FROM daily_logs     WHERE project_id = $1 AND status = 'approved')   AS approved_logs,
			(SELECT COUNT(*) FROM issues         WHERE project_id = $1 AND status != 'resolved')  AS open_issues,
			(SELECT COUNT(*) FROM issues         WHERE project_id = $1 AND status = 'resolved')   AS resolved_issues,
			(SELECT COUNT(*) FROM workers        WHERE project_id = $1)                           AS total_workers,
			(SELECT COUNT(*) FROM attendance     WHERE project_id = $1
			                                       AND check_out IS NULL
			                                       AND DATE(check_in) = CURRENT_DATE)             AS workers_on_site,
			(SELECT COUNT(*) FROM project_members WHERE project_id = $1)                          AS total_members,
			(SELECT COUNT(*) FROM media          WHERE project_id = $1)                           AS total_media`,
		projectID,
	).StructScan(&s)
	if err != nil {
		return nil, err
	}
	total := s.TotalLogs
	if total > 0 {
		s.CompletionPct = float64(s.ApprovedLogs) / float64(total) * 100
	}
	return &s, nil
}

func (r *Repository) GetActivity(projectID string, limit int) ([]ActivityItem, error) {
	var items []ActivityItem
	err := r.db.Select(&items, `
		SELECT id, type, title, actor, created_at FROM (
			SELECT dl.id, 'log' AS type,
			       'Log submitted for ' || TO_CHAR(dl.date,'DD Mon YYYY') AS title,
			       u.name AS actor, dl.created_at
			FROM daily_logs dl JOIN users u ON u.id = dl.created_by
			WHERE dl.project_id = $1

			UNION ALL

			SELECT i.id, 'issue' AS type,
			       'Issue raised: ' || i.title AS title,
			       u.name AS actor, i.created_at
			FROM issues i JOIN users u ON u.id = i.reported_by
			WHERE i.project_id = $1

			UNION ALL

			SELECT m.id, 'media' AS type,
			       'File uploaded: ' || m.filename AS title,
			       u.name AS actor, m.created_at
			FROM media m JOIN users u ON u.id = m.uploaded_by
			WHERE m.project_id = $1

			UNION ALL

			SELECT pm.id, 'member' AS type,
			       u.name || ' added to project' AS title,
			       u.name AS actor, pm.created_at
			FROM project_members pm JOIN users u ON u.id = pm.user_id
			WHERE pm.project_id = $1
		) activity
		ORDER BY created_at DESC LIMIT $2`,
		projectID, limit,
	)
	return items, err
}

// ── Milestones ────────────────────────────────────────────────────────────────

func (r *Repository) CreateMilestone(projectID string, req CreateMilestoneRequest) (*Milestone, error) {
	status := req.Status
	if status == "" {
		status = MilestonePending
	}
	var m Milestone
	err := r.db.QueryRowx(`
		INSERT INTO milestones (project_id, title, due_date, status, notes)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, project_id, title, due_date, status, notes, created_at, updated_at`,
		projectID, req.Title, req.DueDate, status, req.Notes,
	).StructScan(&m)
	return &m, err
}

func (r *Repository) ListMilestones(projectID string) ([]Milestone, error) {
	var milestones []Milestone
	err := r.db.Select(&milestones,
		`SELECT id, project_id, title, due_date, status, notes, created_at, updated_at
		 FROM milestones WHERE project_id = $1 ORDER BY due_date ASC`, projectID)
	return milestones, err
}

func (r *Repository) UpdateMilestone(id string, req UpdateMilestoneRequest) (*Milestone, error) {
	var m Milestone
	err := r.db.QueryRowx(`
		UPDATE milestones SET
			title    = COALESCE(NULLIF($1,''), title),
			due_date = COALESCE($2, due_date),
			status   = COALESCE(NULLIF($3,''), status::text)::milestone_status,
			notes    = COALESCE($4, notes),
			updated_at = NOW()
		WHERE id = $5
		RETURNING id, project_id, title, due_date, status, notes, created_at, updated_at`,
		req.Title, req.DueDate, string(req.Status), req.Notes, id,
	).StructScan(&m)
	return &m, err
}

func (r *Repository) DeleteMilestone(id string) error {
	_, err := r.db.Exec(`DELETE FROM milestones WHERE id = $1`, id)
	return err
}
