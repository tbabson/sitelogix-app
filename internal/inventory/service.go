package inventory

import "github.com/jmoiron/sqlx"

type Service struct{ repo *Repository }

func NewService(db *sqlx.DB) *Service { return &Service{repo: NewRepository(db)} }

func (s *Service) GetHeadOffice() ([]Item, error)             { return s.repo.GetHeadOffice() }
func (s *Service) GetProject(projectID string) ([]Item, error) { return s.repo.GetProject(projectID) }
func (s *Service) AddStock(req AddStockRequest, createdBy string) error {
	return s.repo.AddStock(req, createdBy)
}
func (s *Service) Transfer(req TransferRequest, createdBy string) error {
	return s.repo.Transfer(req, createdBy)
}
func (s *Service) GetTransactions(materialID, projectID string, limit int) ([]Transaction, error) {
	return s.repo.GetTransactions(materialID, projectID, limit)
}
func (s *Service) GetPendingTransfers(projectID string) ([]Transaction, error) {
	return s.repo.GetPendingTransfers(projectID)
}
func (s *Service) ConfirmReceipt(transactionID, receivedBy string) error {
	return s.repo.ConfirmReceipt(transactionID, receivedBy)
}

// Repo exposes the repository so requisition package can share transactions
func (s *Service) Repo() *Repository { return s.repo }
