package repository

import (
	"context"
	"time"

	"github.com/geniussociety/backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// SchoolRepository handles persistence of School documents.
type SchoolRepository struct {
	coll *mongo.Collection
}

// NewSchoolRepository returns a repository bound to the supplied database.
func NewSchoolRepository(db *mongo.Database) *SchoolRepository {
	return &SchoolRepository{
		coll: db.Collection("schools"),
	}
}

// Create inserts a new school document.
func (r *SchoolRepository) Create(ctx context.Context, s *model.School) error {
	s.ID = bson.NewObjectID()
	s.CreatedAt = time.Now()
	s.UpdatedAt = s.CreatedAt
	_, err := r.coll.InsertOne(ctx, s)
	return err
}

// FindByID returns a school by its ObjectID.
func (r *SchoolRepository) FindByID(ctx context.Context, id bson.ObjectID) (*model.School, error) {
	var s model.School
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&s)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

// FindByName returns a school by its name.
func (r *SchoolRepository) FindByName(ctx context.Context, name string) (*model.School, error) {
	var s model.School
	err := r.coll.FindOne(ctx, bson.M{"name": name}).Decode(&s)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

// FindOrCreateByName returns an existing school with the given name or creates one if not present.
func (r *SchoolRepository) FindOrCreateByName(ctx context.Context, name, address string) (*model.School, error) {
	s, err := r.FindByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if s != nil {
		return s, nil
	}
	// Not found – create new.
	s = &model.School{Name: name, Address: address}
	if err := r.Create(ctx, s); err != nil {
		return nil, err
	}
	return s, nil
}
