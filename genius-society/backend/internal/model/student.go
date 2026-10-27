package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Student struct {
	ID             bson.ObjectID `bson:"_id,omitempty" json:"id"`
	SchoolID       bson.ObjectID `bson:"schoolId,omitempty" json:"schoolId,omitempty"`
	UserID         bson.ObjectID `bson:"userId,omitempty" json:"userId,omitempty"`
	NIS            string        `bson:"nis" json:"nis"`
	NISN           string        `bson:"nisn,omitempty" json:"nisn,omitempty"`
	Gender         string        `bson:"gender,omitempty" json:"gender,omitempty"` // "L" / "P"
	BirthPlace     string        `bson:"birthPlace,omitempty" json:"birthPlace,omitempty"`
	BirthDate      string        `bson:"birthDate,omitempty" json:"birthDate,omitempty"`
	Address        string        `bson:"address,omitempty" json:"address,omitempty"`
	ParentName     string        `bson:"parentName,omitempty" json:"parentName,omitempty"`
	ParentPhone    string        `bson:"parentPhone,omitempty" json:"parentPhone,omitempty"`
	EnrollmentYear int           `bson:"enrollmentYear,omitempty" json:"enrollmentYear,omitempty"`
	PhotoURL       string        `bson:"photoUrl,omitempty" json:"photoUrl,omitempty"`
	ClassID        bson.ObjectID `bson:"classId,omitempty" json:"classId,omitempty"`
	Status         string        `bson:"status" json:"status"` // active, inactive, graduated, transferred, archived
	GraduationYear *int          `bson:"graduationYear,omitempty" json:"graduationYear,omitempty"`
	GraduatedAt    *time.Time    `bson:"graduatedAt,omitempty" json:"graduatedAt,omitempty"`
	CreatedAt      time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time     `bson:"updatedAt" json:"updatedAt"`
}

// ID, UserID, NIS/NISN, Gender, BirthPlace, BirthDate, Address, ParentName, ParentPhone, EnrollmentYear, PhotoURL, ClassID, Status
