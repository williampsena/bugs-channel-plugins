package sentry

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/williampsena/bugs-channel-plugins/pkg/service"
)

func TestAuthKeyMiddlewareSuccess(t *testing.T) {
	serviceFetcher := buildServiceFetcher(t)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	req := httptest.NewRequest(http.MethodGet, "/api/1/envelope", nil)
	req.Header.Set("x-sentry-auth", buildSentryAuthHeader("key"))
	w := httptest.NewRecorder()

	authMiddleware := AuthKeyMiddleware(serviceFetcher)(nextHandler)
	authMiddleware.ServeHTTP(w, req)

	res := w.Result()

	assert.Equal(t, res.StatusCode, http.StatusOK)
}

func TestAuthKeyMiddlewareErrorMissingHeader(t *testing.T) {
	serviceFetcher := buildServiceFetcher(t)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	req := httptest.NewRequest(http.MethodGet, "/api/1/envelope", nil)
	w := httptest.NewRecorder()

	authMiddleware := AuthKeyMiddleware(serviceFetcher)(nextHandler)
	authMiddleware.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(io.Reader(res.Body))

	require.Nil(t, err)

	assert.Equal(t, res.StatusCode, http.StatusUnauthorized)
	assert.Contains(t, string(data), service.ErrServiceNotFound.Error())
}

func TestAuthKeyMiddlewareErrorInvalidKey(t *testing.T) {
	serviceFetcher := buildServiceFetcher(t)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	req := httptest.NewRequest(http.MethodGet, "/api/1/envelope", nil)
	req.Header.Set("x-sentry-auth", buildSentryAuthHeader("foo"))
	w := httptest.NewRecorder()

	authMiddleware := AuthKeyMiddleware(serviceFetcher)(nextHandler)
	authMiddleware.ServeHTTP(w, req)

	res := w.Result()

	assert.Equal(t, res.StatusCode, http.StatusUnauthorized)
}

func buildSentryAuthHeader(key string) string {
	return fmt.Sprintf("Sentry sentry_key=%v, sentry_version=7, sentry_client=sentry.python/1.30.0", key)
}
