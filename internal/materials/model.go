package materials

import "time"

type Material struct {
	ID          string    `db:"id"          json:"id"`
	Name        string    `db:"name"        json:"name"`
	Unit        string    `db:"unit"        json:"unit"`
	Description *string   `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at"  json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"  json:"updated_at"`
}

type CreateMaterialRequest struct {
	Name        string  `json:"name"        binding:"required,min=2"`
	Unit        string  `json:"unit"        binding:"required,min=1"`
	Description *string `json:"description"`
}

type UpdateMaterialRequest struct {
	Name        string  `json:"name"`
	Unit        string  `json:"unit"`
	Description *string `json:"description"`
}
