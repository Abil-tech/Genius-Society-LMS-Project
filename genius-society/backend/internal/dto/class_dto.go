package dto

import (
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type CreateClassReq struct {
	Code              string  `json:"code" validate:"required"`
	Name              string  `json:"name" validate:"required"`
	EducationLevel    string  `json:"educationLevel" validate:"required"`
	Department        *string `json:"department,omitempty"`
	AcademicYearID    string  `json:"academicYearId" validate:"required,objectid"`
	HomeroomTeacherID *string `json:"homeroomTeacherId,omitempty" validate:"omitempty,objectid"`
	Capacity          int     `json:"capacity" validate:"required,gt=0"`
	Status            string  `json:"status" validate:"required,oneof=active inactive archived"`
}

type UpdateClassReq struct {
	Name              *string `json:"name,omitempty"`
	EducationLevel    *string `json:"educationLevel,omitempty"`
	Department        *string `json:"department,omitempty"`
	AcademicYearID    *string `json:"academicYearId,omitempty" validate:"omitempty,objectid"`
	HomeroomTeacherID *string `json:"homeroomTeacherId,omitempty" validate:"omitempty,objectid"`
	Capacity          *int    `json:"capacity,omitempty" validate:"omitempty,gt=0"`
	Status            *string `json:"status,omitempty" validate:"omitempty,oneof=active inactive archived"`
}

type ClassResp struct {
	ID                 string  `json:"id"`
	Code               string  `json:"code"`
	Name               string  `json:"name"`
	EducationLevel     string  `json:"educationLevel"`
	Department         *string `json:"department,omitempty"`
	AcademicYearID     string  `json:"academicYearId"`
	HomeroomTeacherID  *string `json:"homeroomTeacherId,omitempty"`
	Capacity           int     `json:"capacity"`
	Status             string  `json:"status"`
	ActiveStudentCount int     `json:"activeStudentCount,omitempty"`
	CreatedAt          string  `json:"createdAt"`
	UpdatedAt          string  `json:"updatedAt"`
}

type ClassListResp struct {
	Classes []ClassResp `json:"classes"`
	Total   int64       `json:"total"`
	Page    int64       `json:"page"`
	Limit   int64       `json:"limit"`
}

func init() {
	v := validator.New()
	_ = v.RegisterValidation("objectid", func(fl validator.FieldLevel) bool {
		s := fl.Field().String()
		if s == "" {
			return true
		}
		_, err := bson.ObjectIDFromHex(s)
		return err == nil
	})
	// Register globally if needed elsewhere
}
