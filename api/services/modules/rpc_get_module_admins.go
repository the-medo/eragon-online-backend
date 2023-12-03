package modules

import (
	"context"
	"encoding/json"
	"github.com/the-medo/talebound-backend/api/apihelpers"
	"github.com/the-medo/talebound-backend/converters"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/util"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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

	fetchInterface := &apihelpers.FetchInterface{
		ModuleIds: []int32{req.GetModuleId()},
		UserIds:   []int32{},
	}

	for i, moduleAdminRow := range moduleAdminRows {
		rsp.ModuleAdmins[i] = converters.ConvertModuleAdminRow(moduleAdminRow)
		fetchInterface.UserIds = util.Upsert(fetchInterface.UserIds, rsp.ModuleAdmins[i].UserId)
	}

	fetchIdsHeader, err := json.Marshal(fetchInterface)

	md := metadata.Pairs(
		"X-Fetch-Ids", string(fetchIdsHeader),
	)

	err = grpc.SendHeader(ctx, md)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

func validateGetModuleAdmins(req *pb.GetModuleAdminsRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateUniversalId(req.GetModuleId()); err != nil {
		violations = append(violations, e.FieldViolation("module_id", err))
	}

	return violations
}
