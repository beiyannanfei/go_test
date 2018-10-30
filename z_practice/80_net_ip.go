package main

//cd /Users/wyq/workspace/go_path/src/github.com/beiyannanfei/go_test/z_practice && go run 80_net_ip.go 127.0.0.1

import (
	"os"
	"fmt"
	"net"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s ip-addr\n", os.Args[0])
		os.Exit(1)
	}

	name := os.Args[1]
	addr := net.ParseIP(name)
	if addr == nil {
		fmt.Println("Invalid address")
		return
	}

	fmt.Println("The address is ", addr.String())
	os.Exit(0)
}

