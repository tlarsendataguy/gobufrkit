package bufr

import "github.com/tlarsendataguy/gobufrkit/table"

type PackingInfo struct {
	Unit   table.Unit
	Scale  int
	Refval float64
	Nbits  int
}
