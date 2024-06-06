package main

import (
	log "github.com/sirupsen/logrus"
	settings "github.com/williampsena/bugs-channel-plugins/internal/settings"
	"github.com/williampsena/bugs-channel-plugins/pkg/config"
	"github.com/williampsena/bugs-channel-plugins/pkg/logger"
	"github.com/williampsena/bugs-channel-plugins/pkg/sentry"
	"github.com/williampsena/bugs-channel-plugins/pkg/service"
)

func init() {
	logger.Setup()
}

func main() {
	configFile, err := settings.BuildConfigFile(config.ConfigFile())

	if err != nil {
		log.Fatal("❌ The configuration file is in incorrect format or does not exist.", err)
	}

	sentry.SetupServer(service.NewServiceFetcher(configFile.Services))
}
