package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/williampsena/bugs-channel-plugins/pkg/test"
)

func TestDispatchSuccess(t *testing.T) {
	buf := test.CaptureLog()

	dispatcher := NewLoggerDispatcher()

	err := dispatcher.Dispatch(Event{
		ID:        "foo",
		ServiceId: "bar",
		Platform:  "python",
	})

	require.Nil(t, err)

	assert.Contains(t, buf.String(), "🍔 Ingest Event: {foo bar  python")

	test.ResetCaptureLog()
}
