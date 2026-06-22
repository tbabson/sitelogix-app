package inventory

import (
	"github.com/gin-gonic/gin"
	"github.com/sitelogix/backend/pkg/middleware"
)

func RegisterRoutes(rg *gin.RouterGroup, h *Handler, jwtSecret string) {
	authMW      := middleware.Auth(jwtSecret)
	adminMW     := middleware.RequireRole("admin")
	supervisorMW := middleware.RequireRole("admin", "supervisor")

	inv := rg.Group("/inventory", authMW)
	{
		inv.GET("/head-office",                              h.GetHeadOffice)
		inv.POST("/head-office/stock",                      adminMW, h.AddStock)
		inv.POST("/transfer",                               adminMW, h.Transfer)
		inv.GET("/projects/:project_id",                    h.GetProject)
		inv.GET("/projects/:project_id/pending-transfers",  h.GetPendingTransfers)
		inv.PATCH("/transactions/:id/receive",              supervisorMW, h.ConfirmReceipt)
		inv.GET("/transactions",                            h.GetTransactions)
	}
}
