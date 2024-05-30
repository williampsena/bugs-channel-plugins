package service

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/williampsena/bugs-channel-plugins/pkg/settings"
)

func TestGetServiceSuccess(t *testing.T) {
	services := buildConfigFileServices(t)
	fetcher := NewServiceFetcher(services)
	service, err := fetcher.GetServiceByAuthKey("key")

	require.Nil(t, err)

	fmt.Println(service)
	assert.Equal(t, Service{Id: "1", Name: "foo bar service"}, service)
}

func TestGetServiceError(t *testing.T) {
	services := buildConfigFileServices(t)
	fetcher := NewServiceFetcher(services)
	_, err := fetcher.GetServiceByAuthKey("expiredKey")

	assert.Equal(t, ErrServiceNotFound, err)
}

func buildConfigFileServices(t *testing.T) []settings.ConfigFileService {
	configFile, err := settings.BuildConfigFile("../../fixtures/settings/config.yml")

	require.Nil(t, err)

	return configFile.Services
}
