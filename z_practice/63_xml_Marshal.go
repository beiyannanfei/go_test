package main

import (
	"encoding/xml"
	"fmt"
	"os"
)

//生成xml

type Servers struct {
	XMLName xml.Name  `xml:"servers"`
	Version string    `xml:"version,attr"`
	Svs     []server1 `xml:"server"`
}

type server1 struct {
	ServerName string `xml:"serverName"`
	ServerIP   string `xml:"serverIP"`
}

func main() {
	v := &Servers{Version: "1"}
	v.Svs = append(v.Svs, server1{"shh_VPN", "127.0.0.1"})
	v.Svs = append(v.Svs, server1{"bj_VPN", "127.0.0.2"})
	output, err := xml.MarshalIndent(v, " ", "    ")
	if err != nil {
		fmt.Println("xml.MarshalIndent err:", err)
		return
	}
	os.Stdout.Write([]byte(xml.Header))
	os.Stdout.Write(output)
}
