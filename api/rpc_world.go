package api

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateWorld(ctx context.Context, req *pb.CreateWorldRequest) (*pb.World, error) {
	violations := validateCreateWorldRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	authPayload, err := server.authorizeUserCookie(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	arg := db.CreateWorldTxParams{
		CreateWorldParams: db.CreateWorldParams{
			Name:        req.GetName(),
			BasedOn:     req.GetBasedOn(),
			Description: req.GetDescription(),
		},
		UserId: authPayload.UserId,
	}

	txResult, err := server.store.CreateWorldTx(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "name already exists: %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create world: %s", err)
	}

	rsp := convertWorld(txResult)

	return rsp, nil
}

func (server *Server) UpdateWorld(ctx context.Context, req *pb.UpdateWorldRequest) (*pb.World, error) {
	violations := validateUpdateWorldRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	err := server.CheckWorldAdmin(ctx, req.GetWorldId(), true)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "failed to update world: %v", err)
	}

	arg := db.UpdateWorldParams{
		WorldID: req.GetWorldId(),
		Name: sql.NullString{
			String: req.GetName(),
			Valid:  req.Name != nil,
		},
		Description: sql.NullString{
			String: req.GetDescription(),
			Valid:  req.Description != nil,
		},
		Public: sql.NullBool{
			Bool:  req.GetPublic(),
			Valid: req.Public != nil,
		},
	}

	_, err = server.store.UpdateWorld(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update world: %v", err)
	}

	world, err := server.store.GetWorldByID(ctx, req.GetWorldId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve updated world: %v", err)
	}

	rsp := convertWorld(world)

	return rsp, nil
}

func validateCreateWorldRequest(req *pb.CreateWorldRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateString(req.GetName(), 3, 32); err != nil {
		violations = append(violations, FieldViolation("name", err))
	}

	if err := validator.ValidateString(req.GetDescription(), 1, 1024); err != nil {
		violations = append(violations, FieldViolation("description", err))
	}

	return violations
}

func validateUpdateWorldRequest(req *pb.UpdateWorldRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if req.Name != nil {
		if err := validator.ValidateString(req.GetName(), 3, 32); err != nil {
			violations = append(violations, FieldViolation("name", err))
		}
	}

	if req.Description != nil {
		if err := validator.ValidateString(req.GetDescription(), 1, 1024); err != nil {
			violations = append(violations, FieldViolation("description", err))
		}
	}

	return violations
}
