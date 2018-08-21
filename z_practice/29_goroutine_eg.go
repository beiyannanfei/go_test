package main

import "fmt"

//var readOnlyChan <-chan int            // 只读chan
//var writeOnlyChan chan<- int           // 只写chan

func sendData(ch chan<- string) { //参数为只写通道
	ch <- "go"
	ch <- "java"
	ch <- "c"
	ch <- "c++"
	ch <- "python"
	close(ch)
}

func getData(ch <-chan string, chChose chan bool) { //第一个参数为只读通道
	for {
		str, ok := <-ch
		if !ok {
			fmt.Println("chan is close.")
			break
		}
		fmt.Println(str)
	}
	chChose <- true
}

func main() {
	ch := make(chan string, 10)
	chChose := make(chan bool, 1)
	go sendData(ch)
	go getData(ch, chChose)
	<-chChose
	close(chChose)
}
