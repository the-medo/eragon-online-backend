package servicecore

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/token"
)

// CheckPostPermissions  - Check if the user has permissions to access the post
//  1. check admin role
//  2. check if the post is entity inside of module (posts can be created outside of modules as well)
//     a. if it is, check if the user has permissions to access the entity
//     b. if it is not, check if the user is the owner of the post
func (core *ServiceCore) CheckPostPermissions(ctx context.Context, postId int32, modulePermissions *ModulePermission) (*token.Payload, error) {
	authPayload, err := core.AuthorizeUserCookie(ctx)
	if err != nil {
		return nil, e.UnauthenticatedError(err)
	}

	err = core.CheckUserRole(ctx, []pb.RoleType{pb.RoleType_admin})
	if err == nil {
		return authPayload, nil
	}

	hasEntity := true
	entity, err := core.Store.GetEntityByPostId(ctx, sql.NullInt32{Int32: postId, Valid: true})

	if err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("failed to entity post : %w", err)
		}
		hasEntity = false
	}

	if hasEntity {
		return core.CheckEntityPermissions(ctx, entity.ID, modulePermissions)
	}

	post, err := core.Store.GetPostById(ctx, postId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("post does not exist")
		}
		return nil, fmt.Errorf("failed to get post : %w", err)
	}

	if post.UserID != authPayload.UserId {
		return nil, fmt.Errorf("user is not the owner of the post")
	}

	return authPayload, nil
}
