package tdcfio_test

import (
	"bytes"
	"github.com/tlarsendataguy/gobufrkit/tdcfio"
	"os"
	"testing"
)

func TestBitWriter(t *testing.T) {

	buffer := bytes.NewBufferString("")
	w := tdcfio.BitWriter(buffer)

	w.WriteBytes([]byte("BUFR"), 4)
	w.WriteUint(94, 24)

	w.WriteInt(4, 8)
	w.WriteUint(22, 24)
	w.WriteUint(0, 8)
	w.WriteUint(1, 16)
	w.WriteUint(0, 16)
	w.WriteUint(0, 8)

	w.WriteBool(false)

	binary, err := tdcfio.NewBinaryFromString("0000000")
	if err != nil {
		t.Fatalf(`got error %v`, err)
	}
	w.WriteBinary(binary, 7)

	f, err := os.Open("../_testdata/contrived.bufr")
	defer f.Close()
	if err != nil {
		t.Fatalf(`got error %v`, err)
	}

	b := make([]byte, 18)
	n, err := f.Read(b)
	if err != nil {
		t.Fatalf(`got error %v`, err)
	}
	assertEqual(t, n, 18)

	assertEqualSlices[byte](t, b, buffer.Bytes())
}
