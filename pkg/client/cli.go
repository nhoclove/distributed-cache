package client

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"discache/pkg/parser"
	"discache/pkg/server"
)

var (
	DefaultParser = parser.New()
)

type Cli struct {
	conns   map[string]*net.Conn
	servers []string
	parser  server.Parser
}

func New(servers []string) *Cli {
	conns := make(map[string]*net.Conn)
	for _, v := range servers {
		conns[v] = nil
	}
	return &Cli{
		conns:   conns,
		servers: servers,
		parser:  DefaultParser,
	}
}

func (cli *Cli) Connect() {
	defer cli.Disconnect()
	for {
		// reader := bufio.NewReader(os.Stdin)
		// fmt.Printf(">> ")
		// text, _ := reader.ReadString('\n')
		text := "SET NGUYEN BAC"
		if text == "exit" {
			return
		}

		command, err := cli.parser.Parse(strings.TrimSuffix(text, "\r\n"))
		if err != nil {
			fmt.Println("->: Invalid command")
		}

		conn, err := cli.getConnection(command.Args[0])
		if err != nil {
			fmt.Println("->: Failed to connect to remote servers" + err.Error())
			continue
		}

		fmt.Fprintf(*conn, "%s\n", text)
		msg, err := bufio.NewReader(*conn).ReadString('\n')
		if err != nil {
			fmt.Println("->: ", err.Error())
			continue
		}
		fmt.Print("->: " + msg)
	}

}

func (cli *Cli) Disconnect() error {
	for _, v := range cli.conns {
		if v != nil {
			(*v).Close()
		}
	}
	return nil
}

func (cli *Cli) getConnection(key string) (*net.Conn, error) {
	h := hashKeyCRC32(key)
	idx := int(h) % len(cli.servers)
	server := cli.servers[idx]
	if conn, ok := cli.conns[server]; ok && conn != nil {
		return conn, nil
	}

	conn, err := net.Dial("tcp", server)
	if err != nil {
		return nil, err
	}
	// Cache this connection for reuse
	cli.conns[server] = &conn
	return &conn, nil
}
