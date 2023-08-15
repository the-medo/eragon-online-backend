package api

import (
	"context"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *Server) GetWorldById(ctx context.Context, req *pb.GetWorldByIdRequest) (*pb.World, error) {
	violations := validateGetWorldById(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	world, err := server.store.GetWorldByID(ctx, req.WorldId)
	if err != nil {
		return nil, err
	}

	rsp := convertWorld(world)

	return rsp, nil
}

func validateGetWorldById(req *pb.GetWorldByIdRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateWorldId(req.GetWorldId()); err != nil {
		violations = append(violations, FieldViolation("world_id", err))
	}

	return violations
}
