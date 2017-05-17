package xlnumfmt

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsWhiteSpace(t *testing.T) {
	assert.True(t, isWhiteSpace(' '))
	assert.True(t, isWhiteSpace('\t'))
	assert.True(t, isWhiteSpace('\n'))
	assert.False(t, isWhiteSpace('a'))
}

func TestIsScientificStartChar(t *testing.T) {
	assert.True(t, isScientificStartChar('e'))
	assert.True(t, isScientificStartChar('E'))
	assert.False(t, isScientificStartChar('f'))
	assert.False(t, isScientificStartChar('F'))
}

func TestIsScientificModifier(t *testing.T) {
	assert.True(t, isScientificModifier('+'))
	assert.True(t, isScientificModifier('-'))
	assert.True(t, isScientificModifier('0'))
}

func TestIsDigit(t *testing.T) {
	assert.True(t, isDigit('0'))

	// This is perculiar to the Excel NumFmt syntax - only 0 is a
	// valid digit.
	assert.False(t, isDigit('1'))

	assert.False(t, isDigit('a'))
}

// NewScanner returns a pointer to a newly allocated Scanner that
// wraps the provided io.Reader.
func TestNewScanner(t *testing.T) {
	scanner := NewScanner(bytes.NewBufferString("foo"))
	assert.NotNil(t, scanner)
	assert.IsType(t, Scanner{}, *scanner)
	result := make([]byte, 3, 3)
	count, err := scanner.r.Read(result)
	assert.Nil(t, err)
	assert.Equal(t, 3, count)
	assert.Equal(t, []byte{'f', 'o', 'o'}, result)
}
