package repository

import (
	"context"
	"errors"

	"github.com/geniussociety/backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// ClassSubjectPair represents an active teaching assignment pairing class and subject.
type ClassSubjectPair struct {
	ClassID   bson.ObjectID `bson:"classId" json:"classId"`
	SubjectID bson.ObjectID `bson:"subjectId" json:"subjectId"`
}

// TeacherScopeRepository is a read-only repository used by TeacherScope middleware
// and teacher domain services. Per project rules, it is strictly read-only.
type TeacherScopeRepository interface {
	FindUserByID(ctx context.Context, userID bson.ObjectID) (*model.User, error)
	FindTeacherByUserID(ctx context.Context, userID bson.ObjectID) (*model.Teacher, error)
	FindActiveAcademicYear(ctx context.Context, schoolID bson.ObjectID) (*model.AcademicYear, error)
	FindTeacherClassPairs(ctx context.Context, schoolID, teacherID, academicYearID bson.ObjectID) ([]ClassSubjectPair, error)
	FindTeacherSubjectIDs(ctx context.Context, schoolID, teacherID bson.ObjectID) ([]bson.ObjectID, error)
}

type teacherScopeRepoImpl struct {
	users           *mongo.Collection
	teachers        *mongo.Collection
	academicYears   *mongo.Collection
	teacherClasses  *mongo.Collection
	teacherSubjects *mongo.Collection
}

// NewTeacherScopeRepository constructs the read-only repository for TeacherScope.
func NewTeacherScopeRepository(db *mongo.Database) TeacherScopeRepository {
	return &teacherScopeRepoImpl{
		users:           db.Collection("users"),
		teachers:        db.Collection("teachers"),
		academicYears:   db.Collection("academic_years"),
		teacherClasses:  db.Collection("teacher_classes"),
		teacherSubjects: db.Collection("teacher_subjects"),
	}
}

func (r *teacherScopeRepoImpl) FindUserByID(ctx context.Context, userID bson.ObjectID) (*model.User, error) {
	var u model.User
	err := r.users.FindOne(ctx, bson.M{"_id": userID}).Decode(&u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *teacherScopeRepoImpl) FindTeacherByUserID(ctx context.Context, userID bson.ObjectID) (*model.Teacher, error) {
	var t model.Teacher
	err := r.teachers.FindOne(ctx, bson.M{"userId": userID}).Decode(&t)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

func (r *teacherScopeRepoImpl) FindActiveAcademicYear(ctx context.Context, schoolID bson.ObjectID) (*model.AcademicYear, error) {
	var ay model.AcademicYear
	filter := bson.M{
		"schoolId": schoolID,
		"isActive": true,
	}
	err := r.academicYears.FindOne(ctx, filter).Decode(&ay)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &ay, nil
}

func (r *teacherScopeRepoImpl) FindTeacherClassPairs(ctx context.Context, schoolID, teacherID, academicYearID bson.ObjectID) ([]ClassSubjectPair, error) {
	filter := bson.M{
		"teacherId":      teacherID,
		"academicYearId": academicYearID,
	}
	if !schoolID.IsZero() {
		filter["$or"] = []bson.M{
			{"schoolId": schoolID},
			{"schoolId": bson.M{"$exists": false}},
		}
	}

	cur, err := r.teacherClasses.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var pairs []ClassSubjectPair
	for cur.Next(ctx) {
		var doc struct {
			ClassID   bson.ObjectID `bson:"classId"`
			SubjectID bson.ObjectID `bson:"subjectId"`
		}
		if err := cur.Decode(&doc); err != nil {
			return nil, err
		}
		pairs = append(pairs, ClassSubjectPair{
			ClassID:   doc.ClassID,
			SubjectID: doc.SubjectID,
		})
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return pairs, nil
}

func (r *teacherScopeRepoImpl) FindTeacherSubjectIDs(ctx context.Context, schoolID, teacherID bson.ObjectID) ([]bson.ObjectID, error) {
	filter := bson.M{
		"teacherId": teacherID,
	}
	if !schoolID.IsZero() {
		filter["$or"] = []bson.M{
			{"schoolId": schoolID},
			{"schoolId": bson.M{"$exists": false}},
		}
	}

	cur, err := r.teacherSubjects.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var subjectIDs []bson.ObjectID
	for cur.Next(ctx) {
		var doc struct {
			SubjectID bson.ObjectID `bson:"subjectId"`
		}
		if err := cur.Decode(&doc); err != nil {
			return nil, err
		}
		subjectIDs = append(subjectIDs, doc.SubjectID)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return subjectIDs, nil
}
