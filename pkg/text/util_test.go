package text

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTruncate(t *testing.T) {
	content := "BugsChannel"

	assert.Equal(t, "BugsChannel", Truncate(content, 11))
	assert.Equal(t, "Bugs", Truncate(content, 4))
	assert.Equal(t, "BugsChannel", content)
}
