package user

import (
	"github.com/gin-gonic/gin"
	"github.com/sitelogix/backend/pkg/middleware"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler, jwtSecret string) {
	authMW := middleware.Auth(jwtSecret)
	adminMW := middleware.RequireRole("admin")

	users := rg.Group("/users", authMW)
	{
		users.GET("/me", h.GetProfile)
		users.PUT("/me", h.UpdateProfile)

		users.GET("", adminMW, h.ListUsers)
		users.DELETE("/:id", adminMW, h.DeactivateUser)
	}
}
