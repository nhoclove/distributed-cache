package parser

import (
	"fmt"
	"strings"
)

// String tries to parse a string to a Command
type String struct{}

// New initializes a new instance
func New() String {
	return String{}
}

// Parse parses a string to a Command
// Invalid command error will be returned if the given string is not parsable
// Command structure: <Op> <Arg> [<Arg>...]
func (p String) Parse(command string) (*Command, error) {
	strs := strings.Split(command, " ")
	if len(strs) < 2 {
		return nil, fmt.Errorf("invalid command: %s", command)
	}
	op, err := OpValue(strs[0])
	if err != nil {
		return nil, err
	}
	return &Command{
		Op:   op,
		Args: strs[1:],
	}, nil
}
