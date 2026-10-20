package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Schedule struct {
	ID             bson.ObjectID `bson:"_id,omitempty" json:"id"`
	ClassID        bson.ObjectID `bson:"classId" json:"classId"`
	SubjectID      bson.ObjectID `bson:"subjectId" json:"subjectId"`
	TeacherID      bson.ObjectID `bson:"teacherId" json:"teacherId"`
	AcademicYearID bson.ObjectID `bson:"academicYearId" json:"academicYearId"`
	Day            int           `bson:"day" json:"day"`       // 1=Monday ... 7=Sunday
	StartTime      string        `bson:"startTime" json:"startTime"` // "08:00"
	EndTime        string        `bson:"endTime" json:"endTime"`     // "09:30"
	Room           string        `bson:"room" json:"room"`
	CreatedAt      time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time     `bson:"updatedAt" json:"updatedAt"`
}
