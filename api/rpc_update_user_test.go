package api

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "github.com/the-medo/talebound-backend/db/mock"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/token"
	"github.com/the-medo/talebound-backend/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUpdateUserAPI(t *testing.T) {
	user, _ := randomUser(t)

	newUsername := util.RandomOwner()
	newEmail := util.RandomEmail()
	invalidEmail := "invalid-email"

	testCases := []struct {
		name          string
		req           *pb.UpdateUserRequest
		buildStubs    func(store *mockdb.MockStore)
		buildContext  func(t *testing.T, tokenMaker token.Maker) context.Context
		checkResponse func(t *testing.T, res *pb.UpdateUserResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.UpdateUserRequest{
				Id:       user.ID,
				Username: &newUsername,
				Email:    &newEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateUserParams{
					ID: user.ID,
					Email: sql.NullString{
						String: newEmail,
						Valid:  true,
					},
					Username: sql.NullString{
						String: newUsername,
						Valid:  true,
					},
				}
				updatedUser := db.User{
					ID:                user.ID,
					Username:          newUsername,
					Email:             newEmail,
					ImgID:             user.ImgID,
					HashedPassword:    user.HashedPassword,
					PasswordChangedAt: user.PasswordChangedAt,
					CreatedAt:         user.CreatedAt,
					IsEmailVerified:   user.IsEmailVerified,
				}

				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(updatedUser, nil)
				store.EXPECT().
					GetImageById(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Image{}, nil)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithCookie(t, tokenMaker, user.ID, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				updatedUser := res.GetUser()
				require.Equal(t, user.ID, updatedUser.Id)
				require.Equal(t, newUsername, updatedUser.Username)
				require.Equal(t, newEmail, updatedUser.Email)
			},
		},
		{
			name: "UserNotFound",
			req: &pb.UpdateUserRequest{
				Id:    user.ID,
				Email: &newEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.User{}, sql.ErrNoRows)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithCookie(t, tokenMaker, user.ID, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.NotFound, st.Code())
			},
		},
		{
			name: "InvalidEmail",
			req: &pb.UpdateUserRequest{
				Id:    user.ID,
				Email: &invalidEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithCookie(t, tokenMaker, user.ID, time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "ExpiredToken",
			req: &pb.UpdateUserRequest{
				Id:    user.ID,
				Email: &newEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithCookie(t, tokenMaker, user.ID, -time.Minute)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())
			},
		},
		{
			name: "NoAuthorization",
			req: &pb.UpdateUserRequest{
				Id:    user.ID,
				Email: &newEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return context.Background()
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			storeCtrl := gomock.NewController(t)
			defer storeCtrl.Finish()
			store := mockdb.NewMockStore(storeCtrl)

			tc.buildStubs(store)
			server := newTestServer(t, store, nil)

			ctx := tc.buildContext(t, server.tokenMaker)
			res, err := server.UpdateUser(ctx, tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func randomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		ID:             util.RandomUserId(),
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		Email:          util.RandomEmail(),
		ImgID: sql.NullInt32{
			Int32: util.RandomImgId(),
			Valid: true,
		},
		IsEmailVerified:   false,
		PasswordChangedAt: time.Now(),
		CreatedAt:         time.Now(),
	}
	return
}

func randomViewUser(t *testing.T, user db.User) (viewUser db.ViewUser) {

	viewUser = db.ViewUser{
		ID:                user.ID,
		Username:          user.Username,
		HashedPassword:    user.HashedPassword,
		Email:             user.Email,
		ImgID:             user.ImgID,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
		IsEmailVerified:   user.IsEmailVerified,
	}
	return
}

func requireMatchUser(t *testing.T, user1 pb.User, user2 db.User) {
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Email, user2.Email)
	//require.Empty(t, user1.HashedPassword)
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)

	require.NoError(t, err)
	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.Email, gotUser.Email)
	require.Empty(t, gotUser.HashedPassword)
}
