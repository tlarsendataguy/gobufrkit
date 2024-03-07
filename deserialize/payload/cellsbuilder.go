package payload

import (
	"github.com/tlarsendataguy/gobufrkit/bufr"
	"github.com/tlarsendataguy/gobufrkit/deserialize/unpack"
)

// CellsBuilder interface wraps operations for building cells list
type CellsBuilder interface {
	Add(node *bufr.ValuedNode, val interface{})

	// Get length of the cells slice
	Len() int

	// Get the cell at given index
	Cell(index int) *bufr.Cell

	// Check the equality of the last cell across all subsets
	// TODO: does not seem to be a great fit for the interface?
	LastCellEquality() bool

	Produce(payload *bufr.Payload, root bufr.Node)
}

type UncompressCellsBuilder struct {
	cells []*bufr.Cell
}

func (ucb *UncompressCellsBuilder) Add(node *bufr.ValuedNode, val interface{}) {
	node.Index = len(ucb.cells)
	cell := bufr.NewCell(node, val)
	ucb.cells = append(ucb.cells, cell)
}

func (ucb *UncompressCellsBuilder) Len() int {
	return len(ucb.cells)
}

func (ucb *UncompressCellsBuilder) Cell(index int) *bufr.Cell {
	return ucb.cells[index]
}

// This is an no-op for UncompressCellsBuilder
func (ucb *UncompressCellsBuilder) LastCellEquality() bool {
	return true
}

func (ucb *UncompressCellsBuilder) Produce(payload *bufr.Payload, root bufr.Node) {
	payload.Add(root, ucb.cells)
}

// Tissue is a slice of cells representing all cells for a single subset
type Tissue = []*bufr.Cell

type CompressedCellsBuilder struct {
	// A slice of tissues represents cells from all subsets
	tissues  []Tissue
	nsubsets int
	isubset  int
}

func (ccb *CompressedCellsBuilder) Add(node *bufr.ValuedNode, val interface{}) {
	node.Index = ccb.Len()
	v := val.(*unpack.CompressedVal)
	node.MinValue = v.MinValue
	node.NbitsDiff = v.NbitsDiff
	for i, value := range v.Values {
		cell := bufr.NewCell(node, value)
		ccb.tissues[i] = append(ccb.tissues[i], cell)
	}
}

func (ccb *CompressedCellsBuilder) Len() int {
	return len(ccb.tissues[0])
}

// Get the cell at given index of the current subset
func (ccb *CompressedCellsBuilder) Cell(index int) *bufr.Cell {
	return ccb.tissues[ccb.isubset][index]
}

func (ccb *CompressedCellsBuilder) LastCellEquality() bool {
	index := ccb.Len() - 1
	cell := ccb.tissues[0][index]
	for _, tissue := range ccb.tissues[1:] {
		if cell.Value() != tissue[index].Value() {
			return false
		}
	}
	return true
}

func (ccb *CompressedCellsBuilder) Produce(payload *bufr.Payload, root bufr.Node) {
	for _, tissue := range ccb.tissues {
		payload.Add(root, tissue)
	}
}
