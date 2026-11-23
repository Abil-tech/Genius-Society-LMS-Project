package routes

import (
	studenthandler "github.com/geniussociety/backend/internal/handler/student"
	"github.com/geniussociety/backend/internal/middleware"
	"github.com/geniussociety/backend/internal/model"
	studentrepo "github.com/geniussociety/backend/internal/repository/student"
	studentsvc "github.com/geniussociety/backend/internal/service/student"
	"github.com/geniussociety/backend/internal/storage"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var studentDB *mongo.Database
var studentStorage storage.FileStorage

// InitStudent initializes the MongoDB connection for student routes.
func InitStudent(db *mongo.Database) {
	studentDB = db
}

// InitStudentStorage initializes file storage for student uploads/downloads.
func InitStudentStorage(fs storage.FileStorage) {
	studentStorage = fs
}

// RegisterStudent registers all HTTP routes under /api/student.
func RegisterStudent(protected *gin.RouterGroup) {
	if studentDB == nil {
		return
	}

	st := studentStorage
	if st == nil {
		st = storage.NewFileStorage(nil)
	}

	repo := studentrepo.NewStudentRepository(studentDB)
	svc := studentsvc.NewStudentService(repo, st)
	h := studenthandler.NewStudentHandler(svc)

	studentGroup := protected.Group("/student")
	studentGroup.Use(middleware.RequireRoles(model.RoleStudent))
	studentGroup.Use(middleware.StudentScopeMiddleware(repo))
	{
		// 1. Dashboard
		studentGroup.GET("/dashboard", h.GetDashboard)

		// 2. Materials
		studentGroup.GET("/materials", h.ListMaterials)
		studentGroup.GET("/materials/:id", h.GetMaterialDetail)
		studentGroup.GET("/materials/:id/files/:index/download", h.DownloadMaterialFile)

		// 3. Assignments
		studentGroup.GET("/assignments", h.ListAssignments)
		studentGroup.GET("/assignments/:id", h.GetAssignmentDetail)
		studentGroup.POST("/assignments/:id/submit", h.SubmitAssignment)
		studentGroup.PUT("/assignments/:id/resubmit", h.ResubmitAssignment)
		studentGroup.GET("/assignments/:id/files/:index/download", h.DownloadAssignmentFile)

		// 4. Projects
		studentGroup.GET("/projects", h.ListProjects)
		studentGroup.GET("/projects/:id", h.GetProjectDetail)
		studentGroup.POST("/projects/:id/submit", h.SubmitProject)
		studentGroup.PUT("/projects/:id/resubmit", h.ResubmitProject)
		studentGroup.GET("/projects/:id/files/:index/download", h.DownloadProjectFile)

		// 5. Assessments
		studentGroup.GET("/assessments", h.ListAssessments)
		studentGroup.GET("/assessments/:id", h.GetAssessmentDetail)
		studentGroup.POST("/assessments/:id/start", h.StartAssessment)
		studentGroup.POST("/assessments/:id/submit", h.SubmitAssessment)
		studentGroup.GET("/assessments/:id/result", h.GetAssessmentResult)

		// 6. Grades
		studentGroup.GET("/grades", h.GetGrades)
		studentGroup.GET("/grades/:subjectId", h.GetSubjectGradeDetail)

		// 7. Schedule
		studentGroup.GET("/schedule", h.GetSchedule)

		// 8. Ranking
		studentGroup.GET("/ranking", h.GetRanking)
	}
}
