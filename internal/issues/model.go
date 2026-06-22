package issues

import "time"

type Priority string
type IssueStatus string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"

	IssueOpen       IssueStatus = "open"
	IssueInProgress IssueStatus = "in_progress"
	IssueResolved   IssueStatus = "resolved"
)

type Issue struct {
	ID          string      `db:"id" json:"id"`
	ProjectID   string      `db:"project_id" json:"project_id"`
	Title       string      `db:"title" json:"title"`
	Description *string     `db:"description" json:"description"`
	Priority    Priority    `db:"priority" json:"priority"`
	Status      IssueStatus `db:"status" json:"status"`
	AssignedTo  *string     `db:"assigned_to" json:"assigned_to"`
	ReportedBy  string      `db:"reported_by" json:"reported_by"`
	GPSLat      *float64    `db:"gps_lat" json:"gps_lat"`
	GPSLng      *float64    `db:"gps_lng" json:"gps_lng"`
	CreatedAt   time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time   `db:"updated_at" json:"updated_at"`
}

type CreateIssueRequest struct {
	ProjectID   string   `json:"project_id" binding:"required,uuid"`
	Title       string   `json:"title" binding:"required,min=3"`
	Description *string  `json:"description"`
	Priority    Priority `json:"priority" binding:"required,oneof=low medium high"`
	AssignedTo  *string  `json:"assigned_to"`
	GPSLat      *float64 `json:"gps_lat"`
	GPSLng      *float64 `json:"gps_lng"`
}

type UpdateIssueRequest struct {
	Title       string      `json:"title" binding:"omitempty,min=3"`
	Description *string     `json:"description"`
	Priority    Priority    `json:"priority" binding:"omitempty,oneof=low medium high"`
	Status      IssueStatus `json:"status" binding:"omitempty,oneof=open in_progress resolved"`
	AssignedTo  *string     `json:"assigned_to"`
	GPSLat      *float64    `json:"gps_lat"`
	GPSLng      *float64    `json:"gps_lng"`
}

type IssueFilter struct {
	ProjectID  string
	AssignedTo string
	Status     string
	Priority   string
	Page       int
	Limit      int
}
