package sentry

import (
	"context"
	"net/http"
	"regexp"

	"github.com/williampsena/bugs-channel-plugins/pkg/service"
)

type key int

const (
	AuthKeyContextValue key = iota
	ServiceContextValue key = iota
)

func AuthKeyMiddleware(s service.ServiceFetcher) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authKey := fetchSentryAuthKey(r)
			service, err := s.GetServiceByAuthKey(authKey)

			if err != nil {
				HandleErrors(w, err, http.StatusUnauthorized)
				return
			}

			ctx := r.Context()
			req := r.WithContext(
				context.WithValue(ctx, AuthKeyContextValue, authKey),
			).WithContext(
				context.WithValue(ctx, ServiceContextValue, service),
			)

			*r = *req

			next.ServeHTTP(w, r)
		})
	}
}

func fetchSentryAuthKey(r *http.Request) string {
	authKey := r.Header[http.CanonicalHeaderKey("x-sentry-auth")]

	if len(authKey) > 0 {
		return extractAuthKeyHeader(authKey[0])
	} else {
		return fetchSentryAuthKeyFromQuery(r)
	}
}

func extractAuthKeyHeader(headerValue string) string {
	r := regexp.MustCompile("sentry_key=(?<value>.+?),")
	matches := r.FindStringSubmatch(headerValue)

	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}

func fetchSentryAuthKeyFromQuery(r *http.Request) string {
	return r.URL.Query().Get("sentry_key")
}
