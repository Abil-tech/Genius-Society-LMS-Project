package config

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// EnsureTeacherIndexes creates the required database indexes for the Teacher module.
func EnsureTeacherIndexes(ctx context.Context, db *mongo.Database) error {
	type indexSpec struct {
		collection string
		model      mongo.IndexModel
	}

	indexes := []indexSpec{
		// assignment_submissions: unique {assignmentId, studentId, version}
		{
			collection: "assignment_submissions",
			model: mongo.IndexModel{
				Keys: bson.D{
					{Key: "assignmentId", Value: 1},
					{Key: "studentId", Value: 1},
					{Key: "version", Value: 1},
				},
				Options: options.Index().SetUnique(true),
			},
		},
		// assessment_attempts: unique {assessmentId, studentId, attemptNumber}
		{
			collection: "assessment_attempts",
			model: mongo.IndexModel{
				Keys: bson.D{
					{Key: "assessmentId", Value: 1},
					{Key: "studentId", Value: 1},
					{Key: "attemptNumber", Value: 1},
				},
				Options: options.Index().SetUnique(true),
			},
		},
		// assessment_answers: unique {attemptId, questionId}
		{
			collection: "assessment_answers",
			model: mongo.IndexModel{
				Keys: bson.D{
					{Key: "attemptId", Value: 1},
					{Key: "questionId", Value: 1},
				},
				Options: options.Index().SetUnique(true),
			},
		},
		// grades: unique {studentId, subjectId, academicYearId, semester, type, sourceId}
		{
			collection: "grades",
			model: mongo.IndexModel{
				Keys: bson.D{
					{Key: "studentId", Value: 1},
					{Key: "subjectId", Value: 1},
					{Key: "academicYearId", Value: 1},
					{Key: "semester", Value: 1},
					{Key: "type", Value: 1},
					{Key: "sourceId", Value: 1},
				},
				Options: options.Index().SetUnique(true),
			},
		},
		// grading_configs: unique {schoolId, academicYearId, subjectId}
		{
			collection: "grading_configs",
			model: mongo.IndexModel{
				Keys: bson.D{
					{Key: "schoolId", Value: 1},
					{Key: "academicYearId", Value: 1},
					{Key: "subjectId", Value: 1},
				},
				Options: options.Index().SetUnique(true),
			},
		},
		// learning_outcomes: unique {schoolId, academicYearId, semester, subjectId, teacherId, titleKey}
		{
			collection: "learning_outcomes",
			model: mongo.IndexModel{
				Keys: bson.D{
					{Key: "schoolId", Value: 1},
					{Key: "academicYearId", Value: 1},
					{Key: "semester", Value: 1},
					{Key: "subjectId", Value: 1},
					{Key: "teacherId", Value: 1},
					{Key: "titleKey", Value: 1},
				},
				Options: options.Index().SetUnique(true),
			},
		},
		// teacher_classes: unique {teacherId, classId, subjectId, academicYearId}
		{
			collection: "teacher_classes",
			model: mongo.IndexModel{
				Keys: bson.D{
					{Key: "teacherId", Value: 1},
					{Key: "classId", Value: 1},
					{Key: "subjectId", Value: 1},
					{Key: "academicYearId", Value: 1},
				},
				Options: options.Index().SetUnique(true),
			},
		},
		// notifications: {userId, isRead, createdAt: -1}
		{
			collection: "notifications",
			model: mongo.IndexModel{
				Keys: bson.D{
					{Key: "userId", Value: 1},
					{Key: "isRead", Value: 1},
					{Key: "createdAt", Value: -1},
				},
			},
		},
		// notifications: unique {userId, type, refType, refId}
		{
			collection: "notifications",
			model: mongo.IndexModel{
				Keys: bson.D{
					{Key: "userId", Value: 1},
					{Key: "type", Value: 1},
					{Key: "refType", Value: 1},
					{Key: "refId", Value: 1},
				},
				Options: options.Index().SetUnique(true),
			},
		},
		// materials: {schoolId, teacherId, createdAt: -1}
		{
			collection: "materials",
			model: mongo.IndexModel{
				Keys: bson.D{
					{Key: "schoolId", Value: 1},
					{Key: "teacherId", Value: 1},
					{Key: "createdAt", Value: -1},
				},
			},
		},
		// assignments: {schoolId, teacherId, createdAt: -1}
		{
			collection: "assignments",
			model: mongo.IndexModel{
				Keys: bson.D{
					{Key: "schoolId", Value: 1},
					{Key: "teacherId", Value: 1},
					{Key: "createdAt", Value: -1},
				},
			},
		},
		// projects: {schoolId, teacherId, createdAt: -1}
		{
			collection: "projects",
			model: mongo.IndexModel{
				Keys: bson.D{
					{Key: "schoolId", Value: 1},
					{Key: "teacherId", Value: 1},
					{Key: "createdAt", Value: -1},
				},
			},
		},
		// assessments: {schoolId, teacherId, createdAt: -1}
		{
			collection: "assessments",
			model: mongo.IndexModel{
				Keys: bson.D{
					{Key: "schoolId", Value: 1},
					{Key: "teacherId", Value: 1},
					{Key: "createdAt", Value: -1},
				},
			},
		},
	}

	for _, idx := range indexes {
		_, err := db.Collection(idx.collection).Indexes().CreateOne(ctx, idx.model)
		if err != nil {
			return err
		}
	}

	return nil
}
