package api

import (
	"context"
	"github.com/the-medo/talebound-backend/api/converters"
	"github.com/the-medo/talebound-backend/api/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *Server) GetWorldAdmins(ctx context.Context, req *pb.GetWorldAdminsRequest) (*pb.GetWorldAdminsResponse, error) {
	violations := validateGetWorldAdmins(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	worldAdminRows, err := server.Store.GetWorldAdmins(ctx, req.GetWorldId())
	if err != nil {
		return nil, err
	}

	rsp := &pb.GetWorldAdminsResponse{
		WorldAdmins: make([]*pb.WorldAdmin, len(worldAdminRows)),
	}

	for i, worldAdminRow := range worldAdminRows {
		rsp.WorldAdmins[i] = converters.ConvertWorldAdminRow(worldAdminRow)
	}

	return rsp, nil
}

func validateGetWorldAdmins(req *pb.GetWorldAdminsRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateWorldId(req.GetWorldId()); err != nil {
		violations = append(violations, e.FieldViolation("world_id", err))
	}

	return violations
}
