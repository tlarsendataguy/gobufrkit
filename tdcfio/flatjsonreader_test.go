package tdcfio_test

import (
	"github.com/tlarsendataguy/gobufrkit/tdcfio"
	"strings"
	"testing"
)

var jstr = `
[["BUFR",94,4],[22,0,1,0,0,false,"0000000",2,4,0,18,0,2016,2,18,23,0,0],
[25,"00000000",2,true,false,"000000",[301001,105002,102000,31001,8002,20011,8002,301011,20011]],
[35,"00000000",[[94,461,2,1,2,3,4,21,3,5,6,7,8,9,10,22,2016,2,18,1],[94,461,2,1,2,3,4,21,3,5,6,7,8,9,10,22,2016,2,18,1]]],
["7777"]]
`

func TestFlatJsonReader(t *testing.T) {
	var (
		err error
		b   []byte
		u   uint
		i   int
		q   bool
		s   *tdcfio.Binary
	)

	r := tdcfio.NewFlatJsonReader(strings.NewReader(jstr))
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

func TestPeekableFlatJsonReader_PeekBytes(t *testing.T) {
	var (
		err error
		b   []byte
		u   uint
		i   int
		q   bool
		s   *tdcfio.Binary
	)

	r := tdcfio.NewPeekableFlatJsonReader(strings.NewReader(jstr))
	b, err = r.PeekBytes(0, 4)
	if err != nil {
		t.Fatalf(`got error %v`, err)
	}
	assertEqual(t, string(b), "BUFR")

	u, err = r.PeekUint(3, 24)
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

	u, err = r.PeekUint(5, 8)
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
}
