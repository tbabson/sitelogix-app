package issues

import (
	"github.com/gin-gonic/gin"
	"github.com/sitelogix/backend/pkg/middleware"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler, jwtSecret string) {
	authMW := middleware.Auth(jwtSecret)

	iss := rg.Group("/issues", authMW)
	{
		iss.GET("", h.List)
		iss.POST("", h.Create)
		iss.GET("/:id", h.Get)
		iss.PUT("/:id", h.Update)
		iss.DELETE("/:id", middleware.RequireRole("admin", "supervisor"), h.Delete)
	}
}
