package main

import (
	"encoding/xml"
	"os"
	"fmt"
	"io/ioutil"
)

type server struct {
	XMLName    xml.Name `xml:"server"`
	ServerName string   `xml:"serverName"`
	ServerIP   string   `xml:"serverIP"`
}

type Recurlyservers struct {
	XMLName     xml.Name `xml:"servers"`
	Version     string   `xml:"version,attr"`
	Svs         []server `xml:"server"`
	Description string   `xml:",innerxml"`
}

func main() {
	file, err := os.Open("/Users/wyq/workspace/go_path/src/github.com/beiyannanfei/go_test/z_practice/62_xml.xml")
	if err != nil {
		fmt.Println("open file err:", err)
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("read all err:", err)
		return
	}

	v := Recurlyservers{}
	err = xml.Unmarshal(data, &v)
	if err != nil {
		fmt.Println("xml Unmarshal err:", err)
		return
	}

	fmt.Println(v)
}
