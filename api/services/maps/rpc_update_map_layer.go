package maps

import (
	"context"
	"database/sql"
	"github.com/the-medo/talebound-backend/converters"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *ServiceMaps) UpdateMapLayer(ctx context.Context, request *pb.UpdateMapLayerRequest) (*pb.ViewMapLayer, error) {
	violations := validateUpdateMapLayer(request)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	_, err := server.CheckEntityTypePermissions(ctx, db.EntityTypeMap, request.GetMapId(), nil)
	if err != nil {
		return nil, err
	}

	if request.Name != nil || request.ImageId != nil || request.Enabled != nil || request.Position != nil {

		if request.ImageId != nil {
			mapRow, err := server.Store.GetMapById(ctx, request.GetMapId())
			if err != nil {
				return nil, err
			}
			imageRow, err := server.Store.GetImageById(ctx, request.GetImageId())
			if err != nil {
				return nil, err
			}

			if (mapRow.Width != imageRow.Width) || (mapRow.Height != imageRow.Height) {
				mapLayers, err := server.Store.GetMapLayers(ctx, request.GetMapId())
				if err != nil {
					return nil, err
				}

				// updating image to a new image that is of different size than map is allowed only in case of a single layer
				if len(mapLayers) > 1 {
					return nil, e.InvalidArgumentError([]*errdetails.BadRequest_FieldViolation{
						{
							Field:       "image_id",
							Description: "size of map layer image must be the same size as the map",
						},
					})
				} else { //if we are updating the only layer in the map, we need to update the map size as well
					updateMapSizeArgs := db.UpdateMapParams{
						ID: request.GetMapId(),
						Width: sql.NullInt32{
							Int32: imageRow.Width,
							Valid: true,
						},
						Height: sql.NullInt32{
							Int32: imageRow.Height,
							Valid: true,
						},
					}

					_, err := server.Store.UpdateMap(ctx, updateMapSizeArgs)
					if err != nil {
						return nil, err
					}
				}
			}
		}

		argLayer := db.UpdateMapLayerParams{
			ID: request.GetLayerId(),
			Name: sql.NullString{
				String: request.GetName(),
				Valid:  request.Name != nil,
			},
			ImageID: sql.NullInt32{
				Int32: request.GetImageId(),
				Valid: request.ImageId != nil,
			},
			Enabled: sql.NullBool{
				Bool:  request.GetEnabled(),
				Valid: request.Enabled != nil,
			},
			Position: sql.NullInt32{
				Int32: request.GetPosition(),
				Valid: request.Position != nil,
			},
		}

		_, err := server.Store.UpdateMapLayer(ctx, argLayer)
		if err != nil {
			return nil, err
		}

	}

	viewMapLayer, err := server.Store.GetMapLayerByID(ctx, request.GetLayerId())
	if err != nil {
		return nil, err
	}

	rsp := converters.ConvertViewMapLayer(viewMapLayer)

	return rsp, nil
}

func validateUpdateMapLayer(req *pb.UpdateMapLayerRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidateMapId(req.GetMapId()); err != nil {
		violations = append(violations, e.FieldViolation("map_id", err))
	}

	if err := validator.ValidateUniversalId(req.GetLayerId()); err != nil {
		violations = append(violations, e.FieldViolation("layer_id", err))
	}

	if req.Name != nil {
		if err := validator.ValidateUniversalName(req.GetName()); err != nil {
			violations = append(violations, e.FieldViolation("name", err))
		}
	}

	if req.ImageId != nil {
		if err := validator.ValidateImageId(req.GetImageId()); err != nil {
			violations = append(violations, e.FieldViolation("image_id", err))
		}
	}

	return violations
}
