package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/the-medo/talebound-backend/api/converters"
	"github.com/the-medo/talebound-backend/api/e"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *Server) UpdateWorldAdmin(ctx context.Context, req *pb.UpdateWorldAdminRequest) (*pb.WorldAdmin, error) {
	violations := validateUpdateWorldAdmin(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	authPayload, err := server.authorizeUserCookie(ctx)
	if err != nil {
		return nil, e.UnauthenticatedError(err)
	}

	isAdmin, err := server.Store.IsWorldAdmin(ctx, db.IsWorldAdminParams{
		UserID:  authPayload.UserId,
		WorldID: req.GetWorldId(),
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user is not an admin of this world")
		}
		return nil, fmt.Errorf("failed to authorize world admin: %w", err)
	}

	if (req.SuperAdmin != nil || req.Approved != nil) && !isAdmin.SuperAdmin {
		return nil, fmt.Errorf("WORLD SUPER ADMIN role required for this action")
	}

	arg := db.UpdateWorldAdminParams{
		WorldID: req.GetWorldId(),
		UserID:  req.GetUserId(),

		SuperAdmin: sql.NullBool{
			Bool:  req.GetSuperAdmin(),
			Valid: req.SuperAdmin != nil,
		},
		Approved: sql.NullInt32{
			Int32: req.GetApproved(),
			Valid: req.Approved != nil,
		},
		MotivationalLetter: sql.NullString{
			String: req.GetMotivationalLetter(),
			Valid:  req.MotivationalLetter != nil,
		},
	}

	worldAdmin, err := server.Store.UpdateWorldAdmin(ctx, arg)
	if err != nil {
		return nil, err
	}

	user, err := server.Store.GetUserById(ctx, req.GetUserId())

	rsp := converters.ConvertWorldAdmin(worldAdmin, user)

	return rsp, nil
}

func validateUpdateWorldAdmin(req *pb.UpdateWorldAdminRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := validator.ValidateUserId(req.GetUserId()); err != nil {
		violations = append(violations, e.FieldViolation("user_id", err))
	}

	if err := validator.ValidateWorldId(req.GetWorldId()); err != nil {
		violations = append(violations, e.FieldViolation("world_id", err))
	}

	if req.Approved != nil {
		if err := validator.ValidateWorldAdminApproved(req.GetApproved()); err != nil {
			violations = append(violations, e.FieldViolation("approved", err))
		}
	}

	if req.MotivationalLetter != nil {
		if err := validator.ValidateWorldAdminMotivationalLetter(req.GetMotivationalLetter()); err != nil {
			violations = append(violations, e.FieldViolation("motivational_letter", err))
		}
	}

	return violations
}
