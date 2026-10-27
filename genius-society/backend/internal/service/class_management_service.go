package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/geniussociety/backend/internal/dto"
	"github.com/geniussociety/backend/internal/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ClassRepository interface {
	Create(context.Context, *model.Class) error
	ExistsByCode(context.Context, bson.ObjectID, bson.ObjectID, string) (bool, error)
	FindByID(context.Context, bson.ObjectID) (*model.Class, error)
	Update(context.Context, bson.ObjectID, bson.M) error
	Archive(context.Context, bson.ObjectID) error
	List(context.Context, bson.M, int64, int64) ([]*model.Class, int64, error)
	CountByFilter(context.Context, bson.M) (int64, error)
}

type ClassService struct {
	classRepo ClassRepository
}

func NewClassService(classRepo ClassRepository) *ClassService {
	return &ClassService{classRepo: classRepo}
}

// CreateClass creates a new class after validating uniqueness.
func (s *ClassService) CreateClass(ctx context.Context, schoolID, academicYearID bson.ObjectID, classDTO *dto.ClassCreateDTO) (*model.Class, error) {
	if classDTO == nil {
		return nil, errors.New("class data is required")
	}
	// Ensure class code is unique within the school and academic year.
	exists, err := s.classRepo.ExistsByCode(ctx, schoolID, academicYearID, classDTO.Code)
	if err != nil {
		return nil, fmt.Errorf("checking class code existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("class with code %s already exists", classDTO.Code)
	}

	// Parse optional ObjectID fields
	var homeroomTeacherID bson.ObjectID
	if classDTO.HomeroomTeacherID != nil && *classDTO.HomeroomTeacherID != "" {
		var err error
		homeroomTeacherID, err = bson.ObjectIDFromHex(*classDTO.HomeroomTeacherID)
		if err != nil {
			return nil, fmt.Errorf("invalid homeroomTeacherId: %w", err)
		}
	}

	newClass := &model.Class{
		SchoolID:          schoolID,
		AcademicYearID:    academicYearID,
		Code:              classDTO.Code,
		Name:              classDTO.Name,
		EducationLevel:    classDTO.EducationLevel,
		Department:        classDTO.Department,
		HomeroomTeacherID: homeroomTeacherID,
		Capacity:          classDTO.Capacity,
		Status:            classDTO.Status,
		CreatedAt:         time.Now().UTC(),
	}
	if err := s.classRepo.Create(ctx, newClass); err != nil {
		return nil, fmt.Errorf("creating class: %w", err)
	}
	return newClass, nil
}

// GetClassDetail retrieves a class by its ID ensuring it belongs to the school.
func (s *ClassService) GetClassDetail(ctx context.Context, schoolID, classID bson.ObjectID) (*model.Class, error) {
	if classID.IsZero() {
		return nil, errors.New("invalid class ID")
	}
	class, err := s.classRepo.FindByID(ctx, classID)
	if err != nil {
		return nil, fmt.Errorf("fetching class: %w", err)
	}
	if class == nil || class.SchoolID != schoolID {
		return nil, errors.New("class not found or does not belong to this school")
	}
	return class, nil
}

// UpdateClass updates mutable fields of a class.
func (s *ClassService) UpdateClass(ctx context.Context, schoolID, classID bson.ObjectID, updateDTO *dto.ClassUpdateDTO) (*model.Class, error) {
	if classID.IsZero() {
		return nil, errors.New("invalid class ID")
	}
	existing, err := s.classRepo.FindByID(ctx, classID)
	if err != nil {
		return nil, fmt.Errorf("fetching class: %w", err)
	}
	if existing == nil || existing.SchoolID != schoolID {
		return nil, errors.New("class not found or does not belong to this school")
	}
	// Build update map
	update := bson.M{}
	if updateDTO.Name != nil {
		existing.Name = *updateDTO.Name
		update["name"] = *updateDTO.Name
	}
	if updateDTO.EducationLevel != nil {
		existing.EducationLevel = *updateDTO.EducationLevel
		update["educationLevel"] = *updateDTO.EducationLevel
	}
	if updateDTO.Department != nil {
		existing.Department = updateDTO.Department
		update["department"] = updateDTO.Department
	}
	if updateDTO.Capacity != nil {
		existing.Capacity = *updateDTO.Capacity
		update["capacity"] = *updateDTO.Capacity
	}
	if updateDTO.Status != nil {
		existing.Status = *updateDTO.Status
		update["status"] = *updateDTO.Status
	}
	if updateDTO.AcademicYearID != nil {
		ayID, err := bson.ObjectIDFromHex(*updateDTO.AcademicYearID)
		if err != nil {
			return nil, fmt.Errorf("invalid academicYearId: %w", err)
		}
		existing.AcademicYearID = ayID
		update["academicYearId"] = ayID
	}
	if updateDTO.HomeroomTeacherID != nil {
		if *updateDTO.HomeroomTeacherID == "" {
			existing.HomeroomTeacherID = bson.ObjectID{}
			update["homeroomTeacherId"] = bson.ObjectID{}
		} else {
			htID, err := bson.ObjectIDFromHex(*updateDTO.HomeroomTeacherID)
			if err != nil {
				return nil, fmt.Errorf("invalid homeroomTeacherId: %w", err)
			}
			existing.HomeroomTeacherID = htID
			update["homeroomTeacherId"] = htID
		}
	}
	if len(update) == 0 {
		return existing, nil // nothing to update
	}
	if err := s.classRepo.Update(ctx, classID, update); err != nil {
		return nil, fmt.Errorf("updating class: %w", err)
	}
	return existing, nil
}

// ArchiveClass soft-deletes (archives) a class.
func (s *ClassService) ArchiveClass(ctx context.Context, schoolID, classID bson.ObjectID) error {
	if classID.IsZero() {
		return errors.New("invalid class ID")
	}
	class, err := s.classRepo.FindByID(ctx, classID)
	if err != nil {
		return fmt.Errorf("fetching class: %w", err)
	}
	if class == nil || class.SchoolID != schoolID {
		return errors.New("class not found or does not belong to this school")
	}
	if err := s.classRepo.Archive(ctx, classID); err != nil {
		return fmt.Errorf("archiving class: %w", err)
	}
	return nil
}

// ListClasses returns paginated list of classes for a school.
func (s *ClassService) ListClasses(ctx context.Context, schoolID bson.ObjectID, filter *dto.ClassListFilter) (*dto.PaginatedResult[*model.Class], error) {
	if filter == nil {
		filter = &dto.ClassListFilter{}
	}
	// Ensure school filter is set
	filter.SchoolID = schoolID
	// Build BSON filter map
	bsonFilter := bson.M{"schoolId": schoolID}
	if filter.Code != nil && *filter.Code != "" {
		bsonFilter["code"] = *filter.Code
	}
	if filter.Name != nil && *filter.Name != "" {
		bsonFilter["name"] = bson.M{"$regex": *filter.Name, "$options": "i"}
	}
	if filter.Status != nil && *filter.Status != "" {
		bsonFilter["status"] = *filter.Status
	}
	if filter.EducationLevel != nil && *filter.EducationLevel != "" {
		bsonFilter["educationLevel"] = *filter.EducationLevel
	}
	if filter.Department != nil && *filter.Department != "" {
		bsonFilter["department"] = *filter.Department
	}
	if filter.AcademicYearID != nil && *filter.AcademicYearID != "" {
		ayID, err := bson.ObjectIDFromHex(*filter.AcademicYearID)
		if err == nil {
			bsonFilter["academicYearId"] = ayID
		}
	}
	// Set defaults for pagination if not provided
	page := filter.Page
	limit := filter.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	classes, total, err := s.classRepo.List(ctx, bsonFilter, page, limit)
	if err != nil {
		return nil, fmt.Errorf("listing classes: %w", err)
	}
	result := &dto.PaginatedResult[*model.Class]{
		Items: classes,
		Total: total,
		Page:  page,
		Limit: limit,
	}
	return result, nil
}
