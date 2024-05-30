package sentry

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/williampsena/bugs-channel-plugins/pkg/test"
)

func TestTranslateEventsSuccess(t *testing.T) {
	envelope := bytes.NewReader(test.ReadFixtureFile(t, "../../fixtures/sentry/sentry_envelope.zip"))

	events, err := TranslateEvents("1", envelope)
	require.Nil(t, err)

	assert.Len(t, events, 1)
}

func TestTranslateEventsFailed(t *testing.T) {
	envelope := bytes.NewReader(test.ReadFixtureFile(t, "../../fixtures/sentry/sentry_envelope_discarded.zip"))

	events, err := TranslateEvents("1", envelope)
	require.Nil(t, err)

	assert.Len(t, events, 0)
}
