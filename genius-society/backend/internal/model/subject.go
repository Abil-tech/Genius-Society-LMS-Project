package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Subject struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Code        string        `bson:"code" json:"code"`                                 // contoh: MTK-X
	Name        string        `bson:"name" json:"name"`                                 // contoh: Matematika Wajib
	Category    string        `bson:"category,omitempty" json:"category,omitempty"`     // Umum, Peminatan IPA, Peminatan IPS, Muatan Lokal
	GradeLevel  string        `bson:"gradeLevel,omitempty" json:"gradeLevel,omitempty"` // X, XI, XII, Semua
	KKM         int           `bson:"kkm,omitempty" json:"kkm,omitempty"`               // Kriteria Ketuntasan Minimal, contoh 75
	Description string        `bson:"description,omitempty" json:"description,omitempty"`
	Status      string        `bson:"status,omitempty" json:"status,omitempty"` // active, inactive
	CreatedAt   time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time     `bson:"updatedAt" json:"updatedAt"`
}

// ID, Code, Name, Category, GradeLevel, KKM, Description, Status, CreatedAt, UpdatedAt
