package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/vibecode/ecommerce/backend/pkg/response"
)

func Recovery(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Error("panic recovered",
					"error", err,
					"stack", string(debug.Stack()),
					"request_id", c.GetString("request_id"),
				)
				response.Fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
			}
		}()
		c.Next()
	}
}
