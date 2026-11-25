package teachersvc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/geniussociety/backend/internal/handler"
	"github.com/geniussociety/backend/internal/middleware"
	"github.com/geniussociety/backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type mockAssignmentRepo struct {
	findAssignmentsFunc    func(ctx context.Context, teacherID bson.ObjectID, classID, subjectID *bson.ObjectID, search string, page, limit int64) ([]*model.Assignment, int64, error)
	findByIDFunc           func(ctx context.Context, teacherID, id bson.ObjectID) (*model.Assignment, error)
	createFunc             func(ctx context.Context, a *model.Assignment) error
	updateFunc             func(ctx context.Context, a *model.Assignment) error
	deleteFunc             func(ctx context.Context, teacherID, id bson.ObjectID) error
	countSubmissionsFunc   func(ctx context.Context, assignmentID bson.ObjectID) (int64, error)
	findSubmissionsFunc    func(ctx context.Context, assignmentID bson.ObjectID) ([]*model.AssignmentSubmission, error)
	findSubmissionByIDFunc func(ctx context.Context, submissionID bson.ObjectID) (*model.AssignmentSubmission, error)
	findStudentsFunc       func(ctx context.Context, classIDs []bson.ObjectID) ([]*model.Student, error)
	findStudentByIDFunc    func(ctx context.Context, studentID bson.ObjectID) (*model.Student, error)
	findUsersMapFunc       func(ctx context.Context, userIDs []bson.ObjectID) (map[bson.ObjectID]*model.User, error)
	findClassesMapFunc     func(ctx context.Context, classIDs []bson.ObjectID) (map[bson.ObjectID]*model.Class, error)
	findSubjectsMapFunc    func(ctx context.Context, subjectIDs []bson.ObjectID) (map[bson.ObjectID]*model.Subject, error)
	gradeInTxFunc          func(ctx context.Context, teacherID bson.ObjectID, submission *model.AssignmentSubmission, grade *model.Grade, audit *model.AuditLog) error
}

func (m *mockAssignmentRepo) FindAssignments(ctx context.Context, teacherID bson.ObjectID, classID, subjectID *bson.ObjectID, search string, page, limit int64) ([]*model.Assignment, int64, error) {
	if m.findAssignmentsFunc != nil {
		return m.findAssignmentsFunc(ctx, teacherID, classID, subjectID, search, page, limit)
	}
	return nil, 0, nil
}

func (m *mockAssignmentRepo) FindAssignmentByID(ctx context.Context, teacherID, id bson.ObjectID) (*model.Assignment, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, teacherID, id)
	}
	return nil, nil
}

func (m *mockAssignmentRepo) CreateAssignment(ctx context.Context, a *model.Assignment) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, a)
	}
	return nil
}

func (m *mockAssignmentRepo) UpdateAssignment(ctx context.Context, a *model.Assignment) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, a)
	}
	return nil
}

func (m *mockAssignmentRepo) DeleteAssignment(ctx context.Context, teacherID, id bson.ObjectID) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, teacherID, id)
	}
	return nil
}

func (m *mockAssignmentRepo) CountSubmissionsByAssignmentID(ctx context.Context, assignmentID bson.ObjectID) (int64, error) {
	if m.countSubmissionsFunc != nil {
		return m.countSubmissionsFunc(ctx, assignmentID)
	}
	return 0, nil
}

func (m *mockAssignmentRepo) FindSubmissions(ctx context.Context, assignmentID bson.ObjectID) ([]*model.AssignmentSubmission, error) {
	if m.findSubmissionsFunc != nil {
		return m.findSubmissionsFunc(ctx, assignmentID)
	}
	return nil, nil
}

func (m *mockAssignmentRepo) FindSubmissionByID(ctx context.Context, submissionID bson.ObjectID) (*model.AssignmentSubmission, error) {
	if m.findSubmissionByIDFunc != nil {
		return m.findSubmissionByIDFunc(ctx, submissionID)
	}
	return nil, nil
}

func (m *mockAssignmentRepo) FindStudentsByClassIDs(ctx context.Context, classIDs []bson.ObjectID) ([]*model.Student, error) {
	if m.findStudentsFunc != nil {
		return m.findStudentsFunc(ctx, classIDs)
	}
	return nil, nil
}

func (m *mockAssignmentRepo) FindStudentByID(ctx context.Context, studentID bson.ObjectID) (*model.Student, error) {
	if m.findStudentByIDFunc != nil {
		return m.findStudentByIDFunc(ctx, studentID)
	}
	return nil, nil
}

func (m *mockAssignmentRepo) FindUsersMap(ctx context.Context, userIDs []bson.ObjectID) (map[bson.ObjectID]*model.User, error) {
	if m.findUsersMapFunc != nil {
		return m.findUsersMapFunc(ctx, userIDs)
	}
	return make(map[bson.ObjectID]*model.User), nil
}

func (m *mockAssignmentRepo) FindClassesMap(ctx context.Context, classIDs []bson.ObjectID) (map[bson.ObjectID]*model.Class, error) {
	if m.findClassesMapFunc != nil {
		return m.findClassesMapFunc(ctx, classIDs)
	}
	return make(map[bson.ObjectID]*model.Class), nil
}

func (m *mockAssignmentRepo) FindSubjectsMap(ctx context.Context, subjectIDs []bson.ObjectID) (map[bson.ObjectID]*model.Subject, error) {
	if m.findSubjectsMapFunc != nil {
		return m.findSubjectsMapFunc(ctx, subjectIDs)
	}
	return make(map[bson.ObjectID]*model.Subject), nil
}

func (m *mockAssignmentRepo) GradeSubmissionInTransaction(ctx context.Context, teacherID bson.ObjectID, submission *model.AssignmentSubmission, grade *model.Grade, audit *model.AuditLog) error {
	if m.gradeInTxFunc != nil {
		return m.gradeInTxFunc(ctx, teacherID, submission, grade, audit)
	}
	return nil
}

func TestAssignmentService_Delete_ConflictWhenHasSubmissions(t *testing.T) {
	teacherID := bson.NewObjectID()
	asgID := bson.NewObjectID()

	scope := &middleware.TeacherScope{
		TeacherID: teacherID,
	}

	repo := &mockAssignmentRepo{
		findByIDFunc: func(ctx context.Context, tID, id bson.ObjectID) (*model.Assignment, error) {
			return &model.Assignment{ID: asgID, TeacherID: teacherID}, nil
		},
		countSubmissionsFunc: func(ctx context.Context, aID bson.ObjectID) (int64, error) {
			return 3, nil // Has 3 submissions!
		},
	}

	svc := NewAssignmentService(repo, &mockFileStorage{}, nil, nil)
	err := svc.DeleteAssignment(context.Background(), scope, asgID)
	if err == nil {
		t.Fatal("expected 409 conflict when deleting assignment with submissions")
	}

	var dErr *handler.DomainError
	if !errors.As(err, &dErr) || dErr.Status != 409 {
		t.Errorf("expected DomainError with status 409, got: %v", err)
	}
}

func TestAssignmentService_GradeSubmission_Success(t *testing.T) {

	teacherID := bson.NewObjectID()
	asgID := bson.NewObjectID()
	subID := bson.NewObjectID()
	studentID := bson.NewObjectID()
	classID := bson.NewObjectID()
	subjectID := bson.NewObjectID()
	academicYearID := bson.NewObjectID()

	scope := &middleware.TeacherScope{
		TeacherID:      teacherID,
		AcademicYearID: academicYearID,
		Semester:       1,
	}

	var gradedTxCalled bool
	var savedGrade *model.Grade
	repo := &mockAssignmentRepo{
		findByIDFunc: func(ctx context.Context, tID, id bson.ObjectID) (*model.Assignment, error) {
			return &model.Assignment{
				ID: asgID,

				TeacherID: teacherID,
				SubjectID: subjectID,
				Title:     "Tugas 1",
			}, nil
		},
		findSubmissionByIDFunc: func(ctx context.Context, id bson.ObjectID) (*model.AssignmentSubmission, error) {
			return &model.AssignmentSubmission{
				ID:           subID,
				AssignmentID: asgID,
				StudentID:    studentID,
				Version:      1,
			}, nil
		},
		findStudentByIDFunc: func(ctx context.Context, id bson.ObjectID) (*model.Student, error) {
			return &model.Student{
				ID:      studentID,
				ClassID: classID,
				UserID:  bson.NewObjectID(),
			}, nil
		},
		gradeInTxFunc: func(ctx context.Context, tID bson.ObjectID, submission *model.AssignmentSubmission, grade *model.Grade, audit *model.AuditLog) error {
			gradedTxCalled = true
			savedGrade = grade
			if *submission.Score != 85.0 {
				t.Errorf("expected score 85.0 on submission, got %f", *submission.Score)
			}
			return nil
		},
	}

	svc := NewAssignmentService(repo, &mockFileStorage{}, nil, nil)
	err := svc.GradeSubmission(context.Background(), scope, asgID, subID, 85.0, "Kerja bagus!")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !gradedTxCalled {
		t.Fatal("expected GradeSubmissionInTransaction to be called")
	}
	if savedGrade.Score != 85.0 || savedGrade.Type != model.GradeAssignment || savedGrade.Source != "submission" {
		t.Errorf("unexpected grade details: %+v", savedGrade)
	}
}

func TestAssignmentService_ListSubmissions_StatusFilter(t *testing.T) {

	teacherID := bson.NewObjectID()
	asgID := bson.NewObjectID()
	classID := bson.NewObjectID()
	st1 := bson.NewObjectID()
	st2 := bson.NewObjectID()

	scope := &middleware.TeacherScope{
		TeacherID: teacherID,
	}

	deadline := time.Now().Add(-1 * time.Hour) // deadline in the past
	score := 90.0

	repo := &mockAssignmentRepo{
		findByIDFunc: func(ctx context.Context, tID, id bson.ObjectID) (*model.Assignment, error) {
			return &model.Assignment{
				ID: asgID,

				TeacherID: teacherID,
				ClassIDs:  []bson.ObjectID{classID},
				Deadline:  deadline,
			}, nil
		},
		findStudentsFunc: func(ctx context.Context, cIDs []bson.ObjectID) ([]*model.Student, error) {
			return []*model.Student{
				{ID: st1, ClassID: classID, NIS: "101"},
				{ID: st2, ClassID: classID, NIS: "102"},
			}, nil
		},
		findSubmissionsFunc: func(ctx context.Context, aID bson.ObjectID) ([]*model.AssignmentSubmission, error) {
			now := time.Now()
			return []*model.AssignmentSubmission{
				{
					ID:           bson.NewObjectID(),
					AssignmentID: asgID,
					StudentID:    st1,
					Version:      1,
					SubmittedAt:  deadline.Add(-10 * time.Minute),
					Score:        &score,
					GradedAt:     &now,
				},
				// st2 has no submission -> missing!
			}, nil
		},
	}

	svc := NewAssignmentService(repo, &mockFileStorage{}, nil, nil)

	// Filter graded
	gradedList, _, err := svc.ListSubmissions(context.Background(), scope, asgID, "graded", "", "", 1, 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(gradedList) != 1 || gradedList[0].StudentID != st1.Hex() {
		t.Errorf("expected 1 graded submission (st1), got %d", len(gradedList))
	}

	// Filter missing
	missingList, _, err := svc.ListSubmissions(context.Background(), scope, asgID, "missing", "", "", 1, 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(missingList) != 1 || missingList[0].StudentID != st2.Hex() {
		t.Errorf("expected 1 missing submission (st2), got %d", len(missingList))
	}
}
