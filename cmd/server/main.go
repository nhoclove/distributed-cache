package main

import (
	"flag"
	"strconv"

	"discache/pkg/cmd"
	"discache/pkg/engine"
	"discache/pkg/parser"
	"discache/pkg/server"
)

func main() {
	var port int
	flag.IntVar(&port, "p", 8881, "The port to client to connect to")
	flag.Parse()

	parser := parser.New()
	engine := engine.New()
	cmdHandler := cmd.New(engine)

	// Register command handler
	for _, h := range cmdHandler.CMDHandlers() {
		server.HandleFunc(h.Op, h.Handler)
	}

	// Start TCP server
	server := server.TCP{
		Addr:       ":" + strconv.Itoa(port),
		Parser:     parser,
		CMDHandler: nil,
	}
	server.Serve()
}
