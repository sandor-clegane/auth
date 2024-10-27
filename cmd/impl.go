package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	desc "github.com/sandor-clegane/auth/internal/generated/user_v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	db *pgxpool.Pool

	desc.UnimplementedUserV1Server
}

// Create ...
func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	dbRole, err := roleToDB(req.GetInfo().GetRole())
	if err != nil {
		slog.Error("failed to convert role to db", slog.Any("error", err))
		return nil, err
	}

	insertBuilder := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("name", "email", "role").
		Values(
			req.GetInfo().GetName(),
			req.GetInfo().GetEmail(),
			dbRole,
		).
		Suffix("RETURNING id")

	query, args, err := insertBuilder.ToSql()
	if err != nil {
		slog.Error("failed to build insert query", slog.Any("error", err))
		return nil, err
	}

	var userID int64
	err = s.db.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		slog.Error("failed to insert user", slog.Any("error", err))
		return nil, err
	}

	slog.Info("user inserted", slog.Any("user_id", userID))

	return &desc.CreateResponse{
		Id: userID,
	}, nil
}

// Get ...
func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	var (
		userID    int64
		name      string
		email     string
		role      Role
		createdAt time.Time
		updatedAt sql.NullTime
	)

	selectBuilder := sq.Select("id", "name", "email",
		"role", "created_at", "updated_at").
		PlaceholderFormat(sq.Dollar).
		From("users").
		Where(sq.Eq{"id": req.GetId()})

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		slog.Error("failed to build select query", slog.Any("error", err))
		return nil, err
	}

	err = s.db.QueryRow(ctx, query, args...).
		Scan(
			&userID, &name, &email,
			&role, &createdAt, &updatedAt,
		)
	if err != nil {
		slog.Error("failed to get user", slog.Any("error", err))
		return nil, err
	}

	slog.Info("get user", slog.Any("user_id", req.GetId()))

	return &desc.GetResponse{
		User: &desc.User{
			Id: userID,
			Info: &desc.UserInfo{
				Name:  name,
				Email: email,
				Role:  roleFromDB(role),
			},
			CreatedAt: timestamppb.New(createdAt),
			UpdatedAt: func(t sql.NullTime) *timestamppb.Timestamp {
				if t.Valid {
					return timestamppb.New(t.Time)
				}

				return nil
			}(updatedAt),
		},
	}, nil
}

func buildUpdatesMap(req *desc.UpdateRequest) (map[string]interface{}, bool) {
	updates := make(map[string]interface{})

	if email := req.GetInfo().GetEmail().GetValue(); email != "" {
		updates["email"] = email
	}

	if name := req.GetInfo().GetName().GetValue(); name != "" {
		updates["name"] = name
	}

	if role, err := roleToDB(req.GetInfo().GetRole()); err == nil {
		updates["role"] = role
	}

	if len(updates) != 0 {
		updates["updated_at"] = time.Now()
		return updates, false
	}

	return updates, true
}

// Update ...
func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	updatedMap, noUpdates := buildUpdatesMap(req)
	if noUpdates {
		return &emptypb.Empty{}, nil
	}

	updateBuilder := sq.Update("users").
		SetMap(updatedMap).
		Where(sq.Eq{"id": req.GetId()}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		slog.Error("failed to build update query", slog.Any("error", err))
		return nil, err
	}

	fmt.Println(query)

	_, err = s.db.Exec(ctx, query, args...)
	if err != nil {
		slog.Error("failed to update user", slog.Any("error", err))
		return nil, err
	}

	slog.Info("user updated", slog.Any("user_id", req.GetId()))

	return &emptypb.Empty{}, nil
}

// Delete ...
func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	deleteBuilder := sq.Delete("users").
		Where(sq.Eq{"id": req.GetId()}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := deleteBuilder.ToSql()
	if err != nil {
		slog.Error("failed to build delete query", slog.Any("error", err))
		return nil, err
	}

	_, err = s.db.Exec(ctx, query, args...)
	if err != nil {
		slog.Error("failed to delete user", slog.Any("error", err))
		return nil, err
	}

	slog.Info("user deleted", slog.Any("user_id", req.GetId()))

	return &emptypb.Empty{}, nil
}
