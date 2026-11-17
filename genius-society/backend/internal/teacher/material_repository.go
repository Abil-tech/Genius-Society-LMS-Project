package teacher

import (
	"context"
	"errors"
	"regexp"

	"github.com/geniussociety/backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// MaterialRepository defines database operations for materials.
type MaterialRepository interface {
	FindMaterials(ctx context.Context, schoolID, teacherID bson.ObjectID, classID, subjectID *bson.ObjectID, search string, page, limit int64) ([]*model.Material, int64, error)
	FindMaterialByID(ctx context.Context, schoolID, teacherID, id bson.ObjectID) (*model.Material, error)
	CreateMaterial(ctx context.Context, m *model.Material) error
	UpdateMaterial(ctx context.Context, m *model.Material) error
	DeleteMaterial(ctx context.Context, schoolID, teacherID, id bson.ObjectID) error
	FindClassesMap(ctx context.Context, schoolID bson.ObjectID, classIDs []bson.ObjectID) (map[bson.ObjectID]*model.Class, error)
	FindSubjectsMap(ctx context.Context, subjectIDs []bson.ObjectID) (map[bson.ObjectID]*model.Subject, error)
}

type materialRepoImpl struct {
	materials *mongo.Collection
	classes   *mongo.Collection
	subjects  *mongo.Collection
}

// NewMaterialRepository creates a new MaterialRepository instance.
func NewMaterialRepository(db *mongo.Database) MaterialRepository {
	return &materialRepoImpl{
		materials: db.Collection("materials"),
		classes:   db.Collection("classes"),
		subjects:  db.Collection("subjects"),
	}
}

func (r *materialRepoImpl) FindMaterials(ctx context.Context, schoolID, teacherID bson.ObjectID, classID, subjectID *bson.ObjectID, search string, page, limit int64) ([]*model.Material, int64, error) {
	filter := bson.M{
		"schoolId":  schoolID,
		"teacherId": teacherID,
	}

	if classID != nil && !classID.IsZero() {
		filter["classIds"] = *classID
	}
	if subjectID != nil && !subjectID.IsZero() {
		filter["subjectId"] = *subjectID
	}
	if search != "" {
		safeSearch := regexp.QuoteMeta(search)
		filter["title"] = bson.M{"$regex": safeSearch, "$options": "i"}
	}

	total, err := r.materials.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	skip := (page - 1) * limit
	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}}).
		SetSkip(skip).
		SetLimit(limit)

	cur, err := r.materials.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(ctx)

	var list []*model.Material
	if err := cur.All(ctx, &list); err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *materialRepoImpl) FindMaterialByID(ctx context.Context, schoolID, teacherID, id bson.ObjectID) (*model.Material, error) {
	filter := bson.M{
		"_id":       id,
		"schoolId":  schoolID,
		"teacherId": teacherID,
	}

	var m model.Material
	err := r.materials.FindOne(ctx, filter).Decode(&m)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r *materialRepoImpl) CreateMaterial(ctx context.Context, m *model.Material) error {
	res, err := r.materials.InsertOne(ctx, m)
	if err != nil {
		return err
	}
	if oid, ok := res.InsertedID.(bson.ObjectID); ok {
		m.ID = oid
	}
	return nil
}

func (r *materialRepoImpl) UpdateMaterial(ctx context.Context, m *model.Material) error {
	filter := bson.M{
		"_id":       m.ID,
		"schoolId":  m.SchoolID,
		"teacherId": m.TeacherID,
	}
	_, err := r.materials.ReplaceOne(ctx, filter, m)
	return err
}

func (r *materialRepoImpl) DeleteMaterial(ctx context.Context, schoolID, teacherID, id bson.ObjectID) error {
	filter := bson.M{
		"_id":       id,
		"schoolId":  schoolID,
		"teacherId": teacherID,
	}
	res, err := r.materials.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *materialRepoImpl) FindClassesMap(ctx context.Context, schoolID bson.ObjectID, classIDs []bson.ObjectID) (map[bson.ObjectID]*model.Class, error) {
	result := make(map[bson.ObjectID]*model.Class)
	if len(classIDs) == 0 {
		return result, nil
	}
	filter := bson.M{
		"_id":      bson.M{"$in": classIDs},
		"schoolId": schoolID,
	}
	cur, err := r.classes.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var c model.Class
		if err := cur.Decode(&c); err == nil {
			result[c.ID] = &c
		}
	}
	return result, nil
}

func (r *materialRepoImpl) FindSubjectsMap(ctx context.Context, subjectIDs []bson.ObjectID) (map[bson.ObjectID]*model.Subject, error) {
	result := make(map[bson.ObjectID]*model.Subject)
	if len(subjectIDs) == 0 {
		return result, nil
	}
	filter := bson.M{"_id": bson.M{"$in": subjectIDs}}
	cur, err := r.subjects.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var s model.Subject
		if err := cur.Decode(&s); err == nil {
			result[s.ID] = &s
		}
	}
	return result, nil
}
