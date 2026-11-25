package studentsvc

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"sort"
	"strings"
	"time"

	studentdto "github.com/geniussociety/backend/internal/dto/student"
	"github.com/geniussociety/backend/internal/handler"
	"github.com/geniussociety/backend/internal/middleware"
	"github.com/geniussociety/backend/internal/model"
	studentrepo "github.com/geniussociety/backend/internal/repository/student"
	"github.com/geniussociety/backend/internal/storage"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Allowed MIME types per instruction
var allowedMIMETypes = map[string]bool{
	"application/pdf":                                                           true,
	"application/zip":                                                           true,
	"application/x-zip-compressed":                                              true,
	"application/msword":                                                        true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document":   true,
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         true,
	"application/vnd.openxmlformats-officedocument.presentationml.presentation": true,
}

// StudentService defines all business operations for student portal.
type StudentService interface {
	GetDashboard(ctx context.Context, scope *middleware.StudentScope) (*studentdto.DashboardResponseDTO, error)

	// Materials
	ListMaterials(ctx context.Context, scope *middleware.StudentScope, subjectIDStr string) ([]studentdto.MaterialListDTO, error)
	GetMaterialDetail(ctx context.Context, scope *middleware.StudentScope, materialID bson.ObjectID) (*studentdto.MaterialDetailDTO, error)
	GetMaterialDownloadURL(ctx context.Context, scope *middleware.StudentScope, materialID bson.ObjectID, fileIndex int) (string, error)

	// Assignments
	ListAssignments(ctx context.Context, scope *middleware.StudentScope, subjectIDStr string) ([]studentdto.AssignmentListDTO, error)
	GetAssignmentDetail(ctx context.Context, scope *middleware.StudentScope, assignmentID bson.ObjectID) (*studentdto.AssignmentDetailDTO, error)
	SubmitAssignment(ctx context.Context, scope *middleware.StudentScope, assignmentID bson.ObjectID, file *multipart.FileHeader, note string, isResubmit bool) (*studentdto.SubmissionInfoDTO, error)
	GetAssignmentFileDownloadURL(ctx context.Context, scope *middleware.StudentScope, assignmentID bson.ObjectID, fileIndex int) (string, error)

	// Projects
	ListProjects(ctx context.Context, scope *middleware.StudentScope, subjectIDStr string) ([]studentdto.ProjectListDTO, error)
	GetProjectDetail(ctx context.Context, scope *middleware.StudentScope, projectID bson.ObjectID) (*studentdto.ProjectDetailDTO, error)
	SubmitProject(ctx context.Context, scope *middleware.StudentScope, projectID bson.ObjectID, file *multipart.FileHeader, isResubmit bool) (*studentdto.SubmissionInfoDTO, error)
	GetProjectFileDownloadURL(ctx context.Context, scope *middleware.StudentScope, projectID bson.ObjectID, fileIndex int) (string, error)

	// Assessments
	ListAssessments(ctx context.Context, scope *middleware.StudentScope, subjectIDStr string) ([]studentdto.AssessmentListDTO, error)
	GetAssessmentDetail(ctx context.Context, scope *middleware.StudentScope, assessmentID bson.ObjectID) (*studentdto.AssessmentDetailDTO, error)
	StartAssessment(ctx context.Context, scope *middleware.StudentScope, assessmentID bson.ObjectID) (*studentdto.AssessmentStartResponseDTO, error)
	SubmitAssessment(ctx context.Context, scope *middleware.StudentScope, assessmentID bson.ObjectID, payload *studentdto.AssessmentSubmitPayload) (*studentdto.AssessmentResultResponseDTO, error)
	GetAssessmentResult(ctx context.Context, scope *middleware.StudentScope, assessmentID bson.ObjectID) (*studentdto.AssessmentResultResponseDTO, error)

	// Grades
	GetGrades(ctx context.Context, scope *middleware.StudentScope) ([]studentdto.SubjectGradeDTO, error)
	GetSubjectGradeDetail(ctx context.Context, scope *middleware.StudentScope, subjectID bson.ObjectID) (*studentdto.SubjectGradeDetailDTO, error)

	// Schedule & Ranking
	GetSchedule(ctx context.Context, scope *middleware.StudentScope) ([]studentdto.ScheduleItemDTO, error)
	GetRanking(ctx context.Context, scope *middleware.StudentScope) (*studentdto.ClassRankingResponseDTO, error)
}

type studentServiceImpl struct {
	repo    studentrepo.StudentRepository
	storage storage.FileStorage
}

// NewStudentService constructs a StudentService implementation.
func NewStudentService(repo studentrepo.StudentRepository, st storage.FileStorage) StudentService {
	return &studentServiceImpl{
		repo:    repo,
		storage: st,
	}
}

// -------------------------------------------------------------
// 1. Dashboard
// -------------------------------------------------------------

func (s *studentServiceImpl) GetDashboard(ctx context.Context, scope *middleware.StudentScope) (*studentdto.DashboardResponseDTO, error) {
	now := time.Now().UTC()

	// 1. Assignments
	assignments, err := s.repo.FindAssignmentsByClass(ctx, scope.SchoolID, scope.ClassID, nil)
	if err != nil {
		return nil, err
	}

	activeAssignments := 0
	var deadlines []studentdto.DeadlineItem
	for _, a := range assignments {
		latestSub, _ := s.repo.FindLatestAssignmentSubmission(ctx, scope.SchoolID, a.ID, scope.StudentID)
		submitted := latestSub != nil

		if !submitted && a.Deadline.After(now) {
			activeAssignments++
		}

		if a.Deadline.After(now.Add(-24 * time.Hour)) {
			deadlines = append(deadlines, studentdto.DeadlineItem{
				ID:        a.ID.Hex(),
				Type:      "assignment",
				Title:     a.Title,
				Deadline:  a.Deadline,
				Submitted: submitted,
			})
		}
	}

	// 2. Projects
	projects, _ := s.repo.FindProjectsByClass(ctx, scope.SchoolID, scope.ClassID, nil)
	for _, p := range projects {
		latestSub, _ := s.repo.FindLatestProjectSubmission(ctx, scope.SchoolID, p.ID, scope.StudentID)
		submitted := latestSub != nil

		if p.Deadline.After(now.Add(-24 * time.Hour)) {
			deadlines = append(deadlines, studentdto.DeadlineItem{
				ID:        p.ID.Hex(),
				Type:      "project",
				Title:     p.Title,
				Deadline:  p.Deadline,
				Submitted: submitted,
			})
		}
	}

	// Sort deadlines ascending by deadline date, take at most 3
	sort.Slice(deadlines, func(i, j int) bool {
		return deadlines[i].Deadline.Before(deadlines[j].Deadline)
	})
	if len(deadlines) > 3 {
		deadlines = deadlines[:3]
	}

	// 3. Active Assessments (status published, endAt > now or nil, no finished attempt)
	assessments, _ := s.repo.FindAssessmentsByClass(ctx, scope.SchoolID, scope.ClassID, nil)
	activeAssessments := 0
	for _, asm := range assessments {
		if asm.Status != model.AssessmentPublished {
			continue
		}
		if asm.EndAt != nil && asm.EndAt.Before(now) {
			continue
		}
		attempt, _ := s.repo.FindAttemptByStudent(ctx, scope.SchoolID, asm.ID, scope.StudentID)
		if attempt == nil || attempt.Status == model.AttemptInProgress {
			activeAssessments++
		}
	}

	// 4. Recent Grades (last 5 graded entries)
	grades, _ := s.repo.FindGradesByStudent(ctx, scope.SchoolID, scope.StudentID, scope.AcademicYearID)
	sort.Slice(grades, func(i, j int) bool {
		return grades[i].CreatedAt.After(grades[j].CreatedAt)
	})

	var recentGrades []studentdto.RecentGradeItem
	var subOIDs []bson.ObjectID
	for _, g := range grades {
		if !g.SubjectID.IsZero() {
			subOIDs = append(subOIDs, g.SubjectID)
		}
	}
	subMap, _ := s.repo.FindSubjectsMap(ctx, subOIDs)

	for i, g := range grades {
		if i >= 5 {
			break
		}
		subName := "Mata Pelajaran"
		if sm, ok := subMap[g.SubjectID]; ok && sm != nil {
			subName = sm.Name
		}
		title := fmt.Sprintf("Nilai %s", strings.ToUpper(string(g.Type)))
		recentGrades = append(recentGrades, studentdto.RecentGradeItem{
			ID:       g.ID.Hex(),
			Type:     string(g.Type),
			Title:    title,
			Subject:  subName,
			Score:    g.Score,
			GradedAt: g.CreatedAt,
		})
	}
	if recentGrades == nil {
		recentGrades = []studentdto.RecentGradeItem{}
	}

	// 5. Today Schedule (based on current weekday)
	// Go weekday: Sunday=0, Monday=1, ... Saturday=6.
	// Model Schedule: 1=Monday ... 7=Sunday.
	weekday := int(now.Weekday())
	schedDay := weekday
	if weekday == 0 {
		schedDay = 7
	}

	allSchedules, _ := s.repo.FindSchedulesByClass(ctx, scope.SchoolID, scope.ClassID, scope.AcademicYearID)
	var todaySchedules []studentdto.ScheduleItemDTO
	var schedTeacherIDs []bson.ObjectID
	var schedSubIDs []bson.ObjectID
	for _, sc := range allSchedules {
		if sc.Day == schedDay {
			schedTeacherIDs = append(schedTeacherIDs, sc.TeacherID)
			schedSubIDs = append(schedSubIDs, sc.SubjectID)
		}
	}
	tNames, _ := s.repo.FindTeachersMap(ctx, schedTeacherIDs)
	sNames, _ := s.repo.FindSubjectsMap(ctx, schedSubIDs)

	for _, sc := range allSchedules {
		if sc.Day == schedDay {
			sName := "Mata Pelajaran"
			if sm, ok := sNames[sc.SubjectID]; ok && sm != nil {
				sName = sm.Name
			}
			tName := tNames[sc.TeacherID]
			todaySchedules = append(todaySchedules, studentdto.ScheduleItemDTO{
				ID:        sc.ID.Hex(),
				Day:       sc.Day,
				DayName:   getDayName(sc.Day),
				StartTime: sc.StartTime,
				EndTime:   sc.EndTime,
				Subject:   sName,
				Teacher:   tName,
				Room:      sc.Room,
			})
		}
	}
	if todaySchedules == nil {
		todaySchedules = []studentdto.ScheduleItemDTO{}
	}

	return &studentdto.DashboardResponseDTO{
		ActiveAssignmentsCount: activeAssignments,
		ActiveAssessmentsCount: activeAssessments,
		UpcomingDeadlines:      deadlines,
		RecentGrades:           recentGrades,
		TodaySchedule:          todaySchedules,
	}, nil
}

// -------------------------------------------------------------
// 2. Materials
// -------------------------------------------------------------

func (s *studentServiceImpl) ListMaterials(ctx context.Context, scope *middleware.StudentScope, subjectIDStr string) ([]studentdto.MaterialListDTO, error) {
	now := time.Now().UTC()
	var subOID *bson.ObjectID
	if subjectIDStr != "" {
		id, err := bson.ObjectIDFromHex(subjectIDStr)
		if err == nil {
			subOID = &id
		}
	}

	materials, err := s.repo.FindMaterialsByClass(ctx, scope.SchoolID, scope.ClassID, subOID, now)
	if err != nil {
		return nil, err
	}

	var subIDs, teacherIDs []bson.ObjectID
	for _, m := range materials {
		if !m.SubjectID.IsZero() {
			subIDs = append(subIDs, m.SubjectID)
		}
		if !m.TeacherID.IsZero() {
			teacherIDs = append(teacherIDs, m.TeacherID)
		}
	}
	subMap, _ := s.repo.FindSubjectsMap(ctx, subIDs)
	teacherMap, _ := s.repo.FindTeachersMap(ctx, teacherIDs)

	list := make([]studentdto.MaterialListDTO, 0, len(materials))
	for _, m := range materials {
		sName := "Mata Pelajaran"
		if sm, ok := subMap[m.SubjectID]; ok && sm != nil {
			sName = sm.Name
		}
		tName := teacherMap[m.TeacherID]

		fileCount := len(m.Attachments)
		if fileCount == 0 && m.File != nil {
			fileCount = 1
		}

		list = append(list, studentdto.MaterialListDTO{
			ID:          m.ID.Hex(),
			Title:       m.Title,
			Description: m.Description,
			SubjectID:   m.SubjectID.Hex(),
			SubjectName: sName,
			TeacherName: tName,
			FileCount:   fileCount,
			Link:        m.Link,
			PublishAt:   m.PublishAt,
			CreatedAt:   m.CreatedAt,
		})
	}
	return list, nil
}

func (s *studentServiceImpl) GetMaterialDetail(ctx context.Context, scope *middleware.StudentScope, materialID bson.ObjectID) (*studentdto.MaterialDetailDTO, error) {
	m, err := s.repo.FindMaterialByID(ctx, scope.SchoolID, materialID)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, handler.ErrNotFound
	}

	// Verify class membership
	if !containsObjectID(m.ClassIDs, scope.ClassID) {
		return nil, handler.NewForbiddenError("Materi ini tidak ditujukan untuk kelas Anda")
	}

	subMap, _ := s.repo.FindSubjectsMap(ctx, []bson.ObjectID{m.SubjectID})
	teacherMap, _ := s.repo.FindTeachersMap(ctx, []bson.ObjectID{m.TeacherID})

	sName := "Mata Pelajaran"
	if sm, ok := subMap[m.SubjectID]; ok && sm != nil {
		sName = sm.Name
	}
	tName := teacherMap[m.TeacherID]

	attachments := m.Attachments
	if len(attachments) == 0 && m.File != nil {
		attachments = []model.FileRef{*m.File}
	}

	return &studentdto.MaterialDetailDTO{
		ID:          m.ID.Hex(),
		Title:       m.Title,
		Description: m.Description,
		SubjectID:   m.SubjectID.Hex(),
		SubjectName: sName,
		TeacherName: tName,
		Attachments: attachments,
		Link:        m.Link,
		PublishAt:   m.PublishAt,
		CreatedAt:   m.CreatedAt,
	}, nil
}

func (s *studentServiceImpl) GetMaterialDownloadURL(ctx context.Context, scope *middleware.StudentScope, materialID bson.ObjectID, fileIndex int) (string, error) {
	m, err := s.repo.FindMaterialByID(ctx, scope.SchoolID, materialID)
	if err != nil {
		return "", err
	}
	if m == nil {
		return "", handler.ErrNotFound
	}
	if !containsObjectID(m.ClassIDs, scope.ClassID) {
		return "", handler.NewForbiddenError("Materi ini tidak ditujukan untuk kelas Anda")
	}

	attachments := m.Attachments
	if len(attachments) == 0 && m.File != nil {
		attachments = []model.FileRef{*m.File}
	}

	if fileIndex < 0 || fileIndex >= len(attachments) {
		return "", handler.NewValidationError("Index berkas lampiran tidak valid")
	}

	key := attachments[fileIndex].Key
	return s.storage.PresignGet(ctx, key, 5*time.Minute)
}

// -------------------------------------------------------------
// 3. Assignments
// -------------------------------------------------------------

func (s *studentServiceImpl) ListAssignments(ctx context.Context, scope *middleware.StudentScope, subjectIDStr string) ([]studentdto.AssignmentListDTO, error) {
	var subOID *bson.ObjectID
	if subjectIDStr != "" {
		id, err := bson.ObjectIDFromHex(subjectIDStr)
		if err == nil {
			subOID = &id
		}
	}

	assignments, err := s.repo.FindAssignmentsByClass(ctx, scope.SchoolID, scope.ClassID, subOID)
	if err != nil {
		return nil, err
	}

	var subIDs, teacherIDs []bson.ObjectID
	for _, a := range assignments {
		if !a.SubjectID.IsZero() {
			subIDs = append(subIDs, a.SubjectID)
		}
		if !a.TeacherID.IsZero() {
			teacherIDs = append(teacherIDs, a.TeacherID)
		}
	}
	subMap, _ := s.repo.FindSubjectsMap(ctx, subIDs)
	teacherMap, _ := s.repo.FindTeachersMap(ctx, teacherIDs)

	now := time.Now().UTC()
	list := make([]studentdto.AssignmentListDTO, 0, len(assignments))
	for _, a := range assignments {
		sName := "Mata Pelajaran"
		if sm, ok := subMap[a.SubjectID]; ok && sm != nil {
			sName = sm.Name
		}
		tName := teacherMap[a.TeacherID]

		latestSub, _ := s.repo.FindLatestAssignmentSubmission(ctx, scope.SchoolID, a.ID, scope.StudentID)
		status := "pending"
		var subInfo *studentdto.SubmissionInfoDTO

		if latestSub != nil {
			subInfo = &studentdto.SubmissionInfoDTO{
				ID:          latestSub.ID.Hex(),
				Submitted:   true,
				SubmittedAt: &latestSub.SubmittedAt,
				Files:       latestSub.Files,
				Note:        latestSub.Note,
				Version:     latestSub.Version,
				Score:       latestSub.Score,
				Feedback:    latestSub.Feedback,
				GradedAt:    latestSub.GradedAt,
				IsLate:      latestSub.SubmittedAt.After(a.Deadline),
			}
			if latestSub.Score != nil {
				status = "graded"
			} else {
				status = "submitted"
			}
		} else if now.After(a.Deadline) {
			status = "overdue"
		}

		list = append(list, studentdto.AssignmentListDTO{
			ID:            a.ID.Hex(),
			Title:         a.Title,
			SubjectID:     a.SubjectID.Hex(),
			SubjectName:   sName,
			TeacherName:   tName,
			Deadline:      a.Deadline,
			AllowResubmit: a.AllowResubmit,
			Status:        status,
			Submission:    subInfo,
			CreatedAt:     a.CreatedAt,
		})
	}
	return list, nil
}

func (s *studentServiceImpl) GetAssignmentDetail(ctx context.Context, scope *middleware.StudentScope, assignmentID bson.ObjectID) (*studentdto.AssignmentDetailDTO, error) {
	a, err := s.repo.FindAssignmentByID(ctx, scope.SchoolID, assignmentID)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, handler.ErrNotFound
	}
	if !containsObjectID(a.ClassIDs, scope.ClassID) {
		return nil, handler.NewForbiddenError("Tugas ini tidak ditujukan untuk kelas Anda")
	}

	subMap, _ := s.repo.FindSubjectsMap(ctx, []bson.ObjectID{a.SubjectID})
	teacherMap, _ := s.repo.FindTeachersMap(ctx, []bson.ObjectID{a.TeacherID})

	sName := "Mata Pelajaran"
	if sm, ok := subMap[a.SubjectID]; ok && sm != nil {
		sName = sm.Name
	}
	tName := teacherMap[a.TeacherID]

	latestSub, _ := s.repo.FindLatestAssignmentSubmission(ctx, scope.SchoolID, a.ID, scope.StudentID)
	status := "pending"
	var subInfo *studentdto.SubmissionInfoDTO

	now := time.Now().UTC()
	if latestSub != nil {
		subInfo = &studentdto.SubmissionInfoDTO{
			ID:          latestSub.ID.Hex(),
			Submitted:   true,
			SubmittedAt: &latestSub.SubmittedAt,
			Files:       latestSub.Files,
			Note:        latestSub.Note,
			Version:     latestSub.Version,
			Score:       latestSub.Score,
			Feedback:    latestSub.Feedback,
			GradedAt:    latestSub.GradedAt,
			IsLate:      latestSub.SubmittedAt.After(a.Deadline),
		}
		if latestSub.Score != nil {
			status = "graded"
		} else {
			status = "submitted"
		}
	} else if now.After(a.Deadline) {
		status = "overdue"
	}

	return &studentdto.AssignmentDetailDTO{
		ID:            a.ID.Hex(),
		Title:         a.Title,
		Description:   a.Description,
		SubjectID:     a.SubjectID.Hex(),
		SubjectName:   sName,
		TeacherName:   tName,
		Deadline:      a.Deadline,
		Attachments:   a.Attachments,
		AllowResubmit: a.AllowResubmit,
		Status:        status,
		Submission:    subInfo,
		CreatedAt:     a.CreatedAt,
	}, nil
}

func (s *studentServiceImpl) SubmitAssignment(ctx context.Context, scope *middleware.StudentScope, assignmentID bson.ObjectID, file *multipart.FileHeader, note string, isResubmit bool) (*studentdto.SubmissionInfoDTO, error) {
	a, err := s.repo.FindAssignmentByID(ctx, scope.SchoolID, assignmentID)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, handler.ErrNotFound
	}
	if !containsObjectID(a.ClassIDs, scope.ClassID) {
		return nil, handler.NewForbiddenError("Tugas ini tidak ditujukan untuk kelas Anda")
	}

	latestSub, err := s.repo.FindLatestAssignmentSubmission(ctx, scope.SchoolID, assignmentID, scope.StudentID)
	if err != nil {
		return nil, err
	}

	version := 1
	if latestSub != nil {
		if !a.AllowResubmit {
			return nil, handler.NewForbiddenError("Tugas ini tidak mengizinkan pengumpulan ulang (resubmit)")
		}
		if latestSub.Score != nil {
			return nil, handler.NewForbiddenError("Tugas telah dinilai oleh guru dan tidak dapat dikumpulkan ulang")
		}
		version = latestSub.Version + 1
	}

	var files []model.FileRef
	if file != nil {
		ref, err := s.uploadStudentFile(ctx, "assignment", scope.StudentID.Hex(), assignmentID.Hex(), file)
		if err != nil {
			return nil, err
		}
		files = append(files, *ref)
	}

	now := time.Now().UTC()
	sub := &model.AssignmentSubmission{
		ID:           bson.NewObjectID(),
		SchoolID:     scope.SchoolID,
		AssignmentID: assignmentID,
		StudentID:    scope.StudentID,
		Files:        files,
		Note:         note,
		SubmittedAt:  now,
		Version:      version,
	}

	if err := s.repo.InsertAssignmentSubmission(ctx, sub); err != nil {
		return nil, err
	}

	return &studentdto.SubmissionInfoDTO{
		ID:          sub.ID.Hex(),
		Submitted:   true,
		SubmittedAt: &now,
		Files:       files,
		Note:        note,
		Version:     version,
		IsLate:      now.After(a.Deadline),
	}, nil
}

func (s *studentServiceImpl) GetAssignmentFileDownloadURL(ctx context.Context, scope *middleware.StudentScope, assignmentID bson.ObjectID, fileIndex int) (string, error) {
	a, err := s.repo.FindAssignmentByID(ctx, scope.SchoolID, assignmentID)
	if err != nil {
		return "", err
	}
	if a == nil {
		return "", handler.ErrNotFound
	}
	if !containsObjectID(a.ClassIDs, scope.ClassID) {
		return "", handler.NewForbiddenError("Tugas ini tidak ditujukan untuk kelas Anda")
	}

	if fileIndex < 0 || fileIndex >= len(a.Attachments) {
		return "", handler.NewValidationError("Index berkas soal tidak valid")
	}

	key := a.Attachments[fileIndex].Key
	return s.storage.PresignGet(ctx, key, 5*time.Minute)
}

// -------------------------------------------------------------
// 4. Projects
// -------------------------------------------------------------

func (s *studentServiceImpl) ListProjects(ctx context.Context, scope *middleware.StudentScope, subjectIDStr string) ([]studentdto.ProjectListDTO, error) {
	var subOID *bson.ObjectID
	if subjectIDStr != "" {
		id, err := bson.ObjectIDFromHex(subjectIDStr)
		if err == nil {
			subOID = &id
		}
	}

	projects, err := s.repo.FindProjectsByClass(ctx, scope.SchoolID, scope.ClassID, subOID)
	if err != nil {
		return nil, err
	}

	var subIDs, teacherIDs []bson.ObjectID
	for _, p := range projects {
		if !p.SubjectID.IsZero() {
			subIDs = append(subIDs, p.SubjectID)
		}
		if !p.TeacherID.IsZero() {
			teacherIDs = append(teacherIDs, p.TeacherID)
		}
	}
	subMap, _ := s.repo.FindSubjectsMap(ctx, subIDs)
	teacherMap, _ := s.repo.FindTeachersMap(ctx, teacherIDs)

	now := time.Now().UTC()
	list := make([]studentdto.ProjectListDTO, 0, len(projects))
	for _, p := range projects {
		sName := "Mata Pelajaran"
		if sm, ok := subMap[p.SubjectID]; ok && sm != nil {
			sName = sm.Name
		}
		tName := teacherMap[p.TeacherID]

		latestSub, _ := s.repo.FindLatestProjectSubmission(ctx, scope.SchoolID, p.ID, scope.StudentID)
		status := "pending"
		var subInfo *studentdto.SubmissionInfoDTO

		if latestSub != nil {
			subInfo = &studentdto.SubmissionInfoDTO{
				ID:          latestSub.ID.Hex(),
				Submitted:   true,
				SubmittedAt: &latestSub.SubmittedAt,
				Files:       latestSub.Files,
				Version:     latestSub.Version,
				Score:       latestSub.Score,
				Feedback:    latestSub.Feedback,
				GradedAt:    latestSub.GradedAt,
				IsLate:      latestSub.SubmittedAt.After(p.Deadline),
			}
			if latestSub.Score != nil {
				status = "graded"
			} else {
				status = "submitted"
			}
		} else if now.After(p.Deadline) {
			status = "overdue"
		}

		list = append(list, studentdto.ProjectListDTO{
			ID:            p.ID.Hex(),
			Title:         p.Title,
			SubjectID:     p.SubjectID.Hex(),
			SubjectName:   sName,
			TeacherName:   tName,
			Deadline:      p.Deadline,
			AllowResubmit: p.AllowResubmit,
			Status:        status,
			Submission:    subInfo,
			CreatedAt:     p.CreatedAt,
		})
	}
	return list, nil
}

func (s *studentServiceImpl) GetProjectDetail(ctx context.Context, scope *middleware.StudentScope, projectID bson.ObjectID) (*studentdto.ProjectDetailDTO, error) {
	p, err := s.repo.FindProjectByID(ctx, scope.SchoolID, projectID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, handler.ErrNotFound
	}
	if !containsObjectID(p.ClassIDs, scope.ClassID) {
		return nil, handler.NewForbiddenError("Projek ini tidak ditujukan untuk kelas Anda")
	}

	subMap, _ := s.repo.FindSubjectsMap(ctx, []bson.ObjectID{p.SubjectID})
	teacherMap, _ := s.repo.FindTeachersMap(ctx, []bson.ObjectID{p.TeacherID})

	sName := "Mata Pelajaran"
	if sm, ok := subMap[p.SubjectID]; ok && sm != nil {
		sName = sm.Name
	}
	tName := teacherMap[p.TeacherID]

	latestSub, _ := s.repo.FindLatestProjectSubmission(ctx, scope.SchoolID, p.ID, scope.StudentID)
	status := "pending"
	var subInfo *studentdto.SubmissionInfoDTO

	now := time.Now().UTC()
	if latestSub != nil {
		subInfo = &studentdto.SubmissionInfoDTO{
			ID:          latestSub.ID.Hex(),
			Submitted:   true,
			SubmittedAt: &latestSub.SubmittedAt,
			Files:       latestSub.Files,
			Version:     latestSub.Version,
			Score:       latestSub.Score,
			Feedback:    latestSub.Feedback,
			GradedAt:    latestSub.GradedAt,
			IsLate:      latestSub.SubmittedAt.After(p.Deadline),
		}
		if latestSub.Score != nil {
			status = "graded"
		} else {
			status = "submitted"
		}
	} else if now.After(p.Deadline) {
		status = "overdue"
	}

	return &studentdto.ProjectDetailDTO{
		ID:            p.ID.Hex(),
		Title:         p.Title,
		Description:   p.Description,
		Instructions:  p.Instructions,
		SubjectID:     p.SubjectID.Hex(),
		SubjectName:   sName,
		TeacherName:   tName,
		Deadline:      p.Deadline,
		Attachments:   p.Attachments,
		AllowResubmit: p.AllowResubmit,
		Status:        status,
		Submission:    subInfo,
		CreatedAt:     p.CreatedAt,
	}, nil
}

func (s *studentServiceImpl) SubmitProject(ctx context.Context, scope *middleware.StudentScope, projectID bson.ObjectID, file *multipart.FileHeader, isResubmit bool) (*studentdto.SubmissionInfoDTO, error) {
	p, err := s.repo.FindProjectByID(ctx, scope.SchoolID, projectID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, handler.ErrNotFound
	}
	if !containsObjectID(p.ClassIDs, scope.ClassID) {
		return nil, handler.NewForbiddenError("Projek ini tidak ditujukan untuk kelas Anda")
	}

	latestSub, err := s.repo.FindLatestProjectSubmission(ctx, scope.SchoolID, projectID, scope.StudentID)
	if err != nil {
		return nil, err
	}

	version := 1
	if latestSub != nil {
		if !p.AllowResubmit {
			return nil, handler.NewForbiddenError("Projek ini tidak mengizinkan pengumpulan ulang (resubmit)")
		}
		if latestSub.Score != nil {
			return nil, handler.NewForbiddenError("Projek telah dinilai oleh guru dan tidak dapat dikumpulkan ulang")
		}
		version = latestSub.Version + 1
	}

	var files []model.FileRef
	if file != nil {
		ref, err := s.uploadStudentFile(ctx, "project", scope.StudentID.Hex(), projectID.Hex(), file)
		if err != nil {
			return nil, err
		}
		files = append(files, *ref)
	}

	now := time.Now().UTC()
	sub := &model.ProjectSubmission{
		ID:          bson.NewObjectID(),
		SchoolID:    scope.SchoolID,
		ProjectID:   projectID,
		StudentID:   scope.StudentID,
		Files:       files,
		SubmittedAt: now,
		Version:     version,
	}

	if err := s.repo.InsertProjectSubmission(ctx, sub); err != nil {
		return nil, err
	}

	return &studentdto.SubmissionInfoDTO{
		ID:          sub.ID.Hex(),
		Submitted:   true,
		SubmittedAt: &now,
		Files:       files,
		Version:     version,
		IsLate:      now.After(p.Deadline),
	}, nil
}

func (s *studentServiceImpl) GetProjectFileDownloadURL(ctx context.Context, scope *middleware.StudentScope, projectID bson.ObjectID, fileIndex int) (string, error) {
	p, err := s.repo.FindProjectByID(ctx, scope.SchoolID, projectID)
	if err != nil {
		return "", err
	}
	if p == nil {
		return "", handler.ErrNotFound
	}
	if !containsObjectID(p.ClassIDs, scope.ClassID) {
		return "", handler.NewForbiddenError("Projek ini tidak ditujukan untuk kelas Anda")
	}

	if fileIndex < 0 || fileIndex >= len(p.Attachments) {
		return "", handler.NewValidationError("Index berkas instruksi projek tidak valid")
	}

	key := p.Attachments[fileIndex].Key
	return s.storage.PresignGet(ctx, key, 5*time.Minute)
}

// -------------------------------------------------------------
// 5. Assessments
// -------------------------------------------------------------

func (s *studentServiceImpl) ListAssessments(ctx context.Context, scope *middleware.StudentScope, subjectIDStr string) ([]studentdto.AssessmentListDTO, error) {
	var subOID *bson.ObjectID
	if subjectIDStr != "" {
		id, err := bson.ObjectIDFromHex(subjectIDStr)
		if err == nil {
			subOID = &id
		}
	}

	assessments, err := s.repo.FindAssessmentsByClass(ctx, scope.SchoolID, scope.ClassID, subOID)
	if err != nil {
		return nil, err
	}

	var subIDs, teacherIDs []bson.ObjectID
	for _, a := range assessments {
		if !a.SubjectID.IsZero() {
			subIDs = append(subIDs, a.SubjectID)
		}
		if !a.TeacherID.IsZero() {
			teacherIDs = append(teacherIDs, a.TeacherID)
		}
	}
	subMap, _ := s.repo.FindSubjectsMap(ctx, subIDs)
	teacherMap, _ := s.repo.FindTeachersMap(ctx, teacherIDs)

	now := time.Now().UTC()
	list := make([]studentdto.AssessmentListDTO, 0, len(assessments))
	for _, a := range assessments {
		sName := "Mata Pelajaran"
		if sm, ok := subMap[a.SubjectID]; ok && sm != nil {
			sName = sm.Name
		}
		tName := teacherMap[a.TeacherID]

		attempt, _ := s.repo.FindAttemptByStudent(ctx, scope.SchoolID, a.ID, scope.StudentID)
		status := "available"
		attemptID := ""
		var score *float64

		if attempt != nil {
			attemptID = attempt.ID.Hex()
			if attempt.Status == model.AttemptInProgress {
				// Check timer expiration
				if now.After(attempt.ExpiresAt) {
					status = "completed"
					score = &attempt.TotalScore
				} else {
					status = "in_progress"
				}
			} else {
				status = "completed"
				score = &attempt.TotalScore
			}
		} else {
			if a.Status == model.AssessmentClosed || (a.EndAt != nil && a.EndAt.Before(now)) {
				status = "closed"
			}
		}

		list = append(list, studentdto.AssessmentListDTO{
			ID:              a.ID.Hex(),
			Title:           a.Title,
			Description:     a.Description,
			DurationMinutes: a.DurationMinutes,
			SubjectID:       a.SubjectID.Hex(),
			SubjectName:     sName,
			TeacherName:     tName,
			Category:        string(a.Category),
			StartAt:         a.StartAt,
			EndAt:           a.EndAt,
			Status:          status,
			AttemptID:       attemptID,
			Score:           score,
		})
	}
	return list, nil
}

func (s *studentServiceImpl) GetAssessmentDetail(ctx context.Context, scope *middleware.StudentScope, assessmentID bson.ObjectID) (*studentdto.AssessmentDetailDTO, error) {
	a, err := s.repo.FindAssessmentByID(ctx, scope.SchoolID, assessmentID)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, handler.ErrNotFound
	}
	if !containsObjectID(a.ClassIDs, scope.ClassID) {
		return nil, handler.NewForbiddenError("Assessment ini tidak ditujukan untuk kelas Anda")
	}

	questions, _ := s.repo.FindQuestionsByAssessmentID(ctx, assessmentID)
	attempt, _ := s.repo.FindAttemptByStudent(ctx, scope.SchoolID, assessmentID, scope.StudentID)

	subMap, _ := s.repo.FindSubjectsMap(ctx, []bson.ObjectID{a.SubjectID})
	teacherMap, _ := s.repo.FindTeachersMap(ctx, []bson.ObjectID{a.TeacherID})

	sName := "Mata Pelajaran"
	if sm, ok := subMap[a.SubjectID]; ok && sm != nil {
		sName = sm.Name
	}
	tName := teacherMap[a.TeacherID]

	attemptStatus := "not_started"
	attemptID := ""
	var startedAt, expiresAt *time.Time
	var totalScore *float64

	if attempt != nil {
		attemptStatus = string(attempt.Status)
		attemptID = attempt.ID.Hex()
		startedAt = &attempt.StartedAt
		expiresAt = &attempt.ExpiresAt
		if attempt.Status != model.AttemptInProgress {
			totalScore = &attempt.TotalScore
		}
	}

	return &studentdto.AssessmentDetailDTO{
		ID:              a.ID.Hex(),
		Title:           a.Title,
		Description:     a.Description,
		DurationMinutes: a.DurationMinutes,
		SubjectID:       a.SubjectID.Hex(),
		SubjectName:     sName,
		TeacherName:     tName,
		Category:        string(a.Category),
		StartAt:         a.StartAt,
		EndAt:           a.EndAt,
		QuestionCount:   len(questions),
		AttemptStatus:   attemptStatus,
		AttemptID:       attemptID,
		StartedAt:       startedAt,
		ExpiresAt:       expiresAt,
		TotalScore:      totalScore,
		ShowResult:      a.ShowResult,
	}, nil
}

func (s *studentServiceImpl) StartAssessment(ctx context.Context, scope *middleware.StudentScope, assessmentID bson.ObjectID) (*studentdto.AssessmentStartResponseDTO, error) {
	a, err := s.repo.FindAssessmentByID(ctx, scope.SchoolID, assessmentID)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, handler.ErrNotFound
	}
	if !containsObjectID(a.ClassIDs, scope.ClassID) {
		return nil, handler.NewForbiddenError("Assessment ini tidak ditujukan untuk kelas Anda")
	}
	if a.Status != model.AssessmentPublished {
		return nil, handler.NewForbiddenError("Assessment belum dibuka atau telah ditutup")
	}

	now := time.Now().UTC()
	if a.StartAt != nil && a.StartAt.After(now) {
		return nil, handler.NewForbiddenError("Waktu pelaksanaan assessment belum dimulai")
	}
	if a.EndAt != nil && a.EndAt.Before(now) {
		return nil, handler.NewForbiddenError("Waktu pelaksanaan assessment telah berakhir")
	}

	// Single attempt rule: if an attempt already exists
	attempt, err := s.repo.FindAttemptByStudent(ctx, scope.SchoolID, assessmentID, scope.StudentID)
	if err != nil {
		return nil, err
	}

	if attempt != nil {
		if attempt.Status != model.AttemptInProgress {
			return nil, handler.NewForbiddenError("Anda sudah menyelesaikan assessment ini")
		}
		// If in progress but expired
		if now.After(attempt.ExpiresAt) {
			attempt.Status = model.AttemptSubmitted
			_ = s.repo.UpdateAttempt(ctx, attempt)
			return nil, handler.NewForbiddenError("Waktu pengerjaan assessment Anda telah habis")
		}
	} else {
		// Create new attempt
		duration := time.Duration(a.DurationMinutes) * time.Minute
		if duration <= 0 {
			duration = 60 * time.Minute
		}
		expiresAt := now.Add(duration)
		if a.EndAt != nil && a.EndAt.Before(expiresAt) {
			expiresAt = *a.EndAt
		}

		attempt = &model.AssessmentAttempt{
			ID:           bson.NewObjectID(),
			SchoolID:     scope.SchoolID,
			AssessmentID: assessmentID,
			StudentID:    scope.StudentID,
			StartedAt:    now,
			ExpiresAt:    expiresAt,
			Status:       model.AttemptInProgress,
		}
		if err := s.repo.InsertAttempt(ctx, attempt); err != nil {
			return nil, err
		}
	}

	// Fetch questions (DO NOT return CorrectOption to student!)
	rawQuestions, err := s.repo.FindQuestionsByAssessmentID(ctx, assessmentID)
	if err != nil {
		return nil, err
	}

	// Load existing answers if re-entering in_progress attempt
	existingAnswers, _ := s.repo.FindAnswersByAttemptID(ctx, attempt.ID)
	answerMap := make(map[bson.ObjectID]string)
	for _, ans := range existingAnswers {
		answerMap[ans.QuestionID] = ans.Answer
	}

	var questions []studentdto.QuestionItemDTO
	for _, q := range rawQuestions {
		ans := answerMap[q.ID]
		questions = append(questions, studentdto.QuestionItemDTO{
			ID:      q.ID.Hex(),
			Type:    string(q.Type),
			Text:    q.Text,
			Options: q.Options,
			Weight:  q.Weight,
			Order:   q.Order,
			Answer:  ans,
		})
	}

	return &studentdto.AssessmentStartResponseDTO{
		AttemptID:       attempt.ID.Hex(),
		AssessmentID:    a.ID.Hex(),
		Title:           a.Title,
		DurationMinutes: a.DurationMinutes,
		StartedAt:       attempt.StartedAt,
		ExpiresAt:       attempt.ExpiresAt,
		Questions:       questions,
	}, nil
}

func (s *studentServiceImpl) SubmitAssessment(ctx context.Context, scope *middleware.StudentScope, assessmentID bson.ObjectID, payload *studentdto.AssessmentSubmitPayload) (*studentdto.AssessmentResultResponseDTO, error) {
	a, err := s.repo.FindAssessmentByID(ctx, scope.SchoolID, assessmentID)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, handler.ErrNotFound
	}

	attempt, err := s.repo.FindAttemptByStudent(ctx, scope.SchoolID, assessmentID, scope.StudentID)
	if err != nil {
		return nil, err
	}
	if attempt == nil {
		return nil, handler.NewForbiddenError("Anda belum memulai assessment ini")
	}
	if attempt.Status != model.AttemptInProgress {
		return nil, handler.NewForbiddenError("Assessment sudah dikumpulkan sebelumnya")
	}

	now := time.Now().UTC()
	questions, err := s.repo.FindQuestionsByAssessmentID(ctx, assessmentID)
	if err != nil {
		return nil, err
	}

	qMap := make(map[string]model.AssessmentQuestion)
	var totalWeight float64
	for _, q := range questions {
		qMap[q.ID.Hex()] = q
		totalWeight += q.Weight
	}
	if totalWeight <= 0 {
		totalWeight = float64(len(questions))
	}

	var autoScore float64
	hasEssay := false
	var answersResult []studentdto.AnswerResultDTO

	if payload != nil && len(payload.Answers) > 0 {
		for _, ans := range payload.Answers {
			q, exists := qMap[ans.QuestionID]
			if !exists {
				continue
			}

			qOID, _ := bson.ObjectIDFromHex(ans.QuestionID)
			ansModel := &model.AssessmentAnswer{
				AttemptID:  attempt.ID,
				QuestionID: qOID,
				Answer:     ans.Answer,
				SchoolID:   scope.SchoolID,
			}

			if q.Type == model.QuestionMCQ {
				isCorrect := false
				var score float64
				if q.CorrectOption != nil {
					// 0-based option check: student may send option text or "0", "1", etc.
					expectedText := ""
					if *q.CorrectOption >= 0 && *q.CorrectOption < len(q.Options) {
						expectedText = q.Options[*q.CorrectOption]
					}
					idxStr := fmt.Sprintf("%d", *q.CorrectOption)
					if strings.TrimSpace(ans.Answer) == expectedText || strings.TrimSpace(ans.Answer) == idxStr {
						isCorrect = true
						weight := q.Weight
						if weight <= 0 {
							weight = 1.0
						}
						score = weight
						autoScore += score
					}
				}
				ansModel.IsCorrect = &isCorrect
				ansModel.Score = score
			} else {
				// Essay: teacher will grade manually later
				hasEssay = true
				ansModel.Score = 0
			}

			_ = s.repo.SaveAnswer(ctx, ansModel)

			// Prepare answer result if needed
			answersResult = append(answersResult, studentdto.AnswerResultDTO{
				QuestionID:    q.ID.Hex(),
				Type:          string(q.Type),
				Text:          q.Text,
				Options:       q.Options,
				StudentAnswer: ans.Answer,
				Score:         ansModel.Score,
				MaxWeight:     q.Weight,
			})
		}
	}

	// Normalize score to scale 0 - 100
	normalizedAutoScore := (autoScore / totalWeight) * 100.0
	if normalizedAutoScore > 100.0 {
		normalizedAutoScore = 100.0
	}

	status := model.AttemptSubmitted
	if !hasEssay {
		status = model.AttemptGraded
	}

	attempt.Status = status
	attempt.AutoScore = normalizedAutoScore
	attempt.ManualScore = 0
	attempt.TotalScore = normalizedAutoScore
	attempt.SubmittedAt = &now
	if err := s.repo.UpdateAttempt(ctx, attempt); err != nil {
		return nil, err
	}

	subMap, _ := s.repo.FindSubjectsMap(ctx, []bson.ObjectID{a.SubjectID})
	sName := "Mata Pelajaran"
	if sm, ok := subMap[a.SubjectID]; ok && sm != nil {
		sName = sm.Name
	}

	return &studentdto.AssessmentResultResponseDTO{
		AssessmentID: a.ID.Hex(),
		Title:        a.Title,
		SubjectName:  sName,
		TotalScore:   normalizedAutoScore,
		AutoScore:    normalizedAutoScore,
		ManualScore:  0,
		Status:       string(status),
		SubmittedAt:  &now,
		Answers:      answersResult,
	}, nil
}

func (s *studentServiceImpl) GetAssessmentResult(ctx context.Context, scope *middleware.StudentScope, assessmentID bson.ObjectID) (*studentdto.AssessmentResultResponseDTO, error) {
	a, err := s.repo.FindAssessmentByID(ctx, scope.SchoolID, assessmentID)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, handler.ErrNotFound
	}

	attempt, err := s.repo.FindAttemptByStudent(ctx, scope.SchoolID, assessmentID, scope.StudentID)
	if err != nil {
		return nil, err
	}
	if attempt == nil || attempt.Status == model.AttemptInProgress {
		return nil, handler.NewForbiddenError("Assessment belum dikumpulkan")
	}

	// Result visibility rule: show if assessment is closed OR showResult is true
	now := time.Now().UTC()
	isClosed := a.Status == model.AssessmentClosed || (a.EndAt != nil && a.EndAt.Before(now))
	if !isClosed && !a.ShowResult {
		return nil, handler.NewForbiddenError("Hasil assessment belum dipublikasikan oleh guru")
	}

	answers, _ := s.repo.FindAnswersByAttemptID(ctx, attempt.ID)
	questions, _ := s.repo.FindQuestionsByAssessmentID(ctx, assessmentID)
	qMap := make(map[bson.ObjectID]model.AssessmentQuestion)
	for _, q := range questions {
		qMap[q.ID] = q
	}

	var answerResults []studentdto.AnswerResultDTO
	for _, ans := range answers {
		q, ok := qMap[ans.QuestionID]
		if !ok {
			continue
		}
		answerResults = append(answerResults, studentdto.AnswerResultDTO{
			QuestionID:    q.ID.Hex(),
			Type:          string(q.Type),
			Text:          q.Text,
			Options:       q.Options,
			StudentAnswer: ans.Answer,
			CorrectOption: q.CorrectOption,
			IsCorrect:     ans.IsCorrect,
			Score:         ans.Score,
			MaxWeight:     q.Weight,
			Feedback:      ans.Feedback,
		})
	}

	subMap, _ := s.repo.FindSubjectsMap(ctx, []bson.ObjectID{a.SubjectID})
	sName := "Mata Pelajaran"
	if sm, ok := subMap[a.SubjectID]; ok && sm != nil {
		sName = sm.Name
	}

	return &studentdto.AssessmentResultResponseDTO{
		AssessmentID: a.ID.Hex(),
		Title:        a.Title,
		SubjectName:  sName,
		TotalScore:   attempt.TotalScore,
		AutoScore:    attempt.AutoScore,
		ManualScore:  attempt.ManualScore,
		Status:       string(attempt.Status),
		SubmittedAt:  attempt.SubmittedAt,
		Answers:      answerResults,
	}, nil
}

// -------------------------------------------------------------
// 6. Grades (Rumus: (tugas + assessment_harian) * 60% + UTS * 20% + UAS * 20%)
// -------------------------------------------------------------

func (s *studentServiceImpl) GetGrades(ctx context.Context, scope *middleware.StudentScope) ([]studentdto.SubjectGradeDTO, error) {
	// Strict filter: only student's own grades in active academic year
	grades, err := s.repo.FindGradesByStudent(ctx, scope.SchoolID, scope.StudentID, scope.AcademicYearID)
	if err != nil {
		return nil, err
	}

	// Group grades by SubjectID
	type subjectAccum struct {
		assignmentScores []float64
		assessmentScores []float64
		projectScores    []float64
		utsScore         *float64
		uasScore         *float64
	}
	accumMap := make(map[bson.ObjectID]*subjectAccum)
	var subjectIDs []bson.ObjectID

	for _, g := range grades {
		acc, ok := accumMap[g.SubjectID]
		if !ok {
			acc = &subjectAccum{}
			accumMap[g.SubjectID] = acc
			subjectIDs = append(subjectIDs, g.SubjectID)
		}
		switch g.Type {
		case model.GradeAssignment:
			acc.assignmentScores = append(acc.assignmentScores, g.Score)
		case model.GradeAssessment:
			acc.assessmentScores = append(acc.assessmentScores, g.Score)
		case model.GradeProject:
			acc.projectScores = append(acc.projectScores, g.Score)
		case model.GradeUTS:
			score := g.Score
			acc.utsScore = &score
		case model.GradeUAS:
			score := g.Score
			acc.uasScore = &score
		}
	}

	subMap, _ := s.repo.FindSubjectsMap(ctx, subjectIDs)

	var results []studentdto.SubjectGradeDTO = make([]studentdto.SubjectGradeDTO, 0)
	for _, sid := range subjectIDs {
		sub, ok := subMap[sid]
		if !ok || sub == nil {
			continue
		}
		acc := accumMap[sid]

		// Calculate averages
		var asgAvg, asmAvg, prjAvg, dailyAvg *float64

		if len(acc.assignmentScores) > 0 {
			v := average(acc.assignmentScores)
			asgAvg = &v
		}
		if len(acc.assessmentScores) > 0 {
			v := average(acc.assessmentScores)
			asmAvg = &v
		}
		if len(acc.projectScores) > 0 {
			v := average(acc.projectScores)
			prjAvg = &v
		}

		// Combined daily (Tugas + Assessment Harian)
		var combinedScores []float64
		combinedScores = append(combinedScores, acc.assignmentScores...)
		combinedScores = append(combinedScores, acc.assessmentScores...)
		combinedScores = append(combinedScores, acc.projectScores...)
		if len(combinedScores) > 0 {
			v := average(combinedScores)
			dailyAvg = &v
		}

		// Config weights (default: 60/20/20)
		cfg, _ := s.repo.FindGradingConfig(ctx, scope.SchoolID, scope.AcademicYearID, &sid)
		dailyWeight := 0.6
		utsWeight := 0.2
		uasWeight := 0.2
		if cfg != nil && (cfg.DailyWeight+cfg.UTSWeight+cfg.UASWeight) > 0 {
			tot := cfg.DailyWeight + cfg.UTSWeight + cfg.UASWeight
			dailyWeight = cfg.DailyWeight / tot
			utsWeight = cfg.UTSWeight / tot
			uasWeight = cfg.UASWeight / tot
		}

		// Final score calculation
		var finalScore *float64
		if dailyAvg != nil {
			score := *dailyAvg * dailyWeight
			if acc.utsScore != nil {
				score += *acc.utsScore * utsWeight
			}
			if acc.uasScore != nil {
				score += *acc.uasScore * uasWeight
			}
			finalScore = &score
		}

		letter := "-"
		if finalScore != nil {
			letter = getGradeLetter(*finalScore)
		}

		results = append(results, studentdto.SubjectGradeDTO{
			SubjectID:         sid.Hex(),
			SubjectName:       sub.Name,
			SubjectCode:       sub.Code,
			AssignmentAverage: asgAvg,
			AssessmentAverage: asmAvg,
			ProjectAverage:    prjAvg,
			DailyAverage:      dailyAvg,
			UTSScore:          acc.utsScore,
			UASScore:          acc.uasScore,
			FinalScore:        finalScore,
			GradeLetter:       letter,
		})
	}

	return results, nil
}

func (s *studentServiceImpl) GetSubjectGradeDetail(ctx context.Context, scope *middleware.StudentScope, subjectID bson.ObjectID) (*studentdto.SubjectGradeDetailDTO, error) {
	subMap, _ := s.repo.FindSubjectsMap(ctx, []bson.ObjectID{subjectID})
	sub, ok := subMap[subjectID]
	if !ok || sub == nil {
		return nil, handler.ErrNotFound
	}

	grades, err := s.repo.FindGradesByStudent(ctx, scope.SchoolID, scope.StudentID, scope.AcademicYearID)
	if err != nil {
		return nil, err
	}

	var asgItems, asmItems, prjItems []studentdto.GradeDetailItem
	var utsItem, uasItem *studentdto.GradeDetailItem
	var allDailyScores []float64
	var utsScore, uasScore *float64

	for _, g := range grades {
		if g.SubjectID != subjectID {
			continue
		}
		item := studentdto.GradeDetailItem{
			ID:       g.ID.Hex(),
			Type:     string(g.Type),
			Score:    g.Score,
			GradedAt: g.CreatedAt,
			Feedback: g.Reason,
		}

		switch g.Type {
		case model.GradeAssignment:
			item.Title = "Tugas"
			asgItems = append(asgItems, item)
			allDailyScores = append(allDailyScores, g.Score)
		case model.GradeAssessment:
			item.Title = "Assessment Harian"
			asmItems = append(asmItems, item)
			allDailyScores = append(allDailyScores, g.Score)
		case model.GradeProject:
			item.Title = "Projek"
			prjItems = append(prjItems, item)
			allDailyScores = append(allDailyScores, g.Score)
		case model.GradeUTS:
			item.Title = "Ujian Tengah Semester (UTS)"
			utsItem = &item
			s := g.Score
			utsScore = &s
		case model.GradeUAS:
			item.Title = "Ujian Akhir Semester (UAS)"
			uasItem = &item
			s := g.Score
			uasScore = &s
		}
	}

	var finalScore *float64
	if len(allDailyScores) > 0 {
		dailyAvg := average(allDailyScores)
		score := dailyAvg * 0.6
		if utsScore != nil {
			score += *utsScore * 0.2
		}
		if uasScore != nil {
			score += *uasScore * 0.2
		}
		finalScore = &score
	}

	letter := "-"
	if finalScore != nil {
		letter = getGradeLetter(*finalScore)
	}

	return &studentdto.SubjectGradeDetailDTO{
		SubjectID:   subjectID.Hex(),
		SubjectName: sub.Name,
		FinalScore:  finalScore,
		GradeLetter: letter,
		Assignments: asgItems,
		Assessments: asmItems,
		Projects:    prjItems,
		UTS:         utsItem,
		UAS:         uasItem,
	}, nil
}

// -------------------------------------------------------------
// 7. Schedule
// -------------------------------------------------------------

func (s *studentServiceImpl) GetSchedule(ctx context.Context, scope *middleware.StudentScope) ([]studentdto.ScheduleItemDTO, error) {
	schedules, err := s.repo.FindSchedulesByClass(ctx, scope.SchoolID, scope.ClassID, scope.AcademicYearID)
	if err != nil {
		return nil, err
	}

	var subIDs, teacherIDs []bson.ObjectID
	for _, sc := range schedules {
		if !sc.SubjectID.IsZero() {
			subIDs = append(subIDs, sc.SubjectID)
		}
		if !sc.TeacherID.IsZero() {
			teacherIDs = append(teacherIDs, sc.TeacherID)
		}
	}
	subMap, _ := s.repo.FindSubjectsMap(ctx, subIDs)
	teacherMap, _ := s.repo.FindTeachersMap(ctx, teacherIDs)

	results := make([]studentdto.ScheduleItemDTO, 0, len(schedules))
	for _, sc := range schedules {
		sName := "Mata Pelajaran"
		if sm, ok := subMap[sc.SubjectID]; ok && sm != nil {
			sName = sm.Name
		}
		tName := teacherMap[sc.TeacherID]

		results = append(results, studentdto.ScheduleItemDTO{
			ID:        sc.ID.Hex(),
			Day:       sc.Day,
			DayName:   getDayName(sc.Day),
			StartTime: sc.StartTime,
			EndTime:   sc.EndTime,
			Subject:   sName,
			Teacher:   tName,
			Room:      sc.Room,
		})
	}
	return results, nil
}

// -------------------------------------------------------------
// 8. Ranking
// -------------------------------------------------------------

func (s *studentServiceImpl) GetRanking(ctx context.Context, scope *middleware.StudentScope) (*studentdto.ClassRankingResponseDTO, error) {
	// Only students in the same class
	students, err := s.repo.FindClassStudents(ctx, scope.SchoolID, scope.ClassID)
	if err != nil {
		return nil, err
	}

	cl, _ := s.repo.FindClassByID(ctx, scope.ClassID)
	className := "Kelas"
	if cl != nil {
		className = cl.Name
	}

	var userIDs []bson.ObjectID
	for _, st := range students {
		if !st.UserID.IsZero() {
			userIDs = append(userIDs, st.UserID)
		}
	}
	userMap, _ := s.repo.FindUsersMap(ctx, userIDs)

	type studentScoreEntry struct {
		studentID bson.ObjectID
		name      string
		nis       string
		score     float64
		isSelf    bool
	}

	var entries []studentScoreEntry
	for _, st := range students {
		grades, _ := s.repo.FindGradesByStudent(ctx, scope.SchoolID, st.ID, scope.AcademicYearID)

		// Calculate average of all subjects for student
		var allScores []float64
		for _, g := range grades {
			allScores = append(allScores, g.Score)
		}
		finalAvg := 0.0
		if len(allScores) > 0 {
			finalAvg = average(allScores)
		}

		name := st.NIS
		if u, ok := userMap[st.UserID]; ok && u != nil && u.Name != "" {
			name = u.Name
		}

		entries = append(entries, studentScoreEntry{
			studentID: st.ID,
			name:      name,
			nis:       st.NIS,
			score:     finalAvg,
			isSelf:    st.ID == scope.StudentID,
		})
	}

	// Sort descending by score
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].score > entries[j].score
	})

	myRank := 0
	rankings := make([]studentdto.RankingItemDTO, 0, len(entries))
	for i, e := range entries {
		rank := i + 1
		if e.isSelf {
			myRank = rank
		}
		rankings = append(rankings, studentdto.RankingItemDTO{
			Rank:       rank,
			StudentID:  e.studentID.Hex(),
			Name:       e.name,
			NIS:        e.nis,
			FinalScore: e.score,
			IsSelf:     e.isSelf,
		})
	}

	return &studentdto.ClassRankingResponseDTO{
		ClassName: className,
		MyRank:    myRank,
		Total:     len(entries),
		Rankings:  rankings,
	}, nil
}

// -------------------------------------------------------------
// Helper Methods
// -------------------------------------------------------------

func (s *studentServiceImpl) uploadStudentFile(ctx context.Context, resourceType, studentIDStr, resourceIDStr string, fh *multipart.FileHeader) (*model.FileRef, error) {
	file, cleanFilename, mimeType, size, err := storage.ValidateAndSniffMultipart(fh)
	if err != nil {
		return nil, handler.NewValidationError(err.Error())
	}
	defer file.Close()

	kind := fmt.Sprintf("student-submissions/%s/%s/%s", studentIDStr, resourceType, resourceIDStr)
	ref, err := s.storage.Upload(ctx, kind, cleanFilename, mimeType, size, file)
	if err != nil {
		if errors.Is(err, storage.ErrStorageNotConfigured) {
			return nil, handler.ErrStorageNotConfigured
		}
		return nil, err
	}
	return ref, nil
}

func isMimeAllowed(mime string) bool {
	if strings.HasPrefix(mime, "image/") {
		return true
	}
	return allowedMIMETypes[mime]
}

func containsObjectID(slice []bson.ObjectID, val bson.ObjectID) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func average(numbers []float64) float64 {
	if len(numbers) == 0 {
		return 0
	}
	var total float64
	for _, n := range numbers {
		total += n
	}
	return total / float64(len(numbers))
}

func getGradeLetter(score float64) string {
	switch {
	case score >= 85:
		return "A"
	case score >= 75:
		return "B"
	case score >= 60:
		return "C"
	default:
		return "D"
	}
}

func getDayName(day int) string {
	switch day {
	case 1:
		return "Senin"
	case 2:
		return "Selasa"
	case 3:
		return "Rabu"
	case 4:
		return "Kamis"
	case 5:
		return "Jumat"
	case 6:
		return "Sabtu"
	case 7:
		return "Minggu"
	default:
		return ""
	}
}
