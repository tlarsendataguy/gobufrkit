package ast

import (
	"github.com/tlarsendataguy/gobufrkit/table"
	"reflect"
	"testing"
)

var (
	tableGroup table.TableGroup
	visitor    = &LintVisitor{}
)

func init() {
	var err error
	tableGroup, err = table.NewSingleTableGroup(
		"../../_definitions/tables",
		0, 0, 0, 28)
	if err != nil {
		panic(err)
	}
}

func lookup(id table.ID) table.Descriptor {
	descriptor, err := tableGroup.Lookup(id)
	if err != nil {
		panic(err)
	}
	return descriptor
}

func isLintError(err error) bool {
	return reflect.TypeOf(err) == reflect.TypeOf(&LintError{})
}

func TestLintVisitor_VisitFixedReplicationNode(t *testing.T) {
	tree := &FixedReplicationNode{BaseNode: NewBaseNode(lookup(102002))}
	tree.SetMembers([]Node{&ElementNode{BaseNode: NewBaseNode(lookup(1001))}})
	err := tree.Accept(visitor)
	if !isLintError(err) {
		t.Fatalf(`not lint error`)
	}
}

func TestLintVisitor_VisitOpAssocFieldNode(t *testing.T) {
	tree := &OpAssocFieldNode{BaseNode: NewBaseNode(lookup(204002))}
	err := tree.Accept(visitor)
	if !isLintError(err) {
		t.Fatalf(`not lint error`)
	}
}
