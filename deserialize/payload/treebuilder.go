package payload

import (
	"fmt"
	"github.com/tlarsendataguy/gobufrkit/bufr"
)

// TreeBuilder provides operations for building the hierarchical BUFR nodes.
type TreeBuilder struct {
	Verbose bool

	node  bufr.Node
	stack []bufr.Node
}

func NewTreeBuilder(node bufr.Node) *TreeBuilder {
	return &TreeBuilder{node: node}
}

// Add given node as a child to the current parent node
func (tb *TreeBuilder) Add(node bufr.Node) {
	tb.node.AddMember(node)
}

// Add given node as a child to the current parent node and push it as the new parent
func (tb *TreeBuilder) Push(node bufr.Node) {
	tb.stack = append(tb.stack, tb.node)
	tb.node = node
}

// Pop replace the current parent with its own parent
func (tb *TreeBuilder) Pop() {
	tb.stack, tb.node = tb.stack[:len(tb.stack)-1], tb.stack[len(tb.stack)-1]
}

func (tb *TreeBuilder) Root() (bufr.Node, error) {
	if len(tb.stack) != 0 {
		return nil, fmt.Errorf("tree builder stack is not empty when root node is required")
	}
	return tb.node, nil
}
