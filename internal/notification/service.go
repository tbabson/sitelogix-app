package notification

import (
	"github.com/jmoiron/sqlx"
	"github.com/sitelogix/backend/pkg/fcm"
)

type Service struct {
	repo *Repository
	fcm  *fcm.Client
}

func NewService(db *sqlx.DB, fcmClient *fcm.Client) *Service {
	return &Service{
		repo: NewRepository(db),
		fcm:  fcmClient,
	}
}

// Notify creates a DB notification and attempts FCM push delivery.
func (s *Service) Notify(req SendRequest) {
	n := &Notification{
		UserID:        req.UserID,
		Title:         req.Title,
		Body:          req.Body,
		Type:          req.Type,
		ReferenceID:   req.ReferenceID,
		ReferenceType: req.ReferenceType,
	}

	saved, err := s.repo.Create(n)
	if err != nil || saved == nil {
		return
	}

	tokens, err := s.repo.GetUserTokens(req.UserID)
	if err != nil || len(tokens) == 0 {
		return
	}

	data := map[string]string{
		"notification_id": saved.ID,
		"type":            string(req.Type),
	}
	if req.ReferenceID != nil {
		data["reference_id"] = *req.ReferenceID
	}

	// Fire-and-forget: don't block the caller on push failures
	go func() {
		s.fcm.SendMulticast(tokens, req.Title, req.Body, data)
	}()
}

func (s *Service) List(userID string, onlyUnread bool, page, limit int) ([]Notification, int, error) {
	if limit <= 0 {
		limit = 20
	}
	offset := (page - 1) * limit
	return s.repo.ListForUser(userID, onlyUnread, limit, offset)
}

func (s *Service) MarkRead(id, userID string) error {
	return s.repo.MarkRead(id, userID)
}

func (s *Service) MarkAllRead(userID string) error {
	return s.repo.MarkAllRead(userID)
}

func (s *Service) UnreadCount(userID string) (int, error) {
	return s.repo.UnreadCount(userID)
}

func (s *Service) RegisterToken(userID string, req RegisterTokenRequest) error {
	return s.repo.SaveDeviceToken(userID, req.Token, req.Platform)
}

func (s *Service) RemoveToken(token string) error {
	return s.repo.DeleteDeviceToken(token)
}
