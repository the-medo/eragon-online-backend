package util

import (
	"context"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"net/http"
	"strings"
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

func RemoveCookie(config Config, name string, httpOnly bool) http.Cookie {
	cookie := http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		HttpOnly: httpOnly,
		Domain:   config.CookieDomain,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
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

		accessToken := headers.Get("Access-Token")
		accessTokenExpiresAt := headers.Get("Access-Token-Expires-At")

		if accessToken != "" && accessTokenExpiresAt != "" {
			if accessToken != "null" {
				expiresAt, err := time.Parse(TimeLayout, accessTokenExpiresAt)
				if err != nil {
					http.Error(w, "Failed to parse access token expiry", http.StatusInternalServerError)
					return err
				}

				cookie := CreateCookie(config, "access_token", accessToken, true, expiresAt)
				headers.Add("Set-Cookie", cookie.String())
				cookie = CreateCookie(config, "access_token_present", "true", false, expiresAt)
				headers.Add("Set-Cookie", cookie.String())
			} else {
				cookie := RemoveCookie(config, "access_token", true)
				headers.Add("Set-Cookie", cookie.String())
				cookie = RemoveCookie(config, "access_token_present", false)
				headers.Add("Set-Cookie", cookie.String())
			}
		}

		refreshToken := headers.Get("Refresh-Token")
		refreshTokenExpiresAt := headers.Get("Refresh-Token-Expires-At")

		if refreshToken != "" && refreshTokenExpiresAt != "" {
			if refreshToken != "null" {
				expiresAt, err := time.Parse(TimeLayout, refreshTokenExpiresAt)
				if err != nil {
					http.Error(w, "Failed to parse refresh token expiry", http.StatusInternalServerError)
					return err
				}
				cookie := CreateCookie(config, "refresh_token", refreshToken, true, expiresAt)
				headers.Add("Set-Cookie", cookie.String())
				cookie = CreateCookie(config, "refresh_token_present", "true", false, expiresAt)
				headers.Add("Set-Cookie", cookie.String())
			} else {
				cookie := RemoveCookie(config, "refresh_token", true)
				headers.Add("Set-Cookie", cookie.String())
				cookie = RemoveCookie(config, "refresh_token_present", false)
				headers.Add("Set-Cookie", cookie.String())
			}
		}

		log.Info().Msgf("my filter: %v", resp)

		return nil
	}
}

func CookieAnnotator(ctx context.Context, req *http.Request) metadata.MD {
	cookieHeaders := make([]string, len(req.Cookies()))
	for i, cookie := range req.Cookies() {
		cookieHeaders[i] = cookie.String()
	}
	return metadata.Pairs("grpc-gateway-cookie", strings.Join(cookieHeaders, ";"))
}
