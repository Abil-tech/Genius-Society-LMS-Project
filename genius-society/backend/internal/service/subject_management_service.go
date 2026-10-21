package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/geniussociety/backend/internal/dto"
	"github.com/geniussociety/backend/internal/model"
	"github.com/geniussociety/backend/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SubjectManagementService struct {
	subjectRepo *repository.SubjectRepository
	teacherRepo *repository.TeacherRepository
	db          *mongo.Database
}

func NewSubjectManagementService(
	subjectRepo *repository.SubjectRepository,
	teacherRepo *repository.TeacherRepository,
	db *mongo.Database,
) *SubjectManagementService {
	return &SubjectManagementService{
		subjectRepo: subjectRepo,
		teacherRepo: teacherRepo,
		db:          db,
	}
}

// ListSubjects with search, filters & pagination
func (s *SubjectManagementService) ListSubjects(
	ctx context.Context,
	search string,
	category string,
	gradeLevel string,
	status string,
	page int64,
	limit int64,
) ([]dto.SubjectItemResponse, int64, error) {
	filter := bson.M{}

	if status != "" {
		filter["status"] = status
	}

	if category != "" {
		filter["category"] = category
	}

	if gradeLevel != "" {
		filter["gradeLevel"] = gradeLevel
	}

	if search != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": search, "$options": "i"}},
			{"code": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	subjects, total, err := s.subjectRepo.List(ctx, filter, page, limit)
	if err != nil {
		return nil, 0, err
	}

	items := make([]dto.SubjectItemResponse, 0, len(subjects))
	for _, sub := range subjects {
		// Hitung guru pengampu mapel ini
		tCount, _ := s.teacherRepo.Count(ctx, bson.M{
			"subject": bson.M{"$regex": "^" + sub.Name + "$", "$options": "i"},
		})

		kkm := sub.KKM
		if kkm == 0 {
			kkm = 75 // default standard KKM
		}

		st := sub.Status
		if st == "" {
			st = "active"
		}

		cat := sub.Category
		if cat == "" {
			cat = "Umum"
		}

		items = append(items, dto.SubjectItemResponse{
			ID:           sub.ID.Hex(),
			Code:         sub.Code,
			Name:         sub.Name,
			Category:     cat,
			GradeLevel:   sub.GradeLevel,
			KKM:          kkm,
			Description:  sub.Description,
			Status:       st,
			TeacherCount: tCount,
			CreatedAt:    sub.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return items, total, nil
}

// GetSubjectDetail
func (s *SubjectManagementService) GetSubjectDetail(ctx context.Context, id bson.ObjectID) (*dto.SubjectItemResponse, error) {
	sub, err := s.subjectRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("mata pelajaran tidak ditemukan")
	}

	tCount, _ := s.teacherRepo.Count(ctx, bson.M{
		"subject": bson.M{"$regex": "^" + sub.Name + "$", "$options": "i"},
	})

	kkm := sub.KKM
	if kkm == 0 {
		kkm = 75
	}

	st := sub.Status
	if st == "" {
		st = "active"
	}

	cat := sub.Category
	if cat == "" {
		cat = "Umum"
	}

	return &dto.SubjectItemResponse{
		ID:           sub.ID.Hex(),
		Code:         sub.Code,
		Name:         sub.Name,
		Category:     cat,
		GradeLevel:   sub.GradeLevel,
		KKM:          kkm,
		Description:  sub.Description,
		Status:       st,
		TeacherCount: tCount,
		CreatedAt:    sub.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// CreateSubject
func (s *SubjectManagementService) CreateSubject(ctx context.Context, req dto.CreateSubjectRequest) (*dto.SubjectItemResponse, error) {
	if ex, _ := s.subjectRepo.FindByCode(ctx, req.Code); ex != nil {
		return nil, fmt.Errorf("kode mata pelajaran %s sudah terdaftar", req.Code)
	}

	kkm := req.KKM
	if kkm <= 0 {
		kkm = 75
	}

	st := req.Status
	if st == "" {
		st = "active"
	}

	cat := req.Category
	if cat == "" {
		cat = "Umum"
	}

	sub := &model.Subject{
		Code:        req.Code,
		Name:        req.Name,
		Category:    cat,
		GradeLevel:  req.GradeLevel,
		KKM:         kkm,
		Description: req.Description,
		Status:      st,
	}

	if err := s.subjectRepo.Create(ctx, sub); err != nil {
		return nil, err
	}

	return &dto.SubjectItemResponse{
		ID:           sub.ID.Hex(),
		Code:         sub.Code,
		Name:         sub.Name,
		Category:     sub.Category,
		GradeLevel:   sub.GradeLevel,
		KKM:          sub.KKM,
		Description:  sub.Description,
		Status:       sub.Status,
		TeacherCount: 0,
		CreatedAt:    time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

// UpdateSubject
func (s *SubjectManagementService) UpdateSubject(ctx context.Context, id bson.ObjectID, req dto.UpdateSubjectRequest) error {
	sub, err := s.subjectRepo.FindByID(ctx, id)
	if err != nil {
		return errors.New("mata pelajaran tidak ditemukan")
	}

	if req.Code != "" && req.Code != sub.Code {
		if ex, _ := s.subjectRepo.FindByCode(ctx, req.Code); ex != nil && ex.ID != id {
			return fmt.Errorf("kode mata pelajaran %s sudah digunakan", req.Code)
		}
	}

	kkm := req.KKM
	if kkm <= 0 {
		kkm = sub.KKM
	}

	update := bson.M{
		"code":        req.Code,
		"name":        req.Name,
		"category":    req.Category,
		"gradeLevel":  req.GradeLevel,
		"kkm":         kkm,
		"description": req.Description,
		"status":      req.Status,
	}

	return s.subjectRepo.Update(ctx, id, update)
}

// DeleteSubject
func (s *SubjectManagementService) DeleteSubject(ctx context.Context, id bson.ObjectID) error {
	sub, err := s.subjectRepo.FindByID(ctx, id)
	if err != nil {
		return errors.New("mata pelajaran tidak ditemukan")
	}

	// Peringatan jika masih ada guru yang mengampu
	tCount, _ := s.teacherRepo.Count(ctx, bson.M{
		"subject": bson.M{"$regex": "^" + sub.Name + "$", "$options": "i"},
	})
	if tCount > 0 {
		return fmt.Errorf("mata pelajaran tidak dapat dihapus karena masih diampu oleh %d guru", tCount)
	}

	return s.subjectRepo.Delete(ctx, id)
}

// GetStats
func (s *SubjectManagementService) GetStats(ctx context.Context) (*dto.SubjectStatsResponse, error) {
	total, _ := s.subjectRepo.Count(ctx, bson.M{})
	active, _ := s.subjectRepo.Count(ctx, bson.M{"$or": []bson.M{{"status": "active"}, {"status": bson.M{"$exists": false}}}})
	general, _ := s.subjectRepo.Count(ctx, bson.M{"category": "Umum"})
	specialty, _ := s.subjectRepo.Count(ctx, bson.M{
		"$or": []bson.M{
			{"category": "Peminatan IPA"},
			{"category": "Peminatan IPS"},
			{"category": "Muatan Lokal"},
		},
	})

	return &dto.SubjectStatsResponse{
		TotalSubjects:     total,
		ActiveSubjects:    active,
		GeneralSubjects:   general,
		SpecialtySubjects: specialty,
	}, nil
}
