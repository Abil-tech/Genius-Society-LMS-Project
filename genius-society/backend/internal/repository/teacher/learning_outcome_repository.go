package teacherrepo

import (
	"context"
	"time"

	"github.com/geniussociety/backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// LearningOutcomeRepository handles database queries for Pencapaian Pembelajaran (PP).
type LearningOutcomeRepository interface {
	Create(ctx context.Context, lo *model.LearningOutcome) error
	FindByID(ctx context.Context, teacherID, id bson.ObjectID) (*model.LearningOutcome, error)
	FindList(ctx context.Context, teacherID, academicYearID bson.ObjectID, semester int, subjectID *bson.ObjectID, classID *bson.ObjectID) ([]model.LearningOutcome, error)
	Update(ctx context.Context, lo *model.LearningOutcome) error
	Delete(ctx context.Context, teacherID, id bson.ObjectID) error
	CountItemsByOutcomeID(ctx context.Context, outcomeID bson.ObjectID) (int64, error)

	FindAssignmentByID(ctx context.Context, teacherID, id bson.ObjectID) (*model.Assignment, error)
	UpdateAssignmentOutcome(ctx context.Context, teacherID, id bson.ObjectID, outcomeID *bson.ObjectID) error

	FindProjectByID(ctx context.Context, teacherID, id bson.ObjectID) (*model.Project, error)
	UpdateProjectOutcome(ctx context.Context, teacherID, id bson.ObjectID, outcomeID *bson.ObjectID) error

	FindAssessmentByID(ctx context.Context, teacherID, id bson.ObjectID) (*model.Assessment, error)
	UpdateAssessmentOutcome(ctx context.Context, teacherID, id bson.ObjectID, outcomeID *bson.ObjectID) error

	FindClassesMap(ctx context.Context, classIDs []bson.ObjectID) (map[bson.ObjectID]*model.Class, error)
	FindSubjectsMap(ctx context.Context, subjectIDs []bson.ObjectID) (map[bson.ObjectID]*model.Subject, error)
}

type learningOutcomeRepositoryImpl struct {
	db          *mongo.Database
	collection  *mongo.Collection
	assignments *mongo.Collection
	projects    *mongo.Collection
	assessments *mongo.Collection
	classes     *mongo.Collection
	subjects    *mongo.Collection
}

// NewLearningOutcomeRepository creates a new LearningOutcomeRepository instance.
func NewLearningOutcomeRepository(db *mongo.Database) LearningOutcomeRepository {
	return &learningOutcomeRepositoryImpl{
		db:          db,
		collection:  db.Collection("learning_outcomes"),
		assignments: db.Collection("assignments"),
		projects:    db.Collection("projects"),
		assessments: db.Collection("assessments"),
		classes:     db.Collection("classes"),
		subjects:    db.Collection("subjects"),
	}
}

func (r *learningOutcomeRepositoryImpl) Create(ctx context.Context, lo *model.LearningOutcome) error {
	_, err := r.collection.InsertOne(ctx, lo)
	return err
}

func (r *learningOutcomeRepositoryImpl) FindByID(ctx context.Context, teacherID, id bson.ObjectID) (*model.LearningOutcome, error) {
	filter := bson.M{
		"_id":      id,
		
	}
	if !teacherID.IsZero() {
		filter["teacherId"] = teacherID
	}
	var lo model.LearningOutcome
	err := r.collection.FindOne(ctx, filter).Decode(&lo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &lo, nil
}

func (r *learningOutcomeRepositoryImpl) FindList(ctx context.Context, teacherID, academicYearID bson.ObjectID, semester int, subjectID *bson.ObjectID, classID *bson.ObjectID) ([]model.LearningOutcome, error) {
	filter := bson.M{
		
	}
	if !teacherID.IsZero() {
		filter["teacherId"] = teacherID
	}
	if !academicYearID.IsZero() {
		filter["academicYearId"] = academicYearID
	}
	if semester > 0 {
		filter["semester"] = semester
	}
	if subjectID != nil && !subjectID.IsZero() {
		filter["subjectId"] = *subjectID
	}
	if classID != nil && !classID.IsZero() {
		filter["classIds"] = *classID
	}

	opts := options.Find().SetSort(bson.D{{Key: "order", Value: 1}, {Key: "createdAt", Value: 1}})
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var list []model.LearningOutcome
	if err := cursor.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *learningOutcomeRepositoryImpl) Update(ctx context.Context, lo *model.LearningOutcome) error {
	filter := bson.M{
		"_id":       lo.ID,
		"schoolId":  lo.SchoolID,
		"teacherId": lo.TeacherID,
	}
	lo.UpdatedAt = time.Now().UTC()
	update := bson.M{"$set": lo}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *learningOutcomeRepositoryImpl) Delete(ctx context.Context, teacherID, id bson.ObjectID) error {
	filter := bson.M{
		"_id":       id,
		
		"teacherId": teacherID,
	}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}

func (r *learningOutcomeRepositoryImpl) CountItemsByOutcomeID(ctx context.Context, outcomeID bson.ObjectID) (int64, error) {
	asgCount, err := r.assignments.CountDocuments(ctx, bson.M{ "learningOutcomeId": outcomeID})
	if err != nil {
		return 0, err
	}
	prjCount, err := r.projects.CountDocuments(ctx, bson.M{ "learningOutcomeId": outcomeID})
	if err != nil {
		return 0, err
	}
	asmCount, err := r.assessments.CountDocuments(ctx, bson.M{ "learningOutcomeId": outcomeID})
	if err != nil {
		return 0, err
	}
	return asgCount + prjCount + asmCount, nil
}

func (r *learningOutcomeRepositoryImpl) FindAssignmentByID(ctx context.Context, teacherID, id bson.ObjectID) (*model.Assignment, error) {
	var asg model.Assignment
	err := r.assignments.FindOne(ctx, bson.M{"_id": id,  "teacherId": teacherID}).Decode(&asg)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &asg, nil
}

func (r *learningOutcomeRepositoryImpl) UpdateAssignmentOutcome(ctx context.Context, teacherID, id bson.ObjectID, outcomeID *bson.ObjectID) error {
	update := bson.M{
		"$set": bson.M{
			"learningOutcomeId": outcomeID,
			"updatedAt":         time.Now().UTC(),
		},
	}
	_, err := r.assignments.UpdateOne(ctx, bson.M{"_id": id,  "teacherId": teacherID}, update)
	return err
}

func (r *learningOutcomeRepositoryImpl) FindProjectByID(ctx context.Context, teacherID, id bson.ObjectID) (*model.Project, error) {
	var prj model.Project
	err := r.projects.FindOne(ctx, bson.M{"_id": id,  "teacherId": teacherID}).Decode(&prj)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &prj, nil
}

func (r *learningOutcomeRepositoryImpl) UpdateProjectOutcome(ctx context.Context, teacherID, id bson.ObjectID, outcomeID *bson.ObjectID) error {
	update := bson.M{
		"$set": bson.M{
			"learningOutcomeId": outcomeID,
			"updatedAt":         time.Now().UTC(),
		},
	}
	_, err := r.projects.UpdateOne(ctx, bson.M{"_id": id,  "teacherId": teacherID}, update)
	return err
}

func (r *learningOutcomeRepositoryImpl) FindAssessmentByID(ctx context.Context, teacherID, id bson.ObjectID) (*model.Assessment, error) {
	var asm model.Assessment
	err := r.assessments.FindOne(ctx, bson.M{"_id": id,  "teacherId": teacherID}).Decode(&asm)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &asm, nil
}

func (r *learningOutcomeRepositoryImpl) UpdateAssessmentOutcome(ctx context.Context, teacherID, id bson.ObjectID, outcomeID *bson.ObjectID) error {
	update := bson.M{
		"$set": bson.M{
			"learningOutcomeId": outcomeID,
			"updatedAt":         time.Now().UTC(),
		},
	}
	_, err := r.assessments.UpdateOne(ctx, bson.M{"_id": id,  "teacherId": teacherID}, update)
	return err
}

func (r *learningOutcomeRepositoryImpl) FindClassesMap(ctx context.Context, classIDs []bson.ObjectID) (map[bson.ObjectID]*model.Class, error) {
	res := make(map[bson.ObjectID]*model.Class)
	if len(classIDs) == 0 {
		return res, nil
	}
	cursor, err := r.classes.Find(ctx, bson.M{"_id": bson.M{"$in": classIDs}})
	if err != nil {
		return res, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var c model.Class
		if err := cursor.Decode(&c); err == nil {
			res[c.ID] = &c
		}
	}
	return res, nil
}

func (r *learningOutcomeRepositoryImpl) FindSubjectsMap(ctx context.Context, subjectIDs []bson.ObjectID) (map[bson.ObjectID]*model.Subject, error) {
	res := make(map[bson.ObjectID]*model.Subject)
	if len(subjectIDs) == 0 {
		return res, nil
	}
	cursor, err := r.subjects.Find(ctx, bson.M{"_id": bson.M{"$in": subjectIDs}})
	if err != nil {
		return res, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var s model.Subject
		if err := cursor.Decode(&s); err == nil {
			res[s.ID] = &s
		}
	}
	return res, nil
}
