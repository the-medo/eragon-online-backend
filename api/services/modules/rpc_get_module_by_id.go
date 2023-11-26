package modules

import (
	"context"
	"github.com/the-medo/talebound-backend/converters"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *ServiceModules) GetModuleById(ctx context.Context, req *pb.GetModuleByIdRequest) (*pb.ViewModule, error) {
	violations := validateGetModuleByIdRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	module, err := server.Store.GetModuleById(ctx, req.GetModuleId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get module: %v", err)
	}

	return converters.ConvertViewModule(module), nil
}

func validateGetModuleByIdRequest(req *pb.GetModuleByIdRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateUniversalId(req.GetModuleId()); err != nil {
		violations = append(violations, e.FieldViolation("module_id", err))
	}
	return violations
}
