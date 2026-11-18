package teachersvc

import (
	"context"
	"testing"
	"time"

	teacherdto "github.com/geniussociety/backend/internal/dto/teacher"
	"github.com/geniussociety/backend/internal/middleware"
	"github.com/geniussociety/backend/internal/model"
	"github.com/geniussociety/backend/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type mockGradeRepo struct {
	grades   []*model.Grade
	students []*model.Student
	users    map[bson.ObjectID]*model.User
	config   *model.GradingConfig
}

func (m *mockGradeRepo) FindGradesByClassSubject(ctx context.Context, schoolID, classID, subjectID, academicYearID bson.ObjectID) ([]*model.Grade, error) {
	return m.grades, nil
}
func (m *mockGradeRepo) FindStudentsByClassID(ctx context.Context, schoolID, classID bson.ObjectID) ([]*model.Student, error) {
	return m.students, nil
}
func (m *mockGradeRepo) FindUsersByIDs(ctx context.Context, userIDs []bson.ObjectID) (map[bson.ObjectID]*model.User, error) {
	return m.users, nil
}
func (m *mockGradeRepo) FindGradingConfig(ctx context.Context, schoolID, academicYearID, subjectID bson.ObjectID) (*model.GradingConfig, error) {
	return m.config, nil
}
func (m *mockGradeRepo) UpsertExamGrade(ctx context.Context, grade *model.Grade) (*model.Grade, error) {
	grade.ID = bson.NewObjectID()
	grade.CreatedAt = time.Now()
	return grade, nil
}
func (m *mockGradeRepo) FindStudentByID(ctx context.Context, studentID bson.ObjectID) (*model.Student, error) {
	for _, s := range m.students {
		if s.ID == studentID {
			return s, nil
		}
	}
	return nil, nil
}
func (m *mockGradeRepo) FindClassByID(ctx context.Context, classID bson.ObjectID) (*model.Class, error) {
	return &model.Class{ID: classID, Name: "X-IPA-1"}, nil
}
func (m *mockGradeRepo) FindSubjectByID(ctx context.Context, subjectID bson.ObjectID) (*model.Subject, error) {
	return &model.Subject{ID: subjectID, Name: "Matematika"}, nil
}
func (m *mockGradeRepo) BulkUpsertSTS(ctx context.Context, schoolID, teacherID bson.ObjectID, grades []*model.Grade, audits []*model.AuditLog) error {
	return nil
}
func (m *mockGradeRepo) FindExistingExamGrade(ctx context.Context, schoolID, studentID, subjectID, academicYearID bson.ObjectID, semester int, gradeType model.GradeType) (*model.Grade, error) {
	return nil, nil
}
func (m *mockGradeRepo) FindAssignment(ctx context.Context, schoolID, teacherID, id bson.ObjectID) (*model.Assignment, error) {
	return nil, nil
}
func (m *mockGradeRepo) HasAssignmentSubmission(ctx context.Context, schoolID, assignmentID, studentID bson.ObjectID) (bool, error) {
	return false, nil
}
func (m *mockGradeRepo) FindProject(ctx context.Context, schoolID, teacherID, id bson.ObjectID) (*model.Project, error) {
	return nil, nil
}
func (m *mockGradeRepo) HasProjectSubmission(ctx context.Context, schoolID, projectID, studentID bson.ObjectID) (bool, error) {
	return false, nil
}
func (m *mockGradeRepo) FindAssessment(ctx context.Context, schoolID, teacherID, id bson.ObjectID) (*model.Assessment, error) {
	return nil, nil
}
func (m *mockGradeRepo) HasAssessmentAttempt(ctx context.Context, schoolID, assessmentID, studentID bson.ObjectID) (bool, error) {
	return false, nil
}
func (m *mockGradeRepo) SetManualItemGrade(ctx context.Context, grade *model.Grade, audit *model.AuditLog) error {
	return nil
}
func (m *mockGradeRepo) DeleteManualItemGrade(ctx context.Context, schoolID, sourceID, studentID bson.ObjectID, audit *model.AuditLog) error {
	return nil
}
func (m *mockGradeRepo) FindManualItemGrade(ctx context.Context, schoolID, sourceID, studentID bson.ObjectID) (*model.Grade, error) {
	return nil, nil
}
func (m *mockGradeRepo) FindLearningOutcomesForRecap(ctx context.Context, schoolID, academicYearID, subjectID, classID bson.ObjectID, semester int) ([]*model.LearningOutcome, error) {
	return nil, nil
}
func (m *mockGradeRepo) FindItemsForRecap(ctx context.Context, schoolID, academicYearID, subjectID, classID bson.ObjectID) ([]*model.Assignment, []*model.Project, []*model.Assessment, error) {
	return nil, nil, nil, nil
}

func TestGradeService_GetGrades(t *testing.T) {
	schoolID := bson.NewObjectID()
	classID := bson.NewObjectID()
	subjectID := bson.NewObjectID()
	teacherID := bson.NewObjectID()
	ayID := bson.NewObjectID()
	studentID := bson.NewObjectID()
	userID := bson.NewObjectID()

	mockRepo := &mockGradeRepo{
		students: []*model.Student{
			{ID: studentID, UserID: userID, NIS: "1001"},
		},
		users: map[bson.ObjectID]*model.User{
			userID: {ID: userID, Name: "Ahmad Fauzi"},
		},
		grades: []*model.Grade{
			{StudentID: studentID, Type: model.GradeAssignment, Score: 80},
			{StudentID: studentID, Type: model.GradeAssessment, Score: 90},
			{StudentID: studentID, Type: model.GradeUTS, Score: 85},
			{StudentID: studentID, Type: model.GradeUAS, Score: 75},
		},
	}

	svc := NewGradeService(mockRepo)
	scope := &middleware.TeacherScope{
		SchoolID:       schoolID,
		TeacherID:      teacherID,
		AcademicYearID: ayID,
		TeachingPairs: []repository.ClassSubjectPair{
			{ClassID: classID, SubjectID: subjectID},
		},
	}

	rows, err := svc.GetGrades(context.Background(), scope, classID, subjectID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}
	row := rows[0]
	if row.StudentName != "Ahmad Fauzi" {
		t.Errorf("expected student name Ahmad Fauzi, got %s", row.StudentName)
	}
	if row.DailyAvg == nil || *row.DailyAvg != 85.0 {
		t.Errorf("expected daily avg 85.0, got %v", row.DailyAvg)
	}
	// final score with 60/20/20: (85 * 0.6) + (85 * 0.2) + (75 * 0.2) = 51 + 17 + 15 = 83.0
	if row.FinalScore == nil || *row.FinalScore != 83.0 {
		t.Errorf("expected final score 83.0, got %v", row.FinalScore)
	}
}

func TestGradeService_UpdateGrade(t *testing.T) {
	schoolID := bson.NewObjectID()
	classID := bson.NewObjectID()
	subjectID := bson.NewObjectID()
	teacherID := bson.NewObjectID()
	ayID := bson.NewObjectID()
	studentID := bson.NewObjectID()

	mockRepo := &mockGradeRepo{
		students: []*model.Student{
			{ID: studentID, NIS: "1001"},
		},
		users: map[bson.ObjectID]*model.User{},
	}

	svc := NewGradeService(mockRepo)
	scope := &middleware.TeacherScope{
		SchoolID:       schoolID,
		TeacherID:      teacherID,
		AcademicYearID: ayID,
		Semester:       1,
		TeachingPairs: []repository.ClassSubjectPair{
			{ClassID: classID, SubjectID: subjectID},
		},
	}

	res, err := svc.UpdateGrade(context.Background(), scope, teacherdto.UpdateStudentGradePayload{
		StudentID: studentID.Hex(),
		SubjectID: subjectID.Hex(),
		ClassID:   classID.Hex(),
		Type:      "uts",
		Score:     95,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Score != 95 {
		t.Errorf("expected score 95, got %v", res.Score)
	}
	if res.Type != "uts" {
		t.Errorf("expected type uts, got %s", res.Type)
	}
}
