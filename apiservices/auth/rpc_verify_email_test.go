package auth

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/the-medo/talebound-backend/apiservices/testutils"
	mockdb "github.com/the-medo/talebound-backend/db/mock"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"github.com/the-medo/talebound-backend/util"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
	"time"
)

func TestVerifyEmailAPI(t *testing.T) {
	user, _ := testutils.RandomUser(t)
	verifyEmail := randomVerifyEmail(t, user)

	testCases := []struct {
		name          string
		req           *pb.VerifyEmailRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, res *pb.VerifyEmailResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.VerifyEmailRequest{
				EmailId:    verifyEmail.ID,
				SecretCode: verifyEmail.SecretCode,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.VerifyEmailTxParams{
					EmailId:    verifyEmail.ID,
					SecretCode: verifyEmail.SecretCode,
				}

				verifyEmailRes := db.VerifyEmailTxResult{
					User:        user,
					VerifyEmail: verifyEmail,
				}

				verifyEmailRes.User.IsEmailVerified = true
				verifyEmailRes.VerifyEmail.IsUsed = true

				store.EXPECT().
					VerifyEmailTx(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(verifyEmailRes, nil)
			},
			checkResponse: func(t *testing.T, res *pb.VerifyEmailResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				isVerified := res.GetIsVerified()
				require.Equal(t, true, isVerified)
			},
		},
		{
			name: "InvalidEmailId",
			req: &pb.VerifyEmailRequest{
				EmailId:    -1, //cannot be negative
				SecretCode: verifyEmail.SecretCode,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.VerifyEmailResponse, err error) {
				require.Error(t, err)
				require.Nil(t, res)

				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())

				violations := st.Details()[0].(*errdetails.BadRequest).GetFieldViolations()
				require.Len(t, violations, 1)
				require.Equal(t, "email_id", violations[0].GetField())
			},
		},
		{
			name: "InvalidSecretCode",
			req: &pb.VerifyEmailRequest{
				EmailId:    verifyEmail.ID, //cannot be negative
				SecretCode: "short",        // must be at least 32 characters
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.VerifyEmailResponse, err error) {
				require.Error(t, err)
				require.Nil(t, res)

				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())

				violations := st.Details()[0].(*errdetails.BadRequest).GetFieldViolations()
				require.Len(t, violations, 1)
				require.Equal(t, "secret_code", violations[0].GetField())
			},
		},
		{
			name: "InvalidEmailIdAndSecretCode",
			req: &pb.VerifyEmailRequest{
				EmailId:    -1,
				SecretCode: "short", // must be at least 32 characters
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.VerifyEmailResponse, err error) {
				require.Error(t, err)
				require.Nil(t, res)

				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())

				violations := st.Details()[0].(*errdetails.BadRequest).GetFieldViolations()
				require.Len(t, violations, 2)
				require.Equal(t, "email_id", violations[0].GetField())
				require.Equal(t, "secret_code", violations[1].GetField())
			},
		},
		{
			name: "WrongEmailId",
			req: &pb.VerifyEmailRequest{
				EmailId:    verifyEmail.ID + 1,
				SecretCode: verifyEmail.SecretCode,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					VerifyEmailTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.VerifyEmailTxResult{}, status.Error(codes.Internal, "failed to auth email"))
			},
			checkResponse: func(t *testing.T, res *pb.VerifyEmailResponse, err error) {
				println(res, err)
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Internal, st.Code())
			},
		},
		{
			name: "WrongSecretCode",
			req: &pb.VerifyEmailRequest{
				EmailId:    verifyEmail.ID,
				SecretCode: util.RandomString(32),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					VerifyEmailTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.VerifyEmailTxResult{}, status.Error(codes.Internal, "failed to auth email"))
			},
			checkResponse: func(t *testing.T, res *pb.VerifyEmailResponse, err error) {
				println(res, err)
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Internal, st.Code())
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

			res, err := server.VerifyEmail(context.Background(), tc.req)
			tc.checkResponse(t, res, err)
		})
	}

}

func randomVerifyEmail(t *testing.T, user db.User) (verifyEmail db.VerifyEmail) {
	verifyEmail = db.VerifyEmail{
		ID:         util.RandomInt(1, 1000),
		UserID:     user.ID,
		Email:      user.Email,
		SecretCode: util.RandomString(32),
		IsUsed:     false,
		CreatedAt:  time.Now(),
		ExpiredAt:  time.Now().Add(15 * time.Minute),
	}
	return
}
