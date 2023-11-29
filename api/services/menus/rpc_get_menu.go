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

func (server *ServiceMenus) GetMenu(ctx context.Context, req *pb.GetMenuRequest) (*pb.ViewMenu, error) {
	violations := validateGetMenuRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	menu, err := server.Store.GetMenu(ctx, req.MenuId)
	if err != nil {
		return nil, err
	}

	rsp := converters.ConvertViewMenu(menu)

	fetchInterface := &apihelpers.FetchInterface{
		ImageIds: []int32{},
	}

	if menu.MenuHeaderImgID.Valid {
		fetchInterface.ImageIds = append(fetchInterface.ImageIds, menu.MenuHeaderImgID.Int32)
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

func validateGetMenuRequest(req *pb.GetMenuRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateMenuId(req.GetMenuId()); err != nil {
		violations = append(violations, e.FieldViolation("menu_id", err))
	}

	return violations
}
