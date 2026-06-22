package search

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type Service struct {
	repo *Repository
}

func NewService(db *sqlx.DB) *Service {
	return &Service{repo: NewRepository(db)}
}

func (s *Service) Search(query, userID, role string) (*SearchResponse, error) {
	query = strings.TrimSpace(query)
	if len(query) < 2 {
		return &SearchResponse{Query: query, Results: []Result{}}, nil
	}

	results := []Result{}

	projects, _ := s.repo.SearchProjects(query, userID, role)
	for _, p := range projects {
		excerpt := ""
		if p.Description != nil {
			excerpt = truncate(*p.Description, 100)
		}
		results = append(results, Result{
			Type:    TypeProject,
			ID:      p.ID,
			Title:   p.Name,
			Excerpt: excerpt,
			Meta:    map[string]string{"status": p.Status},
		})
	}

	logs, _ := s.repo.SearchLogs(query, userID, role)
	for _, l := range logs {
		excerpt := ""
		if l.Notes != nil {
			excerpt = truncate(*l.Notes, 100)
		}
		results = append(results, Result{
			Type:    TypeLog,
			ID:      l.ID,
			Title:   l.ProjectName + " — " + l.Date,
			Excerpt: excerpt,
			Meta:    map[string]string{"status": l.Status, "project": l.ProjectName},
		})
	}

	issues, _ := s.repo.SearchIssues(query, userID, role)
	for _, i := range issues {
		excerpt := ""
		if i.Description != nil {
			excerpt = truncate(*i.Description, 100)
		}
		results = append(results, Result{
			Type:    TypeIssue,
			ID:      i.ID,
			Title:   i.Title,
			Excerpt: excerpt,
			Meta:    map[string]string{"priority": i.Priority, "status": i.Status, "project": i.ProjectName},
		})
	}

	workers, _ := s.repo.SearchWorkers(query, userID, role)
	for _, w := range workers {
		trade := ""
		if w.Trade != nil {
			trade = *w.Trade
		}
		results = append(results, Result{
			Type:    TypeWorker,
			ID:      w.ID,
			Title:   w.Name,
			Excerpt: trade,
			Meta:    map[string]string{"project": w.ProjectName},
		})
	}

	return &SearchResponse{
		Query:   query,
		Total:   len(results),
		Results: results,
	}, nil
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
