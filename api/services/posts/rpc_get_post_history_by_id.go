package posts

import (
	"context"
	"github.com/the-medo/talebound-backend/converters"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *ServicePosts) GetPostHistoryById(ctx context.Context, req *pb.GetPostHistoryByIdRequest) (*pb.PostHistory, error) {
	violations := validateGetPostHistoryById(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	authPayload, err := server.AuthorizeUserCookie(ctx)
	if err != nil {
		return nil, e.UnauthenticatedError(err)
	}

	post, err := server.Store.GetPostHistoryById(ctx, req.GetPostHistoryId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get post history: %v", err)
	}

	if post.UserID != authPayload.UserId {
		err := server.CheckUserRole(ctx, []pb.RoleType{pb.RoleType_admin})
		if err != nil {
			return nil, status.Errorf(codes.PermissionDenied, "can not get post history - you are not creator or admin: %v", err)
		}
	}

	rsp := converters.ConvertPostHistory(post)

	return rsp, nil
}

func validateGetPostHistoryById(req *pb.GetPostHistoryByIdRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidatePostId(req.GetPostId()); err != nil {
		violations = append(violations, e.FieldViolation("post_id", err))
	}

	if err := validator.ValidatePostHistoryId(req.GetPostHistoryId()); err != nil {
		violations = append(violations, e.FieldViolation("post_history_id", err))
	}

	return violations
}
