package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type AcademicYear struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string        `bson:"name" json:"name"`       // e.g. "2024/2025"
	Semester  int           `bson:"semester" json:"semester"` // 1 or 2
	StartDate time.Time     `bson:"startDate" json:"startDate"`
	EndDate   time.Time     `bson:"endDate" json:"endDate"`
	IsActive  bool          `bson:"isActive" json:"isActive"`
	CreatedAt time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time     `bson:"updatedAt" json:"updatedAt"`
}
