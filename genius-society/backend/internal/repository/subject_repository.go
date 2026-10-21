package repository

import (
	"context"
	"time"

	"github.com/geniussociety/backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type SubjectRepository struct {
	col *mongo.Collection
}

func NewSubjectRepository(db *mongo.Database) *SubjectRepository {
	return &SubjectRepository{col: db.Collection("subjects")}
}

func (r *SubjectRepository) FindByID(ctx context.Context, id bson.ObjectID) (*model.Subject, error) {
	var s model.Subject
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SubjectRepository) FindByCode(ctx context.Context, code string) (*model.Subject, error) {
	if code == "" {
		return nil, mongo.ErrNoDocuments
	}
	var s model.Subject
	err := r.col.FindOne(ctx, bson.M{"code": code}).Decode(&s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SubjectRepository) Create(ctx context.Context, s *model.Subject) error {
	now := time.Now()
	s.CreatedAt = now
	s.UpdatedAt = now
	res, err := r.col.InsertOne(ctx, s)
	if err != nil {
		return err
	}
	s.ID = res.InsertedID.(bson.ObjectID)
	return nil
}

func (r *SubjectRepository) Update(ctx context.Context, id bson.ObjectID, update bson.M) error {
	update["updatedAt"] = time.Now()
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	return err
}

func (r *SubjectRepository) Delete(ctx context.Context, id bson.ObjectID) error {
	_, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *SubjectRepository) List(ctx context.Context, filter bson.M, page, limit int64) ([]*model.Subject, int64, error) {
	total, _ := r.col.CountDocuments(ctx, filter)
	skip := (page - 1) * limit
	opts := options.Find().SetSkip(skip).SetLimit(limit).SetSort(bson.D{{Key: "code", Value: 1}})
	cur, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(ctx)
	var subjects []*model.Subject
	cur.All(ctx, &subjects)
	return subjects, total, nil
}

func (r *SubjectRepository) Count(ctx context.Context, filter bson.M) (int64, error) {
	return r.col.CountDocuments(ctx, filter)
}
