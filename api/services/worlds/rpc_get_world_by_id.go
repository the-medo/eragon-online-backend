package worlds

import (
	"context"
	"encoding/json"
	"github.com/the-medo/talebound-backend/api/apihelpers"
	"github.com/the-medo/talebound-backend/converters"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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

	fetchInterface := &apihelpers.FetchInterface{
		PostIds:  []int32{},
		ImageIds: []int32{},
	}

	if world.DescriptionPostID.Valid {
		fetchInterface.PostIds = append(fetchInterface.PostIds, world.DescriptionPostID.Int32)
	}

	fetchIdsHeader, err := json.Marshal(fetchInterface)

	md := metadata.Pairs(
		"X-Fetch-Ids", string(fetchIdsHeader),
	)

	err = grpc.SendHeader(ctx, md)
	if err != nil {
		return nil, err
	}

	return converters.ConvertWorld(world), nil
}

func validateGetWorldById(req *pb.GetWorldByIdRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateWorldId(req.GetWorldId()); err != nil {
		violations = append(violations, e.FieldViolation("world_id", err))
	}

	return violations
}
