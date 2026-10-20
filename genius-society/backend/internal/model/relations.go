package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ClassStudent struct {
	ID             bson.ObjectID `bson:"_id,omitempty" json:"id"`
	ClassID        bson.ObjectID `bson:"classId" json:"classId"`
	StudentID      bson.ObjectID `bson:"studentId" json:"studentId"`
	AcademicYearID bson.ObjectID `bson:"academicYearId" json:"academicYearId"`
	CreatedAt      time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time     `bson:"updatedAt" json:"updatedAt"`
}

type TeacherSubject struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	TeacherID bson.ObjectID `bson:"teacherId" json:"teacherId"`
	SubjectID bson.ObjectID `bson:"subjectId" json:"subjectId"`
	CreatedAt time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time     `bson:"updatedAt" json:"updatedAt"`
}

type TeacherClass struct {
	ID             bson.ObjectID `bson:"_id,omitempty" json:"id"`
	TeacherID      bson.ObjectID `bson:"teacherId" json:"teacherId"`
	ClassID        bson.ObjectID `bson:"classId" json:"classId"`
	SubjectID      bson.ObjectID `bson:"subjectId" json:"subjectId"`
	AcademicYearID bson.ObjectID `bson:"academicYearId" json:"academicYearId"`
	CreatedAt      time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time     `bson:"updatedAt" json:"updatedAt"`
}
