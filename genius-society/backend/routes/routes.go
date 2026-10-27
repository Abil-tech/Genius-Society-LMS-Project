package routes

import (
	"strings"

	"github.com/geniussociety/backend/internal/config"
	"github.com/geniussociety/backend/internal/handler"
	"github.com/geniussociety/backend/internal/middleware"
	"github.com/geniussociety/backend/internal/model"
	"github.com/geniussociety/backend/internal/repository"
	"github.com/geniussociety/backend/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Setup(
	cfg *config.Config,
	authSvc *service.AuthService,
	userSvc *service.UserService,
	adminDashboardSvc *service.AdminDashboardService,
	studentManagementSvc *service.StudentManagementService,
	teacherManagementSvc *service.TeacherManagementService,
	subjectManagementSvc *service.SubjectManagementService,
	classManagementSvc *service.ClassManagementService,
	academicYearSvc *service.AcademicYearService,
	userRepo *repository.UserRepository,
) *gin.Engine {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.RedirectTrailingSlash = false
	r.RedirectFixedPath = false
	r.Use(gin.Logger(), gin.Recovery())

	// CORS
	corsConfig := cors.DefaultConfig()
	corsOrigin := strings.TrimSpace(cfg.CORSOrigin)
	if corsOrigin == "" {
		corsOrigin = "http://localhost:5173"
	}
	origins := []string{corsOrigin}
	if cfg.AppEnv != "production" {
		if corsOrigin == "http://localhost:5173" {
			origins = append(origins, "http://127.0.0.1:5173")
		} else if corsOrigin == "http://127.0.0.1:5173" {
			origins = append(origins, "http://localhost:5173")
		}
	}
	corsConfig.AllowOrigins = origins
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.ExposeHeaders = []string{"Content-Length"}
	r.Use(cors.New(corsConfig))

	// Handlers
	authHandler := handler.NewAuthHandler(authSvc, cfg)
	userHandler := handler.NewUserHandler(userSvc)
	adminDashboardHandler := handler.NewAdminDashboardHandler(adminDashboardSvc)

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
	protected.Use(middleware.SetSchoolID(userRepo))
	protected.Use(middleware.PrincipalReadOnly())

	// Admin Dashboard endpoint (Hanya Admin dan Super Admin)
	protected.GET("/admin/dashboard", middleware.RequireRoles(model.RoleAdmin, model.RoleSuperAdmin), adminDashboardHandler.GetMetrics)

	// Student Management routes (Admin, Super Admin, Kepala Sekolah read-only)
	studentHandler := handler.NewStudentManagementHandler(studentManagementSvc)
	students := protected.Group("/students")
	{
		students.GET("", studentHandler.List)
		students.GET("/stats", studentHandler.Stats)
		students.GET("/:id", studentHandler.GetOne)
		students.POST("", middleware.RequireRoles(model.RoleAdmin, model.RoleSuperAdmin), studentHandler.Create)
		students.PUT("/:id", middleware.RequireRoles(model.RoleAdmin, model.RoleSuperAdmin), studentHandler.Update)
		students.POST("/bulk/class", middleware.RequireRoles(model.RoleAdmin, model.RoleSuperAdmin), studentHandler.BulkClass)
		students.POST("/bulk/status", middleware.RequireRoles(model.RoleAdmin, model.RoleSuperAdmin), studentHandler.BulkStatus)
		students.POST("/:id/account", middleware.RequireRoles(model.RoleAdmin, model.RoleSuperAdmin), studentHandler.CreateAccount)
		students.PATCH("/:id/account/toggle", middleware.RequireRoles(model.RoleAdmin, model.RoleSuperAdmin), studentHandler.ToggleAccount)
		students.DELETE("/:id", middleware.RequireRoles(model.RoleAdmin, model.RoleSuperAdmin), studentHandler.Delete)
	}

	// Teacher Management routes (Admin, Super Admin, Kepala Sekolah read-only)
	teacherHandler := handler.NewTeacherManagementHandler(teacherManagementSvc)
	teachers := protected.Group("/teachers")
	{
		teachers.GET("", teacherHandler.List)
		teachers.GET("/stats", teacherHandler.Stats)
		teachers.GET("/:id", teacherHandler.GetOne)
		teachers.POST("", middleware.RequireRoles(model.RoleAdmin, model.RoleSuperAdmin, model.RoleCurriculum), teacherHandler.Create)
		teachers.PUT("/:id", middleware.RequireRoles(model.RoleAdmin, model.RoleSuperAdmin, model.RoleCurriculum), teacherHandler.Update)
		teachers.DELETE("/:id", middleware.RequireRoles(model.RoleAdmin, model.RoleSuperAdmin), teacherHandler.Delete)
		teachers.POST("/:id/account", middleware.RequireRoles(model.RoleAdmin, model.RoleSuperAdmin), teacherHandler.CreateAccount)
		teachers.PATCH("/:id/account/toggle", middleware.RequireRoles(model.RoleAdmin, model.RoleSuperAdmin), teacherHandler.ToggleAccount)
	}

	// Subject Management routes (Admin, Super Admin, Kurikulum, Kepala Sekolah read-only)
	subjectHandler := handler.NewSubjectManagementHandler(subjectManagementSvc)
	subjects := protected.Group("/subjects")
	{
		subjects.GET("", subjectHandler.List)
		subjects.GET("/stats", subjectHandler.Stats)
		subjects.GET("/:id", subjectHandler.GetOne)
		subjects.POST("", middleware.RequireRoles(model.RoleAdmin, model.RoleSuperAdmin, model.RoleCurriculum), subjectHandler.Create)
		subjects.PUT("/:id", middleware.RequireRoles(model.RoleAdmin, model.RoleSuperAdmin, model.RoleCurriculum), subjectHandler.Update)
		subjects.DELETE("/:id", middleware.RequireRoles(model.RoleAdmin, model.RoleSuperAdmin, model.RoleCurriculum), subjectHandler.Delete)
	}

	// Class Management routes (Admin only, school scoped)
	classHandler := handler.NewClassManagementHandler(classManagementSvc)
	classes := protected.Group("/classes")
	{
		classes.GET("", middleware.RequireRoles(model.RoleAdmin), classHandler.List)
		classes.GET("/:id", middleware.RequireRoles(model.RoleAdmin), classHandler.GetOne)
		classes.POST("", middleware.RequireRoles(model.RoleAdmin), classHandler.Create)
		classes.PUT("/:id", middleware.RequireRoles(model.RoleAdmin), classHandler.Update)
		classes.DELETE("/:id", middleware.RequireRoles(model.RoleAdmin), classHandler.Archive)
	}

	// Academic Year routes (Admin only, school scoped)
	academicYearHandler := handler.NewAcademicYearHandler(academicYearSvc)
	academicYears := protected.Group("/academic-years")
	{
		academicYears.GET("", middleware.RequireRoles(model.RoleAdmin), academicYearHandler.List)
		academicYears.GET("/", middleware.RequireRoles(model.RoleAdmin), academicYearHandler.List)
		academicYears.GET("/:id", middleware.RequireRoles(model.RoleAdmin), academicYearHandler.GetOne)
		academicYears.POST("", middleware.RequireRoles(model.RoleAdmin), academicYearHandler.Create)
		academicYears.POST("/", middleware.RequireRoles(model.RoleAdmin), academicYearHandler.Create)
		academicYears.PUT("/:id", middleware.RequireRoles(model.RoleAdmin), academicYearHandler.Update)
		academicYears.PATCH("/:id/activate", middleware.RequireRoles(model.RoleAdmin), academicYearHandler.Activate)
		academicYears.PATCH("/:id/archive", middleware.RequireRoles(model.RoleAdmin), academicYearHandler.Archive)
		academicYears.DELETE("/:id", middleware.RequireRoles(model.RoleAdmin), academicYearHandler.Delete)

		// Promotion (Kenaikan Kelas)
		academicYears.GET("/:id/promotion/candidates", middleware.RequireRoles(model.RoleAdmin), academicYearHandler.GetPromotionCandidates)
		academicYears.POST("/:id/promotion/preview", middleware.RequireRoles(model.RoleAdmin), academicYearHandler.PreviewPromotion)
		academicYears.POST("/:id/promotion/execute", middleware.RequireRoles(model.RoleAdmin), academicYearHandler.ExecutePromotion)
		academicYears.GET("/:id/promotion/history", middleware.RequireRoles(model.RoleAdmin), academicYearHandler.GetPromotionHistory)

		// Graduation (Kelulusan)
		academicYears.GET("/:id/graduation/candidates", middleware.RequireRoles(model.RoleAdmin), academicYearHandler.GetGraduationCandidates)
		academicYears.POST("/:id/graduation/preview", middleware.RequireRoles(model.RoleAdmin), academicYearHandler.PreviewGraduation)
		academicYears.POST("/:id/graduation/execute", middleware.RequireRoles(model.RoleAdmin), academicYearHandler.ExecuteGraduation)
		academicYears.GET("/:id/graduation/history", middleware.RequireRoles(model.RoleAdmin), academicYearHandler.GetGraduationHistory)

		// Suspended decisions
		academicYears.GET("/:id/decisions/suspended", middleware.RequireRoles(model.RoleAdmin), academicYearHandler.GetSuspendedDecisions)
	}

	// Academic Decisions review & resolve (Admin only, school scoped)
	academicDecisions := protected.Group("/academic-decisions")
	{
		academicDecisions.GET("/:decisionId", middleware.RequireRoles(model.RoleAdmin), academicYearHandler.GetDecisionDetail)
		academicDecisions.POST("/:decisionId/review", middleware.RequireRoles(model.RoleAdmin), academicYearHandler.ReviewDecision)
		academicDecisions.POST("/:decisionId/resolve", middleware.RequireRoles(model.RoleAdmin), academicYearHandler.ResolveDecision)
	}

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
