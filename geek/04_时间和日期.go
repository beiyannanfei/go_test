package main

import (
	"time"
	"fmt"
)

func main() {
	t := time.Now() //获取当前时间
	fmt.Println(t)  //2018-02-07 23:45:16.380920295 +0800 CST
	var dateTimeStr string = "%4d-%02d-%02d %02d:%02d:%02d\n"
	fmt.Printf(dateTimeStr, t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second()) //2018-02-07 23:45:16

	t = time.Now().UTC()
	fmt.Println(t)          // 2018-02-07 15:45:16.38105969 +0000 UTC
	fmt.Println(time.Now()) // 2018-02-07 23:45:16.381062806 +0800 CST
	// calculating times:
	var week time.Duration = 60 * 60 * 24 * 7 * 1e9 // must be in nanosec
	week_from_now := t.Add(week)                    //延后一周
	fmt.Println(week_from_now)                      // 2018-02-14 15:45:16.38105969 +0000 UTC
	// formatting times:
	fmt.Println(t.Format(time.RFC822))         // 07 Feb 18 15:45 UTC
	fmt.Println(t.Format(time.ANSIC))          // Wed Feb  7 15:45:16 2018
	fmt.Println(t.Format("02 Jan 2006 15:04")) // 21 Dec 2011 08:52
	s := t.Format("2006-01-02")
	fmt.Println(t, "=>", s)
	// Wed Dec 21 08:52:14 +0000 UTC 2011 => 20111221

	timeStamp := t.Unix()
	fmt.Printf("timeStamp: %v\n", timeStamp) //秒级时间戳

	timeStamp = t.UnixNano()
	fmt.Printf("timeStamp: %v\n", timeStamp) //纳秒级时间戳
}
