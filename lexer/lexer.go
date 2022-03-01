package lexer

import (
	"bufio"
	"bytes"
	"io"
	"unicode"

	"github.com/momaee/WL/token"
)

// LexReader is an interface that wraps Read method.
// Read method reads and return the next rune from the input.
type LexReader interface {
	read() rune
}

// LexScanner is the interface that adds Unread method to the
// basic LexReader.
//
// Unread causes the next call to the Read method return the same
// rune as the same previous call to Read.
type LexScanner interface {
	LexReader
	unread() error
	Scan() *token.Token
}

// Scanner implements a tokenizer.
type scanner struct {
	r *bufio.Reader
}

// NewScanner returns a new instance of Scanner.
func NewScanner(r io.Reader) LexScanner {
	return &scanner{
		bufio.NewReader(r),
	}
}

// Read method reads the next rune from r.
// err != nil only if there is no more rune to read.
func (s *scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return token.EOF
	}
	return ch
}

// unread re-buffers the last read data.
func (s *scanner) unread() error {
	if err := s.r.UnreadRune(); err != nil {
		return err
	}
	return nil
}

// scanWhitespace consumes all subsequent whitespace.
func (s *scanner) scanWhitespace() *token.Token {
	var buf bytes.Buffer
	for {
		ch := s.read()
		if ch == token.EOF {
			break
		} else if !isWhitespace(ch) {
			_ = s.unread()
			break
		}
		_, _ = buf.WriteRune(ch)
	}
	return &token.Token{Tok: token.WhitespaceToken, Value: buf.String()}
}

// scanLetterDigit ignores all subsequent letters or digits.
func (s *scanner) scanLetterDigit() *token.Token {
	var buff bytes.Buffer
	for {
		ch := s.read()
		if ch == token.EOF {
			break
		} else if !isLetterDigit(ch) {
			_ = s.unread()
			break
		}
		buff.WriteRune(ch)
	}

	return &token.Token{Tok: token.IllegalToken, Value: buff.String()}
}

// Scan prepare and returns the next Token.
func (s *scanner) Scan() *token.Token {

	// read next rune
	ch := s.read()

	// If whitespace code point found, then consume all contiguous whitespaces.
	if isWhitespace(ch) {
		return s.scanWhitespace()
	}

	// If letter, digit code point found, then consume all letters, digits
	if isLetterDigit(ch) {
		return s.scanLetterDigit()
	}

	return s.next(ch)
}

func (s *scanner) next(ch rune) *token.Token {
	// Check against individual code points next.

	for symbol, tok := range token.AllTokens {
		if ch == symbol {
			return tok
		}
	}

	return &token.Token{Tok: token.IllegalToken, Value: string(ch)}
}

// isWhitespace returns True if ch is space, tab, new-line.
func isWhitespace(ch rune) bool {
	return unicode.IsSpace(ch)
}

// isLetterDigit returns True if ch is letter or digit.
func isLetterDigit(ch rune) bool {
	return unicode.IsLetter(ch) || unicode.IsDigit(ch)
}
