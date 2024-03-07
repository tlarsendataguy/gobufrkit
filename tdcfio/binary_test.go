package tdcfio_test

import (
	"fmt"
	"github.com/tlarsendataguy/gobufrkit/tdcfio"
	"reflect"
	"testing"
)

func TestBinary_String(t *testing.T) {
	var binary *tdcfio.Binary

	binary, _ = tdcfio.NewBinary([]byte{1}, 1)
	assertEqual(t, fmt.Sprint(binary), "1")

	binary, _ = tdcfio.NewBinary([]byte{1}, 2)
	assertEqual(t, binary.String(), "01")

	binary, _ = tdcfio.NewBinary([]byte{1}, 8)
	assertEqual(t, binary.String(), "00000001")

	binary, _ = tdcfio.NewBinary([]byte{127, 1}, 10)
	assertEqual(t, binary.String(), "0111111101")
}

func TestBinary_AtBit(t *testing.T) {
	var binary *tdcfio.Binary

	binary, _ = tdcfio.NewBinary([]byte{85}, 8)
	expected := false
	for i := 0; i < 8; i++ {
		assertEqual(t, binary.Bit(i), expected)
		expected = !expected
	}

	binary, _ = tdcfio.NewBinary([]byte{85, 85}, 16)
	expected = false
	for i := 8; i < 16; i++ {
		assertEqual(t, binary.Bit(i), expected)
		expected = !expected
	}

	binary, _ = tdcfio.NewBinary([]byte{80, 2}, 10)
	assertEqual(t, binary.Bit(8), true)
	assertEqual(t, binary.Bit(9), false)
}

func TestNewBinary(t *testing.T) {
	var err error

	_, err = tdcfio.NewBinary([]byte{1}, 1)
	if err != nil {
		t.Fatalf(`got error %v`, err)
	}

	_, err = tdcfio.NewBinary([]byte{1}, 8)
	if err != nil {
		t.Fatalf(`got error %v`, err)
	}

	_, err = tdcfio.NewBinary([]byte{1}, 9)
	if err == nil {
		t.Fatalf(`should have gotten an error but got none`)
	}

	_, err = tdcfio.NewBinary([]byte{1, 2}, 9)
	if err != nil {
		t.Fatalf(`got error %v`, err)
	}

	_, err = tdcfio.NewBinary([]byte{1, 2}, 16)
	if err != nil {
		t.Fatalf(`got error %v`, err)
	}

	_, err = tdcfio.NewBinary([]byte{1, 2}, 8)
	if err == nil {
		t.Fatalf(`should have gotten an error but got none`)
	}

	_, err = tdcfio.NewBinary([]byte{1, 2}, 17)
	if err == nil {
		t.Fatalf(`should have gotten an error but got none`)
	}
}

func TestBinary_UnmarshalJSON(t *testing.T) {
	bin := &tdcfio.Binary{}

	bin.UnmarshalJSON([]byte("\"00110000\""))

	assertEqual(t, bin.String(), "00110000")
}

func assertEqual(t *testing.T, actual, expected any) {
	if expected != actual {
		t.Fatalf(`expected %v but got %v`, expected, actual)
	}
}

func assertEqualSlices[E comparable](t *testing.T, actual, expected []E) {
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf(`slices are different`)
	}
}
