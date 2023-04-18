package api

import (
	"context"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) GetChatMessages(ctx context.Context, req *pb.GetChatMessagesRequest) (*pb.GetChatMessagesResponse, error) {
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
	_, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	message, err := server.store.AddChatMessage(ctx, db.AddChatMessageParams{
		Text:   req.GetText(),
		UserID: req.GetUserId(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to Add message: %v", err)
	}

	user, err := server.store.GetUserById(ctx, message.UserID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
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
