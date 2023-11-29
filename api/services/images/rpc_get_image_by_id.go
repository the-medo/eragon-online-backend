package images

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

func (server *ServiceImages) GetImageById(ctx context.Context, req *pb.GetImageByIdRequest) (*pb.Image, error) {
	violations := validateGetImageByIdRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	image, err := server.Store.GetImageById(ctx, req.GetImageId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get image: %v", err)
	}

	fetchInterface := &apihelpers.FetchInterface{
		UserIds: []int32{image.UserID},
	}

	fetchIdsHeader, err := json.Marshal(fetchInterface)

	md := metadata.Pairs(
		"X-Fetch-Ids", string(fetchIdsHeader),
	)

	err = grpc.SendHeader(ctx, md)
	if err != nil {
		return nil, err
	}

	return converters.ConvertImage(image), nil
}

func validateGetImageByIdRequest(req *pb.GetImageByIdRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateUniversalId(req.GetImageId()); err != nil {
		violations = append(violations, e.FieldViolation("image_id", err))
	}
	return violations
}
