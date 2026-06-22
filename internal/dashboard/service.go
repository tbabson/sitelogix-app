package dashboard

import (
	"github.com/jmoiron/sqlx"
)

type Service struct {
	db *sqlx.DB
}

func NewService(db *sqlx.DB) *Service {
	return &Service{db: db}
}

func (s *Service) GetStats(role, userID string) (*Stats, error) {
	var stats Stats

	projectQuery := `SELECT COUNT(*) FROM projects WHERE status = 'active'`
	workersQuery := `SELECT COUNT(*) FROM attendance WHERE check_out IS NULL AND DATE(check_in) = CURRENT_DATE`
	issuesQuery := `SELECT COUNT(*) FROM issues WHERE status != 'resolved'`
	logsQuery := `SELECT COUNT(*) FROM daily_logs WHERE DATE(created_at) = CURRENT_DATE`
	pendingQuery := `SELECT COUNT(*) FROM daily_logs WHERE status = 'submitted'`
	membersQuery := `SELECT COUNT(*) FROM users WHERE is_active = true`

	if role != "admin" {
		projectQuery = `
			SELECT COUNT(DISTINCT p.id) FROM projects p
			JOIN project_members pm ON pm.project_id = p.id
			WHERE pm.user_id = $1 AND p.status = 'active'`
		workersQuery = `
			SELECT COUNT(*) FROM attendance a
			JOIN projects p ON p.id = a.project_id
			JOIN project_members pm ON pm.project_id = p.id
			WHERE pm.user_id = $1 AND a.check_out IS NULL AND DATE(a.check_in) = CURRENT_DATE`
		issuesQuery = `
			SELECT COUNT(*) FROM issues i
			JOIN project_members pm ON pm.project_id = i.project_id
			WHERE pm.user_id = $1 AND i.status != 'resolved'`
		logsQuery = `
			SELECT COUNT(*) FROM daily_logs dl
			JOIN project_members pm ON pm.project_id = dl.project_id
			WHERE pm.user_id = $1 AND DATE(dl.created_at) = CURRENT_DATE`
		pendingQuery = `
			SELECT COUNT(*) FROM daily_logs dl
			JOIN project_members pm ON pm.project_id = dl.project_id
			WHERE pm.user_id = $1 AND dl.status = 'submitted'`

		_ = s.db.Get(&stats.ActiveProjects, projectQuery, userID)
		_ = s.db.Get(&stats.WorkersOnSite, workersQuery, userID)
		_ = s.db.Get(&stats.OpenIssues, issuesQuery, userID)
		_ = s.db.Get(&stats.LogsToday, logsQuery, userID)
		_ = s.db.Get(&stats.PendingApprovals, pendingQuery, userID)
		return &stats, nil
	}

	_ = s.db.Get(&stats.ActiveProjects, projectQuery)
	_ = s.db.Get(&stats.WorkersOnSite, workersQuery)
	_ = s.db.Get(&stats.OpenIssues, issuesQuery)
	_ = s.db.Get(&stats.LogsToday, logsQuery)
	_ = s.db.Get(&stats.PendingApprovals, pendingQuery)
	_ = s.db.Get(&stats.TotalMembers, membersQuery)

	return &stats, nil
}

func (s *Service) GetRecentActivity(role, userID string) (*RecentActivity, error) {
	var activity RecentActivity

	logsSQL := `
		SELECT dl.id, p.name AS project_name,
		       TO_CHAR(dl.date, 'YYYY-MM-DD') AS date,
		       dl.status, u.name AS created_by_name,
		       dl.created_at
		FROM daily_logs dl
		JOIN projects p ON p.id = dl.project_id
		JOIN users u ON u.id = dl.created_by`

	issuesSQL := `
		SELECT i.id, p.name AS project_name, i.title, i.priority, i.status
		FROM issues i
		JOIN projects p ON p.id = i.project_id`

	if role != "admin" {
		logsSQL += ` JOIN project_members pm ON pm.project_id = dl.project_id WHERE pm.user_id = $1 ORDER BY dl.created_at DESC LIMIT 5`
		issuesSQL += ` JOIN project_members pm ON pm.project_id = i.project_id WHERE pm.user_id = $1 ORDER BY i.created_at DESC LIMIT 5`
		_ = s.db.Select(&activity.Logs, logsSQL, userID)
		_ = s.db.Select(&activity.Issues, issuesSQL, userID)
	} else {
		logsSQL += ` ORDER BY dl.created_at DESC LIMIT 5`
		issuesSQL += ` ORDER BY i.created_at DESC LIMIT 5`
		_ = s.db.Select(&activity.Logs, logsSQL)
		_ = s.db.Select(&activity.Issues, issuesSQL)
	}

	if activity.Logs == nil {
		activity.Logs = []LogSummary{}
	}
	if activity.Issues == nil {
		activity.Issues = []IssueSummary{}
	}

	return &activity, nil
}
