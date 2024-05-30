// This package provides environment variable reading support.
package config

import (
	"errors"
	"os"
	"strconv"
)

// This occurs when an invalid server port number is provided.
var ErrInvalidPort = errors.New("the port is invalid")

const (
	// Production environment enumerator
	production = "production"
)

// Returns the current log level
func LogLevel() string {
	return os.Getenv("LOG_LEVEL")
}

// Returns the current environment
func Env() string {
	return os.Getenv("GO_ENV")
}

// Define if the current environment is production
func IsProduction() bool {
	return Env() == production
}

// Returns the listen sentry port application
func SentryPort() int {
	env := getEnv("SENTRY_PORT", "4001")

	port, err := strconv.Atoi(env)

	if err != nil {
		panic(errors.Join(ErrInvalidPort, err))
	}

	return port
}

// Returns the config file path
func ConfigFile() string {
	return os.Getenv("CONFIG_FILE")
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}
