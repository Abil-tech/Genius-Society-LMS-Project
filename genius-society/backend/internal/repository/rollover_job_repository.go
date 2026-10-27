package repository

import (
	"context"
	"time"

	"github.com/geniussociety/backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type RolloverJobRepository struct {
	coll *mongo.Collection
}

func NewRolloverJobRepository(db *mongo.Database) *RolloverJobRepository {
	return &RolloverJobRepository{
		coll: db.Collection("rollover_jobs"),
	}
}

func (r *RolloverJobRepository) Create(ctx context.Context, job *model.RolloverJob) error {
	now := time.Now()
	job.CreatedAt = now
	job.UpdatedAt = now
	res, err := r.coll.InsertOne(ctx, job)
	if err != nil {
		return err
	}
	job.ID = res.InsertedID.(bson.ObjectID)
	return nil
}

func (r *RolloverJobRepository) FindByID(ctx context.Context, schoolID, id bson.ObjectID) (*model.RolloverJob, error) {
	var job model.RolloverJob
	filter := bson.M{"_id": id, "schoolId": schoolID}
	if err := r.coll.FindOne(ctx, filter).Decode(&job); err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *RolloverJobRepository) FindByIdempotencyKey(ctx context.Context, schoolID bson.ObjectID, key string) (*model.RolloverJob, error) {
	var job model.RolloverJob
	filter := bson.M{"schoolId": schoolID, "idempotencyKey": key}
	if err := r.coll.FindOne(ctx, filter).Decode(&job); err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *RolloverJobRepository) FindByYearAndType(ctx context.Context, schoolID, academicYearID bson.ObjectID, processType model.ProcessType) ([]*model.RolloverJob, error) {
	filter := bson.M{
		"schoolId":             schoolID,
		"sourceAcademicYearId": academicYearID,
		"processType":          processType,
	}
	cur, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var list []*model.RolloverJob
	if err := cur.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *RolloverJobRepository) Update(ctx context.Context, schoolID, id bson.ObjectID, update bson.M) error {
	update["updatedAt"] = time.Now()
	filter := bson.M{"_id": id, "schoolId": schoolID}
	_, err := r.coll.UpdateOne(ctx, filter, bson.M{"$set": update})
	return err
}
