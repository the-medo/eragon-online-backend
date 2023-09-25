package api

import (
	"context"
	"github.com/the-medo/talebound-backend/api/converters"
	"github.com/the-medo/talebound-backend/api/e"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *Server) GetMenuItemPostsByMenuId(ctx context.Context, req *pb.GetMenuItemPostsByMenuIdRequest) (*pb.GetMenuItemPostsByMenuIdResponse, error) {
	violations := validateGetMenuItemPostsByMenuIdRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	menuItemPostRows, err := server.Store.GetMenuItemPostsByMenuId(ctx, req.GetMenuId())
	if err != nil {
		return nil, err
	}

	rsp := &pb.GetMenuItemPostsByMenuIdResponse{
		MenuItemPosts: make([]*pb.MenuItemPost, len(menuItemPostRows)),
	}

	for i, menuItemPostRow := range menuItemPostRows {

		mip := db.MenuItemPost{
			MenuID:     menuItemPostRow.MenuID,
			MenuItemID: menuItemPostRow.MenuItemID,
			PostID:     menuItemPostRow.PostID,
			Position:   menuItemPostRow.Position,
		}

		post := db.ViewPost{
			ID:                 menuItemPostRow.PostID,
			PostTypeID:         menuItemPostRow.PostTypeID,
			UserID:             menuItemPostRow.UserID,
			Title:              menuItemPostRow.Title,
			Content:            "",
			CreatedAt:          menuItemPostRow.CreatedAt,
			DeletedAt:          menuItemPostRow.DeletedAt,
			LastUpdatedAt:      menuItemPostRow.LastUpdatedAt,
			LastUpdatedUserID:  menuItemPostRow.LastUpdatedUserID,
			IsDraft:            menuItemPostRow.IsDraft,
			IsPrivate:          menuItemPostRow.IsPrivate,
			Description:        menuItemPostRow.Description,
			ThumbnailImgID:     menuItemPostRow.ThumbnailImgID,
			PostTypeName:       menuItemPostRow.PostTypeName,
			PostTypeDraftable:  menuItemPostRow.PostTypeDraftable,
			PostTypePrivatable: menuItemPostRow.PostTypePrivatable,
			ThumbnailImgUrl:    menuItemPostRow.ThumbnailImgUrl,
		}

		rsp.MenuItemPosts[i] = converters.ConvertMenuItemPost(mip, post)
	}

	return rsp, nil
}

func validateGetMenuItemPostsByMenuIdRequest(req *pb.GetMenuItemPostsByMenuIdRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateMenuId(req.GetMenuId()); err != nil {
		violations = append(violations, e.FieldViolation("menu_id", err))
	}

	return violations
}
