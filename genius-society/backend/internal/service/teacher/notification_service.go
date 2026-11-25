package teachersvc

import (
	"context"
	"errors"
	"math"
	"time"

	teacherdto "github.com/geniussociety/backend/internal/dto/teacher"
	"github.com/geniussociety/backend/internal/handler"
	"github.com/geniussociety/backend/internal/model"
	teacherrepo "github.com/geniussociety/backend/internal/repository/teacher"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// NotificationService defines business logic for notifications.
type NotificationService interface {
	GetNotifications(ctx context.Context, userID bson.ObjectID, page, limit int64) ([]teacherdto.NotificationDTO, *teacherdto.NotificationMeta, error)
	GetUnreadCount(ctx context.Context, userID bson.ObjectID) (*teacherdto.UnreadCountDTO, error)
	MarkAsRead(ctx context.Context, userID, notifID bson.ObjectID) error
	MarkAllAsRead(ctx context.Context, userID bson.ObjectID) error
	NotifyClasses(ctx context.Context, classIDs []bson.ObjectID, notifType, title, body, refType string, refID bson.ObjectID) error
}

type notificationServiceImpl struct {
	repo teacherrepo.NotificationRepository
}

// NewNotificationService creates a new NotificationService instance.
func NewNotificationService(repo teacherrepo.NotificationRepository) NotificationService {
	return &notificationServiceImpl{repo: repo}
}

func (s *notificationServiceImpl) GetNotifications(ctx context.Context, userID bson.ObjectID, page, limit int64) ([]teacherdto.NotificationDTO, *teacherdto.NotificationMeta, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if page <= 0 {
		page = 1
	}

	list, total, err := s.repo.FindByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, nil, err
	}

	dtos := make([]teacherdto.NotificationDTO, 0, len(list))
	for _, n := range list {
		refIDStr := ""
		if n.RefID != nil && !n.RefID.IsZero() {
			refIDStr = n.RefID.Hex()
		}
		dtos = append(dtos, teacherdto.NotificationDTO{
			ID:        n.ID.Hex(),
			UserID:    n.UserID.Hex(),
			Type:      n.Type,
			Title:     n.Title,
			Body:      n.Body,
			RefType:   n.RefType,
			RefID:     refIDStr,
			IsRead:    n.IsRead,
			ReadAt:    n.ReadAt,
			CreatedAt: n.CreatedAt,
		})
	}

	totalPages := int64(0)
	if total > 0 {
		totalPages = int64(math.Ceil(float64(total) / float64(limit)))
	}

	meta := &teacherdto.NotificationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	return dtos, meta, nil
}

func (s *notificationServiceImpl) GetUnreadCount(ctx context.Context, userID bson.ObjectID) (*teacherdto.UnreadCountDTO, error) {
	count, err := s.repo.CountUnread(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &teacherdto.UnreadCountDTO{Count: count}, nil
}

func (s *notificationServiceImpl) MarkAsRead(ctx context.Context, userID, notifID bson.ObjectID) error {
	err := s.repo.MarkRead(ctx, userID, notifID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return handler.ErrNotFound
		}
		return err
	}
	return nil
}

func (s *notificationServiceImpl) MarkAllAsRead(ctx context.Context, userID bson.ObjectID) error {
	return s.repo.MarkAllRead(ctx, userID)
}

func (s *notificationServiceImpl) NotifyClasses(ctx context.Context, classIDs []bson.ObjectID, notifType, title, body, refType string, refID bson.ObjectID) error {
	if len(classIDs) == 0 {
		return nil
	}

	studentUserIDs, err := s.repo.FindActiveStudentUserIDsByClassIDs(ctx, classIDs)
	if err != nil {
		return err
	}

	if len(studentUserIDs) == 0 {
		return nil
	}

	now := time.Now().UTC()
	notifs := make([]*model.Notification, 0, len(studentUserIDs))
	for _, uid := range studentUserIDs {
		notifs = append(notifs, &model.Notification{
			ID:        bson.NewObjectID(),
			UserID:    uid,
			Type:      notifType,
			Title:     title,
			Body:      body,
			RefType:   refType,
			RefID:     &refID,
			IsRead:    false,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}

	return s.repo.InsertMany(ctx, notifs)
}
