package teacherrepo

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

// ProjectRepository defines database operations for projects and project submissions.
type ProjectRepository interface {
	FindProjects(ctx context.Context, teacherID bson.ObjectID, classID, subjectID *bson.ObjectID, search string, page, limit int64) ([]*model.Project, int64, error)
	FindProjectByID(ctx context.Context, teacherID, id bson.ObjectID) (*model.Project, error)
	CreateProject(ctx context.Context, p *model.Project) error
	UpdateProject(ctx context.Context, p *model.Project) error
	DeleteProject(ctx context.Context, teacherID, id bson.ObjectID) error
	CountSubmissionsByProjectID(ctx context.Context, projectID bson.ObjectID) (int64, error)
	FindSubmissions(ctx context.Context, projectID bson.ObjectID) ([]*model.ProjectSubmission, error)
	FindSubmissionByID(ctx context.Context, submissionID bson.ObjectID) (*model.ProjectSubmission, error)
	FindStudentsByClassIDs(ctx context.Context, classIDs []bson.ObjectID) ([]*model.Student, error)
	FindStudentByID(ctx context.Context, studentID bson.ObjectID) (*model.Student, error)
	FindUsersMap(ctx context.Context, userIDs []bson.ObjectID) (map[bson.ObjectID]*model.User, error)
	FindClassesMap(ctx context.Context, classIDs []bson.ObjectID) (map[bson.ObjectID]*model.Class, error)
	FindSubjectsMap(ctx context.Context, subjectIDs []bson.ObjectID) (map[bson.ObjectID]*model.Subject, error)
	GradeSubmissionInTransaction(ctx context.Context, teacherID bson.ObjectID, submission *model.ProjectSubmission, grade *model.Grade, audit *model.AuditLog) error
}

type projectRepoImpl struct {
	db          *mongo.Database
	projects    *mongo.Collection
	submissions *mongo.Collection
	students    *mongo.Collection
	users       *mongo.Collection
	classes     *mongo.Collection
	subjects    *mongo.Collection
	grades      *mongo.Collection
	auditLogs   *mongo.Collection
}

// NewProjectRepository creates a new ProjectRepository instance.
func NewProjectRepository(db *mongo.Database) ProjectRepository {
	return &projectRepoImpl{
		db:          db,
		projects:    db.Collection("projects"),
		submissions: db.Collection("project_submissions"),
		students:    db.Collection("students"),
		users:       db.Collection("users"),
		classes:     db.Collection("classes"),
		subjects:    db.Collection("subjects"),
		grades:      db.Collection("grades"),
		auditLogs:   db.Collection("audit_logs"),
	}
}

func (r *projectRepoImpl) FindProjects(ctx context.Context, teacherID bson.ObjectID, classID, subjectID *bson.ObjectID, search string, page, limit int64) ([]*model.Project, int64, error) {
	filter := bson.M{
		
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

	total, err := r.projects.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	skip := (page - 1) * limit
	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}}).
		SetSkip(skip).
		SetLimit(limit)

	cur, err := r.projects.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(ctx)

	var list []*model.Project
	if err := cur.All(ctx, &list); err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *projectRepoImpl) FindProjectByID(ctx context.Context, teacherID, id bson.ObjectID) (*model.Project, error) {
	filter := bson.M{
		"_id":       id,
		
		"teacherId": teacherID,
	}
	var p model.Project
	err := r.projects.FindOne(ctx, filter).Decode(&p)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *projectRepoImpl) CreateProject(ctx context.Context, p *model.Project) error {
	res, err := r.projects.InsertOne(ctx, p)
	if err != nil {
		return err
	}
	if oid, ok := res.InsertedID.(bson.ObjectID); ok {
		p.ID = oid
	}
	return nil
}

func (r *projectRepoImpl) UpdateProject(ctx context.Context, p *model.Project) error {
	filter := bson.M{
		"_id":       p.ID,
		
		"teacherId": p.TeacherID,
	}
	_, err := r.projects.ReplaceOne(ctx, filter, p)
	return err
}

func (r *projectRepoImpl) DeleteProject(ctx context.Context, teacherID, id bson.ObjectID) error {
	filter := bson.M{
		"_id":       id,
		
		"teacherId": teacherID,
	}
	res, err := r.projects.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *projectRepoImpl) CountSubmissionsByProjectID(ctx context.Context, projectID bson.ObjectID) (int64, error) {
	filter := bson.M{
		"projectId": projectID,
		
	}
	return r.submissions.CountDocuments(ctx, filter)
}

func (r *projectRepoImpl) FindSubmissions(ctx context.Context, projectID bson.ObjectID) ([]*model.ProjectSubmission, error) {
	filter := bson.M{
		"projectId": projectID,
		
	}
	opts := options.Find().SetSort(bson.D{{Key: "version", Value: -1}})
	cur, err := r.submissions.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var list []*model.ProjectSubmission
	if err := cur.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *projectRepoImpl) FindSubmissionByID(ctx context.Context, submissionID bson.ObjectID) (*model.ProjectSubmission, error) {
	filter := bson.M{
		"_id":      submissionID,
		
	}
	var sub model.ProjectSubmission
	err := r.submissions.FindOne(ctx, filter).Decode(&sub)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &sub, nil
}

func (r *projectRepoImpl) FindStudentsByClassIDs(ctx context.Context, classIDs []bson.ObjectID) ([]*model.Student, error) {
	if len(classIDs) == 0 {
		return []*model.Student{}, nil
	}
	filter := bson.M{
		
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

func (r *projectRepoImpl) FindStudentByID(ctx context.Context, studentID bson.ObjectID) (*model.Student, error) {
	filter := bson.M{
		"_id":      studentID,
		
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

func (r *projectRepoImpl) FindUsersMap(ctx context.Context, userIDs []bson.ObjectID) (map[bson.ObjectID]*model.User, error) {
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

func (r *projectRepoImpl) FindClassesMap(ctx context.Context, classIDs []bson.ObjectID) (map[bson.ObjectID]*model.Class, error) {
	result := make(map[bson.ObjectID]*model.Class)
	if len(classIDs) == 0 {
		return result, nil
	}
	filter := bson.M{
		"_id":      bson.M{"$in": classIDs},
		
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

func (r *projectRepoImpl) FindSubjectsMap(ctx context.Context, subjectIDs []bson.ObjectID) (map[bson.ObjectID]*model.Subject, error) {
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

func (r *projectRepoImpl) GradeSubmissionInTransaction(ctx context.Context, teacherID bson.ObjectID, submission *model.ProjectSubmission, grade *model.Grade, audit *model.AuditLog) error {
	executeOps := func(sc context.Context) error {
		// 1. Update project submission
		subFilter := bson.M{
			"_id":      submission.ID,
			
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
		return executeOps(ctx)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sc context.Context) (interface{}, error) {
		return nil, executeOps(sc)
	})
	return err
}
