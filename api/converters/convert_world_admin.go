package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ConvertWorldAdmin(dbWorldAdmin db.WorldAdmin, dbViewUser db.ViewUser) *pb.WorldAdmin {
	pbWorldAdmin := &pb.WorldAdmin{
		WorldId:            dbWorldAdmin.WorldID,
		UserId:             dbWorldAdmin.UserID,
		User:               ConvertViewUser(dbViewUser),
		CreatedAt:          timestamppb.New(dbWorldAdmin.CreatedAt),
		SuperAdmin:         dbWorldAdmin.SuperAdmin,
		Approved:           dbWorldAdmin.Approved,
		MotivationalLetter: dbWorldAdmin.MotivationalLetter,
	}

	return pbWorldAdmin
}
