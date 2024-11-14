package characters

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

func (server *ServiceCharacters) UpdateCharacter(ctx context.Context, req *pb.UpdateCharacterRequest) (*pb.Character, error) {
	violations := validateUpdateCharacterRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	var needsEntityPermission []db.EntityType

	_, _, err := server.CheckModuleTypePermissions(ctx, db.ModuleTypeCharacter, req.GetCharacterId(), &servicecore.ModulePermission{
		NeedsSuperAdmin:       true,
		NeedsEntityPermission: &needsEntityPermission,
	})

	if err != nil {
		return nil, err
	}

	// TODO - make WorldID and SystemID updateable only when they are "Universal world" or "Universal system"
	// otherwise, another related entities can be out of sync with what is possible

	arg := db.UpdateCharacterParams{
		CharacterID:      req.GetCharacterId(),
		Name:             sql.NullString{String: req.GetName(), Valid: req.Name != nil},
		ShortDescription: sql.NullString{String: req.GetShortDescription(), Valid: req.ShortDescription != nil},
		WorldID:          sql.NullInt32{Int32: req.GetWorldId(), Valid: req.WorldId != nil},
		SystemID:         sql.NullInt32{Int32: req.GetSystemId(), Valid: req.SystemId != nil},
	}

	character, err := server.Store.UpdateCharacter(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update character: %s", err)
	}

	return converters.ConvertCharacter(character), nil
}

func validateUpdateCharacterRequest(req *pb.UpdateCharacterRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateModuleId(req.GetCharacterId()); err != nil {
		violations = append(violations, e.FieldViolation("character_id", err))
	}

	if req.Name != nil {
		if err := validator.ValidateModuleName(req.GetName()); err != nil {
			violations = append(violations, e.FieldViolation("name", err))
		}
	}

	if req.ShortDescription != nil {
		if err := validator.ValidateModuleShortDescription(req.GetShortDescription()); err != nil {
			violations = append(violations, e.FieldViolation("short_description", err))
		}
	}

	if req.WorldId != nil {
		if err := validator.ValidateUniversalId(req.GetWorldId()); err != nil {
			violations = append(violations, e.FieldViolation("world_id", err))
		}
	}

	if req.SystemId != nil {
		if err := validator.ValidateUniversalId(req.GetSystemId()); err != nil {
			violations = append(violations, e.FieldViolation("system_id", err))
		}
	}

	return violations
}
