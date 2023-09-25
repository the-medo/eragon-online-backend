package api

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/the-medo/talebound-backend/api/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) DeleteChatMessage(ctx context.Context, req *pb.DeleteChatMessageRequest) (*pb.DeleteChatMessageResponse, error) {
	violations := validateDeleteChatMessages(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	err := server.CheckUserRole(ctx, []pb.RoleType{pb.RoleType_admin, pb.RoleType_moderator})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "failed to delete chat message: %v", err)
	}

	_, err = server.Store.GetChatMessage(ctx, req.GetId())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "message not found: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to get chat message: %v", err)
	}

	err = server.Store.DeleteChatMessage(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete chat message: %v", err)
	}

	return &pb.DeleteChatMessageResponse{
		Success: true,
		Message: fmt.Sprintf("Chat message %v deleted successfully.", req.GetId()),
	}, nil
}

func validateDeleteChatMessages(req *pb.DeleteChatMessageRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateChatMessageId(req.GetId()); err != nil {
		violations = append(violations, e.FieldViolation("id", err))
	}

	return violations
}
