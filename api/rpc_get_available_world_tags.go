package api

import (
	"context"
	"github.com/the-medo/talebound-backend/api/converters"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *Server) GetAvailableWorldTags(ctx context.Context, req *emptypb.Empty) (*pb.GetAvailableWorldTagsResponse, error) {

	tags, err := server.store.GetWorldTagsAvailable(ctx)
	if err != nil {
		return nil, err
	}

	rsp := &pb.GetAvailableWorldTagsResponse{
		Tags: make([]*pb.Tag, len(tags)),
	}

	for i, dbTag := range tags {
		rsp.Tags[i] = converters.ConvertTag(dbTag)
	}

	return rsp, nil
}
