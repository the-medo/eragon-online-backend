package api

import (
	"context"
	"github.com/hibiken/asynq"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"github.com/the-medo/talebound-backend/worker"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func (server *Server) ResetPasswordSendCode(ctx context.Context, req *pb.ResetPasswordSendCodeRequest) (*pb.ResetPasswordSendCodeResponse, error) {
	violations := validateResetPasswordSendCode(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	arg := db.ResetPasswordRequestTxParams{
		Email: req.Email,
		AfterCreate: func(user db.ViewUser) error {
			taskPayload := &worker.PayloadSendResetPasswordEmail{
				Email: req.Email,
			}

			opts := []asynq.Option{
				asynq.MaxRetry(3),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}

			return server.taskDistributor.DistributeTaskSendResetPasswordEmail(ctx, taskPayload, opts...)
		},
	}

	_, err := server.store.ResetPasswordRequestTx(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to reset password: %s", err)
	}

	rsp := &pb.ResetPasswordSendCodeResponse{
		Success: true,
		Message: "Email is queued",
	}

	return rsp, nil
}

func validateResetPasswordSendCode(req *pb.ResetPasswordSendCodeRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, FieldViolation("email", err))
	}

	return violations
}
