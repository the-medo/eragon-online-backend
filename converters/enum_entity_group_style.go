package converters

import (
	"github.com/the-medo/talebound-backend/pb"
)

var entityGroupStyleToPB = map[string]pb.EntityGroupStyle{
	"unknown":    pb.EntityGroupStyle_ENTITY_GROUP_STYLE_UNKNOWN,
	"framed":     pb.EntityGroupStyle_ENTITY_GROUP_STYLE_FRAMED,
	"not-framed": pb.EntityGroupStyle_ENTITY_GROUP_STYLE_NOT_FRAMED,
}

var entityGroupStyleToDB = map[pb.EntityGroupStyle]string{
	pb.EntityGroupStyle_ENTITY_GROUP_STYLE_UNKNOWN:    "unknown",
	pb.EntityGroupStyle_ENTITY_GROUP_STYLE_FRAMED:     "framed",
	pb.EntityGroupStyle_ENTITY_GROUP_STYLE_NOT_FRAMED: "not-framed",
}

func ConvertEntityGroupStyleToPB(shape string) pb.EntityGroupStyle {
	if val, ok := entityGroupStyleToPB[shape]; ok {
		return val
	}
	return pb.EntityGroupStyle_ENTITY_GROUP_STYLE_UNKNOWN
}

func ConvertEntityGroupStyleToDB(shape pb.EntityGroupStyle) string {
	if val, ok := entityGroupStyleToDB[shape]; ok {
		return val
	}
	return "unknown"
}
