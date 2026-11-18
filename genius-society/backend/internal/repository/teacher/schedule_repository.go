package teacherrepo

import (
	"context"

	"github.com/geniussociety/backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// ScheduleRepository handles database operations for teacher schedule queries.
type ScheduleRepository interface {
	FindTeacherSchedules(ctx context.Context, teacherID, academicYearID bson.ObjectID) ([]*model.Schedule, error)
	FindSubjectsByIDs(ctx context.Context, ids []bson.ObjectID) (map[bson.ObjectID]*model.Subject, error)
	FindClassesByIDs(ctx context.Context, ids []bson.ObjectID) (map[bson.ObjectID]*model.Class, error)
}

type scheduleRepoImpl struct {
	schedules *mongo.Collection
	subjects  *mongo.Collection
	classes   *mongo.Collection
}

// NewScheduleRepository creates a new ScheduleRepository instance.
func NewScheduleRepository(db *mongo.Database) ScheduleRepository {
	return &scheduleRepoImpl{
		schedules: db.Collection("schedules"),
		subjects:  db.Collection("subjects"),
		classes:   db.Collection("classes"),
	}
}

func (r *scheduleRepoImpl) FindTeacherSchedules(ctx context.Context, teacherID, academicYearID bson.ObjectID) ([]*model.Schedule, error) {
	filter := bson.M{
		"teacherId":      teacherID,
		"academicYearId": academicYearID,
	}
	opts := options.Find().SetSort(bson.D{{Key: "day", Value: 1}, {Key: "startTime", Value: 1}})
	cursor, err := r.schedules.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var items []*model.Schedule
	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *scheduleRepoImpl) FindSubjectsByIDs(ctx context.Context, ids []bson.ObjectID) (map[bson.ObjectID]*model.Subject, error) {
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

func (r *scheduleRepoImpl) FindClassesByIDs(ctx context.Context, ids []bson.ObjectID) (map[bson.ObjectID]*model.Class, error) {
	result := make(map[bson.ObjectID]*model.Class)
	if len(ids) == 0 {
		return result, nil
	}
	cursor, err := r.classes.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	var items []*model.Class
	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	for _, c := range items {
		result[c.ID] = c
	}
	return result, nil
}
