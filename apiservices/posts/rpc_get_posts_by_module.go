package posts

import (
	"context"
	"database/sql"
	"github.com/the-medo/talebound-backend/api"
	"github.com/the-medo/talebound-backend/api/converters"
	"github.com/the-medo/talebound-backend/api/e"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *api.Server) GetPostsByModule(ctx context.Context, req *pb.GetPostsByModuleRequest) (*pb.GetPostsByModuleResponse, error) {
	violations := validateGetPostsByModule(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	postRows, err := server.Store.GetPostsByModule(ctx, db.GetPostsByModuleParams{
		PageOffset: req.GetOffset(),
		PageLimit:  req.GetLimit(),
		WorldID: sql.NullInt32{
			Int32: req.GetWorldId(),
			Valid: req.WorldId != nil,
		},
	})
	if err != nil {
		return nil, err
	}

	rsp := &pb.GetPostsByModuleResponse{
		Posts:      make([]*pb.Post, len(postRows)),
		TotalCount: 0,
	}

	if len(postRows) > 0 {
		rsp.TotalCount = postRows[0].TotalCount
	}

	for i, postRow := range postRows {
		rsp.Posts[i] = converters.ConvertViewPostByModuleToPost(postRow)
	}

	return rsp, nil
}

func validateGetPostsByModule(req *pb.GetPostsByModuleRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateLimit(req.GetLimit()); err != nil {
		violations = append(violations, e.FieldViolation("limit", err))
	}

	if err := validator.ValidateOffset(req.GetOffset()); err != nil {
		violations = append(violations, e.FieldViolation("offset", err))
	}

	if err := validator.ValidatePostModuleExtended(req.WorldId, req.QuestId, req.SystemId, req.CharacterId); err != nil {
		violations = append(violations, e.FieldViolation("modules", err))
	}

	return violations
}
