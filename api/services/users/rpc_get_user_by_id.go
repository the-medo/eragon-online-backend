package users

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

func (server *ServiceUsers) GetUserById(ctx context.Context, req *pb.GetUserByIdRequest) (*pb.User, error) {

	violations := validateGetUserById(req)
	if violations != nil {
		return nil, e.InvalidArgumentError(violations)
	}

	user, err := server.Store.GetUserById(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	fetchInterface := &apihelpers.FetchInterface{
		PostIds:  []int32{},
		ImageIds: []int32{},
	}

	if user.IntroductionPostID.Valid {
		fetchInterface.PostIds = append(fetchInterface.PostIds, user.IntroductionPostID.Int32)
	}

	fetchIdsHeader, err := json.Marshal(fetchInterface)

	md := metadata.Pairs(
		"X-Fetch-Ids", string(fetchIdsHeader),
	)

	err = grpc.SendHeader(ctx, md)
	if err != nil {
		return nil, err
	}

	return converters.ConvertUser(user), nil
}

func validateGetUserById(req *pb.GetUserByIdRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := validator.ValidateUserId(req.GetUserId()); err != nil {
		violations = append(violations, e.FieldViolation("user_id", err))
	}

	return violations
}
