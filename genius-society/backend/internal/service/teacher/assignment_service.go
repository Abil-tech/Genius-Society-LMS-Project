package teachersvc

import (
	"context"
	"errors"
	"fmt"
	"math"
	"mime/multipart"
	"strings"
	"time"

	teacherdto "github.com/geniussociety/backend/internal/dto/teacher"
	"github.com/geniussociety/backend/internal/handler"
	"github.com/geniussociety/backend/internal/middleware"
	"github.com/geniussociety/backend/internal/model"
	teacherrepo "github.com/geniussociety/backend/internal/repository/teacher"
	"github.com/geniussociety/backend/internal/storage"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// AssignmentService defines teacher operations for assignments and submissions.
type AssignmentService interface {
	ListAssignments(ctx context.Context, scope *middleware.TeacherScope, classIDStr, subjectIDStr, search string, page, limit int64) ([]teacherdto.AssignmentListDTO, *teacherdto.AssignmentMeta, error)
	GetAssignmentDetail(ctx context.Context, scope *middleware.TeacherScope, id bson.ObjectID) (*teacherdto.AssignmentResponseDTO, error)
	CreateAssignment(ctx context.Context, scope *middleware.TeacherScope, payload *teacherdto.AssignmentCreatePayload, files []*multipart.FileHeader) (*teacherdto.AssignmentResponseDTO, error)
	UpdateAssignment(ctx context.Context, scope *middleware.TeacherScope, id bson.ObjectID, payload *teacherdto.AssignmentUpdatePayload, files []*multipart.FileHeader) (*teacherdto.AssignmentResponseDTO, error)
	DeleteAssignment(ctx context.Context, scope *middleware.TeacherScope, id bson.ObjectID) error
	GetDownloadURL(ctx context.Context, scope *middleware.TeacherScope, id bson.ObjectID, fileIndex int) (string, error)
	ListSubmissions(ctx context.Context, scope *middleware.TeacherScope, assignmentID bson.ObjectID, status, classIDStr, search string, page, limit int64) ([]teacherdto.AssignmentSubmissionDTO, *teacherdto.SubmissionsMeta, error)
	GradeSubmission(ctx context.Context, scope *middleware.TeacherScope, assignmentID, submissionID bson.ObjectID, score float64, feedback string) error
}

type assignmentServiceImpl struct {
	repo      teacherrepo.AssignmentRepository
	storage   storage.FileStorage
	notifSvc  NotificationService
	notifRepo teacherrepo.NotificationRepository
}

// NewAssignmentService creates a new AssignmentService instance.
func NewAssignmentService(repo teacherrepo.AssignmentRepository, st storage.FileStorage, notifSvc NotificationService, notifRepo teacherrepo.NotificationRepository) AssignmentService {
	return &assignmentServiceImpl{
		repo:      repo,
		storage:   st,
		notifSvc:  notifSvc,
		notifRepo: notifRepo,
	}
}

func (s *assignmentServiceImpl) ListAssignments(ctx context.Context, scope *middleware.TeacherScope, classIDStr, subjectIDStr, search string, page, limit int64) ([]teacherdto.AssignmentListDTO, *teacherdto.AssignmentMeta, error) {
	if scope == nil {
		return []teacherdto.AssignmentListDTO{}, &teacherdto.AssignmentMeta{}, nil
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

	list, total, err := s.repo.FindAssignments(ctx, scope.TeacherID, classID, subjectID, search, page, limit)
	if err != nil {
		return nil, nil, err
	}

	// Bulk lookup class and subject names
	classIDSet := make(map[bson.ObjectID]bool)
	subjectIDSet := make(map[bson.ObjectID]bool)
	for _, a := range list {
		if !a.SubjectID.IsZero() {
			subjectIDSet[a.SubjectID] = true
		}
		for _, cid := range a.ClassIDs {
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

	classMap, _ := s.repo.FindClassesMap(ctx, allClassIDs)
	subjectMap, _ := s.repo.FindSubjectsMap(ctx, allSubjectIDs)

	result := make([]teacherdto.AssignmentListDTO, 0, len(list))
	for _, a := range list {
		subName := ""
		if sub, ok := subjectMap[a.SubjectID]; ok && sub != nil {
			subName = sub.Name
		}

		var cIDs []string
		var cNames []string
		for _, cid := range a.ClassIDs {
			cIDs = append(cIDs, cid.Hex())
			if cl, ok := classMap[cid]; ok && cl != nil {
				cNames = append(cNames, cl.Name)
			}
		}

		// Count submissions and graded
		submissions, _ := s.repo.FindSubmissions(ctx, a.ID)
		latestSubmissions := make(map[bson.ObjectID]*model.AssignmentSubmission)
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

		students, _ := s.repo.FindStudentsByClassIDs(ctx, a.ClassIDs)

		result = append(result, teacherdto.AssignmentListDTO{
			ID:             a.ID.Hex(),
			Title:          a.Title,
			SubjectID:      a.SubjectID.Hex(),
			SubjectName:    subName,
			ClassIDs:       cIDs,
			ClassNames:     cNames,
			Deadline:       a.Deadline,
			TotalStudents:  int64(len(students)),
			SubmittedCount: int64(len(latestSubmissions)),
			GradedCount:    gradedCount,
			CreatedAt:      a.CreatedAt,
		})
	}

	totalPages := int64(0)
	if total > 0 {
		totalPages = int64(math.Ceil(float64(total) / float64(limit)))
	}

	meta := &teacherdto.AssignmentMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	return result, meta, nil
}

func (s *assignmentServiceImpl) GetAssignmentDetail(ctx context.Context, scope *middleware.TeacherScope, id bson.ObjectID) (*teacherdto.AssignmentResponseDTO, error) {
	if scope == nil {
		return nil, handler.ErrNotFound
	}

	a, err := s.repo.FindAssignmentByID(ctx, scope.TeacherID, id)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, handler.ErrNotFound
	}

	classMap, _ := s.repo.FindClassesMap(ctx, a.ClassIDs)
	subjectMap, _ := s.repo.FindSubjectsMap(ctx, []bson.ObjectID{a.SubjectID})

	subName := ""
	if sub, ok := subjectMap[a.SubjectID]; ok && sub != nil {
		subName = sub.Name
	}

	var cIDs []string
	var cNames []string
	for _, cid := range a.ClassIDs {
		cIDs = append(cIDs, cid.Hex())
		if cl, ok := classMap[cid]; ok && cl != nil {
			cNames = append(cNames, cl.Name)
		}
	}

	submissions, _ := s.repo.FindSubmissions(ctx, a.ID)
	latestSubmissions := make(map[bson.ObjectID]*model.AssignmentSubmission)
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

	students, _ := s.repo.FindStudentsByClassIDs(ctx, a.ClassIDs)

	return &teacherdto.AssignmentResponseDTO{
		ID:             a.ID.Hex(),
		Title:          a.Title,
		Description:    a.Description,
		SubjectID:      a.SubjectID.Hex(),
		SubjectName:    subName,
		ClassIDs:       cIDs,
		ClassNames:     cNames,
		Deadline:       a.Deadline,
		Attachments:    a.Attachments,
		AllowResubmit:  a.AllowResubmit,
		TotalStudents:  int64(len(students)),
		SubmittedCount: int64(len(latestSubmissions)),
		GradedCount:    gradedCount,
		CreatedAt:      a.CreatedAt,
		UpdatedAt:      a.UpdatedAt,
	}, nil
}

func (s *assignmentServiceImpl) CreateAssignment(ctx context.Context, scope *middleware.TeacherScope, payload *teacherdto.AssignmentCreatePayload, files []*multipart.FileHeader) (*teacherdto.AssignmentResponseDTO, error) {
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

		ref, err := s.storage.Upload(ctx, "assignments", cleanName, mimeType, size, f)
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
	asg := &model.Assignment{
		ID:            bson.NewObjectID(),
		TeacherID:     scope.TeacherID,
		SubjectID:     subjectID,
		ClassIDs:      classOIDs,
		Title:         title,
		Description:   payload.Description,
		Deadline:      payload.Deadline,
		Attachments:   uploadedRefs,
		AllowResubmit: payload.AllowResubmit,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := s.repo.CreateAssignment(ctx, asg); err != nil {
		s.cleanupUploadedKeys(ctx, uploadedKeys)
		return nil, err
	}

	// Trigger notification for students
	if s.notifSvc != nil {
		_ = s.notifSvc.NotifyClasses(ctx, classOIDs, "new_assignment", "Tugas Baru: "+title, fmt.Sprintf("Tugas baru telah diberikan dengan batas waktu %s", payload.Deadline.Format("02 Jan 2006 15:04")), "assignment", asg.ID)
	}

	return s.GetAssignmentDetail(ctx, scope, asg.ID)
}

func (s *assignmentServiceImpl) UpdateAssignment(ctx context.Context, scope *middleware.TeacherScope, id bson.ObjectID, payload *teacherdto.AssignmentUpdatePayload, files []*multipart.FileHeader) (*teacherdto.AssignmentResponseDTO, error) {
	if scope == nil {
		return nil, handler.ErrUnauthorized
	}

	asg, err := s.repo.FindAssignmentByID(ctx, scope.TeacherID, id)
	if err != nil {
		return nil, err
	}
	if asg == nil {
		return nil, handler.ErrNotFound
	}

	if payload.Title != "" {
		title := strings.TrimSpace(payload.Title)
		if len(title) < 3 || len(title) > 200 {
			return nil, handler.NewValidationError("Judul harus antara 3 hingga 200 karakter")
		}
		asg.Title = title
	}

	if payload.Description != "" {
		if len(payload.Description) > 5000 {
			return nil, handler.NewValidationError("Deskripsi maksimal 5000 karakter")
		}
		asg.Description = payload.Description
	}

	if payload.Deadline != nil && !payload.Deadline.IsZero() {
		asg.Deadline = *payload.Deadline
	}

	if payload.AllowResubmit != nil {
		asg.AllowResubmit = *payload.AllowResubmit
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

			if !scope.HasClassSubject(cid, asg.SubjectID) {
				return nil, handler.ErrNotAssigned
			}
			newClassOIDs = append(newClassOIDs, cid)
		}

		// Check if any removed class already has submissions -> 409 Conflict
		submissions, _ := s.repo.FindSubmissions(ctx, asg.ID)
		if len(submissions) > 0 {
			newClassSet := make(map[bson.ObjectID]bool)
			for _, cid := range newClassOIDs {
				newClassSet[cid] = true
			}
			students, _ := s.repo.FindStudentsByClassIDs(ctx, asg.ClassIDs)
			studentClassMap := make(map[bson.ObjectID]bson.ObjectID)
			for _, st := range students {
				studentClassMap[st.ID] = st.ClassID
			}

			for _, sub := range submissions {
				cID := studentClassMap[sub.StudentID]
				if !cID.IsZero() && !newClassSet[cID] {
					return nil, handler.NewConflictError("Tidak dapat menghapus kelas yang siswanya sudah mengumpulkan tugas")
				}
			}
		}
		asg.ClassIDs = newClassOIDs
	}

	var oldKeysToDelete []string
	if payload.KeepExistingFiles != nil && !*payload.KeepExistingFiles {
		for _, att := range asg.Attachments {
			if att.Key != "" {
				oldKeysToDelete = append(oldKeysToDelete, att.Key)
			}
		}
		asg.Attachments = nil
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

			ref, err := s.storage.Upload(ctx, "assignments", cleanName, mimeType, size, f)
			if err != nil {
				s.cleanupUploadedKeys(ctx, newUploadedKeys)
				return nil, err
			}
			newUploadedRefs = append(newUploadedRefs, *ref)
			newUploadedKeys = append(newUploadedKeys, ref.Key)
		}
		asg.Attachments = append(asg.Attachments, newUploadedRefs...)
	}

	asg.UpdatedAt = time.Now().UTC()
	if err := s.repo.UpdateAssignment(ctx, asg); err != nil {
		return nil, err
	}

	s.cleanupUploadedKeys(ctx, oldKeysToDelete)
	return s.GetAssignmentDetail(ctx, scope, asg.ID)
}

func (s *assignmentServiceImpl) DeleteAssignment(ctx context.Context, scope *middleware.TeacherScope, id bson.ObjectID) error {
	if scope == nil {
		return handler.ErrUnauthorized
	}

	asg, err := s.repo.FindAssignmentByID(ctx, scope.TeacherID, id)
	if err != nil {
		return err
	}
	if asg == nil {
		return handler.ErrNotFound
	}

	// Rule: Cannot delete assignment if already has submissions -> 409 Conflict
	count, err := s.repo.CountSubmissionsByAssignmentID(ctx, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return handler.NewConflictError("Tugas tidak dapat dihapus karena sudah memiliki pengumpulan dari siswa")
	}

	// 1. Delete from DB first
	if err := s.repo.DeleteAssignment(ctx, scope.TeacherID, id); err != nil {
		return err
	}

	// 2. Delete storage attachments afterwards
	var keys []string
	for _, att := range asg.Attachments {
		if att.Key != "" {
			keys = append(keys, att.Key)
		}
	}
	s.cleanupUploadedKeys(ctx, keys)

	return nil
}

func (s *assignmentServiceImpl) GetDownloadURL(ctx context.Context, scope *middleware.TeacherScope, id bson.ObjectID, fileIndex int) (string, error) {
	if scope == nil {
		return "", handler.ErrUnauthorized
	}

	asg, err := s.repo.FindAssignmentByID(ctx, scope.TeacherID, id)
	if err != nil {
		return "", err
	}
	if asg == nil {
		return "", handler.ErrNotFound
	}

	if fileIndex < 0 || fileIndex >= len(asg.Attachments) {
		return "", handler.ErrNotFound
	}

	targetKey := asg.Attachments[fileIndex].Key
	if targetKey == "" {
		return "", handler.ErrNotFound
	}

	return s.storage.PresignGet(ctx, targetKey, storage.MaxPresignDuration)
}

func (s *assignmentServiceImpl) ListSubmissions(ctx context.Context, scope *middleware.TeacherScope, assignmentID bson.ObjectID, status, classIDStr, search string, page, limit int64) ([]teacherdto.AssignmentSubmissionDTO, *teacherdto.SubmissionsMeta, error) {
	if scope == nil {
		return nil, nil, handler.ErrUnauthorized
	}

	asg, err := s.repo.FindAssignmentByID(ctx, scope.TeacherID, assignmentID)
	if err != nil {
		return nil, nil, err
	}
	if asg == nil {
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

	targetClasses := asg.ClassIDs
	if filterClassID != nil {
		targetClasses = []bson.ObjectID{*filterClassID}
	}

	students, err := s.repo.FindStudentsByClassIDs(ctx, targetClasses)
	if err != nil {
		return nil, nil, err
	}

	// Bulk lookup user names
	var userIDs []bson.ObjectID
	for _, st := range students {
		if !st.UserID.IsZero() {
			userIDs = append(userIDs, st.UserID)
		}
	}
	userMap, _ := s.repo.FindUsersMap(ctx, userIDs)
	classMap, _ := s.repo.FindClassesMap(ctx, asg.ClassIDs)

	submissions, err := s.repo.FindSubmissions(ctx, asg.ID)
	if err != nil {
		return nil, nil, err
	}

	// Latest submission per student
	latestSubmissions := make(map[bson.ObjectID]*model.AssignmentSubmission)
	for _, sub := range submissions {
		if existing, ok := latestSubmissions[sub.StudentID]; !ok || sub.Version > existing.Version {
			latestSubmissions[sub.StudentID] = sub
		}
	}

	now := time.Now().UTC()
	var allItems []teacherdto.AssignmentSubmissionDTO

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
			isLate = sub.SubmittedAt.After(asg.Deadline)
			if sub.GradedAt != nil && sub.Score != nil {
				itemStatus = "graded"
			} else {
				itemStatus = "ungraded"
			}
			if isLate && status == "late" {
				itemStatus = "late"
			}
		} else {
			if now.After(asg.Deadline) {
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

		dto := teacherdto.AssignmentSubmissionDTO{
			AssignmentID: asg.ID.Hex(),
			StudentID:    st.ID.Hex(),
			StudentName:  studentName,
			NIS:          st.NIS,
			ClassID:      st.ClassID.Hex(),
			ClassName:    className,
			Status:       itemStatus,
			IsLate:       isLate,
		}

		if hasSubmitted {
			dto.ID = sub.ID.Hex()
			dto.Files = sub.Files
			dto.Note = sub.Note
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

	meta := &teacherdto.SubmissionsMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	return paged, meta, nil
}

func (s *assignmentServiceImpl) GradeSubmission(ctx context.Context, scope *middleware.TeacherScope, assignmentID, submissionID bson.ObjectID, score float64, feedback string) error {
	if scope == nil {
		return handler.ErrUnauthorized
	}

	if score < 0 || score > 100 {
		return handler.NewValidationError("Nilai harus berada di rentang 0 hingga 100")
	}

	if len(feedback) > 2000 {
		return handler.NewValidationError("Catatan guru (feedback) maksimal 2000 karakter")
	}

	asg, err := s.repo.FindAssignmentByID(ctx, scope.TeacherID, assignmentID)
	if err != nil {
		return err
	}
	if asg == nil {
		return handler.ErrNotFound
	}

	sub, err := s.repo.FindSubmissionByID(ctx, submissionID)
	if err != nil {
		return err
	}
	if sub == nil || sub.AssignmentID != asg.ID {
		return handler.ErrNotFound
	}

	student, err := s.repo.FindStudentByID(ctx, sub.StudentID)
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
		SubjectID:      asg.SubjectID,
		ClassID:        student.ClassID,
		AcademicYearID: scope.AcademicYearID,
		Semester:       scope.Semester,
		TeacherID:      scope.TeacherID,
		Type:           model.GradeAssignment,
		SourceID:       asg.ID,
		Score:          score,
		Source:         "submission",
		GradedBy:       scope.TeacherID,
	}

	audit := &model.AuditLog{
		ID:         bson.NewObjectID(),
		UserID:     scope.UserID,
		Action:     "grade_assignment",
		EntityType: "assignment_submission",
		EntityID:   sub.ID.Hex(),
		Details: bson.M{
			"assignmentId": asg.ID.Hex(),
			"studentId":    student.ID.Hex(),
			"score":        score,
		},
		Timestamp: now,
	}

	if err := s.repo.GradeSubmissionInTransaction(ctx, scope.TeacherID, sub, grade, audit); err != nil {
		return err
	}

	// Trigger notification "nilai keluar" to the student
	if s.notifRepo != nil && !student.UserID.IsZero() {
		notif := &model.Notification{
			ID:        bson.NewObjectID(),
			UserID:    student.UserID,
			Type:      "grade_released",
			Title:     "Nilai Keluar: " + asg.Title,
			Body:      fmt.Sprintf("Tugas '%s' telah dinilai. Nilai Anda: %.1f", asg.Title, score),
			RefType:   "assignment",
			RefID:     &asg.ID,
			IsRead:    false,
			CreatedAt: now,
			UpdatedAt: now,
		}
		_ = s.notifRepo.InsertMany(ctx, []*model.Notification{notif})
	}

	return nil
}

func (s *assignmentServiceImpl) cleanupUploadedKeys(ctx context.Context, keys []string) {
	for _, k := range keys {
		if k != "" && s.storage != nil {
			_ = s.storage.Delete(ctx, k)
		}
	}
}
