package requisition

import "time"

type Requisition struct {
	ID            string    `db:"id"             json:"id"`
	ProjectID     string    `db:"project_id"     json:"project_id"`
	RequestedBy   string    `db:"requested_by"   json:"requested_by"`
	Status        string    `db:"status"         json:"status"`
	Notes         *string   `db:"notes"          json:"notes"`
	CreatedAt     time.Time `db:"created_at"     json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"     json:"updated_at"`
	ProjectName   string    `db:"project_name"   json:"project_name"`
	RequesterName string    `db:"requester_name" json:"requester_name"`
	Items         []Item    `db:"-"              json:"items,omitempty"`
}

type Item struct {
	ID                string    `db:"id"                 json:"id"`
	RequisitionID     string    `db:"requisition_id"     json:"requisition_id"`
	MaterialID        string    `db:"material_id"        json:"material_id"`
	QuantityRequested float64   `db:"quantity_requested" json:"quantity_requested"`
	QuantityApproved  *float64  `db:"quantity_approved"  json:"quantity_approved"`
	CreatedAt         time.Time `db:"created_at"         json:"created_at"`
	MaterialName      string    `db:"material_name"      json:"material_name"`
	MaterialUnit      string    `db:"material_unit"      json:"material_unit"`
}

type CreateRequest struct {
	ProjectID string      `json:"project_id" binding:"required,uuid"`
	Notes     *string     `json:"notes"`
	Items     []ItemInput `json:"items"      binding:"required,min=1,dive"`
}

type ItemInput struct {
	MaterialID        string  `json:"material_id"        binding:"required,uuid"`
	QuantityRequested float64 `json:"quantity_requested" binding:"required,gt=0"`
}

type ApproveRequest struct {
	Notes *string       `json:"notes"`
	Items []ApproveItem `json:"items" binding:"required,min=1,dive"`
}

type ApproveItem struct {
	ItemID           string  `json:"item_id"           binding:"required,uuid"`
	QuantityApproved float64 `json:"quantity_approved"  binding:"required,gte=0"`
}

type RejectRequest struct {
	Notes *string `json:"notes"`
}
