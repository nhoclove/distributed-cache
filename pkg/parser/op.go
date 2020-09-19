package parser

import "fmt"

// Op represents enum type of operations
type Op int

const (
	// CmdGet is get operation
	CmdGet Op = iota
	// CmdSet is set operation
	CmdSet
	// CmdDel is deleta operation
	CmdDel
)

var (
	opNames = [...]string{
		"get",
		"set",
		"del",
	}

	opValues = map[string]Op{
		"get": CmdGet,
		"set": CmdSet,
		"del": CmdDel,
	}
)

// String converts op value to its string representation
// To convert string to op use function OpValue
func (o Op) String() string {
	return opNames[o]
}

// OpValue takes a string as an argument and attempt to return the equivalent Op
// If not possible to convert an error will be returned and the value returned should be ignored
func OpValue(name string) (Op, error) {
	if op, found := opValues[name]; found {
		return op, nil
	}
	return CmdSet, fmt.Errorf("op name %s not found in op enum", name)
}
