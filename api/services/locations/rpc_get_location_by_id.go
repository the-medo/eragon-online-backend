package locations

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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (server *ServiceLocations) GetLocationById(ctx context.Context, req *pb.GetLocationByIdRequest) (*pb.Location, error) {
	violations := validateGetLocationByIdRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	m, err := server.Store.GetLocationById(ctx, req.GetLocationId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get location: %v", err)
	}

	fetchInterface := &apihelpers.FetchInterface{
		PostIds:  []int32{},
		ImageIds: []int32{},
	}

	if m.PostID.Valid {
		fetchInterface.PostIds = append(fetchInterface.PostIds, m.PostID.Int32)
	}

	if m.ThumbnailImageID.Valid {
		fetchInterface.ImageIds = append(fetchInterface.ImageIds, m.ThumbnailImageID.Int32)
	}

	fetchIdsHeader, err := json.Marshal(fetchInterface)

	md := metadata.Pairs(
		"X-Fetch-Ids", string(fetchIdsHeader),
	)

	err = grpc.SendHeader(ctx, md)
	if err != nil {
		return nil, err
	}

	return converters.ConvertLocation(m), nil
}

func validateGetLocationByIdRequest(req *pb.GetLocationByIdRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateUniversalId(req.GetLocationId()); err != nil {
		violations = append(violations, e.FieldViolation("location_id", err))
	}
	return violations
}
