package auth

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SaveRefreshToken(userID, tokenHash string, expiresAt time.Time) error {
	_, err := r.db.Exec(`
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)`,
		userID, tokenHash, expiresAt,
	)
	return err
}

func (r *Repository) FindRefreshToken(tokenHash string) (*RefreshToken, error) {
	var t RefreshToken
	err := r.db.Get(&t, `
		SELECT id, user_id, token_hash, expires_at, created_at
		FROM refresh_tokens
		WHERE token_hash = $1 AND expires_at > NOW()`,
		tokenHash,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *Repository) DeleteRefreshToken(tokenHash string) error {
	_, err := r.db.Exec(`DELETE FROM refresh_tokens WHERE token_hash = $1`, tokenHash)
	return err
}

func (r *Repository) DeleteUserRefreshTokens(userID string) error {
	_, err := r.db.Exec(`DELETE FROM refresh_tokens WHERE user_id = $1`, userID)
	return err
}
