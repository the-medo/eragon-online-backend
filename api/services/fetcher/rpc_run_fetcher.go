package fetcher

import (
	"context"
	"github.com/the-medo/talebound-backend/converters"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *ServiceFetcher) RunFetcher(ctx context.Context, req *pb.RunFetcherRequest) (*pb.RunFetcherResponse, error) {
	violations := validateRunFetcherRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	rsp := &pb.RunFetcherResponse{}

	if req.ModuleIds != nil {
		data, err := server.Store.GetModulesByIDs(ctx, req.ModuleIds)
		if err != nil {
			return nil, err
		}

		rsp.Modules = make([]*pb.Module, len(data))
		for i, item := range data {
			rsp.Modules[i] = converters.ConvertModule(item)
		}
	}

	if req.WorldIds != nil {
		data, err := server.Store.GetWorldsByIDs(ctx, req.WorldIds)
		if err != nil {
			return nil, err
		}

		rsp.Worlds = make([]*pb.World, len(data))
		for i, item := range data {
			rsp.Worlds[i] = converters.ConvertWorld(item)
		}
	}

	if req.EntityIds != nil {
		data, err := server.Store.GetEntitiesByIDs(ctx, req.EntityIds)
		if err != nil {
			return nil, err
		}

		rsp.Entities = make([]*pb.Entity, len(data))
		for i, item := range data {
			rsp.Entities[i] = converters.ConvertEntity(item)
		}
	}

	if req.PostIds != nil {
		data, err := server.Store.GetPostsByIDs(ctx, req.PostIds)
		if err != nil {
			return nil, err
		}

		rsp.Posts = make([]*pb.Post, len(data))
		for i, item := range data {
			rsp.Posts[i] = converters.ConvertPost(item)
		}
	}

	if req.ImageIds != nil {
		data, err := server.Store.GetImagesByIDs(ctx, req.ImageIds)
		if err != nil {
			return nil, err
		}

		rsp.Images = make([]*pb.Image, len(data))
		for i, item := range data {
			rsp.Images[i] = converters.ConvertImage(&item)
		}
	}

	if req.MapIds != nil {
		data, err := server.Store.GetMapsByIDs(ctx, req.MapIds)
		if err != nil {
			return nil, err
		}

		rsp.Maps = make([]*pb.Map, len(data))
		for i, item := range data {
			rsp.Maps[i] = converters.ConvertMap(item)
		}
	}

	if req.LocationIds != nil {
		data, err := server.Store.GetLocationsByIDs(ctx, req.LocationIds)
		if err != nil {
			return nil, err
		}

		rsp.Locations = make([]*pb.Location, len(data))
		for i, item := range data {
			rsp.Locations[i] = converters.ConvertLocation(item)
		}
	}

	return rsp, nil
}

func validateRunFetcherRequest(req *pb.RunFetcherRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if req.ModuleIds != nil {
		violations = append(violations, validator.ValidateSliceOfIds(req.ModuleIds, "module_id")...)
	}

	if req.WorldIds != nil {
		violations = append(violations, validator.ValidateSliceOfIds(req.WorldIds, "world_id")...)
	}

	if req.SystemIds != nil {
		violations = append(violations, validator.ValidateSliceOfIds(req.SystemIds, "system_id")...)
	}

	if req.CharacterIds != nil {
		violations = append(violations, validator.ValidateSliceOfIds(req.CharacterIds, "character_id")...)
	}

	if req.QuestIds != nil {
		violations = append(violations, validator.ValidateSliceOfIds(req.QuestIds, "quest_id")...)
	}

	if req.EntityIds != nil {
		violations = append(violations, validator.ValidateSliceOfIds(req.EntityIds, "entity_id")...)
	}

	if req.PostIds != nil {
		violations = append(violations, validator.ValidateSliceOfIds(req.PostIds, "post_id")...)
	}

	if req.ImageIds != nil {
		violations = append(violations, validator.ValidateSliceOfIds(req.ImageIds, "image_id")...)
	}

	if req.MapIds != nil {
		violations = append(violations, validator.ValidateSliceOfIds(req.MapIds, "map_id")...)
	}

	if req.LocationIds != nil {
		violations = append(violations, validator.ValidateSliceOfIds(req.LocationIds, "location_id")...)
	}

	return violations
}
