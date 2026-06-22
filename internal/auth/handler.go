package auth

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/sitelogix/backend/pkg/response"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	u, err := h.svc.Register(req)
	if err != nil {
		response.Conflict(c, "email already in use")
		return
	}

	response.Created(c, gin.H{"user": u.PublicView()})
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	pair, err := h.svc.Login(req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			response.Unauthorized(c, err.Error())
			return
		}
		response.InternalError(c, "login failed")
		return
	}

	response.OK(c, pair)
}

func (h *Handler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	pair, err := h.svc.Refresh(req.RefreshToken)
	if err != nil {
		response.Unauthorized(c, "invalid or expired refresh token")
		return
	}

	response.OK(c, pair)
}

func (h *Handler) Logout(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	_ = h.svc.Logout(req.RefreshToken)
	response.OKMsg(c, "logged out successfully")
}
