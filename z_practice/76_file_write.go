package main

import (
	"os"
	"fmt"
)

func main() {
	userFile := "/Users/wyq/workspace/go_path/src/github.com/beiyannanfei/go_test/z_practice/76_astaxie.txt"
	fout, err := os.Create(userFile)
	defer fout.Close()
	if err != nil {
		fmt.Println("create file err:", err)
		return
	}

	for i := 0; i < 10; i++ {
		fout.WriteString("fout.WriteString! \r\n")
		fout.Write([]byte("fout.Write \r\n"))
	}
}

//写文件函数:
//func (file *File) Write(b []byte) (n int, err Error) 写入byte类型的信息到文件
//func (file *File) WriteAt(b []byte, off int64) (n int, err Error) 在指定位置开始写入byte类型的信息
//func (file *File) WriteString(s string) (ret int, err Error) 写入string信息到文件