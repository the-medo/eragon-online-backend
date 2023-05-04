package util

import (
	"context"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
	"net/http"
	"time"
)

func CreateCookie(config Config, name string, value string, httpOnly bool, expires ...time.Time) http.Cookie {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: httpOnly,
		Domain:   config.CookieDomain,
	}

	if len(expires) > 0 {
		cookie.Expires = expires[0]
	}

	if config.Environment == "production" {
		cookie.Domain = "talebound.net"
		cookie.Secure = true
	}

	return cookie
}

func CreateFilterTokensToCookies(config Config) func(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {
	return func(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {

		headers := w.Header()
		log.Info().Msgf("HEADERS: %v", headers)

		accessToken := w.Header().Get("Access-Token")
		accessTokenExpiresAt := w.Header().Get("Access-Token-Expires-At")

		layout := "2006-01-02 15:04:05.9999999 -0700"

		if accessToken != "" && accessTokenExpiresAt != "" {
			expiresAt, err := time.Parse(layout, accessTokenExpiresAt[:33])
			if err != nil {
				http.Error(w, "Failed to parse access token expiry", http.StatusInternalServerError)
				return err
			}
			cookie := CreateCookie(config, "access_token", accessToken, true, expiresAt)
			w.Header().Add("Set-Cookie", cookie.String())
			cookie = CreateCookie(config, "access_token_present", "true", false, expiresAt)
			w.Header().Add("Set-Cookie", cookie.String())
		}

		refreshToken := w.Header().Get("Refresh-Token")
		refreshTokenExpiresAt := w.Header().Get("Refresh-Token-Expires-At")

		if refreshToken != "" && refreshTokenExpiresAt != "" {
			expiresAt, err := time.Parse(layout, refreshTokenExpiresAt[:33])
			if err != nil {
				http.Error(w, "Failed to parse refresh token expiry", http.StatusInternalServerError)
				return err
			}
			cookie := CreateCookie(config, "refresh_token", refreshToken, true, expiresAt)
			w.Header().Add("Set-Cookie", cookie.String())
			cookie = CreateCookie(config, "refresh_token_present", "true", false, expiresAt)
			w.Header().Add("Set-Cookie", cookie.String())
		}

		log.Info().Msgf("my filter: %v", resp)

		return nil
	}
}
