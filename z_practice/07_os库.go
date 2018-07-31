package main

import (
	"os"
	"fmt"
)

// cd /Users/wyq/workspace/go_path/src/github.com/beiyannanfei/go_test/z_practice && go run 07_os库.go

func main() {
	hostName, _ := os.Hostname()
	fmt.Printf("Hostname返回内核提供的主机名: %v\n", hostName)

	memoryPageSize := os.Getpagesize()
	fmt.Printf("Getpagesize返回底层的系统内存页的尺寸: %v\n", memoryPageSize)

	envList := os.Environ()
	fmt.Printf(`Environ返回表示环境变量的格式为"key=value"的字符串的切片拷贝: %v\n`, envList)

	envKey := "PATH"
	envPath := os.Getenv(envKey)
	fmt.Printf("环境变量%v的值为: %v\n", envKey, envPath)

	mapping := func(s string) string {
		fmt.Printf("mapping func s: %v\n", s)
		m := map[string]string{"widuu": "www.jb51.net", "xiaowei": "widuu"}
		return m[s]
	}
	data := "hello $xiaowei blog address $widuu"
	//Expand函数替换s中的${var}或$var为mapping(var)
	fmt.Printf("%s\n", os.Expand(data, mapping)) //hello widuu blog address www.jb51.net

	uid := os.Getuid()
	fmt.Printf("Geteuid返回调用者的有效用户ID: %v\n", uid)

	gid := os.Getgid()
	fmt.Printf("Getgid返回调用者的组ID: %v\n", gid)

	egid := os.Getegid()
	fmt.Printf("Getegid返回调用者的有效组ID: %v\n", egid)

	groupsId, _ := os.Getgroups()
	fmt.Printf("Getgroups返回调用者所属的所有用户组的组ID: %v\n", groupsId)

	pid := os.Getpid()
	fmt.Printf("Getppid返回调用者所在进程的父进程的进程ID: %v\n", pid)

	fileName := "./06_sort.go"
	fileInfo, _ := os.Stat(fileName) //Stat返回一个描述name指定的文件对象的FileInfo
	fmt.Printf("%v Name: %v, Size: %v, ModTime: %v\n", fileName, fileInfo.Name(), fileInfo.Size(), fileInfo.ModTime())

	fileMode := fileInfo.Mode()
	fmt.Printf("%v IsDir: %v, IsRegular: %v, String: %v\n", fileName, fileMode.IsDir(), fileMode.IsRegular(), fileMode.String())

	pathSeparator := os.IsPathSeparator('/')             //IsPathSeparator返回字符c是否是一个路径分隔符
	fmt.Println(`'/' is pathSeparator: `, pathSeparator) // '/' is pathSeparator:  true
	pathSeparator = os.IsPathSeparator('\\')
	fmt.Println(`'\' is pathSeparator: `, pathSeparator) // '\' is pathSeparator:  false

	fileName = "aaa"
	_, err := os.Stat(fileName)
	if err != nil { //注意：经测试IsExist、IsNotExist函数只有在err不等于nil的情况下才有效
		isExist := os.IsExist(err)       //返回一个布尔值说明该错误是否表示一个文件或目录已经存在
		isNotExist := os.IsNotExist(err) //返回一个布尔值说明该错误是否表示一个文件或目录不存在
		fmt.Printf("file: %v, err: %v, isExist: %#v, isNotExist: %#v\n", fileName, err, isExist, isNotExist)
	}

	currentDir, _ := os.Getwd() //Getwd返回一个对应当前工作目录的根路径。如果当前目录可以经过多条路径抵达（因为硬链接），Getwd会返回其中一个
	fmt.Printf("当前文件路径: %v\n", currentDir)

	fileName = "test.dat"
	file, err := os.Create(fileName) //Create采用模式0666（任何人都可读写，不可执行）创建一个名为name的文件，如果文件已存在会截断它（为空文件）。如果成功，返回的文件对象可用于I/O；对应的文件描述符具有O_RDWR模式。如果出错，错误底层类型是*PathError
	fmt.Printf("create file name: %v\n", file.Name())

	writeStrRet, err := file.WriteString("aaabbbcccddd") //WriteString类似Write，但接受一个字符串参数。
	fmt.Printf("writeStrRet: %v\n", writeStrRet)

	file, err = os.Open(fileName)
	readBuff := make([]byte, 10, 20)
	readCount, err := file.Read(readBuff)	//Read方法从f中读取最多len(b)字节数据并写入b。它返回读取的字节数和可能遇到的任何错误。文件终止标志是读取0个字节且返回值err为io.EOF
	fmt.Printf("readCount: %v, err: %v, readBuff: %v\n", readCount, err, string(readBuff))
	defer file.Close()

	err = os.Remove(fileName) //Remove删除name指定的文件或目录
	fmt.Printf("删除文件: %v 结果: %v\n", fileName, err == nil)

	dirName := "/Users/wyq/workspace/go_path/src/github.com/beiyannanfei/go_test/z_practice/"
	dirFile, _ := os.Open(dirName)
	fileInfoList, _ := dirFile.Readdir(0)
	for _, fileInfo := range fileInfoList {
		fmt.Printf("fileInfoList 文件夹中的文件有: %v\n", fileInfo.Name())
	}
	defer dirFile.Close()

	//Exit让当前程序以给出的状态码code退出。一般来说，状态码0表示成功，非0表示出错。程序会立刻终止，defer的函数不会被执行
	os.Exit(0)
}
