package table

import (
	"slices"
	"testing"
)

func TestLoadTableB(t *testing.T) {
	b, err := LoadTableB("../_definitions/tables/0/0/0/25/TableB.csv")
	if err != nil {
		t.Fatalf(`got error %v`, err)
	}

	descriptor, err := b.Lookup(ID(1001))
	if err != nil {
		t.Fatalf(`got error %v`, err)
	}
	assertEqual(t, descriptor.Id(), ID(1001))
	assertEqual(t, descriptor.Entry().Name(), "WMO BLOCK NUMBER")
}

func TestLoadTableD(t *testing.T) {
	d, err := LoadTableD("../_definitions/tables/0/0/0/25/TableD.csv")
	if err != nil {
		t.Fatalf(`got error %v`, err)
	}

	descriptor, err := d.Lookup(ID(301001))
	if err != nil {
		t.Fatalf(`got error %v`, err)
	}

	assertEqual(t, descriptor.Id(), ID(301001))
	assertEqual(t, descriptor.Entry().Name(), "(WMO block and station numbers)")

	assertEqualSlices[ID](t, descriptor.Entry().(*Dentry).Members, []ID{ID(1001), ID(1002)})

}

func assertEqual(t *testing.T, actual, expected any) {
	if expected != actual {
		t.Fatalf(`expected %v but got %v`, expected, actual)
	}
}

func assertEqualSlices[E comparable](t *testing.T, actual, expected []E) {
	if !slices.Equal(actual, expected) {
		t.Fatalf(`slices are not equal`)
	}
}
