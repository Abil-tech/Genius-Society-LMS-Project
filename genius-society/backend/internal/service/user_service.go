package service

import (
	"context"
	"errors"
	"time"

	"github.com/geniussociety/backend/internal/model"
	"github.com/geniussociety/backend/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo    *repository.UserRepository
	teacherRepo *repository.TeacherRepository
	studentRepo *repository.StudentRepository
}

func NewUserService(userRepo *repository.UserRepository, teacherRepo *repository.TeacherRepository, studentRepo *repository.StudentRepository) *UserService {
	return &UserService{userRepo: userRepo, teacherRepo: teacherRepo, studentRepo: studentRepo}
}

func (s *UserService) CreateUser(ctx context.Context, name, email, username, password string, role model.Role) (*model.User, error) {
	// check email uniqueness
	_, err := s.userRepo.FindByEmail(ctx, email)
	if err == nil {
		return nil, errors.New("email already exists")
	}
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:         name,
		Email:        email,
		Username:     username,
		PasswordHash: string(hash),
		Role:         role,
		IsActive:     true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// create role-specific profile
	switch role {
	case model.RoleTeacher:
		teacher := &model.Teacher{
			UserID:        user.ID,
			NIP:           "",
			IsWalas:       false,
			WalasClassIDs: []bson.ObjectID{},
		}
		_ = s.teacherRepo.Create(ctx, teacher)
	case model.RoleStudent:
		student := &model.Student{
			UserID: user.ID,
			NIS:    "",
			Status: "active",
		}
		_ = s.studentRepo.Create(ctx, student)
	}

	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id bson.ObjectID, name, email, username, password string, isActive *bool) error {
	update := bson.M{}
	if name != "" {
		update["name"] = name
	}
	if email != "" {
		update["email"] = email
	}
	if username != "" {
		update["username"] = username
	}
	if password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		update["passwordHash"] = string(hash)
	}
	if isActive != nil {
		update["isActive"] = *isActive
	}
	update["updatedAt"] = time.Now()
	return s.userRepo.Update(ctx, id, update)
}

func (s *UserService) GetUser(ctx context.Context, id bson.ObjectID) (*model.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

func (s *UserService) DeleteUser(ctx context.Context, id bson.ObjectID) error {
	return s.userRepo.Delete(ctx, id)
}

func (s *UserService) ListUsers(ctx context.Context, roleFilter string, page, limit int64) ([]*model.User, int64, error) {
	filter := bson.M{}
	if roleFilter != "" {
		filter["role"] = roleFilter
	}
	return s.userRepo.List(ctx, filter, page, limit)
}
