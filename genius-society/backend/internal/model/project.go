package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Project struct {
	ID            bson.ObjectID   `bson:"_id,omitempty" json:"id"`
	Title         string          `bson:"title" json:"title"`
	Description   string          `bson:"description" json:"description"`
	Instructions  string          `bson:"instructions" json:"instructions"`
	Deadline      time.Time       `bson:"deadline" json:"deadline"`
	SubjectID     bson.ObjectID   `bson:"subjectId" json:"subjectId"`
	ClassIDs      []bson.ObjectID `bson:"classIds" json:"classIds"`
	Attachments   []FileRef       `bson:"attachments" json:"attachments"`
	AllowResubmit bool            `bson:"allowResubmit" json:"allowResubmit"`
	TeacherID     bson.ObjectID   `bson:"teacherId" json:"teacherId"`
	CreatedAt     time.Time       `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time       `bson:"updatedAt" json:"updatedAt"`
}

type ProjectSubmission struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	ProjectID   bson.ObjectID `bson:"projectId" json:"projectId"`
	StudentID   bson.ObjectID `bson:"studentId" json:"studentId"`
	Files       []FileRef     `bson:"files" json:"files"`
	SubmittedAt time.Time     `bson:"submittedAt" json:"submittedAt"`
	Version     int           `bson:"version" json:"version"`
	Score       *float64      `bson:"score,omitempty" json:"score,omitempty"`
	Feedback    string        `bson:"feedback,omitempty" json:"feedback,omitempty"`
	CreatedAt   time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time     `bson:"updatedAt" json:"updatedAt"`
}
