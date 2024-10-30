package config

import (
	"errors"
	"fmt"
	"os"
)

const (
	pgUserEnvName   = "PG_USER"
	pgDbNameEnvName = "PG_DATABASE_NAME"
	pgPassEnvName   = "PG_PASSWORD"
	pgPortEnvName   = "PG_PORT"
	pgHostEnvName   = "PG_HOST"
)

type pgConfig struct {
	user     string
	dbName   string
	password string
	port     string
	host     string
}

func newPGConfig() (*pgConfig, error) {
	user := os.Getenv(pgUserEnvName)
	if len(user) == 0 {
		return nil, errors.New("pg user not found")
	}

	dbName := os.Getenv(pgDbNameEnvName)
	if len(dbName) == 0 {
		return nil, errors.New("pg db name not found")
	}

	password := os.Getenv(pgPassEnvName)
	if len(password) == 0 {
		return nil, errors.New("pg password not found")
	}

	port := os.Getenv(pgPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("pg port not found")
	}

	host := os.Getenv(pgHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("pg host not found")
	}

	return &pgConfig{
		user:     user,
		dbName:   dbName,
		password: password,
		port:     port,
		host:     host,
	}, nil
}

func (cfg *pgConfig) PGDSN() string {
	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		cfg.host, cfg.port, cfg.dbName, cfg.user, cfg.password)
}
