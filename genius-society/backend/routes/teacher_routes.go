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

	assignmentRepo := teacher.NewAssignmentRepository(teacherDB)
	assignmentSvc := teacher.NewAssignmentService(assignmentRepo, st, notifSvc, notifRepo)
	assignmentHandler := teacher.NewAssignmentHandler(assignmentSvc)

	projectRepo := teacher.NewProjectRepository(teacherDB)
	projectSvc := teacher.NewProjectService(projectRepo, st, notifSvc, notifRepo)
	projectHandler := teacher.NewProjectHandler(projectSvc)

	assessmentRepo := teacher.NewAssessmentRepository(teacherDB)
	assessmentSvc := teacher.NewAssessmentService(assessmentRepo, notifSvc, notifRepo)
	assessmentHandler := teacher.NewAssessmentHandler(assessmentSvc)

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

		// Modul 3.4: Assignments
		teacherGroup.GET("/assignments", assignmentHandler.ListAssignments)
		teacherGroup.POST("/assignments", assignmentHandler.CreateAssignment)
		teacherGroup.GET("/assignments/:id", assignmentHandler.GetAssignmentDetail)
		teacherGroup.PUT("/assignments/:id", assignmentHandler.UpdateAssignment)
		teacherGroup.DELETE("/assignments/:id", assignmentHandler.DeleteAssignment)
		teacherGroup.GET("/assignments/:id/files/:index/download", assignmentHandler.DownloadFile)
		teacherGroup.GET("/assignments/:id/submissions", assignmentHandler.ListSubmissions)
		teacherGroup.PUT("/assignments/:id/submissions/:submissionId/grade", assignmentHandler.GradeSubmission)
		teacherGroup.POST("/assignments/:id/submit", assignmentHandler.SubmitAssignment)

		// Modul 3.5: Projects
		teacherGroup.GET("/projects", projectHandler.ListProjects)
		teacherGroup.POST("/projects", projectHandler.CreateProject)
		teacherGroup.GET("/projects/:id", projectHandler.GetProjectDetail)
		teacherGroup.PUT("/projects/:id", projectHandler.UpdateProject)
		teacherGroup.DELETE("/projects/:id", projectHandler.DeleteProject)
		teacherGroup.GET("/projects/:id/files/:index/download", projectHandler.DownloadFile)
		teacherGroup.GET("/projects/:id/submissions", projectHandler.ListSubmissions)
		teacherGroup.PUT("/projects/:id/submissions/:submissionId/grade", projectHandler.GradeSubmission)
		teacherGroup.POST("/projects/:id/submit", projectHandler.SubmitProject)

		// Modul 3.6: Assessments
		teacherGroup.GET("/assessments", assessmentHandler.ListAssessments)
		teacherGroup.POST("/assessments", assessmentHandler.CreateAssessment)
		teacherGroup.GET("/assessments/:id", assessmentHandler.GetAssessmentDetail)
		teacherGroup.PUT("/assessments/:id", assessmentHandler.UpdateAssessment)
		teacherGroup.DELETE("/assessments/:id", assessmentHandler.DeleteAssessment)
		teacherGroup.POST("/assessments/:id/questions", assessmentHandler.AddQuestion)
		teacherGroup.PUT("/assessments/:id/questions/:questionId", assessmentHandler.UpdateQuestion)
		teacherGroup.DELETE("/assessments/:id/questions/:questionId", assessmentHandler.DeleteQuestion)
		teacherGroup.PUT("/assessments/:id/questions/reorder", assessmentHandler.ReorderQuestions)
		teacherGroup.POST("/assessments/:id/regrade-mcq", assessmentHandler.RegradeMCQ)
		teacherGroup.POST("/assessments/:id/publish", assessmentHandler.PublishAssessment)
		teacherGroup.POST("/assessments/:id/close", assessmentHandler.CloseAssessment)
		teacherGroup.GET("/assessments/:id/attempts", assessmentHandler.ListAttempts)
		teacherGroup.GET("/assessments/:id/attempts/:attemptId", assessmentHandler.GetAttemptDetail)
		teacherGroup.PUT("/assessments/:id/answers/:answerId/grade", assessmentHandler.GradeEssayAnswer)
		teacherGroup.POST("/assessments/:id/start", assessmentHandler.StartAssessment)
		teacherGroup.POST("/assessments/:id/submit", assessmentHandler.SubmitAssessment)
	}
}
