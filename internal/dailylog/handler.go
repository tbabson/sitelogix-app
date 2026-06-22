package dailylog

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sitelogix/backend/internal/notification"
	"github.com/sitelogix/backend/pkg/middleware"
	"github.com/sitelogix/backend/pkg/response"
)

type Handler struct {
	svc   *Service
	notif *notification.Service
}

func NewHandler(svc *Service, notif *notification.Service) *Handler {
	return &Handler{svc: svc, notif: notif}
}

func (h *Handler) GetWeather(c *gin.Context) {
	lat, err1 := strconv.ParseFloat(c.Query("lat"), 64)
	lng, err2 := strconv.ParseFloat(c.Query("lng"), 64)
	if err1 != nil || err2 != nil {
		response.BadRequest(c, "lat and lng query params are required")
		return
	}
	w, err := h.svc.GetWeather(lat, lng)
	if err != nil {
		response.InternalError(c, "weather fetch failed")
		return
	}
	response.OK(c, map[string]string{"weather": w})
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID, _ := c.Get(middleware.ContextUserID)
	log, err := h.svc.Create(req, userID.(string))
	if err != nil {
		response.InternalError(c, "failed to create log")
		return
	}
	response.Created(c, log)
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	f := LogFilter{
		ProjectID: c.Query("project_id"),
		Status:    c.Query("status"),
		DateFrom:  c.Query("date_from"),
		DateTo:    c.Query("date_to"),
		Page:      page,
		Limit:     limit,
	}

	role, _ := c.Get(middleware.ContextRole)
	name, _ := c.Get(middleware.ContextName)
	// Supervisors see only their own logs; admin and client see all
	if role.(string) == "supervisor" {
		f.Username = name.(string)
	}

	logs, total, err := h.svc.List(f)
	if err != nil {
		response.InternalError(c, "failed to list logs")
		return
	}
	response.Paginated(c, logs, total, page, limit)
}

func (h *Handler) Get(c *gin.Context) {
	log, err := h.svc.Get(c.Param("id"))
	if err != nil {
		response.NotFound(c, "log not found")
		return
	}
	response.OK(c, log)
}

func (h *Handler) Update(c *gin.Context) {
	var req UpdateLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID, _ := c.Get(middleware.ContextUserID)
	role, _ := c.Get(middleware.ContextRole)
	log, err := h.svc.Update(c.Param("id"), userID.(string), role.(string), req)
	if err != nil {
		if errors.Is(err, ErrForbidden) {
			response.Forbidden(c, err.Error())
			return
		}
		response.InternalError(c, "failed to update log")
		return
	}
	response.OK(c, log)
}

func (h *Handler) Submit(c *gin.Context) {
	userID, _ := c.Get(middleware.ContextUserID)
	logID := c.Param("id")

	log, err := h.svc.Submit(logID, userID.(string))
	if err != nil {
		if errors.Is(err, ErrForbidden) {
			response.Forbidden(c, err.Error())
			return
		}
		response.InternalError(c, "failed to submit log")
		return
	}

	// Notify admins/clients that a log needs approval
	h.notif.Notify(notification.SendRequest{
		UserID: userID.(string),
		Title:  "Log Submitted for Approval",
		Body:   fmt.Sprintf("A daily log dated %s has been submitted and awaits approval.", log.Date.Format("02 Jan 2006")),
		Type:   notification.TypeApprovalNeeded,
		ReferenceID:   &log.ID,
		ReferenceType: strPtr("daily_log"),
	})

	response.OK(c, log)
}

func (h *Handler) Approve(c *gin.Context) {
	logID := c.Param("id")
	log, err := h.svc.Approve(logID)
	if err != nil {
		response.InternalError(c, "failed to approve log")
		return
	}

	h.notif.Notify(notification.SendRequest{
		UserID: log.CreatedBy,
		Title:  "Daily Log Approved",
		Body:   fmt.Sprintf("Your log dated %s has been approved.", log.Date.Format("02 Jan 2006")),
		Type:   notification.TypeLogApproved,
		ReferenceID:   &log.ID,
		ReferenceType: strPtr("daily_log"),
	})

	response.OK(c, log)
}

func (h *Handler) Reject(c *gin.Context) {
	logID := c.Param("id")
	log, err := h.svc.Reject(logID)
	if err != nil {
		response.InternalError(c, "failed to reject log")
		return
	}

	h.notif.Notify(notification.SendRequest{
		UserID: log.CreatedBy,
		Title:  "Daily Log Rejected",
		Body:   fmt.Sprintf("Your log dated %s was rejected. Please review and resubmit.", log.Date.Format("02 Jan 2006")),
		Type:   notification.TypeLogRejected,
		ReferenceID:   &log.ID,
		ReferenceType: strPtr("daily_log"),
	})

	response.OK(c, log)
}

func (h *Handler) Delete(c *gin.Context) {
	userID, _ := c.Get(middleware.ContextUserID)
	role, _ := c.Get(middleware.ContextRole)
	if err := h.svc.Delete(c.Param("id"), userID.(string), role.(string)); err != nil {
		if errors.Is(err, ErrForbidden) {
			response.Forbidden(c, err.Error())
			return
		}
		response.InternalError(c, "failed to delete log")
		return
	}
	response.OKMsg(c, "log deleted")
}

func (h *Handler) ListLogMaterials(c *gin.Context) {
	mats, err := h.svc.ListLogMaterials(c.Param("id"))
	if err != nil {
		response.InternalError(c, "failed to list log materials")
		return
	}
	response.OK(c, mats)
}

func (h *Handler) AddLogMaterial(c *gin.Context) {
	var req AddLogMaterialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID, _ := c.Get(middleware.ContextUserID)
	role, _ := c.Get(middleware.ContextRole)
	lm, err := h.svc.AddLogMaterial(c.Param("id"), userID.(string), role.(string), req)
	if err != nil {
		if errors.Is(err, ErrForbidden) {
			response.Forbidden(c, err.Error())
			return
		}
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, lm)
}

func (h *Handler) RemoveLogMaterial(c *gin.Context) {
	userID, _ := c.Get(middleware.ContextUserID)
	role, _ := c.Get(middleware.ContextRole)
	err := h.svc.RemoveLogMaterial(c.Param("id"), c.Param("material_id"), userID.(string), role.(string))
	if err != nil {
		if errors.Is(err, ErrForbidden) {
			response.Forbidden(c, err.Error())
			return
		}
		response.BadRequest(c, err.Error())
		return
	}
	response.OKMsg(c, "material removed and inventory restored")
}

func strPtr(s string) *string { return &s }
