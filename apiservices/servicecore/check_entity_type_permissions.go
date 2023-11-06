package servicecore

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/token"
)

func (core *ServiceCore) CheckEntityTypePermissions(ctx context.Context, entityType db.EntityType, id int32, modulePermissions *ModulePermission) (*token.Payload, error) {

	// Check if modulePermissions is nil and initialize if necessary
	if modulePermissions == nil {
		modulePermissions = &ModulePermission{
			NeedsEntityPermission: &[]db.EntityType{entityType}, // Initialize with the current entityType
		}

		// Check if NeedsEntityPermission is not nil
	} else if modulePermissions.NeedsEntityPermission != nil {
		// Check if entityType is not in NeedsEntityPermission slice
		found := false
		for _, et := range *modulePermissions.NeedsEntityPermission {
			if et == entityType {
				found = true
				break
			}
		}
		// If not found, append it
		if !found {
			*modulePermissions.NeedsEntityPermission = append(*modulePermissions.NeedsEntityPermission, entityType)
		}
	} else {
		// If the slice is nil, initialize it with the current entityType
		modulePermissions.NeedsEntityPermission = &[]db.EntityType{entityType}
	}

	var entity db.Entity
	var err error

	if entityType == db.EntityTypePost {
		entity, err = core.Store.GetEntityByPostId(ctx, sql.NullInt32{Int32: id, Valid: true})
	} else if entityType == db.EntityTypeImage {
		entity, err = core.Store.GetEntityByImageId(ctx, sql.NullInt32{Int32: id, Valid: true})
	} else if entityType == db.EntityTypeLocation {
		entity, err = core.Store.GetEntityByLocationId(ctx, sql.NullInt32{Int32: id, Valid: true})
	} else if entityType == db.EntityTypeMap {
		entity, err = core.Store.GetEntityByMapId(ctx, sql.NullInt32{Int32: id, Valid: true})
	} else {
		return nil, fmt.Errorf("entity type not implemented")
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%v does not exist", entityType)
		}
		return nil, fmt.Errorf("failed to get entity : %w", err)
	}

	return core.CheckEntityPermissions(ctx, entity.ID, modulePermissions)
}
