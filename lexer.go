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
	SKIP          // _

	SYMBOL // One of $-+():space

	COLOR // A color name between square brackets, e.g. [RED].
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

func isSkip(ch rune) bool {
	return ch == '_'
}

func isColorStartChar(ch rune) bool {
	return ch == '['
}

func isColorEndChar(ch rune) bool {
	return ch == ']'
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
	case isSkip(ch):
		// The skipchar itself we just throw away, so there's
		// no need to unread it.
		return s.scanSkip()
	case isWhitespace(ch):
		s.unread()
		return s.scanWhitespace()
	case isScientificStartChar(ch):
		s.unread()
		return s.scanScientific()
	case isColorStartChar(ch):
		// We'll throw aawy the start char, so no need to
		// unread.
		return s.scanColor()
	}

	// Otherwise read the individual character.
	switch ch {
	case eof:
		return EOF, ""
	case '0':
		return ZERO, string(ch)
	case ',':
		return COMMA, string(ch)
	case '.':
		return PERIOD, string(ch)
	case ';':
		return SEMICOLON, string(ch)
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

// scanSkip consumes the rune that follows a skipChar.
func (s *Scanner) scanSkip() (tok Token, lit string) {
	ch := s.read()
	if ch == eof {
		return EOF, ""
	}
	return SKIP, string(ch)
}

// scanColor consumes the name of a color and throws away the color
// termination char, ']'.
func (s *Scanner) scanColor() (tok Token, lit string) {
	var buf bytes.Buffer
	for {
		if ch := s.read(); ch == eof {
			break
		} else {
			if isColorEndChar(ch) {
				break
			}
			buf.WriteRune(ch)
		}

	}
	return COLOR, buf.String()
}
