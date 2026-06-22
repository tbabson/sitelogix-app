package reports

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sitelogix/backend/pkg/response"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Generate(c *gin.Context) {
	var req ReportRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	pdf, filename, err := h.svc.Generate(req)
	if err != nil {
		response.InternalError(c, "failed to generate report: "+err.Error())
		return
	}

	c.Header("Content-Disposition", `attachment; filename="`+filename+`"`)
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Length", fmt.Sprintf("%d", len(pdf)))
	c.Data(200, "application/pdf", pdf)
}
