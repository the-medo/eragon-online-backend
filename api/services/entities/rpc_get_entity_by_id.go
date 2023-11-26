package entities

import (
	"context"
	"github.com/the-medo/talebound-backend/converters"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *ServiceEntities) GetEntityById(ctx context.Context, req *pb.GetEntityByIdRequest) (*pb.ViewEntity, error) {
	violations := validateGetEntityByIdRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	entity, err := server.Store.GetEntityByID(ctx, req.GetEntityId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get entity: %v", err)
	}

	return converters.ConvertViewEntity(entity), nil
}

func validateGetEntityByIdRequest(req *pb.GetEntityByIdRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateUniversalId(req.GetEntityId()); err != nil {
		violations = append(violations, e.FieldViolation("entity_id", err))
	}
	return violations
}
