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
		modules, err := server.Store.GetModulesByIDs(ctx, req.ModuleIds)
		if err != nil {
			return nil, err
		}

		rsp.Modules = make([]*pb.Module, len(modules))
		for i, module := range modules {
			rsp.Modules[i] = converters.ConvertModule(module)
		}
	}

	if req.WorldIds != nil {
		worlds, err := server.Store.GetWorldsByIDs(ctx, req.WorldIds)
		if err != nil {
			return nil, err
		}

		rsp.Worlds = make([]*pb.World, len(worlds))
		for i, world := range worlds {
			rsp.Worlds[i] = converters.ConvertWorld(world)
		}
	}

	return rsp, nil
}

func validateRunFetcherRequest(req *pb.RunFetcherRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if req.ModuleIds != nil {
		for _, id := range req.ModuleIds {
			if err := validator.ValidateUniversalId(id); err != nil {
				violations = append(violations, e.FieldViolation("module_id", err))
			}
		}
	}

	if req.WorldIds != nil {
		for _, id := range req.WorldIds {
			if err := validator.ValidateUniversalId(id); err != nil {
				violations = append(violations, e.FieldViolation("world_id", err))
			}
		}
	}

	return violations
}
