package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sandor-clegane/auth/internal/config"
	desc "github.com/sandor-clegane/auth/internal/generated/user_v1"
	"github.com/sandor-clegane/auth/pkg/utils/flags"
)

func main() {
	ctx := context.Background()

	flags, err := flags.ParseFlags()
	if err != nil {
		log.Fatalf("failed to parse flags: %v", err)
	}

	err = config.Load(flags.ConfigPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("failed to create config: %v", err)
	}

	lis, err := net.Listen("tcp", cfg.GRPCAddress())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	db, err := pgxpool.New(ctx, cfg.PGDSN())
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}

	err = db.Ping(ctx)
	if err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	desc.RegisterUserV1Server(s, &server{
		db: db,
	})

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	sig := <-sigs
	log.Printf("received signal %v, initiating graceful shutdown", sig)

	shutdownCtx, cancel := context.WithTimeout(ctx, cfg.GRPCShutdownTimeout())
	defer cancel()

	s.GracefulStop()
	db.Close()

	<-shutdownCtx.Done()
}
