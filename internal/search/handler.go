package search

import (
	"github.com/gin-gonic/gin"
	"github.com/sitelogix/backend/pkg/middleware"
	"github.com/sitelogix/backend/pkg/response"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Search(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		response.BadRequest(c, "query parameter 'q' is required")
		return
	}

	userID, _ := c.Get(middleware.ContextUserID)
	role, _ := c.Get(middleware.ContextRole)

	result, err := h.svc.Search(q, userID.(string), role.(string))
	if err != nil {
		response.InternalError(c, "search failed")
		return
	}

	response.OK(c, result)
}
