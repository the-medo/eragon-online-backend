package srv

import (
	"context"
	"fmt"
	"github.com/the-medo/talebound-backend/db/sqlc"
	"github.com/the-medo/talebound-backend/token"
	"github.com/the-medo/talebound-backend/util"
	"github.com/the-medo/talebound-backend/worker"
	"google.golang.org/grpc/metadata"
	"strings"
)

const GrpcCookieHeader = "grpc-gateway-cookie"
const CookieName = "access_token"

type ServiceCore struct {
	Config          util.Config
	Store           db.Store
	TaskDistributor worker.TaskDistributor
	TokenMaker      token.Maker
}

func NewServiceCore(config util.Config, store db.Store, taskDistributor worker.TaskDistributor, tokenMaker token.Maker) *ServiceCore {
	return &ServiceCore{
		Config:          config,
		Store:           store,
		TaskDistributor: taskDistributor,
		TokenMaker:      tokenMaker,
	}
}

func (core *ServiceCore) AuthorizeUserCookie(ctx context.Context) (*token.Payload, error) {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	cookieHeaders := md.Get(GrpcCookieHeader)
	if len(cookieHeaders) == 0 {
		return nil, fmt.Errorf("missing cookie header")
	}

	var accessToken string
	for _, cookieHeader := range cookieHeaders {
		cookies := strings.Split(cookieHeader, ";")
		for _, cookie := range cookies {
			trimmedCookie := strings.TrimSpace(cookie)
			if strings.HasPrefix(trimmedCookie, fmt.Sprintf("%s=", CookieName)) {
				accessToken = strings.TrimPrefix(trimmedCookie, fmt.Sprintf("%s=", CookieName))
				break
			}
		}
	}

	if accessToken == "" {
		return nil, fmt.Errorf("missing access token in cookies")
	}

	payload, err := core.TokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid access token: %s", err)
	}

	return payload, nil
}
