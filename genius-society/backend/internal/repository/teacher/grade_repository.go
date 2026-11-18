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

	BulkUpsertSTS(ctx context.Context, schoolID, teacherID bson.ObjectID, grades []*model.Grade, audits []*model.AuditLog) error
	FindExistingExamGrade(ctx context.Context, schoolID, studentID, subjectID, academicYearID bson.ObjectID, semester int, gradeType model.GradeType) (*model.Grade, error)

	FindAssignment(ctx context.Context, schoolID, teacherID, id bson.ObjectID) (*model.Assignment, error)
	HasAssignmentSubmission(ctx context.Context, schoolID, assignmentID, studentID bson.ObjectID) (bool, error)

	FindProject(ctx context.Context, schoolID, teacherID, id bson.ObjectID) (*model.Project, error)
	HasProjectSubmission(ctx context.Context, schoolID, projectID, studentID bson.ObjectID) (bool, error)

	FindAssessment(ctx context.Context, schoolID, teacherID, id bson.ObjectID) (*model.Assessment, error)
	HasAssessmentAttempt(ctx context.Context, schoolID, assessmentID, studentID bson.ObjectID) (bool, error)

	SetManualItemGrade(ctx context.Context, grade *model.Grade, audit *model.AuditLog) error
	DeleteManualItemGrade(ctx context.Context, schoolID, sourceID, studentID bson.ObjectID, audit *model.AuditLog) error
	FindManualItemGrade(ctx context.Context, schoolID, sourceID, studentID bson.ObjectID) (*model.Grade, error)

	FindLearningOutcomesForRecap(ctx context.Context, schoolID, academicYearID, subjectID, classID bson.ObjectID, semester int) ([]*model.LearningOutcome, error)
	FindItemsForRecap(ctx context.Context, schoolID, academicYearID, subjectID, classID bson.ObjectID) ([]*model.Assignment, []*model.Project, []*model.Assessment, error)
}

type gradeRepoImpl struct {
	db                    *mongo.Database
	grades                *mongo.Collection
	gradingConfigs        *mongo.Collection
	students              *mongo.Collection
	users                 *mongo.Collection
	classes               *mongo.Collection
	subjects              *mongo.Collection
	audits                *mongo.Collection
	assignmentSubmissions *mongo.Collection
	projectSubmissions    *mongo.Collection
	assessmentAttempts    *mongo.Collection
	assignments           *mongo.Collection
	projects              *mongo.Collection
	assessments           *mongo.Collection
	learningOutcomes      *mongo.Collection
}

// NewGradeRepository creates a new GradeRepository instance.
func NewGradeRepository(db *mongo.Database) GradeRepository {
	return &gradeRepoImpl{
		db:                    db,
		grades:                db.Collection("grades"),
		gradingConfigs:        db.Collection("grading_configs"),
		students:              db.Collection("students"),
		users:                 db.Collection("users"),
		classes:               db.Collection("classes"),
		subjects:              db.Collection("subjects"),
		audits:                db.Collection("audit_logs"),
		assignmentSubmissions: db.Collection("assignment_submissions"),
		projectSubmissions:    db.Collection("project_submissions"),
		assessmentAttempts:    db.Collection("assessment_attempts"),
		assignments:           db.Collection("assignments"),
		projects:              db.Collection("projects"),
		assessments:           db.Collection("assessments"),
		learningOutcomes:      db.Collection("learning_outcomes"),
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

func (r *gradeRepoImpl) BulkUpsertSTS(ctx context.Context, schoolID, teacherID bson.ObjectID, grades []*model.Grade, audits []*model.AuditLog) error {
	session, err := r.db.Client().StartSession()
	if err != nil {
		// Fallback without transaction if standalone
		for _, g := range grades {
			_, err := r.UpsertExamGrade(ctx, g)
			if err != nil {
				return err
			}
		}
		for _, a := range audits {
			_, _ = r.audits.InsertOne(ctx, a)
		}
		return nil
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sc context.Context) (interface{}, error) {
		for _, g := range grades {
			_, err := r.UpsertExamGrade(sc, g)
			if err != nil {
				return nil, err
			}
		}
		for _, a := range audits {
			_, err := r.audits.InsertOne(sc, a)
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	})
	return err
}

func (r *gradeRepoImpl) FindExistingExamGrade(ctx context.Context, schoolID, studentID, subjectID, academicYearID bson.ObjectID, semester int, gradeType model.GradeType) (*model.Grade, error) {
	filter := bson.M{
		"schoolId":       schoolID,
		"studentId":      studentID,
		"subjectId":      subjectID,
		"academicYearId": academicYearID,
		"semester":       semester,
		"type":           gradeType,
	}
	var g model.Grade
	err := r.grades.FindOne(ctx, filter).Decode(&g)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &g, nil
}

func (r *gradeRepoImpl) FindAssignment(ctx context.Context, schoolID, teacherID, id bson.ObjectID) (*model.Assignment, error) {
	var asg model.Assignment
	err := r.assignments.FindOne(ctx, bson.M{"_id": id, "schoolId": schoolID, "teacherId": teacherID}).Decode(&asg)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &asg, nil
}

func (r *gradeRepoImpl) HasAssignmentSubmission(ctx context.Context, schoolID, assignmentID, studentID bson.ObjectID) (bool, error) {
	count, err := r.assignmentSubmissions.CountDocuments(ctx, bson.M{
		"schoolId":     schoolID,
		"assignmentId": assignmentID,
		"studentId":    studentID,
	})
	return count > 0, err
}

func (r *gradeRepoImpl) FindProject(ctx context.Context, schoolID, teacherID, id bson.ObjectID) (*model.Project, error) {
	var prj model.Project
	err := r.projects.FindOne(ctx, bson.M{"_id": id, "schoolId": schoolID, "teacherId": teacherID}).Decode(&prj)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &prj, nil
}

func (r *gradeRepoImpl) HasProjectSubmission(ctx context.Context, schoolID, projectID, studentID bson.ObjectID) (bool, error) {
	count, err := r.projectSubmissions.CountDocuments(ctx, bson.M{
		"schoolId":  schoolID,
		"projectId": projectID,
		"studentId": studentID,
	})
	return count > 0, err
}

func (r *gradeRepoImpl) FindAssessment(ctx context.Context, schoolID, teacherID, id bson.ObjectID) (*model.Assessment, error) {
	var asm model.Assessment
	err := r.assessments.FindOne(ctx, bson.M{"_id": id, "schoolId": schoolID, "teacherId": teacherID}).Decode(&asm)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &asm, nil
}

func (r *gradeRepoImpl) HasAssessmentAttempt(ctx context.Context, schoolID, assessmentID, studentID bson.ObjectID) (bool, error) {
	count, err := r.assessmentAttempts.CountDocuments(ctx, bson.M{
		"schoolId":     schoolID,
		"assessmentId": assessmentID,
		"studentId":    studentID,
	})
	return count > 0, err
}

func (r *gradeRepoImpl) SetManualItemGrade(ctx context.Context, grade *model.Grade, audit *model.AuditLog) error {
	session, err := r.db.Client().StartSession()
	if err != nil {
		now := time.Now().UTC()
		filter := bson.M{
			"schoolId":       grade.SchoolID,
			"studentId":      grade.StudentID,
			"subjectId":      grade.SubjectID,
			"academicYearId": grade.AcademicYearID,
			"semester":       grade.Semester,
			"type":           grade.Type,
			"sourceId":       grade.SourceID,
		}
		update := bson.M{
			"$set": bson.M{
				"score":     grade.Score,
				"source":    "manual",
				"reason":    grade.Reason,
				"gradedBy":  grade.GradedBy,
				"teacherId": grade.TeacherID,
				"classId":   grade.ClassID,
				"updatedAt": now,
			},
			"$setOnInsert": bson.M{
				"_id":       bson.NewObjectID(),
				"createdAt": now,
			},
		}
		opts := options.UpdateOne().SetUpsert(true)
		_, err := r.grades.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			return err
		}
		if audit != nil {
			_, _ = r.audits.InsertOne(ctx, audit)
		}
		return nil
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sc context.Context) (interface{}, error) {
		now := time.Now().UTC()
		filter := bson.M{
			"schoolId":       grade.SchoolID,
			"studentId":      grade.StudentID,
			"subjectId":      grade.SubjectID,
			"academicYearId": grade.AcademicYearID,
			"semester":       grade.Semester,
			"type":           grade.Type,
			"sourceId":       grade.SourceID,
		}
		update := bson.M{
			"$set": bson.M{
				"score":     grade.Score,
				"source":    "manual",
				"reason":    grade.Reason,
				"gradedBy":  grade.GradedBy,
				"teacherId": grade.TeacherID,
				"classId":   grade.ClassID,
				"updatedAt": now,
			},
			"$setOnInsert": bson.M{
				"_id":       bson.NewObjectID(),
				"createdAt": now,
			},
		}
		opts := options.UpdateOne().SetUpsert(true)
		_, err := r.grades.UpdateOne(sc, filter, update, opts)
		if err != nil {
			return nil, err
		}
		if audit != nil {
			_, err := r.audits.InsertOne(sc, audit)
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	})
	return err
}

func (r *gradeRepoImpl) DeleteManualItemGrade(ctx context.Context, schoolID, sourceID, studentID bson.ObjectID, audit *model.AuditLog) error {
	session, err := r.db.Client().StartSession()
	if err != nil {
		_, err := r.grades.DeleteOne(ctx, bson.M{"schoolId": schoolID, "sourceId": sourceID, "studentId": studentID, "source": "manual"})
		if err != nil {
			return err
		}
		if audit != nil {
			_, _ = r.audits.InsertOne(ctx, audit)
		}
		return nil
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sc context.Context) (interface{}, error) {
		_, err := r.grades.DeleteOne(sc, bson.M{"schoolId": schoolID, "sourceId": sourceID, "studentId": studentID, "source": "manual"})
		if err != nil {
			return nil, err
		}
		if audit != nil {
			_, err := r.audits.InsertOne(sc, audit)
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	})
	return err
}

func (r *gradeRepoImpl) FindManualItemGrade(ctx context.Context, schoolID, sourceID, studentID bson.ObjectID) (*model.Grade, error) {
	var g model.Grade
	err := r.grades.FindOne(ctx, bson.M{"schoolId": schoolID, "sourceId": sourceID, "studentId": studentID}).Decode(&g)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &g, nil
}

func (r *gradeRepoImpl) FindLearningOutcomesForRecap(ctx context.Context, schoolID, academicYearID, subjectID, classID bson.ObjectID, semester int) ([]*model.LearningOutcome, error) {
	filter := bson.M{
		"schoolId":       schoolID,
		"academicYearId": academicYearID,
		"subjectId":      subjectID,
		"classIds":       classID,
		"semester":       semester,
	}
	opts := options.Find().SetSort(bson.D{{Key: "order", Value: 1}, {Key: "createdAt", Value: 1}})
	cursor, err := r.learningOutcomes.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var list []*model.LearningOutcome
	if err := cursor.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *gradeRepoImpl) FindItemsForRecap(ctx context.Context, schoolID, academicYearID, subjectID, classID bson.ObjectID) ([]*model.Assignment, []*model.Project, []*model.Assessment, error) {
	asgCursor, err := r.assignments.Find(ctx, bson.M{"schoolId": schoolID, "subjectId": subjectID, "classIds": classID})
	if err != nil {
		return nil, nil, nil, err
	}
	defer asgCursor.Close(ctx)
	var assignments []*model.Assignment
	_ = asgCursor.All(ctx, &assignments)

	prjCursor, err := r.projects.Find(ctx, bson.M{"schoolId": schoolID, "subjectId": subjectID, "classIds": classID})
	if err != nil {
		return nil, nil, nil, err
	}
	defer prjCursor.Close(ctx)
	var projects []*model.Project
	_ = prjCursor.All(ctx, &projects)

	asmCursor, err := r.assessments.Find(ctx, bson.M{"schoolId": schoolID, "subjectId": subjectID, "classIds": classID, "academicYearId": academicYearID, "category": model.CategoryDaily})
	if err != nil {
		return nil, nil, nil, err
	}
	defer asmCursor.Close(ctx)
	var assessments []*model.Assessment
	_ = asmCursor.All(ctx, &assessments)

	return assignments, projects, assessments, nil
}

