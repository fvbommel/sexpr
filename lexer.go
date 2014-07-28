// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package sexpr

import (
	"bytes"
	"fmt"
	"unicode"
	"unicode/utf8"
)

const EOF = -1

// A Lexer turns s-expression source into a stream of tokens.
type Lexer struct {
	data   []byte    // Data to be parsed.
	syntax *Syntax   // Set of syntax rules
	line   [2]int    // Current line and line where token started.
	col    [2]uint16 // Current column and column where token started.
	start  int       // Start of current token.
	pos    int       // Current position in buffer.
	size   int       // Size of last read rune.
	prlnsz uint16    // Size of previous line. Needed for accurate line/col tracking when Rewinding.
}

// New creates a new lexer for the given input data.
// The meaning of tokens this lexer looks for can be configured through the
// supplied Syntax struct.
func NewLexer(data []byte, syntax *Syntax) *Lexer {
	l := new(Lexer)
	l.data = data
	l.line[0], l.line[1] = 1, 1
	l.col[0], l.col[1] = 1, 1
	l.syntax = syntax
	return l
}

func (l *Lexer) errorf(tok *Token, f string, argv ...interface{}) {
	tok.Type = TokErr
	tok.Line = l.line[1]
	tok.Col = l.col[1]
	tok.Data = []byte(fmt.Sprintf(f, argv...))
	l.Ignore()
}

func (l *Lexer) emit(tok *Token, tt TokenType) {
	tok.Type = tt
	tok.Line = l.line[1]
	tok.Col = l.col[1]
	tok.Data = l.data[l.start:l.pos]
	l.Ignore()
}

// Next returns the next token. If there are none available,
// this yields a token with Type set to TokEof.
// TokErr denotes that an error occurred.
func (l *Lexer) Next(tok *Token) {
	var ret, i int

	// Consume leading whitespace.
	if ret = l.AcceptSpace(); ret == EOF {
		l.emit(tok, TokEof)
		return
	}

	l.Ignore() // Ignore whatever is in the buffer.

	// Do we have a single-line comment?
	if len(l.syntax.SingleLineComment) > 0 {
		if ret = l.AcceptLiteral(l.syntax.SingleLineComment); ret == EOF {
			l.emit(tok, TokEof)
			return
		} else if ret == 1 {
			l.Ignore()
			l.AcceptUntil("\r\n")
			l.emit(tok, TokComment)
			return
		}
	}

	// Do we have a multi-line comment?
	if len(l.syntax.MultiLineComment) > 1 {
		if ret = l.lexPair(tok, l.syntax.MultiLineComment, TokComment, "comment"); ret == EOF {
			l.emit(tok, TokEof)
		}

		if ret != 0 {
			return
		}
	}

	// Do we have a char literal?
	if len(l.syntax.CharLit) > 1 {
		if ret = l.lexPair(tok, l.syntax.CharLit, TokChar, "char"); ret == EOF {
			l.emit(tok, TokEof)
		}

		if ret != 0 {
			return
		}
	}

	// Do we have a normal string literal?
	if len(l.syntax.StringLit) > 1 {
		if ret = l.lexPair(tok, l.syntax.StringLit, TokString, "string"); ret == EOF {
			l.emit(tok, TokEof)
		}

		if ret != 0 {
			return
		}
	}

	// Do we have a raw string literal?
	if len(l.syntax.RawStringLit) > 1 {
		if ret = l.lexPair(tok, l.syntax.RawStringLit, TokRawString, "raw string"); ret == EOF {
			l.emit(tok, TokEof)
		}

		if ret != 0 {
			return
		}
	}

	// Do we have a list delimiter?
	for i = range l.syntax.Delimiters {
		if ret = l.AcceptLiteral(l.syntax.Delimiters[i][0]); ret == EOF {
			l.emit(tok, TokEof)
			return
		} else if ret == 1 {
			l.emit(tok, TokListOpen)
			return
		}

		if ret = l.AcceptLiteral(l.syntax.Delimiters[i][1]); ret == EOF {
			l.emit(tok, TokEof)
			return
		} else if ret == 1 {
			l.emit(tok, TokListClose)
			return
		}
	}

	// See if we have a boolean literal.
	if l.syntax.BooleanFunc != nil {
		if ret = l.syntax.BooleanFunc(l); ret == EOF {
			l.emit(tok, TokEof)
			return
		} else if ret == 1 {
			l.emit(tok, TokBoolean)
			return
		}
	}

	// Or a numeric literal...
	if l.syntax.NumberFunc != nil {
		if ret = l.syntax.NumberFunc(l); ret == EOF {
			l.emit(tok, TokEof)
			return
		} else if ret == 1 {
			l.emit(tok, TokNumber)
			return
		}
	}

	// An ident then?
	if ret = l.AcceptIdent(); ret == EOF {
		l.emit(tok, TokEof)
		return
	} else if ret == 1 {
		l.emit(tok, TokIdent)
		return
	}

	// This will only occur with really unorthodox runes.'
	// They will most likely indicate an utf8 decoding error or
	// that we have been reading a binary file.
	l.errorf(tok, "Unexpected character %q", l.NextRune())
}

// NextRune retuns the nextrune unicode rune in the input.
func (l *Lexer) NextRune() (r rune) {
	if l.pos >= len(l.data) {
		return EOF
	}

	r, l.size = utf8.DecodeRune(l.data[l.pos:])
	l.pos += l.size

	if r == '\n' {
		l.line[0]++
		l.prlnsz, l.col[0] = l.col[0], 0
	}

	l.col[0]++
	return r
}

// Ignore the input so far.
func (l *Lexer) Ignore() {
	l.start = l.pos
	l.line[1] = l.line[0]
	l.col[1] = l.col[0]
}

// Rewind Rewinds to the last rune.
// Can be called only once per NextRune() call.
func (l *Lexer) Rewind() {
	l.pos -= l.size
	if l.col[0] > 1 {
		l.col[0]--
	} else {
		l.line[0]--
		l.col[0] = l.prlnsz
	}
}

// Skip Skips the NextRune character.
func (l *Lexer) Skip() {
	l.NextRune()
	l.Ignore()
}

// Accept consumes the next rune if it is contained in the supplied string.
func (l *Lexer) Accept(valid string) int {
	r := l.NextRune()

	if r == EOF {
		return EOF
	}

	if indexRune(valid, r) == -1 {
		l.Rewind()
		return 0
	}

	return 1
}

// AcceptRun consumes runes for as long they are contained in the supplied string.
// It returns the number of runes consumed or EOF.
func (l *Lexer) AcceptRun(valid string) int {
	var r rune
	var n int

	for {
		if r = l.NextRune(); r == EOF {
			return EOF
		}

		if indexRune(valid, r) == -1 {
			l.Rewind()
			break
		}

		n++
	}

	return n
}

// AcceptUntil consumes runes for as long they are NOT contained in the
// supplied string.
func (l *Lexer) AcceptUntil(valid string) int {
	var r rune

	for {
		if r = l.NextRune(); r == EOF {
			return EOF
		}

		if indexRune(valid, r) != -1 {
			l.Rewind()
			return 1
		}
	}
}

// AcceptLiteral consumes runes if they are an exact, rune-for-rune match with
// the supplied string.
func (l *Lexer) AcceptLiteral(valid string) int {
	if len(valid) == 0 || l.pos+len(valid) >= len(l.data) {
		return 0
	}

	// This is orders of magnitude faster than using bytes.Index().
	for i := range valid {
		if l.data[l.pos+i] != valid[i] {
			return 0
		}
	}

	// Update line/col info by consuming runes.
	c := utf8.RuneCount(l.data[l.pos : l.pos+len(valid)])
	for i := 0; i < c; i++ {
		l.NextRune()
	}

	return 1
}

// AcceptUntilLiteral consumes runes for as long they are not an exact,
// rune-for-rune match with the supplied string.
func (l *Lexer) AcceptUntilLiteral(valid string) int {
	idx := bytes.Index(l.data[l.pos:], []byte(valid))
	if idx == -1 {
		return 0
	}

	// Update line/col info by consuming runes.
	end := l.pos + idx
	for l.pos != end {
		l.NextRune()
	}

	return 1
}

// AcceptIdent consumes runes until it hits anything that does not
// qualify as a valid identifier, or is one of the reserved tokens in our
// syntax struct.
func (l *Lexer) AcceptIdent() int {
	var r rune

	for {
		if r = l.NextRune(); r == EOF {
			return EOF
		}

		if l.syntax.IsReserved(r) || unicode.IsSpace(r) || !unicode.IsGraphic(r) {
			l.Rewind()
			return 1
		}
	}
}

// AcceptSpace consumes runes for as long as they are whitespace.
func (l *Lexer) AcceptSpace() int {
	var r rune
	var n int

	for {
		if r = l.NextRune(); r == EOF {
			return EOF
		}

		if !unicode.IsSpace(r) {
			l.Rewind()
			break
		}

		n++
	}

	return n
}

func (l *Lexer) lexPair(tok *Token, pair []string, tt TokenType, name string) int {
	ret := l.AcceptLiteral(pair[0])
	if ret == EOF || ret == 0 {
		return ret
	}

	l.Ignore()

	if ret = l.AcceptUntilLiteral(pair[1]); ret <= 0 {
		l.errorf(tok, "Missing %s delimiter", name)
		return ret
	}

	l.emit(tok, tt)
	l.AcceptLiteral(pair[1])
	l.Ignore()
	return 1
}

func indexRune(s string, r rune) int {
	for i, tr := range s {
		if tr == r {
			return i
		}
	}
	return -1
}

// TestBoolean is a builtin function which tests if the given input might
// qualify as a boolean. This looks for literals 'true' and 'false'.
//
// Assign this function to Syntax.BooleanFunc if you want default behaviour.
func LexBoolean(l *Lexer) int {
	if ret := l.AcceptLiteral("true"); ret != 0 {
		return ret
	}

	return l.AcceptLiteral("false")
}

// TestNumber is a builtin function which tests if the given input might
// qualify as a number. This is not a guarantee, but tests for a reasonable
// likeness.
//
// This finds numbers of the following formats:
//
//   1234
//   12.34
//   -0.1234
//   +12.34
//   12e-12
//   +1E+32
//   0xff12AE (hexadecimal)
//   0b010110101 (binary)
//   0644 (octal)
//
// Assign this function to Syntax.NumberFunc if you want default behaviour.
func LexNumber(l *Lexer) (ret int) {
	const Digits = "0123456789abcdefABCDEF"

	l.Accept("+-")
	d := 10

	if ret = l.Accept("0"); ret == 1 {
		if l.Accept("xX") == 1 {
			d = len(Digits)

		} else if l.Accept("bB") == 1 {
			d = 2

		} else {
			d = 8
			l.Rewind()
		}
	}

	if ret = l.AcceptRun(Digits[:d]); ret <= 0 {
		return
	}

	// Floating point number?
	switch ret = l.Accept("."); ret {
	case EOF:
		return

	case 1:
		if ret = l.AcceptRun(Digits[:d]); ret == EOF {
			return
		}
	}

	// Do we have an exponent?
	if ret = l.Accept("eE"); ret != 1 {
		return 1
	}

	l.Accept("+-")
	return l.AcceptRun("0123456789")
}
