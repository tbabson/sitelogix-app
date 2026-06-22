package inventory

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sitelogix/backend/pkg/middleware"
	"github.com/sitelogix/backend/pkg/response"
)

type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) GetHeadOffice(c *gin.Context) {
	items, err := h.svc.GetHeadOffice()
	if err != nil {
		response.InternalError(c, "failed to load inventory")
		return
	}
	response.OK(c, items)
}

func (h *Handler) GetProject(c *gin.Context) {
	items, err := h.svc.GetProject(c.Param("project_id"))
	if err != nil {
		response.InternalError(c, "failed to load project inventory")
		return
	}
	response.OK(c, items)
}

func (h *Handler) AddStock(c *gin.Context) {
	var req AddStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID, _ := c.Get(middleware.ContextUserID)
	if err := h.svc.AddStock(req, userID.(string)); err != nil {
		response.InternalError(c, "failed to add stock")
		return
	}
	response.OKMsg(c, "stock added")
}

func (h *Handler) Transfer(c *gin.Context) {
	var req TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID, _ := c.Get(middleware.ContextUserID)
	if err := h.svc.Transfer(req, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OKMsg(c, "transfer completed")
}

func (h *Handler) GetPendingTransfers(c *gin.Context) {
	txs, err := h.svc.GetPendingTransfers(c.Param("project_id"))
	if err != nil {
		response.InternalError(c, "failed to load pending transfers")
		return
	}
	response.OK(c, txs)
}

func (h *Handler) ConfirmReceipt(c *gin.Context) {
	userID, _ := c.Get(middleware.ContextUserID)
	if err := h.svc.ConfirmReceipt(c.Param("id"), userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OKMsg(c, "receipt confirmed")
}

func (h *Handler) GetTransactions(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	txs, err := h.svc.GetTransactions(c.Query("material_id"), c.Query("project_id"), limit)
	if err != nil {
		response.InternalError(c, "failed to load transactions")
		return
	}
	response.OK(c, txs)
}
