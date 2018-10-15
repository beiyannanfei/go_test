package main

import (
	"github.com/gin-gonic/gin/json"
	"fmt"
)

//生成json

type Server67 struct {
	ServerName string
	ServerIP   string
}

type Serverslice67 struct {
	Servers []Server67
}

func main() {
	var s Serverslice67
	s.Servers = append(s.Servers, Server67{ServerName: "shh_vpn", ServerIP: "127.0.0.1"})
	s.Servers = append(s.Servers, Server67{ServerName: "bj_vpn", ServerIP: "127.0.0.2"})

	b, err := json.Marshal(s)
	if err != nil {
		fmt.Println("json Marshal err:", err)
		return
	}
	fmt.Println("b:", string(b))
}
