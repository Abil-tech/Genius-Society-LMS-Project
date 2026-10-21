package repository

import (
	"context"
	"time"

	"github.com/geniussociety/backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ClassRepository struct {
	col *mongo.Collection
}

func NewClassRepository(db *mongo.Database) *ClassRepository {
	return &ClassRepository{col: db.Collection("classes")}
}

func (r *ClassRepository) FindByID(ctx context.Context, id bson.ObjectID) (*model.Class, error) {
	var c model.Class
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ClassRepository) Create(ctx context.Context, c *model.Class) error {
	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now
	res, err := r.col.InsertOne(ctx, c)
	if err != nil {
		return err
	}
	c.ID = res.InsertedID.(bson.ObjectID)
	return nil
}

func (r *ClassRepository) List(ctx context.Context, filter bson.M, page, limit int64) ([]*model.Class, int64, error) {
	total, _ := r.col.CountDocuments(ctx, filter)
	skip := (page - 1) * limit
	opts := options.Find().SetSkip(skip).SetLimit(limit).SetSort(bson.D{{Key: "name", Value: 1}})
	cur, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(ctx)

	var classes []*model.Class
	if err := cur.All(ctx, &classes); err != nil {
		return nil, 0, err
	}
	return classes, total, nil
}

func (r *ClassRepository) Count(ctx context.Context, filter bson.M) (int64, error) {
	return r.col.CountDocuments(ctx, filter)
}
