package server

import (
	"context"
	"fmt"
	"net"

	"discache/pkg/parser"

	log "github.com/sirupsen/logrus"
)

var (
	exit = "exit"
)

type (
	// Parser receives a string and attempt to parse it to a valid Command
	// Error is returned if the given string is not a valid
	Parser interface {
		Parse(command string) (*parser.Command, error)
	}

	// CMDHandler executes the logic of a command
	// Inspired by HTTP server handler from golang official library
	CMDHandler interface {
		ServeCMD(net.Conn, parser.Command)
	}
)

// HandlerFunc is an adapter
type HandlerFunc func(net.Conn, parser.Command)

// ServeCMD calls f(conn, args...)
func (f HandlerFunc) ServeCMD(conn net.Conn, cmd parser.Command) {
	f(conn, cmd)
}

// TCP uses TCP socket for client-server connection
type TCP struct {
	Addr       string
	Parser     Parser
	CMDHandler CMDHandler
}

// Serve starts a TCP server then waits for client connections
func (t *TCP) Serve() error {
	l, err := net.Listen("tcp", t.Addr)
	if err != nil {
		return err
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			log.Errorf("Failed to accept a connection, %s", err)
			continue
		}
		log.Infof("Accepted a connection from: %s", c.RemoteAddr().String())
		conn := t.newConn(c)
		go conn.serve(context.TODO())
	}
}

func (t *TCP) newConn(c net.Conn) *conn {
	return &conn{
		server: t,
		conn:   c,
	}
}

// defaultServeCMDHandler is the default CMDHandler
var defaultServeCMDHandler = &ServeCMDHandler{
	m: make(map[parser.Op]CMDHandler),
}

// ServeCMDHandler to serve CMD
type ServeCMDHandler struct {
	m map[parser.Op]CMDHandler
}

// ServeCMD finds a CMD Handler for the given command
func (s *ServeCMDHandler) ServeCMD(conn net.Conn, cmd parser.Command) {
	h, ok := s.m[cmd.Op]
	if !ok {
		conn.Write([]byte(fmt.Sprintf("Invalid command: %s\n", cmd.Op.String())))
		return
	}
	h.ServeCMD(conn, cmd)
}

// HandleFunc registers the handler for the given Op
func (s *ServeCMDHandler) HandleFunc(op parser.Op, handler func(net.Conn, parser.Command)) {
	s.m[op] = HandlerFunc(handler)
}

// HandleFunc registers the handler function for the given Op
// in the DefaultServeCMDHandler.
func HandleFunc(op parser.Op, handler func(net.Conn, parser.Command)) {
	defaultServeCMDHandler.HandleFunc(op, handler)
}
