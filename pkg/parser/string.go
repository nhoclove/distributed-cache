package parser

import (
	"fmt"
	"strings"

	"discache/pkg/server"
)

type String struct {
}

func New() String {
	return String{}
}

func (p String) Parse(command string) (*server.Command, error) {
	strs := strings.Split(command, " ")
	if len(strs) < 2 {
		return nil, fmt.Errorf("invalid command: %s", command)
	}
	return &server.Command{
		Op:   strs[0],
		Args: strs[1:],
	}, nil
}
