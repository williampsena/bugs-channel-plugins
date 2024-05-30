package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/williampsena/bugs-channel-plugins/pkg/config"
	"github.com/williampsena/bugs-channel-plugins/pkg/logger"
	"github.com/williampsena/bugs-channel-plugins/pkg/sentry"
	"github.com/williampsena/bugs-channel-plugins/pkg/service"
	"github.com/williampsena/bugs-channel-plugins/pkg/settings"
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
