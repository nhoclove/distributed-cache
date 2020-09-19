package cmd

import (
	"discache/pkg/parser"
	"discache/pkg/server"
)

// CMDHandler store a Operation and its associated Handler
type CMDHandler struct {
	Op      parser.Op
	Handler server.HandlerFunc
}

// CMDHandlers returns a list of all supported commands
func (h Handler) CMDHandlers() []CMDHandler {
	return []CMDHandler{
		{
			Op:      parser.CmdGet,
			Handler: h.CMDGet,
		},
		{
			Op:      parser.CmdSet,
			Handler: h.CMDSet,
		},
		{
			Op:      parser.CmdDel,
			Handler: h.CMDDel,
		},
	}
}
