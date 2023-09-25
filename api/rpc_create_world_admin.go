package api

import (
	"context"
	"github.com/the-medo/talebound-backend/api/converters"
	"github.com/the-medo/talebound-backend/api/e"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *Server) CreateWorldAdmin(ctx context.Context, request *pb.CreateWorldAdminRequest) (*pb.WorldAdmin, error) {
	violations := validateCreateWorldAdmin(request)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	authPayload, err := server.authorizeUserCookie(ctx)
	if err != nil {
		return nil, e.UnauthenticatedError(err)
	}

	arg := db.InsertWorldAdminParams{
		WorldID:            request.GetWorldId(),
		UserID:             authPayload.UserId,
		SuperAdmin:         false,
		Approved:           2,
		MotivationalLetter: request.GetMotivationalLetter(),
	}

	worldAdmin, err := server.Store.InsertWorldAdmin(ctx, arg)
	if err != nil {
		return nil, err
	}

	user, err := server.Store.GetUserById(ctx, authPayload.UserId)

	rsp := converters.ConvertWorldAdmin(worldAdmin, user)

	return rsp, nil
}

func validateCreateWorldAdmin(req *pb.CreateWorldAdminRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateWorldId(req.GetWorldId()); err != nil {
		violations = append(violations, e.FieldViolation("world_id", err))
	}

	if err := validator.ValidateWorldAdminMotivationalLetter(req.GetMotivationalLetter()); err != nil {
		violations = append(violations, e.FieldViolation("motivational_letter", err))
	}

	return violations
}
