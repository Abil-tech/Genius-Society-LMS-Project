package repository

import (
	"context"
	"time"

	"github.com/geniussociety/backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type AcademicDecisionRepository struct {
	coll *mongo.Collection
}

func NewAcademicDecisionRepository(db *mongo.Database) *AcademicDecisionRepository {
	return &AcademicDecisionRepository{
		coll: db.Collection("academic_decisions"),
	}
}

func (r *AcademicDecisionRepository) Create(ctx context.Context, d *model.AcademicDecision) error {
	now := time.Now()
	d.CreatedAt = now
	d.UpdatedAt = now
	res, err := r.coll.InsertOne(ctx, d)
	if err != nil {
		return err
	}
	d.ID = res.InsertedID.(bson.ObjectID)
	return nil
}

func (r *AcademicDecisionRepository) BulkCreate(ctx context.Context, decisions []*model.AcademicDecision) error {
	if len(decisions) == 0 {
		return nil
	}
	docs := make([]interface{}, len(decisions))
	now := time.Now()
	for i, d := range decisions {
		d.CreatedAt = now
		d.UpdatedAt = now
		docs[i] = d
	}
	_, err := r.coll.InsertMany(ctx, docs)
	return err
}

func (r *AcademicDecisionRepository) FindByID(ctx context.Context, schoolID, id bson.ObjectID) (*model.AcademicDecision, error) {
	var d model.AcademicDecision
	filter := bson.M{"_id": id, "schoolId": schoolID}
	if err := r.coll.FindOne(ctx, filter).Decode(&d); err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *AcademicDecisionRepository) FindSuspendedByYear(ctx context.Context, schoolID, academicYearID bson.ObjectID) ([]*model.AcademicDecision, error) {
	filter := bson.M{
		"schoolId":             schoolID,
		"sourceAcademicYearId": academicYearID,
		"decisionStatus":       model.DecisionStatusSuspended,
		"reviewStatus":         bson.M{"$ne": model.ReviewStatusResolved},
	}
	cur, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var list []*model.AcademicDecision
	if err := cur.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *AcademicDecisionRepository) FindByStudentAndYear(ctx context.Context, schoolID, studentID, sourceYearID bson.ObjectID) (*model.AcademicDecision, error) {
	var d model.AcademicDecision
	filter := bson.M{
		"schoolId":             schoolID,
		"studentId":            studentID,
		"sourceAcademicYearId": sourceYearID,
	}
	opts := options.FindOne().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	if err := r.coll.FindOne(ctx, filter, opts).Decode(&d); err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *AcademicDecisionRepository) Update(ctx context.Context, schoolID, id bson.ObjectID, update bson.M) error {
	update["updatedAt"] = time.Now()
	filter := bson.M{"_id": id, "schoolId": schoolID}
	_, err := r.coll.UpdateOne(ctx, filter, bson.M{"$set": update})
	return err
}
