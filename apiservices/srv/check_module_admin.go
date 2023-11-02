package srv

import (
	"context"
	"fmt"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
)

func (core *ServiceCore) CheckModuleAdmin(ctx context.Context, moduleAdmin *db.ViewModuleAdmin, modulePermissions *ModulePermission) error {
	err := core.CheckUserRole(ctx, []pb.RoleType{pb.RoleType_admin})
	if err == nil {
		return nil
	}

	if moduleAdmin == nil {
		return fmt.Errorf("failed to authorize module admin: not found")
	}

	if moduleAdmin.Approved != 1 {
		return fmt.Errorf("failed to authorize module admin: not approved")
	}

	if modulePermissions != nil {
		if modulePermissions.NeedsMenuPermission {
			if moduleAdmin.AllowedMenu {
				return nil
			} else {
				return fmt.Errorf("MENU permission required for this action")
			}
		}

		if modulePermissions.NeedsEntityPermission != nil {
			for _, need := range *modulePermissions.NeedsEntityPermission {
				permissionFound := false
				for _, entityType := range moduleAdmin.AllowedEntityTypes {
					if entityType == need {
						permissionFound = true
						break
					}
				}
				if !permissionFound {
					return fmt.Errorf("%s permission required for this action", need)
				}
			}
		}

		if modulePermissions.NeedsSuperAdmin {
			if moduleAdmin.SuperAdmin {
				return nil
			} else {
				return fmt.Errorf("SUPER ADMIN role required for this action")
			}
		}
	}

	return nil
}
