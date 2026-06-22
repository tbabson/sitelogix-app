package project

import (
	"github.com/gin-gonic/gin"
	"github.com/sitelogix/backend/pkg/middleware"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler, jwtSecret string) {
	authMW := middleware.Auth(jwtSecret)
	adminSupervisor := middleware.RequireRole("admin", "supervisor")
	adminOnly := middleware.RequireRole("admin")

	projects := rg.Group("/projects", authMW)
	{
		projects.GET("", h.List)
		projects.POST("", adminSupervisor, h.Create)
		projects.GET("/:id", h.Get)
		projects.PUT("/:id", adminSupervisor, h.Update)
		projects.DELETE("/:id", adminOnly, h.Delete)

		// lifecycle
		projects.POST("/:id/archive", adminOnly, h.Archive)
		projects.POST("/:id/restore", adminOnly, h.Restore)
		projects.POST("/:id/duplicate", adminSupervisor, h.Duplicate)

		// insights
		projects.GET("/:id/stats", h.GetStats)
		projects.GET("/:id/activity", h.GetActivity)

		// members
		projects.GET("/:id/members", h.ListMembers)
		projects.POST("/:id/members", adminSupervisor, h.AddMember)
		projects.DELETE("/:id/members/:userID", adminSupervisor, h.RemoveMember)
		projects.PATCH("/:id/members/:userID/role", adminSupervisor, h.UpdateMemberRole)

		// milestones
		projects.GET("/:id/milestones", h.ListMilestones)
		projects.POST("/:id/milestones", adminSupervisor, h.CreateMilestone)
		projects.PUT("/:id/milestones/:milestoneID", adminSupervisor, h.UpdateMilestone)
		projects.DELETE("/:id/milestones/:milestoneID", adminSupervisor, h.DeleteMilestone)
	}
}
