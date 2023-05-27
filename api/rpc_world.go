package api

import (
	"context"
	"github.com/lib/pq"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateWorld(ctx context.Context, req *pb.CreateWorldRequest) (*pb.World, error) {
	authPayload, err := server.authorizeUserCookie(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}
	violations := validateCreateWorldRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	arg := db.CreateWorldTxParams{
		CreateWorldParams: db.CreateWorldParams{
			Name:        req.GetName(),
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

func validateCreateWorldRequest(req *pb.CreateWorldRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateString(req.GetName(), 3, 32); err != nil {
		violations = append(violations, FieldViolation("name", err))
	}

	if err := validator.ValidateString(req.GetDescription(), 1, 1024); err != nil {
		violations = append(violations, FieldViolation("description", err))
	}

	return violations
}
