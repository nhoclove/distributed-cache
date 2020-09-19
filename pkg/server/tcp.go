package server

import (
	"context"
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
		ServeCMD(net.Conn, ...string)
	}
)

// HandlerFunc is an adapter
type HandlerFunc func(net.Conn, ...string)

// ServeCMD calls f(conn, args...)
func (f HandlerFunc) ServeCMD(conn net.Conn, args ...string) {
	f(conn, args...)
}

// TCP uses TCP socket for client-server connection
type TCP struct {
	addr        string
	cmdHandlers map[parser.Op]CMDHandler
}

// NewTCP initializes a new TCP server instance
func NewTCP(addr string, parser Parser) *TCP {
	return &TCP{
		addr:        addr,
		cmdHandlers: make(map[parser.Op]CMDHandler, 0),
	}
}

// RegisterCMD registers a command
func (t *TCP) RegisterCMD(op parser.Op, handler CMDHandler) {
	t.cmdHandlers[op] = handler
}

// Serve starts a TCP server then waits for client connections
func (t *TCP) Serve() error {
	l, err := net.Listen("tcp", t.addr)
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
