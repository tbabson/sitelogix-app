package materials

import (
	"github.com/gin-gonic/gin"
	"github.com/sitelogix/backend/pkg/response"
)

type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Create(c *gin.Context) {
	var req CreateMaterialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	m, err := h.svc.Create(req)
	if err != nil {
		response.InternalError(c, "failed to create material")
		return
	}
	response.Created(c, m)
}

func (h *Handler) List(c *gin.Context) {
	ms, err := h.svc.List()
	if err != nil {
		response.InternalError(c, "failed to list materials")
		return
	}
	response.OK(c, ms)
}

func (h *Handler) Get(c *gin.Context) {
	m, err := h.svc.Get(c.Param("id"))
	if err != nil {
		response.NotFound(c, "material not found")
		return
	}
	response.OK(c, m)
}

func (h *Handler) Update(c *gin.Context) {
	var req UpdateMaterialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	m, err := h.svc.Update(c.Param("id"), req)
	if err != nil {
		response.InternalError(c, "failed to update material")
		return
	}
	response.OK(c, m)
}

func (h *Handler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Param("id")); err != nil {
		response.InternalError(c, "failed to delete material")
		return
	}
	response.OKMsg(c, "material deleted")
}
