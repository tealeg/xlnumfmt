package xlnumfmt

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ParserSuite struct {
	suite.Suite
}

// NewParser returns a pointer to a new Parser containing a new Scanner
func (s *ParserSuite) TestNewParser() {
	parser := NewParser(bytes.NewBufferString(""))
	s.NotNil(parser)
	s.IsType(&Parser{}, parser)
	s.NotNil(parser.s)
	s.IsType(&Scanner{}, parser.s)
}

// Scan reads the next token and literal from the Scanner.
func (s *ParserSuite) TestScan() {
	parser := NewParser(bytes.NewBufferString(";"))
	tok, lit := parser.scan()
	s.Equal(SEMICOLON, tok)
	s.Equal(";", lit)

}

func TestParserSuite(t *testing.T) {
	suite.Run(t, new(ParserSuite))
}
