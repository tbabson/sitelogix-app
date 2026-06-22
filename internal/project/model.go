package project

import "time"

type Status string

const (
	StatusActive    Status = "active"
	StatusOnHold    Status = "on_hold"
	StatusCompleted Status = "completed"
	StatusArchived  Status = "archived"
)

type Project struct {
	ID          string    `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description *string   `db:"description" json:"description"`
	Location    *string   `db:"location" json:"location"`
	StartDate   *time.Time `db:"start_date" json:"start_date"`
	EndDate     *time.Time `db:"end_date" json:"end_date"`
	Status      Status    `db:"status" json:"status"`
	CreatedBy   string    `db:"created_by" json:"created_by"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type Member struct {
	ID        string    `db:"id" json:"id"`
	ProjectID string    `db:"project_id" json:"project_id"`
	UserID    string    `db:"user_id" json:"user_id"`
	Role      string    `db:"role" json:"role"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type MemberUser struct {
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	Role      string  `json:"role"`
	AvatarURL *string `json:"avatar_url"`
}

type MemberDetail struct {
	ID        string      `json:"id"`
	ProjectID string      `json:"project_id"`
	UserID    string      `json:"user_id"`
	Role      string      `json:"role"`
	CreatedAt time.Time   `json:"created_at"`
	User      *MemberUser `json:"user,omitempty"`
}

type CreateProjectRequest struct {
	Name        string     `json:"name" binding:"required,min=2"`
	Description *string    `json:"description"`
	Location    *string    `json:"location"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
}

type UpdateProjectRequest struct {
	Name        string     `json:"name" binding:"omitempty,min=2"`
	Description *string    `json:"description"`
	Location    *string    `json:"location"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	Status      Status     `json:"status" binding:"omitempty,oneof=active on_hold completed archived"`
}

type AddMemberRequest struct {
	UserID string `json:"user_id" binding:"required,uuid"`
	Role   string `json:"role" binding:"required,oneof=admin supervisor client"`
}

type UpdateMemberRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=admin supervisor client"`
}

type DuplicateProjectRequest struct {
	Name string `json:"name" binding:"required,min=2"`
}

// ── Stats ─────────────────────────────────────────────────────────────────────

type ProjectStats struct {
	TotalLogs        int     `db:"total_logs" json:"total_logs"`
	SubmittedLogs    int     `db:"submitted_logs" json:"submitted_logs"`
	ApprovedLogs     int     `db:"approved_logs" json:"approved_logs"`
	OpenIssues       int     `db:"open_issues" json:"open_issues"`
	ResolvedIssues   int     `db:"resolved_issues" json:"resolved_issues"`
	TotalWorkers     int     `db:"total_workers" json:"total_workers"`
	WorkersOnSite    int     `db:"workers_on_site" json:"workers_on_site"`
	TotalMembers     int     `db:"total_members" json:"total_members"`
	TotalMedia       int     `db:"total_media" json:"total_media"`
	CompletionPct    float64 `json:"completion_pct"`
}

// ── Activity feed ─────────────────────────────────────────────────────────────

type ActivityItem struct {
	ID        string    `db:"id" json:"id"`
	Type      string    `db:"type" json:"type"`
	Title     string    `db:"title" json:"title"`
	Actor     string    `db:"actor" json:"actor"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// ── Milestones ────────────────────────────────────────────────────────────────

type MilestoneStatus string

const (
	MilestonePending    MilestoneStatus = "pending"
	MilestoneInProgress MilestoneStatus = "in_progress"
	MilestoneDone       MilestoneStatus = "done"
	MilestoneMissed     MilestoneStatus = "missed"
)

type Milestone struct {
	ID        string          `db:"id" json:"id"`
	ProjectID string          `db:"project_id" json:"project_id"`
	Title     string          `db:"title" json:"title"`
	DueDate   time.Time       `db:"due_date" json:"due_date"`
	Status    MilestoneStatus `db:"status" json:"status"`
	Notes     *string         `db:"notes" json:"notes"`
	CreatedAt time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt time.Time       `db:"updated_at" json:"updated_at"`
}

type CreateMilestoneRequest struct {
	Title   string          `json:"title" binding:"required,min=2"`
	DueDate time.Time       `json:"due_date" binding:"required"`
	Status  MilestoneStatus `json:"status" binding:"omitempty,oneof=pending in_progress done missed"`
	Notes   *string         `json:"notes"`
}

type UpdateMilestoneRequest struct {
	Title   string          `json:"title" binding:"omitempty,min=2"`
	DueDate *time.Time      `json:"due_date"`
	Status  MilestoneStatus `json:"status" binding:"omitempty,oneof=pending in_progress done missed"`
	Notes   *string         `json:"notes"`
}
