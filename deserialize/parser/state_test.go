package parser

import (
	"github.com/tlarsendataguy/gobufrkit/table"
	"testing"
)

func getKeeper() *idsKeeper {
	return newIdsKeeper([]table.ID{
		236000,
		101000,
		31002,
		31031,
		1031,
		1032,
		101000,
		31002,
		33007,
	})
}

func TestIdsKeeper_TakeWhile(t *testing.T) {

	ids := getKeeper().takeWhile(func(id table.ID) bool {
		return id == table.ID(31031)
	})

	assertEqual(t, 4, len(ids))
	assertEqual(t, table.ID(236000), ids[0])
	assertEqual(t, table.ID(31031), ids[3])
}

func TestIdsKeeper_TakeTill(t *testing.T) {

	ids := getKeeper().takeTill(func(id table.ID) bool {
		return id == table.ID(1032)
	})

	assertEqual(t, 5, len(ids))
	assertEqual(t, table.ID(236000), ids[0])
	assertEqual(t, table.ID(1031), ids[4])
}

func TestParser_State(t *testing.T) {
	parser := NewParser(nil)

	if parser.getState(stateOpDataNotPresent) {
		t.Fatalf(`state should be false`)
	}
	parser.setState(stateOpDataNotPresent)
	if !parser.getState(stateOpDataNotPresent) {
		t.Fatalf(`state should be true`)
	}
	parser.unsetState(stateOpDataNotPresent)
	if parser.getState(stateOpDataNotPresent) {
		t.Fatalf(`state should be false`)
	}
}

func assertEqual(t *testing.T, expected, actual any) {
	if expected != actual {
		t.Fatalf(`expected %v but got %v`, expected, actual)
	}
}
