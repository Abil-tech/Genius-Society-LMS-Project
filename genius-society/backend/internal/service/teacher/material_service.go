package teachersvc

import (
	"context"
	"errors"
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

// MaterialService defines teacher operations for learning materials.
type MaterialService interface {
	ListMaterials(ctx context.Context, scope *middleware.TeacherScope, classIDStr, subjectIDStr, search string, page, limit int64) ([]teacherdto.MaterialListDTO, *teacherdto.MaterialMeta, error)
	GetMaterialDetail(ctx context.Context, scope *middleware.TeacherScope, id bson.ObjectID) (*teacherdto.MaterialResponseDTO, error)
	CreateMaterial(ctx context.Context, scope *middleware.TeacherScope, payload *teacherdto.MaterialCreatePayload, files []*multipart.FileHeader) (*teacherdto.MaterialResponseDTO, error)
	UpdateMaterial(ctx context.Context, scope *middleware.TeacherScope, id bson.ObjectID, payload *teacherdto.MaterialUpdatePayload, files []*multipart.FileHeader) (*teacherdto.MaterialResponseDTO, error)
	DeleteMaterial(ctx context.Context, scope *middleware.TeacherScope, id bson.ObjectID) error
	GetDownloadURL(ctx context.Context, scope *middleware.TeacherScope, id bson.ObjectID, fileIndex int) (string, error)
}

type materialServiceImpl struct {
	repo     teacherrepo.MaterialRepository
	storage  storage.FileStorage
	notifSvc NotificationService
}

// NewMaterialService creates a new MaterialService instance.
func NewMaterialService(repo teacherrepo.MaterialRepository, st storage.FileStorage, notifSvc NotificationService) MaterialService {
	return &materialServiceImpl{
		repo:     repo,
		storage:  st,
		notifSvc: notifSvc,
	}
}

func (s *materialServiceImpl) ListMaterials(ctx context.Context, scope *middleware.TeacherScope, classIDStr, subjectIDStr, search string, page, limit int64) ([]teacherdto.MaterialListDTO, *teacherdto.MaterialMeta, error) {
	if scope == nil {
		return []teacherdto.MaterialListDTO{}, &teacherdto.MaterialMeta{}, nil
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

	list, total, err := s.repo.FindMaterials(ctx, scope.TeacherID, classID, subjectID, search, page, limit)
	if err != nil {
		return nil, nil, err
	}

	// Bulk lookup class and subject names
	classIDSet := make(map[bson.ObjectID]bool)
	subjectIDSet := make(map[bson.ObjectID]bool)
	for _, m := range list {
		if !m.SubjectID.IsZero() {
			subjectIDSet[m.SubjectID] = true
		}
		for _, cid := range m.ClassIDs {
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

	now := time.Now().UTC()
	result := make([]teacherdto.MaterialListDTO, 0, len(list))
	for _, m := range list {
		subName := ""
		if sub, ok := subjectMap[m.SubjectID]; ok && sub != nil {
			subName = sub.Name
		}

		var cIDs []string
		var cNames []string
		for _, cid := range m.ClassIDs {
			cIDs = append(cIDs, cid.Hex())
			if cl, ok := classMap[cid]; ok && cl != nil {
				cNames = append(cNames, cl.Name)
			}
		}

		isPub := m.PublishAt != nil && !m.PublishAt.After(now)
		hasFile := m.File != nil || len(m.Attachments) > 0

		result = append(result, teacherdto.MaterialListDTO{
			ID:          m.ID.Hex(),
			Title:       m.Title,
			SubjectID:   m.SubjectID.Hex(),
			SubjectName: subName,
			ClassIDs:    cIDs,
			ClassNames:  cNames,
			HasFile:     hasFile,
			HasLink:     m.Link != "",
			PublishAt:   m.PublishAt,
			IsPublished: isPub,
			CreatedAt:   m.CreatedAt,
		})
	}

	totalPages := int64(0)
	if total > 0 {
		totalPages = int64(math.Ceil(float64(total) / float64(limit)))
	}

	meta := &teacherdto.MaterialMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	return result, meta, nil
}

func (s *materialServiceImpl) GetMaterialDetail(ctx context.Context, scope *middleware.TeacherScope, id bson.ObjectID) (*teacherdto.MaterialResponseDTO, error) {
	if scope == nil {
		return nil, handler.ErrNotFound
	}

	m, err := s.repo.FindMaterialByID(ctx, scope.TeacherID, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, handler.ErrNotFound
	}

	classMap, _ := s.repo.FindClassesMap(ctx, m.ClassIDs)
	subjectMap, _ := s.repo.FindSubjectsMap(ctx, []bson.ObjectID{m.SubjectID})

	subName := ""
	if sub, ok := subjectMap[m.SubjectID]; ok && sub != nil {
		subName = sub.Name
	}

	var cIDs []string
	var cNames []string
	for _, cid := range m.ClassIDs {
		cIDs = append(cIDs, cid.Hex())
		if cl, ok := classMap[cid]; ok && cl != nil {
			cNames = append(cNames, cl.Name)
		}
	}

	now := time.Now().UTC()
	isPub := m.PublishAt != nil && !m.PublishAt.After(now)

	return &teacherdto.MaterialResponseDTO{
		ID:          m.ID.Hex(),
		Title:       m.Title,
		Description: m.Description,
		SubjectID:   m.SubjectID.Hex(),
		SubjectName: subName,
		ClassIDs:    cIDs,
		ClassNames:  cNames,
		File:        m.File,
		Attachments: m.Attachments,
		Link:        m.Link,
		PublishAt:   m.PublishAt,
		IsPublished: isPub,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}, nil
}

func (s *materialServiceImpl) CreateMaterial(ctx context.Context, scope *middleware.TeacherScope, payload *teacherdto.MaterialCreatePayload, files []*multipart.FileHeader) (*teacherdto.MaterialResponseDTO, error) {
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

	link := strings.TrimSpace(payload.Link)
	if link != "" && !strings.HasPrefix(link, "http://") && !strings.HasPrefix(link, "https://") {
		return nil, handler.NewValidationError("Link harus menggunakan protokol http:// atau https://")
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

	// Upload files if any
	var uploadedRefs []model.FileRef
	var uploadedKeys []string
	for _, fh := range files {
		f, cleanName, mimeType, size, err := storage.ValidateAndSniffMultipart(fh)
		if err != nil {
			s.cleanupUploadedKeys(ctx, uploadedKeys)
			return nil, handler.NewValidationError(err.Error())
		}
		defer f.Close()

		ref, err := s.storage.Upload(ctx, "materials", cleanName, mimeType, size, f)
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
	var publishAt *time.Time
	if payload.PublishNow {
		publishAt = &now
	} else if payload.PublishAt != nil {
		publishAt = payload.PublishAt
	}

	var primaryFile *model.FileRef
	if len(uploadedRefs) > 0 {
		primaryFile = &uploadedRefs[0]
	}

	mat := &model.Material{
		ID:          bson.NewObjectID(),
		TeacherID:   scope.TeacherID,
		SubjectID:   subjectID,
		ClassIDs:    classOIDs,
		Title:       title,
		Description: payload.Description,
		File:        primaryFile,
		Attachments: uploadedRefs,
		Link:        link,
		PublishAt:   publishAt,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.CreateMaterial(ctx, mat); err != nil {
		// Rollback uploaded files
		s.cleanupUploadedKeys(ctx, uploadedKeys)
		return nil, err
	}

	// If published, trigger notification
	if mat.PublishAt != nil && !mat.PublishAt.After(now) && s.notifSvc != nil {
		_ = s.notifSvc.NotifyClasses(ctx, classOIDs, "new_material", "Materi Baru: "+title, "Materi baru telah diunggah untuk kelas Anda", "material", mat.ID)
	}

	return s.GetMaterialDetail(ctx, scope, mat.ID)
}

func (s *materialServiceImpl) UpdateMaterial(ctx context.Context, scope *middleware.TeacherScope, id bson.ObjectID, payload *teacherdto.MaterialUpdatePayload, files []*multipart.FileHeader) (*teacherdto.MaterialResponseDTO, error) {
	if scope == nil {
		return nil, handler.ErrUnauthorized
	}

	mat, err := s.repo.FindMaterialByID(ctx, scope.TeacherID, id)
	if err != nil {
		return nil, err
	}
	if mat == nil {
		return nil, handler.ErrNotFound
	}

	if payload.Title != "" {
		title := strings.TrimSpace(payload.Title)
		if len(title) < 3 || len(title) > 200 {
			return nil, handler.NewValidationError("Judul harus antara 3 hingga 200 karakter")
		}
		mat.Title = title
	}

	if payload.Description != "" {
		if len(payload.Description) > 5000 {
			return nil, handler.NewValidationError("Deskripsi maksimal 5000 karakter")
		}
		mat.Description = payload.Description
	}

	if payload.Link != "" {
		link := strings.TrimSpace(payload.Link)
		if !strings.HasPrefix(link, "http://") && !strings.HasPrefix(link, "https://") {
			return nil, handler.NewValidationError("Link harus menggunakan protokol http:// atau https://")
		}
		mat.Link = link
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

			if !scope.HasClassSubject(cid, mat.SubjectID) {
				return nil, handler.ErrNotAssigned
			}
			classOIDs = append(classOIDs, cid)
		}
		mat.ClassIDs = classOIDs
	}

	now := time.Now().UTC()
	wasPublished := mat.PublishAt != nil && !mat.PublishAt.After(now)

	if payload.PublishNow != nil && *payload.PublishNow {
		mat.PublishAt = &now
	} else if payload.PublishAt != nil {
		mat.PublishAt = payload.PublishAt
	}

	var oldKeysToDelete []string
	if payload.KeepExistingFile != nil && !*payload.KeepExistingFile {
		if mat.File != nil && mat.File.Key != "" {
			oldKeysToDelete = append(oldKeysToDelete, mat.File.Key)
		}
		for _, att := range mat.Attachments {
			if att.Key != "" {
				oldKeysToDelete = append(oldKeysToDelete, att.Key)
			}
		}
		mat.File = nil
		mat.Attachments = nil
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

			ref, err := s.storage.Upload(ctx, "materials", cleanName, mimeType, size, f)
			if err != nil {
				s.cleanupUploadedKeys(ctx, newUploadedKeys)
				return nil, err
			}
			newUploadedRefs = append(newUploadedRefs, *ref)
			newUploadedKeys = append(newUploadedKeys, ref.Key)
		}
		mat.Attachments = append(mat.Attachments, newUploadedRefs...)
		if mat.File == nil && len(newUploadedRefs) > 0 {
			mat.File = &newUploadedRefs[0]
		}
	}

	mat.UpdatedAt = now
	if err := s.repo.UpdateMaterial(ctx, mat); err != nil {
		return nil, err
	}

	// Delete replaced old files from storage
	s.cleanupUploadedKeys(ctx, oldKeysToDelete)

	// If newly published, trigger notification
	isPublishedNow := mat.PublishAt != nil && !mat.PublishAt.After(now)
	if !wasPublished && isPublishedNow && s.notifSvc != nil {
		_ = s.notifSvc.NotifyClasses(ctx, mat.ClassIDs, "new_material", "Materi Baru: "+mat.Title, "Materi baru telah diunggah untuk kelas Anda", "material", mat.ID)
	}

	return s.GetMaterialDetail(ctx, scope, mat.ID)
}

func (s *materialServiceImpl) DeleteMaterial(ctx context.Context, scope *middleware.TeacherScope, id bson.ObjectID) error {
	if scope == nil {
		return handler.ErrUnauthorized
	}

	mat, err := s.repo.FindMaterialByID(ctx, scope.TeacherID, id)
	if err != nil {
		return err
	}
	if mat == nil {
		return handler.ErrNotFound
	}

	// 1. Delete from DB first
	if err := s.repo.DeleteMaterial(ctx, scope.TeacherID, id); err != nil {
		return err
	}

	// 2. Delete storage objects afterwards
	var keys []string
	if mat.File != nil && mat.File.Key != "" {
		keys = append(keys, mat.File.Key)
	}
	for _, att := range mat.Attachments {
		if att.Key != "" {
			keys = append(keys, att.Key)
		}
	}
	s.cleanupUploadedKeys(ctx, keys)

	return nil
}

func (s *materialServiceImpl) GetDownloadURL(ctx context.Context, scope *middleware.TeacherScope, id bson.ObjectID, fileIndex int) (string, error) {
	if scope == nil {
		return "", handler.ErrUnauthorized
	}

	mat, err := s.repo.FindMaterialByID(ctx, scope.TeacherID, id)
	if err != nil {
		return "", err
	}
	if mat == nil {
		return "", handler.ErrNotFound
	}

	var targetKey string
	if len(mat.Attachments) > fileIndex && fileIndex >= 0 {
		targetKey = mat.Attachments[fileIndex].Key
	} else if mat.File != nil && fileIndex == 0 {
		targetKey = mat.File.Key
	}

	if targetKey == "" {
		return "", handler.ErrNotFound
	}

	return s.storage.PresignGet(ctx, targetKey, storage.MaxPresignDuration)
}

func (s *materialServiceImpl) cleanupUploadedKeys(ctx context.Context, keys []string) {
	for _, k := range keys {
		if k != "" && s.storage != nil {
			_ = s.storage.Delete(ctx, k)
		}
	}
}
