package api

import (
	"context"
	"database/sql"
	"fmt"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math"
)

func StringToEvaluationType(value string) (pb.EvaluationType, error) {
	if enumValue, ok := pb.EvaluationType_value[value]; ok {
		return pb.EvaluationType(enumValue), nil
	}
	return pb.EvaluationType_self, fmt.Errorf("invalid EvaluationType: %s", value)
}

func (server *Server) GetEvaluationById(ctx context.Context, req *pb.GetEvaluationByIdRequest) (*pb.GetEvaluationByIdResponse, error) {
	evaluation, err := server.store.GetEvaluationById(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user roles: %v", err)
	}

	evaluationType, err := StringToEvaluationType(string(evaluation.EvaluationType))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert evaluation type: %v", err)
	}

	rsp := &pb.GetEvaluationByIdResponse{
		Evaluation: &pb.Evaluation{
			Id:          evaluation.ID,
			Name:        evaluation.Name,
			Description: evaluation.Description,
			Type:        evaluationType,
		},
	}
	return rsp, nil
}

func (server *Server) GetEvaluationsByType(ctx context.Context, req *pb.GetEvaluationsByTypeRequest) (*pb.GetEvaluationsByTypeResponse, error) {
	evaluations, err := server.store.GetEvaluationsByType(ctx, db.EvaluationType(req.GetType()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get evaluation list: %v", err)
	}

	rsp := &pb.GetEvaluationsByTypeResponse{
		Evaluation: make([]*pb.Evaluation, len(evaluations)),
	}

	for i, evaluation := range evaluations {
		evaluationType, err := StringToEvaluationType(string(evaluation.EvaluationType))
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to convert evaluation type: %v", err)
		}

		rsp.Evaluation[i] = &pb.Evaluation{
			Id:          evaluation.ID,
			Name:        evaluation.Name,
			Description: evaluation.Description,
			Type:        evaluationType,
		}
	}

	return rsp, nil
}

func (server *Server) CreateOrUpdateEvaluationVote(ctx context.Context, req *pb.CreateOrUpdateEvaluationVoteRequest) (*pb.CreateOrUpdateEvaluationVoteResponse, error) {
	arg := db.GetEvaluationVoteByEvaluationIdUserIdAndVoterParams{
		EvaluationID: req.GetEvaluationId(),
		UserID:       req.GetUserId(),
		UserIDVoter:  req.GetUserIdVoter(),
	}

	var alreadyExists bool
	alreadyExists = true

	_, err := server.store.GetEvaluationVoteByEvaluationIdUserIdAndVoter(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			alreadyExists = false
		} else {
			return nil, status.Errorf(codes.Internal, "failed to get evaluation vote: %v", err)
		}
	}

	var rsp *pb.CreateOrUpdateEvaluationVoteResponse

	if alreadyExists {

		evaluationVoteNew, err := server.store.UpdateEvaluationVote(ctx, db.UpdateEvaluationVoteParams{
			EvaluationID: req.GetEvaluationId(),
			UserID:       req.GetUserId(),
			UserIDVoter:  req.GetUserIdVoter(),
			Value:        req.GetValue(),
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to update evaluation vote: %v", err)
		}

		rsp.EvaluationVote = &pb.EvaluationVote{
			EvaluationId: evaluationVoteNew.EvaluationID,
			UserId:       evaluationVoteNew.UserID,
			UserIdVoter:  evaluationVoteNew.UserIDVoter,
			Value:        evaluationVoteNew.Value,
			CreatedAt:    timestamppb.New(evaluationVoteNew.CreatedAt),
		}
	} else {
		evaluationVoteNew, err := server.store.CreateEvaluationVote(ctx, db.CreateEvaluationVoteParams{
			EvaluationID: req.GetEvaluationId(),
			UserID:       req.GetUserId(),
			UserIDVoter:  req.GetUserIdVoter(),
			Value:        req.GetValue(),
		})
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to create evaluation vote: %v", err)
		}

		rsp.EvaluationVote = &pb.EvaluationVote{
			EvaluationId: evaluationVoteNew.EvaluationID,
			UserId:       evaluationVoteNew.UserID,
			UserIdVoter:  evaluationVoteNew.UserIDVoter,
			Value:        evaluationVoteNew.Value,
			CreatedAt:    timestamppb.New(evaluationVoteNew.CreatedAt),
		}
	}

	return rsp, nil
}

func (server *Server) GetEvaluationVotesByUserId(ctx context.Context, req *pb.GetEvaluationVotesByUserIdRequest) (*pb.GetEvaluationVotesByUserIdResponse, error) {
	evaluationVotes, err := server.store.GetEvaluationVotesByUserId(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get evaluation votes : %v", err)
	}

	rsp := &pb.GetEvaluationVotesByUserIdResponse{
		EvaluationVote: make([]*pb.EvaluationVote, len(evaluationVotes)),
	}

	for i, evaluationVote := range evaluationVotes {
		rsp.EvaluationVote[i] = convertEvaluationVote(evaluationVote)
	}
	return rsp, nil
}

func (server *Server) GetEvaluationVotesByUserIdAndVoter(ctx context.Context, req *pb.GetEvaluationVotesByUserIdAndVoterRequest) (*pb.GetEvaluationVotesByUserIdAndVoterResponse, error) {
	arg := db.GetEvaluationVotesByUserIdAndVoterParams{
		UserID:      req.GetUserId(),
		UserIDVoter: req.GetUserIdVoter(),
	}

	evaluationVotes, err := server.store.GetEvaluationVotesByUserIdAndVoter(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get evaluation votes : %v", err)
	}

	rsp := &pb.GetEvaluationVotesByUserIdAndVoterResponse{
		EvaluationVote: make([]*pb.EvaluationVote, len(evaluationVotes)),
	}

	for i, evaluationVote := range evaluationVotes {
		rsp.EvaluationVote[i] = convertEvaluationVote(evaluationVote)
	}
	return rsp, nil
}

func (server *Server) DeleteEvaluationVote(ctx context.Context, req *pb.DeleteEvaluationVoteRequest) (*pb.DeleteEvaluationVoteResponse, error) {
	arg := db.DeleteEvaluationVoteParams{
		EvaluationID: req.GetEvaluationId(),
		UserID:       req.GetUserId(),
		UserIDVoter:  req.GetUserIdVoter(),
	}

	err := server.store.DeleteEvaluationVote(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get evaluation votes : %v", err)
	}

	rsp := &pb.DeleteEvaluationVoteResponse{
		Success: true,
		Message: "Evaluation vote deleted successfully",
	}

	return rsp, nil
}

func (server *Server) GetAverageUserEvaluationsByType(ctx context.Context, req *pb.GetAverageUserEvaluationsByTypeRequest) (*pb.GetAverageUserEvaluationsByTypeResponse, error) {
	evaluationType := db.EvaluationType(req.GetType())

	arg := db.GetAverageUserEvaluationsByTypeParams{
		UserID:         req.GetUserId(),
		EvaluationType: evaluationType,
	}

	avgEvaluationVotes, err := server.store.GetAverageUserEvaluationsByType(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get evaluation averages: %v", err)
	}

	rsp := &pb.GetAverageUserEvaluationsByTypeResponse{
		AverageEvaluationVote: make([]*pb.AverageEvaluationVote, len(avgEvaluationVotes)),
	}

	for i, e := range avgEvaluationVotes {
		evaluationType, err := StringToEvaluationType(string(e.EvaluationType))
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to convert evaluation type: %v", err)
		}

		rsp.AverageEvaluationVote[i] = &pb.AverageEvaluationVote{
			EvaluationId: e.EvaluationID,
			UserId:       req.GetUserId(),
			Name:         e.Name,
			Description:  e.Description,
			Type:         evaluationType,
			Average:      float32(math.Round(e.AvgValue*100) / 100),
		}
	}

	return rsp, nil
}
