package repository

import (
	"context"
	"time"

	"github.com/geniussociety/backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type AuditLogRepository struct {
	coll *mongo.Collection
}

func NewAuditLogRepository(db *mongo.Database) *AuditLogRepository {
	return &AuditLogRepository{
		coll: db.Collection("audit_logs"),
	}
}

func (r *AuditLogRepository) Log(ctx context.Context, userID bson.ObjectID, action, entityType, entityID string, details bson.M) error {
	log := &model.AuditLog{
		
		UserID:     userID,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		Details:    details,
		Timestamp:  time.Now(),
	}
	res, err := r.coll.InsertOne(ctx, log)
	if err != nil {
		return err
	}
	log.ID = res.InsertedID.(bson.ObjectID)
	return nil
}

func (r *AuditLogRepository) ListBySchool(ctx context.Context, limit int64) ([]*model.AuditLog, error) {
	filter := bson.M{}
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}}).SetLimit(limit)
	cur, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var list []*model.AuditLog
	if err := cur.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}
