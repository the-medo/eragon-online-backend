package converters

import (
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
)

var pinShapeToPB = map[db.PinShape]pb.PinShape{
	db.PinShapeNone:     pb.PinShape_NONE,
	db.PinShapeSquare:   pb.PinShape_SQUARE,
	db.PinShapeTriangle: pb.PinShape_TRIANGLE,
	db.PinShapePin:      pb.PinShape_PIN,
	db.PinShapeCircle:   pb.PinShape_CIRCLE,
	db.PinShapeHexagon:  pb.PinShape_HEXAGON,
	db.PinShapeOctagon:  pb.PinShape_OCTAGON,
	db.PinShapeStar:     pb.PinShape_STAR,
	db.PinShapeDiamond:  pb.PinShape_DIAMOND,
	db.PinShapePentagon: pb.PinShape_PENTAGON,
	db.PinShapeHeart:    pb.PinShape_HEART,
	db.PinShapeCloud:    pb.PinShape_CLOUD,
}

var pinShapeToDB = map[pb.PinShape]db.PinShape{
	pb.PinShape_NONE:     db.PinShapeNone,
	pb.PinShape_SQUARE:   db.PinShapeSquare,
	pb.PinShape_TRIANGLE: db.PinShapeTriangle,
	pb.PinShape_PIN:      db.PinShapePin,
	pb.PinShape_CIRCLE:   db.PinShapeCircle,
	pb.PinShape_HEXAGON:  db.PinShapeHexagon,
	pb.PinShape_OCTAGON:  db.PinShapeOctagon,
	pb.PinShape_STAR:     db.PinShapeStar,
	pb.PinShape_DIAMOND:  db.PinShapeDiamond,
	pb.PinShape_PENTAGON: db.PinShapePentagon,
	pb.PinShape_HEART:    db.PinShapeHeart,
	pb.PinShape_CLOUD:    db.PinShapeCloud,
}

func ConvertPinShapeToPB(shape db.PinShape) pb.PinShape {
	if val, ok := pinShapeToPB[shape]; ok {
		return val
	}
	return pb.PinShape_NONE
}

func ConvertPinShapeToDB(shape pb.PinShape) db.PinShape {
	if val, ok := pinShapeToDB[shape]; ok {
		return val
	}
	return db.PinShapeNone
}
