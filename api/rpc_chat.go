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
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) GetChatMessages(ctx context.Context, req *pb.GetChatMessagesRequest) (*pb.GetChatMessagesResponse, error) {
	violations := validateGetChatMessages(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	limit := req.GetLimit()
	if limit == 0 {
		limit = 100
	}

	offset := req.GetOffset()
	if offset < 0 {
		offset = 0
	}

	messages, err := server.store.GetChatMessages(ctx, db.GetChatMessagesParams{
		PageLimit:  limit,
		PageOffset: offset,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get messages: %v", err)
	}

	rsp := &pb.GetChatMessagesResponse{
		Messages: make([]*pb.ChatMessage, len(messages)),
	}

	for i, message := range messages {
		rsp.Messages[i] = convertChatMessage(message)
	}

	return rsp, nil

}

func (server *Server) AddChatMessage(ctx context.Context, req *pb.AddChatMessageRequest) (*pb.AddChatMessageResponse, error) {
	violations := validateAddChatMessage(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	user, err := server.store.GetUserById(ctx, authPayload.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	message, err := server.store.AddChatMessage(ctx, db.AddChatMessageParams{
		Text:   req.GetText(),
		UserID: user.ID,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to Add message: %v", err)
	}

	return &pb.AddChatMessageResponse{
		Message: &pb.ChatMessage{
			Id:        message.ID,
			UserId:    message.UserID,
			Username:  user.Username,
			Text:      message.Text,
			CreatedAt: timestamppb.New(message.CreatedAt),
		},
	}, nil
}

func (server *Server) DeleteChatMessage(ctx context.Context, req *pb.DeleteChatMessageRequest) (*pb.DeleteChatMessageResponse, error) {
	violations := validateDeleteChatMessages(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	err := server.CheckUserRole(ctx, []pb.RoleType{pb.RoleType_admin, pb.RoleType_moderator})
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "failed to delete chat message: %v", err)
	}

	_, err = server.store.GetChatMessage(ctx, req.GetId())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "message not found: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to get chat message: %v", err)
	}

	err = server.store.DeleteChatMessage(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete chat message: %v", err)
	}

	return &pb.DeleteChatMessageResponse{
		Success: true,
		Message: fmt.Sprintf("Chat message %v deleted successfully.", req.GetId()),
	}, nil
}

func validateGetChatMessages(req *pb.GetChatMessagesRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateLimitOrOffset(req.GetLimit(), 1000); err != nil {
		violations = append(violations, FieldViolation("limit", err))
	}
	if err := validator.ValidateLimitOrOffset(req.GetOffset()); err != nil {
		violations = append(violations, FieldViolation("offset", err))
	}

	return violations
}

func validateAddChatMessage(req *pb.AddChatMessageRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateString(req.GetText(), 1, 1024); err != nil {
		violations = append(violations, FieldViolation("text", err))
	}

	return violations
}

func validateDeleteChatMessages(req *pb.DeleteChatMessageRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateInt64(req.GetId(), 1); err != nil {
		violations = append(violations, FieldViolation("id", err))
	}

	return violations
}
