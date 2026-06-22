package media

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sitelogix/backend/pkg/storage"
)

var ErrUnsupportedType = errors.New("unsupported file type")

type Service struct {
	repo    *Repository
	storage *storage.GitHubStorage
}

func NewService(db *sqlx.DB, gh *storage.GitHubStorage) *Service {
	return &Service{repo: NewRepository(db), storage: gh}
}

func (s *Service) Upload(fileData []byte, filename, contentType, projectID, uploaderID string, logID *string, issueID *string) (*Media, error) {
	folder, ok := AllowedTypes[contentType]
	if !ok {
		return nil, ErrUnsupportedType
	}

	if int64(len(fileData)) > MaxUploadSize {
		return nil, fmt.Errorf("file exceeds maximum size of 100MB")
	}

	result, err := s.storage.Upload(fileData, filename, fmt.Sprintf("%s/%s", projectID, folder))
	if err != nil {
		return nil, fmt.Errorf("upload to github: %w", err)
	}

	m := &Media{
		ProjectID:   projectID,
		UploadedBy:  uploaderID,
		Filename:    filename,
		ContentType: contentType,
		URL:         result.URL,
		GitHubSHA:   result.SHA,
		GitHubKey:   result.Key,
		Size:        int64(len(fileData)),
		LogID:       logID,
		IssueID:     issueID,
	}

	return s.repo.Save(m)
}

func (s *Service) List(f MediaFilter) ([]Media, int, error) {
	if f.Limit <= 0 {
		f.Limit = 20
	}
	if f.Page <= 0 {
		f.Page = 1
	}
	return s.repo.List(f)
}

func (s *Service) Get(id string) (*Media, error) {
	return s.repo.FindByID(id)
}

func (s *Service) Delete(id string) error {
	m, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	_ = s.storage.Delete(m.GitHubKey, m.GitHubSHA)
	return s.repo.Delete(id)
}

func (s *Service) Proxy(id string) ([]byte, string, error) {
	m, err := s.repo.FindByID(id)
	if err != nil {
		return nil, "", err
	}
	data, err := s.storage.Download(m.GitHubKey)
	if err != nil {
		return nil, "", err
	}
	return data, m.ContentType, nil
}

func (s *Service) AttachToLog(mediaID, logID string) error {
	return s.repo.AttachToLog(mediaID, logID)
}

func (s *Service) AttachToIssue(mediaID, issueID string) error {
	return s.repo.AttachToIssue(mediaID, issueID)
}
