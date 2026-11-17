package teacher

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/geniussociety/backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// AssignmentRepository defines database operations for assignments and submissions.
type AssignmentRepository interface {
	FindAssignments(ctx context.Context, schoolID, teacherID bson.ObjectID, classID, subjectID *bson.ObjectID, search string, page, limit int64) ([]*model.Assignment, int64, error)
	FindAssignmentByID(ctx context.Context, schoolID, teacherID, id bson.ObjectID) (*model.Assignment, error)
	CreateAssignment(ctx context.Context, a *model.Assignment) error
	UpdateAssignment(ctx context.Context, a *model.Assignment) error
	DeleteAssignment(ctx context.Context, schoolID, teacherID, id bson.ObjectID) error
	CountSubmissionsByAssignmentID(ctx context.Context, schoolID, assignmentID bson.ObjectID) (int64, error)
	FindSubmissions(ctx context.Context, schoolID, assignmentID bson.ObjectID) ([]*model.AssignmentSubmission, error)
	FindSubmissionByID(ctx context.Context, schoolID, submissionID bson.ObjectID) (*model.AssignmentSubmission, error)
	FindStudentsByClassIDs(ctx context.Context, schoolID bson.ObjectID, classIDs []bson.ObjectID) ([]*model.Student, error)
	FindStudentByID(ctx context.Context, schoolID, studentID bson.ObjectID) (*model.Student, error)
	FindUsersMap(ctx context.Context, userIDs []bson.ObjectID) (map[bson.ObjectID]*model.User, error)
	FindClassesMap(ctx context.Context, schoolID bson.ObjectID, classIDs []bson.ObjectID) (map[bson.ObjectID]*model.Class, error)
	FindSubjectsMap(ctx context.Context, subjectIDs []bson.ObjectID) (map[bson.ObjectID]*model.Subject, error)
	GradeSubmissionInTransaction(ctx context.Context, schoolID, teacherID bson.ObjectID, submission *model.AssignmentSubmission, grade *model.Grade, audit *model.AuditLog) error
}

type assignmentRepoImpl struct {
	db          *mongo.Database
	assignments *mongo.Collection
	submissions *mongo.Collection
	students    *mongo.Collection
	users       *mongo.Collection
	classes     *mongo.Collection
	subjects    *mongo.Collection
	grades      *mongo.Collection
	auditLogs   *mongo.Collection
}

// NewAssignmentRepository creates a new AssignmentRepository instance.
func NewAssignmentRepository(db *mongo.Database) AssignmentRepository {
	return &assignmentRepoImpl{
		db:          db,
		assignments: db.Collection("assignments"),
		submissions: db.Collection("assignment_submissions"),
		students:    db.Collection("students"),
		users:       db.Collection("users"),
		classes:     db.Collection("classes"),
		subjects:    db.Collection("subjects"),
		grades:      db.Collection("grades"),
		auditLogs:   db.Collection("audit_logs"),
	}
}

func (r *assignmentRepoImpl) FindAssignments(ctx context.Context, schoolID, teacherID bson.ObjectID, classID, subjectID *bson.ObjectID, search string, page, limit int64) ([]*model.Assignment, int64, error) {
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

	total, err := r.assignments.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	skip := (page - 1) * limit
	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}}).
		SetSkip(skip).
		SetLimit(limit)

	cur, err := r.assignments.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(ctx)

	var list []*model.Assignment
	if err := cur.All(ctx, &list); err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *assignmentRepoImpl) FindAssignmentByID(ctx context.Context, schoolID, teacherID, id bson.ObjectID) (*model.Assignment, error) {
	filter := bson.M{
		"_id":       id,
		"schoolId":  schoolID,
		"teacherId": teacherID,
	}
	var a model.Assignment
	err := r.assignments.FindOne(ctx, filter).Decode(&a)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}

func (r *assignmentRepoImpl) CreateAssignment(ctx context.Context, a *model.Assignment) error {
	res, err := r.assignments.InsertOne(ctx, a)
	if err != nil {
		return err
	}
	if oid, ok := res.InsertedID.(bson.ObjectID); ok {
		a.ID = oid
	}
	return nil
}

func (r *assignmentRepoImpl) UpdateAssignment(ctx context.Context, a *model.Assignment) error {
	filter := bson.M{
		"_id":       a.ID,
		"schoolId":  a.SchoolID,
		"teacherId": a.TeacherID,
	}
	_, err := r.assignments.ReplaceOne(ctx, filter, a)
	return err
}

func (r *assignmentRepoImpl) DeleteAssignment(ctx context.Context, schoolID, teacherID, id bson.ObjectID) error {
	filter := bson.M{
		"_id":       id,
		"schoolId":  schoolID,
		"teacherId": teacherID,
	}
	res, err := r.assignments.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *assignmentRepoImpl) CountSubmissionsByAssignmentID(ctx context.Context, schoolID, assignmentID bson.ObjectID) (int64, error) {
	filter := bson.M{
		"assignmentId": assignmentID,
		"schoolId":     schoolID,
	}
	return r.submissions.CountDocuments(ctx, filter)
}

func (r *assignmentRepoImpl) FindSubmissions(ctx context.Context, schoolID, assignmentID bson.ObjectID) ([]*model.AssignmentSubmission, error) {
	filter := bson.M{
		"assignmentId": assignmentID,
		"schoolId":     schoolID,
	}
	opts := options.Find().SetSort(bson.D{{Key: "version", Value: -1}})
	cur, err := r.submissions.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var list []*model.AssignmentSubmission
	if err := cur.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *assignmentRepoImpl) FindSubmissionByID(ctx context.Context, schoolID, submissionID bson.ObjectID) (*model.AssignmentSubmission, error) {
	filter := bson.M{
		"_id":      submissionID,
		"schoolId": schoolID,
	}
	var sub model.AssignmentSubmission
	err := r.submissions.FindOne(ctx, filter).Decode(&sub)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &sub, nil
}

func (r *assignmentRepoImpl) FindStudentsByClassIDs(ctx context.Context, schoolID bson.ObjectID, classIDs []bson.ObjectID) ([]*model.Student, error) {
	if len(classIDs) == 0 {
		return []*model.Student{}, nil
	}
	filter := bson.M{
		"schoolId": schoolID,
		"classId":  bson.M{"$in": classIDs},
		"status":   "active",
	}
	cur, err := r.students.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var list []*model.Student
	if err := cur.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *assignmentRepoImpl) FindStudentByID(ctx context.Context, schoolID, studentID bson.ObjectID) (*model.Student, error) {
	filter := bson.M{
		"_id":      studentID,
		"schoolId": schoolID,
	}
	var s model.Student
	err := r.students.FindOne(ctx, filter).Decode(&s)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *assignmentRepoImpl) FindUsersMap(ctx context.Context, userIDs []bson.ObjectID) (map[bson.ObjectID]*model.User, error) {
	result := make(map[bson.ObjectID]*model.User)
	if len(userIDs) == 0 {
		return result, nil
	}
	filter := bson.M{"_id": bson.M{"$in": userIDs}}
	cur, err := r.users.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var u model.User
		if err := cur.Decode(&u); err == nil {
			result[u.ID] = &u
		}
	}
	return result, nil
}

func (r *assignmentRepoImpl) FindClassesMap(ctx context.Context, schoolID bson.ObjectID, classIDs []bson.ObjectID) (map[bson.ObjectID]*model.Class, error) {
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

func (r *assignmentRepoImpl) FindSubjectsMap(ctx context.Context, subjectIDs []bson.ObjectID) (map[bson.ObjectID]*model.Subject, error) {
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

func (r *assignmentRepoImpl) GradeSubmissionInTransaction(ctx context.Context, schoolID, teacherID bson.ObjectID, submission *model.AssignmentSubmission, grade *model.Grade, audit *model.AuditLog) error {
	executeOps := func(sc context.Context) error {
		// 1. Update assignment submission
		subFilter := bson.M{
			"_id":      submission.ID,
			"schoolId": schoolID,
		}
		subUpdate := bson.M{
			"$set": bson.M{
				"score":     submission.Score,
				"feedback":  submission.Feedback,
				"gradedAt":  submission.GradedAt,
				"updatedAt": time.Now().UTC(),
			},
		}
		if _, err := r.submissions.UpdateOne(sc, subFilter, subUpdate); err != nil {
			return err
		}

		// 2. Upsert grade
		gradeFilter := bson.M{
			"studentId":      grade.StudentID,
			"subjectId":      grade.SubjectID,
			"academicYearId": grade.AcademicYearID,
			"semester":       grade.Semester,
			"type":           grade.Type,
			"sourceId":       grade.SourceID,
		}
		now := time.Now().UTC()
		gradeUpdate := bson.M{
			"$set": bson.M{
				"score":     grade.Score,
				"schoolId":  schoolID,
				"teacherId": teacherID,
				"classId":   grade.ClassID,
				"source":    grade.Source,
				"reason":    grade.Reason,
				"gradedBy":  teacherID,
				"updatedAt": now,
			},
			"$setOnInsert": bson.M{
				"_id":       bson.NewObjectID(),
				"createdAt": now,
			},
		}
		opts := options.UpdateOne().SetUpsert(true)
		if _, err := r.grades.UpdateOne(sc, gradeFilter, gradeUpdate, opts); err != nil {
			return err
		}

		// 3. Write audit log
		if audit != nil {
			if _, err := r.auditLogs.InsertOne(sc, audit); err != nil {
				return err
			}
		}

		return nil
	}

	session, err := r.db.Client().StartSession()
	if err != nil {
		// If sessions not supported, execute directly
		return executeOps(ctx)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sc context.Context) (interface{}, error) {
		return nil, executeOps(sc)
	})
	return err
}
