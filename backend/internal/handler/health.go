package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Health struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewHealth(db *gorm.DB, rdb *redis.Client) *Health {
	return &Health{db: db, rdb: rdb}
}

func (h *Health) Live(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Health) Ready(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	checks := gin.H{"status": "ok"}
	status := http.StatusOK

	if sqlDB, err := h.db.DB(); err != nil {
		checks["database"] = "down"
		status = http.StatusServiceUnavailable
	} else if err := sqlDB.PingContext(ctx); err != nil {
		checks["database"] = "down"
		status = http.StatusServiceUnavailable
	} else {
		checks["database"] = "up"
	}

	if h.rdb != nil {
		if err := h.rdb.Ping(ctx).Err(); err != nil {
			checks["redis"] = "down"
		} else {
			checks["redis"] = "up"
		}
	}

	c.JSON(status, checks)
}
