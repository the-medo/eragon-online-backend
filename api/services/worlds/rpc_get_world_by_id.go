package worlds

import (
	"context"
	"github.com/the-medo/talebound-backend/converters"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *ServiceWorlds) GetWorldById(ctx context.Context, req *pb.GetWorldByIdRequest) (*pb.World, error) {
	violations := validateGetWorldById(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	world, err := server.Store.GetWorldByID(ctx, req.WorldId)
	if err != nil {
		return nil, err
	}

	return converters.ConvertWorld(world), nil
}

func validateGetWorldById(req *pb.GetWorldByIdRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateModuleId(req.GetWorldId()); err != nil {
		violations = append(violations, e.FieldViolation("world_id", err))
	}

	return violations
}
