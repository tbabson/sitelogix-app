package materials

import "github.com/jmoiron/sqlx"

type Repository struct{ db *sqlx.DB }

func NewRepository(db *sqlx.DB) *Repository { return &Repository{db: db} }

func (r *Repository) Create(req CreateMaterialRequest) (*Material, error) {
	var m Material
	err := r.db.QueryRowx(`
		INSERT INTO materials (name, unit, description)
		VALUES ($1, $2, $3)
		RETURNING id, name, unit, description, created_at, updated_at`,
		req.Name, req.Unit, req.Description,
	).StructScan(&m)
	return &m, err
}

func (r *Repository) List() ([]Material, error) {
	var ms []Material
	err := r.db.Select(&ms,
		`SELECT id, name, unit, description, created_at, updated_at FROM materials ORDER BY name`)
	return ms, err
}

func (r *Repository) FindByID(id string) (*Material, error) {
	var m Material
	err := r.db.Get(&m,
		`SELECT id, name, unit, description, created_at, updated_at FROM materials WHERE id = $1`, id)
	return &m, err
}

func (r *Repository) Update(id string, req UpdateMaterialRequest) (*Material, error) {
	var m Material
	err := r.db.QueryRowx(`
		UPDATE materials SET
			name        = CASE WHEN $1::text <> '' THEN $1::text ELSE name END,
			unit        = CASE WHEN $2::text <> '' THEN $2::text ELSE unit END,
			description = COALESCE($3::text, description),
			updated_at  = NOW()
		WHERE id = $4::uuid
		RETURNING id, name, unit, description, created_at, updated_at`,
		req.Name, req.Unit, req.Description, id,
	).StructScan(&m)
	return &m, err
}

func (r *Repository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM materials WHERE id = $1`, id)
	return err
}
