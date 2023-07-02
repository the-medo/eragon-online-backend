package api

import (
	"context"
	"database/sql"
	"fmt"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

const (
	PostTypeUniversal            = 100
	PostTypeQuestPost            = 200
	PostTypeWorldDescription     = 300
	PostTypeRuleSetDescription   = 400
	PostTypeQuestDescription     = 500
	PostTypeCharacterDescription = 600
	PostTypeNews                 = 700
	PostTypeUserIntroduction     = 800
)

func (server *Server) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.Post, error) {
	violations := validateCreatePostRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	authPayload, err := server.authorizeUserCookie(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	arg := db.CreatePostParams{
		UserID:     authPayload.UserId,
		PostTypeID: req.GetPostTypeId(),
		Title:      req.GetTitle(),
		Content:    req.GetContent(),
	}

	//in case we got a nil value for IsDraft or IsPrivate, we need to get the default value from the post type
	postTypeNeeded := false
	if req.IsDraft == nil || req.IsPrivate == nil {
		postTypeNeeded = true
	}

	if postTypeNeeded {
		postType, err := server.store.GetPostTypeById(ctx, req.GetPostTypeId())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get post type: %v", err)
		}
		if req.IsDraft == nil {
			arg.IsDraft = postType.Draftable
		}
		if req.IsPrivate == nil {
			arg.IsPrivate = postType.Privatable
		}
	} else {
		arg.IsDraft = req.GetIsDraft()
		arg.IsPrivate = req.GetIsPrivate()
	}

	postResult, err := server.store.CreatePost(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create post: %s", err)
	}

	postType, err := server.store.GetPostTypeById(ctx, postResult.PostTypeID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get post type: %v", err)
	}

	rsp := convertPostAndPostType(postResult, postType)

	return rsp, nil
}

func (server *Server) UpdatePost(ctx context.Context, req *pb.UpdatePostRequest) (*pb.Post, error) {
	violations := validateUpdatePostRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}
	authPayload, err := server.authorizeUserCookie(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	_, err = server.store.InsertPostHistory(ctx, req.GetPostId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to insert post history: %v", err)
	}

	arg := db.UpdatePostParams{
		PostID: req.GetPostId(),
		Title: sql.NullString{
			String: req.GetTitle(),
			Valid:  req.Title != nil,
		},
		Content: sql.NullString{
			String: req.GetContent(),
			Valid:  req.Title != nil,
		},
		PostTypeID: sql.NullInt32{
			Int32: req.GetPostTypeId(),
			Valid: req.PostTypeId != nil,
		},
		LastUpdatedUserID: sql.NullInt32{
			Int32: authPayload.UserId,
			Valid: true,
		},
		IsDraft: sql.NullBool{
			Bool:  req.GetIsDraft(),
			Valid: req.IsDraft != nil,
		},
		IsPrivate: sql.NullBool{
			Bool:  req.GetIsPrivate(),
			Valid: req.IsPrivate != nil,
		},
	}

	post, err := server.store.UpdatePost(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update post: %v", err)
	}

	postType, err := server.store.GetPostTypeById(ctx, post.PostTypeID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get post type: %v", err)
	}

	rsp := convertPostAndPostType(post, postType)

	return rsp, nil
}

func (server *Server) DeletePost(ctx context.Context, req *pb.DeletePostRequest) (*pb.DeletePostResponse, error) {
	violations := validateDeletePostRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	authPayload, err := server.authorizeUserCookie(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	post, err := server.store.GetPostById(ctx, req.GetPostId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get post: %v", err)
	}

	if post.UserID != authPayload.UserId {
		err := server.CheckUserRole(ctx, []pb.RoleType{pb.RoleType_admin})
		if err != nil {
			return nil, status.Errorf(codes.PermissionDenied, "can not delete this post - you are not creator or admin: %v", err)
		}
	}

	err = server.store.DeletePost(ctx, req.GetPostId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete post: %v", err)
	}

	return &pb.DeletePostResponse{
		Success: true,
		Message: fmt.Sprintf("Post %s (%v) deleted successfully.", post.Title, req.GetPostId()),
	}, nil
}

func (server *Server) GetPostById(ctx context.Context, req *pb.GetPostByIdRequest) (*pb.Post, error) {
	violations := validateGetPostByIdRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	//authPayload, err := server.authorizeUserCookie(ctx)
	//if err != nil {
	//	return nil, unauthenticatedError(err)
	//}

	post, err := server.store.GetPostById(ctx, req.GetPostId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get post: %v", err)
	}

	//if post.UserID != authPayload.UserId {
	//	err := server.CheckUserRole(ctx, []pb.RoleType{pb.RoleType_admin})
	//	if err != nil {
	//		return nil, status.Errorf(codes.PermissionDenied, "can not update this post - you are not creator or admin: %v", err)
	//	}
	//}

	rsp := convertViewPost(post)

	time.Sleep(2 * time.Second)

	return rsp, nil
}

func (server *Server) GetPostTypes(ctx context.Context, req *emptypb.Empty) (*pb.GetPostTypesResponse, error) {

	postTypes, err := server.store.GetPostTypes(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get posts: %v", err)
	}

	rsp := &pb.GetPostTypesResponse{
		PostTypes: make([]*pb.DataPostType, len(postTypes)),
	}

	for i, postType := range postTypes {
		rsp.PostTypes[i] = convertPostType(postType)
	}

	return rsp, nil
}

func (server *Server) GetUserPosts(ctx context.Context, req *pb.GetUserPostsRequest) (*pb.GetUserPostsResponse, error) {
	violations := validateGetUserPostsRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	limit, offset := GetDefaultQueryBoundaries(req.GetLimit(), req.GetOffset())

	authPayload, err := server.authorizeUserCookie(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	if req.GetUserId() != authPayload.UserId {
		err := server.CheckUserRole(ctx, []pb.RoleType{pb.RoleType_admin})
		if err != nil {
			return nil, status.Errorf(codes.PermissionDenied, "can not get list of posts - you are not creator or admin: %v", err)
		}
	}

	arg := db.GetPostsByUserIdParams{
		UserID: req.GetUserId(),
		PostTypeID: sql.NullInt32{
			Int32: req.GetPostTypeId(),
			Valid: req.PostTypeId != nil,
		},
		PageLimit:  limit,
		PageOffset: offset,
	}

	posts, err := server.store.GetPostsByUserId(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get posts: %v", err)
	}

	rsp := &pb.GetUserPostsResponse{
		Posts: make([]*pb.Post, len(posts)),
	}

	for i, post := range posts {
		rsp.Posts[i] = convertViewPost(post)
	}

	return rsp, nil
}

// ================= VALIDATION =================

func validateCreatePostRequest(req *pb.CreatePostRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidatePostTitle(req.GetTitle()); err != nil {
		violations = append(violations, FieldViolation("title", err))
	}

	return violations
}

func validateUpdatePostRequest(req *pb.UpdatePostRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidatePostId(req.GetPostId()); err != nil {
		violations = append(violations, FieldViolation("post_id", err))
	}

	if req.Title != nil {
		if err := validator.ValidatePostTitle(req.GetTitle()); err != nil {
			violations = append(violations, FieldViolation("title", err))
		}
	}

	if req.PostTypeId != nil {
		if err := validator.ValidatePostTypeId(req.GetPostTypeId()); err != nil {
			violations = append(violations, FieldViolation("post_type_id", err))
		}
	}

	return violations
}

func validateDeletePostRequest(req *pb.DeletePostRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidatePostId(req.GetPostId()); err != nil {
		violations = append(violations, FieldViolation("post_id", err))
	}

	return violations
}

func validateGetPostByIdRequest(req *pb.GetPostByIdRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidatePostId(req.GetPostId()); err != nil {
		violations = append(violations, FieldViolation("post_id", err))
	}

	return violations
}

func validateGetUserPostsRequest(req *pb.GetUserPostsRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateUserId(req.GetUserId()); err != nil {
		violations = append(violations, FieldViolation("user_id", err))
	}

	if req.PostTypeId != nil {
		if err := validator.ValidatePostTypeId(req.GetPostTypeId()); err != nil {
			violations = append(violations, FieldViolation("post_type_id", err))
		}
	}

	return violations
}
