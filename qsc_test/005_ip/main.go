package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		_, _ = fmt.Fprintf(os.Stderr, "usage: %s ip-addr\n", os.Args[0])
		os.Exit(1)
	}

	name := os.Args[1]
	addr := net.ParseIP(name)
	if addr == nil {
		fmt.Println("Invalid address")
	} else {
		fmt.Println("The address is ", addr.String())
	}

	os.Exit(0)
}

// cd /Users/wyq/workspace/go_path/src/github.com/beiyannanfei/go_test/qsc_test/005_ip && go run main.go 2002:c0e8:82e7:0:0:0:c0e8:82e7
// cd /Users/wyq/workspace/go_path/src/github.com/beiyannanfei/go_test/qsc_test/005_ip && go run main.go 127.0.0.1
