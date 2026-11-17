package teacher

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

type mockProjectRepo struct {
	findProjectsFunc       func(ctx context.Context, schoolID, teacherID bson.ObjectID, classID, subjectID *bson.ObjectID, search string, page, limit int64) ([]*model.Project, int64, error)
	findByIDFunc           func(ctx context.Context, schoolID, teacherID, id bson.ObjectID) (*model.Project, error)
	createFunc             func(ctx context.Context, p *model.Project) error
	updateFunc             func(ctx context.Context, p *model.Project) error
	deleteFunc             func(ctx context.Context, schoolID, teacherID, id bson.ObjectID) error
	countSubmissionsFunc   func(ctx context.Context, schoolID, projectID bson.ObjectID) (int64, error)
	findSubmissionsFunc    func(ctx context.Context, schoolID, projectID bson.ObjectID) ([]*model.ProjectSubmission, error)
	findSubmissionByIDFunc func(ctx context.Context, schoolID, submissionID bson.ObjectID) (*model.ProjectSubmission, error)
	findStudentsFunc       func(ctx context.Context, schoolID bson.ObjectID, classIDs []bson.ObjectID) ([]*model.Student, error)
	findStudentByIDFunc    func(ctx context.Context, schoolID, studentID bson.ObjectID) (*model.Student, error)
	findUsersMapFunc       func(ctx context.Context, userIDs []bson.ObjectID) (map[bson.ObjectID]*model.User, error)
	findClassesMapFunc     func(ctx context.Context, schoolID bson.ObjectID, classIDs []bson.ObjectID) (map[bson.ObjectID]*model.Class, error)
	findSubjectsMapFunc    func(ctx context.Context, subjectIDs []bson.ObjectID) (map[bson.ObjectID]*model.Subject, error)
	gradeInTxFunc          func(ctx context.Context, schoolID, teacherID bson.ObjectID, submission *model.ProjectSubmission, grade *model.Grade, audit *model.AuditLog) error
}

func (m *mockProjectRepo) FindProjects(ctx context.Context, schoolID, teacherID bson.ObjectID, classID, subjectID *bson.ObjectID, search string, page, limit int64) ([]*model.Project, int64, error) {
	if m.findProjectsFunc != nil {
		return m.findProjectsFunc(ctx, schoolID, teacherID, classID, subjectID, search, page, limit)
	}
	return nil, 0, nil
}

func (m *mockProjectRepo) FindProjectByID(ctx context.Context, schoolID, teacherID, id bson.ObjectID) (*model.Project, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, schoolID, teacherID, id)
	}
	return nil, nil
}

func (m *mockProjectRepo) CreateProject(ctx context.Context, p *model.Project) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, p)
	}
	return nil
}

func (m *mockProjectRepo) UpdateProject(ctx context.Context, p *model.Project) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, p)
	}
	return nil
}

func (m *mockProjectRepo) DeleteProject(ctx context.Context, schoolID, teacherID, id bson.ObjectID) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, schoolID, teacherID, id)
	}
	return nil
}

func (m *mockProjectRepo) CountSubmissionsByProjectID(ctx context.Context, schoolID, projectID bson.ObjectID) (int64, error) {
	if m.countSubmissionsFunc != nil {
		return m.countSubmissionsFunc(ctx, schoolID, projectID)
	}
	return 0, nil
}

func (m *mockProjectRepo) FindSubmissions(ctx context.Context, schoolID, projectID bson.ObjectID) ([]*model.ProjectSubmission, error) {
	if m.findSubmissionsFunc != nil {
		return m.findSubmissionsFunc(ctx, schoolID, projectID)
	}
	return nil, nil
}

func (m *mockProjectRepo) FindSubmissionByID(ctx context.Context, schoolID, submissionID bson.ObjectID) (*model.ProjectSubmission, error) {
	if m.findSubmissionByIDFunc != nil {
		return m.findSubmissionByIDFunc(ctx, schoolID, submissionID)
	}
	return nil, nil
}

func (m *mockProjectRepo) FindStudentsByClassIDs(ctx context.Context, schoolID bson.ObjectID, classIDs []bson.ObjectID) ([]*model.Student, error) {
	if m.findStudentsFunc != nil {
		return m.findStudentsFunc(ctx, schoolID, classIDs)
	}
	return nil, nil
}

func (m *mockProjectRepo) FindStudentByID(ctx context.Context, schoolID, studentID bson.ObjectID) (*model.Student, error) {
	if m.findStudentByIDFunc != nil {
		return m.findStudentByIDFunc(ctx, schoolID, studentID)
	}
	return nil, nil
}

func (m *mockProjectRepo) FindUsersMap(ctx context.Context, userIDs []bson.ObjectID) (map[bson.ObjectID]*model.User, error) {
	if m.findUsersMapFunc != nil {
		return m.findUsersMapFunc(ctx, userIDs)
	}
	return make(map[bson.ObjectID]*model.User), nil
}

func (m *mockProjectRepo) FindClassesMap(ctx context.Context, schoolID bson.ObjectID, classIDs []bson.ObjectID) (map[bson.ObjectID]*model.Class, error) {
	if m.findClassesMapFunc != nil {
		return m.findClassesMapFunc(ctx, schoolID, classIDs)
	}
	return make(map[bson.ObjectID]*model.Class), nil
}

func (m *mockProjectRepo) FindSubjectsMap(ctx context.Context, subjectIDs []bson.ObjectID) (map[bson.ObjectID]*model.Subject, error) {
	if m.findSubjectsMapFunc != nil {
		return m.findSubjectsMapFunc(ctx, subjectIDs)
	}
	return make(map[bson.ObjectID]*model.Subject), nil
}

func (m *mockProjectRepo) GradeSubmissionInTransaction(ctx context.Context, schoolID, teacherID bson.ObjectID, submission *model.ProjectSubmission, grade *model.Grade, audit *model.AuditLog) error {
	if m.gradeInTxFunc != nil {
		return m.gradeInTxFunc(ctx, schoolID, teacherID, submission, grade, audit)
	}
	return nil
}

func TestProjectService_Delete_ConflictWhenHasSubmissions(t *testing.T) {
	schoolID := bson.NewObjectID()
	teacherID := bson.NewObjectID()
	prjID := bson.NewObjectID()

	scope := &middleware.TeacherScope{
		SchoolID:  schoolID,
		TeacherID: teacherID,
	}

	repo := &mockProjectRepo{
		findByIDFunc: func(ctx context.Context, sID, tID, id bson.ObjectID) (*model.Project, error) {
			return &model.Project{ID: prjID, SchoolID: schoolID, TeacherID: teacherID}, nil
		},
		countSubmissionsFunc: func(ctx context.Context, sID, pID bson.ObjectID) (int64, error) {
			return 2, nil // Has 2 submissions!
		},
	}

	svc := NewProjectService(repo, &mockFileStorage{}, nil, nil)
	err := svc.DeleteProject(context.Background(), scope, prjID)
	if err == nil {
		t.Fatal("expected 409 conflict when deleting project with submissions")
	}

	var dErr *handler.DomainError
	if !errors.As(err, &dErr) || dErr.Status != 409 {
		t.Errorf("expected DomainError with status 409, got: %v", err)
	}
}

func TestProjectService_GradeSubmission_Success(t *testing.T) {
	schoolID := bson.NewObjectID()
	teacherID := bson.NewObjectID()
	prjID := bson.NewObjectID()
	subID := bson.NewObjectID()
	studentID := bson.NewObjectID()
	classID := bson.NewObjectID()
	subjectID := bson.NewObjectID()
	academicYearID := bson.NewObjectID()

	scope := &middleware.TeacherScope{
		SchoolID:       schoolID,
		TeacherID:      teacherID,
		AcademicYearID: academicYearID,
		Semester:       1,
	}

	var gradedTxCalled bool
	var savedGrade *model.Grade
	repo := &mockProjectRepo{
		findByIDFunc: func(ctx context.Context, sID, tID, id bson.ObjectID) (*model.Project, error) {
			return &model.Project{
				ID:        prjID,
				SchoolID:  schoolID,
				TeacherID: teacherID,
				SubjectID: subjectID,
				Title:     "Projek Akhir",
			}, nil
		},
		findSubmissionByIDFunc: func(ctx context.Context, sID, id bson.ObjectID) (*model.ProjectSubmission, error) {
			return &model.ProjectSubmission{
				ID:        subID,
				ProjectID: prjID,
				StudentID: studentID,
				SchoolID:  schoolID,
				Version:   2,
			}, nil
		},
		findStudentByIDFunc: func(ctx context.Context, sID, id bson.ObjectID) (*model.Student, error) {
			return &model.Student{
				ID:       studentID,
				ClassID:  classID,
				SchoolID: schoolID,
				UserID:   bson.NewObjectID(),
			}, nil
		},
		gradeInTxFunc: func(ctx context.Context, sID, tID bson.ObjectID, submission *model.ProjectSubmission, grade *model.Grade, audit *model.AuditLog) error {
			gradedTxCalled = true
			savedGrade = grade
			if *submission.Score != 92.5 {
				t.Errorf("expected score 92.5 on submission, got %f", *submission.Score)
			}
			return nil
		},
	}

	svc := NewProjectService(repo, &mockFileStorage{}, nil, nil)
	err := svc.GradeSubmission(context.Background(), scope, prjID, subID, 92.5, "Projek sangat kreatif!")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !gradedTxCalled {
		t.Fatal("expected GradeSubmissionInTransaction to be called")
	}
	if savedGrade.Score != 92.5 || savedGrade.Type != model.GradeProject || savedGrade.Source != "submission" {
		t.Errorf("unexpected grade details: %+v", savedGrade)
	}
}

func TestProjectService_ListSubmissions_VersionTracking(t *testing.T) {
	schoolID := bson.NewObjectID()
	teacherID := bson.NewObjectID()
	prjID := bson.NewObjectID()
	classID := bson.NewObjectID()
	st1 := bson.NewObjectID()

	scope := &middleware.TeacherScope{
		SchoolID:  schoolID,
		TeacherID: teacherID,
	}

	deadline := time.Now().Add(1 * time.Hour)
	v1Score := 70.0
	v2Score := 95.0
	now := time.Now()

	repo := &mockProjectRepo{
		findByIDFunc: func(ctx context.Context, sID, tID, id bson.ObjectID) (*model.Project, error) {
			return &model.Project{
				ID:        prjID,
				SchoolID:  schoolID,
				TeacherID: teacherID,
				ClassIDs:  []bson.ObjectID{classID},
				Deadline:  deadline,
			}, nil
		},
		findStudentsFunc: func(ctx context.Context, sID bson.ObjectID, cIDs []bson.ObjectID) ([]*model.Student, error) {
			return []*model.Student{
				{ID: st1, ClassID: classID, NIS: "101"},
			}, nil
		},
		findSubmissionsFunc: func(ctx context.Context, sID, pID bson.ObjectID) ([]*model.ProjectSubmission, error) {
			return []*model.ProjectSubmission{
				{
					ID:          bson.NewObjectID(),
					ProjectID:   prjID,
					StudentID:   st1,
					Version:     1,
					SubmittedAt: now.Add(-2 * time.Hour),
					Score:       &v1Score,
				},
				{
					ID:          bson.NewObjectID(),
					ProjectID:   prjID,
					StudentID:   st1,
					Version:     2,
					SubmittedAt: now.Add(-10 * time.Minute),
					Score:       &v2Score,
					GradedAt:    &now,
				},
			}, nil
		},
	}

	svc := NewProjectService(repo, &mockFileStorage{}, nil, nil)
	list, _, err := svc.ListSubmissions(context.Background(), scope, prjID, "", "", "", 1, 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 submission (latest version 2), got %d", len(list))
	}
	if list[0].Version != 2 || *list[0].Score != 95.0 {
		t.Errorf("expected latest version 2 with score 95.0, got version %d score %v", list[0].Version, list[0].Score)
	}
}
