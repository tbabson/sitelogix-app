package dashboard

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

func (h *Handler) Get(c *gin.Context) {
	role, _ := c.Get(middleware.ContextRole)
	userID, _ := c.Get(middleware.ContextUserID)

	stats, err := h.svc.GetStats(role.(string), userID.(string))
	if err != nil {
		response.InternalError(c, "failed to load stats")
		return
	}

	activity, err := h.svc.GetRecentActivity(role.(string), userID.(string))
	if err != nil {
		response.InternalError(c, "failed to load activity")
		return
	}

	response.OK(c, DashboardResponse{
		Stats:          *stats,
		RecentActivity: *activity,
	})
}
