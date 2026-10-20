package config

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ConnectMongo(uri, dbName string) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	db := client.Database(dbName)
	if err := EnsureIndexes(ctx, db); err != nil {
		return nil, err
	}

	return db, nil
}

func EnsureIndexes(ctx context.Context, db *mongo.Database) error {
	type indexSpec struct {
		collection string
		model      mongo.IndexModel
	}

	indexes := []indexSpec{
		// users
		{
			collection: "users",
			model: mongo.IndexModel{
				Keys:    bson.D{{Key: "email", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
		},
		{
			collection: "users",
			model: mongo.IndexModel{
				Keys:    bson.D{{Key: "username", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
		},
		// students
		{
			collection: "students",
			model: mongo.IndexModel{
				Keys:    bson.D{{Key: "userId", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
		},
		{
			collection: "students",
			model: mongo.IndexModel{
				Keys:    bson.D{{Key: "nis", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
		},
		// teachers
		{
			collection: "teachers",
			model: mongo.IndexModel{
				Keys:    bson.D{{Key: "userId", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
		},
		{
			collection: "teachers",
			model: mongo.IndexModel{
				Keys:    bson.D{{Key: "nip", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
		},
		// classes
		{
			collection: "classes",
			model: mongo.IndexModel{
				Keys:    bson.D{{Key: "name", Value: 1}, {Key: "academicYearId", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
		},
		// subjects
		{
			collection: "subjects",
			model: mongo.IndexModel{
				Keys:    bson.D{{Key: "code", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
		},
		// academic_years
		{
			collection: "academic_years",
			model: mongo.IndexModel{
				Keys:    bson.D{{Key: "name", Value: 1}, {Key: "semester", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
		},
		// class_students
		{
			collection: "class_students",
			model: mongo.IndexModel{
				Keys:    bson.D{{Key: "classId", Value: 1}, {Key: "studentId", Value: 1}, {Key: "academicYearId", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
		},
		// teacher_subjects
		{
			collection: "teacher_subjects",
			model: mongo.IndexModel{
				Keys:    bson.D{{Key: "teacherId", Value: 1}, {Key: "subjectId", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
		},
		// teacher_classes
		{
			collection: "teacher_classes",
			model: mongo.IndexModel{
				Keys:    bson.D{{Key: "teacherId", Value: 1}, {Key: "classId", Value: 1}, {Key: "subjectId", Value: 1}, {Key: "academicYearId", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
		},
		// schedules
		{
			collection: "schedules",
			model: mongo.IndexModel{
				Keys: bson.D{{Key: "teacherId", Value: 1}, {Key: "day", Value: 1}, {Key: "startTime", Value: 1}},
			},
		},
		{
			collection: "schedules",
			model: mongo.IndexModel{
				Keys: bson.D{{Key: "classId", Value: 1}, {Key: "day", Value: 1}, {Key: "startTime", Value: 1}},
			},
		},
		// assignment_submissions
		{
			collection: "assignment_submissions",
			model: mongo.IndexModel{
				Keys:    bson.D{{Key: "assignmentId", Value: 1}, {Key: "studentId", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
		},
		// assessment_attempts
		{
			collection: "assessment_attempts",
			model: mongo.IndexModel{
				Keys:    bson.D{{Key: "assessmentId", Value: 1}, {Key: "studentId", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
		},
		// assessment_answers
		{
			collection: "assessment_answers",
			model: mongo.IndexModel{
				Keys:    bson.D{{Key: "attemptId", Value: 1}, {Key: "questionId", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
		},
		// project_submissions
		{
			collection: "project_submissions",
			model: mongo.IndexModel{
				Keys:    bson.D{{Key: "projectId", Value: 1}, {Key: "studentId", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
		},
		// grades
		{
			collection: "grades",
			model: mongo.IndexModel{
				Keys:    bson.D{{Key: "studentId", Value: 1}, {Key: "subjectId", Value: 1}, {Key: "academicYearId", Value: 1}, {Key: "type", Value: 1}, {Key: "sourceId", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
		},
		// grading_configs
		{
			collection: "grading_configs",
			model: mongo.IndexModel{
				Keys:    bson.D{{Key: "academicYearId", Value: 1}, {Key: "subjectId", Value: 1}},
				Options: options.Index().SetUnique(true).SetSparse(true),
			},
		},
		// notifications
		{
			collection: "notifications",
			model: mongo.IndexModel{
				Keys: bson.D{{Key: "userId", Value: 1}, {Key: "isRead", Value: 1}, {Key: "createdAt", Value: -1}},
			},
		},
	}

	for _, idx := range indexes {
		coll := db.Collection(idx.collection)
		_, err := coll.Indexes().CreateOne(ctx, idx.model)
		if err != nil {
			// ignore duplicate key index error (already exists)
			if mongo.IsDuplicateKeyError(err) {
				continue
			}
			return err
		}
	}
	return nil
}
