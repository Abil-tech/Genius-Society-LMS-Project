package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Role string

const (
	RoleSuperAdmin Role = "super_admin"
	RoleAdmin      Role = "admin"
	RoleTeacher    Role = "teacher"
	RoleStudent    Role = "student"
	RoleCurriculum Role = "curriculum"
	RolePrincipal  Role = "principal"
)

type User struct {
	ID           bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name         string        `bson:"name" json:"name"`
	Email        string        `bson:"email" json:"email"`
	Username     string        `bson:"username" json:"username"`
	PasswordHash string        `bson:"passwordHash" json:"-"`
	Role         Role          `bson:"role" json:"role"`
	IsActive     bool          `bson:"isActive" json:"isActive"`
	SchoolID     bson.ObjectID `bson:"schoolId,omitempty" json:"schoolId,omitempty"`
	CreatedAt    time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time     `bson:"updatedAt" json:"updatedAt"`
}
