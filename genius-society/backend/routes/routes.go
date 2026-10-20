package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/geniussociety/backend/internal/config"
	"github.com/geniussociety/backend/internal/handler"
	"github.com/geniussociety/backend/internal/middleware"
	"github.com/geniussociety/backend/internal/model"
	"github.com/geniussociety/backend/internal/service"
)

func Setup(
	cfg *config.Config,
	authSvc *service.AuthService,
	userSvc *service.UserService,
) *gin.Engine {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{cfg.CORSOrigin}
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	r.Use(cors.New(corsConfig))

	// Handlers
	authHandler := handler.NewAuthHandler(authSvc, cfg)
	userHandler := handler.NewUserHandler(userSvc)

	api := r.Group("/api")

	// Auth routes (public)
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/logout", authHandler.Logout)
		authGroup.GET("/me", middleware.Auth(authSvc), authHandler.Me)
	}

	// Protected routes
	protected := api.Group("/")
	protected.Use(middleware.Auth(authSvc))
	protected.Use(middleware.PrincipalReadOnly())

	// User management
	// Super Admin: CRUD all users
	// Admin: CRUD non-admin users
	users := protected.Group("/users")
	{
		users.GET("", middleware.RequireRoles(model.RoleSuperAdmin, model.RoleAdmin), userHandler.List)
		users.POST("", middleware.RequireRoles(model.RoleSuperAdmin, model.RoleAdmin), userHandler.Create)
		users.GET("/:id", middleware.RequireRoles(model.RoleSuperAdmin, model.RoleAdmin), userHandler.GetOne)
		users.PUT("/:id", middleware.RequireRoles(model.RoleSuperAdmin, model.RoleAdmin), userHandler.Update)
		users.DELETE("/:id", middleware.RequireRoles(model.RoleSuperAdmin, model.RoleAdmin), userHandler.Delete)
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r
}
