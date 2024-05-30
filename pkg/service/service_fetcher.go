package service

import (
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/williampsena/bugs-channel-plugins/pkg/settings"
)

// Represents an error when service is not found
var ErrServiceNotFound = errors.New("an error occurred when attempting to fetch the service")

// The service entity
type Service struct {
	// The service identifier
	Id string

	// The name of service
	Name string
}

// Service Fetcher contracts
type ServiceFetcher interface {
	// Fetches a service by unique identifier
	GetServiceByAuthKey(string) (Service, error)
}

type serviceFetcher struct {
	services []settings.ConfigFileService
}

func (s *serviceFetcher) GetServiceByAuthKey(authKey string) (Service, error) {
	if authKey == "" {
		return Service{}, ErrServiceNotFound
	}

	for _, s := range s.services {
		for _, a := range s.AuthKeys {
			if a.Key == authKey && !a.Disabled && !isAuthKeyExpired(a.ExpiredAt) {
				return Service{Id: s.Id, Name: s.Name}, nil
			}
		}
	}

	log.Debugf("AuthKey: %v", authKey)

	return Service{}, ErrServiceNotFound
}

func isAuthKeyExpired(expiredAt int64) bool {
	if expiredAt == 0 {
		return false
	}

	return expiredAt < time.Now().Unix()
}

// Build a new service fetcher instance
func NewServiceFetcher(services []settings.ConfigFileService) ServiceFetcher {
	return &serviceFetcher{services}
}
