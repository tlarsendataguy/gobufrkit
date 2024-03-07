package tdcfio_test

import (
	"github.com/tlarsendataguy/gobufrkit/tdcfio"
	"os"
	"testing"
)

func TestBitReader(t *testing.T) {
	var (
		err error
		b   []byte
		u   uint
		i   int
		q   bool
		s   *tdcfio.Binary
	)

	f, err := os.Open("../_testdata/contrived.bufr")
	defer f.Close()
	if err != nil {
		t.Fatalf(`got error %v`, err)
	}
	r := tdcfio.NewBitReader(f)
	b, err = r.ReadBytes(4)
	if err != nil {
		t.Fatalf(`got error %v`, err)
	}
	assertEqual(t, string(b), "BUFR")

	u, err = r.ReadUint(24)
	if err != nil {
		t.Fatalf(`got error %v`, err)
	}
	assertEqual(t, u, uint(94))

	i, err = r.ReadInt(8)
	if err != nil {
		t.Fatalf(`got error %v`, err)
	}
	assertEqual(t, i, 4)

	r.ReadUint(24)
	r.ReadUint(8)
	r.ReadUint(16)
	r.ReadUint(16)
	r.ReadUint(8)

	q, err = r.ReadBool()
	if err != nil {
		t.Fatalf(`got error %v`, err)
	}
	if q {
		t.Fatalf(`value should be false`)
	}

	s, err = r.ReadBinary(7)
	if err != nil {
		t.Fatalf(`got error %v`, err)
	}
	assertEqual(t, s.String(), "0000000")
}

func TestPeekableBitReader(t *testing.T) {
	var (
		err error
		b   []byte
		u   uint
		i   int
		q   bool
		s   *tdcfio.Binary
	)

	f, err := os.Open("../_testdata/contrived.bufr")
	defer f.Close()
	if err != nil {
		t.Fatalf(`got error %v`, err)
	}
	r := tdcfio.NewPeekableBitReader(f)

	b, err = r.PeekBytes(0, 4)
	if err != nil {
		t.Fatalf(`got error %v`, err)
	}
	assertEqual(t, string(b), "BUFR")

	u, err = r.PeekUint(8, 24)
	if err != nil {
		t.Fatalf(`got error %v`, err)
	}
	assertEqual(t, u, uint(22))

	b, err = r.ReadBytes(4)
	if err != nil {
		t.Fatalf(`got error %v`, err)
	}
	assertEqual(t, string(b), "BUFR")

	u, err = r.ReadUint(24)
	if err != nil {
		t.Fatalf(`got error %v`, err)
	}
	assertEqual(t, u, uint(94))

	i, err = r.ReadInt(8)
	if err != nil {
		t.Fatalf(`got error %v`, err)
	}
	assertEqual(t, i, 4)

	r.ReadUint(24)
	r.ReadUint(8)
	r.ReadUint(16)
	r.ReadUint(16)
	r.ReadUint(8)

	u, err = r.PeekUint(4, 8)
	if err != nil {
		t.Fatalf(`got error %v`, err)
	}
	assertEqual(t, u, uint(18))

	q, err = r.ReadBool()
	if err != nil {
		t.Fatalf(`got error %v`, err)
	}
	if q {
		t.Fatalf(`value should be false`)
	}

	s, err = r.ReadBinary(7)
	if err != nil {
		t.Fatalf(`got error %v`, err)
	}
	assertEqual(t, s.String(), "0000000")

	u, err = r.PeekUint(7, 8)
	if err != nil {
		t.Fatalf(`got error %v`, err)
	}
	assertEqual(t, u, uint(2))

}
