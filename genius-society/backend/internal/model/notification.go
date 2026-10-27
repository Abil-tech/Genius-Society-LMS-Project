package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Notification struct {
	ID        bson.ObjectID  `bson:"_id,omitempty" json:"id"`
	UserID    bson.ObjectID  `bson:"userId" json:"userId"`
	Type      string         `bson:"type" json:"type"` // new_assignment, new_material, etc.
	Title     string         `bson:"title" json:"title"`
	Body      string         `bson:"body" json:"body"`
	RefType   string         `bson:"refType,omitempty" json:"refType,omitempty"` // assignment, assessment, etc.
	RefID     *bson.ObjectID `bson:"refId,omitempty" json:"refId,omitempty"`
	IsRead    bool           `bson:"isRead" json:"isRead"`
	ReadAt    *time.Time     `bson:"readAt,omitempty" json:"readAt,omitempty"`
	CreatedAt time.Time      `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time      `bson:"updatedAt" json:"updatedAt"`
}
