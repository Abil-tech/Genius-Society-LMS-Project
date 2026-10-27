package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type AuditLog struct {
	ID         bson.ObjectID `bson:"_id,omitempty" json:"id"`
	SchoolID   bson.ObjectID `bson:"schoolId" json:"schoolId"`
	UserID     bson.ObjectID `bson:"userId" json:"userId"`
	Action     string        `bson:"action" json:"action"`
	EntityType string        `bson:"entityType" json:"entityType"`
	EntityID   string        `bson:"entityId" json:"entityId"`
	Details    bson.M        `bson:"details,omitempty" json:"details,omitempty"`
	Timestamp  time.Time     `bson:"timestamp" json:"timestamp"`
}
