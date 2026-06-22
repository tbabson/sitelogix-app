package attendance

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

type Service struct {
	repo *Repository
}

func NewService(db *sqlx.DB) *Service {
	return &Service{repo: NewRepository(db)}
}

func (s *Service) CreateWorker(req CreateWorkerRequest) (*Worker, error) {
	return s.repo.CreateWorker(req)
}

func (s *Service) ListWorkers(projectID string) ([]Worker, error) {
	return s.repo.ListWorkers(projectID)
}

func (s *Service) CheckIn(req CheckInRequest, userID string) (*Attendance, error) {
	return s.repo.CheckIn(req, userID)
}

func (s *Service) CheckOut(id string, req CheckOutRequest) (*Attendance, error) {
	return s.repo.CheckOut(id, req)
}

func (s *Service) LinkUser(workerID, userID string) (*Worker, error) {
	return s.repo.LinkUser(workerID, userID)
}

func (s *Service) MyWorkerProfiles(userID string) ([]Worker, error) {
	return s.repo.FindByUserID(userID)
}

func (s *Service) SelfCheckIn(userID string, req SelfCheckInRequest) (*Attendance, error) {
	worker, err := s.repo.FindByUserAndProject(userID, req.ProjectID)
	if err != nil {
		return nil, errors.New("no worker record found for this project")
	}
	return s.repo.CheckIn(CheckInRequest{
		ProjectID: req.ProjectID,
		WorkerID:  worker.ID,
		CheckIn:   time.Now(),
		Notes:     req.Notes,
	}, userID)
}

func (s *Service) SelfCheckOut(userID string) (*Attendance, error) {
	workers, err := s.repo.FindByUserID(userID)
	if err != nil || len(workers) == 0 {
		return nil, errors.New("no worker record found")
	}
	for _, w := range workers {
		a, err := s.repo.GetOpenAttendance(w.ID)
		if err == nil {
			return s.repo.CheckOut(a.ID, CheckOutRequest{CheckOut: time.Now()})
		}
	}
	return nil, errors.New("no open check-in found")
}

func (s *Service) List(f AttendanceFilter) ([]AttendanceWithWorker, int, error) {
	if f.Limit <= 0 {
		f.Limit = 50
	}
	if f.Page <= 0 {
		f.Page = 1
	}
	return s.repo.List(f)
}
