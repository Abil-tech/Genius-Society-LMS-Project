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
	"golang.org/x/crypto/bcrypt"
)

type StudentManagementService struct {
	studentRepo *repository.StudentRepository
	userRepo    *repository.UserRepository
	classRepo   *repository.ClassRepository
	db          *mongo.Database
}

func NewStudentManagementService(
	studentRepo *repository.StudentRepository,
	userRepo *repository.UserRepository,
	classRepo *repository.ClassRepository,
	db *mongo.Database,
) *StudentManagementService {
	return &StudentManagementService{
		studentRepo: studentRepo,
		userRepo:    userRepo,
		classRepo:   classRepo,
		db:          db,
	}
}

// ListStudents with search, filters, sorting & pagination
func (s *StudentManagementService) ListStudents(
	ctx context.Context,
	search string,
	classIDStr string,
	status string,
	page int64,
	limit int64,
) ([]dto.StudentItemResponse, int64, error) {
	filter := bson.M{}

	if status != "" {
		filter["status"] = status
	}

	if classIDStr != "" {
		if classIDStr == "none" {
			filter["$or"] = []bson.M{
				{"classId": bson.M{"$exists": false}},
				{"classId": bson.NilObjectID},
			}
		} else if cid, err := bson.ObjectIDFromHex(classIDStr); err == nil {
			filter["classId"] = cid
		}
	}

	if search != "" {
		// Cari di NIS atau NISN langsung
		searchFilter := bson.M{
			"$or": []bson.M{
				{"nis": bson.M{"$regex": search, "$options": "i"}},
				{"nisn": bson.M{"$regex": search, "$options": "i"}},
			},
		}
		// Atau cari user dengan nama yang cocok
		userCur, err := s.db.Collection("users").Find(ctx, bson.M{
			"name": bson.M{"$regex": search, "$options": "i"},
			"role": model.RoleStudent,
		})
		if err == nil {
			var matchingUsers []model.User
			if err := userCur.All(ctx, &matchingUsers); err == nil && len(matchingUsers) > 0 {
				userIDs := make([]bson.ObjectID, len(matchingUsers))
				for i, u := range matchingUsers {
					userIDs[i] = u.ID
				}
				searchFilter["$or"] = append(searchFilter["$or"].([]bson.M), bson.M{"userId": bson.M{"$in": userIDs}})
			}
		}
		filter["$and"] = []bson.M{searchFilter}
	}

	students, total, err := s.studentRepo.List(ctx, filter, page, limit)
	if err != nil {
		return nil, 0, err
	}

	// Kumpulkan ID untuk populate
	userIDs := make([]bson.ObjectID, 0)
	classIDs := make([]bson.ObjectID, 0)
	for _, st := range students {
		if !st.UserID.IsZero() {
			userIDs = append(userIDs, st.UserID)
		}
		if !st.ClassID.IsZero() {
			classIDs = append(classIDs, st.ClassID)
		}
	}

	userMap := make(map[bson.ObjectID]*model.User)
	if len(userIDs) > 0 {
		cur, err := s.db.Collection("users").Find(ctx, bson.M{"_id": bson.M{"$in": userIDs}})
		if err == nil {
			var users []model.User
			if err := cur.All(ctx, &users); err == nil {
				for i := range users {
					userMap[users[i].ID] = &users[i]
				}
			}
		}
	}

	classMap := make(map[bson.ObjectID]string)
	if len(classIDs) > 0 {
		cur, err := s.db.Collection("classes").Find(ctx, bson.M{"_id": bson.M{"$in": classIDs}})
		if err == nil {
			var classes []model.Class
			if err := cur.All(ctx, &classes); err == nil {
				for _, c := range classes {
					classMap[c.ID] = c.Name
				}
			}
		}
	}

	res := make([]dto.StudentItemResponse, 0, len(students))
	for _, st := range students {
		item := dto.StudentItemResponse{
			ID:             st.ID.Hex(),
			NIS:            st.NIS,
			NISN:           st.NISN,
			Gender:         st.Gender,
			BirthPlace:     st.BirthPlace,
			BirthDate:      st.BirthDate,
			Address:        st.Address,
			ParentName:     st.ParentName,
			ParentPhone:    st.ParentPhone,
			EnrollmentYear: st.EnrollmentYear,
			PhotoURL:       st.PhotoURL,
			Status:         st.Status,
			CreatedAt:      st.CreatedAt.Format("02 Jan 2006"),
		}

		if !st.ClassID.IsZero() {
			item.ClassID = st.ClassID.Hex()
			item.ClassName = classMap[st.ClassID]
		}
		if item.ClassName == "" {
			item.ClassName = "Belum Ada Kelas"
		}

		if !st.UserID.IsZero() {
			item.UserID = st.UserID.Hex()
			if u, exists := userMap[st.UserID]; exists {
				item.Name = u.Name
				item.Username = u.Username
				item.Email = u.Email
				item.HasAccount = true
				item.IsActive = u.IsActive
			}
		}
		if item.Name == "" {
			item.Name = fmt.Sprintf("Siswa %s", st.NIS)
		}

		res = append(res, item)
	}

	return res, total, nil
}

// GetStudentDetail by ID
func (s *StudentManagementService) GetStudentDetail(ctx context.Context, id bson.ObjectID) (*dto.StudentItemResponse, error) {
	st, err := s.studentRepo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("data siswa tidak ditemukan")
	}

	item := &dto.StudentItemResponse{
		ID:             st.ID.Hex(),
		NIS:            st.NIS,
		NISN:           st.NISN,
		Gender:         st.Gender,
		BirthPlace:     st.BirthPlace,
		BirthDate:      st.BirthDate,
		Address:        st.Address,
		ParentName:     st.ParentName,
		ParentPhone:    st.ParentPhone,
		EnrollmentYear: st.EnrollmentYear,
		PhotoURL:       st.PhotoURL,
		Status:         st.Status,
		CreatedAt:      st.CreatedAt.Format("02 Jan 2006"),
	}

	if !st.ClassID.IsZero() {
		item.ClassID = st.ClassID.Hex()
		c, err := s.classRepo.FindByID(ctx, st.ClassID)
		if err == nil && c != nil {
			item.ClassName = c.Name
		}
	}
	if item.ClassName == "" {
		item.ClassName = "Belum Ada Kelas"
	}

	if !st.UserID.IsZero() {
		item.UserID = st.UserID.Hex()
		u, err := s.userRepo.FindByID(ctx, st.UserID)
		if err == nil && u != nil {
			item.Name = u.Name
			item.Username = u.Username
			item.Email = u.Email
			item.HasAccount = true
			item.IsActive = u.IsActive
		}
	}

	return item, nil
}

// CreateStudent with optional user login creation
func (s *StudentManagementService) CreateStudent(ctx context.Context, req dto.CreateStudentRequest) (*dto.StudentItemResponse, error) {
	// 1. Validasi duplikat NIS
	existingNIS, _ := s.studentRepo.FindByNIS(ctx, req.NIS)
	if existingNIS != nil {
		return nil, errors.New("NIS sudah terdaftar di sistem")
	}

	// 2. Validasi duplikat NISN
	if req.NISN != "" {
		existingNISN, _ := s.studentRepo.FindByNISN(ctx, req.NISN)
		if existingNISN != nil {
			return nil, errors.New("NISN sudah terdaftar di sistem")
		}
	}

	var userID bson.ObjectID
	// 3. Buat akun jika dicentang
	if req.CreateAccount && req.Username != "" && req.Password != "" {
		// Cek username
		if exU, _ := s.userRepo.FindByUsername(ctx, req.Username); exU != nil {
			return nil, errors.New("username sudah digunakan")
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		email := req.Email
		if email == "" {
			email = fmt.Sprintf("%s@siswa.genius.ac.id", req.Username)
		}

		user := &model.User{
			Name:         req.Name,
			Username:     req.Username,
			Email:        email,
			PasswordHash: string(hash),
			Role:         model.RoleStudent,
			IsActive:     true,
		}
		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, fmt.Errorf("gagal membuat akun murid: %v", err)
		}
		userID = user.ID
	}

	var classID bson.ObjectID
	if req.ClassID != "" && req.ClassID != "none" {
		if cid, err := bson.ObjectIDFromHex(req.ClassID); err == nil {
			classID = cid
		}
	}

	status := req.Status
	if status == "" {
		status = "active"
	}

	enrollYear := req.EnrollmentYear
	if enrollYear <= 0 {
		enrollYear = time.Now().Year()
	}

	st := &model.Student{
		UserID:         userID,
		NIS:            req.NIS,
		NISN:           req.NISN,
		Gender:         req.Gender,
		BirthPlace:     req.BirthPlace,
		BirthDate:      req.BirthDate,
		Address:        req.Address,
		ParentName:     req.ParentName,
		ParentPhone:    req.ParentPhone,
		EnrollmentYear: enrollYear,
		ClassID:        classID,
		Status:         status,
	}

	if err := s.studentRepo.Create(ctx, st); err != nil {
		return nil, err
	}

	return s.GetStudentDetail(ctx, st.ID)
}

// UpdateStudent
func (s *StudentManagementService) UpdateStudent(ctx context.Context, id bson.ObjectID, req dto.UpdateStudentRequest) error {
	st, err := s.studentRepo.FindByID(ctx, id)
	if err != nil {
		return errors.New("siswa tidak ditemukan")
	}

	// Validasi NIS duplikat
	if req.NIS != "" && req.NIS != st.NIS {
		if ex, _ := s.studentRepo.FindByNIS(ctx, req.NIS); ex != nil {
			return errors.New("NIS baru sudah terdaftar")
		}
	}

	update := bson.M{}
	if req.NIS != "" {
		update["nis"] = req.NIS
	}
	if req.NISN != "" {
		update["nisn"] = req.NISN
	}
	if req.Gender != "" {
		update["gender"] = req.Gender
	}
	if req.BirthPlace != "" {
		update["birthPlace"] = req.BirthPlace
	}
	if req.BirthDate != "" {
		update["birthDate"] = req.BirthDate
	}
	if req.Address != "" {
		update["address"] = req.Address
	}
	if req.ParentName != "" {
		update["parentName"] = req.ParentName
	}
	if req.ParentPhone != "" {
		update["parentPhone"] = req.ParentPhone
	}
	if req.EnrollmentYear > 0 {
		update["enrollmentYear"] = req.EnrollmentYear
	}
	if req.Status != "" {
		update["status"] = req.Status
	}

	if req.ClassID == "none" {
		update["classId"] = bson.NilObjectID
	} else if req.ClassID != "" {
		if cid, err := bson.ObjectIDFromHex(req.ClassID); err == nil {
			update["classId"] = cid
		}
	}

	if err := s.studentRepo.Update(ctx, id, update); err != nil {
		return err
	}

	// Update user name jika ada userID
	if !st.UserID.IsZero() && req.Name != "" {
		_ = s.userRepo.Update(ctx, st.UserID, bson.M{"name": req.Name})
	}

	return nil
}

// BulkAssignClass
func (s *StudentManagementService) BulkAssignClass(ctx context.Context, studentIDs []string, classIDStr string) error {
	var classID bson.ObjectID
	if classIDStr != "none" {
		cid, err := bson.ObjectIDFromHex(classIDStr)
		if err != nil {
			return errors.New("invalid classId")
		}
		classID = cid
	}

	objIDs := make([]bson.ObjectID, 0, len(studentIDs))
	for _, idStr := range studentIDs {
		if oid, err := bson.ObjectIDFromHex(idStr); err == nil {
			objIDs = append(objIDs, oid)
		}
	}

	return s.studentRepo.BulkUpdateClass(ctx, objIDs, classID)
}

// BulkUpdateStatus
func (s *StudentManagementService) BulkUpdateStatus(ctx context.Context, studentIDs []string, status string) error {
	objIDs := make([]bson.ObjectID, 0, len(studentIDs))
	for _, idStr := range studentIDs {
		if oid, err := bson.ObjectIDFromHex(idStr); err == nil {
			objIDs = append(objIDs, oid)
		}
	}
	return s.studentRepo.BulkUpdateStatus(ctx, objIDs, status)
}

// CreateOrLinkAccount
func (s *StudentManagementService) CreateAccount(ctx context.Context, studentID bson.ObjectID, req dto.CreateStudentAccountRequest) error {
	st, err := s.studentRepo.FindByID(ctx, studentID)
	if err != nil {
		return errors.New("siswa tidak ditemukan")
	}

	if !st.UserID.IsZero() {
		return errors.New("siswa sudah memiliki akun aktif")
	}

	if exU, _ := s.userRepo.FindByUsername(ctx, req.Username); exU != nil {
		return errors.New("username sudah digunakan")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u := &model.User{
		Name:         fmt.Sprintf("Siswa %s", st.NIS),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hash),
		Role:         model.RoleStudent,
		IsActive:     true,
	}

	if err := s.userRepo.Create(ctx, u); err != nil {
		return err
	}

	return s.studentRepo.Update(ctx, studentID, bson.M{"userId": u.ID})
}

// DeleteStudent menghapus data siswa dan akun user terkait jika ada
func (s *StudentManagementService) DeleteStudent(ctx context.Context, studentID bson.ObjectID) error {
	st, err := s.studentRepo.FindByID(ctx, studentID)
	if err != nil {
		return errors.New("siswa tidak ditemukan")
	}

	// Jika siswa memiliki akun user yang terhubung, hapus juga akunnya
	if !st.UserID.IsZero() {
		_ = s.userRepo.Delete(ctx, st.UserID)
	}

	return s.studentRepo.Delete(ctx, studentID)
}

// ToggleAccountStatus
func (s *StudentManagementService) ToggleAccountStatus(ctx context.Context, studentID bson.ObjectID) (bool, error) {
	st, err := s.studentRepo.FindByID(ctx, studentID)
	if err != nil || st.UserID.IsZero() {
		return false, errors.New("akun murid tidak ditemukan")
	}

	u, err := s.userRepo.FindByID(ctx, st.UserID)
	if err != nil {
		return false, err
	}

	newStatus := !u.IsActive
	err = s.userRepo.Update(ctx, st.UserID, bson.M{"isActive": newStatus})
	return newStatus, err
}

// GetStats
func (s *StudentManagementService) GetStats(ctx context.Context) (*dto.StudentStatsResponse, error) {
	total, _ := s.studentRepo.Count(ctx, bson.M{})
	active, _ := s.studentRepo.Count(ctx, bson.M{"status": "active"})
	noClass, _ := s.studentRepo.Count(ctx, bson.M{
		"$or": []bson.M{
			{"classId": bson.M{"$exists": false}},
			{"classId": bson.NilObjectID},
		},
	})
	noAccount, _ := s.studentRepo.Count(ctx, bson.M{
		"$or": []bson.M{
			{"userId": bson.M{"$exists": false}},
			{"userId": bson.NilObjectID},
		},
	})

	return &dto.StudentStatsResponse{
		TotalStudents:          total,
		ActiveStudents:         active,
		StudentsWithoutClass:   noClass,
		StudentsWithoutAccount: noAccount,
	}, nil
}
