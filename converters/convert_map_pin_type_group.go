package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
)

// ConvertMapPinTypeGroup converts a db.MapPinTypeGroup to pb.MapPinTypeGroup
func ConvertMapPinTypeGroup(mapPinTypeGroup db.MapPinTypeGroup) *pb.MapPinTypeGroup {
	pbType := &pb.MapPinTypeGroup{
		Id:   mapPinTypeGroup.ID,
		Name: mapPinTypeGroup.Name,
	}
	return pbType
}
