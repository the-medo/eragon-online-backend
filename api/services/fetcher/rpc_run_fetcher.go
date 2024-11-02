package fetcher

import (
	"context"
	"encoding/json"
	"github.com/the-medo/talebound-backend/api/apihelpers"
	"github.com/the-medo/talebound-backend/converters"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/util"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (server *ServiceFetcher) RunFetcher(ctx context.Context, req *pb.RunFetcherRequest) (*pb.RunFetcherResponse, error) {
	violations := validateRunFetcherRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	rsp := &pb.RunFetcherResponse{}

	fetchInterface := &apihelpers.FetchInterface{}

	if req.ModuleIds != nil {
		data, err := server.Store.GetModulesByIDs(ctx, req.ModuleIds)
		if err != nil {
			return nil, err
		}

		rsp.Modules = make([]*pb.ViewModule, len(data))
		for i, item := range data {
			rsp.Modules[i] = converters.ConvertViewModule(item)

			fetchInterface.PostIds = util.Upsert(fetchInterface.PostIds, item.DescriptionPostID)

			if item.WorldID.Valid {
				fetchInterface.WorldIds = util.Upsert(fetchInterface.WorldIds, item.WorldID.Int32)
			}

			if item.SystemID.Valid {
				fetchInterface.SystemIds = util.Upsert(fetchInterface.SystemIds, item.SystemID.Int32)
			}

			if item.CharacterID.Valid {
				fetchInterface.CharacterIds = util.Upsert(fetchInterface.CharacterIds, item.CharacterID.Int32)
			}

			if item.QuestID.Valid {
				fetchInterface.QuestIds = util.Upsert(fetchInterface.QuestIds, item.QuestID.Int32)
			}

			if item.HeaderImgID.Valid {
				fetchInterface.ImageIds = util.Upsert(fetchInterface.ImageIds, item.HeaderImgID.Int32)
			}

			if item.ThumbnailImgID.Valid {
				fetchInterface.ImageIds = util.Upsert(fetchInterface.ImageIds, item.ThumbnailImgID.Int32)
			}

			if item.AvatarImgID.Valid {
				fetchInterface.ImageIds = util.Upsert(fetchInterface.ImageIds, item.AvatarImgID.Int32)
			}

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

	if req.SystemIds != nil {
		data, err := server.Store.GetSystemsByIDs(ctx, req.SystemIds)
		if err != nil {
			return nil, err
		}

		rsp.Systems = make([]*pb.System, len(data))
		for i, item := range data {
			rsp.Systems[i] = converters.ConvertSystem(item)
		}
	}

	if req.CharacterIds != nil {
		data, err := server.Store.GetCharactersByIDs(ctx, req.CharacterIds)
		if err != nil {
			return nil, err
		}

		rsp.Characters = make([]*pb.Character, len(data))
		for i, item := range data {
			rsp.Characters[i] = converters.ConvertCharacter(item)
		}
	}

	if req.QuestIds != nil {
		data, err := server.Store.GetQuestsByIDs(ctx, req.QuestIds)
		if err != nil {
			return nil, err
		}

		rsp.Quests = make([]*pb.Quest, len(data))
		for i, item := range data {
			rsp.Quests[i] = converters.ConvertQuest(item)
		}
	}

	if req.EntityIds != nil {
		data, err := server.Store.GetEntitiesByIDs(ctx, req.EntityIds)
		if err != nil {
			return nil, err
		}

		rsp.Entities = make([]*pb.ViewEntity, len(data))
		for i, item := range data {
			rsp.Entities[i] = converters.ConvertViewEntity(item)

			fetchInterface.ModuleIds = util.Upsert(fetchInterface.ModuleIds, item.ModuleID)

			if item.PostID.Valid {
				fetchInterface.PostIds = util.Upsert(fetchInterface.PostIds, item.PostID.Int32)
			}

			if item.MapID.Valid {
				fetchInterface.MapIds = util.Upsert(fetchInterface.MapIds, item.MapID.Int32)
			}

			if item.LocationID.Valid {
				fetchInterface.LocationIds = util.Upsert(fetchInterface.LocationIds, item.LocationID.Int32)
			}

			if item.ImageID.Valid {
				fetchInterface.ImageIds = util.Upsert(fetchInterface.ImageIds, item.ImageID.Int32)
			}
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

			if item.ThumbnailImgID.Valid {
				fetchInterface.ImageIds = util.Upsert(fetchInterface.ImageIds, item.ThumbnailImgID.Int32)
			}

			fetchInterface.UserIds = util.Upsert(fetchInterface.UserIds, item.UserID)

			if item.LastUpdatedUserID.Valid {
				fetchInterface.UserIds = util.Upsert(fetchInterface.UserIds, item.LastUpdatedUserID.Int32)
			}
		}
	}

	if req.ImageIds != nil {
		data, err := server.Store.GetImagesByIDs(ctx, req.ImageIds)
		if err != nil {
			return nil, err
		}

		rsp.Images = make([]*pb.Image, len(data))
		for i, item := range data {
			rsp.Images[i] = converters.ConvertImage(item)

			fetchInterface.UserIds = util.Upsert(fetchInterface.UserIds, item.UserID)
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

			if item.ThumbnailImageID.Valid {
				fetchInterface.ImageIds = util.Upsert(fetchInterface.ImageIds, item.ThumbnailImageID.Int32)
			}
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

			if item.PostID.Valid {
				fetchInterface.PostIds = util.Upsert(fetchInterface.PostIds, item.PostID.Int32)
			}

			if item.ThumbnailImageID.Valid {
				fetchInterface.ImageIds = util.Upsert(fetchInterface.ImageIds, item.ThumbnailImageID.Int32)
			}
		}
	}

	if req.UserIds != nil {
		data, err := server.Store.GetUsersByIDs(ctx, req.UserIds)
		if err != nil {
			return nil, err
		}

		rsp.Users = make([]*pb.User, len(data))
		for i, item := range data {
			rsp.Users[i] = converters.ConvertUser(item)

			if item.IntroductionPostID.Valid {
				fetchInterface.PostIds = util.Upsert(fetchInterface.PostIds, item.IntroductionPostID.Int32)
			}
		}
	}

	fetchIdsHeader, err := json.Marshal(fetchInterface)

	md := metadata.Pairs(
		"X-Fetch-Ids", string(fetchIdsHeader),
	)
	err = grpc.SendHeader(ctx, md)
	if err != nil {
		return nil, err
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

	if req.UserIds != nil {
		violations = append(violations, validator.ValidateSliceOfIds(req.UserIds, "user_id")...)
	}

	return violations
}
