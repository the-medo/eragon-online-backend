package modules

import (
	"context"
	"github.com/the-medo/talebound-backend/converters"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *ServiceModules) GetModuleAdmins(ctx context.Context, req *pb.GetModuleAdminsRequest) (*pb.GetModuleAdminsResponse, error) {
	violations := validateGetModuleAdmins(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	moduleAdminRows, err := server.Store.GetModuleAdmins(ctx, req.GetModuleId())
	if err != nil {
		return nil, err
	}

	rsp := &pb.GetModuleAdminsResponse{
		ModuleAdmins: make([]*pb.ModuleAdmin, len(moduleAdminRows)),
	}

	for i, moduleAdminRow := range moduleAdminRows {
		rsp.ModuleAdmins[i] = converters.ConvertModuleAdminRow(moduleAdminRow)
	}

	return rsp, nil
}

func validateGetModuleAdmins(req *pb.GetModuleAdminsRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateUniversalId(req.GetModuleId()); err != nil {
		violations = append(violations, e.FieldViolation("module_id", err))
	}

	return violations
}
