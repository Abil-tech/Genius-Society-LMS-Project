package teacherrepo

import (
	"context"
	"time"

	"github.com/geniussociety/backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// GradeRepository defines persistence operations for teacher grade management.
type GradeRepository interface {
	FindGradesByClassSubject(ctx context.Context, schoolID, classID, subjectID, academicYearID bson.ObjectID) ([]*model.Grade, error)
	FindStudentsByClassID(ctx context.Context, schoolID, classID bson.ObjectID) ([]*model.Student, error)
	FindUsersByIDs(ctx context.Context, userIDs []bson.ObjectID) (map[bson.ObjectID]*model.User, error)
	FindGradingConfig(ctx context.Context, schoolID, academicYearID, subjectID bson.ObjectID) (*model.GradingConfig, error)
	UpsertExamGrade(ctx context.Context, grade *model.Grade) (*model.Grade, error)
	FindStudentByID(ctx context.Context, studentID bson.ObjectID) (*model.Student, error)
	FindClassByID(ctx context.Context, classID bson.ObjectID) (*model.Class, error)
	FindSubjectByID(ctx context.Context, subjectID bson.ObjectID) (*model.Subject, error)
}

type gradeRepoImpl struct {
	grades         *mongo.Collection
	gradingConfigs *mongo.Collection
	students       *mongo.Collection
	users          *mongo.Collection
	classes        *mongo.Collection
	subjects       *mongo.Collection
}

// NewGradeRepository creates a new GradeRepository instance.
func NewGradeRepository(db *mongo.Database) GradeRepository {
	return &gradeRepoImpl{
		grades:         db.Collection("grades"),
		gradingConfigs: db.Collection("grading_configs"),
		students:       db.Collection("students"),
		users:          db.Collection("users"),
		classes:        db.Collection("classes"),
		subjects:       db.Collection("subjects"),
	}
}

func (r *gradeRepoImpl) FindGradesByClassSubject(ctx context.Context, schoolID, classID, subjectID, academicYearID bson.ObjectID) ([]*model.Grade, error) {
	filter := bson.M{
		"schoolId":       schoolID,
		"classId":        classID,
		"subjectId":      subjectID,
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

func (r *gradeRepoImpl) FindStudentsByClassID(ctx context.Context, schoolID, classID bson.ObjectID) ([]*model.Student, error) {
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

func (r *gradeRepoImpl) FindUsersByIDs(ctx context.Context, userIDs []bson.ObjectID) (map[bson.ObjectID]*model.User, error) {
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

func (r *gradeRepoImpl) FindGradingConfig(ctx context.Context, schoolID, academicYearID, subjectID bson.ObjectID) (*model.GradingConfig, error) {
	// First attempt: find subject-specific config
	var cfg model.GradingConfig
	err := r.gradingConfigs.FindOne(ctx, bson.M{
		"schoolId":       schoolID,
		"academicYearId": academicYearID,
		"subjectId":      subjectID,
	}).Decode(&cfg)
	if err == nil {
		return &cfg, nil
	}

	// Fallback: school-wide default config (subjectId == nil or doesn't exist)
	err = r.gradingConfigs.FindOne(ctx, bson.M{
		"schoolId":       schoolID,
		"academicYearId": academicYearID,
		"subjectId":      nil,
	}).Decode(&cfg)
	if err == nil {
		return &cfg, nil
	}

	return nil, nil
}

func (r *gradeRepoImpl) UpsertExamGrade(ctx context.Context, grade *model.Grade) (*model.Grade, error) {
	now := time.Now()
	filter := bson.M{
		"schoolId":       grade.SchoolID,
		"studentId":      grade.StudentID,
		"subjectId":      grade.SubjectID,
		"classId":        grade.ClassID,
		"academicYearId": grade.AcademicYearID,
		"semester":       grade.Semester,
		"type":           grade.Type,
	}

	update := bson.M{
		"$set": bson.M{
			"score":     grade.Score,
			"teacherId": grade.TeacherID,
			"source":    "manual",
			"updatedAt": now,
		},
		"$setOnInsert": bson.M{
			"_id":       bson.NewObjectID(),
			"createdAt": now,
		},
	}

	opts := options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After)

	var result model.Grade
	err := r.grades.FindOneAndUpdate(ctx, filter, update, opts).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *gradeRepoImpl) FindStudentByID(ctx context.Context, studentID bson.ObjectID) (*model.Student, error) {
	var st model.Student
	err := r.students.FindOne(ctx, bson.M{"_id": studentID}).Decode(&st)
	if err != nil {
		return nil, err
	}
	return &st, nil
}

func (r *gradeRepoImpl) FindClassByID(ctx context.Context, classID bson.ObjectID) (*model.Class, error) {
	var cl model.Class
	err := r.classes.FindOne(ctx, bson.M{"_id": classID}).Decode(&cl)
	if err != nil {
		return nil, err
	}
	return &cl, nil
}

func (r *gradeRepoImpl) FindSubjectByID(ctx context.Context, subjectID bson.ObjectID) (*model.Subject, error) {
	var sub model.Subject
	err := r.subjects.FindOne(ctx, bson.M{"_id": subjectID}).Decode(&sub)
	if err != nil {
		return nil, err
	}
	return &sub, nil
}
