package notification

import "github.com/jmoiron/sqlx"

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(n *Notification) (*Notification, error) {
	var saved Notification
	err := r.db.QueryRowx(`
		INSERT INTO notifications (user_id, title, body, type, reference_id, reference_type)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, title, body, type, reference_id, reference_type, is_read, created_at`,
		n.UserID, n.Title, n.Body, n.Type, n.ReferenceID, n.ReferenceType,
	).StructScan(&saved)
	return &saved, err
}

func (r *Repository) ListForUser(userID string, onlyUnread bool, limit, offset int) ([]Notification, int, error) {
	where := "user_id = $1"
	if onlyUnread {
		where += " AND is_read = false"
	}
	var total int
	_ = r.db.Get(&total, `SELECT COUNT(*) FROM notifications WHERE `+where, userID)

	var list []Notification
	err := r.db.Select(&list,
		`SELECT id, user_id, title, body, type, reference_id, reference_type, is_read, created_at
		 FROM notifications WHERE `+where+` ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		userID, limit, offset,
	)
	return list, total, err
}

func (r *Repository) MarkRead(id, userID string) error {
	_, err := r.db.Exec(`UPDATE notifications SET is_read = true WHERE id = $1 AND user_id = $2`, id, userID)
	return err
}

func (r *Repository) MarkAllRead(userID string) error {
	_, err := r.db.Exec(`UPDATE notifications SET is_read = true WHERE user_id = $1`, userID)
	return err
}

func (r *Repository) UnreadCount(userID string) (int, error) {
	var count int
	err := r.db.Get(&count, `SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = false`, userID)
	return count, err
}

func (r *Repository) SaveDeviceToken(userID, token, platform string) error {
	_, err := r.db.Exec(`
		INSERT INTO device_tokens (user_id, token, platform)
		VALUES ($1, $2, $3)
		ON CONFLICT (token) DO UPDATE SET user_id = EXCLUDED.user_id, platform = EXCLUDED.platform`,
		userID, token, platform,
	)
	return err
}

func (r *Repository) GetUserTokens(userID string) ([]string, error) {
	var tokens []string
	err := r.db.Select(&tokens, `SELECT token FROM device_tokens WHERE user_id = $1`, userID)
	return tokens, err
}

func (r *Repository) DeleteDeviceToken(token string) error {
	_, err := r.db.Exec(`DELETE FROM device_tokens WHERE token = $1`, token)
	return err
}
