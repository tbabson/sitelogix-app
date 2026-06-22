package materials

import "github.com/jmoiron/sqlx"

type Service struct{ repo *Repository }

func NewService(db *sqlx.DB) *Service { return &Service{repo: NewRepository(db)} }

func (s *Service) Create(req CreateMaterialRequest) (*Material, error) { return s.repo.Create(req) }
func (s *Service) List() ([]Material, error)                           { return s.repo.List() }
func (s *Service) Get(id string) (*Material, error)                    { return s.repo.FindByID(id) }
func (s *Service) Update(id string, req UpdateMaterialRequest) (*Material, error) {
	return s.repo.Update(id, req)
}
func (s *Service) Delete(id string) error { return s.repo.Delete(id) }
