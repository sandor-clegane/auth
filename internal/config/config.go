package config

import "github.com/joho/godotenv"

// Load инициализирует переменные окружения значениями из переданного файла
func Load(path string) error {
	return godotenv.Load(path)
}

// New конструктор конфига
func New() (*Config, error) {
	grpcConfig, err := newGRPCConfig()
	if err != nil {
		return nil, err
	}

	pgConfig, err := newPGConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		grpcConfig: grpcConfig,
		pgConfig:   pgConfig,
	}, nil
}

// Config конфиг сервиса
type Config struct {
	*grpcConfig
	*pgConfig
}
