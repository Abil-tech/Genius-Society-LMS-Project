package studentsvc

import (
	"context"
	"testing"
	"time"

	"github.com/geniussociety/backend/internal/middleware"
	"github.com/geniussociety/backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// MockStudentRepo implements studentrepo.StudentRepository for testing.
type mockStudentRepo struct {
	findGradesByStudentFn func(ctx context.Context, schoolID, studentID, academicYearID bson.ObjectID) ([]model.Grade, error)
	findClassStudentsFn   func(ctx context.Context, schoolID, classID bson.ObjectID) ([]model.Student, error)
	findSubjectsMapFn     func(ctx context.Context, subjectIDs []bson.ObjectID) (map[bson.ObjectID]*model.Subject, error)
	findUsersMapFn        func(ctx context.Context, userIDs []bson.ObjectID) (map[bson.ObjectID]*model.User, error)
}

func (m *mockStudentRepo) FindUserByID(ctx context.Context, id bson.ObjectID) (*model.User, error) {
	return nil, nil
}
func (m *mockStudentRepo) FindStudentByUserID(ctx context.Context, userID bson.ObjectID) (*model.Student, error) {
	return nil, nil
}
func (m *mockStudentRepo) FindActiveAcademicYear(ctx context.Context, schoolID bson.ObjectID) (*model.AcademicYear, error) {
	return nil, nil
}
func (m *mockStudentRepo) FindActiveEnrollment(ctx context.Context, schoolID, studentID, academicYearID bson.ObjectID) (*model.StudentEnrollment, error) {
	return nil, nil
}
func (m *mockStudentRepo) FindClassByID(ctx context.Context, classID bson.ObjectID) (*model.Class, error) {
	return &model.Class{ID: classID, Name: "7-A"}, nil
}
func (m *mockStudentRepo) FindMaterialsByClass(ctx context.Context, schoolID, classID bson.ObjectID, subjectID *bson.ObjectID, now time.Time) ([]model.Material, error) {
	return nil, nil
}
func (m *mockStudentRepo) FindMaterialByID(ctx context.Context, schoolID, materialID bson.ObjectID) (*model.Material, error) {
	return nil, nil
}
func (m *mockStudentRepo) FindAssignmentsByClass(ctx context.Context, schoolID, classID bson.ObjectID, subjectID *bson.ObjectID) ([]model.Assignment, error) {
	return nil, nil
}
func (m *mockStudentRepo) FindAssignmentByID(ctx context.Context, schoolID, assignmentID bson.ObjectID) (*model.Assignment, error) {
	return nil, nil
}
func (m *mockStudentRepo) FindAssignmentSubmissions(ctx context.Context, schoolID, assignmentID, studentID bson.ObjectID) ([]model.AssignmentSubmission, error) {
	return nil, nil
}
func (m *mockStudentRepo) FindLatestAssignmentSubmission(ctx context.Context, schoolID, assignmentID, studentID bson.ObjectID) (*model.AssignmentSubmission, error) {
	return nil, nil
}
func (m *mockStudentRepo) InsertAssignmentSubmission(ctx context.Context, sub *model.AssignmentSubmission) error {
	return nil
}
func (m *mockStudentRepo) UpdateAssignmentSubmission(ctx context.Context, sub *model.AssignmentSubmission) error {
	return nil
}
func (m *mockStudentRepo) FindProjectsByClass(ctx context.Context, schoolID, classID bson.ObjectID, subjectID *bson.ObjectID) ([]model.Project, error) {
	return nil, nil
}
func (m *mockStudentRepo) FindProjectByID(ctx context.Context, schoolID, projectID bson.ObjectID) (*model.Project, error) {
	return nil, nil
}
func (m *mockStudentRepo) FindProjectSubmissions(ctx context.Context, schoolID, projectID, studentID bson.ObjectID) ([]model.ProjectSubmission, error) {
	return nil, nil
}
func (m *mockStudentRepo) FindLatestProjectSubmission(ctx context.Context, schoolID, projectID, studentID bson.ObjectID) (*model.ProjectSubmission, error) {
	return nil, nil
}
func (m *mockStudentRepo) InsertProjectSubmission(ctx context.Context, sub *model.ProjectSubmission) error {
	return nil
}
func (m *mockStudentRepo) UpdateProjectSubmission(ctx context.Context, sub *model.ProjectSubmission) error {
	return nil
}
func (m *mockStudentRepo) FindAssessmentsByClass(ctx context.Context, schoolID, classID bson.ObjectID, subjectID *bson.ObjectID) ([]model.Assessment, error) {
	return nil, nil
}
func (m *mockStudentRepo) FindAssessmentByID(ctx context.Context, schoolID, assessmentID bson.ObjectID) (*model.Assessment, error) {
	return nil, nil
}
func (m *mockStudentRepo) FindQuestionsByAssessmentID(ctx context.Context, assessmentID bson.ObjectID) ([]model.AssessmentQuestion, error) {
	return nil, nil
}
func (m *mockStudentRepo) FindAttemptByStudent(ctx context.Context, schoolID, assessmentID, studentID bson.ObjectID) (*model.AssessmentAttempt, error) {
	return nil, nil
}
func (m *mockStudentRepo) InsertAttempt(ctx context.Context, attempt *model.AssessmentAttempt) error {
	return nil
}
func (m *mockStudentRepo) UpdateAttempt(ctx context.Context, attempt *model.AssessmentAttempt) error {
	return nil
}
func (m *mockStudentRepo) FindAnswersByAttemptID(ctx context.Context, attemptID bson.ObjectID) ([]model.AssessmentAnswer, error) {
	return nil, nil
}
func (m *mockStudentRepo) SaveAnswer(ctx context.Context, answer *model.AssessmentAnswer) error {
	return nil
}
func (m *mockStudentRepo) FindGradesByStudent(ctx context.Context, schoolID, studentID, academicYearID bson.ObjectID) ([]model.Grade, error) {
	if m.findGradesByStudentFn != nil {
		return m.findGradesByStudentFn(ctx, schoolID, studentID, academicYearID)
	}
	return nil, nil
}
func (m *mockStudentRepo) FindGradingConfig(ctx context.Context, schoolID, academicYearID bson.ObjectID, subjectID *bson.ObjectID) (*model.GradingConfig, error) {
	return &model.GradingConfig{
		DailyWeight: 60,
		UTSWeight:   20,
		UASWeight:   20,
	}, nil
}
func (m *mockStudentRepo) FindClassStudents(ctx context.Context, schoolID, classID bson.ObjectID) ([]model.Student, error) {
	if m.findClassStudentsFn != nil {
		return m.findClassStudentsFn(ctx, schoolID, classID)
	}
	return nil, nil
}
func (m *mockStudentRepo) FindSchedulesByClass(ctx context.Context, schoolID, classID, academicYearID bson.ObjectID) ([]model.Schedule, error) {
	return nil, nil
}
func (m *mockStudentRepo) FindSubjectsMap(ctx context.Context, subjectIDs []bson.ObjectID) (map[bson.ObjectID]*model.Subject, error) {
	if m.findSubjectsMapFn != nil {
		return m.findSubjectsMapFn(ctx, subjectIDs)
	}
	res := make(map[bson.ObjectID]*model.Subject)
	for _, id := range subjectIDs {
		res[id] = &model.Subject{ID: id, Name: "Matematika", Code: "MTK"}
	}
	return res, nil
}
func (m *mockStudentRepo) FindTeachersMap(ctx context.Context, teacherIDs []bson.ObjectID) (map[bson.ObjectID]string, error) {
	return make(map[bson.ObjectID]string), nil
}
func (m *mockStudentRepo) FindUsersMap(ctx context.Context, userIDs []bson.ObjectID) (map[bson.ObjectID]*model.User, error) {
	if m.findUsersMapFn != nil {
		return m.findUsersMapFn(ctx, userIDs)
	}
	return make(map[bson.ObjectID]*model.User), nil
}

func TestStudentService_GetGradesFormula(t *testing.T) {
	schoolID := bson.NewObjectID()
	studentID := bson.NewObjectID()
	classID := bson.NewObjectID()
	ayID := bson.NewObjectID()
	subjectID := bson.NewObjectID()

	scope := &middleware.StudentScope{
		SchoolID:       schoolID,
		StudentID:      studentID,
		ClassID:        classID,
		AcademicYearID: ayID,
		Semester:       1,
	}

	repo := &mockStudentRepo{
		findGradesByStudentFn: func(ctx context.Context, sID, stID, aID bson.ObjectID) ([]model.Grade, error) {
			return []model.Grade{
				{Type: model.GradeAssignment, Score: 80, SubjectID: subjectID},
				{Type: model.GradeAssessment, Score: 90, SubjectID: subjectID}, // average daily = 85
				{Type: model.GradeUTS, Score: 80, SubjectID: subjectID},        // UTS = 80
				{Type: model.GradeUAS, Score: 90, SubjectID: subjectID},        // UAS = 90
			}, nil
		},
	}

	svc := NewStudentService(repo, nil)
	grades, err := svc.GetGrades(context.Background(), scope)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if len(grades) != 1 {
		t.Fatalf("expected 1 subject grade, got: %d", len(grades))
	}

	sg := grades[0]
	// Formula: 85 * 0.6 (51) + 80 * 0.2 (16) + 90 * 0.2 (18) = 85.0
	expectedFinal := 85.0
	if sg.FinalScore == nil || *sg.FinalScore != expectedFinal {
		t.Fatalf("expected final score %f, got %v", expectedFinal, sg.FinalScore)
	}
	if sg.GradeLetter != "A" {
		t.Fatalf("expected letter A, got %s", sg.GradeLetter)
	}
}

func TestStudentService_GetRanking(t *testing.T) {
	schoolID := bson.NewObjectID()
	studentID1 := bson.NewObjectID()
	studentID2 := bson.NewObjectID()
	classID := bson.NewObjectID()
	ayID := bson.NewObjectID()

	scope := &middleware.StudentScope{
		SchoolID:       schoolID,
		StudentID:      studentID1,
		ClassID:        classID,
		AcademicYearID: ayID,
	}

	repo := &mockStudentRepo{
		findClassStudentsFn: func(ctx context.Context, sID, cID bson.ObjectID) ([]model.Student, error) {
			return []model.Student{
				{ID: studentID1, NIS: "1001"},
				{ID: studentID2, NIS: "1002"},
			}, nil
		},
		findGradesByStudentFn: func(ctx context.Context, sID, stID, aID bson.ObjectID) ([]model.Grade, error) {
			if stID == studentID1 {
				return []model.Grade{{Score: 85}}, nil
			}
			return []model.Grade{{Score: 95}}, nil
		},
	}

	svc := NewStudentService(repo, nil)
	ranking, err := svc.GetRanking(context.Background(), scope)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if ranking.Total != 2 {
		t.Fatalf("expected 2 students, got %d", ranking.Total)
	}
	// student2 (95) should be rank 1, student1 (85) should be rank 2
	if ranking.MyRank != 2 {
		t.Fatalf("expected myRank to be 2, got %d", ranking.MyRank)
	}
	if ranking.Rankings[0].StudentID != studentID2.Hex() {
		t.Fatalf("expected rank 1 to be student 2, got %s", ranking.Rankings[0].StudentID)
	}
}
