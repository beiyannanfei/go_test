package main

import (
	"time"
	"fmt"
	"reflect"
)

func main() {
	var t time.Time
	t = time.Now()
	fmt.Printf("时间: %v，时区: %v，时间类型: %v\n", t, t.Location(), reflect.TypeOf(t))
	//时间: 2018-08-21 15:47:28.636331313 +0800 CST m=+0.000455518，时区: Local，时间类型: time.Time

	fmt.Printf("时间: %v，时区: %v，时间类型: %v\n", t.UTC(), t.UTC().Location(), reflect.TypeOf(t))
	//时间: 2018-08-21 07:48:20.114690994 +0000 UTC，时区: UTC，时间类型: time.Time

	fmt.Println("---------------- 根据时间戳返回本地时间 ----------------")
	t_by_unix := time.Unix(1534839192, 676*1000*1000)
	fmt.Printf("t_by_unix: %v\n", t_by_unix) //t_by_unix: 2018-08-21 16:13:12.676 +0800 CST

	fmt.Println("----------------------- 返回指定时间 -----------------------")
	t_by_data := time.Date(2018, time.Month(8), 21, 16, 24, 12, 0, time.Local)
	fmt.Printf("t_by_data: %v\n", t_by_data) //t_by_data: 2018-08-21 16:24:12 +0800 CST

	fmt.Println("--------------------- 时间显示 ---------------------")
	t_by_utc := t.UTC()             //获取指定时间在UTC 时区的时间表示
	fmt.Println("t.UTC:", t_by_utc) //t.UTC: 2018-08-21 08:54:42.868850966 +0000 UTC

	t_by_local := t.Local()             //获取本地时间表示
	fmt.Println("t.Local:", t_by_local) //t.Local: 2018-08-21 16:55:30.953676202 +0800 CST

	t_in := t.In(time.UTC)     //时间在指定时区的表示
	fmt.Println("t.In:", t_in) //t.In: 2018-08-21 08:56:50.74447049 +0000 UTC

	t_format_rfc := t.Format(time.RFC3339)     //format
	fmt.Println("t_format_rfc:", t_format_rfc) //t_format_rfc: 2018-08-21T16:59:30+08:00

	t_format_ymd := t.Format("2006-01-02")
	fmt.Println("t_format_ymd:", t_format_ymd) //t_format_ymd: 2018-08-21

	t_format_ymdhms := t.Format("2006-01-02 15:04:05")
	fmt.Println("t_format_ymdhms:", t_format_ymdhms) //t_format_ymdhms: 2018-08-21 18:08:53

	fmt.Println("---------------------- 获取日期信息 ----------------------")
	year, month, day := t.Date()             // 返回时间的日期信息
	fmt.Println("t.Data:", year, month, day) //t.Data: 2018 August 21

	week := t.Weekday()             // 星期
	fmt.Println("t.Weekday:", week) //t.Weekday: Tuesday

	year, week_int := t.ISOWeek()             // 返回年，星期范围编号
	fmt.Println("w.ISOWeek:", year, week_int) //w.ISOWeek: 2018 34

	hour, min, sec := t.Clock()             //返回时间的时分秒
	fmt.Println("t.Clock:", hour, min, sec) //t.Clock: 18 26 4

	fmt.Println("--------------------- 时间比较与计算 ---------------------")
	t_add_1second := t.Add(time.Second)                                                                                    //加1秒
	fmt.Printf("t: %v, t_add_1second: %v\n", t.Format("2006-01-02 15:04:05"), t_add_1second.Format("2006-01-02 15:04:05")) //t: 2018-08-21 18:39:11, t_add_1second: 2018-08-21 18:39:12

	t_add_date := t.AddDate(1, 1, 1)                                                                                 //加一年一月一日
	fmt.Printf("t: %v, t_add_date: %v\n", t.Format("2006-01-02 15:04:05"), t_add_date.Format("2006-01-02 15:04:05")) //t: 2018-08-21 18:41:48, t_add_date: 2019-09-22 18:41:48

	fmt.Println("------------------ 时间序列化 ------------------")
	t_byte, err := t.MarshalJSON()                     // 时间序列化
	fmt.Println("t.MarshalJSON:", string(t_byte), err) //t.MarshalJSON: "2018-08-21T19:06:38.764802122+08:00" <nil>

	var t_un time.Time
	err = t_un.UnmarshalJSON(t_byte)              //时间数据反序列化
	fmt.Println("t_un.UnmarshalJSON:", t_un, err) //t_un.UnmarshalJSON: 2018-08-21 19:08:05.139775192 +0800 CST <nil>

	fmt.Println("-------------------- time.Duration 方法 --------------------")
	d := time.Duration(10000000000000)
	fmt.Printf("string: %v, seconds: %v, minutes: %v, hours: %v\n", d.String(), d.Seconds(), d.Minutes(), d.Hours()) //string: 2h46m40s, seconds: 10000, minutes: 166.66666666666666, hours: 2.7777777777777777

	fmt.Println("----------------------- 其他方法 -----------------------")
	d_second := time.Second
	time.Sleep(d_second)
}
