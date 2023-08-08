package api

import (
	"context"
	"github.com/the-medo/talebound-backend/api/converters"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *Server) CreateWorldAdmin(ctx context.Context, request *pb.CreateWorldAdminRequest) (*pb.WorldAdmin, error) {
	violations := validateCreateWorldAdmin(request)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	authPayload, err := server.authorizeUserCookie(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	arg := db.InsertWorldAdminParams{
		WorldID:            request.GetWorldId(),
		UserID:             authPayload.UserId,
		SuperAdmin:         false,
		Approved:           2,
		MotivationalLetter: request.GetMotivationalLetter(),
	}

	worldAdmin, err := server.store.InsertWorldAdmin(ctx, arg)
	if err != nil {
		return nil, err
	}

	user, err := server.store.GetUserById(ctx, authPayload.UserId)

	rsp := converters.ConvertWorldAdmin(worldAdmin, user)

	return rsp, nil
}

func validateCreateWorldAdmin(req *pb.CreateWorldAdminRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateInt(req.GetWorldId(), 1, 4098); err != nil {
		violations = append(violations, FieldViolation("world_id", err))
	}

	if err := validator.ValidateString(req.GetMotivationalLetter(), 0, 2000); err != nil {
		violations = append(violations, FieldViolation("motivational_letter", err))
	}

	return violations
}
