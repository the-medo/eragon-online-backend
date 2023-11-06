package tags

import (
	"context"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *ServiceTags) GetModuleEntityAvailableTags(ctx context.Context, req *pb.GetModuleEntityAvailableTagsRequest) (*pb.GetModuleEntityAvailableTagsResponse, error) {
	violations := validateGetModuleEntityAvailableTags(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	tagRows, err := server.Store.GetModuleEntityTagsAvailable(ctx, req.GetModuleId())

	if err != nil {
		return nil, err
	}

	rsp := &pb.GetModuleEntityAvailableTagsResponse{
		Tags: make([]*pb.Tag, len(tagRows)),
	}

	for i, dbTag := range tagRows {
		rsp.Tags[i] = &pb.Tag{
			Id:  dbTag.ID,
			Tag: dbTag.Tag,
		}
	}

	return rsp, nil
}

func validateGetModuleEntityAvailableTags(req *pb.GetModuleEntityAvailableTagsRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	return violations
}
