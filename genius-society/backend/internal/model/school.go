package model

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

// School represents a school entity to which users belong.
type School struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string        `bson:"name" json:"name" binding:"required"`
	Address   string        `bson:"address,omitempty" json:"address,omitempty"`
	CreatedAt time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time     `bson:"updatedAt" json:"updatedAt"`
}
