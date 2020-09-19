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
	server *TCP
	conn   net.Conn
}

func (c *conn) serve(ctx context.Context) {
	defer c.conn.Close()

	for {
		netData, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			log.Errorf("Failed to read net data from connection, %v", err)
			return
		}
		// Terminate signal from client
		if string(netData) == exit {
			log.Infof("Received terminated signal from client: %s", c.conn.RemoteAddr().String())
			return
		}

		cmd, err := c.server.Parser.Parse(strings.TrimSuffix(string(netData), "\r\n"))
		if err != nil {
			c.conn.Write([]byte(fmt.Sprintf("%s\n", err)))
			continue
		}

		c.ServeCMD(*cmd)
	}
}

func (c *conn) ServeCMD(cmd parser.Command) {
	handler := c.server.CMDHandler
	if handler == nil {
		// Use default handler
		handler = defaultServeCMDHandler
	}
	handler.ServeCMD(c.conn, cmd)
}
