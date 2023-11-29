package posts

import (
	"context"
	"encoding/json"
	"github.com/the-medo/talebound-backend/api/apihelpers"
	"github.com/the-medo/talebound-backend/converters"
	"github.com/the-medo/talebound-backend/e"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (server *ServicePosts) GetPostById(ctx context.Context, req *pb.GetPostByIdRequest) (*pb.Post, error) {
	violations := validateGetPostByIdRequest(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	post, err := server.Store.GetPostById(ctx, req.GetPostId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get post: %v", err)
	}

	fetchInterface := &apihelpers.FetchInterface{
		UserIds:  []int32{post.UserID},
		ImageIds: []int32{},
	}

	if post.ThumbnailImgID.Valid {
		fetchInterface.ImageIds = append(fetchInterface.ImageIds, post.ThumbnailImgID.Int32)
	}

	fetchIdsHeader, err := json.Marshal(fetchInterface)

	md := metadata.Pairs(
		"X-Fetch-Ids", string(fetchIdsHeader),
	)

	err = grpc.SendHeader(ctx, md)
	if err != nil {
		return nil, err
	}

	rsp := converters.ConvertPost(post)

	return rsp, nil
}

func validateGetPostByIdRequest(req *pb.GetPostByIdRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validator.ValidatePostId(req.GetPostId()); err != nil {
		violations = append(violations, e.FieldViolation("post_id", err))
	}

	return violations
}
