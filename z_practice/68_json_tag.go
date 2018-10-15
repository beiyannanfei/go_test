package main

import (
	"encoding/json"
	"os"
)

//字段的tag是"-"，那么这个字段不会输出到JSON
// tag中带有自定义名称，那么这个自定义名称会出现在JSON的字段名中，例如上面例子中serverName
// tag中如果带有"omitempty"选项，那么如果该字段值为空，就不会输出到JSON串中
//如果字段类型是bool, string, int, int64等，而tag中带有",string"选项，那么这个字段在输出到JSON的 时候会把该字段对应的值转换成JSON字符串

type Server68 struct {
	ID          int    `json:"-"`           //ID不会导出到json中
	ServerName  string `json:"server_name"` //ServerName的值会进行二次json编码
	ServerName2 string `json:"server_name_2,string"`
	ServerIP    string `json:"server_ip,omitempty"` //如果ServerIP为空，则不输出到JSON串中
}

func main() {
	s := Server68{
		ID:          1,
		ServerName:  `Go "1.0" `,
		ServerName2: `Go "1.2" `,
		ServerIP:    ``,
	}
	b, _ := json.Marshal(s)
	os.Stdout.Write(b)
}
