package project

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

func (h *Handler) Create(c *gin.Context) {
	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID, _ := c.Get(middleware.ContextUserID)
	p, err := h.svc.Create(req, userID.(string))
	if err != nil {
		response.InternalError(c, "failed to create project")
		return
	}
	response.Created(c, p)
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	role, _ := c.Get(middleware.ContextRole)
	userID, _ := c.Get(middleware.ContextUserID)

	projects, total, err := h.svc.List(role.(string), userID.(string), page, limit)
	if err != nil {
		response.InternalError(c, "failed to list projects")
		return
	}
	response.Paginated(c, projects, total, page, limit)
}

func (h *Handler) Get(c *gin.Context) {
	p, err := h.svc.Get(c.Param("id"))
	if err != nil {
		response.NotFound(c, "project not found")
		return
	}
	response.OK(c, p)
}

func (h *Handler) Update(c *gin.Context) {
	var req UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	p, err := h.svc.Update(c.Param("id"), req)
	if err != nil {
		response.InternalError(c, "failed to update project")
		return
	}
	response.OK(c, p)
}

func (h *Handler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Param("id")); err != nil {
		response.InternalError(c, "failed to delete project")
		return
	}
	response.OKMsg(c, "project deleted")
}

func (h *Handler) Archive(c *gin.Context) {
	p, err := h.svc.Archive(c.Param("id"))
	if err != nil {
		response.InternalError(c, "failed to archive project")
		return
	}
	response.OK(c, p)
}

func (h *Handler) Restore(c *gin.Context) {
	p, err := h.svc.Restore(c.Param("id"))
	if err != nil {
		response.InternalError(c, "failed to restore project")
		return
	}
	response.OK(c, p)
}

func (h *Handler) Duplicate(c *gin.Context) {
	var req DuplicateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID, _ := c.Get(middleware.ContextUserID)
	p, err := h.svc.Duplicate(c.Param("id"), req, userID.(string))
	if err != nil {
		response.InternalError(c, "failed to duplicate project")
		return
	}
	response.Created(c, p)
}

func (h *Handler) GetStats(c *gin.Context) {
	stats, err := h.svc.GetStats(c.Param("id"))
	if err != nil {
		response.InternalError(c, "failed to get project stats")
		return
	}
	response.OK(c, stats)
}

func (h *Handler) GetActivity(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	items, err := h.svc.GetActivity(c.Param("id"), limit)
	if err != nil {
		response.InternalError(c, "failed to get activity feed")
		return
	}
	response.OK(c, items)
}

// ── Members ───────────────────────────────────────────────────────────────────

func (h *Handler) AddMember(c *gin.Context) {
	var req AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	m, err := h.svc.AddMember(c.Param("id"), req)
	if err != nil {
		response.InternalError(c, "failed to add member")
		return
	}
	response.Created(c, m)
}

func (h *Handler) RemoveMember(c *gin.Context) {
	if err := h.svc.RemoveMember(c.Param("id"), c.Param("userID")); err != nil {
		response.InternalError(c, "failed to remove member")
		return
	}
	response.OKMsg(c, "member removed")
}

func (h *Handler) ListMembers(c *gin.Context) {
	members, err := h.svc.ListMembers(c.Param("id"))
	if err != nil {
		response.InternalError(c, "failed to list members")
		return
	}
	response.OK(c, members)
}

func (h *Handler) UpdateMemberRole(c *gin.Context) {
	var req UpdateMemberRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.UpdateMemberRole(c.Param("id"), c.Param("userID"), req); err != nil {
		response.InternalError(c, "failed to update member role")
		return
	}
	response.OKMsg(c, "role updated")
}

// ── Milestones ────────────────────────────────────────────────────────────────

func (h *Handler) CreateMilestone(c *gin.Context) {
	var req CreateMilestoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	m, err := h.svc.CreateMilestone(c.Param("id"), req)
	if err != nil {
		response.InternalError(c, "failed to create milestone")
		return
	}
	response.Created(c, m)
}

func (h *Handler) ListMilestones(c *gin.Context) {
	milestones, err := h.svc.ListMilestones(c.Param("id"))
	if err != nil {
		response.InternalError(c, "failed to list milestones")
		return
	}
	response.OK(c, milestones)
}

func (h *Handler) UpdateMilestone(c *gin.Context) {
	var req UpdateMilestoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	m, err := h.svc.UpdateMilestone(c.Param("milestoneID"), req)
	if err != nil {
		response.InternalError(c, "failed to update milestone")
		return
	}
	response.OK(c, m)
}

func (h *Handler) DeleteMilestone(c *gin.Context) {
	if err := h.svc.DeleteMilestone(c.Param("milestoneID")); err != nil {
		response.InternalError(c, "failed to delete milestone")
		return
	}
	response.OKMsg(c, "milestone deleted")
}
