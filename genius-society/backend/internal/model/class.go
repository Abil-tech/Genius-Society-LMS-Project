package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Class struct {
	ID                bson.ObjectID `bson:"_id,omitempty" json:"id"`
	SchoolID          bson.ObjectID `bson:"schoolId,omitempty" json:"schoolId,omitempty"`
	Code              string        `bson:"code" json:"code" validate:"required"`
	Name              string        `bson:"name" json:"name" validate:"required"`
	EducationLevel    string        `bson:"educationLevel" json:"educationLevel" validate:"required"`
	Department        *string       `bson:"department,omitempty" json:"department,omitempty"`
	AcademicYearID    bson.ObjectID `bson:"academicYearId" json:"academicYearId" validate:"required"`
	HomeroomTeacherID bson.ObjectID `bson:"homeroomTeacherId,omitempty" json:"homeroomTeacherId,omitempty"`
	Capacity          int           `bson:"capacity" json:"capacity" validate:"gt=0"`
	Status            string        `bson:"status" json:"status" validate:"oneof=active inactive archived"`
	CreatedAt         time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt         time.Time     `bson:"updatedAt" json:"updatedAt"`
}
