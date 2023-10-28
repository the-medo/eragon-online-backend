package tags

import (
	"context"
	"github.com/the-medo/talebound-backend/api/converters"
	"github.com/the-medo/talebound-backend/api/e"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *ServiceTags) GetModuleTypeAvailableTags(ctx context.Context, req *pb.GetModuleTypeAvailableTagsRequest) (*pb.GetModuleTypeAvailableTagsResponse, error) {
	violations := validateGetModuleTypeAvailableTags(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	tagRows, err := server.Store.GetModuleTypeTagsAvailable(ctx, converters.ConvertModuleTypeToDB(req.GetModuleType()))

	if err != nil {
		return nil, err
	}

	rsp := &pb.GetModuleTypeAvailableTagsResponse{
		Tags: make([]*pb.ViewTag, len(tagRows)),
	}

	for i, dbTag := range tagRows {
		rsp.Tags[i] = converters.ConvertViewTag(dbTag)
	}

	return rsp, nil
}

func validateGetModuleTypeAvailableTags(req *pb.GetModuleTypeAvailableTagsRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	return violations
}
