// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package sexpr

import (
	"fmt"
	"io/ioutil"
)

// ParseFile processes the given file and stores all the nodes it finds in
// the given AST instance. The parser uses the given syntax rule set to
// perform the parsing.
func ParseFile(ast *AST, file string, syntax *Syntax) (err error) {
	fileindex, new := ast.addFile(file)

	if !new {
		return fmt.Errorf("Parsing duplicate file %q", file)
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	return parseData(ast, data, syntax, fileindex)
}

// Parse processes the given data and stores all the nodes it finds in
// the given AST instance. The parser uses the given syntax rule set to
// perform the parsing.
func Parse(ast *AST, data []byte, syntax *Syntax) (err error) {
	fileindex, _ := ast.addFile("<raw data>")

	return parseData(ast, data, syntax, fileindex)
}

// ParseString processes the given data and stores all the nodes it finds in
// the given AST instance. The parser uses the given syntax rule set to
// perform the parsing.
func ParseString(ast *AST, data string, syntax *Syntax) (err error) {
	return Parse(ast, []byte(data), syntax)
}

// parseData processes the given data and stores all the nodes it finds in
// the given AST instance. The parser uses the given syntax rule set to
// perform the parsing.
func parseData(ast *AST, data []byte, syntax *Syntax, fileindex int) (err error) {
	var tok Token
	var node *Node

	lex := NewLexer(data, syntax)

	for {
		lex.Next(&tok)

		// This would lead to an infinite loop otherwise.
		if len(tok.Data) == 0 && tok.Type != TokEof && tok.Type != TokErr {
			tok.Type = TokErr
		}

		switch tok.Type {
		case TokEof:
			return

		case TokErr:
			return NewParseError(ast.Files[fileindex], tok.Line, tok.Col, string(tok.Data))

		case TokListOpen:
			n := &Node{
				File: uint8(fileindex),
				Line: tok.Line,
				Col:  tok.Col,
				Data: tok.Data,
				Type: tok.Type,
			}

			if node == nil {
				n.Parent = &ast.Root
				ast.Root.Children = append(ast.Root.Children, n)
				node = n
			} else {
				n.Parent = node
				node.Children = append(node.Children, n)
				node = n
			}

		default:
			if node == nil {
				return NewParseError(ast.Files[fileindex], tok.Line, tok.Col,
					"Unexpected token %s; expected %s", tok, TokListOpen)
			}

			if tok.Type == TokListClose {
				node = node.Parent
				break
			}

			node.Children = append(node.Children, &Node{
				File: uint8(fileindex),
				Line: tok.Line,
				Col:  tok.Col,
				Data: tok.Data,
				Type: tok.Type,
			})
		}
	}
}
