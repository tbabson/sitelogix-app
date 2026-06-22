package user

import "time"

type User struct {
	ID           string    `db:"id" json:"id"`
	Name         string    `db:"name" json:"name"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Role         string    `db:"role" json:"role"`
	AvatarURL    *string   `db:"avatar_url" json:"avatar_url"`
	IsActive     bool      `db:"is_active" json:"is_active"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type PublicUser struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	AvatarURL *string   `json:"avatar_url"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *User) PublicView() PublicUser {
	return PublicUser{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role,
		AvatarURL: u.AvatarURL,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
	}
}

type CreateUserParams struct {
	Name         string
	Email        string
	PasswordHash string
	Role         string
}

type UpdateUserParams struct {
	Name      string  `json:"name" binding:"omitempty,min=2"`
	AvatarURL *string `json:"avatar_url"`
}
