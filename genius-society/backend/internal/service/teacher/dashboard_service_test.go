package teachersvc

import (
	"context"
	"testing"

	"github.com/geniussociety/backend/internal/middleware"
	"github.com/geniussociety/backend/internal/model"
	"github.com/geniussociety/backend/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type mockDashboardRepo struct {
	studentsCount       int
	activeAssignments   int
	activeAssessments   int
	activeProjects      int
	ungradedAssignments int
	ungradedProjects    int
	ungradedAssessments int
	schedules           []*model.Schedule
	recentAssignments   []*model.AssignmentSubmission
	recentProjects      []*model.ProjectSubmission
	assignmentsMap      map[bson.ObjectID]*model.Assignment
	projectsMap         map[bson.ObjectID]*model.Project
	studentsMap         map[bson.ObjectID]*model.Student
	subjectsMap         map[bson.ObjectID]*model.Subject
	classesMap          map[bson.ObjectID]*model.Class
	walasClass          *model.Class
}

func (m *mockDashboardRepo) CountStudentsByClassIDs(ctx context.Context, classIDs []bson.ObjectID) (int, error) {
	return m.studentsCount, nil
}
func (m *mockDashboardRepo) CountActiveAssignments(ctx context.Context, teacherID bson.ObjectID) (int, error) {
	return m.activeAssignments, nil
}
func (m *mockDashboardRepo) CountActiveAssessments(ctx context.Context, teacherID bson.ObjectID) (int, error) {
	return m.activeAssessments, nil
}
func (m *mockDashboardRepo) CountActiveProjects(ctx context.Context, teacherID bson.ObjectID) (int, error) {
	return m.activeProjects, nil
}
func (m *mockDashboardRepo) CountUngradedAssignmentSubmissions(ctx context.Context, teacherID bson.ObjectID) (int, error) {
	return m.ungradedAssignments, nil
}
func (m *mockDashboardRepo) CountUngradedProjectSubmissions(ctx context.Context, teacherID bson.ObjectID) (int, error) {
	return m.ungradedProjects, nil
}
func (m *mockDashboardRepo) CountUngradedAssessmentAttempts(ctx context.Context, teacherID bson.ObjectID) (int, error) {
	return m.ungradedAssessments, nil
}
func (m *mockDashboardRepo) FindTodaySchedules(ctx context.Context, teacherID, academicYearID bson.ObjectID, dayOfWeek int) ([]*model.Schedule, error) {
	return m.schedules, nil
}
func (m *mockDashboardRepo) FindRecentAssignmentSubmissions(ctx context.Context, teacherID bson.ObjectID, limit int) ([]*model.AssignmentSubmission, error) {
	return m.recentAssignments, nil
}
func (m *mockDashboardRepo) FindRecentProjectSubmissions(ctx context.Context, teacherID bson.ObjectID, limit int) ([]*model.ProjectSubmission, error) {
	return m.recentProjects, nil
}
func (m *mockDashboardRepo) FindAssignmentsByIDs(ctx context.Context, ids []bson.ObjectID) (map[bson.ObjectID]*model.Assignment, error) {
	return m.assignmentsMap, nil
}
func (m *mockDashboardRepo) FindProjectsByIDs(ctx context.Context, ids []bson.ObjectID) (map[bson.ObjectID]*model.Project, error) {
	return m.projectsMap, nil
}
func (m *mockDashboardRepo) FindStudentsByIDs(ctx context.Context, ids []bson.ObjectID) (map[bson.ObjectID]*model.Student, error) {
	return m.studentsMap, nil
}
func (m *mockDashboardRepo) FindClassByID(ctx context.Context, classID bson.ObjectID) (*model.Class, error) {
	return m.walasClass, nil
}
func (m *mockDashboardRepo) FindSubjectsByIDs(ctx context.Context, ids []bson.ObjectID) (map[bson.ObjectID]*model.Subject, error) {
	return m.subjectsMap, nil
}
func (m *mockDashboardRepo) FindClassesByIDs(ctx context.Context, ids []bson.ObjectID) (map[bson.ObjectID]*model.Class, error) {
	return m.classesMap, nil
}
func (m *mockDashboardRepo) CountStudentsByClassID(ctx context.Context, classID bson.ObjectID) (int, error) {
	return m.studentsCount, nil
}

func TestDashboardService_GetDashboardData(t *testing.T) {
	teacherID := bson.NewObjectID()
	classID := bson.NewObjectID()
	subjectID := bson.NewObjectID()
	ayID := bson.NewObjectID()

	mockDash := &mockDashboardRepo{
		studentsCount:       30,
		activeAssignments:   2,
		activeAssessments:   1,
		activeProjects:      1,
		ungradedAssignments: 3,
		ungradedProjects:    2,
		ungradedAssessments: 1,
		schedules: []*model.Schedule{
			{
				ID:        bson.NewObjectID(),
				ClassID:   classID,
				SubjectID: subjectID,
				TeacherID: teacherID,
				Day:       1,
				StartTime: "07:30",
				EndTime:   "09:00",
				Room:      "R101",
			},
		},
		subjectsMap: map[bson.ObjectID]*model.Subject{
			subjectID: {ID: subjectID, Name: "Matematika"},
		},
		classesMap: map[bson.ObjectID]*model.Class{
			classID: {ID: classID, Name: "X-IPA-1"},
		},
		walasClass: &model.Class{
			ID:             classID,
			Name:           "X-IPA-1",
			EducationLevel: "SMA",
		},
	}

	mockClass := &mockClassRepo{
		usersFunc: func(ctx context.Context, userIDs []bson.ObjectID) (map[bson.ObjectID]*model.User, error) {
			return map[bson.ObjectID]*model.User{}, nil
		},
	}

	svc := NewDashboardService(mockDash, mockClass)

	scope := &middleware.TeacherScope{
		TeacherID:      teacherID,
		AcademicYearID: ayID,
		IsWalas:        true,
		WalasClassIDs:  []bson.ObjectID{classID},
		TeachingPairs: []repository.ClassSubjectPair{
			{ClassID: classID, SubjectID: subjectID},
		},
	}

	resp, err := svc.GetDashboardData(context.Background(), scope)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if resp.KPI.ClassesCount != 1 {
		t.Errorf("expected ClassesCount 1, got %d", resp.KPI.ClassesCount)
	}
	if resp.KPI.StudentsCount != 30 {
		t.Errorf("expected StudentsCount 30, got %d", resp.KPI.StudentsCount)
	}
	if resp.KPI.PendingGradesCount != 6 {
		t.Errorf("expected PendingGradesCount 6, got %d", resp.KPI.PendingGradesCount)
	}
	if resp.KPI.ActiveAssignmentsCount != 2 {
		t.Errorf("expected ActiveAssignmentsCount 2, got %d", resp.KPI.ActiveAssignmentsCount)
	}
	if !resp.IsWalas {
		t.Errorf("expected IsWalas true, got false")
	}
	if resp.WalasClass == nil || resp.WalasClass.Name != "X-IPA-1" {
		t.Errorf("expected WalasClass name X-IPA-1, got %v", resp.WalasClass)
	}
}
