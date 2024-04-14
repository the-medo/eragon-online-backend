package maps

import (
	"context"
	"github.com/the-medo/talebound-backend/converters"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *ServiceMaps) GetModuleMapPinTypes(ctx context.Context, request *pb.GetModuleMapPinTypesRequest) (*pb.GetModuleMapPinTypesResponse, error) {
	violations := validateGetModuleMapPinTypes(request)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	pinTypeRows, err := server.Store.GetMapPinTypesForModule(ctx, request.GetModuleId())
	if err != nil {
		return nil, err
	}

	pinTypeGroupRows, err := server.Store.GetMapPinTypeGroupsForModule(ctx, request.GetModuleId())
	if err != nil {
		return nil, err
	}

	rsp := &pb.GetModuleMapPinTypesResponse{
		PinTypes:      make([]*pb.MapPinType, len(pinTypeRows)),
		PinTypeGroups: make([]*pb.MapPinTypeGroup, len(pinTypeGroupRows)),
	}

	for i, pinTypeRow := range pinTypeRows {
		rsp.PinTypes[i] = converters.ConvertMapPinType(pinTypeRow)
	}

	for i, pinTypeGroupRow := range pinTypeGroupRows {
		rsp.PinTypeGroups[i] = converters.ConvertMapPinTypeGroup(pinTypeGroupRow)
	}

	return rsp, nil
}

func validateGetModuleMapPinTypes(req *pb.GetModuleMapPinTypesRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateMapId(req.GetModuleId()); err != nil {
		violations = append(violations, e.FieldViolation("module_id", err))
	}

	return violations
}
