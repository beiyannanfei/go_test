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
	fmt.Printf("%s\n", os.Expand(data, mapping))	//hello widuu blog address www.jb51.net

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

	currentDir, _ := os.Getwd()
	fmt.Printf("当前文件路径: %v\n", currentDir)

	fileName := "./06_sort.go"
	fileInfo, _ := os.Stat(fileName)


	fileMode := fileInfo.Mode()
	fmt.Printf("%v IsDir: %v, IsRegular: %v, String: %v\n", fileName, fileMode.IsDir(), fileMode.IsRegular(), fileMode.String())


	//Exit让当前程序以给出的状态码code退出。一般来说，状态码0表示成功，非0表示出错。程序会立刻终止，defer的函数不会被执行
	os.Exit(0)
}
