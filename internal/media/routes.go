package media

import (
	"github.com/gin-gonic/gin"
	"github.com/sitelogix/backend/pkg/middleware"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler, jwtSecret string) {
	authMW := middleware.Auth(jwtSecret)

	med := rg.Group("/media", authMW)
	{
		med.GET("", h.List)
		med.POST("/upload", h.Upload)
		med.GET("/:id", h.Get)
		med.DELETE("/:id", h.Delete)
	}

	// Raw file proxy — no auth so <img src> tags work (UUID is unguessable)
	rg.GET("/media/:id/raw", h.Raw)
}
