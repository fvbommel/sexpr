// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package sexpr

import (
	"fmt"
	"io/ioutil"
)

// Parse processes the given file and stores all the nodes it finds in
// the given AST instance. The parser uses the given syntax rule set to
// perform the parsing.
func Parse(ast *AST, file string, syntax *Syntax) (err error) {
	fileindex := ast.addFile(file)

	if fileindex == -1 {
		return fmt.Errorf("Parsing duplicate file %q", file)
	}

	var tok Token
	var node *Node
	var data []byte

	if data, err = ioutil.ReadFile(file); err != nil {
		return
	}

	lex := NewLexer(data, syntax)

	for {
		lex.Next(&tok)

		switch tok.Type {
		case TokEof:
			return

		case TokErr:
			return NewParseError(file, tok.Line, tok.Col, string(tok.Data))

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
				return NewParseError(file, tok.Line, tok.Col,
					"Unexpected token %s; expected %s", tok, TokListOpen)
			}

			if tok.Type == TokListClose {
				node = node.Parent
				break
			}

			node.Children = append(node.Children, &Node{
				File:  uint8(fileindex),
				Line: tok.Line,
				Col:  tok.Col,
				Data: tok.Data,
				Type: tok.Type,
			})
		}
	}

	return
}
