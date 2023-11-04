package posts

import (
	"context"
	"github.com/the-medo/talebound-backend/api/converters"
	"github.com/the-medo/talebound-backend/api/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *ServicePosts) GetPostHistory(ctx context.Context, req *pb.GetPostHistoryRequest) (*pb.GetPostHistoryResponse, error) {
	violations := validateGetPostHistoryRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	authPayload, err := server.AuthorizeUserCookie(ctx)
	if err != nil {
		return nil, e.UnauthenticatedError(err)
	}

	post, err := server.Store.GetPostById(ctx, req.GetPostId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get post: %v", err)
	}

	if post.UserID != authPayload.UserId {
		err := server.CheckUserRole(ctx, []pb.RoleType{pb.RoleType_admin})
		if err != nil {
			return nil, status.Errorf(codes.PermissionDenied, "can not get list of postHistory - you are not creator or admin: %v", err)
		}
	}

	postHistory, err := server.Store.GetPostHistoryByPostId(ctx, req.GetPostId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get postHistory: %v", err)
	}

	rsp := &pb.GetPostHistoryResponse{
		HistoryPosts: make([]*pb.HistoryPost, len(postHistory)),
	}

	postType, err := server.Store.GetPostTypeById(ctx, post.PostTypeID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get post type: %v", err)
	}

	for i, post := range postHistory {
		rsp.HistoryPosts[i] = converters.ConvertHistoryPostWithoutContent(post, postType)
	}

	return rsp, nil
}

func validateGetPostHistoryRequest(req *pb.GetPostHistoryRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidatePostId(req.GetPostId()); err != nil {
		violations = append(violations, e.FieldViolation("post_id", err))
	}

	return violations
}
