package main

import (
	"fmt"
	"time"
)

func main() {
	// 获取本月开始时间
	year, month, _ := time.Now().Date()
	thisMonthStr := time.Date(year, month, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02")
	fmt.Printf("thisMonthStr: %+v\n", thisMonthStr)

	fmt.Println("=====================================")
	// 获取本周一时间
	now := time.Now()
	weekDay := now.Weekday() //默认是 Sunday 开始到 Saturday 算 0,1,2,3,4,5,6
	offset := int(time.Monday - weekDay)
	if offset > 0 { //如果今天是周日则去上周一，即减六天
		offset = -6
	}

	weekStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	weekStartStr := weekStart.Format("2006-01-02")
	fmt.Printf("weekStartStr: %+v\n", weekStartStr)
}
