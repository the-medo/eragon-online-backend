package posts

import (
	"context"
	"github.com/the-medo/talebound-backend/api/converters"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *ServicePosts) GetPostTypes(ctx context.Context, req *emptypb.Empty) (*pb.GetPostTypesResponse, error) {

	postTypes, err := server.Store.GetPostTypes(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get posts: %v", err)
	}

	rsp := &pb.GetPostTypesResponse{
		PostTypes: make([]*pb.DataPostType, len(postTypes)),
	}

	for i, postType := range postTypes {
		rsp.PostTypes[i] = converters.ConvertPostType(postType)
	}

	return rsp, nil
}
