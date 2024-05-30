package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogLevel(t *testing.T) {
	t.Setenv("LOG_LEVEL", "debug")

	require.Equal(t, LogLevel(), "debug")
}

func TestEnvironment(t *testing.T) {
	t.Setenv("GO_ENV", "development")

	require.Equal(t, Env(), "development")
}

func TestIsProduction(t *testing.T) {
	t.Setenv("GO_ENV", "production")
	require.Equal(t, IsProduction(), true)

	t.Setenv("GO_ENV", "test")
	require.Equal(t, IsProduction(), false)
}

func TestSentryPort(t *testing.T) {
	t.Setenv("SENTRY_PORT", "1000")
	require.Equal(t, SentryPort(), 1000)
}

func TestConfigFille(t *testing.T) {
	t.Setenv("CONFIG_FILE", "/tmp/config.yml")
	require.Equal(t, ConfigFile(), "/tmp/config.yml")
}
