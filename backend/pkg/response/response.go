package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Envelope struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   *Error `json:"error,omitempty"`
	Meta    *Meta  `json:"meta,omitempty"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type Meta struct {
	Page       int   `json:"page,omitempty"`
	PageSize   int   `json:"page_size,omitempty"`
	Total      int64 `json:"total,omitempty"`
	TotalPages int   `json:"total_pages,omitempty"`
}

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Envelope{Success: true, Data: data})
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, Envelope{Success: true, Data: data})
}

func Paginated(c *gin.Context, data any, meta Meta) {
	c.JSON(http.StatusOK, Envelope{Success: true, Data: data, Meta: &meta})
}

func Fail(c *gin.Context, status int, code, message string, details ...any) {
	e := &Error{Code: code, Message: message}
	if len(details) > 0 {
		e.Details = details[0]
	}
	c.AbortWithStatusJSON(status, Envelope{Success: false, Error: e})
}

func BadRequest(c *gin.Context, message string, details ...any) {
	Fail(c, http.StatusBadRequest, "BAD_REQUEST", message, details...)
}

func Unauthorized(c *gin.Context, message string) {
	Fail(c, http.StatusUnauthorized, "UNAUTHORIZED", message)
}

func Forbidden(c *gin.Context, message string) {
	Fail(c, http.StatusForbidden, "FORBIDDEN", message)
}

func NotFound(c *gin.Context, message string) {
	Fail(c, http.StatusNotFound, "NOT_FOUND", message)
}

func Internal(c *gin.Context, message string) {
	Fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", message)
}
