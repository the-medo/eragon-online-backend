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

func (server *Server) UpdateWorldIntroduction(ctx context.Context, req *pb.UpdateWorldIntroductionRequest) (*pb.Post, error) {
	violations := validateUpdateWorldIntroductionRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	authPayload, err := server.CheckWorldAdmin(ctx, req.GetWorldId(), false)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "failed update introduction: %v", err)
	}

	world, err := server.store.GetWorldByID(ctx, req.GetWorldId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "unable to save changes - world not found: %s", err)
	}

	postType, err := server.store.GetPostTypeById(ctx, PostTypeWorldDescription)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get post type: %v", err)
	}

	if !world.DescriptionPostID.Valid {
		//create new post
		createPostArg := db.CreatePostParams{
			UserID:     authPayload.UserId,
			Title:      "World introduction",
			PostTypeID: PostTypeWorldDescription,
			Content:    req.GetContent(),
			IsDraft:    false,
			IsPrivate:  false,
		}

		post, err := server.store.CreatePost(ctx, createPostArg)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to create post: %s", err)
		}

		//if it is post (so first time introduction is created), we put it into users table
		updateWorldArg := db.UpdateWorldParams{
			WorldID: req.WorldId,
			DescriptionPostID: sql.NullInt32{
				Int32: post.ID,
				Valid: true,
			},
		}
		_, err = server.store.UpdateWorld(ctx, updateWorldArg)

		return convertPostAndPostType(post, postType), nil
	} else {
		//update existing post
		arg := db.UpdatePostParams{
			PostID: world.DescriptionPostID.Int32,
			LastUpdatedUserID: sql.NullInt32{
				Int32: authPayload.UserId,
				Valid: true,
			},
			Content: sql.NullString{
				String: req.GetContent(),
				Valid:  true,
			},
		}
		post, err := server.store.UpdatePost(ctx, arg)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to update post: %s", err)
		}

		return convertPostAndPostType(post, postType), nil
	}
}

func validateUpdateWorldIntroductionRequest(req *pb.UpdateWorldIntroductionRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateWorldId(req.GetWorldId()); err != nil {
		violations = append(violations, FieldViolation("world_id", err))
	}

	return violations
}
