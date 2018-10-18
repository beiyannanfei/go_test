package main

import (
	"os"
	"fmt"
)

//读文件函数:
//func (file *File) Read(b []byte) (n int, err Error) 读取数据到b中
//func (file *File) ReadAt(b []byte, off int64) (n int, err Error) 从off开始读取数据到b中

func main() {
	userFile := "/Users/wyq/workspace/go_path/src/github.com/beiyannanfei/go_test/z_practice/77_file_read.go"
	fl, err := os.Open(userFile)
	defer fl.Close()
	if err != nil {
		fmt.Println("open file err:", err)
		return
	}

	buf := make([]byte, 1024)
	for {
		n, _ := fl.Read(buf)
		if n == 0 {
			break
		}

		os.Stdout.Write(buf[:n])
	}
}
