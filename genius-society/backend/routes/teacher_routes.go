package routes

import (
	"github.com/geniussociety/backend/internal/middleware"
	"github.com/geniussociety/backend/internal/model"
	"github.com/geniussociety/backend/internal/repository"
	"github.com/geniussociety/backend/internal/storage"
	"github.com/geniussociety/backend/internal/teacher"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var teacherDB *mongo.Database
var teacherStorage storage.FileStorage

// InitTeacher initializes the database connection for teacher routes.
func InitTeacher(db *mongo.Database) {
	teacherDB = db
}

// InitTeacherStorage initializes the file storage for teacher routes.
func InitTeacherStorage(fs storage.FileStorage) {
	teacherStorage = fs
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

	st := teacherStorage
	if st == nil {
		st = storage.NewFileStorage(nil)
	}

	materialRepo := teacher.NewMaterialRepository(teacherDB)
	materialSvc := teacher.NewMaterialService(materialRepo, st, notifSvc)
	materialHandler := teacher.NewMaterialHandler(materialSvc)

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

		// Modul 3.3: Materials
		teacherGroup.GET("/materials", materialHandler.ListMaterials)
		teacherGroup.POST("/materials", materialHandler.CreateMaterial)
		teacherGroup.GET("/materials/:id", materialHandler.GetMaterialDetail)
		teacherGroup.PUT("/materials/:id", materialHandler.UpdateMaterial)
		teacherGroup.DELETE("/materials/:id", materialHandler.DeleteMaterial)
		teacherGroup.GET("/materials/:id/files/:index/download", materialHandler.DownloadFile)
	}
}
