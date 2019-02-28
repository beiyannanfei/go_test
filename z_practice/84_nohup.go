package main

import (
	"os"
	"fmt"
	"time"
)

func main() {
	userFile := "/Users/wyq/workspace/go_path/src/github.com/beiyannanfei/go_test/z_practice/84_log"
	fout, err := os.Create(userFile)
	defer fout.Close()
	if err != nil {
		fmt.Println("create file err:", err)
		return
	}

	for {
		fout.WriteString("fout.WriteString! \r\n")
		fout.Write([]byte("fout.Write \r\n"))
		time.Sleep(time.Second)
	}
}

/*
1. go build 84_nohup.go    -- 生成可执行程序
2. nohup ./84_nohup &      --让程序在后台执行
3. ps aux | grep 84_nohup  --查看后台程序的pid
4. kill 11273			   --杀掉后台进程
*/

