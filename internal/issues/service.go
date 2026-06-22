package issues

import "github.com/jmoiron/sqlx"

type Service struct {
	repo *Repository
}

func NewService(db *sqlx.DB) *Service {
	return &Service{repo: NewRepository(db)}
}

func (s *Service) Create(req CreateIssueRequest, userID string) (*Issue, error) {
	return s.repo.Create(req, userID)
}

func (s *Service) Get(id string) (*Issue, error) {
	return s.repo.FindByID(id)
}

func (s *Service) List(f IssueFilter) ([]Issue, int, error) {
	if f.Limit <= 0 {
		f.Limit = 20
	}
	if f.Page <= 0 {
		f.Page = 1
	}
	return s.repo.List(f)
}

func (s *Service) Update(id string, req UpdateIssueRequest) (*Issue, error) {
	return s.repo.Update(id, req)
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}
