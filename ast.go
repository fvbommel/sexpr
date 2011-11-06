// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package sexpr

import (
	"fmt"
	"bytes"
)

// An abstract syntax tree.
type AST struct {
	// Root node.
	Root Node

	// Name of the source files this AST was built from.
	// A single AST can be used as input for multiple parse sessions.
	// The generated data is then merged with the existing AST.
	//
	// Each node retains line/column information from the source it came from.
	// Additionally, it will have an integer index into this list of 
	// file names.
	Files []string
}

// hasFile returns false if the given file is already present in the
// ast.Files list.
func (a *AST) addFile(file string) int {
	for i := range a.Files {
		if a.Files[i] == file {
			return -1
		}
	}

	a.Files = append(a.Files, file)
	return len(a.Files) - 1
}

func (a *AST) String() string {
	var b bytes.Buffer

	fmt.Fprintf(&b, "Files:\n")
	for i := range a.Files {
		fmt.Fprintf(&b, "- %d: %q\n", i, a.Files[i])
	}

	if len(a.Root.Children) == 0 {
		return ""
	}

	fmt.Fprintf(&b, "Nodes:\n")
	for i := range a.Root.Children {
		a.printNode(a.Root.Children[i], "  ", &b)
	}

	return b.String()
}

func (a *AST) printNode(n *Node, pad string, b *bytes.Buffer) {
	fmt.Fprintf(b, "%s%d:%03d:%03d %s(%q)\n",
		pad, n.File, n.Line, n.Col, n.Type, n.Data)

	for i := range n.Children {
		a.printNode(n.Children[i], pad+"  ", b)
	}
}
