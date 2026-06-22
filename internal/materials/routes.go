package materials

import (
	"github.com/gin-gonic/gin"
	"github.com/sitelogix/backend/pkg/middleware"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler, jwtSecret string) {
	authMW  := middleware.Auth(jwtSecret)
	adminMW := middleware.RequireRole("admin")

	mat := rg.Group("/materials", authMW)
	{
		mat.GET("",        h.List)
		mat.GET("/:id",    h.Get)
		mat.POST("",       adminMW, h.Create)
		mat.PUT("/:id",    adminMW, h.Update)
		mat.DELETE("/:id", adminMW, h.Delete)
	}
}
