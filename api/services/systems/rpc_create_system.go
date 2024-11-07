package systems

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

func (server *ServiceSystems) CreateSystem(ctx context.Context, req *pb.CreateSystemRequest) (*pb.CreateSystemResponse, error) {
	violations := validateCreateSystemRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	authPayload, err := server.AuthorizeUserCookie(ctx)
	if err != nil {
		return nil, e.UnauthenticatedError(err)
	}

	arg := db.CreateSystemTxParams{
		CreateSystemParams: db.CreateSystemParams{
			Name:             req.GetName(),
			BasedOn:          req.GetBasedOn(),
			ShortDescription: req.GetShortDescription(),
		},
		UserId: authPayload.UserId,
	}

	txResult, err := server.Store.CreateSystemTx(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "name already exists: %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create system: %s", err)
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

	rsp := &pb.CreateSystemResponse{
		System: converters.ConvertSystem(*txResult.System),
		Module: converters.ConvertViewModule(*txResult.Module),
	}

	return rsp, nil
}

func validateCreateSystemRequest(req *pb.CreateSystemRequest) (violations []*errdetails.BadRequest_FieldViolation) {
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
