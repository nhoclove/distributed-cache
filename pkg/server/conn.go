package server

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strings"

	"discache/pkg/parser"

	log "github.com/sirupsen/logrus"
)

type conn struct {
	Parser
	Engine
	net.Conn
}

func (c *conn) serve(ctx context.Context) {
	defer c.Conn.Close()

	for {
		netData, err := bufio.NewReader(c.Conn).ReadString('\n')
		if err != nil {
			log.Errorf("Failed to read net data from connection, %v", err)
			return
		}
		// Terminate signal from client
		if string(netData) == exit {
			log.Infof("Received terminated signal from client: %s", c.Conn.RemoteAddr().String())
			return
		}

		cmd, err := c.Parser.Parse(strings.TrimSuffix(string(netData), "\r\n"))
		if err != nil {
			c.Conn.Write([]byte(fmt.Sprintf("%s\n", err)))
			continue
		}

		c.serveCommand(*cmd)
	}
}

func (c *conn) serveCommand(cmd parser.Command) {
	switch cmd.Op {
	case parser.CmdGet:
		c.serveCommandGet(cmd)
	case parser.CmdSet:
		c.serveCommandSet(cmd)
	case parser.CmdDel:
		c.serveCommandDel(cmd)
	default:
		c.Conn.Write([]byte(fmt.Sprintf("Invalid operation: %s\n", cmd.Op.String())))
	}
}

func (c *conn) serveCommandGet(cmd parser.Command) {
	if len(cmd.Args) < 1 {
		c.Conn.Write([]byte("GET command requires a key argument\n"))
		return
	}
	if len(cmd.Args) > 1 {
		c.Conn.Write([]byte("Too many arguments\n"))
		return
	}
	val, err := c.Engine.Get(cmd.Args[0])
	if err != nil {
		c.Conn.Write([]byte(fmt.Sprintf("%s\n", err.Error())))
		return
	}
	_, err = c.Conn.Write(append(val, byte('\n')))
	if err != nil {
		log.Errorf("Failed to write in GET, %v", err)
	}
}

func (c *conn) serveCommandSet(cmd parser.Command) {
	if len(cmd.Args) < 2 {
		c.Conn.Write([]byte("SET command requires a key and value argument\n"))
		return
	}
	if len(cmd.Args) > 2 {
		c.Conn.Write([]byte("Too many arguments\n"))
		return
	}
	err := c.Engine.Set(cmd.Args[0], []byte(cmd.Args[1]))
	if err != nil {
		c.Conn.Write([]byte(err.Error()))
		return
	}
	_, err = c.Conn.Write([]byte("ok\n"))
	if err != nil {
		log.Errorf("Failed to write in SET, %v", err)
	}
}

func (c *conn) serveCommandDel(cmd parser.Command) {
	if len(cmd.Args) < 1 {
		c.Conn.Write([]byte("Del command requires a key argument\n"))
		return
	}
	if len(cmd.Args) > 1 {
		c.Conn.Write([]byte("Too many arguments\n"))
		return
	}
	err := c.Engine.Del(cmd.Args[0])
	if err != nil {
		c.Conn.Write([]byte(err.Error()))
	}
	_, err = c.Conn.Write([]byte("ok\n"))
	if err != nil {
		log.Errorf("Failed to write in DEL, %v", err)
	}
}
