package api

import (
	"context"
	"fmt"
	"github.com/the-medo/talebound-backend/token"
	"google.golang.org/grpc/metadata"
	"strings"
)

const (
	authorizationHeader = "authorization"
	authorizationBearer = "bearer"
)

func (server *Server) authorizeUser(ctx context.Context) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	values := md.Get(authorizationHeader)
	if len(values) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}

	authHeader := values[0]
	// <authorization-type> <authorization-data>
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	authType := strings.ToLower(fields[0])
	if authType != authorizationBearer {
		return nil, fmt.Errorf("unsupported authorization type: %s", authType)
	}

	accessToken := fields[1]
	payload, err := server.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid access token: %s", err)
	}

	return payload, nil
}

const (
	cookieName = "access_token"
)

func (server *Server) authorizeUserCookie(ctx context.Context) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	cookieHeaders := md.Get("cookie")
	if len(cookieHeaders) == 0 {
		return nil, fmt.Errorf("missing cookie header")
	}

	var accessToken string
	for _, cookieHeader := range cookieHeaders {
		cookies := strings.Split(cookieHeader, ";")
		for _, cookie := range cookies {
			trimmedCookie := strings.TrimSpace(cookie)
			if strings.HasPrefix(trimmedCookie, fmt.Sprintf("%s=", cookieName)) {
				accessToken = strings.TrimPrefix(trimmedCookie, fmt.Sprintf("%s=", cookieName))
				break
			}
		}
	}

	if accessToken == "" {
		return nil, fmt.Errorf("missing access token in cookies")
	}

	payload, err := server.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid access token: %s", err)
	}

	return payload, nil
}
