// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package sexpr

import "fmt"

type TokenType uint8

const (
	TokListOpen TokenType = iota
	TokListClose
	TokComment
	TokIdent
	TokString
	TokRawString
	TokChar
	TokNumber
	TokBoolean
	TokEof
	TokErr
)

func (tt TokenType) String() string {
	switch tt {
	case TokEof:
		return "Eof"
	case TokErr:
		return "Err"
	case TokListOpen:
		return "List"
	case TokListClose:
		return "List"
	case TokComment:
		return "Comment"
	case TokIdent:
		return "Ident"
	case TokString:
		return "String"
	case TokRawString:
		return "RawString"
	case TokChar:
		return "Char"
	case TokNumber:
		return "Number"
	case TokBoolean:
		return "Bool"
	}

	return "Unknown"
}

type Token struct {
	Data []byte
	Line int
	Col  uint16
	Type TokenType
}

func (t Token) String() string {
	if len(t.Data) > 30 {
		return fmt.Sprintf("%s(%.30q...)", t.Type, t.Data)
	}
	return fmt.Sprintf("%s(%q)", t.Type, t.Data)
}
