package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Student struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    bson.ObjectID `bson:"userId" json:"userId"`
	NIS       string        `bson:"nis" json:"nis"`
	ClassID   bson.ObjectID `bson:"classId,omitempty" json:"classId,omitempty"`
	Status    string        `bson:"status" json:"status"` // active, inactive, graduated
	CreatedAt time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time     `bson:"updatedAt" json:"updatedAt"`
}
