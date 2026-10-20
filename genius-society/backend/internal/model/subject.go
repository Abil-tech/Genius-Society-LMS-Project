package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Subject struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string        `bson:"name" json:"name"`
	Code      string        `bson:"code" json:"code"`
	CreatedAt time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time     `bson:"updatedAt" json:"updatedAt"`
}
