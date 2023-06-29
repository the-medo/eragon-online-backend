package api

import (
	"context"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) GetPostHistory(ctx context.Context, req *pb.GetPostHistoryRequest) (*pb.GetPostHistoryResponse, error) {
	violations := validateGetPostHistoryRequest(req)
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
			return nil, status.Errorf(codes.PermissionDenied, "can not get list of postHistory - you are not creator or admin: %v", err)
		}
	}

	postHistory, err := server.store.GetPostHistoryByPostId(ctx, req.GetPostId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get postHistory: %v", err)
	}

	rsp := &pb.GetPostHistoryResponse{
		HistoryPosts: make([]*pb.HistoryPost, len(postHistory)),
	}

	postType, err := server.store.GetPostTypeById(ctx, post.PostTypeID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get post type: %v", err)
	}

	for i, post := range postHistory {
		rsp.HistoryPosts[i] = convertHistoryPostWithoutContent(post, postType)
	}

	return rsp, nil
}

func (server *Server) GetPostHistoryById(ctx context.Context, req *pb.GetPostHistoryByIdRequest) (*pb.HistoryPost, error) {
	violations := validateGetPostHistoryByIdRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	authPayload, err := server.authorizeUserCookie(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	post, err := server.store.GetPostHistoryById(ctx, req.GetPostHistoryId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get post history: %v", err)
	}

	if post.UserID != authPayload.UserId {
		err := server.CheckUserRole(ctx, []pb.RoleType{pb.RoleType_admin})
		if err != nil {
			return nil, status.Errorf(codes.PermissionDenied, "can not get post history - you are not creator or admin: %v", err)
		}
	}

	postType, err := server.store.GetPostTypeById(ctx, post.PostTypeID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get post type: %v", err)
	}

	rsp := convertHistoryPost(post, postType)

	return rsp, nil
}

func validateGetPostHistoryRequest(req *pb.GetPostHistoryRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidatePostId(req.GetPostId()); err != nil {
		violations = append(violations, FieldViolation("post_id", err))
	}

	return violations
}

func validateGetPostHistoryByIdRequest(req *pb.GetPostHistoryByIdRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidatePostId(req.GetPostId()); err != nil {
		violations = append(violations, FieldViolation("post_id", err))
	}
	if err := validator.ValidatePostHistoryId(req.GetPostHistoryId()); err != nil {
		violations = append(violations, FieldViolation("post_history_id", err))
	}

	return violations
}
