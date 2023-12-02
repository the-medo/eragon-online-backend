package posts

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
	"google.golang.org/grpc/metadata"
)

func (server *ServicePosts) GetPosts(ctx context.Context, req *pb.GetPostsRequest) (*pb.GetPostsResponse, error) {
	violations := validateGetPosts(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	limit, offset := apihelpers.GetDefaultQueryBoundaries(req.GetLimit(), req.GetOffset())

	arg := db.GetPostsParams{
		PageLimit:  limit,
		PageOffset: offset,
		ModuleID: sql.NullInt32{
			Int32: req.GetModuleId(),
			Valid: req.ModuleId != nil,
		},
		UserID: sql.NullInt32{
			Int32: req.GetUserId(),
			Valid: req.UserId != nil,
		},
		IsPrivate: sql.NullBool{
			Bool:  req.GetIsPrivate(),
			Valid: req.IsPrivate != nil,
		},
		IsDraft: sql.NullBool{
			Bool:  req.GetIsDraft(),
			Valid: req.IsDraft != nil,
		},
		OrderBy: sql.NullString{
			String: req.GetOrderBy(),
			Valid:  req.OrderBy != nil,
		},
		Tags: req.GetTags(),
	}

	postRows, err := server.Store.GetPosts(ctx, arg)

	if err != nil {
		return nil, err
	}

	rsp := &pb.GetPostsResponse{
		Posts:      make([]*pb.Post, len(postRows)),
		TotalCount: 0,
	}

	if len(postRows) > 0 {
		rsp.TotalCount = postRows[0].TotalCount
	}

	fetchInterface := &apihelpers.FetchInterface{
		UserIds:   []int32{},
		ModuleIds: []int32{},
		ImageIds:  []int32{},
	}

	for i, postRow := range postRows {
		rsp.Posts[i] = converters.ConvertGetPostsRow(postRow)

		fetchInterface.UserIds = util.Upsert(fetchInterface.UserIds, postRow.UserID)

		if postRow.ModuleID.Valid {
			fetchInterface.ModuleIds = util.Upsert(fetchInterface.ModuleIds, postRow.ModuleID.Int32)
		}

		if postRow.ThumbnailImgID.Valid {
			fetchInterface.ImageIds = util.Upsert(fetchInterface.ImageIds, postRow.ThumbnailImgID.Int32)
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

func validateGetPosts(req *pb.GetPostsRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	fields := []string{"title", "created_at", "description"}

	if req.OrderBy != nil {
		if validator.StringInSlice(req.GetOrderBy(), fields) == false {
			violations = append(violations, e.FieldViolation("order_by", fmt.Errorf("invalid field to order by %s", req.GetOrderBy())))
		}
	}

	if err := validator.ValidateLimit(req.GetLimit()); err != nil {
		violations = append(violations, e.FieldViolation("limit", err))
	}

	if err := validator.ValidateOffset(req.GetOffset()); err != nil {
		violations = append(violations, e.FieldViolation("offset", err))
	}

	if err := validator.ValidateUniversalId(req.GetModuleId()); err != nil {
		violations = append(violations, e.FieldViolation("module_id", err))
	}

	return violations
}
