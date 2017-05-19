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

// Parse returns an XLNumFmt when 4 sections are present
func (s *ParserSuite) TestParse() {
	parser := NewParser(bytes.NewBufferString("#.###.00_);[RED](#,###.00);0.00;\"sales \"@"))
	xlNumFmt, err := parser.Parse()
	s.Nil(err)
	s.NotNil(xlNumFmt)
	s.IsType(&XLNumFmt{}, xlNumFmt)
	s.NotNil(xlNumFmt.Positive)
	s.NotNil(xlNumFmt.Negative)
	s.NotNil(xlNumFmt.Zero)
	s.NotNil(xlNumFmt.Text)
}

func TestParserSuite(t *testing.T) {
	suite.Run(t, new(ParserSuite))
}
