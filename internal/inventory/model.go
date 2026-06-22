package inventory

import "time"

type Item struct {
	ID           string    `db:"id"            json:"id"`
	MaterialID   string    `db:"material_id"   json:"material_id"`
	LocationType string    `db:"location_type" json:"location_type"`
	ProjectID    *string   `db:"project_id"    json:"project_id"`
	Quantity     float64   `db:"quantity"      json:"quantity"`
	UnitCost     *float64  `db:"unit_cost"     json:"unit_cost"`
	CreatedAt    time.Time `db:"created_at"    json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"    json:"updated_at"`
	MaterialName string    `db:"material_name" json:"material_name"`
	MaterialUnit string    `db:"material_unit" json:"material_unit"`
	ProjectName  *string   `db:"project_name"  json:"project_name,omitempty"`
}

type Transaction struct {
	ID              string     `db:"id"               json:"id"`
	MaterialID      string     `db:"material_id"      json:"material_id"`
	LocationType    string     `db:"location_type"    json:"location_type"`
	ProjectID       *string    `db:"project_id"       json:"project_id"`
	QuantityChange  float64    `db:"quantity_change"  json:"quantity_change"`
	TransactionType string     `db:"transaction_type" json:"transaction_type"`
	ReferenceID     *string    `db:"reference_id"     json:"reference_id"`
	ReferenceType   *string    `db:"reference_type"   json:"reference_type"`
	Notes           *string    `db:"notes"            json:"notes"`
	CreatedBy       *string    `db:"created_by"       json:"created_by"`
	CreatedAt       time.Time  `db:"created_at"       json:"created_at"`
	Received        bool       `db:"received"         json:"received"`
	ReceivedBy      *string    `db:"received_by"      json:"received_by,omitempty"`
	ReceivedAt      *time.Time `db:"received_at"      json:"received_at,omitempty"`
	MaterialName    string     `db:"material_name"    json:"material_name"`
	MaterialUnit    string     `db:"material_unit"    json:"material_unit"`
	ProjectName     *string    `db:"project_name"     json:"project_name,omitempty"`
	CreatorName     *string    `db:"creator_name"     json:"creator_name,omitempty"`
}

type AddStockRequest struct {
	MaterialID string   `json:"material_id" binding:"required,uuid"`
	Quantity   float64  `json:"quantity"    binding:"required,gt=0"`
	UnitCost   *float64 `json:"unit_cost"`
	Notes      *string  `json:"notes"`
}

type TransferRequest struct {
	MaterialID string  `json:"material_id" binding:"required,uuid"`
	ProjectID  string  `json:"project_id"  binding:"required,uuid"`
	Quantity   float64 `json:"quantity"    binding:"required,gt=0"`
	Notes      *string `json:"notes"`
}
