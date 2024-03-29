package main

import (
	"time"
	"fmt"
)

// 13~23 例子来源 https://juejin.im/post/5b5bd2ddf265da0f716c2fea?utm_source=gold_browser_extension

func main() {
	time := time.Now()
	fmt.Println(time)

	time1 := time.Format("20060102") //相当于Ymd
	fmt.Println(time1)

	time1 = time.Format("2006-01-02") //相当于Y-m-d
	fmt.Println(time1)

	time1 = time.Format("2006-01-02 15:04:05") //相当于Y-m-d H:i:s
	fmt.Println(time1)

	time1 = time.Format("2006-01-02 00:00:00") //相当于Y-m-d 00:00:00
	fmt.Println(time1)
}
