package main

import (
	"fmt"

	desc "github.com/sandor-clegane/auth/internal/generated/user_v1"
)

// Role роль пользователя
type Role string

const (
	// Admin роль администратора
	Admin Role = "ADMIN"
	// User роль пользователя
	User Role = "USER"
)

// TODO: separate layers
func roleToDB(role desc.Role) (Role, error) {
	switch role {
	case desc.Role_ADMIN:
		return Admin, nil
	case desc.Role_USER:
		return User, nil
	default:
		return "", fmt.Errorf("failed to convert role %v to db", role)
	}
}

func roleFromDB(role Role) desc.Role {
	switch role {
	case Admin:
		return desc.Role_ADMIN
	case User:
		return desc.Role_USER
	default:
		return desc.Role_UNSPECIFIED
	}
}
