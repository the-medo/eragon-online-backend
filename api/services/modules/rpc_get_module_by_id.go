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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
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

	fetchInterface := &apihelpers.FetchInterface{
		WorldIds: []int32{},
		ImageIds: []int32{},
		PostIds:  []int32{module.DescriptionPostID},
	}

	if module.WorldID.Valid {
		fetchInterface.WorldIds = append(fetchInterface.WorldIds, module.WorldID.Int32)
	}

	if module.HeaderImgID.Valid {
		fetchInterface.ImageIds = util.Upsert(fetchInterface.ImageIds, module.HeaderImgID.Int32)
	}

	if module.AvatarImgID.Valid {
		fetchInterface.ImageIds = util.Upsert(fetchInterface.ImageIds, module.AvatarImgID.Int32)
	}

	if module.AvatarImgID.Valid {
		fetchInterface.ImageIds = util.Upsert(fetchInterface.ImageIds, module.AvatarImgID.Int32)
	}

	fetchIdsHeader, err := json.Marshal(fetchInterface)

	md := metadata.Pairs(
		"X-Fetch-Ids", string(fetchIdsHeader),
	)

	err = grpc.SendHeader(ctx, md)
	if err != nil {
		return nil, err
	}

	return converters.ConvertViewModule(module), nil
}

func validateGetModuleByIdRequest(req *pb.GetModuleByIdRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateUniversalId(req.GetModuleId()); err != nil {
		violations = append(violations, e.FieldViolation("module_id", err))
	}
	return violations
}
