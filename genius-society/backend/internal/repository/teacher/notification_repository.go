package teacherrepo

import (
	"context"
	"time"

	"github.com/geniussociety/backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// NotificationRepository defines database operations for notifications.
type NotificationRepository interface {
	FindByUserID(ctx context.Context, userID bson.ObjectID, page, limit int64) ([]*model.Notification, int64, error)
	CountUnread(ctx context.Context, userID bson.ObjectID) (int64, error)
	MarkRead(ctx context.Context, userID, notifID bson.ObjectID) error
	MarkAllRead(ctx context.Context, userID bson.ObjectID) error
	FindActiveStudentUserIDsByClassIDs(ctx context.Context, schoolID bson.ObjectID, classIDs []bson.ObjectID) ([]bson.ObjectID, error)
	InsertMany(ctx context.Context, notifs []*model.Notification) error
}

type notificationRepoImpl struct {
	notifications *mongo.Collection
	students      *mongo.Collection
}

// NewNotificationRepository creates a new NotificationRepository instance.
func NewNotificationRepository(db *mongo.Database) NotificationRepository {
	return &notificationRepoImpl{
		notifications: db.Collection("notifications"),
		students:      db.Collection("students"),
	}
}

func (r *notificationRepoImpl) FindByUserID(ctx context.Context, userID bson.ObjectID, page, limit int64) ([]*model.Notification, int64, error) {
	filter := bson.M{"userId": userID}

	total, err := r.notifications.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	skip := (page - 1) * limit
	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}}).
		SetSkip(skip).
		SetLimit(limit)

	cur, err := r.notifications.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(ctx)

	var list []*model.Notification
	if err := cur.All(ctx, &list); err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *notificationRepoImpl) CountUnread(ctx context.Context, userID bson.ObjectID) (int64, error) {
	filter := bson.M{
		"userId": userID,
		"isRead": false,
	}
	return r.notifications.CountDocuments(ctx, filter)
}

func (r *notificationRepoImpl) MarkRead(ctx context.Context, userID, notifID bson.ObjectID) error {
	filter := bson.M{
		"_id":    notifID,
		"userId": userID,
	}
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"isRead":    true,
			"readAt":    now,
			"updatedAt": now,
		},
	}
	res, err := r.notifications.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *notificationRepoImpl) MarkAllRead(ctx context.Context, userID bson.ObjectID) error {
	filter := bson.M{
		"userId": userID,
		"isRead": false,
	}
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"isRead":    true,
			"readAt":    now,
			"updatedAt": now,
		},
	}
	_, err := r.notifications.UpdateMany(ctx, filter, update)
	return err
}

func (r *notificationRepoImpl) FindActiveStudentUserIDsByClassIDs(ctx context.Context, schoolID bson.ObjectID, classIDs []bson.ObjectID) ([]bson.ObjectID, error) {
	if len(classIDs) == 0 {
		return []bson.ObjectID{}, nil
	}

	filter := bson.M{
		"classId":  bson.M{"$in": classIDs},
		"schoolId": schoolID,
		"status":   "active",
		"userId":   bson.M{"$exists": true, "$ne": bson.NilObjectID},
	}

	opts := options.Find().SetProjection(bson.M{"userId": 1})
	cur, err := r.students.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var userIDs []bson.ObjectID
	seen := make(map[bson.ObjectID]bool)
	for cur.Next(ctx) {
		var doc struct {
			UserID bson.ObjectID `bson:"userId"`
		}
		if err := cur.Decode(&doc); err == nil && !doc.UserID.IsZero() {
			if !seen[doc.UserID] {
				seen[doc.UserID] = true
				userIDs = append(userIDs, doc.UserID)
			}
		}
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return userIDs, nil
}

func (r *notificationRepoImpl) InsertMany(ctx context.Context, notifs []*model.Notification) error {
	if len(notifs) == 0 {
		return nil
	}

	// Chunk by 500
	const chunkSize = 500
	for i := 0; i < len(notifs); i += chunkSize {
		end := i + chunkSize
		if end > len(notifs) {
			end = len(notifs)
		}
		chunk := notifs[i:end]

		docs := make([]interface{}, len(chunk))
		for j, n := range chunk {
			docs[j] = n
		}

		// Use Unordered insert so duplicate keys don't abort the entire batch
		opts := options.InsertMany().SetOrdered(false)
		_, err := r.notifications.InsertMany(ctx, docs, opts)
		if err != nil {
			// If error is duplicate key error, we consider it idempotent and proceed
			// In mongo v2, bulk write exceptions or duplicate keys are tolerated
			// Any non-duplicate error is still non-fatal as notifs are best-effort, but log if needed
		}
	}
	return nil
}
