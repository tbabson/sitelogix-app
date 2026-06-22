package media

import (
	"errors"
	"io"
	"net/http"
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

func (h *Handler) Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.BadRequest(c, "file is required")
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		response.InternalError(c, "failed to read file")
		return
	}

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	projectID := c.PostForm("project_id")
	if projectID == "" {
		response.BadRequest(c, "project_id is required")
		return
	}

	var logID *string
	if lid := c.PostForm("log_id"); lid != "" {
		logID = &lid
	}

	var issueID *string
	if iid := c.PostForm("issue_id"); iid != "" {
		issueID = &iid
	}

	userID, _ := c.Get(middleware.ContextUserID)
	m, err := h.svc.Upload(data, header.Filename, contentType, projectID, userID.(string), logID, issueID)
	if err != nil {
		if errors.Is(err, ErrUnsupportedType) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, "upload failed: "+err.Error())
		return
	}
	response.Created(c, m)
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	list, total, err := h.svc.List(MediaFilter{
		ProjectID: c.Query("project_id"),
		LogID:     c.Query("log_id"),
		IssueID:   c.Query("issue_id"),
		DateFrom:  c.Query("date_from"),
		DateTo:    c.Query("date_to"),
		Page:      page,
		Limit:     limit,
	})
	if err != nil {
		response.InternalError(c, "failed to list media")
		return
	}
	response.Paginated(c, list, total, page, limit)
}

func (h *Handler) Get(c *gin.Context) {
	m, err := h.svc.Get(c.Param("id"))
	if err != nil {
		response.NotFound(c, "media not found")
		return
	}
	response.OK(c, m)
}

func (h *Handler) Raw(c *gin.Context) {
	data, contentType, err := h.svc.Proxy(c.Param("id"))
	if err != nil {
		response.NotFound(c, "media not found")
		return
	}
	c.Header("Cache-Control", "public, max-age=31536000, immutable")
	c.Data(http.StatusOK, contentType, data)
}

func (h *Handler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Param("id")); err != nil {
		response.InternalError(c, "failed to delete media")
		return
	}
	response.OKMsg(c, "media deleted")
}
