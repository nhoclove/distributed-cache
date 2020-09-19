package main

import (
	"discache/pkg/cmd"
	"discache/pkg/engine"
	"discache/pkg/parser"
	"discache/pkg/server"
	"flag"
	"strconv"
)

func main() {
	var port int
	flag.IntVar(&port, "p", 8881, "The port to client to connect to")
	flag.Parse()

	parser := parser.New()
	engine := engine.New()
	cmdHandler := cmd.New(engine)
	server := server.NewTCP(":"+strconv.Itoa(port), parser)
	// Register command handlers
	for _, h := range cmdHandler.CMDHandlers() {
		server.RegisterCMD(h.Op, h.Handler)
	}
	server.Serve()
}
