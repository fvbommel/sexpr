// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package sexpr

import "fmt"

// Represents a parse error.
type ParseError struct {
	Line int
	Col  uint16
	File string
	Msg  string
}

// NewParseError creates a new parse error from the given values.
func NewParseError(file string, line int, col uint16, f string, argv ...interface{}) *ParseError {
	return &ParseError{line, col, file, fmt.Sprintf(f, argv...)}
}

// Error returns a string representation of this error.
func (e *ParseError) Error() string {
	return fmt.Sprintf("%s:%d:%d %s", e.File, e.Line, e.Col, e.Msg)
}
