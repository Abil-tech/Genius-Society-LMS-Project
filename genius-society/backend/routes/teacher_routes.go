package routes

import (
	teacherhandler "github.com/geniussociety/backend/internal/handler/teacher"
	"github.com/geniussociety/backend/internal/middleware"
	"github.com/geniussociety/backend/internal/model"
	"github.com/geniussociety/backend/internal/repository"
	teacherrepo "github.com/geniussociety/backend/internal/repository/teacher"
	teachersvc "github.com/geniussociety/backend/internal/service/teacher"
	"github.com/geniussociety/backend/internal/storage"
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
	classRepo := teacherrepo.NewClassRepository(teacherDB)
	classSvc := teachersvc.NewClassService(classRepo)
	classHandler := teacherhandler.NewClassHandler(classSvc)

	notifRepo := teacherrepo.NewNotificationRepository(teacherDB)
	notifSvc := teachersvc.NewNotificationService(notifRepo)
	notifHandler := teacherhandler.NewNotificationHandler(notifSvc)

	st := teacherStorage
	if st == nil {
		st = storage.NewFileStorage(nil)
	}

	materialRepo := teacherrepo.NewMaterialRepository(teacherDB)
	materialSvc := teachersvc.NewMaterialService(materialRepo, st, notifSvc)
	materialHandler := teacherhandler.NewMaterialHandler(materialSvc)

	assignmentRepo := teacherrepo.NewAssignmentRepository(teacherDB)
	assignmentSvc := teachersvc.NewAssignmentService(assignmentRepo, st, notifSvc, notifRepo)
	assignmentHandler := teacherhandler.NewAssignmentHandler(assignmentSvc)

	projectRepo := teacherrepo.NewProjectRepository(teacherDB)
	projectSvc := teachersvc.NewProjectService(projectRepo, st, notifSvc, notifRepo)
	projectHandler := teacherhandler.NewProjectHandler(projectSvc)

	assessmentRepo := teacherrepo.NewAssessmentRepository(teacherDB)
	assessmentSvc := teachersvc.NewAssessmentService(assessmentRepo, notifSvc, notifRepo)
	assessmentHandler := teacherhandler.NewAssessmentHandler(assessmentSvc)

	dashboardRepo := teacherrepo.NewDashboardRepository(teacherDB)
	dashboardSvc := teachersvc.NewDashboardService(dashboardRepo, classRepo)
	dashboardHandler := teacherhandler.NewDashboardHandler(dashboardSvc)

	scheduleRepo := teacherrepo.NewScheduleRepository(teacherDB)
	scheduleSvc := teachersvc.NewScheduleService(scheduleRepo)
	scheduleHandler := teacherhandler.NewScheduleHandler(scheduleSvc)

	gradeRepo := teacherrepo.NewGradeRepository(teacherDB)
	gradeSvc := teachersvc.NewGradeService(gradeRepo)
	gradeHandler := teacherhandler.NewGradeHandler(gradeSvc)

	walasRepo := teacherrepo.NewWalasRepository(teacherDB)
	walasSvc := teachersvc.NewWalasService(walasRepo)
	walasHandler := teacherhandler.NewWalasHandler(walasSvc)

	loRepo := teacherrepo.NewLearningOutcomeRepository(teacherDB)
	loSvc := teachersvc.NewLearningOutcomeService(loRepo)
	loHandler := teacherhandler.NewLearningOutcomeHandler(loSvc)

	fileHandler := teacherhandler.NewFileHandler(st)

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
		// Teacher Dashboard
		teacherGroup.GET("/dashboard", dashboardHandler.GetDashboard)

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

		// Modul 3.7: Learning Outcomes (Pencapaian Pembelajaran)
		teacherGroup.GET("/learning-outcomes", loHandler.ListOutcomes)
		teacherGroup.POST("/learning-outcomes", loHandler.CreateOutcome)
		teacherGroup.GET("/learning-outcomes/:id", loHandler.GetOutcomeDetail)
		teacherGroup.PUT("/learning-outcomes/:id", loHandler.UpdateOutcome)
		teacherGroup.DELETE("/learning-outcomes/:id", loHandler.DeleteOutcome)

		// Modul 3.7: Learning Outcome Mapping
		teacherGroup.PUT("/assignments/:id/learning-outcome", loHandler.MapAssignmentOutcome)
		teacherGroup.PUT("/projects/:id/learning-outcome", loHandler.MapProjectOutcome)
		teacherGroup.PUT("/assessments/:id/learning-outcome", loHandler.MapAssessmentOutcome)

		// Modul 3.8: Grades, Exams & Recap
		teacherGroup.GET("/grades", gradeHandler.GetGrades)
		teacherGroup.PUT("/grades", gradeHandler.UpdateGrade)
		teacherGroup.PUT("/grades/exams", gradeHandler.InputSTSExams)
		teacherGroup.GET("/grades/recap", gradeHandler.GetGradeRecap)

		// Modul 3.8: Manual Grading for Students Without Submission
		teacherGroup.PUT("/assignments/:id/students/:studentId/grade", gradeHandler.SetManualAssignmentGrade)
		teacherGroup.DELETE("/assignments/:id/students/:studentId/grade", gradeHandler.DeleteManualAssignmentGrade)
		teacherGroup.PUT("/projects/:id/students/:studentId/grade", gradeHandler.SetManualProjectGrade)
		teacherGroup.DELETE("/projects/:id/students/:studentId/grade", gradeHandler.DeleteManualProjectGrade)
		teacherGroup.PUT("/assessments/:id/students/:studentId/grade", gradeHandler.SetManualAssessmentGrade)
		teacherGroup.DELETE("/assessments/:id/students/:studentId/grade", gradeHandler.DeleteManualAssessmentGrade)

		// Modul 3.9: Schedule
		teacherGroup.GET("/schedule", scheduleHandler.GetTeacherSchedule)

		// Walas: Homeroom features
		teacherGroup.GET("/walas/students", walasHandler.GetWalasStudents)
		teacherGroup.GET("/walas/grades", walasHandler.GetWalasGrades)

		// File: Presigned download URL for any R2 object key
		teacherGroup.GET("/files/download", fileHandler.GetPresignedDownloadURL)
	}

	// Unprefixed aliases for learning outcomes and schedules
	loAlias := protected.Group("/learning-outcomes")
	loAlias.Use(middleware.RequireRoles(model.RoleTeacher))
	loAlias.Use(middleware.TeacherScopeMiddleware(scopeRepo))
	{
		loAlias.GET("", loHandler.ListOutcomes)
		loAlias.POST("", loHandler.CreateOutcome)
		loAlias.GET("/:id", loHandler.GetOutcomeDetail)
		loAlias.PUT("/:id", loHandler.UpdateOutcome)
		loAlias.DELETE("/:id", loHandler.DeleteOutcome)
	}

	scheduleAlias := protected.Group("/schedules")
	scheduleAlias.Use(middleware.RequireRoles(model.RoleTeacher))
	scheduleAlias.Use(middleware.TeacherScopeMiddleware(scopeRepo))
	{
		scheduleAlias.GET("", scheduleHandler.GetTeacherSchedule)
	}
}
