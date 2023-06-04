package api

import (
	"context"
	"database/sql"
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

func (server *Server) ResetPasswordVerifyCode(ctx context.Context, req *pb.ResetPasswordVerifyCodeRequest) (*pb.ResetPasswordVerifyCodeResponse, error) {
	violations := validateResetPasswordVerifyCode(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	arg := db.ResetPasswordVerifyTxParams{
		Code:        req.GetSecretCode(),
		NewPassword: req.GetNewPassword(),
	}

	_, err := server.store.ResetPasswordVerifyTx(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to reset password: %s", err)
	}

	rsp := &pb.ResetPasswordVerifyCodeResponse{
		Success: true,
		Message: "Password changed!",
	}

	return rsp, nil
}

func (server *Server) ResetPasswordVerifyCodeValidity(ctx context.Context, req *pb.ResetPasswordVerifyCodeValidityRequest) (*pb.ResetPasswordVerifyCodeValidityResponse, error) {

	var rsp *pb.ResetPasswordVerifyCodeValidityResponse

	println("ResetPasswordVerifyCodeValidity - ", req.GetSecretCode())

	violations := validateResetPasswordVerifyCodeValidity(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	_, err := server.store.GetUserPasswordReset(ctx, req.GetSecretCode())
	if err != nil {
		if err == sql.ErrNoRows {
			rsp = &pb.ResetPasswordVerifyCodeValidityResponse{
				Success: false,
				Message: "Code is invalid",
			}
			return rsp, nil
		}
		return nil, status.Errorf(codes.Internal, "failed to validate secret code: %s", err)
	}

	rsp = &pb.ResetPasswordVerifyCodeValidityResponse{
		Success: true,
		Message: "Code is valid",
	}

	return rsp, nil
}

func validateResetPasswordSendCode(req *pb.ResetPasswordSendCodeRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, FieldViolation("email", err))
	}

	return violations
}

func validateResetPasswordVerifyCode(req *pb.ResetPasswordVerifyCodeRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateSecretCode(req.GetSecretCode()); err != nil {
		violations = append(violations, FieldViolation("secret_code", err))
	}

	if err := validator.ValidatePassword(req.GetNewPassword()); err != nil {
		violations = append(violations, FieldViolation("new_password", err))
	}

	return violations
}

func validateResetPasswordVerifyCodeValidity(req *pb.ResetPasswordVerifyCodeValidityRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := validator.ValidateSecretCode(req.GetSecretCode()); err != nil {
		violations = append(violations, FieldViolation("secret_code", err))
	}

	return violations
}
