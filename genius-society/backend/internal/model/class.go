package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Class struct {
	ID               bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name             string        `bson:"name" json:"name"`
	Grade            int           `bson:"grade" json:"grade"`
	AcademicYearID   bson.ObjectID `bson:"academicYearId" json:"academicYearId"`
	HomeroomTeacherID bson.ObjectID `bson:"homeroomTeacherId,omitempty" json:"homeroomTeacherId,omitempty"`
	CreatedAt        time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt        time.Time     `bson:"updatedAt" json:"updatedAt"`
}
