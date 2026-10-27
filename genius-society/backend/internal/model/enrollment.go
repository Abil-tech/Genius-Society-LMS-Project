package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type EnrollmentStatus string

const (
	EnrollmentStatusActive      EnrollmentStatus = "active"
	EnrollmentStatusPromoted    EnrollmentStatus = "promoted"
	EnrollmentStatusRetained    EnrollmentStatus = "retained"
	EnrollmentStatusGraduated   EnrollmentStatus = "graduated"
	EnrollmentStatusTransferred EnrollmentStatus = "transferred"
)

type StudentEnrollment struct {
	ID             bson.ObjectID    `bson:"_id,omitempty" json:"id"`
	SchoolID       bson.ObjectID    `bson:"schoolId" json:"schoolId"`
	StudentID      bson.ObjectID    `bson:"studentId" json:"studentId"`
	ClassID        bson.ObjectID    `bson:"classId" json:"classId"`
	AcademicYearID bson.ObjectID    `bson:"academicYearId" json:"academicYearId"`
	Status         EnrollmentStatus `bson:"status" json:"status"`
	CreatedAt      time.Time        `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time        `bson:"updatedAt" json:"updatedAt"`
}
