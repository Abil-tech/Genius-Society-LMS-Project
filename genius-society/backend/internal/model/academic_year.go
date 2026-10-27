package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type AcademicYear struct {
    ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
    SchoolID  bson.ObjectID `bson:"schoolId" json:"schoolId" binding:"required"`
    Name      string        `bson:"name" json:"name" binding:"required"` // e.g. "2026/2027"
    Semester  int           `bson:"semester" json:"semester" binding:"required,oneof=1 2"`
    StartDate time.Time     `bson:"startDate" json:"startDate" binding:"required"`
    EndDate   time.Time     `bson:"endDate" json:"endDate" binding:"required,gtfield=StartDate"`
    IsActive  bool          `bson:"isActive" json:"isActive"`
    Status    string        `bson:"status,omitempty" json:"status,omitempty"` // active, archived, inactive
    CreatedAt time.Time     `bson:"createdAt" json:"createdAt"`
    UpdatedAt time.Time     `bson:"updatedAt" json:"updatedAt"`
}

