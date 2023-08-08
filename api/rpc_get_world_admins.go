package api

import (
	"context"
	"github.com/the-medo/talebound-backend/api/converters"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *Server) GetWorldAdmins(ctx context.Context, request *pb.GetWorldAdminsRequest) (*pb.GetWorldAdminsResponse, error) {
	violations := validateGetWorldAdmins(request)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	worldAdmins, err := server.store.GetWorldAdmins(ctx, request.GetWorldId())
	if err != nil {
		return nil, err
	}

	rsp := &pb.GetWorldAdminsResponse{
		WorldAdmins: make([]*pb.WorldAdmin, len(worldAdmins)),
	}

	for i, world := range worldAdmins {
		rsp.WorldAdmins[i] = converters.ConvertWorldAdmin(world)
	}

	return rsp, nil
}

func validateGetWorldAdmins(req *pb.GetWorldAdminsRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateWorldId(req.GetWorldId()); err != nil {
		violations = append(violations, FieldViolation("world_id", err))
	}

	return violations
}
