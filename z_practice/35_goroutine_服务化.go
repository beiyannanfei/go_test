package main

import "fmt"

//比如我们加载一个网站的时候，例如我们登入新浪微博，我们的消息数据应该来自一个独立的服务，这个服务只负责 返回某个用户的新的消息提醒
//https://blog.csdn.net/skh2015java/article/details/60330975
func get_notification(name string) chan string {
	notifications := make(chan string)
	go func() {
		notifications <- fmt.Sprintf("Hi %s, welcome to weibo.com!", name)
	}()

	return notifications
}

func main() {
	jackNotifi := get_notification("jack")
	joeNotifi := get_notification("joe")

	jackV, ok := <-jackNotifi
	fmt.Printf("jackV: %v, ok: %v\n", jackV, ok)
	joeV, ok := <-joeNotifi
	fmt.Printf("joeV: %v, ok: %v\n", joeV, ok)
}
