package main

import (
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
	server := server.NewTCP(":"+strconv.Itoa(port), parser, engine)
	server.Serve()
}
