// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package sexpr

type SyntaxFunc func(*Lexer) int

// A Syntax struct contains rules on how the lexer should treat
// the characters it encounters in the source. This determines what
// tokens are generated.
type Syntax struct {
	// A set of list delimiters. These are pairs of strings denoting the
	// start and end of an S-expression.
	Delimiters [][2]string

	// This string starts a single line comment.
	// A single line runs until the end of a line.
	SingleLineComment string

	// These strings denotes what a multi-line comment looks starts with
	// and ends with.
	MultiLineComment []string

	// These strings determine how a string literal starts and ends.
	StringLit []string

	// These strings determine how a raw string literal starts and ends.
	// A raw string does not have its escape sequences parsed.
	RawStringLit []string

	// These strings determine how a char literal starts and ends.
	CharLit []string

	// This function should return whether or not the given
	// input qualifies as a boolean.
	BooleanFunc SyntaxFunc

	// This function should return whether or not the given
	// input qualifies as a number.
	NumberFunc SyntaxFunc
}

// IsReserved returns true if the given rune is contained in one of the syntax fields.
func (s *Syntax) IsReserved(r rune) bool {
	if len(s.SingleLineComment) > 0 && indexRune(s.SingleLineComment, r) != -1 {
		return true
	}

	if len(s.MultiLineComment) > 1 {
		if indexRune(s.MultiLineComment[0], r) != -1 {
			return true
		}

		if indexRune(s.MultiLineComment[1], r) != -1 {
			return true
		}
	}

	for i := range s.Delimiters {
		if indexRune(s.Delimiters[i][0], r) != -1 {
			return true
		}

		if indexRune(s.Delimiters[i][1], r) != -1 {
			return true
		}
	}

	if len(s.CharLit) > 1 {
		if indexRune(s.CharLit[0], r) != -1 {
			return true
		}

		if indexRune(s.CharLit[1], r) != -1 {
			return true
		}
	}

	if len(s.StringLit) > 1 {
		if indexRune(s.StringLit[0], r) != -1 {
			return true
		}

		if indexRune(s.StringLit[1], r) != -1 {
			return true
		}
	}

	if len(s.RawStringLit) > 1 {
		if indexRune(s.RawStringLit[0], r) != -1 {
			return true
		}

		if indexRune(s.RawStringLit[1], r) != -1 {
			return true
		}
	}

	return false
}
