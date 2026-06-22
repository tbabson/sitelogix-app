package attendance

import (
	"github.com/gin-gonic/gin"
	"github.com/sitelogix/backend/pkg/middleware"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler, jwtSecret string) {
	authMW := middleware.Auth(jwtSecret)
	adminSupervisor := middleware.RequireRole("admin", "supervisor")
	workerOnly := middleware.RequireRole("worker")

	att := rg.Group("/attendance", authMW)
	{
		att.GET("", adminSupervisor, h.List)
		att.POST("/checkin", adminSupervisor, h.CheckIn)
		att.PUT("/:id/checkout", adminSupervisor, h.CheckOut)

		// Worker self-service
		att.GET("/me", workerOnly, h.MyWorkerProfiles)
		att.POST("/me/checkin", workerOnly, h.SelfCheckIn)
		att.PUT("/me/checkout", workerOnly, h.SelfCheckOut)
	}

	workers := rg.Group("/workers", authMW)
	{
		workers.GET("", h.ListWorkers)
		workers.POST("", adminSupervisor, h.CreateWorker)
		workers.PATCH("/:id/link-user", adminSupervisor, h.LinkUser)
	}
}
