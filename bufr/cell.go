package bufr

import (
	"encoding/json"
	"fmt"
	"math"
)

// Cell is a pair of a Node and its value. It is used to build the flat
// list of decoded values for each subset from data section of a BUFR message.
type Cell struct {
	n *ValuedNode
	v interface{}
}

func NewCell(n *ValuedNode, v interface{}) *Cell {
	return &Cell{n: n, v: v}
}

func (c *Cell) Accept(visitor Visitor) error {
	return visitor.VisitCell(c)
}

func (c *Cell) String() string {
	var valueString string
	if _, ok := c.v.(string); ok {
		valueString = fmt.Sprintf("%q", c.v)
	} else {
		valueString = fmt.Sprintf("%v", c.v)
	}

	return fmt.Sprintf("%-81v %v", c.n, valueString)
}

func (c *Cell) Node() *ValuedNode {
	return c.n
}

func (c *Cell) Value() interface{} {
	return c.v
}

func (c *Cell) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.v)
}

func (c *Cell) StringValue() (string, error) {
	switch value := c.v.(type) {
	case string:
		return value, nil
	case []uint8:
		return string(value), nil
	default:
		return "", fmt.Errorf("cell value is not string, but: %T", c.v)
	}
}

func (c *Cell) UintValue() (uint, error) {
	switch value := c.v.(type) {
	case float64:
		return uint(math.Round(value)), nil
	case int:
		return uint(value), nil
	case uint:
		return value, nil
	default:
		return 0, fmt.Errorf("cell value is not compatible with uint: %T", c.v)
	}
}

func (c *Cell) IntValue() (int, error) {
	switch value := c.v.(type) {
	case float64:
		return int(math.Round(value)), nil
	case int:
		return value, nil
	case uint:
		return int(value), nil
	default:
		return 0, fmt.Errorf("cell value is not compatible with int: %T", c.v)
	}
}

func (c *Cell) FloatValue() (float64, error) {
	switch value := c.v.(type) {
	case float64:
		return value, nil
	case int:
		return float64(value), nil
	case uint:
		return float64(value), nil
	default:
		return 0, fmt.Errorf("cell value is not compatible with float: %T", c.v)
	}
}
