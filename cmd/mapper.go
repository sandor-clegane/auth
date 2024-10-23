package main

import (
	"fmt"

	desc "github.com/sandor-clegane/auth/internal/generated/user_v1"
)

type Role string

const (
	Admin Role = "ADMIN"
	User  Role = "USER"
)

// TODO: separate layers
func roleToDB(role desc.Role) (Role, error) {
	switch role {
	case desc.Role_ROLE_ADMIN:
		return Admin, nil
	case desc.Role_ROLE_USER:
		return User, nil
	default:
		return "", fmt.Errorf("failed to convert role %v to db", role)
	}
}

func roleFromDB(role Role) desc.Role {
	switch role {
	case Admin:
		return desc.Role_ROLE_ADMIN
	case User:
		return desc.Role_ROLE_USER
	default:
		return desc.Role_ROLE_UNSPECIFIED
	}
}
