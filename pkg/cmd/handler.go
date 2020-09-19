package cmd

import (
	"fmt"
	"net"

	"discache/pkg/parser"

	log "github.com/sirupsen/logrus"
)

// Engine is short for StorageEngine provides functionalities to the underlying data structure
// in this case Map data structure is used
type Engine interface {
	Set(key string, val []byte) error
	Get(key string) ([]byte, error)
	Del(key string) error
}

// Handler responsible for handling supported commands
// Currently we only support commands:
// GET: Retrieves a value associating with a key
// SET: Sets a key-value pair in cache
// Del: Deletes a key-value pair in cache
type Handler struct {
	engine Engine
}

// New returns a new instance of CMD Handler
func New(engine Engine) *Handler {
	return &Handler{
		engine: engine,
	}
}

// CMDGet retrieves a value with the provided key
func (h Handler) CMDGet(conn net.Conn, cmd parser.Command) {
	val, err := h.engine.Get(cmd.Args[0])
	if err != nil {
		conn.Write([]byte(fmt.Sprintf("%s\n", err.Error())))
	}
	_, err = conn.Write(append(val, byte('\n')))
	if err != nil {
		log.Errorf("Failed to write in GET, %v", err)
	}
}

// CMDSet sets a key-value pair in cache
func (h Handler) CMDSet(conn net.Conn, cmd parser.Command) {
	err := h.engine.Set(cmd.Args[0], []byte(cmd.Args[1]))
	if err != nil {
		conn.Write([]byte(err.Error()))
	}
	_, err = conn.Write([]byte("ok\n"))
	if err != nil {
		log.Errorf("Failed to write in SET, %v", err)
	}
}

// CMDDel deletes a key-value pair in cache
func (h Handler) CMDDel(conn net.Conn, cmd parser.Command) {
	err := h.engine.Del(cmd.Args[0])
	if err != nil {
		conn.Write([]byte(err.Error()))
	}
	_, err = conn.Write([]byte("ok\n"))
	if err != nil {
		log.Errorf("Failed to write in DEL, %v", err)
	}
}
