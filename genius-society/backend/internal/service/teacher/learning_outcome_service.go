package teachersvc

import (
	"context"
	"strings"
	"time"

	teacherdto "github.com/geniussociety/backend/internal/dto/teacher"
	"github.com/geniussociety/backend/internal/handler"
	"github.com/geniussociety/backend/internal/middleware"
	"github.com/geniussociety/backend/internal/model"
	teacherrepo "github.com/geniussociety/backend/internal/repository/teacher"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// LearningOutcomeService defines business logic for Pencapaian Pembelajaran (PP).
type LearningOutcomeService interface {
	ListOutcomes(ctx context.Context, scope *middleware.TeacherScope, subjectIDStr, classIDStr string) ([]teacherdto.LearningOutcomeResponseDTO, error)
	GetOutcomeDetail(ctx context.Context, scope *middleware.TeacherScope, id bson.ObjectID) (*teacherdto.LearningOutcomeResponseDTO, error)
	CreateOutcome(ctx context.Context, scope *middleware.TeacherScope, payload *teacherdto.LearningOutcomeCreatePayload) (*teacherdto.LearningOutcomeResponseDTO, error)
	UpdateOutcome(ctx context.Context, scope *middleware.TeacherScope, id bson.ObjectID, payload *teacherdto.LearningOutcomeUpdatePayload) (*teacherdto.LearningOutcomeResponseDTO, error)
	DeleteOutcome(ctx context.Context, scope *middleware.TeacherScope, id bson.ObjectID) error

	MapAssignmentOutcome(ctx context.Context, scope *middleware.TeacherScope, assignmentID bson.ObjectID, outcomeIDStr *string) error
	MapProjectOutcome(ctx context.Context, scope *middleware.TeacherScope, projectID bson.ObjectID, outcomeIDStr *string) error
	MapAssessmentOutcome(ctx context.Context, scope *middleware.TeacherScope, assessmentID bson.ObjectID, outcomeIDStr *string) error
}

type learningOutcomeServiceImpl struct {
	repo teacherrepo.LearningOutcomeRepository
}

// NewLearningOutcomeService creates a new LearningOutcomeService instance.
func NewLearningOutcomeService(repo teacherrepo.LearningOutcomeRepository) LearningOutcomeService {
	return &learningOutcomeServiceImpl{repo: repo}
}

func (s *learningOutcomeServiceImpl) ListOutcomes(ctx context.Context, scope *middleware.TeacherScope, subjectIDStr, classIDStr string) ([]teacherdto.LearningOutcomeResponseDTO, error) {
	if scope == nil {
		return []teacherdto.LearningOutcomeResponseDTO{}, nil
	}

	var subjectID *bson.ObjectID
	if subjectIDStr != "" {
		sid, err := bson.ObjectIDFromHex(subjectIDStr)
		if err != nil {
			return nil, handler.NewValidationError("ID mata pelajaran tidak valid")
		}
		subjectID = &sid
	}

	var classID *bson.ObjectID
	if classIDStr != "" {
		cid, err := bson.ObjectIDFromHex(classIDStr)
		if err != nil {
			return nil, handler.NewValidationError("ID kelas tidak valid")
		}
		classID = &cid
	}

	list, err := s.repo.FindList(ctx, scope.SchoolID, scope.TeacherID, scope.AcademicYearID, scope.Semester, subjectID, classID)
	if err != nil {
		return nil, err
	}

	var allClassIDs []bson.ObjectID
	var allSubjectIDs []bson.ObjectID
	for _, lo := range list {
		if !lo.SubjectID.IsZero() {
			allSubjectIDs = append(allSubjectIDs, lo.SubjectID)
		}
		allClassIDs = append(allClassIDs, lo.ClassIDs...)
	}

	classMap, _ := s.repo.FindClassesMap(ctx, scope.SchoolID, allClassIDs)
	subjectMap, _ := s.repo.FindSubjectsMap(ctx, allSubjectIDs)

	res := make([]teacherdto.LearningOutcomeResponseDTO, 0, len(list))
	for _, lo := range list {
		subName := ""
		if sub, ok := subjectMap[lo.SubjectID]; ok && sub != nil {
			subName = sub.Name
		}

		var cIDs []string
		var cNames []string
		for _, cid := range lo.ClassIDs {
			cIDs = append(cIDs, cid.Hex())
			if cl, ok := classMap[cid]; ok && cl != nil {
				cNames = append(cNames, cl.Name)
			}
		}

		itemCount, _ := s.repo.CountItemsByOutcomeID(ctx, scope.SchoolID, lo.ID)

		res = append(res, teacherdto.LearningOutcomeResponseDTO{
			ID:          lo.ID.Hex(),
			SubjectID:   lo.SubjectID.Hex(),
			SubjectName: subName,
			ClassIDs:    cIDs,
			ClassNames:  cNames,
			Title:       lo.Title,
			Description: lo.Description,
			Order:       lo.Order,
			ItemCount:   itemCount,
			CreatedAt:   lo.CreatedAt,
			UpdatedAt:   lo.UpdatedAt,
		})
	}

	return res, nil
}

func (s *learningOutcomeServiceImpl) GetOutcomeDetail(ctx context.Context, scope *middleware.TeacherScope, id bson.ObjectID) (*teacherdto.LearningOutcomeResponseDTO, error) {
	if scope == nil {
		return nil, handler.ErrNotFound
	}

	lo, err := s.repo.FindByID(ctx, scope.SchoolID, scope.TeacherID, id)
	if err != nil {
		return nil, err
	}
	if lo == nil {
		return nil, handler.ErrNotFound
	}

	classMap, _ := s.repo.FindClassesMap(ctx, scope.SchoolID, lo.ClassIDs)
	subjectMap, _ := s.repo.FindSubjectsMap(ctx, []bson.ObjectID{lo.SubjectID})

	subName := ""
	if sub, ok := subjectMap[lo.SubjectID]; ok && sub != nil {
		subName = sub.Name
	}

	var cIDs []string
	var cNames []string
	for _, cid := range lo.ClassIDs {
		cIDs = append(cIDs, cid.Hex())
		if cl, ok := classMap[cid]; ok && cl != nil {
			cNames = append(cNames, cl.Name)
		}
	}

	itemCount, _ := s.repo.CountItemsByOutcomeID(ctx, scope.SchoolID, lo.ID)

	return &teacherdto.LearningOutcomeResponseDTO{
		ID:          lo.ID.Hex(),
		SubjectID:   lo.SubjectID.Hex(),
		SubjectName: subName,
		ClassIDs:    cIDs,
		ClassNames:  cNames,
		Title:       lo.Title,
		Description: lo.Description,
		Order:       lo.Order,
		ItemCount:   itemCount,
		CreatedAt:   lo.CreatedAt,
		UpdatedAt:   lo.UpdatedAt,
	}, nil
}

func (s *learningOutcomeServiceImpl) CreateOutcome(ctx context.Context, scope *middleware.TeacherScope, payload *teacherdto.LearningOutcomeCreatePayload) (*teacherdto.LearningOutcomeResponseDTO, error) {
	if scope == nil {
		return nil, handler.ErrUnauthorized
	}

	title := strings.TrimSpace(payload.Title)
	if len(title) < 3 || len(title) > 200 {
		return nil, handler.NewValidationError("Judul pencapaian pembelajaran harus antara 3 hingga 200 karakter")
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

	now := time.Now().UTC()
	lo := &model.LearningOutcome{
		ID:             bson.NewObjectID(),
		SchoolID:      scope.SchoolID,
		AcademicYearID: scope.AcademicYearID,
		Semester:       scope.Semester,
		SubjectID:      subjectID,
		ClassIDs:       classOIDs,
		TeacherID:      scope.TeacherID,
		Title:          title,
		TitleKey:       strings.ToLower(title),
		Description:    payload.Description,
		Order:          payload.Order,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.repo.Create(ctx, lo); err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "E11000") {
			return nil, handler.NewConflictError("Pencapaian pembelajaran dengan judul ini sudah ada pada semester ini")
		}
		return nil, err
	}

	return s.GetOutcomeDetail(ctx, scope, lo.ID)
}

func (s *learningOutcomeServiceImpl) UpdateOutcome(ctx context.Context, scope *middleware.TeacherScope, id bson.ObjectID, payload *teacherdto.LearningOutcomeUpdatePayload) (*teacherdto.LearningOutcomeResponseDTO, error) {
	if scope == nil {
		return nil, handler.ErrUnauthorized
	}

	lo, err := s.repo.FindByID(ctx, scope.SchoolID, scope.TeacherID, id)
	if err != nil {
		return nil, err
	}
	if lo == nil {
		return nil, handler.ErrNotFound
	}

	if payload.Title != "" {
		title := strings.TrimSpace(payload.Title)
		if len(title) < 3 || len(title) > 200 {
			return nil, handler.NewValidationError("Judul harus antara 3 hingga 200 karakter")
		}
		lo.Title = title
		lo.TitleKey = strings.ToLower(title)
	}

	if payload.Description != "" {
		if len(payload.Description) > 5000 {
			return nil, handler.NewValidationError("Deskripsi maksimal 5000 karakter")
		}
		lo.Description = payload.Description
	}

	if payload.Order != nil {
		lo.Order = *payload.Order
	}

	if len(payload.ClassIDs) > 0 {
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

			if !scope.HasClassSubject(cid, lo.SubjectID) {
				return nil, handler.ErrNotAssigned
			}
			classOIDs = append(classOIDs, cid)
		}
		lo.ClassIDs = classOIDs
	}

	if err := s.repo.Update(ctx, lo); err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "E11000") {
			return nil, handler.NewConflictError("Pencapaian pembelajaran dengan judul ini sudah ada pada semester ini")
		}
		return nil, err
	}

	return s.GetOutcomeDetail(ctx, scope, lo.ID)
}

func (s *learningOutcomeServiceImpl) DeleteOutcome(ctx context.Context, scope *middleware.TeacherScope, id bson.ObjectID) error {
	if scope == nil {
		return handler.ErrUnauthorized
	}

	lo, err := s.repo.FindByID(ctx, scope.SchoolID, scope.TeacherID, id)
	if err != nil {
		return err
	}
	if lo == nil {
		return handler.ErrNotFound
	}

	count, err := s.repo.CountItemsByOutcomeID(ctx, scope.SchoolID, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return handler.NewConflictError("Pencapaian pembelajaran tidak dapat dihapus karena masih memiliki tugas, projek, atau asesmen tertaut")
	}

	return s.repo.Delete(ctx, scope.SchoolID, scope.TeacherID, id)
}

func (s *learningOutcomeServiceImpl) validateOutcomeForMapping(ctx context.Context, scope *middleware.TeacherScope, outcomeID *bson.ObjectID, subjectID bson.ObjectID, itemClassIDs []bson.ObjectID) error {
	if outcomeID == nil {
		return nil
	}

	lo, err := s.repo.FindByID(ctx, scope.SchoolID, scope.TeacherID, *outcomeID)
	if err != nil {
		return err
	}
	if lo == nil {
		return handler.ErrNotFound
	}

	if lo.SubjectID != subjectID {
		return handler.NewValidationError("Mata pelajaran pencapaian pembelajaran tidak cocok dengan item")
	}

	if lo.Semester != scope.Semester {
		return handler.NewValidationError("Semester pencapaian pembelajaran tidak cocok dengan semester aktif")
	}

	// Subset check: all itemClassIDs must exist in lo.ClassIDs
	outcomeClassSet := make(map[bson.ObjectID]bool)
	for _, cid := range lo.ClassIDs {
		outcomeClassSet[cid] = true
	}
	for _, cid := range itemClassIDs {
		if !outcomeClassSet[cid] {
			return handler.NewValidationError("Kelas sasaran item harus merupakan bagian (subset) dari kelas pencapaian pembelajaran")
		}
	}

	return nil
}

func (s *learningOutcomeServiceImpl) parseOutcomeID(outcomeIDStr *string) (*bson.ObjectID, error) {
	if outcomeIDStr == nil || *outcomeIDStr == "" || *outcomeIDStr == "null" {
		return nil, nil
	}
	oid, err := bson.ObjectIDFromHex(*outcomeIDStr)
	if err != nil {
		return nil, handler.NewValidationError("ID pencapaian pembelajaran tidak valid")
	}
	return &oid, nil
}

func (s *learningOutcomeServiceImpl) MapAssignmentOutcome(ctx context.Context, scope *middleware.TeacherScope, assignmentID bson.ObjectID, outcomeIDStr *string) error {
	if scope == nil {
		return handler.ErrUnauthorized
	}

	asg, err := s.repo.FindAssignmentByID(ctx, scope.SchoolID, scope.TeacherID, assignmentID)
	if err != nil {
		return err
	}
	if asg == nil {
		return handler.ErrNotFound
	}

	targetOID, err := s.parseOutcomeID(outcomeIDStr)
	if err != nil {
		return err
	}

	if err := s.validateOutcomeForMapping(ctx, scope, targetOID, asg.SubjectID, asg.ClassIDs); err != nil {
		return err
	}

	return s.repo.UpdateAssignmentOutcome(ctx, scope.SchoolID, scope.TeacherID, assignmentID, targetOID)
}

func (s *learningOutcomeServiceImpl) MapProjectOutcome(ctx context.Context, scope *middleware.TeacherScope, projectID bson.ObjectID, outcomeIDStr *string) error {
	if scope == nil {
		return handler.ErrUnauthorized
	}

	prj, err := s.repo.FindProjectByID(ctx, scope.SchoolID, scope.TeacherID, projectID)
	if err != nil {
		return err
	}
	if prj == nil {
		return handler.ErrNotFound
	}

	targetOID, err := s.parseOutcomeID(outcomeIDStr)
	if err != nil {
		return err
	}

	if err := s.validateOutcomeForMapping(ctx, scope, targetOID, prj.SubjectID, prj.ClassIDs); err != nil {
		return err
	}

	return s.repo.UpdateProjectOutcome(ctx, scope.SchoolID, scope.TeacherID, projectID, targetOID)
}

func (s *learningOutcomeServiceImpl) MapAssessmentOutcome(ctx context.Context, scope *middleware.TeacherScope, assessmentID bson.ObjectID, outcomeIDStr *string) error {
	if scope == nil {
		return handler.ErrUnauthorized
	}

	asm, err := s.repo.FindAssessmentByID(ctx, scope.SchoolID, scope.TeacherID, assessmentID)
	if err != nil {
		return err
	}
	if asm == nil {
		return handler.ErrNotFound
	}

	if asm.Category != model.CategoryDaily {
		return handler.NewValidationError("Hanya asesmen kategori harian yang dapat ditautkan ke pencapaian pembelajaran")
	}

	targetOID, err := s.parseOutcomeID(outcomeIDStr)
	if err != nil {
		return err
	}

	if err := s.validateOutcomeForMapping(ctx, scope, targetOID, asm.SubjectID, asm.ClassIDs); err != nil {
		return err
	}

	return s.repo.UpdateAssessmentOutcome(ctx, scope.SchoolID, scope.TeacherID, assessmentID, targetOID)
}
