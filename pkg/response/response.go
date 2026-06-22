package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type PaginatedResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Meta    PaginMeta   `json:"meta"`
}

type PaginMeta struct {
	Total  int `json:"total"`
	Page   int `json:"page"`
	Limit  int `json:"limit"`
	Pages  int `json:"pages"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Success: true, Data: data})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{Success: true, Data: data})
}

func OKMsg(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{Success: true, Message: msg})
}

func BadRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, Response{Success: false, Error: msg})
}

func Unauthorized(c *gin.Context, msg string) {
	c.JSON(http.StatusUnauthorized, Response{Success: false, Error: msg})
}

func Forbidden(c *gin.Context, msg string) {
	c.JSON(http.StatusForbidden, Response{Success: false, Error: msg})
}

func NotFound(c *gin.Context, msg string) {
	c.JSON(http.StatusNotFound, Response{Success: false, Error: msg})
}

func Conflict(c *gin.Context, msg string) {
	c.JSON(http.StatusConflict, Response{Success: false, Error: msg})
}

func InternalError(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, Response{Success: false, Error: msg})
}

func Paginated(c *gin.Context, data interface{}, total, page, limit int) {
	pages := total / limit
	if total%limit != 0 {
		pages++
	}
	c.JSON(http.StatusOK, PaginatedResponse{
		Success: true,
		Data:    data,
		Meta:    PaginMeta{Total: total, Page: page, Limit: limit, Pages: pages},
	})
}
