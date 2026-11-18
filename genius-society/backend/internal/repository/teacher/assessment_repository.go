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

// AssessmentRepository defines database operations for assessments, questions, attempts, and answers.
type AssessmentRepository interface {
	FindAssessments(ctx context.Context, schoolID, teacherID bson.ObjectID, classID, subjectID *bson.ObjectID, search string, page, limit int64) ([]*model.Assessment, int64, error)
	FindAssessmentByID(ctx context.Context, schoolID, teacherID, id bson.ObjectID) (*model.Assessment, error)
	CreateAssessment(ctx context.Context, a *model.Assessment) error
	UpdateAssessment(ctx context.Context, a *model.Assessment) error
	DeleteAssessment(ctx context.Context, schoolID, teacherID, id bson.ObjectID) error
	CountAttemptsByAssessmentID(ctx context.Context, schoolID, assessmentID bson.ObjectID) (int64, error)

	FindQuestions(ctx context.Context, assessmentID bson.ObjectID) ([]*model.AssessmentQuestion, error)
	FindQuestionByID(ctx context.Context, assessmentID, questionID bson.ObjectID) (*model.AssessmentQuestion, error)
	CreateQuestion(ctx context.Context, q *model.AssessmentQuestion) error
	UpdateQuestion(ctx context.Context, q *model.AssessmentQuestion) error
	DeleteQuestion(ctx context.Context, assessmentID, questionID bson.ObjectID) error
	ReorderQuestions(ctx context.Context, assessmentID bson.ObjectID, orders map[bson.ObjectID]int) error

	FindAttempts(ctx context.Context, schoolID, assessmentID bson.ObjectID) ([]*model.AssessmentAttempt, error)
	FindAttemptByID(ctx context.Context, schoolID, attemptID bson.ObjectID) (*model.AssessmentAttempt, error)
	UpdateAttempt(ctx context.Context, attempt *model.AssessmentAttempt) error

	FindAnswersByAttemptID(ctx context.Context, attemptID bson.ObjectID) ([]*model.AssessmentAnswer, error)
	FindAnswerByID(ctx context.Context, answerID bson.ObjectID) (*model.AssessmentAnswer, error)

	FindStudentsByClassIDs(ctx context.Context, schoolID bson.ObjectID, classIDs []bson.ObjectID) ([]*model.Student, error)
	FindStudentByID(ctx context.Context, schoolID, studentID bson.ObjectID) (*model.Student, error)
	FindUsersMap(ctx context.Context, userIDs []bson.ObjectID) (map[bson.ObjectID]*model.User, error)
	FindClassesMap(ctx context.Context, schoolID bson.ObjectID, classIDs []bson.ObjectID) (map[bson.ObjectID]*model.Class, error)
	FindSubjectsMap(ctx context.Context, subjectIDs []bson.ObjectID) (map[bson.ObjectID]*model.Subject, error)

	GradeEssayInTransaction(ctx context.Context, schoolID, teacherID bson.ObjectID, answer *model.AssessmentAnswer, attempt *model.AssessmentAttempt, grade *model.Grade, audit *model.AuditLog) error
	RegradeMCQInTransaction(ctx context.Context, schoolID, teacherID bson.ObjectID, question *model.AssessmentQuestion, attempts []*model.AssessmentAttempt, answers []*model.AssessmentAnswer, grades []*model.Grade, audit *model.AuditLog) error
}

type assessmentRepoImpl struct {
	db          *mongo.Database
	assessments *mongo.Collection
	questions   *mongo.Collection
	attempts    *mongo.Collection
	answers     *mongo.Collection
	students    *mongo.Collection
	users       *mongo.Collection
	classes     *mongo.Collection
	subjects    *mongo.Collection
	grades      *mongo.Collection
	auditLogs   *mongo.Collection
}

// NewAssessmentRepository creates a new AssessmentRepository instance.
func NewAssessmentRepository(db *mongo.Database) AssessmentRepository {
	return &assessmentRepoImpl{
		db:          db,
		assessments: db.Collection("assessments"),
		questions:   db.Collection("assessment_questions"),
		attempts:    db.Collection("assessment_attempts"),
		answers:     db.Collection("assessment_answers"),
		students:    db.Collection("students"),
		users:       db.Collection("users"),
		classes:     db.Collection("classes"),
		subjects:    db.Collection("subjects"),
		grades:      db.Collection("grades"),
		auditLogs:   db.Collection("audit_logs"),
	}
}

func (r *assessmentRepoImpl) FindAssessments(ctx context.Context, schoolID, teacherID bson.ObjectID, classID, subjectID *bson.ObjectID, search string, page, limit int64) ([]*model.Assessment, int64, error) {
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

	total, err := r.assessments.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	skip := (page - 1) * limit
	opts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}}).
		SetSkip(skip).
		SetLimit(limit)

	cur, err := r.assessments.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(ctx)

	var list []*model.Assessment
	if err := cur.All(ctx, &list); err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *assessmentRepoImpl) FindAssessmentByID(ctx context.Context, schoolID, teacherID, id bson.ObjectID) (*model.Assessment, error) {
	filter := bson.M{
		"_id":       id,
		"schoolId":  schoolID,
		"teacherId": teacherID,
	}
	var a model.Assessment
	err := r.assessments.FindOne(ctx, filter).Decode(&a)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}

func (r *assessmentRepoImpl) CreateAssessment(ctx context.Context, a *model.Assessment) error {
	res, err := r.assessments.InsertOne(ctx, a)
	if err != nil {
		return err
	}
	if oid, ok := res.InsertedID.(bson.ObjectID); ok {
		a.ID = oid
	}
	return nil
}

func (r *assessmentRepoImpl) UpdateAssessment(ctx context.Context, a *model.Assessment) error {
	filter := bson.M{
		"_id":       a.ID,
		"schoolId":  a.SchoolID,
		"teacherId": a.TeacherID,
	}
	_, err := r.assessments.ReplaceOne(ctx, filter, a)
	return err
}

func (r *assessmentRepoImpl) DeleteAssessment(ctx context.Context, schoolID, teacherID, id bson.ObjectID) error {
	filter := bson.M{
		"_id":       id,
		"schoolId":  schoolID,
		"teacherId": teacherID,
	}
	res, err := r.assessments.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	// Also delete associated questions
	_, _ = r.questions.DeleteMany(ctx, bson.M{"assessmentId": id})
	return nil
}

func (r *assessmentRepoImpl) CountAttemptsByAssessmentID(ctx context.Context, schoolID, assessmentID bson.ObjectID) (int64, error) {
	filter := bson.M{
		"assessmentId": assessmentID,
		"schoolId":     schoolID,
	}
	return r.attempts.CountDocuments(ctx, filter)
}

func (r *assessmentRepoImpl) FindQuestions(ctx context.Context, assessmentID bson.ObjectID) ([]*model.AssessmentQuestion, error) {
	filter := bson.M{"assessmentId": assessmentID}
	opts := options.Find().SetSort(bson.D{{Key: "order", Value: 1}, {Key: "createdAt", Value: 1}})
	cur, err := r.questions.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var list []*model.AssessmentQuestion
	if err := cur.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *assessmentRepoImpl) FindQuestionByID(ctx context.Context, assessmentID, questionID bson.ObjectID) (*model.AssessmentQuestion, error) {
	filter := bson.M{
		"_id":          questionID,
		"assessmentId": assessmentID,
	}
	var q model.AssessmentQuestion
	err := r.questions.FindOne(ctx, filter).Decode(&q)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &q, nil
}

func (r *assessmentRepoImpl) CreateQuestion(ctx context.Context, q *model.AssessmentQuestion) error {
	res, err := r.questions.InsertOne(ctx, q)
	if err != nil {
		return err
	}
	if oid, ok := res.InsertedID.(bson.ObjectID); ok {
		q.ID = oid
	}
	return nil
}

func (r *assessmentRepoImpl) UpdateQuestion(ctx context.Context, q *model.AssessmentQuestion) error {
	filter := bson.M{
		"_id":          q.ID,
		"assessmentId": q.AssessmentID,
	}
	_, err := r.questions.ReplaceOne(ctx, filter, q)
	return err
}

func (r *assessmentRepoImpl) DeleteQuestion(ctx context.Context, assessmentID, questionID bson.ObjectID) error {
	filter := bson.M{
		"_id":          questionID,
		"assessmentId": assessmentID,
	}
	res, err := r.questions.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *assessmentRepoImpl) ReorderQuestions(ctx context.Context, assessmentID bson.ObjectID, orders map[bson.ObjectID]int) error {
	for qID, newOrder := range orders {
		filter := bson.M{
			"_id":          qID,
			"assessmentId": assessmentID,
		}
		update := bson.M{
			"$set": bson.M{
				"order":     newOrder,
				"updatedAt": time.Now().UTC(),
			},
		}
		if _, err := r.questions.UpdateOne(ctx, filter, update); err != nil {
			return err
		}
	}
	return nil
}

func (r *assessmentRepoImpl) FindAttempts(ctx context.Context, schoolID, assessmentID bson.ObjectID) ([]*model.AssessmentAttempt, error) {
	filter := bson.M{
		"assessmentId": assessmentID,
		"schoolId":     schoolID,
	}
	opts := options.Find().SetSort(bson.D{{Key: "attemptNumber", Value: -1}, {Key: "createdAt", Value: -1}})
	cur, err := r.attempts.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var list []*model.AssessmentAttempt
	if err := cur.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *assessmentRepoImpl) FindAttemptByID(ctx context.Context, schoolID, attemptID bson.ObjectID) (*model.AssessmentAttempt, error) {
	filter := bson.M{
		"_id":      attemptID,
		"schoolId": schoolID,
	}
	var at model.AssessmentAttempt
	err := r.attempts.FindOne(ctx, filter).Decode(&at)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &at, nil
}

func (r *assessmentRepoImpl) UpdateAttempt(ctx context.Context, attempt *model.AssessmentAttempt) error {
	filter := bson.M{
		"_id":      attempt.ID,
		"schoolId": attempt.SchoolID,
	}
	_, err := r.attempts.ReplaceOne(ctx, filter, attempt)
	return err
}

func (r *assessmentRepoImpl) FindAnswersByAttemptID(ctx context.Context, attemptID bson.ObjectID) ([]*model.AssessmentAnswer, error) {
	filter := bson.M{"attemptId": attemptID}
	cur, err := r.answers.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var list []*model.AssessmentAnswer
	if err := cur.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *assessmentRepoImpl) FindAnswerByID(ctx context.Context, answerID bson.ObjectID) (*model.AssessmentAnswer, error) {
	filter := bson.M{"_id": answerID}
	var a model.AssessmentAnswer
	err := r.answers.FindOne(ctx, filter).Decode(&a)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}

func (r *assessmentRepoImpl) FindStudentsByClassIDs(ctx context.Context, schoolID bson.ObjectID, classIDs []bson.ObjectID) ([]*model.Student, error) {
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

func (r *assessmentRepoImpl) FindStudentByID(ctx context.Context, schoolID, studentID bson.ObjectID) (*model.Student, error) {
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

func (r *assessmentRepoImpl) FindUsersMap(ctx context.Context, userIDs []bson.ObjectID) (map[bson.ObjectID]*model.User, error) {
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

func (r *assessmentRepoImpl) FindClassesMap(ctx context.Context, schoolID bson.ObjectID, classIDs []bson.ObjectID) (map[bson.ObjectID]*model.Class, error) {
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

func (r *assessmentRepoImpl) FindSubjectsMap(ctx context.Context, subjectIDs []bson.ObjectID) (map[bson.ObjectID]*model.Subject, error) {
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

func (r *assessmentRepoImpl) GradeEssayInTransaction(ctx context.Context, schoolID, teacherID bson.ObjectID, answer *model.AssessmentAnswer, attempt *model.AssessmentAttempt, grade *model.Grade, audit *model.AuditLog) error {
	executeOps := func(sc context.Context) error {
		// 1. Update answer
		now := time.Now().UTC()
		ansFilter := bson.M{"_id": answer.ID}
		ansUpdate := bson.M{
			"$set": bson.M{
				"score":     answer.Score,
				"feedback":  answer.Feedback,
				"updatedAt": now,
			},
		}
		if _, err := r.answers.UpdateOne(sc, ansFilter, ansUpdate); err != nil {
			return err
		}

		// 2. Update attempt
		attFilter := bson.M{"_id": attempt.ID, "schoolId": schoolID}
		attUpdate := bson.M{
			"$set": bson.M{
				"manualScore": attempt.ManualScore,
				"totalScore":  attempt.TotalScore,
				"status":      attempt.Status,
				"updatedAt":   now,
			},
		}
		if _, err := r.attempts.UpdateOne(sc, attFilter, attUpdate); err != nil {
			return err
		}

		// 3. Upsert grade if attempt is graded
		if grade != nil {
			gradeFilter := bson.M{
				"studentId":      grade.StudentID,
				"subjectId":      grade.SubjectID,
				"academicYearId": grade.AcademicYearID,
				"semester":       grade.Semester,
				"type":           grade.Type,
				"sourceId":       grade.SourceID,
			}
			gradeUpdate := bson.M{
				"$set": bson.M{
					"score":     grade.Score,
					"schoolId":  schoolID,
					"teacherId": teacherID,
					"classId":   grade.ClassID,
					"source":    grade.Source,
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
		}

		// 4. Audit log
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

func (r *assessmentRepoImpl) RegradeMCQInTransaction(ctx context.Context, schoolID, teacherID bson.ObjectID, question *model.AssessmentQuestion, attempts []*model.AssessmentAttempt, answers []*model.AssessmentAnswer, grades []*model.Grade, audit *model.AuditLog) error {
	executeOps := func(sc context.Context) error {
		now := time.Now().UTC()

		// 1. Update question correctOption
		qFilter := bson.M{"_id": question.ID}
		qUpdate := bson.M{
			"$set": bson.M{
				"correctOption": question.CorrectOption,
				"updatedAt":     now,
			},
		}
		if _, err := r.questions.UpdateOne(sc, qFilter, qUpdate); err != nil {
			return err
		}

		// 2. Update all answers
		for _, ans := range answers {
			ansFilter := bson.M{"_id": ans.ID}
			ansUpdate := bson.M{
				"$set": bson.M{
					"isCorrect": ans.IsCorrect,
					"score":     ans.Score,
					"updatedAt": now,
				},
			}
			if _, err := r.answers.UpdateOne(sc, ansFilter, ansUpdate); err != nil {
				return err
			}
		}

		// 3. Update all attempts
		for _, att := range attempts {
			attFilter := bson.M{"_id": att.ID}
			attUpdate := bson.M{
				"$set": bson.M{
					"autoScore":  att.AutoScore,
					"totalScore": att.TotalScore,
					"updatedAt":  now,
				},
			}
			if _, err := r.attempts.UpdateOne(sc, attFilter, attUpdate); err != nil {
				return err
			}
		}

		// 4. Update all grades
		for _, gr := range grades {
			gradeFilter := bson.M{
				"studentId":      gr.StudentID,
				"subjectId":      gr.SubjectID,
				"academicYearId": gr.AcademicYearID,
				"semester":       gr.Semester,
				"type":           gr.Type,
				"sourceId":       gr.SourceID,
			}
			gradeUpdate := bson.M{
				"$set": bson.M{
					"score":     gr.Score,
					"updatedAt": now,
				},
			}
			if _, err := r.grades.UpdateOne(sc, gradeFilter, gradeUpdate); err != nil {
				return err
			}
		}

		// 5. Audit log
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
