package xlnumfmt

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// isWhitespace correctly identifies whitespace runes
func TestIsWhitespace(t *testing.T) {
	assert.True(t, isWhitespace(' '))
	assert.True(t, isWhitespace('\t'))
	assert.True(t, isWhitespace('\n'))
	assert.False(t, isWhitespace('a'))
}

// isScientificStartChar correctly identifies 'e' notation.
func TestIsScientificStartChar(t *testing.T) {
	assert.True(t, isScientificStartChar('e'))
	assert.True(t, isScientificStartChar('E'))
	assert.False(t, isScientificStartChar('f'))
	assert.False(t, isScientificStartChar('F'))
}

// isScientificModifier correctly identifies modifier parts of 'e'
// notation.
func TestIsScientificModifier(t *testing.T) {
	assert.True(t, isScientificModifier('+'))
	assert.True(t, isScientificModifier('-'))
	assert.True(t, isScientificModifier('0'))
}

// isDigit matches exactly '0', and nothing else.
func TestIsDigit(t *testing.T) {
	assert.True(t, isDigit('0'))

	// This is perculiar to the Excel NumFmt syntax - only 0 is a
	// valid digit.
	assert.False(t, isDigit('1'))

	assert.False(t, isDigit('a'))
}

// isSkip matches only the '_' character.
func TestIsSkip(t *testing.T) {
	assert.True(t, isSkip('_'))
	assert.False(t, isSkip('#'))
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

// Scanner.Scan recognises semicolons
func (s *ScannerSuite) TestScanHandlesSemicolon() {
	scanner := NewScanner(bytes.NewBufferString(";"))
	tok, lit := scanner.Scan()
	s.Equal(SEMICOLON, tok)
	s.Equal(";", lit)
}

// Scanner.Scan recognises opening square braces
func (s *ScannerSuite) TestScanHandlesOpenSquareBrace() {
	scanner := NewScanner(bytes.NewBufferString("["))
	tok, lit := scanner.Scan()
	s.Equal(OPEN_SQUARE_BRACE, tok)
	s.Equal("[", lit)
}

// Scanner.Scan recognises closing square braces
func (s *ScannerSuite) TestScanHandlesCloseSquareBrace() {
	scanner := NewScanner(bytes.NewBufferString("]"))
	tok, lit := scanner.Scan()
	s.Equal(CLOSE_SQUARE_BRACE, tok)
	s.Equal("]", lit)
}

// Scanner.Scan recognises closing square braces
func (s *ScannerSuite) TestScanHandlesHash() {
	scanner := NewScanner(bytes.NewBufferString("#"))
	tok, lit := scanner.Scan()
	s.Equal(HASH, tok)
	s.Equal("#", lit)
}

// Scanner.Scan recognises a bad character
func (s *ScannerSuite) TestScanHandlesBad() {
	scanner := NewScanner(bytes.NewBufferString("¡"))
	tok, lit := scanner.Scan()
	s.Equal(BAD, tok)
	s.Equal("¡", lit)
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

// Scan recognise Skip characters and returns the correct type.
func (s *ScannerSuite) TestScannerHandleSkip() {
	scanner := NewScanner(bytes.NewBufferString("_("))
	tok, lit := scanner.Scan()
	s.Equal(SKIP, tok)
	s.Equal("(", lit)
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
