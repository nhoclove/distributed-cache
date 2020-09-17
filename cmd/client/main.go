package main

import (
	"flag"
	"fmt"
	"strings"

	"discache/pkg/client"
)

func main() {
	var (
		serverList string
	)
	flag.StringVar(&serverList, "servers", "localhost:8801,localhost:8802,localhost:8803", "The list of cache servers")
	flag.Parse()

	if len(strings.TrimSpace(serverList)) == 0 {
		fmt.Printf("Invalid server list")
		return
	}

	servers := strings.Split(serverList, ",")
	// For now only connect to the first server in the provided list
	cli := client.New(servers)
	cli.Connect()
}
