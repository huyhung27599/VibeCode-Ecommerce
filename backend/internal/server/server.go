package server

import (
	"log/slog"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/vibecode/ecommerce/backend/internal/config"
	"github.com/vibecode/ecommerce/backend/internal/domain"
	"github.com/vibecode/ecommerce/backend/internal/handler"
	"github.com/vibecode/ecommerce/backend/internal/middleware"
	"github.com/vibecode/ecommerce/backend/internal/repository"
	"github.com/vibecode/ecommerce/backend/internal/service"
	"gorm.io/gorm"
)

type Server struct {
	cfg    *config.Config
	log    *slog.Logger
	db     *gorm.DB
	rdb    *redis.Client
	engine *gin.Engine
}

func New(cfg *config.Config, log *slog.Logger, db *gorm.DB, rdb *redis.Client) *Server {
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(
		middleware.RequestID(),
		middleware.Logger(log),
		middleware.Recovery(log),
		cors.New(cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Request-ID"},
			ExposeHeaders:    []string{"X-Request-ID"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}),
	)

	s := &Server{cfg: cfg, log: log, db: db, rdb: rdb, engine: engine}
	s.registerRoutes()
	return s
}

func (s *Server) Router() *gin.Engine {
	return s.engine
}

func (s *Server) registerRoutes() {
	userRepo := repository.NewUserRepository(s.db)
	userSvc := service.NewUserService(userRepo)

	health := handler.NewHealth(s.db, s.rdb)
	userH := handler.NewUser(userSvc)

	s.engine.GET("/health", health.Live)
	s.engine.GET("/ready", health.Ready)

	v1 := s.engine.Group("/api/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})

		// Public: user registration
		v1.POST("/users", userH.Create)

		// Authenticated
		authed := v1.Group("", middleware.Auth(s.cfg.JWT.Secret))
		{
			authed.GET("/users/me", userH.Me)
		}

		// Admin-only
		admin := v1.Group("/users",
			middleware.Auth(s.cfg.JWT.Secret),
			middleware.RequireRole(string(domain.RoleAdmin)),
		)
		{
			admin.GET("", userH.List)
			admin.GET("/:id", userH.GetByID)
			admin.PATCH("/:id", userH.Update)
			admin.DELETE("/:id", userH.Delete)
		}
	}
}
