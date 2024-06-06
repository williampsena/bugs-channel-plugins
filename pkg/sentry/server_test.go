package sentry

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	settings "github.com/williampsena/bugs-channel-plugins/internal"
	"github.com/williampsena/bugs-channel-plugins/pkg/service"
)

func TestServer(t *testing.T) {
	router := buildRouter(buildServiceFetcher(t))

	svr := httptest.NewServer(router)

	defer svr.Close()

	requestURL := fmt.Sprintf("%v/health", svr.URL)

	res, err := http.Get(requestURL)

	require.Nil(t, err)

	assert.Equal(t, 200, res.StatusCode)
}

func buildServiceFetcher(t *testing.T) service.ServiceFetcher {
	configFile, err := settings.BuildConfigFile("../../fixtures/settings/config.yml")

	require.Nil(t, err)

	return service.NewServiceFetcher(configFile.Services)
}
