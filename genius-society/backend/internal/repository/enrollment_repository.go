package repository

import (
	"context"
	"time"

	"github.com/geniussociety/backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type EnrollmentRepository struct {
	coll *mongo.Collection
}

func NewEnrollmentRepository(db *mongo.Database) *EnrollmentRepository {
	return &EnrollmentRepository{
		coll: db.Collection("student_enrollments"),
	}
}

func (r *EnrollmentRepository) Create(ctx context.Context, e *model.StudentEnrollment) error {
	now := time.Now()
	e.CreatedAt = now
	e.UpdatedAt = now
	res, err := r.coll.InsertOne(ctx, e)
	if err != nil {
		return err
	}
	e.ID = res.InsertedID.(bson.ObjectID)
	return nil
}

func (r *EnrollmentRepository) BulkCreate(ctx context.Context, enrollments []*model.StudentEnrollment) error {
	if len(enrollments) == 0 {
		return nil
	}
	docs := make([]interface{}, len(enrollments))
	now := time.Now()
	for i, e := range enrollments {
		e.CreatedAt = now
		e.UpdatedAt = now
		docs[i] = e
	}
	_, err := r.coll.InsertMany(ctx, docs)
	return err
}

func (r *EnrollmentRepository) FindByStudentAndYear(ctx context.Context, schoolID, studentID, academicYearID bson.ObjectID) (*model.StudentEnrollment, error) {
	var e model.StudentEnrollment
	filter := bson.M{
		"schoolId":       schoolID,
		"studentId":      studentID,
		"academicYearId": academicYearID,
	}
	if err := r.coll.FindOne(ctx, filter).Decode(&e); err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *EnrollmentRepository) FindByClass(ctx context.Context, schoolID, classID bson.ObjectID) ([]*model.StudentEnrollment, error) {
	filter := bson.M{
		"schoolId": schoolID,
		"classId":  classID,
	}
	cur, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var list []*model.StudentEnrollment
	if err := cur.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *EnrollmentRepository) FindByStudent(ctx context.Context, schoolID, studentID bson.ObjectID) ([]*model.StudentEnrollment, error) {
	filter := bson.M{
		"schoolId":  schoolID,
		"studentId": studentID,
	}
	cur, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var list []*model.StudentEnrollment
	if err := cur.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *EnrollmentRepository) UpdateStatus(ctx context.Context, schoolID, id bson.ObjectID, status model.EnrollmentStatus) error {
	filter := bson.M{"_id": id, "schoolId": schoolID}
	update := bson.M{"$set": bson.M{"status": status, "updatedAt": time.Now()}}
	_, err := r.coll.UpdateOne(ctx, filter, update)
	return err
}
