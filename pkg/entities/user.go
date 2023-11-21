package entities

import "go-api/pkg/auth"

type User struct {
	Id       string      `json:"id,omitempty" bson:"_id,omitempty"`
	Username string      `json:"username"`
	Password string      `json:"password"`
	Roles    []auth.Role `json:"roles,omitempty"`
}

type Role string

const (
	ADMIN  Role = "ADMIN"
	EDITOR Role = "EDITOR"
	VIEWER Role = "VIEWER"
)
