package user

import "github.com/jmoiron/sqlx"

type Service struct {
	repo *Repository
}

func NewService(db *sqlx.DB) *Service {
	return &Service{repo: NewRepository(db)}
}

func (s *Service) GetProfile(userID string) (*PublicUser, error) {
	u, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	v := u.PublicView()
	return &v, nil
}

func (s *Service) UpdateProfile(userID string, p UpdateUserParams) (*PublicUser, error) {
	u, err := s.repo.Update(userID, p)
	if err != nil {
		return nil, err
	}
	v := u.PublicView()
	return &v, nil
}

func (s *Service) ListUsers(page, limit int) ([]PublicUser, int, error) {
	if limit <= 0 {
		limit = 20
	}
	offset := (page - 1) * limit
	users, total, err := s.repo.List(limit, offset)
	if err != nil {
		return nil, 0, err
	}
	pub := make([]PublicUser, len(users))
	for i, u := range users {
		pub[i] = u.PublicView()
	}
	return pub, total, nil
}

func (s *Service) DeactivateUser(userID string) error {
	return s.repo.SetActive(userID, false)
}
