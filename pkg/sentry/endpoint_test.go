package sentry

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/williampsena/bugs-channel-plugins/pkg/event"
	"github.com/williampsena/bugs-channel-plugins/pkg/test"
)

func TestHealthCheckEndPoint(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	HealthCheckEndpoint(w, req)

	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(io.Reader(res.Body))

	require.Nil(t, err)

	assert.Contains(t, string(data), "Keep calm I'm absolutely alive 🐛")
}

func TestPostEventEndpointSuccess(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/1/envelope",
		bytes.NewReader(test.ReadFixtureFile(t, "../../fixtures/sentry/sentry_envelope.zip")),
	)

	w := httptest.NewRecorder()

	PostEventEndpoint(event.NewLoggerDispatcher())(w, req)

	res := w.Result()
	defer res.Body.Close()
	data, err := io.ReadAll(io.Reader(res.Body))

	require.Nil(t, err)

	assert.Equal(t, res.StatusCode, http.StatusNoContent)
	assert.Contains(t, string(data), "")
}

func TestPostEventEndpointError(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodPost,
		"/api/1/envelope",
		nil,
	)

	w := httptest.NewRecorder()

	PostEventEndpoint(event.NewLoggerDispatcher())(w, req)

	res := w.Result()

	assert.Equal(t, res.StatusCode, http.StatusUnprocessableEntity)
}
