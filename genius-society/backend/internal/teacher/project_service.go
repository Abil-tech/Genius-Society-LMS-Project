package teacher

import (
	"context"
	"errors"
	"fmt"
	"math"
	"mime/multipart"
	"strings"
	"time"

	"github.com/geniussociety/backend/internal/handler"
	"github.com/geniussociety/backend/internal/middleware"
	"github.com/geniussociety/backend/internal/model"
	"github.com/geniussociety/backend/internal/storage"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// ProjectService defines teacher operations for projects and project submissions.
type ProjectService interface {
	ListProjects(ctx context.Context, scope *middleware.TeacherScope, classIDStr, subjectIDStr, search string, page, limit int64) ([]ProjectListDTO, *ProjectMeta, error)
	GetProjectDetail(ctx context.Context, scope *middleware.TeacherScope, id bson.ObjectID) (*ProjectResponseDTO, error)
	CreateProject(ctx context.Context, scope *middleware.TeacherScope, payload *ProjectCreatePayload, files []*multipart.FileHeader) (*ProjectResponseDTO, error)
	UpdateProject(ctx context.Context, scope *middleware.TeacherScope, id bson.ObjectID, payload *ProjectUpdatePayload, files []*multipart.FileHeader) (*ProjectResponseDTO, error)
	DeleteProject(ctx context.Context, scope *middleware.TeacherScope, id bson.ObjectID) error
	GetDownloadURL(ctx context.Context, scope *middleware.TeacherScope, id bson.ObjectID, fileIndex int) (string, error)
	ListSubmissions(ctx context.Context, scope *middleware.TeacherScope, projectID bson.ObjectID, status, classIDStr, search string, page, limit int64) ([]ProjectSubmissionDTO, *ProjectSubmissionsMeta, error)
	GradeSubmission(ctx context.Context, scope *middleware.TeacherScope, projectID, submissionID bson.ObjectID, score float64, feedback string) error
}

type projectServiceImpl struct {
	repo      ProjectRepository
	storage   storage.FileStorage
	notifSvc  NotificationService
	notifRepo NotificationRepository
}

// NewProjectService creates a new ProjectService instance.
func NewProjectService(repo ProjectRepository, st storage.FileStorage, notifSvc NotificationService, notifRepo NotificationRepository) ProjectService {
	return &projectServiceImpl{
		repo:      repo,
		storage:   st,
		notifSvc:  notifSvc,
		notifRepo: notifRepo,
	}
}

func (s *projectServiceImpl) ListProjects(ctx context.Context, scope *middleware.TeacherScope, classIDStr, subjectIDStr, search string, page, limit int64) ([]ProjectListDTO, *ProjectMeta, error) {
	if scope == nil {
		return []ProjectListDTO{}, &ProjectMeta{}, nil
	}

	var classID *bson.ObjectID
	if classIDStr != "" {
		cid, err := bson.ObjectIDFromHex(classIDStr)
		if err != nil {
			return nil, nil, handler.NewValidationError("ID kelas tidak valid")
		}
		classID = &cid
	}

	var subjectID *bson.ObjectID
	if subjectIDStr != "" {
		sid, err := bson.ObjectIDFromHex(subjectIDStr)
		if err != nil {
			return nil, nil, handler.NewValidationError("ID mata pelajaran tidak valid")
		}
		subjectID = &sid
	}

	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if page <= 0 {
		page = 1
	}

	list, total, err := s.repo.FindProjects(ctx, scope.SchoolID, scope.TeacherID, classID, subjectID, search, page, limit)
	if err != nil {
		return nil, nil, err
	}

	// Bulk lookup class and subject names
	classIDSet := make(map[bson.ObjectID]bool)
	subjectIDSet := make(map[bson.ObjectID]bool)
	for _, p := range list {
		if !p.SubjectID.IsZero() {
			subjectIDSet[p.SubjectID] = true
		}
		for _, cid := range p.ClassIDs {
			classIDSet[cid] = true
		}
	}

	var allClassIDs []bson.ObjectID
	for cid := range classIDSet {
		allClassIDs = append(allClassIDs, cid)
	}
	var allSubjectIDs []bson.ObjectID
	for sid := range subjectIDSet {
		allSubjectIDs = append(allSubjectIDs, sid)
	}

	classMap, _ := s.repo.FindClassesMap(ctx, scope.SchoolID, allClassIDs)
	subjectMap, _ := s.repo.FindSubjectsMap(ctx, allSubjectIDs)

	result := make([]ProjectListDTO, 0, len(list))
	for _, p := range list {
		subName := ""
		if sub, ok := subjectMap[p.SubjectID]; ok && sub != nil {
			subName = sub.Name
		}

		var cIDs []string
		var cNames []string
		for _, cid := range p.ClassIDs {
			cIDs = append(cIDs, cid.Hex())
			if cl, ok := classMap[cid]; ok && cl != nil {
				cNames = append(cNames, cl.Name)
			}
		}

		submissions, _ := s.repo.FindSubmissions(ctx, scope.SchoolID, p.ID)
		latestSubmissions := make(map[bson.ObjectID]*model.ProjectSubmission)
		var gradedCount int64
		for _, sub := range submissions {
			if existing, ok := latestSubmissions[sub.StudentID]; !ok || sub.Version > existing.Version {
				latestSubmissions[sub.StudentID] = sub
			}
		}
		for _, sub := range latestSubmissions {
			if sub.GradedAt != nil && sub.Score != nil {
				gradedCount++
			}
		}

		students, _ := s.repo.FindStudentsByClassIDs(ctx, scope.SchoolID, p.ClassIDs)

		result = append(result, ProjectListDTO{
			ID:             p.ID.Hex(),
			Title:          p.Title,
			SubjectID:      p.SubjectID.Hex(),
			SubjectName:    subName,
			ClassIDs:       cIDs,
			ClassNames:     cNames,
			Deadline:       p.Deadline,
			TotalStudents:  int64(len(students)),
			SubmittedCount: int64(len(latestSubmissions)),
			GradedCount:    gradedCount,
			CreatedAt:      p.CreatedAt,
		})
	}

	totalPages := int64(0)
	if total > 0 {
		totalPages = int64(math.Ceil(float64(total) / float64(limit)))
	}

	meta := &ProjectMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	return result, meta, nil
}

func (s *projectServiceImpl) GetProjectDetail(ctx context.Context, scope *middleware.TeacherScope, id bson.ObjectID) (*ProjectResponseDTO, error) {
	if scope == nil {
		return nil, handler.ErrNotFound
	}

	p, err := s.repo.FindProjectByID(ctx, scope.SchoolID, scope.TeacherID, id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, handler.ErrNotFound
	}

	classMap, _ := s.repo.FindClassesMap(ctx, scope.SchoolID, p.ClassIDs)
	subjectMap, _ := s.repo.FindSubjectsMap(ctx, []bson.ObjectID{p.SubjectID})

	subName := ""
	if sub, ok := subjectMap[p.SubjectID]; ok && sub != nil {
		subName = sub.Name
	}

	var cIDs []string
	var cNames []string
	for _, cid := range p.ClassIDs {
		cIDs = append(cIDs, cid.Hex())
		if cl, ok := classMap[cid]; ok && cl != nil {
			cNames = append(cNames, cl.Name)
		}
	}

	submissions, _ := s.repo.FindSubmissions(ctx, scope.SchoolID, p.ID)
	latestSubmissions := make(map[bson.ObjectID]*model.ProjectSubmission)
	var gradedCount int64
	for _, sub := range submissions {
		if existing, ok := latestSubmissions[sub.StudentID]; !ok || sub.Version > existing.Version {
			latestSubmissions[sub.StudentID] = sub
		}
	}
	for _, sub := range latestSubmissions {
		if sub.GradedAt != nil && sub.Score != nil {
			gradedCount++
		}
	}

	students, _ := s.repo.FindStudentsByClassIDs(ctx, scope.SchoolID, p.ClassIDs)

	return &ProjectResponseDTO{
		ID:             p.ID.Hex(),
		Title:          p.Title,
		Description:    p.Description,
		Instructions:   p.Instructions,
		SubjectID:      p.SubjectID.Hex(),
		SubjectName:    subName,
		ClassIDs:       cIDs,
		ClassNames:     cNames,
		Deadline:       p.Deadline,
		Attachments:    p.Attachments,
		AllowResubmit:  p.AllowResubmit,
		TotalStudents:  int64(len(students)),
		SubmittedCount: int64(len(latestSubmissions)),
		GradedCount:    gradedCount,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}, nil
}

func (s *projectServiceImpl) CreateProject(ctx context.Context, scope *middleware.TeacherScope, payload *ProjectCreatePayload, files []*multipart.FileHeader) (*ProjectResponseDTO, error) {
	if scope == nil {
		return nil, handler.ErrUnauthorized
	}

	title := strings.TrimSpace(payload.Title)
	if len(title) < 3 || len(title) > 200 {
		return nil, handler.NewValidationError("Judul harus antara 3 hingga 200 karakter")
	}

	if len(payload.Description) > 5000 {
		return nil, handler.NewValidationError("Deskripsi maksimal 5000 karakter")
	}

	if len(payload.Instructions) > 5000 {
		return nil, handler.NewValidationError("Instruksi projek maksimal 5000 karakter")
	}

	if payload.Deadline.IsZero() {
		return nil, handler.NewValidationError("Batas waktu pengumpulan (deadline) harus diisi")
	}

	subjectID, err := bson.ObjectIDFromHex(payload.SubjectID)
	if err != nil {
		return nil, handler.NewValidationError("ID mata pelajaran tidak valid")
	}

	if len(payload.ClassIDs) == 0 {
		return nil, handler.NewValidationError("Minimal satu kelas harus dipilih")
	}

	seenClasses := make(map[bson.ObjectID]bool)
	var classOIDs []bson.ObjectID
	for _, cStr := range payload.ClassIDs {
		cid, err := bson.ObjectIDFromHex(cStr)
		if err != nil {
			return nil, handler.NewValidationError("ID kelas tidak valid")
		}
		if seenClasses[cid] {
			continue
		}
		seenClasses[cid] = true

		if !scope.HasClassSubject(cid, subjectID) {
			return nil, handler.ErrNotAssigned
		}
		classOIDs = append(classOIDs, cid)
	}

	if len(files) > storage.MaxFilesPerResource {
		return nil, storage.ErrFileCountExceeded
	}

	var uploadedRefs []model.FileRef
	var uploadedKeys []string
	for _, fh := range files {
		f, cleanName, mimeType, size, err := storage.ValidateAndSniffMultipart(fh)
		if err != nil {
			s.cleanupUploadedKeys(ctx, uploadedKeys)
			return nil, handler.NewValidationError(err.Error())
		}
		defer f.Close()

		ref, err := s.storage.Upload(ctx, "projects", cleanName, mimeType, size, f)
		if err != nil {
			s.cleanupUploadedKeys(ctx, uploadedKeys)
			if errors.Is(err, storage.ErrStorageNotConfigured) {
				return nil, handler.ErrStorageNotConfigured
			}
			return nil, err
		}
		uploadedRefs = append(uploadedRefs, *ref)
		uploadedKeys = append(uploadedKeys, ref.Key)
	}

	now := time.Now().UTC()
	prj := &model.Project{
		ID:            bson.NewObjectID(),
		SchoolID:      scope.SchoolID,
		TeacherID:     scope.TeacherID,
		SubjectID:     subjectID,
		ClassIDs:      classOIDs,
		Title:         title,
		Description:   payload.Description,
		Instructions:  payload.Instructions,
		Deadline:      payload.Deadline,
		Attachments:   uploadedRefs,
		AllowResubmit: payload.AllowResubmit,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := s.repo.CreateProject(ctx, prj); err != nil {
		s.cleanupUploadedKeys(ctx, uploadedKeys)
		return nil, err
	}

	// Trigger notification for students
	if s.notifSvc != nil {
		_ = s.notifSvc.NotifyClasses(ctx, scope.SchoolID, classOIDs, "new_project", "Projek Baru: "+title, fmt.Sprintf("Projek baru telah diberikan dengan batas waktu %s", payload.Deadline.Format("02 Jan 2006 15:04")), "project", prj.ID)
	}

	return s.GetProjectDetail(ctx, scope, prj.ID)
}

func (s *projectServiceImpl) UpdateProject(ctx context.Context, scope *middleware.TeacherScope, id bson.ObjectID, payload *ProjectUpdatePayload, files []*multipart.FileHeader) (*ProjectResponseDTO, error) {
	if scope == nil {
		return nil, handler.ErrUnauthorized
	}

	prj, err := s.repo.FindProjectByID(ctx, scope.SchoolID, scope.TeacherID, id)
	if err != nil {
		return nil, err
	}
	if prj == nil {
		return nil, handler.ErrNotFound
	}

	if payload.Title != "" {
		title := strings.TrimSpace(payload.Title)
		if len(title) < 3 || len(title) > 200 {
			return nil, handler.NewValidationError("Judul harus antara 3 hingga 200 karakter")
		}
		prj.Title = title
	}

	if payload.Description != "" {
		if len(payload.Description) > 5000 {
			return nil, handler.NewValidationError("Deskripsi maksimal 5000 karakter")
		}
		prj.Description = payload.Description
	}

	if payload.Instructions != "" {
		if len(payload.Instructions) > 5000 {
			return nil, handler.NewValidationError("Instruksi projek maksimal 5000 karakter")
		}
		prj.Instructions = payload.Instructions
	}

	if payload.Deadline != nil && !payload.Deadline.IsZero() {
		prj.Deadline = *payload.Deadline
	}

	if payload.AllowResubmit != nil {
		prj.AllowResubmit = *payload.AllowResubmit
	}

	if len(payload.ClassIDs) > 0 {
		seenClasses := make(map[bson.ObjectID]bool)
		var newClassOIDs []bson.ObjectID
		for _, cStr := range payload.ClassIDs {
			cid, err := bson.ObjectIDFromHex(cStr)
			if err != nil {
				return nil, handler.NewValidationError("ID kelas tidak valid")
			}
			if seenClasses[cid] {
				continue
			}
			seenClasses[cid] = true

			if !scope.HasClassSubject(cid, prj.SubjectID) {
				return nil, handler.ErrNotAssigned
			}
			newClassOIDs = append(newClassOIDs, cid)
		}

		// Check if any removed class already has submissions -> 409 Conflict
		submissions, _ := s.repo.FindSubmissions(ctx, scope.SchoolID, prj.ID)
		if len(submissions) > 0 {
			newClassSet := make(map[bson.ObjectID]bool)
			for _, cid := range newClassOIDs {
				newClassSet[cid] = true
			}
			students, _ := s.repo.FindStudentsByClassIDs(ctx, scope.SchoolID, prj.ClassIDs)
			studentClassMap := make(map[bson.ObjectID]bson.ObjectID)
			for _, st := range students {
				studentClassMap[st.ID] = st.ClassID
			}

			for _, sub := range submissions {
				cID := studentClassMap[sub.StudentID]
				if !cID.IsZero() && !newClassSet[cID] {
					return nil, handler.NewConflictError("Tidak dapat menghapus kelas yang siswanya sudah mengumpulkan projek")
				}
			}
		}
		prj.ClassIDs = newClassOIDs
	}

	var oldKeysToDelete []string
	if payload.KeepExistingFiles != nil && !*payload.KeepExistingFiles {
		for _, att := range prj.Attachments {
			if att.Key != "" {
				oldKeysToDelete = append(oldKeysToDelete, att.Key)
			}
		}
		prj.Attachments = nil
	}

	if len(files) > 0 {
		if len(files) > storage.MaxFilesPerResource {
			return nil, storage.ErrFileCountExceeded
		}
		var newUploadedRefs []model.FileRef
		var newUploadedKeys []string
		for _, fh := range files {
			f, cleanName, mimeType, size, err := storage.ValidateAndSniffMultipart(fh)
			if err != nil {
				s.cleanupUploadedKeys(ctx, newUploadedKeys)
				return nil, handler.NewValidationError(err.Error())
			}
			defer f.Close()

			ref, err := s.storage.Upload(ctx, "projects", cleanName, mimeType, size, f)
			if err != nil {
				s.cleanupUploadedKeys(ctx, newUploadedKeys)
				return nil, err
			}
			newUploadedRefs = append(newUploadedRefs, *ref)
			newUploadedKeys = append(newUploadedKeys, ref.Key)
		}
		prj.Attachments = append(prj.Attachments, newUploadedRefs...)
	}

	prj.UpdatedAt = time.Now().UTC()
	if err := s.repo.UpdateProject(ctx, prj); err != nil {
		return nil, err
	}

	s.cleanupUploadedKeys(ctx, oldKeysToDelete)
	return s.GetProjectDetail(ctx, scope, prj.ID)
}

func (s *projectServiceImpl) DeleteProject(ctx context.Context, scope *middleware.TeacherScope, id bson.ObjectID) error {
	if scope == nil {
		return handler.ErrUnauthorized
	}

	prj, err := s.repo.FindProjectByID(ctx, scope.SchoolID, scope.TeacherID, id)
	if err != nil {
		return err
	}
	if prj == nil {
		return handler.ErrNotFound
	}

	count, err := s.repo.CountSubmissionsByProjectID(ctx, scope.SchoolID, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return handler.NewConflictError("Projek tidak dapat dihapus karena sudah memiliki pengumpulan dari siswa")
	}

	if err := s.repo.DeleteProject(ctx, scope.SchoolID, scope.TeacherID, id); err != nil {
		return err
	}

	var keys []string
	for _, att := range prj.Attachments {
		if att.Key != "" {
			keys = append(keys, att.Key)
		}
	}
	s.cleanupUploadedKeys(ctx, keys)

	return nil
}

func (s *projectServiceImpl) GetDownloadURL(ctx context.Context, scope *middleware.TeacherScope, id bson.ObjectID, fileIndex int) (string, error) {
	if scope == nil {
		return "", handler.ErrUnauthorized
	}

	prj, err := s.repo.FindProjectByID(ctx, scope.SchoolID, scope.TeacherID, id)
	if err != nil {
		return "", err
	}
	if prj == nil {
		return "", handler.ErrNotFound
	}

	if fileIndex < 0 || fileIndex >= len(prj.Attachments) {
		return "", handler.ErrNotFound
	}

	targetKey := prj.Attachments[fileIndex].Key
	if targetKey == "" {
		return "", handler.ErrNotFound
	}

	return s.storage.PresignGet(ctx, targetKey, storage.MaxPresignDuration)
}

func (s *projectServiceImpl) ListSubmissions(ctx context.Context, scope *middleware.TeacherScope, projectID bson.ObjectID, status, classIDStr, search string, page, limit int64) ([]ProjectSubmissionDTO, *ProjectSubmissionsMeta, error) {
	if scope == nil {
		return nil, nil, handler.ErrUnauthorized
	}

	prj, err := s.repo.FindProjectByID(ctx, scope.SchoolID, scope.TeacherID, projectID)
	if err != nil {
		return nil, nil, err
	}
	if prj == nil {
		return nil, nil, handler.ErrNotFound
	}

	var filterClassID *bson.ObjectID
	if classIDStr != "" {
		cid, err := bson.ObjectIDFromHex(classIDStr)
		if err != nil {
			return nil, nil, handler.NewValidationError("ID kelas tidak valid")
		}
		filterClassID = &cid
	}

	targetClasses := prj.ClassIDs
	if filterClassID != nil {
		targetClasses = []bson.ObjectID{*filterClassID}
	}

	students, err := s.repo.FindStudentsByClassIDs(ctx, scope.SchoolID, targetClasses)
	if err != nil {
		return nil, nil, err
	}

	var userIDs []bson.ObjectID
	for _, st := range students {
		if !st.UserID.IsZero() {
			userIDs = append(userIDs, st.UserID)
		}
	}
	userMap, _ := s.repo.FindUsersMap(ctx, userIDs)
	classMap, _ := s.repo.FindClassesMap(ctx, scope.SchoolID, prj.ClassIDs)

	submissions, err := s.repo.FindSubmissions(ctx, scope.SchoolID, prj.ID)
	if err != nil {
		return nil, nil, err
	}

	// Latest submission per student (Version tracking)
	latestSubmissions := make(map[bson.ObjectID]*model.ProjectSubmission)
	for _, sub := range submissions {
		if existing, ok := latestSubmissions[sub.StudentID]; !ok || sub.Version > existing.Version {
			latestSubmissions[sub.StudentID] = sub
		}
	}

	now := time.Now().UTC()
	var allItems []ProjectSubmissionDTO

	searchLower := strings.ToLower(strings.TrimSpace(search))

	for _, st := range students {
		studentName := ""
		if u, ok := userMap[st.UserID]; ok && u != nil {
			studentName = u.Name
		}

		if searchLower != "" {
			nameMatch := strings.Contains(strings.ToLower(studentName), searchLower)
			nisMatch := strings.Contains(strings.ToLower(st.NIS), searchLower)
			if !nameMatch && !nisMatch {
				continue
			}
		}

		className := ""
		if cl, ok := classMap[st.ClassID]; ok && cl != nil {
			className = cl.Name
		}

		sub, hasSubmitted := latestSubmissions[st.ID]

		var itemStatus string
		isLate := false

		if hasSubmitted {
			isLate = sub.SubmittedAt.After(prj.Deadline)
			if sub.GradedAt != nil && sub.Score != nil {
				itemStatus = "graded"
			} else {
				itemStatus = "ungraded"
			}
			if isLate && status == "late" {
				itemStatus = "late"
			}
		} else {
			if now.After(prj.Deadline) {
				itemStatus = "missing"
			} else {
				itemStatus = "unsubmitted"
			}
		}

		if status != "" {
			if status == "late" {
				if !isLate {
					continue
				}
			} else if status == "missing" {
				if itemStatus != "missing" {
					continue
				}
			} else if status == "graded" {
				if itemStatus != "graded" {
					continue
				}
			} else if status == "ungraded" {
				if itemStatus != "ungraded" {
					continue
				}
			}
		}

		dto := ProjectSubmissionDTO{
			ProjectID:   prj.ID.Hex(),
			StudentID:   st.ID.Hex(),
			StudentName: studentName,
			NIS:         st.NIS,
			ClassID:     st.ClassID.Hex(),
			ClassName:   className,
			Status:      itemStatus,
			IsLate:      isLate,
		}

		if hasSubmitted {
			dto.ID = sub.ID.Hex()
			dto.Files = sub.Files
			dto.SubmittedAt = &sub.SubmittedAt
			dto.Version = sub.Version
			dto.Score = sub.Score
			dto.Feedback = sub.Feedback
			dto.GradedAt = sub.GradedAt
		}

		allItems = append(allItems, dto)
	}

	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if page <= 0 {
		page = 1
	}

	total := int64(len(allItems))
	start := (page - 1) * limit
	if start > total {
		start = total
	}
	end := start + limit
	if end > total {
		end = total
	}

	paged := allItems[start:end]
	totalPages := int64(0)
	if total > 0 {
		totalPages = int64(math.Ceil(float64(total) / float64(limit)))
	}

	meta := &ProjectSubmissionsMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	return paged, meta, nil
}

func (s *projectServiceImpl) GradeSubmission(ctx context.Context, scope *middleware.TeacherScope, projectID, submissionID bson.ObjectID, score float64, feedback string) error {
	if scope == nil {
		return handler.ErrUnauthorized
	}

	if score < 0 || score > 100 {
		return handler.NewValidationError("Nilai harus berada di rentang 0 hingga 100")
	}

	if len(feedback) > 2000 {
		return handler.NewValidationError("Catatan guru (feedback) maksimal 2000 karakter")
	}

	prj, err := s.repo.FindProjectByID(ctx, scope.SchoolID, scope.TeacherID, projectID)
	if err != nil {
		return err
	}
	if prj == nil {
		return handler.ErrNotFound
	}

	sub, err := s.repo.FindSubmissionByID(ctx, scope.SchoolID, submissionID)
	if err != nil {
		return err
	}
	if sub == nil || sub.ProjectID != prj.ID {
		return handler.ErrNotFound
	}

	student, err := s.repo.FindStudentByID(ctx, scope.SchoolID, sub.StudentID)
	if err != nil {
		return err
	}
	if student == nil {
		return handler.ErrNotFound
	}

	now := time.Now().UTC()
	sub.Score = &score
	sub.Feedback = feedback
	sub.GradedAt = &now

	grade := &model.Grade{
		StudentID:      student.ID,
		SubjectID:      prj.SubjectID,
		ClassID:        student.ClassID,
		AcademicYearID: scope.AcademicYearID,
		Semester:       scope.Semester,
		TeacherID:      scope.TeacherID,
		Type:           model.GradeProject,
		SourceID:       prj.ID,
		Score:          score,
		Source:         "submission",
		GradedBy:       scope.TeacherID,
	}

	audit := &model.AuditLog{
		ID:         bson.NewObjectID(),
		SchoolID:   scope.SchoolID,
		UserID:     scope.UserID,
		Action:     "grade_project",
		EntityType: "project_submission",
		EntityID:   sub.ID.Hex(),
		Details: bson.M{
			"projectId": prj.ID.Hex(),
			"studentId": student.ID.Hex(),
			"score":     score,
			"version":   sub.Version,
		},
		Timestamp: now,
	}

	if err := s.repo.GradeSubmissionInTransaction(ctx, scope.SchoolID, scope.TeacherID, sub, grade, audit); err != nil {
		return err
	}

	// Trigger notification "nilai keluar" to the student
	if s.notifRepo != nil && !student.UserID.IsZero() {
		notif := &model.Notification{
			ID:        bson.NewObjectID(),
			SchoolID:  scope.SchoolID,
			UserID:    student.UserID,
			Type:      "grade_released",
			Title:     "Nilai Keluar: " + prj.Title,
			Body:      fmt.Sprintf("Projek '%s' telah dinilai. Nilai Anda: %.1f", prj.Title, score),
			RefType:   "project",
			RefID:     &prj.ID,
			IsRead:    false,
			CreatedAt: now,
			UpdatedAt: now,
		}
		_ = s.notifRepo.InsertMany(ctx, []*model.Notification{notif})
	}

	return nil
}

func (s *projectServiceImpl) cleanupUploadedKeys(ctx context.Context, keys []string) {
	for _, k := range keys {
		if k != "" && s.storage != nil {
			_ = s.storage.Delete(ctx, k)
		}
	}
}
