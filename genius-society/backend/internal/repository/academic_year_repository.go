// File: d:/projectME/Genius-Society-Baru/genius-society/backend/internal/repository/academic_year_repository.go
package repository

import (
	"context"
	"time"

	"github.com/geniussociety/backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type AcademicYearRepository struct {
	coll *mongo.Collection
}

func NewAcademicYearRepository(db *mongo.Database) *AcademicYearRepository {
	return &AcademicYearRepository{
		coll: db.Collection("academic_years"),
	}
}

// ---------- CREATE ----------
func (r *AcademicYearRepository) Create(ctx context.Context, ay *model.AcademicYear) error {
	now := time.Now()
	ay.CreatedAt = now
	ay.UpdatedAt = now
	if ay.Status == "" {
		if ay.IsActive {
			ay.Status = "active"
		} else {
			ay.Status = "inactive"
		}
	}

	res, err := r.coll.InsertOne(ctx, ay)
	if err != nil {
		return err
	}
	ay.ID = res.InsertedID.(bson.ObjectID)
	return nil
}

// ---------- FIND ----------
func (r *AcademicYearRepository) FindByID(ctx context.Context, schoolID, id bson.ObjectID) (*model.AcademicYear, error) {
	var ay model.AcademicYear
	filter := bson.M{
		"_id":      id,
		"schoolId": schoolID,
	}
	if err := r.coll.FindOne(ctx, filter).Decode(&ay); err != nil {
		return nil, err
	}
	return &ay, nil
}

func (r *AcademicYearRepository) FindAllBySchool(ctx context.Context, schoolID bson.ObjectID) ([]*model.AcademicYear, error) {
	filter := bson.M{"schoolId": schoolID}
	cur, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "name", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var result []*model.AcademicYear
	if err = cur.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *AcademicYearRepository) FindActiveBySchool(ctx context.Context, schoolID bson.ObjectID) (*model.AcademicYear, error) {
	var ay model.AcademicYear
	filter := bson.M{
		"schoolId": schoolID,
		"isActive": true,
	}
	if err := r.coll.FindOne(ctx, filter).Decode(&ay); err != nil {
		return nil, err
	}
	return &ay, nil
}

func (r *AcademicYearRepository) FindByName(ctx context.Context, schoolID bson.ObjectID, name string) (*model.AcademicYear, error) {
	var ay model.AcademicYear
	filter := bson.M{
		"schoolId": schoolID,
		"name":     name,
	}
	if err := r.coll.FindOne(ctx, filter).Decode(&ay); err != nil {
		return nil, err
	}
	return &ay, nil
}

// ---------- UPDATE ----------
func (r *AcademicYearRepository) Update(ctx context.Context, schoolID, id bson.ObjectID, update bson.M) error {
	update["updatedAt"] = time.Now()
	filter := bson.M{
		"_id":      id,
		"schoolId": schoolID,
	}
	_, err := r.coll.UpdateOne(ctx, filter, bson.M{"$set": update})
	return err
}

// ---------- DELETE ----------
func (r *AcademicYearRepository) Delete(ctx context.Context, schoolID, id bson.ObjectID) error {
	// NOTE: higher‑level service should ensure the year is not referenced by any class.
	filter := bson.M{
		"_id":      id,
		"schoolId": schoolID,
	}
	_, err := r.coll.DeleteOne(ctx, filter)
	return err
}

// ---------- ACTIVATE ----------
func (r *AcademicYearRepository) Activate(ctx context.Context, schoolID, id bson.ObjectID) error {
	// 1. Deactivate other active years for the school
	_, err := r.coll.UpdateMany(
		ctx,
		bson.M{"schoolId": schoolID, "isActive": true},
		bson.M{"$set": bson.M{"isActive": false, "status": "inactive", "updatedAt": time.Now()}},
	)
	if err != nil {
		return err
	}
	// 2. Activate the selected year
	_, err = r.coll.UpdateOne(
		ctx,
		bson.M{"_id": id, "schoolId": schoolID},
		bson.M{"$set": bson.M{"isActive": true, "status": "active", "updatedAt": time.Now()}},
	)
	return err
}

// ---------- ARCHIVE ----------
func (r *AcademicYearRepository) Archive(ctx context.Context, schoolID, id bson.ObjectID) error {
	filter := bson.M{"_id": id, "schoolId": schoolID}
	update := bson.M{
		"$set": bson.M{
			"status":    "archived",
			"isActive":  false,
			"updatedAt": time.Now(),
		},
	}
	_, err := r.coll.UpdateOne(ctx, filter, update)
	return err
}
