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

// Scanner.read returns EOF, finally
func (s *ScannerSuite) TestReadReturnsEOFFinally() {
	scanner := NewScanner(bytes.NewBufferString(""))
	s.Equal(eof, scanner.read())
}

// Scanner.unread sets us back a character
func (s *ScannerSuite) TestUnread() {
	scanner := NewScanner(bytes.NewBufferString("ab"))
	first := scanner.read()
	scanner.unread()
	second := scanner.read()
	s.Equal('a', first)
	s.Equal('a', second)
}

func (s *ScannerSuite) TestScanHandlesWhiteSpeace() {
	scanner := NewScanner(bytes.NewBufferString(" \t\n "))
	tok, lit := scanner.Scan()
	s.Equal(WHITESPACE, tok)
	s.Equal(" \t\n ", lit)
}

func (s *ScannerSuite) TestScanHaldesScientific() {
	scanner := NewScanner(bytes.NewBufferString("E+00"))
	tok, lit := scanner.Scan()
	s.Equal(SCIENTIFIC, tok)
	s.Equal("E+00", lit)
}

func TestScannerSuite(t *testing.T) {
	suite.Run(t, new(ScannerSuite))
}
