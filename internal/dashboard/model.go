package dashboard

type Stats struct {
	ActiveProjects   int `json:"active_projects"`
	WorkersOnSite    int `json:"workers_on_site"`
	OpenIssues       int `json:"open_issues"`
	LogsToday        int `json:"logs_today"`
	PendingApprovals int `json:"pending_approvals"`
	TotalMembers     int `json:"total_members"`
}

type RecentActivity struct {
	Logs   []LogSummary   `json:"recent_logs"`
	Issues []IssueSummary `json:"recent_issues"`
}

type LogSummary struct {
	ID          string `db:"id" json:"id"`
	ProjectName string `db:"project_name" json:"project_name"`
	Date        string `db:"date" json:"date"`
	Status      string `db:"status" json:"status"`
	CreatedBy   string `db:"created_by_name" json:"created_by_name"`
	CreatedAt   string `db:"created_at" json:"created_at"`
}

type IssueSummary struct {
	ID          string `db:"id" json:"id"`
	ProjectName string `db:"project_name" json:"project_name"`
	Title       string `db:"title" json:"title"`
	Priority    string `db:"priority" json:"priority"`
	Status      string `db:"status" json:"status"`
}

type DashboardResponse struct {
	Stats          Stats          `json:"stats"`
	RecentActivity RecentActivity `json:"recent_activity"`
}
