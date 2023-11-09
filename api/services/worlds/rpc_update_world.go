package worlds

import (
	"context"
	"database/sql"
	"github.com/the-medo/talebound-backend/api/servicecore"
	"github.com/the-medo/talebound-backend/converters"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *ServiceWorlds) UpdateWorld(ctx context.Context, req *pb.UpdateWorldRequest) (*pb.ViewWorld, error) {
	violations := validateUpdateWorldRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	var needsEntityPermission []db.EntityType

	if req.DescriptionPostId != nil {
		needsEntityPermission = append(needsEntityPermission, db.EntityTypePost)
	}

	_, _, err := server.CheckModuleTypePermissions(ctx, db.ModuleTypeWorld, req.GetWorldId(), &servicecore.ModulePermission{
		NeedsSuperAdmin:       true,
		NeedsEntityPermission: &needsEntityPermission,
	})

	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "failed to update world: %v", err)
	}

	arg := db.UpdateWorldParams{
		WorldID: req.GetWorldId(),
		Name: sql.NullString{
			String: req.GetName(),
			Valid:  req.Name != nil,
		},
		BasedOn: sql.NullString{
			String: req.GetBasedOn(),
			Valid:  req.BasedOn != nil,
		},
		ShortDescription: sql.NullString{
			String: req.GetShortDescription(),
			Valid:  req.ShortDescription != nil,
		},
		Public: sql.NullBool{
			Bool:  req.GetPublic(),
			Valid: req.Public != nil,
		},
		DescriptionPostID: sql.NullInt32{
			Int32: req.GetDescriptionPostId(),
			Valid: req.DescriptionPostId != nil,
		},
	}

	_, err = server.Store.UpdateWorld(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update world: %v", err)
	}

	world, err := server.Store.GetWorldByID(ctx, req.GetWorldId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve updated world: %v", err)
	}

	rsp := converters.ConvertViewWorld(world)

	return rsp, nil
}

func validateUpdateWorldRequest(req *pb.UpdateWorldRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := validator.ValidateWorldId(req.GetWorldId()); err != nil {
		violations = append(violations, e.FieldViolation("world_id", err))
	}

	if req.Name != nil {
		if err := validator.ValidateWorldName(req.GetName()); err != nil {
			violations = append(violations, e.FieldViolation("name", err))
		}
	}

	if req.ShortDescription != nil {
		if err := validator.ValidateWorldShortDescription(req.GetShortDescription()); err != nil {
			violations = append(violations, e.FieldViolation("short_description", err))
		}
	}

	if req.BasedOn != nil {
		if err := validator.ValidateWorldBasedOn(req.GetBasedOn()); err != nil {
			violations = append(violations, e.FieldViolation("based_on", err))
		}
	}

	if req.DescriptionPostId != nil {
		if err := validator.ValidatePostId(req.GetDescriptionPostId()); err != nil {
			violations = append(violations, e.FieldViolation("description_post_id", err))
		}
	}

	return violations
}
