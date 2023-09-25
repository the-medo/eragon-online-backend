package api

import (
	"context"
	"github.com/the-medo/talebound-backend/api/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *Server) GetWorldsOfCreator(ctx context.Context, req *pb.GetWorldsOfCreatorRequest) (*pb.GetWorldsOfCreatorResponse, error) {
	violations := validateGetWorldsOfCreator(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	worldsWithAdminInfo, err := server.Store.GetWorldsOfUser(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	rsp := &pb.GetWorldsOfCreatorResponse{
		Worlds: make([]*pb.WorldOfCreatorResponse, len(worldsWithAdminInfo)),
	}

	for i, world := range worldsWithAdminInfo {
		rsp.Worlds[i] = &pb.WorldOfCreatorResponse{
			World:      convertWorldOfUser(world),
			SuperAdmin: world.WorldSuperAdmin,
		}
	}

	return rsp, nil
}

func validateGetWorldsOfCreator(req *pb.GetWorldsOfCreatorRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateUserId(req.GetUserId()); err != nil {
		violations = append(violations, e.FieldViolation("user_id", err))
	}

	return violations
}
