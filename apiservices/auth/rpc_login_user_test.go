package auth

import (
	"context"
	"database/sql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/the-medo/talebound-backend/apiservices/testutils"
	mockdb "github.com/the-medo/talebound-backend/db/mock"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestLoginUserAPI(t *testing.T) {
	test := testutils.MockServerTransportStream{}
	ctx := grpc.NewContextWithServerTransportStream(context.Background(), &test)
	user, password := testutils.RandomUser(t)
	viewUser := testutils.RandomViewUser(t, user)

	testCases := []struct {
		name          string
		req           *pb.LoginUserRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, res *pb.LoginUserResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.LoginUserRequest{
				Username: user.Username,
				Password: password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Eq(user.Username)).
					Times(1).
					Return(viewUser, nil)
				store.EXPECT().
					GetImageById(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Image{}, nil)
				store.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Session{}, nil)
			},
			checkResponse: func(t *testing.T, res *pb.LoginUserResponse, err error) {
				testutils.RequireMatchUser(t, *res.GetUser(), user)
			},
		},
		{
			name: "UserNotFound",
			req: &pb.LoginUserRequest{
				Username: "nonexistinguser",
				Password: password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.ViewUser{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, res *pb.LoginUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.NotFound, st.Code())
			},
		},
		{
			name: "WrongPassword",
			req: &pb.LoginUserRequest{
				Username: user.Username,
				Password: "incorrect-password", // password is not correct
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), user.Username).
					Times(1).
					Return(viewUser, nil)
			},
			checkResponse: func(t *testing.T, res *pb.LoginUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.NotFound, st.Code())
			},
		},
		{
			name: "InvalidUsername",
			req: &pb.LoginUserRequest{
				Username: "invalid username", // username cannot contain spaces
				Password: password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.LoginUserResponse, err error) {
				require.Error(t, err)
				require.Nil(t, res)

				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())

				violations := st.Details()[0].(*errdetails.BadRequest).GetFieldViolations()
				require.Len(t, violations, 1)
				require.Equal(t, "username", violations[0].GetField())
			},
		},
		{
			name: "InvalidPassword",
			req: &pb.LoginUserRequest{
				Username: user.Username,
				Password: "short", // password is too short (min 6 chars)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), user.Username).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.LoginUserResponse, err error) {
				require.Error(t, err)
				require.Nil(t, res)

				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())

				violations := st.Details()[0].(*errdetails.BadRequest).GetFieldViolations()
				require.Len(t, violations, 1)
				require.Equal(t, "password", violations[0].GetField())
			},
		},
		{
			name: "InvalidUsernameAndPassword",
			req: &pb.LoginUserRequest{
				Username: "invalid username",
				Password: "short",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUserByUsername(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.LoginUserResponse, err error) {
				require.Error(t, err)
				require.Nil(t, res)

				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())

				violations := st.Details()[0].(*errdetails.BadRequest).GetFieldViolations()
				require.Len(t, violations, 2)
				require.Equal(t, "username", violations[0].GetField())
				require.Equal(t, "password", violations[1].GetField())
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			// Initialize a new mock Store for each test case
			storeCtrl := gomock.NewController(t)
			defer storeCtrl.Finish()
			store := mockdb.NewMockStore(storeCtrl)

			tc.buildStubs(store)
			server := NewTestAuthService(t, store, nil)

			res, err := server.LoginUser(ctx, tc.req)
			tc.checkResponse(t, res, err)
		})
	}

}
