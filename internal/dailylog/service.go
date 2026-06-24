package dailylog

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/sitelogix/backend/internal/inventory"
	"github.com/sitelogix/backend/pkg/weather"
)

var ErrForbidden = errors.New("not authorized to modify this log")

type Service struct {
	repo    *Repository
	weather *weather.Client
}

func NewService(db *sqlx.DB, weatherClient *weather.Client, invSvc *inventory.Service) *Service {
	return &Service{repo: NewRepository(db, invSvc.Repo()), weather: weatherClient}
}

func (s *Service) Create(req CreateLogRequest, userID string) (*DailyLog, error) {
	if req.GPSLat != nil && req.GPSLng != nil && req.Weather == nil {
		if w, err := s.weather.Fetch(*req.GPSLat, *req.GPSLng); err == nil && w != "" {
			req.Weather = &w
		}
	}
	return s.repo.Create(req, userID)
}

func (s *Service) Get(id string) (*DailyLog, error) {
	log, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	mats, _ := s.repo.ListLogMaterials(id)
	log.Materials = mats
	return log, nil
}

func (s *Service) List(f LogFilter) ([]DailyLog, int, error) {
	if f.Limit <= 0 {
		f.Limit = 20
	}
	if f.Page <= 0 {
		f.Page = 1
	}
	return s.repo.List(f)
}

func (s *Service) Update(id, userID, role string, req UpdateLogRequest) (*DailyLog, error) {
	log, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if log.CreatedBy != userID && role != "admin" {
		return nil, ErrForbidden
	}
	return s.repo.Update(id, req)
}

func (s *Service) Submit(id, userID string) (*DailyLog, error) {
	log, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if log.CreatedBy != userID {
		return nil, ErrForbidden
	}
	return s.repo.UpdateStatus(id, StatusSubmitted)
}

func (s *Service) Approve(id string) (*DailyLog, error) {
	return s.repo.UpdateStatus(id, StatusApproved)
}

func (s *Service) Reject(id string) (*DailyLog, error) {
	return s.repo.UpdateStatus(id, StatusRejected)
}

func (s *Service) Delete(id, userID, role string) error {
	log, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if log.CreatedBy != userID && role != "admin" {
		return ErrForbidden
	}
	return s.repo.Delete(id, userID)
}

func (s *Service) ListLogMaterials(logID string) ([]LogMaterial, error) {
	return s.repo.ListLogMaterials(logID)
}

func (s *Service) GetWeather(lat, lng float64) (string, error) {
	return s.weather.Fetch(lat, lng)
}

func (s *Service) AddLogMaterial(logID, userID, role string, req AddLogMaterialRequest) (*LogMaterial, error) {
	log, err := s.repo.FindByID(logID)
	if err != nil {
		return nil, err
	}
	if log.CreatedBy != userID && role != "admin" {
		return nil, ErrForbidden
	}
	return s.repo.AddLogMaterial(logID, userID, req)
}

func (s *Service) RemoveLogMaterial(logID, materialID, userID, role string) error {
	log, err := s.repo.FindByID(logID)
	if err != nil {
		return err
	}
	if log.CreatedBy != userID && role != "admin" {
		return ErrForbidden
	}
	return s.repo.RemoveLogMaterial(logID, materialID, userID)
}
