package posts

import (
	"context"
	"database/sql"
	"github.com/the-medo/talebound-backend/api/converters"
	"github.com/the-medo/talebound-backend/api/e"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *ServicePosts) GetPostsByPlacement(ctx context.Context, req *pb.GetPostsByPlacementRequest) (*pb.GetPostsByPlacementResponse, error) {
	violations := validateGetPostsByPlacement(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	postRows, err := server.Store.GetPostsByPlacement(ctx, db.GetPostsByPlacementParams{
		PageOffset: req.GetOffset(),
		PageLimit:  req.GetLimit(),
		WorldID: sql.NullInt32{
			Int32: req.GetPlacement().GetWorldId(),
			Valid: req.GetPlacement().WorldId != nil,
		},
	})
	if err != nil {
		return nil, err
	}

	rsp := &pb.GetPostsByPlacementResponse{
		Posts:      make([]*pb.Post, len(postRows)),
		TotalCount: 0,
	}

	for i, postRow := range postRows {
		rsp.Posts[i] = converters.ConvertViewPostByPlacementToPost(postRow)
	}

	return rsp, nil
}

func validateGetPostsByPlacement(req *pb.GetPostsByPlacementRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateLimit(req.GetLimit()); err != nil {
		violations = append(violations, e.FieldViolation("limit", err))
	}

	if err := validator.ValidateOffset(req.GetOffset()); err != nil {
		violations = append(violations, e.FieldViolation("offset", err))
	}

	if err := validator.ValidateLocationPlacement(req.GetPlacement()); err != nil {
		violations = append(violations, e.FieldViolation("placements", err))
	}

	return violations
}
