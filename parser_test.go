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

// Parse returns an XLNumFmt with positive, negative, zero and text
// FormatSection instances when 4 sections are present in the string.
func (s *ParserSuite) TestParse4Sections() {
	parser := NewParser(bytes.NewBufferString("#,###.00_);[RED](#,###.00);0.00;\"sales \"@"))
	xlNumFmt, err := parser.Parse()
	s.Nil(err)
	s.NotNil(xlNumFmt)
	s.IsType(&XLNumFmt{}, xlNumFmt)
	s.NotNil(xlNumFmt.Positive)
	s.Equal(9, len(xlNumFmt.Positive.Parts))
	s.Equal(
		[]Part{
			Part{Tok: HASH, Lit: "#"},
			Part{Tok: COMMA, Lit: ","},
			Part{Tok: HASH, Lit: "#"},
			Part{Tok: HASH, Lit: "#"},
			Part{Tok: HASH, Lit: "#"},
			Part{Tok: PERIOD, Lit: "."},
			Part{Tok: ZERO, Lit: "0"},
			Part{Tok: ZERO, Lit: "0"},
			Part{Tok: SKIP, Lit: ")"},
		},
		xlNumFmt.Positive.Parts)
	s.NotNil(xlNumFmt.Negative)
	s.NotNil(xlNumFmt.Zero)
	s.NotNil(xlNumFmt.Text)
}

// Parse returns an XLNumFmt with positive, negative, and zero
// FormatSection instances when 2 sections are present in the string.
func (s *ParserSuite) TestParse2Sections() {
	parser := NewParser(bytes.NewBufferString("#.###.00_);[RED](#,###.00)"))
	xlNumFmt, err := parser.Parse()
	s.Nil(err)
	s.NotNil(xlNumFmt)
	s.IsType(&XLNumFmt{}, xlNumFmt)
	s.NotNil(xlNumFmt.Positive)
	s.NotNil(xlNumFmt.Negative)
	s.NotNil(xlNumFmt.Zero)
	s.Nil(xlNumFmt.Text)
}

func TestParserSuite(t *testing.T) {
	suite.Run(t, new(ParserSuite))
}
