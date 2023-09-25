package api

import (
	"context"
	"github.com/the-medo/talebound-backend/api/e"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) GetChatMessages(ctx context.Context, req *pb.GetChatMessagesRequest) (*pb.GetChatMessagesResponse, error) {
	violations := validateGetChatMessages(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	limit, offset := GetDefaultQueryBoundaries(req.GetLimit(), req.GetOffset())

	messages, err := server.Store.GetChatMessages(ctx, db.GetChatMessagesParams{
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

func validateGetChatMessages(req *pb.GetChatMessagesRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateLimit(req.GetLimit()); err != nil {
		violations = append(violations, e.FieldViolation("limit", err))
	}
	if err := validator.ValidateOffset(req.GetOffset()); err != nil {
		violations = append(violations, e.FieldViolation("offset", err))
	}

	return violations
}
