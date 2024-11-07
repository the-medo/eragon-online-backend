package worlds

import (
	"context"
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

func (server *ServiceWorlds) CreateWorld(ctx context.Context, req *pb.CreateWorldRequest) (*pb.CreateWorldResponse, error) {
	violations := validateCreateWorldRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	authPayload, err := server.AuthorizeUserCookie(ctx)
	if err != nil {
		return nil, e.UnauthenticatedError(err)
	}

	arg := db.CreateWorldTxParams{
		CreateWorldParams: db.CreateWorldParams{
			Name:             req.GetName(),
			BasedOn:          req.GetBasedOn(),
			ShortDescription: req.GetShortDescription(),
		},
		UserId: authPayload.UserId,
	}

	txResult, err := server.Store.CreateWorldTx(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "name already exists: %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create world: %s", err)
	}

	menuId := txResult.Module.MenuID

	pos1, pos2, pos3, pos4, pos5, pos6, isMainTrue, isMainFalse := int32(1), int32(2), int32(3), int32(4), int32(5), int32(6), true, false
	items := []pb.CreateMenuItemRequest{
		{MenuId: menuId, Code: "overview", Name: "Overview", Position: &pos1, IsMain: &isMainTrue},
		{MenuId: menuId, Code: "races", Name: "Races", Position: &pos2, IsMain: &isMainFalse},
		{MenuId: menuId, Code: "flora-and-fauna", Name: "Flora & Fauna", Position: &pos3, IsMain: &isMainFalse},
		{MenuId: menuId, Code: "magic", Name: "Magic", Position: &pos4, IsMain: &isMainFalse},
		{MenuId: menuId, Code: "science-and-technology", Name: "Science & Technology", Position: &pos5, IsMain: &isMainFalse},
		{MenuId: menuId, Code: "history", Name: "History", Position: &pos6, IsMain: &isMainFalse},
	}

	for i := range items {
		_, err := server.SharedCreateMenuItem(ctx, &items[i])
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to create menu item: %s", err)
		}
	}

	rsp := &pb.CreateWorldResponse{
		World:  converters.ConvertWorld(*txResult.World),
		Module: converters.ConvertViewModule(*txResult.Module),
	}

	return rsp, nil
}

func validateCreateWorldRequest(req *pb.CreateWorldRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateModuleName(req.GetName()); err != nil {
		violations = append(violations, e.FieldViolation("name", err))
	}

	if err := validator.ValidateModuleShortDescription(req.GetShortDescription()); err != nil {
		violations = append(violations, e.FieldViolation("short_description", err))
	}

	if err := validator.ValidateModuleBasedOn(req.GetBasedOn()); err != nil {
		violations = append(violations, e.FieldViolation("based_on", err))
	}

	return violations
}
