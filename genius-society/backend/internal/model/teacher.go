package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Teacher struct {
	ID            bson.ObjectID   `bson:"_id,omitempty" json:"id"`
	UserID        bson.ObjectID   `bson:"userId" json:"userId"`
	NIP           string          `bson:"nip" json:"nip"`
	Gender        string          `bson:"gender,omitempty" json:"gender,omitempty"`
	Phone         string          `bson:"phone,omitempty" json:"phone,omitempty"`
	Address       string          `bson:"address,omitempty" json:"address,omitempty"`
	Subject       string          `bson:"subject,omitempty" json:"subject,omitempty"`
	Status        string          `bson:"status,omitempty" json:"status,omitempty"` // active, inactive, leave
	IsWalas       bool            `bson:"isWalas" json:"isWalas"`
	WalasClassIDs []bson.ObjectID `bson:"walasClassIds" json:"walasClassIds"`
	CreatedAt     time.Time       `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time       `bson:"updatedAt" json:"updatedAt"`
}

// userID, nip, gender, phone, address, subject, status, isWalas, walasClassIds
