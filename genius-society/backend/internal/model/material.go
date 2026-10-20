package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// FileRef stores a reference to a file in Backblaze B2
type FileRef struct {
	Key  string `bson:"key" json:"key"`
	Name string `bson:"name" json:"name"`
	Size int64  `bson:"size" json:"size"`
	Mime string `bson:"mime" json:"mime"`
}

type Material struct {
	ID          bson.ObjectID   `bson:"_id,omitempty" json:"id"`
	Title       string          `bson:"title" json:"title"`
	Description string          `bson:"description" json:"description"`
	SubjectID   bson.ObjectID   `bson:"subjectId" json:"subjectId"`
	ClassIDs    []bson.ObjectID `bson:"classIds" json:"classIds"`
	File        *FileRef        `bson:"file,omitempty" json:"file,omitempty"`
	Link        string          `bson:"link,omitempty" json:"link,omitempty"`
	PublishAt   *time.Time      `bson:"publishAt,omitempty" json:"publishAt,omitempty"`
	TeacherID   bson.ObjectID   `bson:"teacherId" json:"teacherId"`
	CreatedAt   time.Time       `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time       `bson:"updatedAt" json:"updatedAt"`
}
