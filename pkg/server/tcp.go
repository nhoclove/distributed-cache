package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	log "github.com/sirupsen/logrus"
)

var (
	exit = "exit"
)

type (
	Parser interface {
		Parse(command string) (*Command, error)
	}

	Engine interface {
		Set(key string, val []byte) error
		Get(key string) ([]byte, error)
		Del(key string) error
	}
)

type TCP struct {
	addr   string
	parser Parser
	engine Engine
}

func NewTCP(addr string, parser Parser, engine Engine) *TCP {
	return &TCP{
		addr:   addr,
		parser: parser,
		engine: engine,
	}
}

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
		go t.handle(c)
	}
}

func (t *TCP) handle(conn net.Conn) {
	defer conn.Close()

	for {
		netData, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Errorf("Failed to read net data from connection, %v", err)
			return
		}
		// Terminate signal from client
		if string(netData) == exit {
			log.Infof("Received terminated signal from client: %s", conn.RemoteAddr().String())
			return
		}

		command, err := t.parser.Parse(strings.TrimSuffix(string(netData), "\r\n"))
		if err != nil {
			conn.Write([]byte(fmt.Sprintf("%s\n", err)))
			continue
		}

		switch command.Op {
		case "GET":
			val, err := t.engine.Get(command.Args[0])
			if err != nil {
				conn.Write([]byte(fmt.Sprintf("%s\n", err.Error())))
				break
			}
			_, err = conn.Write(append(val, byte('\n')))
			if err != nil {
				log.Errorf("Failed to write in GET, %v", err)
			}
		case "SET":
			err := t.engine.Set(command.Args[0], []byte(command.Args[1]))
			if err != nil {
				conn.Write([]byte(err.Error()))
				break
			}
			_, err = conn.Write([]byte("ok\n"))
			if err != nil {
				log.Errorf("Failed to write in SET, %v", err)
			}
		case "DEL":
			err := t.engine.Del(command.Args[0])
			if err != nil {
				conn.Write([]byte(err.Error()))
				break
			}
			_, err = conn.Write([]byte("ok\n"))
			if err != nil {
				log.Errorf("Failed to write in DEL, %v", err)
			}
		default:
			conn.Write([]byte(fmt.Sprintf("Invalid operation: %s\n", command.Op)))
		}
	}
}

func (t *TCP) get(args []string) (interface{}, error) {
	return t.engine.Get(args[0])
}
