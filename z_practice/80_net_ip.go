package main

//cd /Users/wyq/workspace/go_path/src/github.com/beiyannanfei/go_test/z_practice && go run 80_net_ip.go 127.0.0.1

import (
	"net"
	"math/big"
	"fmt"
)

func main() {

	ret := big.NewInt(0)
	ret.SetBytes(net.ParseIP("161.117.143.191").To4())
	fmt.Println(ret.Int64())
}
