package search

import (
	"github.com/gin-gonic/gin"
	"github.com/sitelogix/backend/pkg/middleware"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler, jwtSecret string) {
	rg.GET("/search", middleware.Auth(jwtSecret), h.Search)
}
