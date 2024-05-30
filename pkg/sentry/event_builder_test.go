package sentry

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/williampsena/bugs-channel-plugins/pkg/test"
)

func TestBuildSentryEventsSuccess(t *testing.T) {
	envelope := test.ReadFixtureFileString(t, "../../fixtures/sentry/sentry_envelope.txt")

	events, err := BuildSentryEvents(1, envelope)
	require.Nil(t, err)

	assert.Len(t, events, 1)
}

func TestBuildSentryEventsFailed(t *testing.T) {
	envelope := test.ReadFixtureFileString(t, "../../fixtures/sentry/sentry_envelope_discarded.txt")

	events, err := BuildSentryEvents(1, envelope)
	require.Nil(t, err)

	assert.Len(t, events, 0)
}
