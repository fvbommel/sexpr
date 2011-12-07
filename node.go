// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package sexpr

// An AST node
type Node struct {
	Data     []byte    // Node data.
	Children []*Node   // Optional child nodes.
	Parent   *Node     // Parent node.
	Line     int       // Line in original source file.
	Col      uint16    // Column in original source file.
	File     uint8     // Index of name for original source file.
	Type     TokenType // Type of node.
}
