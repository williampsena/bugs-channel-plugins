package gzip

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/williampsena/bugs-channel-plugins/pkg/test"
)

func TestUnzipSuccess(t *testing.T) {
	compressed := bytes.NewReader(test.ReadFixtureFile(t, "../../fixtures/sentry/sentry_envelope.zip"))

	content, err := UnzipReader(compressed)
	require.Nil(t, err)

	assert.Contains(t, content, "70872dd60777409396331a7f06a38039")
}
