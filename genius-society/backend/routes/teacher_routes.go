package routes

import (
	"github.com/geniussociety/backend/internal/middleware"
	"github.com/geniussociety/backend/internal/model"
	"github.com/geniussociety/backend/internal/repository"
	"github.com/geniussociety/backend/internal/teacher"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var teacherDB *mongo.Database

// InitTeacher initializes the database connection for teacher routes.
func InitTeacher(db *mongo.Database) {
	teacherDB = db
}

// RegisterTeacher registers all routes for the Teacher portal.
func RegisterTeacher(protected *gin.RouterGroup) {
	if teacherDB == nil {
		return
	}

	scopeRepo := repository.NewTeacherScopeRepository(teacherDB)
	classRepo := teacher.NewClassRepository(teacherDB)
	classSvc := teacher.NewClassService(classRepo)
	classHandler := teacher.NewClassHandler(classSvc)

	notifRepo := teacher.NewNotificationRepository(teacherDB)
	notifSvc := teacher.NewNotificationService(notifRepo)
	notifHandler := teacher.NewNotificationHandler(notifSvc)

	// Modul 3.2: Notifications (all authenticated users)
	notifications := protected.Group("/notifications")
	{
		notifications.GET("", notifHandler.ListNotifications)
		notifications.GET("/unread-count", notifHandler.GetUnreadCount)
		notifications.PATCH("/:id/read", notifHandler.MarkAsRead)
		notifications.PATCH("/read-all", notifHandler.MarkAllAsRead)
	}

	teacherGroup := protected.Group("/teacher")
	teacherGroup.Use(middleware.RequireRoles(model.RoleTeacher))
	teacherGroup.Use(middleware.TeacherScopeMiddleware(scopeRepo))
	{
		// Modul 3.1: Classes & Students
		teacherGroup.GET("/classes", classHandler.ListClasses)
		teacherGroup.GET("/classes/:id", classHandler.GetClassDetail)
		teacherGroup.GET("/classes/:id/students", classHandler.ListClassStudents)
		teacherGroup.GET("/subjects", classHandler.ListTeacherSubjects)
	}
}
