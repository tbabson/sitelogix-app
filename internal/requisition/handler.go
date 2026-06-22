package requisition

import (
	"github.com/gin-gonic/gin"
	"github.com/sitelogix/backend/pkg/middleware"
	"github.com/sitelogix/backend/pkg/response"
)

type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID, _ := c.Get(middleware.ContextUserID)
	rq, err := h.svc.Create(req, userID.(string))
	if err != nil {
		response.InternalError(c, "failed to create requisition")
		return
	}
	response.Created(c, rq)
}

func (h *Handler) List(c *gin.Context) {
	role, _ := c.Get(middleware.ContextRole)
	projectID := c.Query("project_id")
	// Non-admins must filter by project or see all their projects; we expose filter to frontend
	_ = role
	rqs, err := h.svc.List(projectID, c.Query("status"))
	if err != nil {
		response.InternalError(c, "failed to list requisitions")
		return
	}
	response.OK(c, rqs)
}

func (h *Handler) Get(c *gin.Context) {
	rq, err := h.svc.Get(c.Param("id"))
	if err != nil {
		response.NotFound(c, "requisition not found")
		return
	}
	response.OK(c, rq)
}

func (h *Handler) Approve(c *gin.Context) {
	var req ApproveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.Approve(c.Param("id"), req); err != nil {
		response.InternalError(c, "failed to approve requisition")
		return
	}
	response.OKMsg(c, "requisition approved")
}

func (h *Handler) Reject(c *gin.Context) {
	var req RejectRequest
	_ = c.ShouldBindJSON(&req)
	if err := h.svc.Reject(c.Param("id"), req.Notes); err != nil {
		response.InternalError(c, "failed to reject requisition")
		return
	}
	response.OKMsg(c, "requisition rejected")
}

func (h *Handler) Fulfill(c *gin.Context) {
	userID, _ := c.Get(middleware.ContextUserID)
	if err := h.svc.Fulfill(c.Param("id"), userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OKMsg(c, "requisition fulfilled — inventory transferred")
}

func (h *Handler) Cancel(c *gin.Context) {
	userID, _ := c.Get(middleware.ContextUserID)
	role, _ := c.Get(middleware.ContextRole)
	if err := h.svc.Cancel(c.Param("id"), userID.(string), role.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OKMsg(c, "requisition cancelled")
}
