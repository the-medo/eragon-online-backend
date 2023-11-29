package menus

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
	"google.golang.org/grpc/metadata"
)

func (server *ServiceMenus) GetMenuItems(ctx context.Context, req *pb.GetMenuItemsRequest) (*pb.GetMenuItemsResponse, error) {
	violations := validateGetMenuItemsRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	menuItemRows, err := server.Store.GetMenuItems(ctx, req.MenuId)
	if err != nil {
		return nil, err
	}

	rsp := &pb.GetMenuItemsResponse{
		MenuItems: make([]*pb.MenuItem, len(menuItemRows)),
	}

	fetchInterface := &apihelpers.FetchInterface{
		PostIds: []int32{},
	}

	for i, menuItemRow := range menuItemRows {
		rsp.MenuItems[i] = converters.ConvertMenuItem(menuItemRow)
		if menuItemRow.DescriptionPostID.Valid {
			fetchInterface.PostIds = append(fetchInterface.PostIds, menuItemRow.DescriptionPostID.Int32)
		}
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

func validateGetMenuItemsRequest(req *pb.GetMenuItemsRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateMenuId(req.GetMenuId()); err != nil {
		violations = append(violations, e.FieldViolation("menu_id", err))
	}

	return violations
}
