package xlnumfmt

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestIsWhiteSpace(t *testing.T) {
	assert.True(t, isWhiteSpace(' '))
	assert.True(t, isWhiteSpace('\t'))
	assert.True(t, isWhiteSpace('\n'))
	assert.False(t, isWhiteSpace('a'))
}
