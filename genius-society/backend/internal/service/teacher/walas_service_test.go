package teachersvc

import (
	"context"
	"testing"

	"github.com/geniussociety/backend/internal/middleware"
	"github.com/geniussociety/backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type mockWalasRepo struct {
	students []*model.Student
	users    map[bson.ObjectID]*model.User
	grades   []*model.Grade
	subjects []*model.Subject
}

func (m *mockWalasRepo) FindStudentsByClassID(ctx context.Context, schoolID, classID bson.ObjectID) ([]*model.Student, error) {
	return m.students, nil
}
func (m *mockWalasRepo) FindUsersByIDs(ctx context.Context, userIDs []bson.ObjectID) (map[bson.ObjectID]*model.User, error) {
	return m.users, nil
}
func (m *mockWalasRepo) FindGradesByClass(ctx context.Context, schoolID, classID, academicYearID bson.ObjectID) ([]*model.Grade, error) {
	return m.grades, nil
}
func (m *mockWalasRepo) FindSubjectsBySchool(ctx context.Context, schoolID bson.ObjectID) ([]*model.Subject, error) {
	return m.subjects, nil
}
func (m *mockWalasRepo) FindSubjectsByIDs(ctx context.Context, ids []bson.ObjectID) (map[bson.ObjectID]*model.Subject, error) {
	result := make(map[bson.ObjectID]*model.Subject)
	for _, s := range m.subjects {
		result[s.ID] = s
	}
	return result, nil
}

func TestWalasService_GetWalasStudents(t *testing.T) {
	schoolID := bson.NewObjectID()
	classID := bson.NewObjectID()
	studentID := bson.NewObjectID()
	userID := bson.NewObjectID()

	mockRepo := &mockWalasRepo{
		students: []*model.Student{
			{ID: studentID, UserID: userID, NIS: "1001", Gender: "L"},
		},
		users: map[bson.ObjectID]*model.User{
			userID: {ID: userID, Name: "Ahmad Fauzi"},
		},
		grades: []*model.Grade{
			{StudentID: studentID, Score: 85},
			{StudentID: studentID, Score: 95},
		},
	}

	svc := NewWalasService(mockRepo)
	scope := &middleware.TeacherScope{
		SchoolID:      schoolID,
		IsWalas:       true,
		WalasClassIDs: []bson.ObjectID{classID},
	}

	students, err := svc.GetWalasStudents(context.Background(), scope)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(students) != 1 {
		t.Fatalf("expected 1 student, got %d", len(students))
	}
	st := students[0]
	if st.Name != "Ahmad Fauzi" {
		t.Errorf("expected name Ahmad Fauzi, got %s", st.Name)
	}
	if st.AverageScore == nil || *st.AverageScore != 90.0 {
		t.Errorf("expected averageScore 90.0, got %v", st.AverageScore)
	}
	if st.Rank != 1 {
		t.Errorf("expected rank 1, got %d", st.Rank)
	}
}
