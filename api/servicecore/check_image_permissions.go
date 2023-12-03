package servicecore

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/token"
)

// CheckImagePermissions  - Check if the user has permissions to access the image
//  1. check admin role
//  2. check if the image is an entity inside of module (images can be created outside of modules as well)
//     a. if it is, check if the user has permissions to access the entity
//     b. if it is not, check if the user is the owner of the image
func (core *ServiceCore) CheckImagePermissions(ctx context.Context, imageId int32, modulePermissions *ModulePermission) (*token.Payload, error) {
	authPayload, err := core.AuthorizeUserCookie(ctx)
	if err != nil {
		return nil, e.UnauthenticatedError(err)
	}

	err = core.CheckUserRole(ctx, []pb.RoleType{pb.RoleType_admin})
	if err == nil {
		return authPayload, nil
	}

	hasEntity := true
	entity, err := core.Store.GetEntityByImageId(ctx, sql.NullInt32{Int32: imageId, Valid: true})

	if err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("failed to entity image : %w", err)
		}
		hasEntity = false
	}

	if hasEntity {
		return core.CheckEntityPermissions(ctx, entity.ID, modulePermissions)
	}

	image, err := core.Store.GetImageById(ctx, imageId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("image does not exist")
		}
		return nil, fmt.Errorf("failed to get image : %w", err)
	}

	if image.UserID != authPayload.UserId {
		return nil, fmt.Errorf("user is not the owner of the image")
	}

	return authPayload, nil
}
