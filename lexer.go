package xlnumfmt

import (
	"bufio"
	"bytes"
	"io"
)

type Token int

var eof = rune(0)

const (
	// Implied tokens
	BAD Token = iota
	EOF
	WHITESPACE

	// Literals
	STRING

	// Control Chars
	SEMICOLON // ;

	// Format Symbols
	ZERO          // 0
	HASH          // #
	QUESTION_MARK // ?
	PERIOD        // .
	PERCENTAGE    // %
	COMMA         // ,
	SCIENTIFIC    // E-, E+, e- or e+

	OPEN_SQUARE_BRACE  // [
	CLOSE_SQUARE_BRACE // ]

)

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isScientificStartChar(ch rune) bool {
	return ch == 'e' || ch == 'E'
}

func isScientificModifier(ch rune) bool {
	return ch == '+' || ch == '-' || isDigit(ch)
}

func isDigit(ch rune) bool {
	return ch == '0'
}

// Lexical Scanner
type Scanner struct {
	r *bufio.Reader
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

func (s *Scanner) unread() {
	_ = s.r.UnreadRune()
}

// Scan returns the next token and literal value.
func (s *Scanner) Scan() (tok Token, lit string) {
	// Read the next rune.
	ch := s.read()

	// If we see whitespace then consume all contiguous whitespace.
	// If we see a letter then consume as an ident or reserved word.
	switch {
	case isWhitespace(ch):
		s.unread()
		return s.scanWhitespace()
	case isScientificStartChar(ch):
		s.unread()
		return s.scanScientific()
	}

	// Otherwise read the individual character.
	switch ch {
	case eof:
		return EOF, ""
	case ';':
		return SEMICOLON, string(ch)
	case '[':
		return OPEN_SQUARE_BRACE, string(ch)
	case ']':
		return CLOSE_SQUARE_BRACE, string(ch)
	case '#':
		return HASH, string(ch)
	}

	return BAD, string(ch)
}

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (s *Scanner) scanWhitespace() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WHITESPACE, buf.String()
}

// scanScientific consumes the current rune and its associated modifiers
func (s *Scanner) scanScientific() (tok Token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else {

			if !isScientificStartChar(ch) && !isScientificModifier(ch) {
				s.unread()
				break
			}
			buf.WriteRune(ch)
		}
	}
	return SCIENTIFIC, buf.String()
}
