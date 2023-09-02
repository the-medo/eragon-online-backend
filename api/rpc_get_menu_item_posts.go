package api

import (
	"context"
	"database/sql"
	"github.com/the-medo/talebound-backend/api/converters"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *Server) GetMenuItemPosts(ctx context.Context, req *pb.GetMenuItemPostsRequest) (*pb.GetMenuItemPostsResponse, error) {
	violations := validateGetMenuItemPostsRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	menuItemId := sql.NullInt32{
		Int32: req.GetMenuItemId(),
		Valid: true,
	}
	menuItemRows, err := server.store.GetMenuItemPosts(ctx, menuItemId)
	if err != nil {
		return nil, err
	}

	rsp := &pb.GetMenuItemPostsResponse{
		MenuItemPosts: make([]*pb.MenuItemPost, len(menuItemRows)),
	}

	for i, menuItemPostRow := range menuItemRows {

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

func validateGetMenuItemPostsRequest(req *pb.GetMenuItemPostsRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateMenuId(req.GetMenuId()); err != nil {
		violations = append(violations, FieldViolation("menu_id", err))
	}

	if err := validator.ValidateMenuItemId(req.GetMenuItemId()); err != nil {
		violations = append(violations, FieldViolation("menu_item_id", err))
	}

	return violations
}
