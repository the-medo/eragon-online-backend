package testutils

import (
	"bytes"
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/the-medo/talebound-backend/constants"
	"github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/pb"
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/the-medo/talebound-backend/token"
	"google.golang.org/grpc/metadata"
)

const (
	AuthorizationHeader = "authorization"
	AuthorizationBearer = "bearer"
)

type MockServerTransportStream struct{}

func (m *MockServerTransportStream) Method() string {
	return "foo"
}

func (m *MockServerTransportStream) SetHeader(md metadata.MD) error {
	return nil
}

func (m *MockServerTransportStream) SendHeader(md metadata.MD) error {
	return nil
}

func (m *MockServerTransportStream) SetTrailer(md metadata.MD) error {
	return nil
}

func NewContextWithBearerToken(t *testing.T, tokenMaker token.Maker, userId int32, duration time.Duration) context.Context {
	accessToken, _, err := tokenMaker.CreateToken(userId, duration)
	require.NoError(t, err)

	bearerToken := fmt.Sprintf("%s %s", AuthorizationBearer, accessToken)
	md := metadata.MD{
		AuthorizationHeader: []string{
			bearerToken,
		},
	}

	return metadata.NewIncomingContext(context.Background(), md)
}

func NewContextWithCookie(t *testing.T, tokenMaker token.Maker, userId int32, duration time.Duration) context.Context {
	accessToken, _, err := tokenMaker.CreateToken(userId, duration)
	require.NoError(t, err)

	cookie := fmt.Sprintf("%s=%s", constants.CookieName, accessToken)
	md := metadata.MD{
		constants.GrpcCookieHeader: []string{
			cookie,
		},
	}

	return metadata.NewIncomingContext(context.Background(), md)
}

func RequireMatchUser(t *testing.T, user1 pb.User, user2 db.User) {
	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.Email, user2.Email)
	//require.Empty(t, user1.HashedPassword)
}

func RequireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)

	require.NoError(t, err)
	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.Email, gotUser.Email)
	require.Empty(t, gotUser.HashedPassword)
}
