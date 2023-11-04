package auth

import (
	"context"
	"database/sql"
	"github.com/the-medo/talebound-backend/api"
	"github.com/the-medo/talebound-backend/api/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *api.Server) ResetPasswordVerifyCodeValidity(ctx context.Context, req *pb.ResetPasswordVerifyCodeValidityRequest) (*pb.ResetPasswordVerifyCodeValidityResponse, error) {

	var rsp *pb.ResetPasswordVerifyCodeValidityResponse

	violations := validateResetPasswordVerifyCodeValidity(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	_, err := server.Store.GetUserPasswordReset(ctx, req.GetSecretCode())
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

func validateResetPasswordVerifyCodeValidity(req *pb.ResetPasswordVerifyCodeValidityRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := validator.ValidateSecretCode(req.GetSecretCode()); err != nil {
		violations = append(violations, e.FieldViolation("secret_code", err))
	}

	return violations
}
