package teachersvc

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	teacherdto "github.com/geniussociety/backend/internal/dto/teacher"
	"github.com/geniussociety/backend/internal/handler"
	"github.com/geniussociety/backend/internal/middleware"
	"github.com/geniussociety/backend/internal/model"
	"github.com/geniussociety/backend/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type mockMaterialRepo struct {
	findMaterialsFunc func(ctx context.Context, teacherID bson.ObjectID, classID, subjectID *bson.ObjectID, search string, page, limit int64) ([]*model.Material, int64, error)
	findByIDFunc      func(ctx context.Context, teacherID, id bson.ObjectID) (*model.Material, error)
	createFunc        func(ctx context.Context, m *model.Material) error
	updateFunc        func(ctx context.Context, m *model.Material) error
	deleteFunc        func(ctx context.Context, teacherID, id bson.ObjectID) error
	classesMapFunc    func(ctx context.Context, classIDs []bson.ObjectID) (map[bson.ObjectID]*model.Class, error)
	subjectsMapFunc   func(ctx context.Context, subjectIDs []bson.ObjectID) (map[bson.ObjectID]*model.Subject, error)
}

func (m *mockMaterialRepo) FindMaterials(ctx context.Context, teacherID bson.ObjectID, classID, subjectID *bson.ObjectID, search string, page, limit int64) ([]*model.Material, int64, error) {
	if m.findMaterialsFunc != nil {
		return m.findMaterialsFunc(ctx, teacherID, classID, subjectID, search, page, limit)
	}
	return nil, 0, nil
}

func (m *mockMaterialRepo) FindMaterialByID(ctx context.Context, teacherID, id bson.ObjectID) (*model.Material, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, teacherID, id)
	}
	return nil, nil
}

func (m *mockMaterialRepo) CreateMaterial(ctx context.Context, mat *model.Material) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, mat)
	}
	return nil
}

func (m *mockMaterialRepo) UpdateMaterial(ctx context.Context, mat *model.Material) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, mat)
	}
	return nil
}

func (m *mockMaterialRepo) DeleteMaterial(ctx context.Context, teacherID, id bson.ObjectID) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, teacherID, id)
	}
	return nil
}

func (m *mockMaterialRepo) FindClassesMap(ctx context.Context, classIDs []bson.ObjectID) (map[bson.ObjectID]*model.Class, error) {
	if m.classesMapFunc != nil {
		return m.classesMapFunc(ctx, classIDs)
	}
	return make(map[bson.ObjectID]*model.Class), nil
}

func (m *mockMaterialRepo) FindSubjectsMap(ctx context.Context, subjectIDs []bson.ObjectID) (map[bson.ObjectID]*model.Subject, error) {
	if m.subjectsMapFunc != nil {
		return m.subjectsMapFunc(ctx, subjectIDs)
	}
	return make(map[bson.ObjectID]*model.Subject), nil
}

type mockFileStorage struct {
	uploadFunc     func(ctx context.Context, kind, filename, mimeType string, size int64, body io.Reader) (*model.FileRef, error)
	presignGetFunc func(ctx context.Context, key string, duration time.Duration) (string, error)
	deleteFunc     func(ctx context.Context, key string) error
}

func (m *mockFileStorage) Upload(ctx context.Context, kind, filename, mimeType string, size int64, body io.Reader) (*model.FileRef, error) {
	if m.uploadFunc != nil {
		return m.uploadFunc(ctx, kind, filename, mimeType, size, body)
	}
	return &model.FileRef{Key: "test-key", Name: filename, Size: size, Mime: mimeType}, nil
}

func (m *mockFileStorage) PresignGet(ctx context.Context, key string, duration time.Duration) (string, error) {
	if m.presignGetFunc != nil {
		return m.presignGetFunc(ctx, key, duration)
	}
	return "https://download.url/" + key, nil
}

func (m *mockFileStorage) Delete(ctx context.Context, key string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, key)
	}
	return nil
}

type mockNotifService struct {
	notifyClassesFunc func(ctx context.Context, classIDs []bson.ObjectID, notifType, title, body, refType string, refID bson.ObjectID) error
}

func (m *mockNotifService) GetNotifications(ctx context.Context, userID bson.ObjectID, page, limit int64) ([]teacherdto.NotificationDTO, *teacherdto.NotificationMeta, error) {
	return nil, nil, nil
}
func (m *mockNotifService) GetUnreadCount(ctx context.Context, userID bson.ObjectID) (*teacherdto.UnreadCountDTO, error) {
	return nil, nil
}
func (m *mockNotifService) MarkAsRead(ctx context.Context, userID, notifID bson.ObjectID) error {
	return nil
}
func (m *mockNotifService) MarkAllAsRead(ctx context.Context, userID bson.ObjectID) error {
	return nil
}
func (m *mockNotifService) NotifyClasses(ctx context.Context, classIDs []bson.ObjectID, notifType, title, body, refType string, refID bson.ObjectID) error {
	if m.notifyClassesFunc != nil {
		return m.notifyClassesFunc(ctx, classIDs, notifType, title, body, refType, refID)
	}
	return nil
}

func TestMaterialService_Create_ValidationAndScope(t *testing.T) {

	teacherID := bson.NewObjectID()
	classAllowed := bson.NewObjectID()
	classForbidden := bson.NewObjectID()
	subjectID := bson.NewObjectID()

	scope := &middleware.TeacherScope{

		TeacherID: teacherID,
		TeachingPairs: []repository.ClassSubjectPair{
			{ClassID: classAllowed, SubjectID: subjectID},
		},
	}

	repo := &mockMaterialRepo{}
	st := &mockFileStorage{}
	notif := &mockNotifService{}
	svc := NewMaterialService(repo, st, notif)

	// Test 1: Title too short
	_, err := svc.CreateMaterial(context.Background(), scope, &teacherdto.MaterialCreatePayload{
		Title:     "ab",
		SubjectID: subjectID.Hex(),
		ClassIDs:  []string{classAllowed.Hex()},
	}, nil)
	if err == nil {
		t.Fatal("expected validation error for short title")
	}

	// Test 2: Invalid Link (non http/https)
	_, err = svc.CreateMaterial(context.Background(), scope, &teacherdto.MaterialCreatePayload{
		Title:     "Materi Matematika",
		SubjectID: subjectID.Hex(),
		ClassIDs:  []string{classAllowed.Hex()},
		Link:      "javascript:alert(1)",
	}, nil)
	if err == nil {
		t.Fatal("expected validation error for invalid link protocol")
	}

	// Test 3: Unassigned class -> 403 NOT_ASSIGNED
	_, err = svc.CreateMaterial(context.Background(), scope, &teacherdto.MaterialCreatePayload{
		Title:     "Materi Matematika",
		SubjectID: subjectID.Hex(),
		ClassIDs:  []string{classForbidden.Hex()},
	}, nil)
	if !errors.Is(err, handler.ErrNotAssigned) {
		t.Fatalf("expected ErrNotAssigned for class outside scope, got: %v", err)
	}
}

func TestMaterialService_CreateAndPublish_TriggersNotification(t *testing.T) {

	teacherID := bson.NewObjectID()
	classID := bson.NewObjectID()
	subjectID := bson.NewObjectID()

	scope := &middleware.TeacherScope{

		TeacherID: teacherID,
		TeachingPairs: []repository.ClassSubjectPair{
			{ClassID: classID, SubjectID: subjectID},
		},
	}

	var createdMat *model.Material
	repo := &mockMaterialRepo{
		createFunc: func(ctx context.Context, m *model.Material) error {
			createdMat = m
			return nil
		},
		findByIDFunc: func(ctx context.Context, tID, id bson.ObjectID) (*model.Material, error) {
			return createdMat, nil
		},
		subjectsMapFunc: func(ctx context.Context, sIDs []bson.ObjectID) (map[bson.ObjectID]*model.Subject, error) {
			return map[bson.ObjectID]*model.Subject{
				subjectID: {ID: subjectID, Name: "Matematika"},
			}, nil
		},
		classesMapFunc: func(ctx context.Context, cIDs []bson.ObjectID) (map[bson.ObjectID]*model.Class, error) {
			return map[bson.ObjectID]*model.Class{
				classID: {ID: classID, Name: "X-IPA-1"},
			}, nil
		},
	}

	var notifSent bool
	notif := &mockNotifService{
		notifyClassesFunc: func(ctx context.Context, cIDs []bson.ObjectID, notifType, title, body, refType string, refID bson.ObjectID) error {
			notifSent = true
			if notifType != "new_material" {
				t.Errorf("expected notifType new_material, got %s", notifType)
			}
			return nil
		},
	}

	svc := NewMaterialService(repo, &mockFileStorage{}, notif)
	res, err := svc.CreateMaterial(context.Background(), scope, &teacherdto.MaterialCreatePayload{
		Title:      "Materi Aljabar Lanjut",
		SubjectID:  subjectID.Hex(),
		ClassIDs:   []string{classID.Hex()},
		PublishNow: true,
	}, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Title != "Materi Aljabar Lanjut" {
		t.Errorf("unexpected response title: %s", res.Title)
	}
	if !notifSent {
		t.Error("expected notification to be sent upon publish")
	}
}

func TestMaterialService_Delete_OrderDBFirstStorageAfter(t *testing.T) {

	teacherID := bson.NewObjectID()
	matID := bson.NewObjectID()

	scope := &middleware.TeacherScope{

		TeacherID: teacherID,
	}

	var callOrder []string
	repo := &mockMaterialRepo{
		findByIDFunc: func(ctx context.Context, tID, id bson.ObjectID) (*model.Material, error) {
			return &model.Material{
				ID: matID,

				TeacherID: teacherID,
				File:      &model.FileRef{Key: "materials/file1.pdf"},
			}, nil
		},
		deleteFunc: func(ctx context.Context, tID, id bson.ObjectID) error {
			callOrder = append(callOrder, "db_delete")
			return nil
		},
	}

	st := &mockFileStorage{
		deleteFunc: func(ctx context.Context, key string) error {
			callOrder = append(callOrder, "storage_delete")
			return nil
		},
	}

	svc := NewMaterialService(repo, st, nil)
	err := svc.DeleteMaterial(context.Background(), scope, matID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(callOrder) != 2 || callOrder[0] != "db_delete" || callOrder[1] != "storage_delete" {
		t.Errorf("expected order [db_delete, storage_delete], got %v", callOrder)
	}
}
