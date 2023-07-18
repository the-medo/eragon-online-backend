package api

import (
	"context"
	"database/sql"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) GetImages(ctx context.Context, req *pb.GetImagesRequest) (*pb.GetImagesResponse, error) {
	violations := validateGetImagesRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	limit, offset := GetDefaultQueryBoundaries(req.GetLimit(), req.GetOffset())

	_, err := server.authorizeUserCookie(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
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

	images, err := server.store.GetImages(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get posts: %v", err)
	}

	imageCountArg := convertGetImagesParamsToGetImagesCountParams(&arg)
	totalCount, err := server.store.GetImagesCount(ctx, imageCountArg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get posts count: %v", err)
	}

	rsp := &pb.GetImagesResponse{
		Images:     make([]*pb.Image, len(images)),
		TotalCount: int32(totalCount),
	}

	for i, image := range images {
		rsp.Images[i] = convertImage(&image)
	}

	return rsp, nil
}

func (server *Server) GetUserImages(ctx context.Context, req *pb.GetUserImagesRequest) (*pb.GetImagesResponse, error) {
	violations := validateGetUserImagesRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	limit := req.GetLimit()
	offset := req.GetOffset()

	getImagesReq := &pb.GetImagesRequest{
		UserId:      &req.UserId,
		ImageTypeId: req.ImageTypeId,
		Limit:       &limit,
		Offset:      &offset,
	}

	return server.GetImages(ctx, getImagesReq)
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
			violations = append(violations, FieldViolation("user_id", err))
		}
	}

	if req.ImageTypeId != nil {
		if err := validator.ValidateImageTypeId(req.GetImageTypeId()); err != nil {
			violations = append(violations, FieldViolation("post_type_id", err))
		}
	}

	return violations
}

func validateGetUserImagesRequest(req *pb.GetUserImagesRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if req.ImageTypeId != nil {
		if err := validator.ValidateImageTypeId(req.GetImageTypeId()); err != nil {
			violations = append(violations, FieldViolation("post_type_id", err))
		}
	}

	return violations
}
