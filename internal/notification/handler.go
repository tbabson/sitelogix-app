package notification

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

func (h *Handler) List(c *gin.Context) {
	userID, _ := c.Get(middleware.ContextUserID)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	onlyUnread := c.Query("unread") == "true"

	list, total, err := h.svc.List(userID.(string), onlyUnread, page, limit)
	if err != nil {
		response.InternalError(c, "failed to list notifications")
		return
	}
	response.Paginated(c, list, total, page, limit)
}

func (h *Handler) UnreadCount(c *gin.Context) {
	userID, _ := c.Get(middleware.ContextUserID)
	count, err := h.svc.UnreadCount(userID.(string))
	if err != nil {
		response.InternalError(c, "failed to get unread count")
		return
	}
	response.OK(c, gin.H{"unread": count})
}

func (h *Handler) MarkRead(c *gin.Context) {
	userID, _ := c.Get(middleware.ContextUserID)
	if err := h.svc.MarkRead(c.Param("id"), userID.(string)); err != nil {
		response.InternalError(c, "failed to mark notification read")
		return
	}
	response.OKMsg(c, "marked as read")
}

func (h *Handler) MarkAllRead(c *gin.Context) {
	userID, _ := c.Get(middleware.ContextUserID)
	if err := h.svc.MarkAllRead(userID.(string)); err != nil {
		response.InternalError(c, "failed to mark all read")
		return
	}
	response.OKMsg(c, "all notifications marked as read")
}

func (h *Handler) RegisterToken(c *gin.Context) {
	var req RegisterTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID, _ := c.Get(middleware.ContextUserID)
	if err := h.svc.RegisterToken(userID.(string), req); err != nil {
		response.InternalError(c, "failed to register device token")
		return
	}
	response.OKMsg(c, "device token registered")
}

func (h *Handler) RemoveToken(c *gin.Context) {
	token := c.Param("token")
	if err := h.svc.RemoveToken(token); err != nil {
		response.InternalError(c, "failed to remove device token")
		return
	}
	response.OKMsg(c, "device token removed")
}
