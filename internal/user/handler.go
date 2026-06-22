package user

import (
	"strconv"

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

func (h *Handler) GetProfile(c *gin.Context) {
	userID, _ := c.Get(middleware.ContextUserID)
	u, err := h.svc.GetProfile(userID.(string))
	if err != nil {
		response.NotFound(c, "user not found")
		return
	}
	response.OK(c, u)
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	var req UpdateUserParams
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID, _ := c.Get(middleware.ContextUserID)
	u, err := h.svc.UpdateProfile(userID.(string), req)
	if err != nil {
		response.InternalError(c, "failed to update profile")
		return
	}
	response.OK(c, u)
}

func (h *Handler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	users, total, err := h.svc.ListUsers(page, limit)
	if err != nil {
		response.InternalError(c, "failed to list users")
		return
	}
	response.Paginated(c, users, total, page, limit)
}

func (h *Handler) DeactivateUser(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeactivateUser(id); err != nil {
		response.InternalError(c, "failed to deactivate user")
		return
	}
	response.OKMsg(c, "user deactivated")
}
