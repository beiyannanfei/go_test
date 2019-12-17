package main

import (
	"encoding/xml"
	"fmt"
)

type Servers struct {
	XMLName xml.Name `xml:"servers"`
	Version string   `xml:"version,attr"`
	Svs     []server `xml:"server"`
}

type server struct {
	ServerName string `xml:"serverName"`
	ServerIP   string `xml:"serverIP"`
}

func main() {
	v := &Servers{Version: "1.0"}
	v.Svs = append(v.Svs, server{"shanghai_VPN", "127.0.0.1"})
	v.Svs = append(v.Svs, server{"beijing_VPN", "127.0.0.2"})

	output, err := xml.MarshalIndent(v, "  ", "  ")
	if err != nil {
		fmt.Println("err: ", err)
		return
	}

	fmt.Printf(xml.Header) //=>	os.Stdout.Write([]byte(xml.Header))

	fmt.Println(string(output)) //=> os.Stdout.Write(output)
}
