package dailylog

import (
	"github.com/gin-gonic/gin"
	"github.com/sitelogix/backend/pkg/middleware"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler, jwtSecret string) {
	authMW      := middleware.Auth(jwtSecret)
	adminMW     := middleware.RequireRole("admin")
	canViewMW   := middleware.RequireRole("admin", "supervisor", "client") // workers excluded
	canWriteMW  := middleware.RequireRole("admin", "supervisor")           // workers + clients excluded

	logs := rg.Group("/logs", authMW)
	{
		logs.GET("",                              canViewMW,  h.List)
		logs.POST("",                             canWriteMW, h.Create)
		logs.GET("/weather",                      canWriteMW, h.GetWeather)
		logs.GET("/:id",                          canViewMW,  h.Get)
		logs.PUT("/:id",                          canWriteMW, h.Update)
		logs.DELETE("/:id",                       canWriteMW, h.Delete)
		logs.POST("/:id/submit",                  canWriteMW, h.Submit)
		logs.POST("/:id/approve",                 adminMW,    h.Approve)
		logs.POST("/:id/reject",                  adminMW,    h.Reject)
		logs.GET("/:id/materials",                canViewMW,  h.ListLogMaterials)
		logs.POST("/:id/materials",               canWriteMW, h.AddLogMaterial)
		logs.DELETE("/:id/materials/:material_id", canWriteMW, h.RemoveLogMaterial)
	}
}
