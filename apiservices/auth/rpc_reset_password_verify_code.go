package auth

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

func (server *ServiceAuth) ResetPasswordVerifyCode(ctx context.Context, req *pb.ResetPasswordVerifyCodeRequest) (*pb.ResetPasswordVerifyCodeResponse, error) {
	violations := validateResetPasswordVerifyCode(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	arg := db.ResetPasswordVerifyTxParams{
		Code:        req.GetSecretCode(),
		NewPassword: req.GetNewPassword(),
	}

	_, err := server.Store.ResetPasswordVerifyTx(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to reset password: %s", err)
	}

	rsp := &pb.ResetPasswordVerifyCodeResponse{
		Success: true,
		Message: "Password changed!",
	}

	return rsp, nil
}

func validateResetPasswordVerifyCode(req *pb.ResetPasswordVerifyCodeRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateSecretCode(req.GetSecretCode()); err != nil {
		violations = append(violations, e.FieldViolation("secret_code", err))
	}

	if err := validator.ValidatePassword(req.GetNewPassword()); err != nil {
		violations = append(violations, e.FieldViolation("new_password", err))
	}

	return violations
}
