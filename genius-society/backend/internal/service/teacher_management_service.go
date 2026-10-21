package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/geniussociety/backend/internal/dto"
	"github.com/geniussociety/backend/internal/model"
	"github.com/geniussociety/backend/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

type TeacherManagementService struct {
	teacherRepo *repository.TeacherRepository
	userRepo    *repository.UserRepository
	classRepo   *repository.ClassRepository
	db          *mongo.Database
}

func NewTeacherManagementService(
	teacherRepo *repository.TeacherRepository,
	userRepo *repository.UserRepository,
	classRepo *repository.ClassRepository,
	db *mongo.Database,
) *TeacherManagementService {
	return &TeacherManagementService{
		teacherRepo: teacherRepo,
		userRepo:    userRepo,
		classRepo:   classRepo,
		db:          db,
	}
}

// ListTeachers with search, filters & pagination
func (s *TeacherManagementService) ListTeachers(
	ctx context.Context,
	search string,
	subject string,
	status string,
	isWalas *bool,
	page int64,
	limit int64,
) ([]dto.TeacherItemResponse, int64, error) {
	filter := bson.M{}

	if status != "" {
		filter["status"] = status
	}

	if subject != "" {
		filter["subject"] = subject
	}

	if isWalas != nil {
		filter["isWalas"] = *isWalas
	}

	if search != "" {
		searchFilter := bson.M{
			"nip": bson.M{"$regex": search, "$options": "i"},
		}

		userCur, err := s.db.Collection("users").Find(ctx, bson.M{
			"name": bson.M{"$regex": search, "$options": "i"},
			"role": model.RoleTeacher,
		})
		if err == nil {
			var matchingUsers []model.User
			if err := userCur.All(ctx, &matchingUsers); err == nil && len(matchingUsers) > 0 {
				userIDs := make([]bson.ObjectID, len(matchingUsers))
				for i, u := range matchingUsers {
					userIDs[i] = u.ID
				}
				filter["$or"] = []bson.M{
					searchFilter,
					{"userId": bson.M{"$in": userIDs}},
				}
			} else {
				filter["nip"] = searchFilter["nip"]
			}
		}
	}

	teachers, total, err := s.teacherRepo.List(ctx, filter, page, limit)
	if err != nil {
		return nil, 0, err
	}

	// Kumpulkan UserID dan ClassID untuk populate
	userIDs := make([]bson.ObjectID, 0)
	classIDs := make([]bson.ObjectID, 0)
	for _, t := range teachers {
		if !t.UserID.IsZero() {
			userIDs = append(userIDs, t.UserID)
		}
		if len(t.WalasClassIDs) > 0 {
			classIDs = append(classIDs, t.WalasClassIDs...)
		}
	}

	userMap := make(map[string]model.User)
	if len(userIDs) > 0 {
		cur, err := s.db.Collection("users").Find(ctx, bson.M{"_id": bson.M{"$in": userIDs}})
		if err == nil {
			var users []model.User
			_ = cur.All(ctx, &users)
			for _, u := range users {
				userMap[u.ID.Hex()] = u
			}
		}
	}

	classMap := make(map[string]string)
	if len(classIDs) > 0 {
		cur, err := s.db.Collection("classes").Find(ctx, bson.M{"_id": bson.M{"$in": classIDs}})
		if err == nil {
			var classes []model.Class
			_ = cur.All(ctx, &classes)
			for _, c := range classes {
				classMap[c.ID.Hex()] = c.Name
			}
		}
	}

	items := make([]dto.TeacherItemResponse, 0, len(teachers))
	for _, t := range teachers {
		item := dto.TeacherItemResponse{
			ID:        t.ID.Hex(),
			NIP:       t.NIP,
			Gender:    t.Gender,
			Phone:     t.Phone,
			Address:   t.Address,
			Subject:   t.Subject,
			Status:    t.Status,
			IsWalas:   t.IsWalas,
			CreatedAt: t.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if item.Status == "" {
			item.Status = "active"
		}

		// Kelas Walas
		if len(t.WalasClassIDs) > 0 {
			names := make([]string, 0, len(t.WalasClassIDs))
			for _, cid := range t.WalasClassIDs {
				if cname, ok := classMap[cid.Hex()]; ok {
					names = append(names, cname)
				}
			}
			item.WalasClasses = names
		}

		if !t.UserID.IsZero() {
			item.UserID = t.UserID.Hex()
			if u, ok := userMap[t.UserID.Hex()]; ok {
				item.Name = u.Name
				item.Username = u.Username
				item.Email = u.Email
				item.HasAccount = true
				item.IsActive = u.IsActive
			}
		}

		if item.Name == "" {
			item.Name = fmt.Sprintf("Guru (NIP: %s)", t.NIP)
		}

		items = append(items, item)
	}

	return items, total, nil
}

// GetTeacherDetail
func (s *TeacherManagementService) GetTeacherDetail(ctx context.Context, id bson.ObjectID) (*dto.TeacherItemResponse, error) {
	t, err := s.teacherRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("data guru tidak ditemukan")
	}

	resp := &dto.TeacherItemResponse{
		ID:        t.ID.Hex(),
		NIP:       t.NIP,
		Gender:    t.Gender,
		Phone:     t.Phone,
		Address:   t.Address,
		Subject:   t.Subject,
		Status:    t.Status,
		IsWalas:   t.IsWalas,
		CreatedAt: t.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	if resp.Status == "" {
		resp.Status = "active"
	}

	if len(t.WalasClassIDs) > 0 {
		cur, err := s.db.Collection("classes").Find(ctx, bson.M{"_id": bson.M{"$in": t.WalasClassIDs}})
		if err == nil {
			var classes []model.Class
			_ = cur.All(ctx, &classes)
			names := make([]string, 0, len(classes))
			for _, c := range classes {
				names = append(names, c.Name)
			}
			resp.WalasClasses = names
		}
	}

	if !t.UserID.IsZero() {
		resp.UserID = t.UserID.Hex()
		u, err := s.userRepo.FindByID(ctx, t.UserID)
		if err == nil && u != nil {
			resp.Name = u.Name
			resp.Username = u.Username
			resp.Email = u.Email
			resp.HasAccount = true
			resp.IsActive = u.IsActive
		}
	}

	if resp.Name == "" {
		resp.Name = fmt.Sprintf("Guru (NIP: %s)", t.NIP)
	}

	return resp, nil
}

// CreateTeacher
func (s *TeacherManagementService) CreateTeacher(ctx context.Context, req dto.CreateTeacherRequest) (*dto.TeacherItemResponse, error) {
	if ex, _ := s.teacherRepo.FindByNIP(ctx, req.NIP); ex != nil {
		return nil, fmt.Errorf("NIP %s sudah digunakan oleh guru lain", req.NIP)
	}

	status := req.Status
	if status == "" {
		status = "active"
	}

	var userID bson.ObjectID
	var hasAccount bool

	if req.CreateAccount {
		if req.Username == "" || req.Email == "" || req.Password == "" {
			return nil, errors.New("username, email, dan password wajib diisi jika membuat akun login")
		}
		if exU, _ := s.userRepo.FindByUsername(ctx, req.Username); exU != nil {
			return nil, fmt.Errorf("username %s sudah digunakan", req.Username)
		}
		if exE, _ := s.userRepo.FindByEmail(ctx, req.Email); exE != nil {
			return nil, fmt.Errorf("email %s sudah digunakan", req.Email)
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		u := &model.User{
			Name:         req.Name,
			Username:     req.Username,
			Email:        req.Email,
			PasswordHash: string(hash),
			Role:         model.RoleTeacher,
			IsActive:     true,
		}
		if err := s.userRepo.Create(ctx, u); err != nil {
			return nil, err
		}
		userID = u.ID
		hasAccount = true
	}

	walasClassIDs := make([]bson.ObjectID, 0)
	for _, cidStr := range req.WalasClassIDs {
		if cid, err := bson.ObjectIDFromHex(cidStr); err == nil {
			walasClassIDs = append(walasClassIDs, cid)
		}
	}

	t := &model.Teacher{
		UserID:        userID,
		NIP:           req.NIP,
		Gender:        req.Gender,
		Phone:         req.Phone,
		Address:       req.Address,
		Subject:       req.Subject,
		Status:        status,
		IsWalas:       req.IsWalas,
		WalasClassIDs: walasClassIDs,
	}

	if err := s.teacherRepo.Create(ctx, t); err != nil {
		return nil, err
	}

	return &dto.TeacherItemResponse{
		ID:         t.ID.Hex(),
		UserID:     userID.Hex(),
		Name:       req.Name,
		Username:   req.Username,
		Email:      req.Email,
		NIP:        t.NIP,
		Gender:     t.Gender,
		Phone:      t.Phone,
		Address:    t.Address,
		Subject:    t.Subject,
		Status:     t.Status,
		IsWalas:    t.IsWalas,
		HasAccount: hasAccount,
		IsActive:   true,
		CreatedAt:  t.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// UpdateTeacher
func (s *TeacherManagementService) UpdateTeacher(ctx context.Context, id bson.ObjectID, req dto.UpdateTeacherRequest) error {
	t, err := s.teacherRepo.FindByID(ctx, id)
	if err != nil {
		return errors.New("data guru tidak ditemukan")
	}

	if req.NIP != "" && req.NIP != t.NIP {
		if ex, _ := s.teacherRepo.FindByNIP(ctx, req.NIP); ex != nil && ex.ID != id {
			return fmt.Errorf("NIP %s sudah digunakan oleh guru lain", req.NIP)
		}
	}

	update := bson.M{
		"nip":     req.NIP,
		"gender":  req.Gender,
		"phone":   req.Phone,
		"address": req.Address,
		"subject": req.Subject,
		"status":  req.Status,
		"isWalas": req.IsWalas,
	}

	walasClassIDs := make([]bson.ObjectID, 0)
	for _, cidStr := range req.WalasClassIDs {
		if cid, err := bson.ObjectIDFromHex(cidStr); err == nil {
			walasClassIDs = append(walasClassIDs, cid)
		}
	}
	update["walasClassIds"] = walasClassIDs

	// Update nama di User jika ada
	if !t.UserID.IsZero() && req.Name != "" {
		_ = s.userRepo.Update(ctx, t.UserID, bson.M{"name": req.Name})
	}

	return s.teacherRepo.Update(ctx, id, update)
}

// DeleteTeacher
func (s *TeacherManagementService) DeleteTeacher(ctx context.Context, id bson.ObjectID) error {
	t, err := s.teacherRepo.FindByID(ctx, id)
	if err != nil {
		return errors.New("guru tidak ditemukan")
	}

	if !t.UserID.IsZero() {
		_ = s.userRepo.Delete(ctx, t.UserID)
	}

	return s.teacherRepo.Delete(ctx, id)
}

// CreateAccount
func (s *TeacherManagementService) CreateAccount(ctx context.Context, id bson.ObjectID, req dto.CreateTeacherAccountRequest) error {
	t, err := s.teacherRepo.FindByID(ctx, id)
	if err != nil {
		return errors.New("guru tidak ditemukan")
	}

	if !t.UserID.IsZero() {
		return errors.New("guru sudah memiliki akun aktif")
	}

	if exU, _ := s.userRepo.FindByUsername(ctx, req.Username); exU != nil {
		return errors.New("username sudah digunakan")
	}
	if exE, _ := s.userRepo.FindByEmail(ctx, req.Email); exE != nil {
		return errors.New("email sudah digunakan")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u := &model.User{
		Name:         fmt.Sprintf("Guru (NIP: %s)", t.NIP),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hash),
		Role:         model.RoleTeacher,
		IsActive:     true,
	}

	if err := s.userRepo.Create(ctx, u); err != nil {
		return err
	}

	return s.teacherRepo.Update(ctx, id, bson.M{"userId": u.ID})
}

// ToggleAccount
func (s *TeacherManagementService) ToggleAccount(ctx context.Context, id bson.ObjectID) (bool, error) {
	t, err := s.teacherRepo.FindByID(ctx, id)
	if err != nil || t.UserID.IsZero() {
		return false, errors.New("akun guru tidak ditemukan")
	}

	u, err := s.userRepo.FindByID(ctx, t.UserID)
	if err != nil {
		return false, err
	}

	newStatus := !u.IsActive
	err = s.userRepo.Update(ctx, t.UserID, bson.M{"isActive": newStatus})
	return newStatus, err
}

// GetStats
func (s *TeacherManagementService) GetStats(ctx context.Context) (*dto.TeacherStatsResponse, error) {
	total, _ := s.teacherRepo.Count(ctx, bson.M{})
	active, _ := s.teacherRepo.Count(ctx, bson.M{"$or": []bson.M{{"status": "active"}, {"status": bson.M{"$exists": false}}}})
	walas, _ := s.teacherRepo.Count(ctx, bson.M{"isWalas": true})
	noSubject, _ := s.teacherRepo.Count(ctx, bson.M{
		"$or": []bson.M{
			{"subject": ""},
			{"subject": bson.M{"$exists": false}},
		},
	})

	return &dto.TeacherStatsResponse{
		TotalTeachers:        total,
		ActiveTeachers:       active,
		WalasCount:           walas,
		TeachersWithoutMapel: noSubject,
	}, nil
}
