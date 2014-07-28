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

func TestShortParse(t *testing.T) {
	s := new(Syntax)
	s.StringLit = []string{"\"", "\""}
	s.Delimiters = [][2]string{{"(", ")"}}
	s.NumberFunc = LexNumber
	s.BooleanFunc = LexBoolean

	var ast AST
	err := ParseString(&ast, "(a)", s)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	n := &ast.Root
	if len(n.Children) != 1 {
		t.Errorf("Expected one child of root node, got %d", len(n.Children))
		t.FailNow()
	}

	n = n.Children[0]
	if len(n.Children) != 1 {
		t.Errorf("Expected one child, got %d", len(n.Children))
	}

	n = n.Children[0]
	if n.Type != TokIdent {
		t.Errorf("Expected an identifier, got %s", n.Type)
	}

	if len(n.Data) != 1 || n.Data[0] != 'a' {
		t.Errorf("Expected identifier `a`, got %+q", n.Data)
	}
}
