package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Assignment struct {
	ID            bson.ObjectID   `bson:"_id,omitempty" json:"id"`
	Title         string          `bson:"title" json:"title"`
	Description   string          `bson:"description" json:"description"`
	SubjectID     bson.ObjectID   `bson:"subjectId" json:"subjectId"`
	ClassIDs      []bson.ObjectID `bson:"classIds" json:"classIds"`
	Deadline      time.Time       `bson:"deadline" json:"deadline"`
	Attachments   []FileRef       `bson:"attachments" json:"attachments"`
	AllowResubmit bool            `bson:"allowResubmit" json:"allowResubmit"`
	TeacherID     bson.ObjectID   `bson:"teacherId" json:"teacherId"`
	CreatedAt     time.Time       `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time       `bson:"updatedAt" json:"updatedAt"`
}

type AssignmentSubmission struct {
	ID           bson.ObjectID `bson:"_id,omitempty" json:"id"`
	AssignmentID bson.ObjectID `bson:"assignmentId" json:"assignmentId"`
	StudentID    bson.ObjectID `bson:"studentId" json:"studentId"`
	Files        []FileRef     `bson:"files" json:"files"`
	Note         string        `bson:"note,omitempty" json:"note,omitempty"`
	SubmittedAt  time.Time     `bson:"submittedAt" json:"submittedAt"`
	Version      int           `bson:"version" json:"version"`
	Score        *float64      `bson:"score,omitempty" json:"score,omitempty"`
	Feedback     string        `bson:"feedback,omitempty" json:"feedback,omitempty"`
	GradedAt     *time.Time    `bson:"gradedAt,omitempty" json:"gradedAt,omitempty"`
	CreatedAt    time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time     `bson:"updatedAt" json:"updatedAt"`
}
