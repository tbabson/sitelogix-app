package project

import "github.com/jmoiron/sqlx"

type Service struct {
	repo *Repository
}

func NewService(db *sqlx.DB) *Service {
	return &Service{repo: NewRepository(db)}
}

func (s *Service) Create(req CreateProjectRequest, createdBy string) (*Project, error) {
	return s.repo.Create(req, createdBy)
}

func (s *Service) Get(id string) (*Project, error) {
	return s.repo.FindByID(id)
}

func (s *Service) List(role, userID string, page, limit int) ([]Project, int, error) {
	if limit <= 0 {
		limit = 20
	}
	offset := (page - 1) * limit
	if role == "admin" {
		return s.repo.List(limit, offset)
	}
	return s.repo.ListForUser(userID, limit, offset)
}

func (s *Service) Update(id string, req UpdateProjectRequest) (*Project, error) {
	return s.repo.Update(id, req)
}

func (s *Service) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *Service) AddMember(projectID string, req AddMemberRequest) (*Member, error) {
	return s.repo.AddMember(projectID, req.UserID, req.Role)
}

func (s *Service) RemoveMember(projectID, userID string) error {
	return s.repo.RemoveMember(projectID, userID)
}

func (s *Service) ListMembers(projectID string) ([]MemberDetail, error) {
	return s.repo.ListMembers(projectID)
}

func (s *Service) UpdateMemberRole(projectID, userID string, req UpdateMemberRoleRequest) error {
	return s.repo.UpdateMemberRole(projectID, userID, req.Role)
}

func (s *Service) Archive(id string) (*Project, error) {
	return s.repo.SetStatus(id, StatusArchived)
}

func (s *Service) Restore(id string) (*Project, error) {
	return s.repo.SetStatus(id, StatusActive)
}

func (s *Service) Duplicate(id string, req DuplicateProjectRequest, createdBy string) (*Project, error) {
	return s.repo.Duplicate(id, req.Name, createdBy)
}

func (s *Service) GetStats(projectID string) (*ProjectStats, error) {
	return s.repo.GetStats(projectID)
}

func (s *Service) GetActivity(projectID string, limit int) ([]ActivityItem, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	return s.repo.GetActivity(projectID, limit)
}

func (s *Service) CreateMilestone(projectID string, req CreateMilestoneRequest) (*Milestone, error) {
	return s.repo.CreateMilestone(projectID, req)
}

func (s *Service) ListMilestones(projectID string) ([]Milestone, error) {
	return s.repo.ListMilestones(projectID)
}

func (s *Service) UpdateMilestone(id string, req UpdateMilestoneRequest) (*Milestone, error) {
	return s.repo.UpdateMilestone(id, req)
}

func (s *Service) DeleteMilestone(id string) error {
	return s.repo.DeleteMilestone(id)
}
