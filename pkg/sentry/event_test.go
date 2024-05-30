package sentry

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/williampsena/bugs-channel-plugins/pkg/test"
)

func TestSentryEventMarshalJSON(t *testing.T) {
	b := []byte(test.ReadFixtureFile(t, "../../fixtures/sentry/event.json"))

	sentryEvent, err := NewSentryEventFromJson(1, "foo", b)

	require.Nil(t, err)

	assert.Equal(t, "ccbf33b5bfd446368efe7ef2113ec78c", sentryEvent.ID)
	assert.Equal(t, 1, sentryEvent.Project)

	stackTraceJson, err := json.Marshal(sentryEvent.StackTrace)

	require.Nil(t, err)

	assert.Equal(t, "ccbf33b5bfd446368efe7ef2113ec78c", sentryEvent.ID)
	assert.Equal(t, "production", sentryEvent.Environment)
	assert.Equal(t, "python", sentryEvent.Platform)
	assert.Equal(t, "foo", sentryEvent.ServerName)
	assert.Equal(t, "2023-09-10T10:42:52.743172Z", sentryEvent.Timestamp)
	assert.Equal(t, "error", sentryEvent.Level)
	assert.Contains(t, string(stackTraceJson), `ValueError`)
}

func TestNewSentryEventsFromEnvelope(t *testing.T) {
	meta, items := readSentryEnvelope(t, "../../fixtures/sentry/sentry_envelope.txt")

	eventMeta, err := NewSentryEventMetaFromJson([]byte(meta))
	require.Nil(t, err)

	assert.Equal(t,
		&SentryEventMeta{ID: "70872dd60777409396331a7f06a38039", SentAt: "2024-05-21T00:27:21.484976Z", Trace: SentryEventTrace{TraceId: "31c10ec554f84c6cb47af91475ae5260", Environment: "production", PublicKey: "key", SampleRate: 1}},
		eventMeta)

	event, err := NewSentryEventsFromEnvelope(1, "foo", items)
	require.Nil(t, err)

	assert.Len(t, event, 1)
}

func TestNewSentryEventsFromEnvelopeDiscarded(t *testing.T) {
	_, items := readSentryEnvelope(t, "../../fixtures/sentry/sentry_envelope_discarded.txt")

	event, err := NewSentryEventsFromEnvelope(1, "bar", items)

	require.Nil(t, err)

	assert.Len(t, event, 0)
}

func readSentryEnvelope(t *testing.T, fixture string) (string, [][]string) {
	lines := test.ReadFixtureFileLines(t, fixture)
	items := [][]string{lines[1:]}

	return lines[0], items
}
