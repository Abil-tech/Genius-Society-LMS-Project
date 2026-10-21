package repository

import (
	"context"
	"time"

	"github.com/geniussociety/backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type StudentRepository struct {
	col *mongo.Collection
}

func NewStudentRepository(db *mongo.Database) *StudentRepository {
	return &StudentRepository{col: db.Collection("students")}
}

func (r *StudentRepository) FindByUserID(ctx context.Context, userID bson.ObjectID) (*model.Student, error) {
	var s model.Student
	err := r.col.FindOne(ctx, bson.M{"userId": userID}).Decode(&s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *StudentRepository) FindByID(ctx context.Context, id bson.ObjectID) (*model.Student, error) {
	var s model.Student
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *StudentRepository) Create(ctx context.Context, s *model.Student) error {
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

func (r *StudentRepository) Update(ctx context.Context, id bson.ObjectID, update bson.M) error {
	update["updatedAt"] = time.Now()
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	return err
}

func (r *StudentRepository) List(ctx context.Context, filter bson.M, page, limit int64) ([]*model.Student, int64, error) {
	total, _ := r.col.CountDocuments(ctx, filter)
	skip := (page - 1) * limit
	opts := options.Find().SetSkip(skip).SetLimit(limit).SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cur, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(ctx)
	var students []*model.Student
	cur.All(ctx, &students)
	return students, total, nil
}

func (r *StudentRepository) FindByNIS(ctx context.Context, nis string) (*model.Student, error) {
	var s model.Student
	err := r.col.FindOne(ctx, bson.M{"nis": nis}).Decode(&s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *StudentRepository) FindByNISN(ctx context.Context, nisn string) (*model.Student, error) {
	if nisn == "" {
		return nil, mongo.ErrNoDocuments
	}
	var s model.Student
	err := r.col.FindOne(ctx, bson.M{"nisn": nisn}).Decode(&s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *StudentRepository) Delete(ctx context.Context, id bson.ObjectID) error {
	_, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *StudentRepository) Count(ctx context.Context, filter bson.M) (int64, error) {
	return r.col.CountDocuments(ctx, filter)
}

func (r *StudentRepository) BulkUpdateClass(ctx context.Context, studentIDs []bson.ObjectID, classID bson.ObjectID) error {
	filter := bson.M{"_id": bson.M{"$in": studentIDs}}
	update := bson.M{
		"$set": bson.M{
			"classId":   classID,
			"updatedAt": time.Now(),
		},
	}
	_, err := r.col.UpdateMany(ctx, filter, update)
	return err
}

func (r *StudentRepository) BulkUpdateStatus(ctx context.Context, studentIDs []bson.ObjectID, status string) error {
	filter := bson.M{"_id": bson.M{"$in": studentIDs}}
	update := bson.M{
		"$set": bson.M{
			"status":    status,
			"updatedAt": time.Now(),
		},
	}
	_, err := r.col.UpdateMany(ctx, filter, update)
	return err
}
