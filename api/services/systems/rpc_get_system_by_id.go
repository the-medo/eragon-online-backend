package systems

import (
	"context"
	"github.com/the-medo/talebound-backend/converters"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *ServiceSystems) GetSystemById(ctx context.Context, req *pb.GetSystemByIdRequest) (*pb.System, error) {
	violations := validateGetSystemById(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	system, err := server.Store.GetSystemByID(ctx, req.SystemId)
	if err != nil {
		return nil, err
	}

	return converters.ConvertSystem(system), nil
}

func validateGetSystemById(req *pb.GetSystemByIdRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateModuleId(req.GetSystemId()); err != nil {
		violations = append(violations, e.FieldViolation("system_id", err))
	}

	return violations
}
