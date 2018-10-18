package main

import (
	"os"
	"fmt"
)

func main() {
	err := os.Mkdir("astaxie", 0777) //创建名称为name的目录
	if err != nil {
		fmt.Println("mkdir astaxie err:", err)
	} else {
		fmt.Println("mkdir astaxie success")
	}

	err = os.MkdirAll("astaxie/test1/test2", 0777) //根据path创建多级子目录
	if err != nil {
		fmt.Println("mkdir astaxie/test1/test2 err:", err)
	} else {
		fmt.Println("mkdir astaxie/test1/test2 success")
	}

	err = os.Remove("astaxie") //删除名称为name的目录，当目录下有文件或者其他目录是会出错
	if err != nil {
		fmt.Println("remove astaxie err:", err)
	}

	err = os.RemoveAll("astaxie") //根据path删除多级子目录，如果path是单个名称，那么该目录不删除
	if err != nil {
		fmt.Println("removeall astaxie err:", err)
	} else {
		fmt.Println("removeall astaxie success")
	}
}
