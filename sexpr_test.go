// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package sexpr

import (
	"io/ioutil"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	s := new(Syntax)
	s.StringLit = []string{"\"", "\""}
	s.Delimiters = [][2]string{{"(", ")"}}
	s.NumberFunc = LexNumber
	s.BooleanFunc = func(l *Lexer) int {
		if ret := l.AcceptLiteral("#t"); ret != 0 {
			return ret
		}
		return l.AcceptLiteral("#f")
	}

	testFile(t, "testdata/palindrome.scm", s)
	testFile(t, "testdata/style.gss", s)
}

func testFile(t *testing.T, file string, syntax *Syntax) {
	var ast AST
	var err error

	st := time.Now()

	if err = ParseFile(&ast, file, syntax); err != nil {
		t.Error(err)
	}

	//println(ast.String())
	t.Log(file, "(file)  ", time.Now().Sub(st))

	ast = AST{} // reset
	st = time.Now()
	if data, err := ioutil.ReadFile(file); err != nil {
		t.Error(err)
	} else if err = Parse(&ast, data, syntax); err != nil {
		t.Error(err)
	}

	//println(ast.String())
	t.Log(file, "([]byte)", time.Now().Sub(st))
}
