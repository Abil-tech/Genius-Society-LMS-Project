package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type GradeType string

const (
	GradeAssignment GradeType = "assignment"
	GradeAssessment GradeType = "assessment"
	GradeProject    GradeType = "project"
	GradeUTS        GradeType = "uts"
	GradeUAS        GradeType = "uas"
)

type Grade struct {
	ID             bson.ObjectID `bson:"_id,omitempty" json:"id"`
	StudentID      bson.ObjectID `bson:"studentId" json:"studentId"`
	SubjectID      bson.ObjectID `bson:"subjectId" json:"subjectId"`
	ClassID        bson.ObjectID `bson:"classId" json:"classId"`
	AcademicYearID bson.ObjectID `bson:"academicYearId" json:"academicYearId"`
	Semester       int           `bson:"semester" json:"semester"`
	TeacherID      bson.ObjectID `bson:"teacherId" json:"teacherId"`
	Type           GradeType     `bson:"type" json:"type"`
	SourceID       bson.ObjectID `bson:"sourceId" json:"sourceId"` // assignmentId, assessmentId, or projectId
	Score          float64       `bson:"score" json:"score"`
	CreatedAt      time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time     `bson:"updatedAt" json:"updatedAt"`
}

type GradingConfig struct {
	ID             bson.ObjectID  `bson:"_id,omitempty" json:"id"`
	AcademicYearID bson.ObjectID  `bson:"academicYearId" json:"academicYearId"`
	SubjectID      *bson.ObjectID `bson:"subjectId,omitempty" json:"subjectId,omitempty"` // nil = global default
	DailyWeight    float64        `bson:"dailyWeight" json:"dailyWeight"`                 // default 60
	UTSWeight      float64        `bson:"utsWeight" json:"utsWeight"`                     // default 20
	UASWeight      float64        `bson:"uasWeight" json:"uasWeight"`                     // default 20
	CreatedAt      time.Time      `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time      `bson:"updatedAt" json:"updatedAt"`
}
