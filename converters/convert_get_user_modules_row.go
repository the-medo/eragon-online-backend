package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
)

func ConvertGetUserModulesRow(userModule db.GetUserModulesRow) *pb.UserModule {
	pbUserModule := &pb.UserModule{
		UserId:              userModule.UserID,
		ModuleId:            userModule.ID,
		Admin:               userModule.Admin,
		Favorite:            userModule.Favorite,
		Following:           userModule.Following,
		EntityNotifications: make([]pb.EntityType, len(userModule.EntityNotifications)),
	}

	for i, entityNotification := range userModule.EntityNotifications {
		pbUserModule.EntityNotifications[i] = ConvertEntityTypeToPB(entityNotification)
	}

	return pbUserModule
}
