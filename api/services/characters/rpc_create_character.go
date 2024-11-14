package characters

import (
	"context"
	"database/sql"
	"github.com/lib/pq"
	"github.com/the-medo/talebound-backend/converters"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *ServiceCharacters) CreateCharacter(ctx context.Context, req *pb.CreateCharacterRequest) (*pb.CreateCharacterResponse, error) {
	violations := validateCreateCharacterRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	authPayload, err := server.AuthorizeUserCookie(ctx)
	if err != nil {
		return nil, e.UnauthenticatedError(err)
	}

	arg := db.CreateCharacterTxParams{
		CreateCharacterParams: db.CreateCharacterParams{
			Name:             req.GetName(),
			ShortDescription: req.GetShortDescription(),
			WorldID: sql.NullInt32{
				Int32: req.WorldId,
				Valid: true,
			},
			SystemID: sql.NullInt32{
				Int32: req.SystemId,
				Valid: true,
			},
		},
		UserId: authPayload.UserId,
	}

	txResult, err := server.Store.CreateCharacterTx(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "name already exists: %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create character: %s", err)
	}

	menuId := txResult.Module.MenuID

	pos1, isMainFalse := int32(1), false
	items := []pb.CreateMenuItemRequest{
		{MenuId: menuId, Code: "overview", Name: "Overview", Position: &pos1, IsMain: &isMainFalse},
	}

	for i := range items {
		_, err := server.SharedCreateMenuItem(ctx, &items[i])
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to create menu item: %s", err)
		}
	}

	rsp := &pb.CreateCharacterResponse{
		Character: converters.ConvertCharacter(*txResult.Character),
		Module:    converters.ConvertViewModule(*txResult.Module),
	}

	return rsp, nil
}

func validateCreateCharacterRequest(req *pb.CreateCharacterRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateUniversalId(req.GetSystemId()); err != nil {
		violations = append(violations, e.FieldViolation("system_id", err))
	}

	if err := validator.ValidateUniversalId(req.GetWorldId()); err != nil {
		violations = append(violations, e.FieldViolation("world_id", err))
	}

	if err := validator.ValidateModuleName(req.GetName()); err != nil {
		violations = append(violations, e.FieldViolation("name", err))
	}

	if err := validator.ValidateModuleShortDescription(req.GetShortDescription()); err != nil {
		violations = append(violations, e.FieldViolation("short_description", err))
	}

	return violations
}
