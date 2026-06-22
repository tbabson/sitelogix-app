package issues

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sitelogix/backend/internal/notification"
	"github.com/sitelogix/backend/pkg/middleware"
	"github.com/sitelogix/backend/pkg/response"
)

type Handler struct {
	svc   *Service
	notif *notification.Service
}

func NewHandler(svc *Service, notif *notification.Service) *Handler {
	return &Handler{svc: svc, notif: notif}
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateIssueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID, _ := c.Get(middleware.ContextUserID)
	issue, err := h.svc.Create(req, userID.(string))
	if err != nil {
		response.InternalError(c, "failed to create issue")
		return
	}

	// Notify assigned user if set
	if issue.AssignedTo != nil {
		h.notif.Notify(notification.SendRequest{
			UserID: *issue.AssignedTo,
			Title:  "Issue Assigned to You",
			Body:   fmt.Sprintf("You have been assigned issue: %s (Priority: %s)", issue.Title, issue.Priority),
			Type:   notification.TypeIssueAssigned,
			ReferenceID:   &issue.ID,
			ReferenceType: strPtr("issue"),
		})
	}

	response.Created(c, issue)
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	issues, total, err := h.svc.List(IssueFilter{
		ProjectID:  c.Query("project_id"),
		AssignedTo: c.Query("assigned_to"),
		Status:     c.Query("status"),
		Priority:   c.Query("priority"),
		Page:       page,
		Limit:      limit,
	})
	if err != nil {
		response.InternalError(c, "failed to list issues")
		return
	}
	response.Paginated(c, issues, total, page, limit)
}

func (h *Handler) Get(c *gin.Context) {
	issue, err := h.svc.Get(c.Param("id"))
	if err != nil {
		response.NotFound(c, "issue not found")
		return
	}
	response.OK(c, issue)
}

func (h *Handler) Update(c *gin.Context) {
	var req UpdateIssueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	issue, err := h.svc.Update(c.Param("id"), req)
	if err != nil {
		response.InternalError(c, "failed to update issue")
		return
	}

	// Notify newly assigned user
	if req.AssignedTo != nil && issue.AssignedTo != nil {
		h.notif.Notify(notification.SendRequest{
			UserID: *issue.AssignedTo,
			Title:  "Issue Assigned to You",
			Body:   fmt.Sprintf("You have been assigned issue: %s", issue.Title),
			Type:   notification.TypeIssueAssigned,
			ReferenceID:   &issue.ID,
			ReferenceType: strPtr("issue"),
		})
	}

	response.OK(c, issue)
}

func (h *Handler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Param("id")); err != nil {
		response.InternalError(c, "failed to delete issue")
		return
	}
	response.OKMsg(c, "issue deleted")
}

func strPtr(s string) *string { return &s }
