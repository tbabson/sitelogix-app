package user

import (
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(p CreateUserParams) (*User, error) {
	var u User
	err := r.db.QueryRowx(`
		INSERT INTO users (name, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, email, password_hash, role, avatar_url, is_active, created_at, updated_at`,
		p.Name, p.Email, p.PasswordHash, p.Role,
	).StructScan(&u)
	return &u, err
}

func (r *Repository) FindByEmail(email string) (*User, error) {
	var u User
	err := r.db.Get(&u,
		`SELECT id, name, email, password_hash, role, avatar_url, is_active, created_at, updated_at
		 FROM users WHERE email = $1 AND is_active = true`, email)
	return &u, err
}

func (r *Repository) FindByID(id string) (*User, error) {
	var u User
	err := r.db.Get(&u,
		`SELECT id, name, email, password_hash, role, avatar_url, is_active, created_at, updated_at
		 FROM users WHERE id = $1`, id)
	return &u, err
}

func (r *Repository) List(limit, offset int) ([]User, int, error) {
	var users []User
	err := r.db.Select(&users,
		`SELECT id, name, email, password_hash, role, avatar_url, is_active, created_at, updated_at
		 FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	var total int
	_ = r.db.Get(&total, `SELECT COUNT(*) FROM users`)
	return users, total, nil
}

func (r *Repository) Update(id string, p UpdateUserParams) (*User, error) {
	var u User
	err := r.db.QueryRowx(`
		UPDATE users SET name = COALESCE(NULLIF($1,''), name),
		                 avatar_url = COALESCE($2, avatar_url),
		                 updated_at = NOW()
		WHERE id = $3
		RETURNING id, name, email, password_hash, role, avatar_url, is_active, created_at, updated_at`,
		p.Name, p.AvatarURL, id,
	).StructScan(&u)
	return &u, err
}

func (r *Repository) SetActive(id string, active bool) error {
	_, err := r.db.Exec(`UPDATE users SET is_active = $1, updated_at = NOW() WHERE id = $2`, active, id)
	return err
}
