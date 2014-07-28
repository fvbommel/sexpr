// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package sexpr

import (
	"bytes"
	"fmt"
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

// addFile tries to add a file to the AST if not yet present. It returns its
// index in the AST.Files list and whether it was newly added.
func (a *AST) addFile(file string) (idx int, new bool) {
	for i := range a.Files {
		if a.Files[i] == file {
			return i, false
		}
	}

	a.Files = append(a.Files, file)
	return len(a.Files) - 1, true
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
