package dailylog

import "time"

type LogStatus string

const (
	StatusDraft    LogStatus = "draft"
	StatusSubmitted LogStatus = "submitted"
	StatusApproved LogStatus = "approved"
	StatusRejected LogStatus = "rejected"
)

type DailyLog struct {
	ID           string        `db:"id" json:"id"`
	ProjectID    string        `db:"project_id" json:"project_id"`
	CreatedBy    string        `db:"created_by" json:"created_by"`
	Date         time.Time     `db:"date" json:"date"`
	Weather      *string       `db:"weather" json:"weather"`
	Notes        *string       `db:"notes" json:"notes"`
	GPSLat       *float64      `db:"gps_lat" json:"gps_lat"`
	GPSLng       *float64      `db:"gps_lng" json:"gps_lng"`
	SignatureURL *string       `db:"signature_url" json:"signature_url"`
	Status       LogStatus     `db:"status" json:"status"`
	CreatedAt    time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time     `db:"updated_at" json:"updated_at"`
	ProjectName  string        `db:"project_name" json:"project_name"`
	CreatorName  string        `db:"creator_name" json:"creator_name"`
	Materials    []LogMaterial `db:"-" json:"materials"`
}

type CreateLogRequest struct {
	ProjectID string                  `json:"project_id" binding:"required,uuid"`
	Date      time.Time               `json:"date" binding:"required"`
	Weather   *string                 `json:"weather"`
	Notes     *string                 `json:"notes"`
	GPSLat    *float64                `json:"gps_lat"`
	GPSLng    *float64                `json:"gps_lng"`
	Materials []AddLogMaterialRequest `json:"materials"`
}

type UpdateLogRequest struct {
	Weather      *string  `json:"weather"`
	Notes        *string  `json:"notes"`
	GPSLat       *float64 `json:"gps_lat"`
	GPSLng       *float64 `json:"gps_lng"`
	SignatureURL *string  `json:"signature_url"`
}

type LogFilter struct {
	ProjectID string
	Username  string
	Status    string
	DateFrom  string
	DateTo    string
	Page      int
	Limit     int
}

type LogMaterial struct {
	ID           string    `db:"id"            json:"id"`
	LogID        string    `db:"log_id"        json:"log_id"`
	MaterialID   string    `db:"material_id"   json:"material_id"`
	QuantityUsed float64   `db:"quantity_used" json:"quantity_used"`
	CreatedAt    time.Time `db:"created_at"    json:"created_at"`
	MaterialName string    `db:"material_name" json:"material_name"`
	MaterialUnit string    `db:"material_unit" json:"material_unit"`
}

type AddLogMaterialRequest struct {
	MaterialID   string  `json:"material_id"  binding:"required,uuid"`
	QuantityUsed float64 `json:"quantity_used" binding:"required,gt=0"`
}
