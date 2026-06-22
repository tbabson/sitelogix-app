package reports

import (
	"github.com/gin-gonic/gin"
	"github.com/sitelogix/backend/pkg/middleware"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler, jwtSecret string) {
	authMW := middleware.Auth(jwtSecret)
	adminSupervisor := middleware.RequireRole("admin", "supervisor")

	rg.GET("/reports/generate", authMW, adminSupervisor, h.Generate)
}
