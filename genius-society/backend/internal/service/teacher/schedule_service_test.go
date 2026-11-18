package teachersvc

import (
	"context"
	"testing"

	"github.com/geniussociety/backend/internal/middleware"
	"github.com/geniussociety/backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type mockScheduleRepo struct {
	schedules []*model.Schedule
	subjects  map[bson.ObjectID]*model.Subject
	classes   map[bson.ObjectID]*model.Class
}

func (m *mockScheduleRepo) FindTeacherSchedules(ctx context.Context, teacherID, academicYearID bson.ObjectID) ([]*model.Schedule, error) {
	return m.schedules, nil
}

func (m *mockScheduleRepo) FindSubjectsByIDs(ctx context.Context, ids []bson.ObjectID) (map[bson.ObjectID]*model.Subject, error) {
	return m.subjects, nil
}

func (m *mockScheduleRepo) FindClassesByIDs(ctx context.Context, ids []bson.ObjectID) (map[bson.ObjectID]*model.Class, error) {
	return m.classes, nil
}

func TestScheduleService_GetTeacherSchedule(t *testing.T) {
	teacherID := bson.NewObjectID()
	ayID := bson.NewObjectID()
	classID := bson.NewObjectID()
	subjectID := bson.NewObjectID()

	mockRepo := &mockScheduleRepo{
		schedules: []*model.Schedule{
			{
				ID:             bson.NewObjectID(),
				ClassID:        classID,
				SubjectID:      subjectID,
				TeacherID:      teacherID,
				AcademicYearID: ayID,
				Day:            1,
				StartTime:      "07:30",
				EndTime:        "09:00",
				Room:           "R101",
			},
		},
		subjects: map[bson.ObjectID]*model.Subject{
			subjectID: {ID: subjectID, Name: "Matematika Wajib"},
		},
		classes: map[bson.ObjectID]*model.Class{
			classID: {ID: classID, Name: "X-IPA-1"},
		},
	}

	svc := NewScheduleService(mockRepo)
	scope := &middleware.TeacherScope{
		TeacherID:      teacherID,
		AcademicYearID: ayID,
	}

	items, err := svc.GetTeacherSchedule(context.Background(), scope)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].Subject != "Matematika Wajib" {
		t.Errorf("expected subject Matematika Wajib, got %s", items[0].Subject)
	}
	if items[0].ClassName != "X-IPA-1" {
		t.Errorf("expected class X-IPA-1, got %s", items[0].ClassName)
	}
}
