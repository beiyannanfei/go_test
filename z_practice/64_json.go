package main

import (
	"encoding/json"
	"fmt"
)

type Server64 struct {
	ServerName string
	ServerIP   string
}

type ServerSlice64 struct {
	Servers []Server64
}

func main() {
	var s ServerSlice64
	str := `{"servers":[{"serverName":"Shanghai_VPN","serverIP":"127.0.0.1"}, {"serverName":"Beijing_VPN","serverIP":"127.0.0.2"}]}`
	err := json.Unmarshal([]byte(str), &s)
	if err != nil {
		fmt.Println("json Unmarshal err:", err)
		return
	}
	fmt.Println("s:", s)
}

//通在上面的示例代码中，我们首先定义了与json数据对应的结构体，数组对应slice，字段名对应JSON里面的KEY，在 解析的时候，
// 如何将json数据与struct字段相匹配呢?例如JSON的key是Foo，那么怎么找对应的字段呢?
//首先查找tag含有Foo的可导出的struct字段(首字母大写)
// 其次查找字段名是Foo的导出字段
// 最后查找类似FOO或者FoO这样的除了首字母之外其他大小写不敏感的导出字段
//聪明的你一定注意到了这一点:能够被赋值的字段必须是可导出字段(即首字母大写)。同时JSON解析的时候只会解 析能找得到的字段，
// 如果找不到的字段会被忽略，这样的一个好处是:当你接收到一个很大的JSON数据结构而你却只 想获取其中的部分数据的时候，
// 你只需将你想要的数据对应的字段名大写，即可轻松解决这个问题。