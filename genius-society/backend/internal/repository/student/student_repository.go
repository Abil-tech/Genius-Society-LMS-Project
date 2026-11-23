package studentrepo

import (
	"context"
	"time"

	"github.com/geniussociety/backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// StudentRepository provides database queries for all student portal operations.
type StudentRepository interface {
	// Scope & Profile
	FindUserByID(ctx context.Context, id bson.ObjectID) (*model.User, error)
	FindStudentByUserID(ctx context.Context, userID bson.ObjectID) (*model.Student, error)
	FindActiveAcademicYear(ctx context.Context, schoolID bson.ObjectID) (*model.AcademicYear, error)
	FindActiveEnrollment(ctx context.Context, schoolID, studentID, academicYearID bson.ObjectID) (*model.StudentEnrollment, error)
	FindClassByID(ctx context.Context, classID bson.ObjectID) (*model.Class, error)

	// Materials
	FindMaterialsByClass(ctx context.Context, schoolID, classID bson.ObjectID, subjectID *bson.ObjectID, now time.Time) ([]model.Material, error)
	FindMaterialByID(ctx context.Context, schoolID, materialID bson.ObjectID) (*model.Material, error)

	// Assignments
	FindAssignmentsByClass(ctx context.Context, schoolID, classID bson.ObjectID, subjectID *bson.ObjectID) ([]model.Assignment, error)
	FindAssignmentByID(ctx context.Context, schoolID, assignmentID bson.ObjectID) (*model.Assignment, error)
	FindAssignmentSubmissions(ctx context.Context, schoolID, assignmentID, studentID bson.ObjectID) ([]model.AssignmentSubmission, error)
	FindLatestAssignmentSubmission(ctx context.Context, schoolID, assignmentID, studentID bson.ObjectID) (*model.AssignmentSubmission, error)
	InsertAssignmentSubmission(ctx context.Context, sub *model.AssignmentSubmission) error
	UpdateAssignmentSubmission(ctx context.Context, sub *model.AssignmentSubmission) error

	// Projects
	FindProjectsByClass(ctx context.Context, schoolID, classID bson.ObjectID, subjectID *bson.ObjectID) ([]model.Project, error)
	FindProjectByID(ctx context.Context, schoolID, projectID bson.ObjectID) (*model.Project, error)
	FindProjectSubmissions(ctx context.Context, schoolID, projectID, studentID bson.ObjectID) ([]model.ProjectSubmission, error)
	FindLatestProjectSubmission(ctx context.Context, schoolID, projectID, studentID bson.ObjectID) (*model.ProjectSubmission, error)
	InsertProjectSubmission(ctx context.Context, sub *model.ProjectSubmission) error
	UpdateProjectSubmission(ctx context.Context, sub *model.ProjectSubmission) error

	// Assessments
	FindAssessmentsByClass(ctx context.Context, schoolID, classID bson.ObjectID, subjectID *bson.ObjectID) ([]model.Assessment, error)
	FindAssessmentByID(ctx context.Context, schoolID, assessmentID bson.ObjectID) (*model.Assessment, error)
	FindQuestionsByAssessmentID(ctx context.Context, assessmentID bson.ObjectID) ([]model.AssessmentQuestion, error)
	FindAttemptByStudent(ctx context.Context, schoolID, assessmentID, studentID bson.ObjectID) (*model.AssessmentAttempt, error)
	InsertAttempt(ctx context.Context, attempt *model.AssessmentAttempt) error
	UpdateAttempt(ctx context.Context, attempt *model.AssessmentAttempt) error
	FindAnswersByAttemptID(ctx context.Context, attemptID bson.ObjectID) ([]model.AssessmentAnswer, error)
	SaveAnswer(ctx context.Context, answer *model.AssessmentAnswer) error

	// Grades & Config
	FindGradesByStudent(ctx context.Context, schoolID, studentID, academicYearID bson.ObjectID) ([]model.Grade, error)
	FindGradingConfig(ctx context.Context, schoolID, academicYearID bson.ObjectID, subjectID *bson.ObjectID) (*model.GradingConfig, error)
	FindClassStudents(ctx context.Context, schoolID, classID bson.ObjectID) ([]model.Student, error)

	// Schedules
	FindSchedulesByClass(ctx context.Context, schoolID, classID, academicYearID bson.ObjectID) ([]model.Schedule, error)

	// Auxiliary maps
	FindSubjectsMap(ctx context.Context, subjectIDs []bson.ObjectID) (map[bson.ObjectID]*model.Subject, error)
	FindTeachersMap(ctx context.Context, teacherIDs []bson.ObjectID) (map[bson.ObjectID]string, error)
	FindUsersMap(ctx context.Context, userIDs []bson.ObjectID) (map[bson.ObjectID]*model.User, error)
}

type studentRepoImpl struct {
	db *mongo.Database
}

// NewStudentRepository constructs a new StudentRepository backed by MongoDB.
func NewStudentRepository(db *mongo.Database) StudentRepository {
	return &studentRepoImpl{db: db}
}

// FindUserByID finds user by ObjectID.
func (r *studentRepoImpl) FindUserByID(ctx context.Context, id bson.ObjectID) (*model.User, error) {
	var u model.User
	err := r.db.Collection("users").FindOne(ctx, bson.M{"_id": id}).Decode(&u)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

// FindStudentByUserID finds student record by userId.
func (r *studentRepoImpl) FindStudentByUserID(ctx context.Context, userID bson.ObjectID) (*model.Student, error) {
	var st model.Student
	err := r.db.Collection("students").FindOne(ctx, bson.M{"userId": userID}).Decode(&st)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &st, nil
}

// FindActiveAcademicYear finds active academic year for a school.
func (r *studentRepoImpl) FindActiveAcademicYear(ctx context.Context, schoolID bson.ObjectID) (*model.AcademicYear, error) {
	var ay model.AcademicYear
	err := r.db.Collection("academic_years").FindOne(ctx, bson.M{
		"schoolId": schoolID,
		"isActive": true,
	}).Decode(&ay)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &ay, nil
}

// FindActiveEnrollment finds enrollment with status "active".
func (r *studentRepoImpl) FindActiveEnrollment(ctx context.Context, schoolID, studentID, academicYearID bson.ObjectID) (*model.StudentEnrollment, error) {
	var en model.StudentEnrollment
	err := r.db.Collection("enrollments").FindOne(ctx, bson.M{
		"schoolId":       schoolID,
		"studentId":      studentID,
		"academicYearId": academicYearID,
		"status":         model.EnrollmentStatusActive,
	}).Decode(&en)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &en, nil
}

// FindClassByID retrieves a class by its ObjectID.
func (r *studentRepoImpl) FindClassByID(ctx context.Context, classID bson.ObjectID) (*model.Class, error) {
	var cl model.Class
	err := r.db.Collection("classes").FindOne(ctx, bson.M{"_id": classID}).Decode(&cl)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &cl, nil
}

// FindMaterialsByClass retrieves materials targeting the given class and published <= now.
func (r *studentRepoImpl) FindMaterialsByClass(ctx context.Context, schoolID, classID bson.ObjectID, subjectID *bson.ObjectID, now time.Time) ([]model.Material, error) {
	filter := bson.M{
		"schoolId": schoolID,
		"classIds": classID,
		"$or": []bson.M{
			{"publishAt": bson.M{"$lte": now}},
			{"publishAt": bson.M{"$exists": false}},
			{"publishAt": nil},
		},
	}
	if subjectID != nil && !subjectID.IsZero() {
		filter["subjectId"] = *subjectID
	}

	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cur, err := r.db.Collection("materials").Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var list []model.Material
	if err := cur.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

// FindMaterialByID finds material by ID and school.
func (r *studentRepoImpl) FindMaterialByID(ctx context.Context, schoolID, materialID bson.ObjectID) (*model.Material, error) {
	var m model.Material
	err := r.db.Collection("materials").FindOne(ctx, bson.M{
		"_id":      materialID,
		"schoolId": schoolID,
	}).Decode(&m)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

// FindAssignmentsByClass retrieves assignments targeting student's class.
func (r *studentRepoImpl) FindAssignmentsByClass(ctx context.Context, schoolID, classID bson.ObjectID, subjectID *bson.ObjectID) ([]model.Assignment, error) {
	filter := bson.M{
		"schoolId": schoolID,
		"classIds": classID,
	}
	if subjectID != nil && !subjectID.IsZero() {
		filter["subjectId"] = *subjectID
	}

	opts := options.Find().SetSort(bson.D{{Key: "deadline", Value: 1}})
	cur, err := r.db.Collection("assignments").Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var list []model.Assignment
	if err := cur.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

// FindAssignmentByID finds an assignment by ID.
func (r *studentRepoImpl) FindAssignmentByID(ctx context.Context, schoolID, assignmentID bson.ObjectID) (*model.Assignment, error) {
	var a model.Assignment
	err := r.db.Collection("assignments").FindOne(ctx, bson.M{
		"_id":      assignmentID,
		"schoolId": schoolID,
	}).Decode(&a)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}

// FindAssignmentSubmissions finds all submissions for a student on an assignment.
func (r *studentRepoImpl) FindAssignmentSubmissions(ctx context.Context, schoolID, assignmentID, studentID bson.ObjectID) ([]model.AssignmentSubmission, error) {
	filter := bson.M{
		"schoolId":     schoolID,
		"assignmentId": assignmentID,
		"studentId":    studentID,
	}
	opts := options.Find().SetSort(bson.D{{Key: "version", Value: -1}})
	cur, err := r.db.Collection("assignment_submissions").Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var list []model.AssignmentSubmission
	if err := cur.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

// FindLatestAssignmentSubmission gets the latest version submission for a student.
func (r *studentRepoImpl) FindLatestAssignmentSubmission(ctx context.Context, schoolID, assignmentID, studentID bson.ObjectID) (*model.AssignmentSubmission, error) {
	filter := bson.M{
		"schoolId":     schoolID,
		"assignmentId": assignmentID,
		"studentId":    studentID,
	}
	opts := options.FindOne().SetSort(bson.D{{Key: "version", Value: -1}})
	var sub model.AssignmentSubmission
	err := r.db.Collection("assignment_submissions").FindOne(ctx, filter, opts).Decode(&sub)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &sub, nil
}

// InsertAssignmentSubmission creates a new submission record.
func (r *studentRepoImpl) InsertAssignmentSubmission(ctx context.Context, sub *model.AssignmentSubmission) error {
	if sub.ID.IsZero() {
		sub.ID = bson.NewObjectID()
	}
	now := time.Now().UTC()
	sub.CreatedAt = now
	sub.UpdatedAt = now
	_, err := r.db.Collection("assignment_submissions").InsertOne(ctx, sub)
	return err
}

// UpdateAssignmentSubmission updates an existing submission.
func (r *studentRepoImpl) UpdateAssignmentSubmission(ctx context.Context, sub *model.AssignmentSubmission) error {
	sub.UpdatedAt = time.Now().UTC()
	filter := bson.M{"_id": sub.ID}
	_, err := r.db.Collection("assignment_submissions").ReplaceOne(ctx, filter, sub)
	return err
}

// FindProjectsByClass retrieves projects targeting student's class.
func (r *studentRepoImpl) FindProjectsByClass(ctx context.Context, schoolID, classID bson.ObjectID, subjectID *bson.ObjectID) ([]model.Project, error) {
	filter := bson.M{
		"schoolId": schoolID,
		"classIds": classID,
	}
	if subjectID != nil && !subjectID.IsZero() {
		filter["subjectId"] = *subjectID
	}

	opts := options.Find().SetSort(bson.D{{Key: "deadline", Value: 1}})
	cur, err := r.db.Collection("projects").Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var list []model.Project
	if err := cur.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

// FindProjectByID finds a project by ID.
func (r *studentRepoImpl) FindProjectByID(ctx context.Context, schoolID, projectID bson.ObjectID) (*model.Project, error) {
	var p model.Project
	err := r.db.Collection("projects").FindOne(ctx, bson.M{
		"_id":      projectID,
		"schoolId": schoolID,
	}).Decode(&p)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

// FindProjectSubmissions finds all submissions for a student on a project.
func (r *studentRepoImpl) FindProjectSubmissions(ctx context.Context, schoolID, projectID, studentID bson.ObjectID) ([]model.ProjectSubmission, error) {
	filter := bson.M{
		"schoolId":  schoolID,
		"projectId": projectID,
		"studentId": studentID,
	}
	opts := options.Find().SetSort(bson.D{{Key: "version", Value: -1}})
	cur, err := r.db.Collection("project_submissions").Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var list []model.ProjectSubmission
	if err := cur.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

// FindLatestProjectSubmission retrieves latest submission for a project.
func (r *studentRepoImpl) FindLatestProjectSubmission(ctx context.Context, schoolID, projectID, studentID bson.ObjectID) (*model.ProjectSubmission, error) {
	filter := bson.M{
		"schoolId":  schoolID,
		"projectId": projectID,
		"studentId": studentID,
	}
	opts := options.FindOne().SetSort(bson.D{{Key: "version", Value: -1}})
	var sub model.ProjectSubmission
	err := r.db.Collection("project_submissions").FindOne(ctx, filter, opts).Decode(&sub)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &sub, nil
}

// InsertProjectSubmission creates a new project submission record.
func (r *studentRepoImpl) InsertProjectSubmission(ctx context.Context, sub *model.ProjectSubmission) error {
	if sub.ID.IsZero() {
		sub.ID = bson.NewObjectID()
	}
	now := time.Now().UTC()
	sub.CreatedAt = now
	sub.UpdatedAt = now
	_, err := r.db.Collection("project_submissions").InsertOne(ctx, sub)
	return err
}

// UpdateProjectSubmission replaces an existing project submission.
func (r *studentRepoImpl) UpdateProjectSubmission(ctx context.Context, sub *model.ProjectSubmission) error {
	sub.UpdatedAt = time.Now().UTC()
	filter := bson.M{"_id": sub.ID}
	_, err := r.db.Collection("project_submissions").ReplaceOne(ctx, filter, sub)
	return err
}

// FindAssessmentsByClass retrieves published or closed assessments for the class.
func (r *studentRepoImpl) FindAssessmentsByClass(ctx context.Context, schoolID, classID bson.ObjectID, subjectID *bson.ObjectID) ([]model.Assessment, error) {
	filter := bson.M{
		"schoolId": schoolID,
		"classIds": classID,
		"status":   bson.M{"$in": []model.AssessmentStatus{model.AssessmentPublished, model.AssessmentClosed}},
	}
	if subjectID != nil && !subjectID.IsZero() {
		filter["subjectId"] = *subjectID
	}

	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cur, err := r.db.Collection("assessments").Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var list []model.Assessment
	if err := cur.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

// FindAssessmentByID finds an assessment by ID.
func (r *studentRepoImpl) FindAssessmentByID(ctx context.Context, schoolID, assessmentID bson.ObjectID) (*model.Assessment, error) {
	var a model.Assessment
	err := r.db.Collection("assessments").FindOne(ctx, bson.M{
		"_id":      assessmentID,
		"schoolId": schoolID,
	}).Decode(&a)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}

// FindQuestionsByAssessmentID finds questions sorted by order.
func (r *studentRepoImpl) FindQuestionsByAssessmentID(ctx context.Context, assessmentID bson.ObjectID) ([]model.AssessmentQuestion, error) {
	opts := options.Find().SetSort(bson.D{{Key: "order", Value: 1}})
	cur, err := r.db.Collection("assessment_questions").Find(ctx, bson.M{"assessmentId": assessmentID}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var list []model.AssessmentQuestion
	if err := cur.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

// FindAttemptByStudent finds existing student attempt on an assessment.
func (r *studentRepoImpl) FindAttemptByStudent(ctx context.Context, schoolID, assessmentID, studentID bson.ObjectID) (*model.AssessmentAttempt, error) {
	var att model.AssessmentAttempt
	err := r.db.Collection("assessment_attempts").FindOne(ctx, bson.M{
		"schoolId":     schoolID,
		"assessmentId": assessmentID,
		"studentId":    studentID,
	}).Decode(&att)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &att, nil
}

// InsertAttempt inserts a new assessment attempt.
func (r *studentRepoImpl) InsertAttempt(ctx context.Context, attempt *model.AssessmentAttempt) error {
	if attempt.ID.IsZero() {
		attempt.ID = bson.NewObjectID()
	}
	now := time.Now().UTC()
	attempt.CreatedAt = now
	attempt.UpdatedAt = now
	_, err := r.db.Collection("assessment_attempts").InsertOne(ctx, attempt)
	return err
}

// UpdateAttempt updates attempt state and score.
func (r *studentRepoImpl) UpdateAttempt(ctx context.Context, attempt *model.AssessmentAttempt) error {
	attempt.UpdatedAt = time.Now().UTC()
	_, err := r.db.Collection("assessment_attempts").ReplaceOne(ctx, bson.M{"_id": attempt.ID}, attempt)
	return err
}

// FindAnswersByAttemptID retrieves all answers recorded for an attempt.
func (r *studentRepoImpl) FindAnswersByAttemptID(ctx context.Context, attemptID bson.ObjectID) ([]model.AssessmentAnswer, error) {
	cur, err := r.db.Collection("assessment_answers").Find(ctx, bson.M{"attemptId": attemptID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var list []model.AssessmentAnswer
	if err := cur.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

// SaveAnswer upserts an answer for an attempt question.
func (r *studentRepoImpl) SaveAnswer(ctx context.Context, answer *model.AssessmentAnswer) error {
	filter := bson.M{
		"attemptId":  answer.AttemptID,
		"questionId": answer.QuestionID,
	}
	now := time.Now().UTC()
	answer.UpdatedAt = now
	if answer.CreatedAt.IsZero() {
		answer.CreatedAt = now
	}
	if answer.ID.IsZero() {
		answer.ID = bson.NewObjectID()
	}

	opts := options.UpdateOne().SetUpsert(true)
	update := bson.M{
		"$set": answer,
	}
	_, err := r.db.Collection("assessment_answers").UpdateOne(ctx, filter, update, opts)
	return err
}

// FindGradesByStudent fetches all grade entries for the student in current academic year.
func (r *studentRepoImpl) FindGradesByStudent(ctx context.Context, schoolID, studentID, academicYearID bson.ObjectID) ([]model.Grade, error) {
	filter := bson.M{
		"schoolId":       schoolID,
		"studentId":      studentID,
		"academicYearId": academicYearID,
	}
	cur, err := r.db.Collection("grades").Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var list []model.Grade
	if err := cur.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

// FindGradingConfig finds grading config for subject or fallback to school/ay default.
func (r *studentRepoImpl) FindGradingConfig(ctx context.Context, schoolID, academicYearID bson.ObjectID, subjectID *bson.ObjectID) (*model.GradingConfig, error) {
	if subjectID != nil && !subjectID.IsZero() {
		var cfg model.GradingConfig
		err := r.db.Collection("grading_configs").FindOne(ctx, bson.M{
			"schoolId":       schoolID,
			"academicYearId": academicYearID,
			"subjectId":      *subjectID,
		}).Decode(&cfg)
		if err == nil {
			return &cfg, nil
		}
	}

	var globalCfg model.GradingConfig
	err := r.db.Collection("grading_configs").FindOne(ctx, bson.M{
		"schoolId":       schoolID,
		"academicYearId": academicYearID,
		"$or": []bson.M{
			{"subjectId": nil},
			{"subjectId": bson.M{"$exists": false}},
		},
	}).Decode(&globalCfg)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &globalCfg, nil
}

// FindClassStudents returns all active students enrolled in class.
func (r *studentRepoImpl) FindClassStudents(ctx context.Context, schoolID, classID bson.ObjectID) ([]model.Student, error) {
	filter := bson.M{
		"schoolId": schoolID,
		"classId":  classID,
		"status":   "active",
	}
	cur, err := r.db.Collection("students").Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var list []model.Student
	if err := cur.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

// FindSchedulesByClass returns class schedules sorted by day and startTime.
func (r *studentRepoImpl) FindSchedulesByClass(ctx context.Context, schoolID, classID, academicYearID bson.ObjectID) ([]model.Schedule, error) {
	filter := bson.M{
		"schoolId":       schoolID,
		"classId":        classID,
		"academicYearId": academicYearID,
	}
	opts := options.Find().SetSort(bson.D{
		{Key: "day", Value: 1},
		{Key: "startTime", Value: 1},
	})
	cur, err := r.db.Collection("schedules").Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var list []model.Schedule
	if err := cur.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

// FindSubjectsMap builds a map of subjects by ObjectID.
func (r *studentRepoImpl) FindSubjectsMap(ctx context.Context, subjectIDs []bson.ObjectID) (map[bson.ObjectID]*model.Subject, error) {
	res := make(map[bson.ObjectID]*model.Subject)
	if len(subjectIDs) == 0 {
		return res, nil
	}
	cur, err := r.db.Collection("subjects").Find(ctx, bson.M{"_id": bson.M{"$in": subjectIDs}})
	if err != nil {
		return res, err
	}
	defer cur.Close(ctx)

	var list []model.Subject
	if err := cur.All(ctx, &list); err != nil {
		return res, err
	}
	for i := range list {
		res[list[i].ID] = &list[i]
	}
	return res, nil
}

// FindTeachersMap returns a map of teacher names indexed by teacherID.
func (r *studentRepoImpl) FindTeachersMap(ctx context.Context, teacherIDs []bson.ObjectID) (map[bson.ObjectID]string, error) {
	res := make(map[bson.ObjectID]string)
	if len(teacherIDs) == 0 {
		return res, nil
	}
	cur, err := r.db.Collection("teachers").Find(ctx, bson.M{"_id": bson.M{"$in": teacherIDs}})
	if err != nil {
		return res, err
	}
	defer cur.Close(ctx)

	var teachers []model.Teacher
	if err := cur.All(ctx, &teachers); err != nil {
		return res, err
	}

	var userIDs []bson.ObjectID
	for _, t := range teachers {
		if !t.UserID.IsZero() {
			userIDs = append(userIDs, t.UserID)
		}
	}

	usersMap, _ := r.FindUsersMap(ctx, userIDs)
	for _, t := range teachers {
		name := t.NIP
		if u, ok := usersMap[t.UserID]; ok && u != nil && u.Name != "" {
			name = u.Name
		}
		res[t.ID] = name
	}
	return res, nil
}

// FindUsersMap returns a map of users indexed by UserID.
func (r *studentRepoImpl) FindUsersMap(ctx context.Context, userIDs []bson.ObjectID) (map[bson.ObjectID]*model.User, error) {
	res := make(map[bson.ObjectID]*model.User)
	if len(userIDs) == 0 {
		return res, nil
	}
	cur, err := r.db.Collection("users").Find(ctx, bson.M{"_id": bson.M{"$in": userIDs}})
	if err != nil {
		return res, err
	}
	defer cur.Close(ctx)

	var list []model.User
	if err := cur.All(ctx, &list); err != nil {
		return res, err
	}
	for i := range list {
		res[list[i].ID] = &list[i]
	}
	return res, nil
}
