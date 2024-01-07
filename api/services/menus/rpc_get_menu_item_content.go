package menus

import (
	"context"
	"database/sql"
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

func (server *ServiceMenus) GetMenuItemContent(ctx context.Context, req *pb.GetMenuItemContentRequest) (*pb.GetMenuItemContentResponse, error) {
	violations := validateGetMenuItemContentRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	menuItem, err := server.Store.GetMenuItemById(ctx, req.GetMenuItemId())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "menu item not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get menu item: %v", err)
	}

	if !menuItem.EntityGroupID.Valid {
		return nil, status.Errorf(codes.Internal, "menu item has no entity group %d", req.GetMenuItemId())
	}

	recursiveEntities, err := server.Store.GetRecursiveEntities(ctx, menuItem.EntityGroupID.Int32)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "no content not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get menu item content: %v", err)
	}

	groupIDs := []int32{menuItem.EntityGroupID.Int32}
	entityIDs := make([]int32, 0)
	contents := make([]*pb.EntityGroupContent, len(recursiveEntities))

	for i, entityGroupContent := range recursiveEntities {
		if entityGroupContent.ContentEntityGroupID.Valid {
			groupIDs = append(groupIDs, entityGroupContent.ContentEntityGroupID.Int32)
		}
		if entityGroupContent.ContentEntityID.Valid {
			entityIDs = append(entityIDs, entityGroupContent.ContentEntityID.Int32)
		}
		contents[i] = converters.ConvertEntityGroupContent(entityGroupContent)
	}

	entityGroups, err := server.Store.GetEntityGroupsByIDs(ctx, groupIDs)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "no entity groups found for this menu item")
		}
		return nil, status.Errorf(codes.Internal, "failed to get entity groups for this menu item: %v", err)
	}

	rsp := &pb.GetMenuItemContentResponse{
		MainGroupId: req.GetMenuItemId(),
		Groups:      make([]*pb.EntityGroup, len(entityGroups)),
		Contents:    contents,
	}

	for i, entityGroup := range entityGroups {
		rsp.Groups[i] = converters.ConvertEntityGroup(entityGroup)
	}

	fetchInterface := &apihelpers.FetchInterface{
		EntityIds: entityIDs,
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

func validateGetMenuItemContentRequest(req *pb.GetMenuItemContentRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateMenuId(req.GetMenuId()); err != nil {
		violations = append(violations, e.FieldViolation("menu_id", err))
	}
	if err := validator.ValidateMenuItemId(req.GetMenuItemId()); err != nil {
		violations = append(violations, e.FieldViolation("menu_item_id", err))
	}

	return violations
}
