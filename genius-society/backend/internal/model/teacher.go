package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Teacher struct {
	ID             bson.ObjectID   `bson:"_id,omitempty" json:"id"`
	UserID         bson.ObjectID   `bson:"userId" json:"userId"`
	NIP            string          `bson:"nip" json:"nip"`
	IsWalas        bool            `bson:"isWalas" json:"isWalas"`
	WalasClassIDs  []bson.ObjectID `bson:"walasClassIds" json:"walasClassIds"`
	CreatedAt      time.Time       `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time       `bson:"updatedAt" json:"updatedAt"`
}
