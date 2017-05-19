package xlnumfmt

import (
	"io"
	"fmt"
)

type Parser struct {
	s   *Scanner
	buf struct {
		tok Token
		lit string
		n   int
	}
}

func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() (tok Token, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	// Otherwise read the next token from the scanner.
	tok, lit = p.s.Scan()

	// Save it to the buffer in case we unscan later.
	p.buf.tok, p.buf.lit = tok, lit

	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() {
	p.buf.n = 1
}

type Part struct {
	Tok Token
	Lit string
}

type FormatSection struct {
	Parts []Part 
}

func NewFormatSection() *FormatSection {
	return &FormatSection{Parts: make([]Part, 10)}
}

type XLNumFmt struct {
	Positive *FormatSection
	Negative *FormatSection
	Zero     *FormatSection
	Text     *FormatSection
}

func (p *Parser) Parse() (*XLNumFmt, error) {
	var sections = make([]*FormatSection, 1, 4)
	var section = NewFormatSection()
	for {
		tok, lit := p.scan()
		if tok == SEMICOLON || tok == EOF {
			if len(section.Parts) > 0 {
				sections = append(sections, section)
			}
			if tok == EOF {
				break
			}
			section = NewFormatSection()			
		}
		part := Part{Tok: tok, Lit: lit}
		section.Parts = append(section.Parts, part)
	}
	switch len(sections) {
	case 0:
		return nil, fmt.Errorf("No sections found")
	case 1:
		// If only one section is specified, it is used for all numbers.
		numFmt := &XLNumFmt{
			Positive: sections[0],
			Negative: sections[0],
			Zero: sections[0],
		}
		return numFmt, nil
	case 2:
		// If only two sections are specified, the first is
		// used for positive numbers and zeros, and the second
		// is used for negative numbers.
		numFmt := &XLNumFmt{
			Positive: sections[0],
			Negative: sections[1],
			Zero: sections[0],
		}
		return numFmt, nil
	case 4:
		// The format codes, separated by semicolons, define
		// the formats for positive numbers, negative numbers,
		// zero values, and text, in that order.
		numFmt := &XLNumFmt{
			Positive: sections[0],
			Negative: sections[1],
			Zero: sections[2],
			Text: sections[3],
		}
		return numFmt, nil
	}
	return nil, fmt.Errorf("An Excel number format must have 1, 2 or 4 semicolon separated sections.")
}
