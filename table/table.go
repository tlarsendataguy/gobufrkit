//go:generate stringer -type=Unit
package table

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// TODO: move this to bufr
type Unit int

const (
	NUMERIC Unit = iota
	STRING
	NONNEG_CODE // non-negative code, this covers all the descriptors of code unit
	CODE        // full range code
	FLAG
	BINARY
)

// Convert the unit description in table records to Normalised Unit.
func unitOf(s string) (unit Unit) {
	switch {
	case s == "CCITT IA5" || s == "Character":
		unit = STRING
	case s == "FLAG TABLE":
		unit = FLAG
	case s == "CODE TABLE" || strings.HasPrefix(s, "Common CODE TABLE"):
		unit = NONNEG_CODE
	default:
		unit = NUMERIC
	}
	return
}

// Entry is a basic interface representing additional information about a descriptor
// that is available in tables. It is useful to provide a common interface to different
// meanings that different type of descriptors may have.
type Entry interface {
	Name() string
}

// Bentry is the information about an Elemental descriptor that can be found in table B.
type Bentry struct {
	name string

	UnitString string
	Unit       Unit
	Scale      int
	Refval     int
	Nbits      int

	CrexUnitString string
	CrexUnit       Unit
	CrexScale      int
	CrexNchars     int
}

func (e *Bentry) Name() string {
	return e.name
}

// Rentry is a placeholder entry for replication descriptor
type Rentry struct {
	name string
}

func (e *Rentry) Name() string {
	return e.name
}

// Centry represents the entry for operator descriptor
type Centry struct {
	name string
}

func (e *Centry) Name() string {
	return e.name
}

// Dentry represents the information about Sequence Descriptor that is
// available in Table D.
type Dentry struct {
	name string

	// The expanded member IDs of the sequence descriptor
	Members []ID
}

func (e *Dentry) Name() string {
	return e.name
}

// Table is the basic interface representing a BUFR table.
type Table interface {
	// Get back a descriptor given its ID.
	Lookup(id ID) (Descriptor, error)
}

// B represents a BUFR Table B deserialised from an input CSV file.
type B struct {
	// Path to the input file.
	path string

	// Each non-comment line in the input file is converted to one Bentry
	// with the key being the corresponding ID.
	entries map[ID]*Bentry
}

func (b *B) Lookup(id ID) (Descriptor, error) {
	entry, ok := b.entries[id]
	if ok {
		return NewElementDescriptor(id, entry), nil
	}
	return nil, fmt.Errorf("ID not found: %s", id)
}

// LoadTableB build a Table B by reading the given input file.
func LoadTableB(tablePath string) (*B, error) {
	ins, err := os.Open(tablePath)
	defer ins.Close()
	if err != nil {
		return nil, err
	}

	records := map[string][]any{}
	r := json.NewDecoder(ins)
	err = r.Decode(&records)
	if err != nil {
		return nil, fmt.Errorf(`error loading table B from %v: %w`, tablePath, err)
	}

	entries := make(map[ID]*Bentry, len(records))

	for key, record := range records {
		id, err := strconv.Atoi(key)
		if err != nil {
			return nil, fmt.Errorf(`error loading table B from %v: %w`, tablePath, err)
		}
		entry, err := recordToBentry(record)
		if err != nil {
			return nil, fmt.Errorf(`error loading table B from %v: %w`, tablePath, err)
		}
		entries[ID(id)] = entry
	}
	return &B{path: tablePath, entries: entries}, nil
}

// D represents a single BUFR Table D deserialised from an input CSV file.
type D struct {
	// Path to the input file
	path string

	// Each non-commment line from the input file is converted to a Dentry
	// indexed by the corresponding ID.
	entries map[ID]*Dentry
}

func (d *D) Lookup(id ID) (Descriptor, error) {
	entry, ok := d.entries[id]
	if ok {
		return NewSequenceDescriptor(id, entry), nil
	}
	return nil, fmt.Errorf("ID not found: %s", id)
}

// Build a Table D from the given input file.
func LoadTableD(tablePath string) (*D, error) {
	ins, err := os.Open(tablePath)
	defer ins.Close()
	if err != nil {
		return nil, err
	}

	r := json.NewDecoder(ins)
	records := map[string][]any{}
	err = r.Decode(&records)
	if err != nil {
		return nil, fmt.Errorf(`error loading table D from %v: %w`, tablePath, err)
	}

	entries := make(map[ID]*Dentry, len(records))

	for key, record := range records {
		id, err := strconv.Atoi(key)
		if err != nil {
			return nil, fmt.Errorf(`error loading table D from %v: %w`, tablePath, err)
		}
		entry, err := recordToDentry(record)
		if err != nil {
			return nil, fmt.Errorf(`error loading table D from %v: %w`, tablePath, err)
		}
		entries[ID(id)] = entry
	}
	return &D{path: tablePath, entries: entries}, nil
}

// recordToBentry is a helper function that convert a list of string
// reading from input CSV file to a Bentry.
func recordToBentry(record []any) (*Bentry, error) {
	name, ok := record[0].(string)
	if !ok {
		return nil, errors.New(`name is not a string`)
	}
	unitString, ok := record[1].(string)
	if !ok {
		return nil, errors.New(`unitString is not a string`)
	}
	scale, ok := record[2].(float64)
	if !ok {
		return nil, errors.New(`scale is not a number`)
	}
	refval, ok := record[3].(float64)
	if !ok {
		return nil, errors.New(`refval is not a number`)
	}
	nbits, ok := record[4].(float64)
	if !ok {
		return nil, errors.New(`nbits is not a number`)
	}
	crexUnitString, ok := record[5].(string)
	if !ok {
		return nil, errors.New(`crexUnitString is not a number`)
	}
	crexScale, ok := record[6].(float64)
	if !ok {
		return nil, errors.New(`crexScale is not a number`)
	}
	crexNchars, ok := record[7].(float64)
	if !ok {
		return nil, errors.New(`crexNchars is not a number`)
	}
	return &Bentry{
		name:           name,
		UnitString:     unitString,
		Unit:           unitOf(unitString),
		Scale:          int(scale),
		Refval:         int(refval),
		Nbits:          int(nbits),
		CrexUnitString: crexUnitString,
		CrexUnit:       unitOf(crexUnitString),
		CrexScale:      int(crexScale),
		CrexNchars:     int(crexNchars),
	}, nil
}

// recordToDentry is a helper function that converts a list of string
// reading form input CSV file to a Dentry.
func recordToDentry(record []any) (*Dentry, error) {
	name, ok := record[0].(string)
	if !ok {
		return nil, errors.New(`name is not a string`)
	}
	fields, ok := record[1].([]any)
	if !ok {
		return nil, errors.New(`fields is not a list`)
	}
	ids := make([]ID, len(fields))

	for i, idany := range fields {
		idstring, ok := idany.(string)
		if !ok {
			return nil, fmt.Errorf(`id at index %v is not a string`, i)
		}
		id, err := strconv.Atoi(idstring)
		if err != nil {
			return nil, fmt.Errorf(`id at index %v did not convert to an integer: %w`, i, err)
		}
		ids[i] = ID(id)
	}
	return &Dentry{name: name, Members: ids}, nil
}
