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

	// Engine is short for StorageEngine provides functionalities to the underlying data structure
	// in this case Map data structure is used
	Engine interface {
		Set(key string, val []byte) error
		Get(key string) ([]byte, error)
		Del(key string) error
	}
)

// TCP uses TCP socket for client-server connection
type TCP struct {
	addr   string
	parser Parser
	engine Engine
}

// NewTCP initializes a new TCP server instance
func NewTCP(addr string, parser Parser, engine Engine) *TCP {
	return &TCP{
		parser: parser,
		engine: engine,
		addr:   addr,
	}
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
		conn := t.newConn(c, t.parser, t.engine)
		go conn.serve(context.TODO())
	}
}

func (t *TCP) newConn(c net.Conn, parser Parser, engine Engine) *conn {
	return &conn{
		Parser: parser,
		Engine: engine,
		Conn:   c,
	}
}
