package attendance

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sitelogix/backend/pkg/middleware"
	"github.com/sitelogix/backend/pkg/response"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateWorker(c *gin.Context) {
	var req CreateWorkerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	w, err := h.svc.CreateWorker(req)
	if err != nil {
		response.InternalError(c, "failed to create worker")
		return
	}
	response.Created(c, w)
}

func (h *Handler) ListWorkers(c *gin.Context) {
	workers, err := h.svc.ListWorkers(c.Query("project_id"))
	if err != nil {
		response.InternalError(c, "failed to list workers")
		return
	}
	response.OK(c, workers)
}

func (h *Handler) CheckIn(c *gin.Context) {
	var req CheckInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID, _ := c.Get(middleware.ContextUserID)
	a, err := h.svc.CheckIn(req, userID.(string))
	if err != nil {
		response.InternalError(c, "failed to record check-in")
		return
	}
	response.Created(c, a)
}

func (h *Handler) CheckOut(c *gin.Context) {
	var req CheckOutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	a, err := h.svc.CheckOut(c.Param("id"), req)
	if err != nil {
		response.InternalError(c, "failed to record check-out")
		return
	}
	response.OK(c, a)
}

func (h *Handler) LinkUser(c *gin.Context) {
	var req LinkUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	w, err := h.svc.LinkUser(c.Param("id"), req.UserID)
	if err != nil {
		response.InternalError(c, "failed to link user account")
		return
	}
	response.OK(c, w)
}

func (h *Handler) MyWorkerProfiles(c *gin.Context) {
	userID, _ := c.Get(middleware.ContextUserID)
	workers, err := h.svc.MyWorkerProfiles(userID.(string))
	if err != nil {
		response.InternalError(c, "failed to get worker profiles")
		return
	}
	response.OK(c, workers)
}

func (h *Handler) SelfCheckIn(c *gin.Context) {
	var req SelfCheckInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID, _ := c.Get(middleware.ContextUserID)
	a, err := h.svc.SelfCheckIn(userID.(string), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, a)
}

func (h *Handler) SelfCheckOut(c *gin.Context) {
	userID, _ := c.Get(middleware.ContextUserID)
	a, err := h.svc.SelfCheckOut(userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OK(c, a)
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	records, total, err := h.svc.List(AttendanceFilter{
		ProjectID: c.Query("project_id"),
		WorkerID:  c.Query("worker_id"),
		DateFrom:  c.Query("date_from"),
		DateTo:    c.Query("date_to"),
		Page:      page,
		Limit:     limit,
	})
	if err != nil {
		response.InternalError(c, "failed to list attendance")
		return
	}
	response.Paginated(c, records, total, page, limit)
}
