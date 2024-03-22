package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
)

var deleteEntityGroupContentActionToPB = map[db.DeleteEntityGroupContentAction]pb.DeleteEntityGroupContentAction{
	db.DeleteEntityGroupContentActionUnknown:        pb.DeleteEntityGroupContentAction_DELETE_EGC_ACTION_UNKNOWN,
	db.DeleteEntityGroupContentActionMoveChildren:   pb.DeleteEntityGroupContentAction_DELETE_EGC_ACTION_MOVE_CHILDREN,
	db.DeleteEntityGroupContentActionDeleteChildren: pb.DeleteEntityGroupContentAction_DELETE_EGC_ACTION_DELETE_CHILDREN,
}

var deleteEntityGroupContentActionToDB = map[pb.DeleteEntityGroupContentAction]db.DeleteEntityGroupContentAction{
	pb.DeleteEntityGroupContentAction_DELETE_EGC_ACTION_UNKNOWN:         db.DeleteEntityGroupContentActionUnknown,
	pb.DeleteEntityGroupContentAction_DELETE_EGC_ACTION_MOVE_CHILDREN:   db.DeleteEntityGroupContentActionMoveChildren,
	pb.DeleteEntityGroupContentAction_DELETE_EGC_ACTION_DELETE_CHILDREN: db.DeleteEntityGroupContentActionDeleteChildren,
}

func ConvertDeleteEntityGroupContentActionToPB(action db.DeleteEntityGroupContentAction) pb.DeleteEntityGroupContentAction {
	if val, ok := deleteEntityGroupContentActionToPB[action]; ok {
		return val
	}
	return pb.DeleteEntityGroupContentAction_DELETE_EGC_ACTION_UNKNOWN
}

func ConvertDeleteEntityGroupContentActionToDB(action pb.DeleteEntityGroupContentAction) db.DeleteEntityGroupContentAction {
	if val, ok := deleteEntityGroupContentActionToDB[action]; ok {
		return val
	}
	return db.DeleteEntityGroupContentActionUnknown
}
