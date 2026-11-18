package teacherrepo

import (
	"context"
	"time"

	"github.com/geniussociety/backend/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// DashboardRepository defines read operations for aggregating teacher dashboard data.
type DashboardRepository interface {
	// CountStudentsByClassIDs returns total active students across the given classes.
	CountStudentsByClassIDs(ctx context.Context, schoolID bson.ObjectID, classIDs []bson.ObjectID) (int, error)

	// CountActiveAssignments returns assignments owned by teacher whose deadline is still in the future.
	CountActiveAssignments(ctx context.Context, teacherID bson.ObjectID) (int, error)

	// CountActiveAssessments returns published assessments owned by teacher.
	CountActiveAssessments(ctx context.Context, teacherID bson.ObjectID) (int, error)

	// CountActiveProjects returns projects owned by teacher whose deadline is still in the future.
	CountActiveProjects(ctx context.Context, teacherID bson.ObjectID) (int, error)

	// CountUngradedAssignmentSubmissions counts assignment submissions (score == null) for teacher's assignments.
	CountUngradedAssignmentSubmissions(ctx context.Context, teacherID bson.ObjectID) (int, error)

	// CountUngradedProjectSubmissions counts project submissions (score == null) for teacher's projects.
	CountUngradedProjectSubmissions(ctx context.Context, teacherID bson.ObjectID) (int, error)

	// CountUngradedAssessmentAttempts counts assessment attempts with status "submitted" for teacher's assessments.
	CountUngradedAssessmentAttempts(ctx context.Context, teacherID bson.ObjectID) (int, error)

	// FindTodaySchedules returns schedules matching teacher + academic year for a given day-of-week.
	FindTodaySchedules(ctx context.Context, teacherID, academicYearID bson.ObjectID, dayOfWeek int) ([]*model.Schedule, error)

	// FindRecentSubmissions returns the most recent assignment + project submissions for teacher's items.
	// Results are sorted by submittedAt desc, limited to `limit`.
	FindRecentAssignmentSubmissions(ctx context.Context, teacherID bson.ObjectID, limit int) ([]*model.AssignmentSubmission, error)
	FindRecentProjectSubmissions(ctx context.Context, teacherID bson.ObjectID, limit int) ([]*model.ProjectSubmission, error)

	// FindAssignmentsByIDs returns assignments keyed by their ID.
	FindAssignmentsByIDs(ctx context.Context, ids []bson.ObjectID) (map[bson.ObjectID]*model.Assignment, error)

	// FindProjectsByIDs returns projects keyed by their ID.
	FindProjectsByIDs(ctx context.Context, ids []bson.ObjectID) (map[bson.ObjectID]*model.Project, error)

	// FindStudentsByIDs returns students keyed by their ID.
	FindStudentsByIDs(ctx context.Context, ids []bson.ObjectID) (map[bson.ObjectID]*model.Student, error)

	// FindClassByID returns a single class.
	FindClassByID(ctx context.Context, classID bson.ObjectID) (*model.Class, error)

	// FindSubjectsByIDs returns subjects keyed by their ID.
	FindSubjectsByIDs(ctx context.Context, ids []bson.ObjectID) (map[bson.ObjectID]*model.Subject, error)

	// FindClassesByIDs returns classes keyed by their ID.
	FindClassesByIDs(ctx context.Context, ids []bson.ObjectID) (map[bson.ObjectID]*model.Class, error)

	// CountStudentsByClassID returns number of active students in a single class.
	CountStudentsByClassID(ctx context.Context, schoolID, classID bson.ObjectID) (int, error)
}

type dashboardRepoImpl struct {
	students               *mongo.Collection
	assignments            *mongo.Collection
	assignmentSubmissions  *mongo.Collection
	assessments            *mongo.Collection
	assessmentAttempts     *mongo.Collection
	projects               *mongo.Collection
	projectSubmissions     *mongo.Collection
	schedules              *mongo.Collection
	classes                *mongo.Collection
	subjects               *mongo.Collection
}

// NewDashboardRepository creates a new DashboardRepository instance.
func NewDashboardRepository(db *mongo.Database) DashboardRepository {
	return &dashboardRepoImpl{
		students:              db.Collection("students"),
		assignments:           db.Collection("assignments"),
		assignmentSubmissions: db.Collection("assignment_submissions"),
		assessments:           db.Collection("assessments"),
		assessmentAttempts:    db.Collection("assessment_attempts"),
		projects:              db.Collection("projects"),
		projectSubmissions:    db.Collection("project_submissions"),
		schedules:             db.Collection("schedules"),
		classes:               db.Collection("classes"),
		subjects:              db.Collection("subjects"),
	}
}

func (r *dashboardRepoImpl) CountStudentsByClassIDs(ctx context.Context, schoolID bson.ObjectID, classIDs []bson.ObjectID) (int, error) {
	if len(classIDs) == 0 {
		return 0, nil
	}
	filter := bson.M{
		"classId":  bson.M{"$in": classIDs},
		"schoolId": schoolID,
		"status":   "active",
	}
	count, err := r.students.CountDocuments(ctx, filter)
	return int(count), err
}

func (r *dashboardRepoImpl) CountActiveAssignments(ctx context.Context, teacherID bson.ObjectID) (int, error) {
	now := time.Now()
	filter := bson.M{
		"teacherId": teacherID,
		"deadline":  bson.M{"$gte": now},
	}
	count, err := r.assignments.CountDocuments(ctx, filter)
	return int(count), err
}

func (r *dashboardRepoImpl) CountActiveAssessments(ctx context.Context, teacherID bson.ObjectID) (int, error) {
	filter := bson.M{
		"teacherId": teacherID,
		"status":    string(model.AssessmentPublished),
	}
	count, err := r.assessments.CountDocuments(ctx, filter)
	return int(count), err
}

func (r *dashboardRepoImpl) CountActiveProjects(ctx context.Context, teacherID bson.ObjectID) (int, error) {
	now := time.Now()
	filter := bson.M{
		"teacherId": teacherID,
		"deadline":  bson.M{"$gte": now},
	}
	count, err := r.projects.CountDocuments(ctx, filter)
	return int(count), err
}

func (r *dashboardRepoImpl) CountUngradedAssignmentSubmissions(ctx context.Context, teacherID bson.ObjectID) (int, error) {
	// First get all assignment IDs for this teacher
	assignmentIDs, err := r.findAssignmentIDsByTeacher(ctx, teacherID)
	if err != nil || len(assignmentIDs) == 0 {
		return 0, err
	}
	filter := bson.M{
		"assignmentId": bson.M{"$in": assignmentIDs},
		"score":        nil,
	}
	count, err := r.assignmentSubmissions.CountDocuments(ctx, filter)
	return int(count), err
}

func (r *dashboardRepoImpl) CountUngradedProjectSubmissions(ctx context.Context, teacherID bson.ObjectID) (int, error) {
	projectIDs, err := r.findProjectIDsByTeacher(ctx, teacherID)
	if err != nil || len(projectIDs) == 0 {
		return 0, err
	}
	filter := bson.M{
		"projectId": bson.M{"$in": projectIDs},
		"score":     nil,
	}
	count, err := r.projectSubmissions.CountDocuments(ctx, filter)
	return int(count), err
}

func (r *dashboardRepoImpl) CountUngradedAssessmentAttempts(ctx context.Context, teacherID bson.ObjectID) (int, error) {
	assessmentIDs, err := r.findAssessmentIDsByTeacher(ctx, teacherID)
	if err != nil || len(assessmentIDs) == 0 {
		return 0, err
	}
	filter := bson.M{
		"assessmentId": bson.M{"$in": assessmentIDs},
		"status":       string(model.AttemptSubmitted),
	}
	count, err := r.assessmentAttempts.CountDocuments(ctx, filter)
	return int(count), err
}

func (r *dashboardRepoImpl) FindTodaySchedules(ctx context.Context, teacherID, academicYearID bson.ObjectID, dayOfWeek int) ([]*model.Schedule, error) {
	filter := bson.M{
		"teacherId":      teacherID,
		"academicYearId": academicYearID,
		"day":            dayOfWeek,
	}
	opts := options.Find().SetSort(bson.D{{Key: "startTime", Value: 1}})
	cursor, err := r.schedules.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var schedules []*model.Schedule
	if err := cursor.All(ctx, &schedules); err != nil {
		return nil, err
	}
	return schedules, nil
}

func (r *dashboardRepoImpl) FindRecentAssignmentSubmissions(ctx context.Context, teacherID bson.ObjectID, limit int) ([]*model.AssignmentSubmission, error) {
	assignmentIDs, err := r.findAssignmentIDsByTeacher(ctx, teacherID)
	if err != nil || len(assignmentIDs) == 0 {
		return nil, err
	}
	filter := bson.M{
		"assignmentId": bson.M{"$in": assignmentIDs},
	}
	opts := options.Find().
		SetSort(bson.D{{Key: "submittedAt", Value: -1}}).
		SetLimit(int64(limit))
	cursor, err := r.assignmentSubmissions.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var subs []*model.AssignmentSubmission
	if err := cursor.All(ctx, &subs); err != nil {
		return nil, err
	}
	return subs, nil
}

func (r *dashboardRepoImpl) FindRecentProjectSubmissions(ctx context.Context, teacherID bson.ObjectID, limit int) ([]*model.ProjectSubmission, error) {
	projectIDs, err := r.findProjectIDsByTeacher(ctx, teacherID)
	if err != nil || len(projectIDs) == 0 {
		return nil, err
	}
	filter := bson.M{
		"projectId": bson.M{"$in": projectIDs},
	}
	opts := options.Find().
		SetSort(bson.D{{Key: "submittedAt", Value: -1}}).
		SetLimit(int64(limit))
	cursor, err := r.projectSubmissions.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var subs []*model.ProjectSubmission
	if err := cursor.All(ctx, &subs); err != nil {
		return nil, err
	}
	return subs, nil
}

func (r *dashboardRepoImpl) FindAssignmentsByIDs(ctx context.Context, ids []bson.ObjectID) (map[bson.ObjectID]*model.Assignment, error) {
	result := make(map[bson.ObjectID]*model.Assignment)
	if len(ids) == 0 {
		return result, nil
	}
	cursor, err := r.assignments.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	var items []*model.Assignment
	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	for _, a := range items {
		result[a.ID] = a
	}
	return result, nil
}

func (r *dashboardRepoImpl) FindProjectsByIDs(ctx context.Context, ids []bson.ObjectID) (map[bson.ObjectID]*model.Project, error) {
	result := make(map[bson.ObjectID]*model.Project)
	if len(ids) == 0 {
		return result, nil
	}
	cursor, err := r.projects.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	var items []*model.Project
	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	for _, p := range items {
		result[p.ID] = p
	}
	return result, nil
}

func (r *dashboardRepoImpl) FindStudentsByIDs(ctx context.Context, ids []bson.ObjectID) (map[bson.ObjectID]*model.Student, error) {
	result := make(map[bson.ObjectID]*model.Student)
	if len(ids) == 0 {
		return result, nil
	}
	cursor, err := r.students.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	var items []*model.Student
	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	for _, s := range items {
		result[s.ID] = s
	}
	return result, nil
}

func (r *dashboardRepoImpl) FindClassByID(ctx context.Context, classID bson.ObjectID) (*model.Class, error) {
	var cl model.Class
	err := r.classes.FindOne(ctx, bson.M{"_id": classID}).Decode(&cl)
	if err != nil {
		return nil, err
	}
	return &cl, nil
}

func (r *dashboardRepoImpl) FindSubjectsByIDs(ctx context.Context, ids []bson.ObjectID) (map[bson.ObjectID]*model.Subject, error) {
	result := make(map[bson.ObjectID]*model.Subject)
	if len(ids) == 0 {
		return result, nil
	}
	cursor, err := r.subjects.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	var items []*model.Subject
	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	for _, s := range items {
		result[s.ID] = s
	}
	return result, nil
}

func (r *dashboardRepoImpl) FindClassesByIDs(ctx context.Context, ids []bson.ObjectID) (map[bson.ObjectID]*model.Class, error) {
	result := make(map[bson.ObjectID]*model.Class)
	if len(ids) == 0 {
		return result, nil
	}
	cursor, err := r.classes.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	var items []*model.Class
	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	for _, c := range items {
		result[c.ID] = c
	}
	return result, nil
}

func (r *dashboardRepoImpl) CountStudentsByClassID(ctx context.Context, schoolID, classID bson.ObjectID) (int, error) {
	filter := bson.M{
		"classId":  classID,
		"schoolId": schoolID,
		"status":   "active",
	}
	count, err := r.students.CountDocuments(ctx, filter)
	return int(count), err
}

// --- internal helpers ---

func (r *dashboardRepoImpl) findAssignmentIDsByTeacher(ctx context.Context, teacherID bson.ObjectID) ([]bson.ObjectID, error) {
	opts := options.Find().SetProjection(bson.M{"_id": 1})
	cursor, err := r.assignments.Find(ctx, bson.M{"teacherId": teacherID}, opts)
	if err != nil {
		return nil, err
	}
	var results []struct {
		ID bson.ObjectID `bson:"_id"`
	}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	ids := make([]bson.ObjectID, len(results))
	for i, r := range results {
		ids[i] = r.ID
	}
	return ids, nil
}

func (r *dashboardRepoImpl) findProjectIDsByTeacher(ctx context.Context, teacherID bson.ObjectID) ([]bson.ObjectID, error) {
	opts := options.Find().SetProjection(bson.M{"_id": 1})
	cursor, err := r.projects.Find(ctx, bson.M{"teacherId": teacherID}, opts)
	if err != nil {
		return nil, err
	}
	var results []struct {
		ID bson.ObjectID `bson:"_id"`
	}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	ids := make([]bson.ObjectID, len(results))
	for i, r := range results {
		ids[i] = r.ID
	}
	return ids, nil
}

func (r *dashboardRepoImpl) findAssessmentIDsByTeacher(ctx context.Context, teacherID bson.ObjectID) ([]bson.ObjectID, error) {
	opts := options.Find().SetProjection(bson.M{"_id": 1})
	cursor, err := r.assessments.Find(ctx, bson.M{"teacherId": teacherID}, opts)
	if err != nil {
		return nil, err
	}
	var results []struct {
		ID bson.ObjectID `bson:"_id"`
	}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	ids := make([]bson.ObjectID, len(results))
	for i, r := range results {
		ids[i] = r.ID
	}
	return ids, nil
}
