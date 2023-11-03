package api

import (
	"context"
	"fmt"
	"github.com/the-medo/talebound-backend/apiservices/srv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	db "github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/token"
	"github.com/the-medo/talebound-backend/util"
	"github.com/the-medo/talebound-backend/worker"
	"google.golang.org/grpc/metadata"
)

type mockServerTransportStream struct{}

func (m *mockServerTransportStream) Method() string {
	return "foo"
}

func (m *mockServerTransportStream) SetHeader(md metadata.MD) error {
	return nil
}

func (m *mockServerTransportStream) SendHeader(md metadata.MD) error {
	return nil
}

func (m *mockServerTransportStream) SetTrailer(md metadata.MD) error {
	return nil
}

func NewTestServer(t *testing.T, store db.Store, taskDistributor worker.TaskDistributor) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store, taskDistributor)
	require.NoError(t, err)

	return server
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

	cookie := fmt.Sprintf("%s=%s", srv.CookieName, accessToken)
	md := metadata.MD{
		srv.GrpcCookieHeader: []string{
			cookie,
		},
	}

	return metadata.NewIncomingContext(context.Background(), md)
}
