package teachersvc

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	teacherdto "github.com/geniussociety/backend/internal/dto/teacher"
	"github.com/geniussociety/backend/internal/handler"
	"github.com/geniussociety/backend/internal/middleware"
	"github.com/geniussociety/backend/internal/model"
	teacherrepo "github.com/geniussociety/backend/internal/repository/teacher"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// AssessmentService defines teacher operations for assessments, questions, attempts, and grading.
type AssessmentService interface {
	ListAssessments(ctx context.Context, scope *middleware.TeacherScope, classIDStr, subjectIDStr, search string, page, limit int64) ([]teacherdto.AssessmentListDTO, *teacherdto.AssessmentMeta, error)
	GetAssessmentDetail(ctx context.Context, scope *middleware.TeacherScope, id bson.ObjectID) (*teacherdto.AssessmentResponseDTO, error)
	CreateAssessment(ctx context.Context, scope *middleware.TeacherScope, payload *teacherdto.AssessmentCreatePayload) (*teacherdto.AssessmentResponseDTO, error)
	UpdateAssessment(ctx context.Context, scope *middleware.TeacherScope, id bson.ObjectID, payload *teacherdto.AssessmentUpdatePayload) (*teacherdto.AssessmentResponseDTO, error)
	DeleteAssessment(ctx context.Context, scope *middleware.TeacherScope, id bson.ObjectID) error

	AddQuestion(ctx context.Context, scope *middleware.TeacherScope, assessmentID bson.ObjectID, payload *teacherdto.QuestionCreateUpdatePayload) (*teacherdto.QuestionTeacherDTO, error)
	UpdateQuestion(ctx context.Context, scope *middleware.TeacherScope, assessmentID, questionID bson.ObjectID, payload *teacherdto.QuestionCreateUpdatePayload) (*teacherdto.QuestionTeacherDTO, error)
	DeleteQuestion(ctx context.Context, scope *middleware.TeacherScope, assessmentID, questionID bson.ObjectID) error
	ReorderQuestions(ctx context.Context, scope *middleware.TeacherScope, assessmentID bson.ObjectID, payload *teacherdto.QuestionReorderPayload) error

	RegradeMCQ(ctx context.Context, scope *middleware.TeacherScope, assessmentID bson.ObjectID, payload *teacherdto.RegradeMCQPayload) error
	PublishAssessment(ctx context.Context, scope *middleware.TeacherScope, id bson.ObjectID) (*teacherdto.AssessmentResponseDTO, error)
	CloseAssessment(ctx context.Context, scope *middleware.TeacherScope, id bson.ObjectID) (*teacherdto.AssessmentResponseDTO, error)

	ListAttempts(ctx context.Context, scope *middleware.TeacherScope, assessmentID bson.ObjectID, status, classIDStr, search string, page, limit int64) ([]teacherdto.AssessmentAttemptDTO, *teacherdto.AttemptsMeta, error)
	GetAttemptDetail(ctx context.Context, scope *middleware.TeacherScope, assessmentID, attemptID bson.ObjectID) (*teacherdto.AttemptDetailResponseDTO, error)
	GradeEssayAnswer(ctx context.Context, scope *middleware.TeacherScope, assessmentID, answerID bson.ObjectID, score float64, feedback string) error

	FinalizeExpiredAttempt(ctx context.Context, attemptID bson.ObjectID) error
}

type assessmentServiceImpl struct {
	repo      teacherrepo.AssessmentRepository
	notifSvc  NotificationService
	notifRepo teacherrepo.NotificationRepository
}

// NewAssessmentService creates a new AssessmentService instance.
func NewAssessmentService(repo teacherrepo.AssessmentRepository, notifSvc NotificationService, notifRepo teacherrepo.NotificationRepository) AssessmentService {
	return &assessmentServiceImpl{
		repo:      repo,
		notifSvc:  notifSvc,
		notifRepo: notifRepo,
	}
}

func (s *assessmentServiceImpl) ListAssessments(ctx context.Context, scope *middleware.TeacherScope, classIDStr, subjectIDStr, search string, page, limit int64) ([]teacherdto.AssessmentListDTO, *teacherdto.AssessmentMeta, error) {
	if scope == nil {
		return []teacherdto.AssessmentListDTO{}, &teacherdto.AssessmentMeta{}, nil
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

	list, total, err := s.repo.FindAssessments(ctx, scope.SchoolID, scope.TeacherID, classID, subjectID, search, page, limit)
	if err != nil {
		return nil, nil, err
	}

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

	classMap, _ := s.repo.FindClassesMap(ctx, scope.SchoolID, allClassIDs)
	subjectMap, _ := s.repo.FindSubjectsMap(ctx, allSubjectIDs)

	result := make([]teacherdto.AssessmentListDTO, 0, len(list))
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

		questions, _ := s.repo.FindQuestions(ctx, a.ID)
		var totalScore float64
		for _, q := range questions {
			totalScore += q.Weight
		}

		attemptCount, _ := s.repo.CountAttemptsByAssessmentID(ctx, scope.SchoolID, a.ID)

		result = append(result, teacherdto.AssessmentListDTO{
			ID:              a.ID.Hex(),
			Title:           a.Title,
			SubjectID:       a.SubjectID.Hex(),
			SubjectName:     subName,
			ClassIDs:        cIDs,
			ClassNames:      cNames,
			DurationMinutes: a.DurationMinutes,
			Category:        a.Category,
			Status:          a.Status,
			StartAt:         a.StartAt,
			EndAt:           a.EndAt,
			QuestionCount:   len(questions),
			TotalScore:      totalScore,
			AttemptCount:    attemptCount,
			CreatedAt:       a.CreatedAt,
		})
	}

	totalPages := int64(0)
	if total > 0 {
		totalPages = int64(math.Ceil(float64(total) / float64(limit)))
	}

	meta := &teacherdto.AssessmentMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	return result, meta, nil
}

func (s *assessmentServiceImpl) GetAssessmentDetail(ctx context.Context, scope *middleware.TeacherScope, id bson.ObjectID) (*teacherdto.AssessmentResponseDTO, error) {
	if scope == nil {
		return nil, handler.ErrNotFound
	}

	a, err := s.repo.FindAssessmentByID(ctx, scope.SchoolID, scope.TeacherID, id)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, handler.ErrNotFound
	}

	classMap, _ := s.repo.FindClassesMap(ctx, scope.SchoolID, a.ClassIDs)
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

	questions, err := s.repo.FindQuestions(ctx, a.ID)
	if err != nil {
		return nil, err
	}

	var totalScore float64
	qDTOs := make([]teacherdto.QuestionTeacherDTO, 0, len(questions))
	for _, q := range questions {
		totalScore += q.Weight
		qDTOs = append(qDTOs, teacherdto.QuestionTeacherDTO{
			ID:            q.ID.Hex(),
			AssessmentID:  q.AssessmentID.Hex(),
			Type:          q.Type,
			Text:          q.Text,
			Options:       q.Options,
			CorrectOption: q.CorrectOption,
			Weight:        q.Weight,
			Order:         q.Order,
			CreatedAt:     q.CreatedAt,
		})
	}

	attemptCount, _ := s.repo.CountAttemptsByAssessmentID(ctx, scope.SchoolID, a.ID)

	return &teacherdto.AssessmentResponseDTO{
		ID:               a.ID.Hex(),
		Title:            a.Title,
		Description:      a.Description,
		DurationMinutes:  a.DurationMinutes,
		SubjectID:        a.SubjectID.Hex(),
		SubjectName:      subName,
		ClassIDs:         cIDs,
		ClassNames:       cNames,
		StartAt:          a.StartAt,
		EndAt:            a.EndAt,
		Category:         a.Category,
		Status:           a.Status,
		MaxAttempts:      a.MaxAttempts,
		ShowResultMode:   a.ShowResultMode,
		ShuffleQuestions: a.ShuffleQuestions,
		TotalScore:       totalScore,
		QuestionCount:    len(questions),
		Questions:        qDTOs,
		AttemptCount:     attemptCount,
		CreatedAt:        a.CreatedAt,
		UpdatedAt:        a.UpdatedAt,
	}, nil
}

func (s *assessmentServiceImpl) CreateAssessment(ctx context.Context, scope *middleware.TeacherScope, payload *teacherdto.AssessmentCreatePayload) (*teacherdto.AssessmentResponseDTO, error) {
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

	if payload.StartAt != nil && payload.EndAt != nil {
		if !payload.StartAt.Before(*payload.EndAt) {
			return nil, handler.NewValidationError("Waktu mulai (startAt) harus sebelum waktu selesai (endAt)")
		}
	}

	category := payload.Category
	if category == "" {
		category = model.CategoryDaily
	}

	maxAttempts := payload.MaxAttempts
	if maxAttempts <= 0 {
		maxAttempts = 1
	}

	showResultMode := payload.ShowResultMode
	if showResultMode == "" {
		showResultMode = "after_grading"
	}

	now := time.Now().UTC()
	as := &model.Assessment{
		ID:               bson.NewObjectID(),
		SchoolID:         scope.SchoolID,
		TeacherID:        scope.TeacherID,
		SubjectID:        subjectID,
		ClassIDs:         classOIDs,
		Title:            title,
		Description:      payload.Description,
		DurationMinutes:  payload.DurationMinutes,
		Category:         category,
		Status:           model.AssessmentDraft,
		StartAt:          payload.StartAt,
		EndAt:            payload.EndAt,
		MaxAttempts:      maxAttempts,
		ShowResultMode:   showResultMode,
		ShuffleQuestions: payload.ShuffleQuestions,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := s.repo.CreateAssessment(ctx, as); err != nil {
		return nil, err
	}

	return s.GetAssessmentDetail(ctx, scope, as.ID)
}

func (s *assessmentServiceImpl) UpdateAssessment(ctx context.Context, scope *middleware.TeacherScope, id bson.ObjectID, payload *teacherdto.AssessmentUpdatePayload) (*teacherdto.AssessmentResponseDTO, error) {
	if scope == nil {
		return nil, handler.ErrUnauthorized
	}

	as, err := s.repo.FindAssessmentByID(ctx, scope.SchoolID, scope.TeacherID, id)
	if err != nil {
		return nil, err
	}
	if as == nil {
		return nil, handler.ErrNotFound
	}

	if payload.Title != "" {
		title := strings.TrimSpace(payload.Title)
		if len(title) < 3 || len(title) > 200 {
			return nil, handler.NewValidationError("Judul harus antara 3 hingga 200 karakter")
		}
		as.Title = title
	}

	if payload.Description != "" {
		if len(payload.Description) > 5000 {
			return nil, handler.NewValidationError("Deskripsi maksimal 5000 karakter")
		}
		as.Description = payload.Description
	}

	if payload.DurationMinutes != nil {
		as.DurationMinutes = *payload.DurationMinutes
	}

	if payload.StartAt != nil {
		as.StartAt = payload.StartAt
	}
	if payload.EndAt != nil {
		as.EndAt = payload.EndAt
	}
	if as.StartAt != nil && as.EndAt != nil && !as.StartAt.Before(*as.EndAt) {
		return nil, handler.NewValidationError("Waktu mulai (startAt) harus sebelum waktu selesai (endAt)")
	}

	if payload.MaxAttempts != nil {
		as.MaxAttempts = *payload.MaxAttempts
	}
	if payload.ShowResultMode != "" {
		as.ShowResultMode = payload.ShowResultMode
	}
	if payload.ShuffleQuestions != nil {
		as.ShuffleQuestions = *payload.ShuffleQuestions
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

			if !scope.HasClassSubject(cid, as.SubjectID) {
				return nil, handler.ErrNotAssigned
			}
			newClassOIDs = append(newClassOIDs, cid)
		}
		as.ClassIDs = newClassOIDs
	}

	as.UpdatedAt = time.Now().UTC()
	if err := s.repo.UpdateAssessment(ctx, as); err != nil {
		return nil, err
	}

	return s.GetAssessmentDetail(ctx, scope, as.ID)
}

func (s *assessmentServiceImpl) DeleteAssessment(ctx context.Context, scope *middleware.TeacherScope, id bson.ObjectID) error {
	if scope == nil {
		return handler.ErrUnauthorized
	}

	as, err := s.repo.FindAssessmentByID(ctx, scope.SchoolID, scope.TeacherID, id)
	if err != nil {
		return err
	}
	if as == nil {
		return handler.ErrNotFound
	}

	// Rule: Delete only if no attempts exist, otherwise 409
	count, err := s.repo.CountAttemptsByAssessmentID(ctx, scope.SchoolID, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return handler.NewConflictError("Assessment tidak dapat dihapus karena sudah memiliki attempt pengerjaan dari siswa")
	}

	return s.repo.DeleteAssessment(ctx, scope.SchoolID, scope.TeacherID, id)
}

func (s *assessmentServiceImpl) AddQuestion(ctx context.Context, scope *middleware.TeacherScope, assessmentID bson.ObjectID, payload *teacherdto.QuestionCreateUpdatePayload) (*teacherdto.QuestionTeacherDTO, error) {
	if scope == nil {
		return nil, handler.ErrUnauthorized
	}

	as, err := s.repo.FindAssessmentByID(ctx, scope.SchoolID, scope.TeacherID, assessmentID)
	if err != nil {
		return nil, err
	}
	if as == nil {
		return nil, handler.ErrNotFound
	}

	// Rule: Setelah ada attempt: soal, bobot, kunci tidak boleh diubah (409)
	count, err := s.repo.CountAttemptsByAssessmentID(ctx, scope.SchoolID, assessmentID)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, handler.NewConflictError("Tidak dapat menambah soal setelah ada pengerjaan dari siswa")
	}

	if payload.Weight <= 0 {
		return nil, handler.NewValidationError("Bobot soal harus lebih besar dari 0")
	}

	text := strings.TrimSpace(payload.Text)
	if text == "" {
		return nil, handler.NewValidationError("Teks soal tidak boleh kosong")
	}

	qType := payload.Type
	if qType != model.QuestionMCQ && qType != model.QuestionEssay {
		return nil, handler.NewValidationError("Tipe soal harus 'mcq' atau 'essay'")
	}

	var options []string
	var correctOption *int

	if qType == model.QuestionMCQ {
		if len(payload.Options) < 2 || len(payload.Options) > 6 {
			return nil, handler.NewValidationError("Pilihan ganda harus memiliki 2 hingga 6 opsi")
		}
		options = payload.Options
		if payload.CorrectOption == nil || *payload.CorrectOption < 0 || *payload.CorrectOption >= len(options) {
			return nil, handler.NewValidationError("Kunci jawaban (correctOption) wajib dipilih dan harus salah satu opsi")
		}
		correctOption = payload.CorrectOption
	}

	existingQuestions, _ := s.repo.FindQuestions(ctx, assessmentID)
	nextOrder := len(existingQuestions) + 1
	if payload.Order > 0 {
		nextOrder = payload.Order
	}

	now := time.Now().UTC()
	q := &model.AssessmentQuestion{
		ID:            bson.NewObjectID(),
		AssessmentID:  as.ID,
		SchoolID:      scope.SchoolID,
		Type:          qType,
		Text:          text,
		Options:       options,
		CorrectOption: correctOption,
		Weight:        payload.Weight,
		Order:         nextOrder,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := s.repo.CreateQuestion(ctx, q); err != nil {
		return nil, err
	}

	return &teacherdto.QuestionTeacherDTO{
		ID:            q.ID.Hex(),
		AssessmentID:  q.AssessmentID.Hex(),
		Type:          q.Type,
		Text:          q.Text,
		Options:       q.Options,
		CorrectOption: q.CorrectOption,
		Weight:        q.Weight,
		Order:         q.Order,
		CreatedAt:     q.CreatedAt,
	}, nil
}

func (s *assessmentServiceImpl) UpdateQuestion(ctx context.Context, scope *middleware.TeacherScope, assessmentID, questionID bson.ObjectID, payload *teacherdto.QuestionCreateUpdatePayload) (*teacherdto.QuestionTeacherDTO, error) {
	if scope == nil {
		return nil, handler.ErrUnauthorized
	}

	as, err := s.repo.FindAssessmentByID(ctx, scope.SchoolID, scope.TeacherID, assessmentID)
	if err != nil {
		return nil, err
	}
	if as == nil {
		return nil, handler.ErrNotFound
	}

	// Rule: Setelah ada attempt: soal, bobot, kunci tidak boleh diubah (409)
	count, err := s.repo.CountAttemptsByAssessmentID(ctx, scope.SchoolID, assessmentID)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, handler.NewConflictError("Soal tidak dapat diubah setelah ada siswa yang mengerjakan. Gunakan fitur regrade MCQ untuk mengoreksi kunci jawaban.")
	}

	q, err := s.repo.FindQuestionByID(ctx, assessmentID, questionID)
	if err != nil {
		return nil, err
	}
	if q == nil {
		return nil, handler.ErrNotFound
	}

	if payload.Weight > 0 {
		q.Weight = payload.Weight
	}
	if strings.TrimSpace(payload.Text) != "" {
		q.Text = strings.TrimSpace(payload.Text)
	}

	if q.Type == model.QuestionMCQ {
		if len(payload.Options) > 0 {
			if len(payload.Options) < 2 || len(payload.Options) > 6 {
				return nil, handler.NewValidationError("Pilihan ganda harus memiliki 2 hingga 6 opsi")
			}
			q.Options = payload.Options
		}
		if payload.CorrectOption != nil {
			if *payload.CorrectOption < 0 || *payload.CorrectOption >= len(q.Options) {
				return nil, handler.NewValidationError("Kunci jawaban tidak valid untuk daftar opsi")
			}
			q.CorrectOption = payload.CorrectOption
		}
	}

	q.UpdatedAt = time.Now().UTC()
	if err := s.repo.UpdateQuestion(ctx, q); err != nil {
		return nil, err
	}

	return &teacherdto.QuestionTeacherDTO{
		ID:            q.ID.Hex(),
		AssessmentID:  q.AssessmentID.Hex(),
		Type:          q.Type,
		Text:          q.Text,
		Options:       q.Options,
		CorrectOption: q.CorrectOption,
		Weight:        q.Weight,
		Order:         q.Order,
		CreatedAt:     q.CreatedAt,
	}, nil
}

func (s *assessmentServiceImpl) DeleteQuestion(ctx context.Context, scope *middleware.TeacherScope, assessmentID, questionID bson.ObjectID) error {
	if scope == nil {
		return handler.ErrUnauthorized
	}

	as, err := s.repo.FindAssessmentByID(ctx, scope.SchoolID, scope.TeacherID, assessmentID)
	if err != nil {
		return err
	}
	if as == nil {
		return handler.ErrNotFound
	}

	count, err := s.repo.CountAttemptsByAssessmentID(ctx, scope.SchoolID, assessmentID)
	if err != nil {
		return err
	}
	if count > 0 {
		return handler.NewConflictError("Soal tidak dapat dihapus setelah ada pengerjaan dari siswa")
	}

	return s.repo.DeleteQuestion(ctx, assessmentID, questionID)
}

func (s *assessmentServiceImpl) ReorderQuestions(ctx context.Context, scope *middleware.TeacherScope, assessmentID bson.ObjectID, payload *teacherdto.QuestionReorderPayload) error {
	if scope == nil {
		return handler.ErrUnauthorized
	}

	as, err := s.repo.FindAssessmentByID(ctx, scope.SchoolID, scope.TeacherID, assessmentID)
	if err != nil {
		return err
	}
	if as == nil {
		return handler.ErrNotFound
	}

	count, err := s.repo.CountAttemptsByAssessmentID(ctx, scope.SchoolID, assessmentID)
	if err != nil {
		return err
	}
	if count > 0 {
		return handler.NewConflictError("Urutan soal tidak dapat diubah setelah ada pengerjaan dari siswa")
	}

	orderMap := make(map[bson.ObjectID]int)
	for _, it := range payload.Questions {
		qID, err := bson.ObjectIDFromHex(it.ID)
		if err != nil {
			return handler.NewValidationError("ID soal tidak valid: " + it.ID)
		}
		orderMap[qID] = it.Order
	}

	return s.repo.ReorderQuestions(ctx, assessmentID, orderMap)
}

func (s *assessmentServiceImpl) RegradeMCQ(ctx context.Context, scope *middleware.TeacherScope, assessmentID bson.ObjectID, payload *teacherdto.RegradeMCQPayload) error {
	if scope == nil {
		return handler.ErrUnauthorized
	}

	as, err := s.repo.FindAssessmentByID(ctx, scope.SchoolID, scope.TeacherID, assessmentID)
	if err != nil {
		return err
	}
	if as == nil {
		return handler.ErrNotFound
	}

	qID, err := bson.ObjectIDFromHex(payload.QuestionID)
	if err != nil {
		return handler.NewValidationError("ID soal tidak valid")
	}

	question, err := s.repo.FindQuestionByID(ctx, assessmentID, qID)
	if err != nil {
		return err
	}
	if question == nil {
		return handler.ErrNotFound
	}

	if question.Type != model.QuestionMCQ {
		return handler.NewValidationError("Hanya soal pilihan ganda (MCQ) yang dapat di-regrade")
	}

	if payload.NewCorrectOption < 0 || payload.NewCorrectOption >= len(question.Options) {
		return handler.NewValidationError("Kunci jawaban baru tidak valid")
	}

	// Update question correctOption
	question.CorrectOption = &payload.NewCorrectOption

	// Find all attempts on this assessment
	attempts, err := s.repo.FindAttempts(ctx, scope.SchoolID, assessmentID)
	if err != nil {
		return err
	}

	var updatedAnswers []*model.AssessmentAnswer
	var updatedGrades []*model.Grade

	for _, att := range attempts {
		answers, err := s.repo.FindAnswersByAttemptID(ctx, att.ID)
		if err != nil {
			return err
		}

		var newAutoScore float64
		for _, ans := range answers {
			if ans.QuestionID == question.ID {
				// Check answer
				chosenIndex, parseErr := strconv.Atoi(ans.Answer)
				isCorr := (parseErr == nil && chosenIndex == payload.NewCorrectOption)
				score := 0.0
				if isCorr {
					score = question.Weight
				}
				ans.IsCorrect = &isCorr
				ans.Score = score
				updatedAnswers = append(updatedAnswers, ans)
			}

			// Add to autoScore if question was MCQ
			if ans.QuestionID == question.ID {
				newAutoScore += ans.Score
			} else {
				// Check if another MCQ
				if ans.IsCorrect != nil && *ans.IsCorrect {
					newAutoScore += ans.Score
				}
			}
		}

		att.AutoScore = newAutoScore
		att.TotalScore = att.AutoScore + att.ManualScore

		// If attempt was already graded, update grade record
		if att.Status == model.AttemptGraded {
			student, _ := s.repo.FindStudentByID(ctx, scope.SchoolID, att.StudentID)
			if student != nil {
				updatedGrades = append(updatedGrades, &model.Grade{
					StudentID:      student.ID,
					SubjectID:      as.SubjectID,
					ClassID:        student.ClassID,
					AcademicYearID: scope.AcademicYearID,
					Semester:       scope.Semester,
					Type:           model.GradeAssessment,
					SourceID:       as.ID,
					Score:          att.TotalScore,
				})
			}
		}
	}

	audit := &model.AuditLog{
		ID:         bson.NewObjectID(),
		SchoolID:   scope.SchoolID,
		UserID:     scope.UserID,
		Action:     "regrade_mcq",
		EntityType: "assessment_question",
		EntityID:   question.ID.Hex(),
		Details: bson.M{
			"assessmentId":     as.ID.Hex(),
			"newCorrectOption": payload.NewCorrectOption,
			"affectedAttempts": len(attempts),
		},
		Timestamp: time.Now().UTC(),
	}

	return s.repo.RegradeMCQInTransaction(ctx, scope.SchoolID, scope.TeacherID, question, attempts, updatedAnswers, updatedGrades, audit)
}

func (s *assessmentServiceImpl) PublishAssessment(ctx context.Context, scope *middleware.TeacherScope, id bson.ObjectID) (*teacherdto.AssessmentResponseDTO, error) {
	if scope == nil {
		return nil, handler.ErrUnauthorized
	}

	as, err := s.repo.FindAssessmentByID(ctx, scope.SchoolID, scope.TeacherID, id)
	if err != nil {
		return nil, err
	}
	if as == nil {
		return nil, handler.ErrNotFound
	}

	questions, err := s.repo.FindQuestions(ctx, id)
	if err != nil {
		return nil, err
	}
	if len(questions) == 0 {
		return nil, handler.NewValidationError("Assessment harus memiliki minimal 1 soal untuk dipublikasikan")
	}

	for i, q := range questions {
		if q.Type == model.QuestionMCQ {
			if q.CorrectOption == nil || *q.CorrectOption < 0 || *q.CorrectOption >= len(q.Options) {
				return nil, handler.NewValidationError(fmt.Sprintf("Soal nomor %d (MCQ) belum memiliki kunci jawaban yang valid", i+1))
			}
		}
	}

	if as.StartAt == nil || as.EndAt == nil {
		return nil, handler.NewValidationError("Waktu mulai (startAt) dan waktu selesai (endAt) wajib diisi untuk publikasi")
	}
	if !as.StartAt.Before(*as.EndAt) {
		return nil, handler.NewValidationError("Waktu mulai harus sebelum waktu selesai")
	}
	if as.DurationMinutes <= 0 {
		return nil, handler.NewValidationError("Durasi assessment harus lebih besar dari 0 menit")
	}

	for _, cid := range as.ClassIDs {
		if !scope.HasClassSubject(cid, as.SubjectID) {
			return nil, handler.ErrNotAssigned
		}
	}

	now := time.Now().UTC()
	as.Status = model.AssessmentPublished
	as.UpdatedAt = now

	if err := s.repo.UpdateAssessment(ctx, as); err != nil {
		return nil, err
	}

	// Trigger notification
	if s.notifSvc != nil {
		_ = s.notifSvc.NotifyClasses(ctx, scope.SchoolID, as.ClassIDs, "new_assessment", "Assessment Baru: "+as.Title, fmt.Sprintf("Assessment baru telah dipublikasikan dengan batas waktu %s", as.EndAt.Format("02 Jan 2006 15:04")), "assessment", as.ID)
	}

	return s.GetAssessmentDetail(ctx, scope, as.ID)
}

func (s *assessmentServiceImpl) CloseAssessment(ctx context.Context, scope *middleware.TeacherScope, id bson.ObjectID) (*teacherdto.AssessmentResponseDTO, error) {
	if scope == nil {
		return nil, handler.ErrUnauthorized
	}

	as, err := s.repo.FindAssessmentByID(ctx, scope.SchoolID, scope.TeacherID, id)
	if err != nil {
		return nil, err
	}
	if as == nil {
		return nil, handler.ErrNotFound
	}

	as.Status = model.AssessmentClosed
	as.UpdatedAt = time.Now().UTC()

	if err := s.repo.UpdateAssessment(ctx, as); err != nil {
		return nil, err
	}

	return s.GetAssessmentDetail(ctx, scope, as.ID)
}

func (s *assessmentServiceImpl) ListAttempts(ctx context.Context, scope *middleware.TeacherScope, assessmentID bson.ObjectID, status, classIDStr, search string, page, limit int64) ([]teacherdto.AssessmentAttemptDTO, *teacherdto.AttemptsMeta, error) {
	if scope == nil {
		return nil, nil, handler.ErrUnauthorized
	}

	as, err := s.repo.FindAssessmentByID(ctx, scope.SchoolID, scope.TeacherID, assessmentID)
	if err != nil {
		return nil, nil, err
	}
	if as == nil {
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

	targetClasses := as.ClassIDs
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
	classMap, _ := s.repo.FindClassesMap(ctx, scope.SchoolID, as.ClassIDs)

	attempts, err := s.repo.FindAttempts(ctx, scope.SchoolID, assessmentID)
	if err != nil {
		return nil, nil, err
	}

	studentMap := make(map[bson.ObjectID]*model.Student)
	for _, st := range students {
		studentMap[st.ID] = st
	}

	questions, _ := s.repo.FindQuestions(ctx, assessmentID)
	questionMap := make(map[bson.ObjectID]*model.AssessmentQuestion)
	for _, q := range questions {
		questionMap[q.ID] = q
	}

	searchLower := strings.ToLower(strings.TrimSpace(search))
	var allItems []teacherdto.AssessmentAttemptDTO

	for _, att := range attempts {
		st, ok := studentMap[att.StudentID]
		if !ok {
			// Student outside filtered class or inactive
			continue
		}

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

		if status != "" && string(att.Status) != status {
			continue
		}

		className := ""
		if cl, ok := classMap[st.ClassID]; ok && cl != nil {
			className = cl.Name
		}

		// Count pending essays
		answers, _ := s.repo.FindAnswersByAttemptID(ctx, att.ID)
		pendingEssays := 0
		for _, ans := range answers {
			if q, ok := questionMap[ans.QuestionID]; ok && q.Type == model.QuestionEssay {
				if ans.Score == 0 && ans.Feedback == "" && att.Status != model.AttemptGraded {
					pendingEssays++
				}
			}
		}

		allItems = append(allItems, teacherdto.AssessmentAttemptDTO{
			ID:            att.ID.Hex(),
			AssessmentID:  as.ID.Hex(),
			StudentID:     st.ID.Hex(),
			StudentName:   studentName,
			NIS:           st.NIS,
			ClassID:       st.ClassID.Hex(),
			ClassName:     className,
			AttemptNumber: att.AttemptNumber,
			Status:        att.Status,
			StartedAt:     att.StartedAt,
			ExpiresAt:     att.ExpiresAt,
			SubmittedAt:   att.SubmittedAt,
			AutoScore:     att.AutoScore,
			ManualScore:   att.ManualScore,
			TotalScore:    att.TotalScore,
			PendingEssays: pendingEssays,
		})
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

	meta := &teacherdto.AttemptsMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	return paged, meta, nil
}

func (s *assessmentServiceImpl) GetAttemptDetail(ctx context.Context, scope *middleware.TeacherScope, assessmentID, attemptID bson.ObjectID) (*teacherdto.AttemptDetailResponseDTO, error) {
	if scope == nil {
		return nil, handler.ErrUnauthorized
	}

	as, err := s.repo.FindAssessmentByID(ctx, scope.SchoolID, scope.TeacherID, assessmentID)
	if err != nil {
		return nil, err
	}
	if as == nil {
		return nil, handler.ErrNotFound
	}

	att, err := s.repo.FindAttemptByID(ctx, scope.SchoolID, attemptID)
	if err != nil {
		return nil, err
	}
	if att == nil || att.AssessmentID != as.ID {
		return nil, handler.ErrNotFound
	}

	st, err := s.repo.FindStudentByID(ctx, scope.SchoolID, att.StudentID)
	if err != nil {
		return nil, err
	}
	if st == nil {
		return nil, handler.ErrNotFound
	}

	userMap, _ := s.repo.FindUsersMap(ctx, []bson.ObjectID{st.UserID})
	studentName := ""
	if u, ok := userMap[st.UserID]; ok && u != nil {
		studentName = u.Name
	}

	classMap, _ := s.repo.FindClassesMap(ctx, scope.SchoolID, []bson.ObjectID{st.ClassID})
	className := ""
	if cl, ok := classMap[st.ClassID]; ok && cl != nil {
		className = cl.Name
	}

	questions, err := s.repo.FindQuestions(ctx, as.ID)
	if err != nil {
		return nil, err
	}
	questionMap := make(map[bson.ObjectID]*model.AssessmentQuestion)
	for _, q := range questions {
		questionMap[q.ID] = q
	}

	answers, err := s.repo.FindAnswersByAttemptID(ctx, att.ID)
	if err != nil {
		return nil, err
	}

	pendingEssays := 0
	answerDTOs := make([]teacherdto.AnswerDetailDTO, 0, len(answers))
	for _, ans := range answers {
		q, ok := questionMap[ans.QuestionID]
		if !ok {
			continue
		}

		if q.Type == model.QuestionEssay && ans.Score == 0 && ans.Feedback == "" && att.Status != model.AttemptGraded {
			pendingEssays++
		}

		answerDTOs = append(answerDTOs, teacherdto.AnswerDetailDTO{
			ID:            ans.ID.Hex(),
			QuestionID:    q.ID.Hex(),
			QuestionType:  q.Type,
			QuestionText:  q.Text,
			Options:       q.Options,
			CorrectOption: q.CorrectOption,
			Weight:        q.Weight,
			Answer:        ans.Answer,
			IsCorrect:     ans.IsCorrect,
			Score:         ans.Score,
			Feedback:      ans.Feedback,
			Order:         q.Order,
		})
	}

	attDTO := teacherdto.AssessmentAttemptDTO{
		ID:            att.ID.Hex(),
		AssessmentID:  as.ID.Hex(),
		StudentID:     st.ID.Hex(),
		StudentName:   studentName,
		NIS:           st.NIS,
		ClassID:       st.ClassID.Hex(),
		ClassName:     className,
		AttemptNumber: att.AttemptNumber,
		Status:        att.Status,
		StartedAt:     att.StartedAt,
		ExpiresAt:     att.ExpiresAt,
		SubmittedAt:   att.SubmittedAt,
		AutoScore:     att.AutoScore,
		ManualScore:   att.ManualScore,
		TotalScore:    att.TotalScore,
		PendingEssays: pendingEssays,
	}

	return &teacherdto.AttemptDetailResponseDTO{
		Attempt: attDTO,
		Answers: answerDTOs,
	}, nil
}

func (s *assessmentServiceImpl) GradeEssayAnswer(ctx context.Context, scope *middleware.TeacherScope, assessmentID, answerID bson.ObjectID, score float64, feedback string) error {
	if scope == nil {
		return handler.ErrUnauthorized
	}

	as, err := s.repo.FindAssessmentByID(ctx, scope.SchoolID, scope.TeacherID, assessmentID)
	if err != nil {
		return err
	}
	if as == nil {
		return handler.ErrNotFound
	}

	ans, err := s.repo.FindAnswerByID(ctx, answerID)
	if err != nil {
		return err
	}
	if ans == nil {
		return handler.ErrNotFound
	}

	att, err := s.repo.FindAttemptByID(ctx, scope.SchoolID, ans.AttemptID)
	if err != nil {
		return err
	}
	if att == nil || att.AssessmentID != as.ID {
		return handler.ErrNotFound
	}

	q, err := s.repo.FindQuestionByID(ctx, as.ID, ans.QuestionID)
	if err != nil {
		return err
	}
	if q == nil {
		return handler.ErrNotFound
	}

	if q.Type != model.QuestionEssay {
		return handler.NewValidationError("Hanya jawaban essay yang dapat dinilai manual")
	}

	if score < 0 || score > q.Weight {
		return handler.NewValidationError(fmt.Sprintf("Skor essay harus berada antara 0 hingga bobot soal (%.1f)", q.Weight))
	}

	if len(feedback) > 2000 {
		return handler.NewValidationError("Feedback guru maksimal 2000 karakter")
	}

	// Update answer
	ans.Score = score
	ans.Feedback = feedback

	// Recalculate manualScore and totalScore on attempt
	allAnswers, err := s.repo.FindAnswersByAttemptID(ctx, att.ID)
	if err != nil {
		return err
	}

	allQuestions, err := s.repo.FindQuestions(ctx, as.ID)
	if err != nil {
		return err
	}
	qMap := make(map[bson.ObjectID]*model.AssessmentQuestion)
	for _, question := range allQuestions {
		qMap[question.ID] = question
	}

	var newManualScore float64
	allEssaysGraded := true
	for _, a := range allAnswers {
		question := qMap[a.QuestionID]
		if question != nil && question.Type == model.QuestionEssay {
			if a.ID == ans.ID {
				newManualScore += score
			} else {
				newManualScore += a.Score
				if a.Score == 0 && a.Feedback == "" && att.Status != model.AttemptGraded {
					allEssaysGraded = false
				}
			}
		}
	}

	att.ManualScore = newManualScore
	att.TotalScore = att.AutoScore + att.ManualScore

	var grade *model.Grade
	if allEssaysGraded {
		att.Status = model.AttemptGraded
		student, _ := s.repo.FindStudentByID(ctx, scope.SchoolID, att.StudentID)
		if student != nil {
			grade = &model.Grade{
				StudentID:      student.ID,
				SubjectID:      as.SubjectID,
				ClassID:        student.ClassID,
				AcademicYearID: scope.AcademicYearID,
				Semester:       scope.Semester,
				TeacherID:      scope.TeacherID,
				Type:           model.GradeAssessment,
				SourceID:       as.ID,
				Score:          att.TotalScore,
				Source:         "submission",
				GradedBy:       scope.TeacherID,
			}
		}
	}

	now := time.Now().UTC()
	audit := &model.AuditLog{
		ID:         bson.NewObjectID(),
		SchoolID:   scope.SchoolID,
		UserID:     scope.UserID,
		Action:     "grade_essay",
		EntityType: "assessment_answer",
		EntityID:   ans.ID.Hex(),
		Details: bson.M{
			"attemptId":    att.ID.Hex(),
			"assessmentId": as.ID.Hex(),
			"questionId":   q.ID.Hex(),
			"score":        score,
		},
		Timestamp: now,
	}

	if err := s.repo.GradeEssayInTransaction(ctx, scope.SchoolID, scope.TeacherID, ans, att, grade, audit); err != nil {
		return err
	}

	// If all essays are graded, notify the student
	if allEssaysGraded && s.notifRepo != nil {
		student, _ := s.repo.FindStudentByID(ctx, scope.SchoolID, att.StudentID)
		if student != nil && !student.UserID.IsZero() {
			notif := &model.Notification{
				ID:        bson.NewObjectID(),
				SchoolID:  scope.SchoolID,
				UserID:    student.UserID,
				Type:      "grade_released",
				Title:     "Nilai Keluar: " + as.Title,
				Body:      fmt.Sprintf("Assessment '%s' telah selesai dinilai. Nilai akhir Anda: %.1f", as.Title, att.TotalScore),
				RefType:   "assessment",
				RefID:     &as.ID,
				IsRead:    false,
				CreatedAt: now,
				UpdatedAt: now,
			}
			_ = s.notifRepo.InsertMany(ctx, []*model.Notification{notif})
		}
	}

	return nil
}

// FinalizeExpiredAttempt auto-finalizes attempts whose timer has expired.
func (s *assessmentServiceImpl) FinalizeExpiredAttempt(ctx context.Context, attemptID bson.ObjectID) error {
	att, err := s.repo.FindAttemptByID(ctx, bson.NilObjectID, attemptID)
	if err != nil {
		return err
	}
	if att == nil || att.Status != model.AttemptInProgress {
		return nil
	}

	now := time.Now().UTC()
	if !now.After(att.ExpiresAt) {
		return nil
	}

	as, err := s.repo.FindAssessmentByID(ctx, att.SchoolID, bson.NilObjectID, att.AssessmentID)
	if err != nil || as == nil {
		return err
	}

	questions, err := s.repo.FindQuestions(ctx, as.ID)
	if err != nil {
		return err
	}
	hasEssay := false
	for _, q := range questions {
		if q.Type == model.QuestionEssay {
			hasEssay = true
			break
		}
	}

	att.SubmittedAt = &att.ExpiresAt
	if hasEssay {
		att.Status = model.AttemptSubmitted
	} else {
		att.Status = model.AttemptGraded
		att.TotalScore = att.AutoScore
	}
	att.UpdatedAt = now

	return s.repo.UpdateAttempt(ctx, att)
}
