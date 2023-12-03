package entities

import (
	"context"
	"encoding/json"
	"github.com/the-medo/talebound-backend/api/apihelpers"
	"github.com/the-medo/talebound-backend/converters"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (server *ServiceEntities) GetEntityById(ctx context.Context, req *pb.GetEntityByIdRequest) (*pb.ViewEntity, error) {
	violations := validateGetEntityByIdRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	entity, err := server.Store.GetEntityByID(ctx, req.GetEntityId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get entity: %v", err)
	}

	fetchInterface := &apihelpers.FetchInterface{
		ModuleIds:   []int32{entity.ModuleID},
		PostIds:     []int32{},
		ImageIds:    []int32{},
		MapIds:      []int32{},
		LocationIds: []int32{},
	}

	if entity.PostID.Valid {
		fetchInterface.PostIds = append(fetchInterface.PostIds, entity.PostID.Int32)
	}

	if entity.ImageID.Valid {
		fetchInterface.ImageIds = append(fetchInterface.ImageIds, entity.ImageID.Int32)
	}

	if entity.MapID.Valid {
		fetchInterface.MapIds = append(fetchInterface.MapIds, entity.MapID.Int32)
	}

	if entity.LocationID.Valid {
		fetchInterface.LocationIds = append(fetchInterface.LocationIds, entity.LocationID.Int32)
	}

	fetchIdsHeader, err := json.Marshal(fetchInterface)

	md := metadata.Pairs(
		"X-Fetch-Ids", string(fetchIdsHeader),
	)

	err = grpc.SendHeader(ctx, md)
	if err != nil {
		return nil, err
	}

	return converters.ConvertViewEntity(entity), nil
}

func validateGetEntityByIdRequest(req *pb.GetEntityByIdRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateUniversalId(req.GetEntityId()); err != nil {
		violations = append(violations, e.FieldViolation("entity_id", err))
	}
	return violations
}
