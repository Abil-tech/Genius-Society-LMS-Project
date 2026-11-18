package teachersvc

import (
	"context"
	"testing"
	"time"

	"github.com/geniussociety/backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type mockNotificationRepo struct {
	findByUserIDFunc func(ctx context.Context, userID bson.ObjectID, page, limit int64) ([]*model.Notification, int64, error)
	countUnreadFunc  func(ctx context.Context, userID bson.ObjectID) (int64, error)
	markReadFunc     func(ctx context.Context, userID, notifID bson.ObjectID) error
	markAllReadFunc  func(ctx context.Context, userID bson.ObjectID) error
	findStudentsFunc func(ctx context.Context, schoolID bson.ObjectID, classIDs []bson.ObjectID) ([]bson.ObjectID, error)
	insertManyFunc   func(ctx context.Context, notifs []*model.Notification) error
}

func (m *mockNotificationRepo) FindByUserID(ctx context.Context, userID bson.ObjectID, page, limit int64) ([]*model.Notification, int64, error) {
	if m.findByUserIDFunc != nil {
		return m.findByUserIDFunc(ctx, userID, page, limit)
	}
	return nil, 0, nil
}

func (m *mockNotificationRepo) CountUnread(ctx context.Context, userID bson.ObjectID) (int64, error) {
	if m.countUnreadFunc != nil {
		return m.countUnreadFunc(ctx, userID)
	}
	return 0, nil
}

func (m *mockNotificationRepo) MarkRead(ctx context.Context, userID, notifID bson.ObjectID) error {
	if m.markReadFunc != nil {
		return m.markReadFunc(ctx, userID, notifID)
	}
	return nil
}

func (m *mockNotificationRepo) MarkAllRead(ctx context.Context, userID bson.ObjectID) error {
	if m.markAllReadFunc != nil {
		return m.markAllReadFunc(ctx, userID)
	}
	return nil
}

func (m *mockNotificationRepo) FindActiveStudentUserIDsByClassIDs(ctx context.Context, schoolID bson.ObjectID, classIDs []bson.ObjectID) ([]bson.ObjectID, error) {
	if m.findStudentsFunc != nil {
		return m.findStudentsFunc(ctx, schoolID, classIDs)
	}
	return nil, nil
}

func (m *mockNotificationRepo) InsertMany(ctx context.Context, notifs []*model.Notification) error {
	if m.insertManyFunc != nil {
		return m.insertManyFunc(ctx, notifs)
	}
	return nil
}

func TestNotificationService_GetNotifications(t *testing.T) {
	userID := bson.NewObjectID()
	refID := bson.NewObjectID()
	now := time.Now()

	repo := &mockNotificationRepo{
		findByUserIDFunc: func(ctx context.Context, uid bson.ObjectID, page, limit int64) ([]*model.Notification, int64, error) {
			if uid != userID {
				t.Errorf("expected userID %s, got %s", userID.Hex(), uid.Hex())
			}
			return []*model.Notification{
				{
					ID:        bson.NewObjectID(),
					UserID:    userID,
					Type:      "new_material",
					Title:     "Materi Baru",
					Body:      "Materi Aljabar telah ditambahkan",
					RefType:   "material",
					RefID:     &refID,
					IsRead:    false,
					CreatedAt: now,
				},
			}, 1, nil
		},
	}

	svc := NewNotificationService(repo)
	notifs, meta, err := svc.GetNotifications(context.Background(), userID, 1, 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(notifs) != 1 {
		t.Fatalf("expected 1 notification, got %d", len(notifs))
	}
	if notifs[0].Title != "Materi Baru" || notifs[0].RefID != refID.Hex() {
		t.Errorf("unexpected notif content: %+v", notifs[0])
	}
	if meta.Total != 1 || meta.TotalPages != 1 {
		t.Errorf("unexpected meta: %+v", meta)
	}
}

func TestNotificationService_GetUnreadCount(t *testing.T) {
	userID := bson.NewObjectID()
	repo := &mockNotificationRepo{
		countUnreadFunc: func(ctx context.Context, uid bson.ObjectID) (int64, error) {
			return 5, nil
		},
	}

	svc := NewNotificationService(repo)
	res, err := svc.GetUnreadCount(context.Background(), userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Count != 5 {
		t.Errorf("expected count 5, got %d", res.Count)
	}
}

func TestNotificationService_MarkAsRead_NotFound(t *testing.T) {
	userID := bson.NewObjectID()
	notifID := bson.NewObjectID()
	repo := &mockNotificationRepo{
		markReadFunc: func(ctx context.Context, uid, nid bson.ObjectID) error {
			return mongo.ErrNoDocuments
		},
	}

	svc := NewNotificationService(repo)
	err := svc.MarkAsRead(context.Background(), userID, notifID)
	if err == nil {
		t.Fatal("expected error for non-existent notification, got nil")
	}
}

func TestNotificationService_NotifyClasses(t *testing.T) {
	schoolID := bson.NewObjectID()
	class1 := bson.NewObjectID()
	studentUser1 := bson.NewObjectID()
	studentUser2 := bson.NewObjectID()
	refID := bson.NewObjectID()

	var insertedCount int
	repo := &mockNotificationRepo{
		findStudentsFunc: func(ctx context.Context, sID bson.ObjectID, cIDs []bson.ObjectID) ([]bson.ObjectID, error) {
			return []bson.ObjectID{studentUser1, studentUser2}, nil
		},
		insertManyFunc: func(ctx context.Context, notifs []*model.Notification) error {
			insertedCount = len(notifs)
			return nil
		},
	}

	svc := NewNotificationService(repo)
	err := svc.NotifyClasses(context.Background(), schoolID, []bson.ObjectID{class1}, "new_assignment", "Tugas Baru", "Kerjakan tugas 1", "assignment", refID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if insertedCount != 2 {
		t.Errorf("expected 2 inserted notifications, got %d", insertedCount)
	}
}
