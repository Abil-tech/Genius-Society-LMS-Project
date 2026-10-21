package repository

import (
	"context"
	"time"

	"github.com/geniussociety/backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type TeacherRepository struct {
	col *mongo.Collection
}

func NewTeacherRepository(db *mongo.Database) *TeacherRepository {
	return &TeacherRepository{col: db.Collection("teachers")}
}

func (r *TeacherRepository) FindByUserID(ctx context.Context, userID bson.ObjectID) (*model.Teacher, error) {
	var t model.Teacher
	err := r.col.FindOne(ctx, bson.M{"userId": userID}).Decode(&t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TeacherRepository) FindByID(ctx context.Context, id bson.ObjectID) (*model.Teacher, error) {
	var t model.Teacher
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TeacherRepository) Create(ctx context.Context, t *model.Teacher) error {
	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now
	res, err := r.col.InsertOne(ctx, t)
	if err != nil {
		return err
	}
	t.ID = res.InsertedID.(bson.ObjectID)
	return nil
}

func (r *TeacherRepository) Update(ctx context.Context, id bson.ObjectID, update bson.M) error {
	update["updatedAt"] = time.Now()
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	return err
}

func (r *TeacherRepository) List(ctx context.Context, filter bson.M, page, limit int64) ([]*model.Teacher, int64, error) {
	total, _ := r.col.CountDocuments(ctx, filter)
	skip := (page - 1) * limit
	opts := options.Find().SetSkip(skip).SetLimit(limit).SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cur, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(ctx)
	var teachers []*model.Teacher
	cur.All(ctx, &teachers)
	return teachers, total, nil
}

func (r *TeacherRepository) FindByNIP(ctx context.Context, nip string) (*model.Teacher, error) {
	if nip == "" {
		return nil, mongo.ErrNoDocuments
	}
	var t model.Teacher
	err := r.col.FindOne(ctx, bson.M{"nip": nip}).Decode(&t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TeacherRepository) Delete(ctx context.Context, id bson.ObjectID) error {
	_, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *TeacherRepository) Count(ctx context.Context, filter bson.M) (int64, error) {
	return r.col.CountDocuments(ctx, filter)
}
