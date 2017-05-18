package xlnumfmt

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestIsWhitespace(t *testing.T) {
	assert.True(t, isWhitespace(' '))
	assert.True(t, isWhitespace('\t'))
	assert.True(t, isWhitespace('\n'))
	assert.False(t, isWhitespace('a'))
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

// Scan recognises whitespace and returns the correct type.
func (s *ScannerSuite) TestScanHandlesWhiteSpeace() {
	scanner := NewScanner(bytes.NewBufferString(" \t\n "))
	tok, lit := scanner.Scan()
	s.Equal(WHITESPACE, tok)
	s.Equal(" \t\n ", lit)
}

// Scan recognises Scientific notation and returns the correct type.
func (s *ScannerSuite) TestScanHandlesScientific() {
	scanner := NewScanner(bytes.NewBufferString("E+00"))
	tok, lit := scanner.Scan()
	s.Equal(SCIENTIFIC, tok)
	s.Equal("E+00", lit)
}

// scanScientific returns the full scientific part
func (s *ScannerSuite) TestScanScientific() {
	scanner := NewScanner(bytes.NewBufferString("e-0"))
	tok, lit := scanner.scanScientific()
	s.Equal(SCIENTIFIC, tok)
	s.Equal("e-0", lit)
}

// scanScientific terminates its output when it his a non-scientific
// character.
func (s *ScannerSuite) TestScanScientificTerminatesOnNonScientific() {
	scanner := NewScanner(bytes.NewBufferString("E-00;"))
	tok, lit := scanner.scanScientific()
	s.Equal(SCIENTIFIC, tok)
	s.Equal("E-00", lit)
}

// scanScientific terminates its output on EOF
func (s *ScannerSuite) TestScanScientificStopsOnEOF() {
	scanner := NewScanner(bytes.NewBufferString("E"))
	tok, lit := scanner.scanScientific()
	s.Equal(SCIENTIFIC, tok)
	s.Equal("E", lit)
	tok, lit = scanner.Scan()
	s.Equal(EOF, tok)
}

// scanWhitespace returns the entire contiguous block of whitespace
func (s *ScannerSuite) TestScanWhitespace() {
	scanner := NewScanner(bytes.NewBufferString(" \t\n\t E"))
	tok, lit := scanner.scanWhitespace()
	s.Equal(WHITESPACE, tok)
	s.Equal(" \t\n\t ", lit)
}

// scanWhitespace stops at EOF
func (s *ScannerSuite) TestScanWhitespaceStopsAtEOF() {
	scanner := NewScanner(bytes.NewBufferString(" "))
	tok, lit := scanner.scanWhitespace()
	s.Equal(WHITESPACE, tok)
	s.Equal(" ", lit)
	tok, lit = scanner.Scan()
	s.Equal(EOF, tok)
}

func TestScannerSuite(t *testing.T) {
	suite.Run(t, new(ScannerSuite))
}
