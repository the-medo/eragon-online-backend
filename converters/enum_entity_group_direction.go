package converters

import (
	"github.com/the-medo/talebound-backend/pb"
)

var entityGroupDirectionToPB = map[string]pb.EntityGroupDirection{
	"unknown":    pb.EntityGroupDirection_ENTITY_GROUP_DIRECTION_UNKNOWN,
	"horizontal": pb.EntityGroupDirection_ENTITY_GROUP_DIRECTION_HORIZONTAL,
	"vertical":   pb.EntityGroupDirection_ENTITY_GROUP_DIRECTION_VERTICAL,
}

var entityGroupDirectionToDB = map[pb.EntityGroupDirection]string{
	pb.EntityGroupDirection_ENTITY_GROUP_DIRECTION_UNKNOWN:    "unknown",
	pb.EntityGroupDirection_ENTITY_GROUP_DIRECTION_HORIZONTAL: "horizontal",
	pb.EntityGroupDirection_ENTITY_GROUP_DIRECTION_VERTICAL:   "vertical",
}

func ConvertEntityGroupDirectionToPB(shape string) pb.EntityGroupDirection {
	if val, ok := entityGroupDirectionToPB[shape]; ok {
		return val
	}
	return pb.EntityGroupDirection_ENTITY_GROUP_DIRECTION_UNKNOWN
}

func ConvertEntityGroupDirectionToDB(shape pb.EntityGroupDirection) string {
	if val, ok := entityGroupDirectionToDB[shape]; ok {
		return val
	}
	return "unknown"
}
