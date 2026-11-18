package teacherrepo

import (
	"context"

	"github.com/geniussociety/backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// WalasRepository handles persistence queries for homeroom teacher operations.
type WalasRepository interface {
	FindStudentsByClassID(ctx context.Context, schoolID, classID bson.ObjectID) ([]*model.Student, error)
	FindUsersByIDs(ctx context.Context, userIDs []bson.ObjectID) (map[bson.ObjectID]*model.User, error)
	FindGradesByClass(ctx context.Context, schoolID, classID, academicYearID bson.ObjectID) ([]*model.Grade, error)
	FindSubjectsBySchool(ctx context.Context, schoolID bson.ObjectID) ([]*model.Subject, error)
	FindSubjectsByIDs(ctx context.Context, ids []bson.ObjectID) (map[bson.ObjectID]*model.Subject, error)
}

type walasRepoImpl struct {
	students *mongo.Collection
	users    *mongo.Collection
	grades   *mongo.Collection
	subjects *mongo.Collection
}

// NewWalasRepository creates a new WalasRepository instance.
func NewWalasRepository(db *mongo.Database) WalasRepository {
	return &walasRepoImpl{
		students: db.Collection("students"),
		users:    db.Collection("users"),
		grades:   db.Collection("grades"),
		subjects: db.Collection("subjects"),
	}
}

func (r *walasRepoImpl) FindStudentsByClassID(ctx context.Context, schoolID, classID bson.ObjectID) ([]*model.Student, error) {
	filter := bson.M{
		"schoolId": schoolID,
		"classId":  classID,
		"status":   "active",
	}
	opts := options.Find().SetSort(bson.D{{Key: "nis", Value: 1}})
	cursor, err := r.students.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var items []*model.Student
	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *walasRepoImpl) FindUsersByIDs(ctx context.Context, userIDs []bson.ObjectID) (map[bson.ObjectID]*model.User, error) {
	result := make(map[bson.ObjectID]*model.User)
	if len(userIDs) == 0 {
		return result, nil
	}
	cursor, err := r.users.Find(ctx, bson.M{"_id": bson.M{"$in": userIDs}})
	if err != nil {
		return nil, err
	}
	var users []*model.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	for _, u := range users {
		result[u.ID] = u
	}
	return result, nil
}

func (r *walasRepoImpl) FindGradesByClass(ctx context.Context, schoolID, classID, academicYearID bson.ObjectID) ([]*model.Grade, error) {
	filter := bson.M{
		"schoolId":       schoolID,
		"classId":        classID,
		"academicYearId": academicYearID,
	}
	cursor, err := r.grades.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var items []*model.Grade
	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *walasRepoImpl) FindSubjectsBySchool(ctx context.Context, schoolID bson.ObjectID) ([]*model.Subject, error) {
	cursor, err := r.subjects.Find(ctx, bson.M{"schoolId": schoolID})
	if err != nil {
		return nil, err
	}
	var items []*model.Subject
	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *walasRepoImpl) FindSubjectsByIDs(ctx context.Context, ids []bson.ObjectID) (map[bson.ObjectID]*model.Subject, error) {
	result := make(map[bson.ObjectID]*model.Subject)
	if len(ids) == 0 {
		return result, nil
	}
	cursor, err := r.subjects.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	var items []*model.Subject
	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	for _, s := range items {
		result[s.ID] = s
	}
	return result, nil
}
