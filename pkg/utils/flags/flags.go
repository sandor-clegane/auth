package flags

import (
	"errors"
	"flag"
)

const (
	configPathFlag = "config-path"
)

// CMDFlags флаговые аргументы приложения
type CMDFlags struct {
	ConfigPath string
}

// ParseFlags парсит флаги
func ParseFlags() (*CMDFlags, error) {
	configPath := flag.String(configPathFlag, ".env", "Configuration file path")

	flag.Parse()

	if *configPath == "" {
		return nil, errors.New("Configuration file path was not found in application flags")
	}

	return &CMDFlags{ConfigPath: *configPath}, nil
}
