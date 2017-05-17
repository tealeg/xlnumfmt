package xlnumfmt

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
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

type ScannerSuite struct {
	suite.Suite
}

// NewScanner returns a pointer to a newly allocated Scanner that
// wraps the provided io.Reader.
func (s *ScannerSuite) TestNewScanner() {
	scanner := NewScanner(bytes.NewBufferString("foo"))
	s.NotNil(scanner)
	s.IsType(Scanner{}, *scanner)
	result := make([]byte, 3, 3)
	count, err := scanner.r.Read(result)
	s.Nil(err)
	s.Equal(3, count)
	s.Equal([]byte{'f', 'o', 'o'}, result)
}

// Scanner.read can read a char
func (s *ScannerSuite) TestRead() {
	scanner := NewScanner(bytes.NewBufferString("a"))
	s.Equal('a', scanner.read())

}

func TestScannerSuite(t *testing.T) {
	suite.Run(t, new(ScannerSuite))
}
