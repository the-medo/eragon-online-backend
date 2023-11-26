package images

import (
	"context"
	"database/sql"
	"github.com/the-medo/talebound-backend/api/apihelpers"
	"github.com/the-medo/talebound-backend/converters"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *ServiceImages) GetImages(ctx context.Context, req *pb.GetImagesRequest) (*pb.GetImagesResponse, error) {
	violations := validateGetImagesRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	limit, offset := apihelpers.GetDefaultQueryBoundaries(req.GetLimit(), req.GetOffset())

	_, err := server.AuthorizeUserCookie(ctx)
	if err != nil {
		return nil, e.UnauthenticatedError(err)
	}

	arg := db.GetImagesParams{
		UserID: sql.NullInt32{
			Int32: req.GetUserId(),
			Valid: req.UserId != nil,
		},
		ImageTypeID: sql.NullInt32{
			Int32: req.GetImageTypeId(),
			Valid: req.ImageTypeId != nil,
		},
		PageLimit:  limit,
		PageOffset: offset,
	}

	images, err := server.Store.GetImages(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get posts: %v", err)
	}

	imageCountArg := convertGetImagesParamsToGetImagesCountParams(&arg)
	totalCount, err := server.Store.GetImagesCount(ctx, imageCountArg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get posts count: %v", err)
	}

	rsp := &pb.GetImagesResponse{
		Images:     make([]*pb.Image, len(images)),
		TotalCount: int32(totalCount),
	}

	for i, image := range images {
		rsp.Images[i] = converters.ConvertImage(image)
	}

	return rsp, nil
}

func convertGetImagesParamsToGetImagesCountParams(arg *db.GetImagesParams) db.GetImagesCountParams {
	return db.GetImagesCountParams{
		UserID:      arg.UserID,
		ImageTypeID: arg.ImageTypeID,
	}
}

func validateGetImagesRequest(req *pb.GetImagesRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if req.UserId != nil {
		if err := validator.ValidateUserId(req.GetUserId()); err != nil {
			violations = append(violations, e.FieldViolation("user_id", err))
		}
	}

	if req.ImageTypeId != nil {
		if err := validator.ValidateImageTypeId(req.GetImageTypeId()); err != nil {
			violations = append(violations, e.FieldViolation("image_type_id", err))
		}
	}

	if req.Limit != nil {
		if err := validator.ValidateLimit(req.GetLimit()); err != nil {
			violations = append(violations, e.FieldViolation("limit", err))
		}
	}

	if req.Offset != nil {
		if err := validator.ValidateOffset(req.GetOffset()); err != nil {
			violations = append(violations, e.FieldViolation("offset", err))
		}
	}

	return violations
}
