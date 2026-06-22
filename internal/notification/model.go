package notification

import "time"

type NotifType string

const (
	TypeLogSubmitted    NotifType = "log_submitted"
	TypeLogApproved     NotifType = "log_approved"
	TypeLogRejected     NotifType = "log_rejected"
	TypeIssueAssigned   NotifType = "issue_assigned"
	TypeIssueUpdated    NotifType = "issue_updated"
	TypeApprovalNeeded  NotifType = "approval_needed"
)

type Notification struct {
	ID            string    `db:"id" json:"id"`
	UserID        string    `db:"user_id" json:"user_id"`
	Title         string    `db:"title" json:"title"`
	Body          string    `db:"body" json:"body"`
	Type          NotifType `db:"type" json:"type"`
	ReferenceID   *string   `db:"reference_id" json:"reference_id"`
	ReferenceType *string   `db:"reference_type" json:"reference_type"`
	IsRead        bool      `db:"is_read" json:"is_read"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}

type DeviceToken struct {
	ID        string    `db:"id" json:"id"`
	UserID    string    `db:"user_id" json:"user_id"`
	Token     string    `db:"token" json:"token"`
	Platform  string    `db:"platform" json:"platform"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type RegisterTokenRequest struct {
	Token    string `json:"token" binding:"required"`
	Platform string `json:"platform" binding:"required,oneof=android ios web"`
}

type SendRequest struct {
	UserID        string    `json:"user_id"`
	Title         string    `json:"title"`
	Body          string    `json:"body"`
	Type          NotifType `json:"type"`
	ReferenceID   *string   `json:"reference_id"`
	ReferenceType *string   `json:"reference_type"`
}
