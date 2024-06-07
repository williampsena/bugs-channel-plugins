package main

import (
	"context"

	log "github.com/sirupsen/logrus"
	settings "github.com/williampsena/bugs-channel-plugins/internal/settings"
	"github.com/williampsena/bugs-channel-plugins/pkg/config"
	"github.com/williampsena/bugs-channel-plugins/pkg/event"
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

	serverContext := sentry.ServerContext{
		Context:          context.Background(),
		ServiceFetcher:   service.NewServiceFetcher(configFile.Services),
		EventsDispatcher: event.NewLoggerDispatcher(),
	}

	sentry.SetupServer(&serverContext)
}
