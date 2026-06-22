package requisition

import (
	"github.com/gin-gonic/gin"
	"github.com/sitelogix/backend/pkg/middleware"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler, jwtSecret string) {
	authMW  := middleware.Auth(jwtSecret)
	adminMW := middleware.RequireRole("admin")

	req := rg.Group("/requisitions", authMW)
	{
		req.GET("",              h.List)
		req.POST("",             h.Create)
		req.GET("/:id",          h.Get)
		req.DELETE("/:id",       h.Cancel)
		req.POST("/:id/approve", adminMW, h.Approve)
		req.POST("/:id/reject",  adminMW, h.Reject)
		req.POST("/:id/fulfill", adminMW, h.Fulfill)
	}
}
