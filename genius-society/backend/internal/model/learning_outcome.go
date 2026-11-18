package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// LearningOutcome represents a Pencapaian Pembelajaran (PP) defined by a teacher for a subject and classes.
type LearningOutcome struct {
	ID             bson.ObjectID   `bson:"_id,omitempty" json:"id"`
	SchoolID       bson.ObjectID   `bson:"schoolId" json:"schoolId"`
	AcademicYearID bson.ObjectID   `bson:"academicYearId" json:"academicYearId"`
	Semester       int             `bson:"semester" json:"semester"`
	SubjectID      bson.ObjectID   `bson:"subjectId" json:"subjectId"`
	ClassIDs       []bson.ObjectID `bson:"classIds" json:"classIds"`
	TeacherID      bson.ObjectID   `bson:"teacherId" json:"teacherId"`
	Title          string          `bson:"title" json:"title"`
	TitleKey       string          `bson:"titleKey" json:"titleKey"` // Normalized lower key for unique indexing
	Description    string          `bson:"description,omitempty" json:"description,omitempty"`
	Order          int             `bson:"order" json:"order"`
	CreatedAt      time.Time       `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time       `bson:"updatedAt" json:"updatedAt"`
}
