package worlds

import (
	"context"
	"github.com/lib/pq"
	"github.com/the-medo/talebound-backend/converters"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *ServiceWorlds) CreateWorld(ctx context.Context, req *pb.CreateWorldRequest) (*pb.CreateWorldResponse, error) {
	violations := validateCreateWorldRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	authPayload, err := server.AuthorizeUserCookie(ctx)
	if err != nil {
		return nil, e.UnauthenticatedError(err)
	}

	arg := db.CreateWorldTxParams{
		CreateWorldParams: db.CreateWorldParams{
			Name:             req.GetName(),
			BasedOn:          req.GetBasedOn(),
			ShortDescription: req.GetShortDescription(),
		},
		UserId: authPayload.UserId,
	}

	txResult, err := server.Store.CreateWorldTx(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "name already exists: %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create world: %s", err)
	}

	rsp := &pb.CreateWorldResponse{
		World:  converters.ConvertWorld(*txResult.World),
		Module: converters.ConvertViewModule(*txResult.Module),
	}

	return rsp, nil
}

func validateCreateWorldRequest(req *pb.CreateWorldRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateModuleName(req.GetName()); err != nil {
		violations = append(violations, e.FieldViolation("name", err))
	}

	if err := validator.ValidateModuleShortDescription(req.GetShortDescription()); err != nil {
		violations = append(violations, e.FieldViolation("short_description", err))
	}

	if err := validator.ValidateModuleBasedOn(req.GetBasedOn()); err != nil {
		violations = append(violations, e.FieldViolation("based_on", err))
	}

	return violations
}
