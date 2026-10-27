package service

import (
	"context"
	"errors"
	"time"

	"github.com/geniussociety/backend/internal/config"
	"github.com/geniussociety/backend/internal/dto"
	"github.com/geniussociety/backend/internal/model"
	"github.com/geniussociety/backend/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	UserID  string     `json:"sub"`
	Role    model.Role `json:"role"`
	IsWalas bool       `json:"isWalas"`
	jwt.RegisteredClaims
}

type AuthService struct {
	userRepo    *repository.UserRepository
	teacherRepo *repository.TeacherRepository
	cfg         *config.Config
}

func NewAuthService(userRepo *repository.UserRepository, teacherRepo *repository.TeacherRepository, cfg *config.Config) *AuthService {
	return &AuthService{userRepo: userRepo, teacherRepo: teacherRepo, cfg: cfg}
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*model.User, string, error) {
	var user *model.User
	var err error

	// Try email first, then username
	user, err = s.userRepo.FindByEmail(ctx, req.Login)
	if errors.Is(err, mongo.ErrNoDocuments) {
		user, err = s.userRepo.FindByUsername(ctx, req.Login)
	}
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	if !user.IsActive {
		return nil, "", errors.New("account is inactive")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	token, err := s.generateToken(ctx, user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) generateToken(ctx context.Context, user *model.User) (string, error) {
	isWalas := false
	if user.Role == model.RoleTeacher {
		teacher, err := s.teacherRepo.FindByUserID(ctx, user.ID)
		if err == nil {
			isWalas = teacher.IsWalas
		}
	}

	claims := Claims{
		UserID:  user.ID.Hex(),
		Role:    user.Role,
		IsWalas: isWalas,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

func (s *AuthService) ValidateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func (s *AuthService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid user id")
	}
	return s.userRepo.FindByID(ctx, oid)
}
