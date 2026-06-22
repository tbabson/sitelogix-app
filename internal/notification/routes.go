package notification

import (
	"github.com/gin-gonic/gin"
	"github.com/sitelogix/backend/pkg/middleware"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler, jwtSecret string) {
	authMW := middleware.Auth(jwtSecret)

	notifs := rg.Group("/notifications", authMW)
	{
		notifs.GET("", h.List)
		notifs.GET("/unread-count", h.UnreadCount)
		notifs.PUT("/:id/read", h.MarkRead)
		notifs.PUT("/read-all", h.MarkAllRead)
	}

	tokens := rg.Group("/device-tokens", authMW)
	{
		tokens.POST("", h.RegisterToken)
		tokens.DELETE("/:token", h.RemoveToken)
	}
}
