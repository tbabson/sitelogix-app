package requisition

import (
	"github.com/jmoiron/sqlx"
	"github.com/sitelogix/backend/internal/inventory"
)

type Service struct{ repo *Repository }

func NewService(db *sqlx.DB, invSvc *inventory.Service) *Service {
	return &Service{repo: NewRepository(db, invSvc.Repo())}
}

func (s *Service) Create(req CreateRequest, requestedBy string) (*Requisition, error) {
	return s.repo.Create(req, requestedBy)
}
func (s *Service) Get(id string) (*Requisition, error) { return s.repo.FindByID(id) }
func (s *Service) List(projectID, status string) ([]Requisition, error) {
	return s.repo.List(projectID, status)
}
func (s *Service) Approve(id string, req ApproveRequest) error { return s.repo.Approve(id, req) }
func (s *Service) Reject(id string, notes *string) error       { return s.repo.Reject(id, notes) }
func (s *Service) Fulfill(id, createdBy string) error          { return s.repo.Fulfill(id, createdBy) }
func (s *Service) Cancel(id, userID, role string) error        { return s.repo.Cancel(id, userID, role) }
