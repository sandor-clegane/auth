package config

import (
	"errors"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	grpcHostEnvName            = "GRPC_HOST"
	grpcPortEnvName            = "GRPC_PORT"
	grpcShutdownTimeoutEnvName = "GRPC_SHUTDOWN_TIMEOUT"
)

type grpcConfig struct {
	host            string
	port            string
	shutdowmTimeout time.Duration
}

func newGRPCConfig() (*grpcConfig, error) {
	host := os.Getenv(grpcHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("grpc host not found")
	}

	port := os.Getenv(grpcPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("grpc port not found")
	}

	shutdownTimeoutStr := os.Getenv(grpcShutdownTimeoutEnvName)
	if len(shutdownTimeoutStr) == 0 {
		return nil, errors.New("grpc shutdown timeout not found")
	}

	shutdownTimeoutInt, err := strconv.ParseInt(shutdownTimeoutStr, 10, 64)
	if err != nil {
		return nil, err
	}

	return &grpcConfig{
		host:            host,
		port:            port,
		shutdowmTimeout: time.Second * time.Duration(shutdownTimeoutInt),
	}, nil
}

func (cfg *grpcConfig) GRPCAddress() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}

func (cfg *grpcConfig) GRPCShutdownTimeout() time.Duration {
	return cfg.shutdowmTimeout
}
