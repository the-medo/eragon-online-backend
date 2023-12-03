package images

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/the-medo/talebound-backend/api/apihelpers"
	"github.com/the-medo/talebound-backend/converters"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/util"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
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
		PageLimit:  limit,
		PageOffset: offset,
		UserID: sql.NullInt32{
			Int32: req.GetUserId(),
			Valid: req.UserId != nil,
		},
		ModuleID: sql.NullInt32{
			Int32: req.GetModuleId(),
			Valid: req.ModuleId != nil,
		},
		Width: sql.NullInt32{
			Int32: req.GetWidth(),
			Valid: req.Width != nil,
		},
		Height: sql.NullInt32{
			Int32: req.GetHeight(),
			Valid: req.Height != nil,
		},
		OrderBy: sql.NullString{
			String: req.GetOrderBy(),
			Valid:  req.OrderBy != nil,
		},
		Tags: req.GetTags(),
	}

	images, err := server.Store.GetImages(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get posts: %v", err)
	}

	rsp := &pb.GetImagesResponse{
		Images:     make([]*pb.Image, len(images)),
		TotalCount: 0,
	}

	if len(images) > 0 {
		rsp.TotalCount = images[0].TotalCount
	}

	fetchInterface := &apihelpers.FetchInterface{
		UserIds:   []int32{},
		ModuleIds: []int32{},
	}

	for i, image := range images {
		rsp.Images[i] = converters.ConvertGetImagesRow(image)

		fetchInterface.UserIds = util.Upsert(fetchInterface.UserIds, image.UserID)

		if image.ModuleID.Valid {
			fetchInterface.ModuleIds = util.Upsert(fetchInterface.ModuleIds, image.ModuleID.Int32)
		}
	}

	fetchIdsHeader, err := json.Marshal(fetchInterface)

	md := metadata.Pairs(
		"X-Fetch-Ids", string(fetchIdsHeader),
	)

	err = grpc.SendHeader(ctx, md)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

func validateGetImagesRequest(req *pb.GetImagesRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	fields := []string{"name", "created_at", "description", "width", "height"}

	if req.OrderBy != nil {
		if validator.StringInSlice(req.GetOrderBy(), fields) == false {
			violations = append(violations, e.FieldViolation("order_by", fmt.Errorf("invalid field to order by %s", req.GetOrderBy())))
		}
	}

	if req.UserId != nil {
		if err := validator.ValidateUserId(req.GetUserId()); err != nil {
			violations = append(violations, e.FieldViolation("user_id", err))
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
