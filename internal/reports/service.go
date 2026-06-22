package reports

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type Service struct {
	repo *Repository
}

func NewService(db *sqlx.DB) *Service {
	return &Service{repo: NewRepository(db)}
}

func (s *Service) Generate(req ReportRequest) ([]byte, string, error) {
	project, err := s.repo.GetProject(req.ProjectID)
	if err != nil {
		return nil, "", err
	}

	logs, err := s.repo.GetLogs(req.ProjectID, req.DateFrom, req.DateTo)
	if err != nil {
		return nil, "", err
	}

	issues, err := s.repo.GetIssues(req.ProjectID)
	if err != nil {
		return nil, "", err
	}

	attendance, err := s.repo.GetAttendance(req.ProjectID, req.DateFrom, req.DateTo)
	if err != nil {
		return nil, "", err
	}

	dateFrom := ""
	dateTo := ""
	if req.DateFrom != nil {
		dateFrom = req.DateFrom.Format("02 Jan 2006")
	}
	if req.DateTo != nil {
		dateTo = req.DateTo.Format("02 Jan 2006")
	}

	data := &ReportData{
		Project:     *project,
		Logs:        logs,
		Issues:      issues,
		Attendance:  attendance,
		GeneratedAt: time.Now(),
		DateFrom:    dateFrom,
		DateTo:      dateTo,
	}

	pdf, err := GeneratePDF(data)
	if err != nil {
		return nil, "", err
	}

	filename := "sitelogix-report-" + project.Name + "-" + time.Now().Format("20060102") + ".pdf"
	return pdf, filename, nil
}
