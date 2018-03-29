package base

import (
	"time"
	log "github.com/sirupsen/logrus"
	"reflect"
	"fmt"
)

func Time_test() {
	//解析时间
	timeStr := time.Now().Format("2006-01-02 15:04:05") //注意：记忆时间格式的方法，从月开始数字1-5(必须是这个字符串，否则时间不准确)
	log.Infof("now timeStr: %v, type: %v", timeStr, reflect.TypeOf(timeStr))

	//格式化时间
	timeParse, _ := time.Parse("2006-01-02 15:04:05", "2008-12-01 10:05:49")
	log.Infof("timeParse: %v", timeParse)

	fmt.Println(time.Now().AddDate(0, 1, 0))  //下个月的今天
	fmt.Println(time.Now().AddDate(0, 0, -1)) //昨天

	durationH, _ := time.ParseDuration("-1h") //一小时前
	fmt.Println(time.Now().Add(durationH))

	durationM, _ := time.ParseDuration("-1m") //一分钟前
	fmt.Println(time.Now().Add(durationM))

	durationS, _ := time.ParseDuration("-1s") //一秒钟前
	fmt.Println(time.Now().Add(durationS))

	timeStamp := time.Now().UnixNano() / 1e6 //毫秒级时间戳
	fmt.Println(timeStamp)

	timeStamp = time.Now().Unix() //秒级时间戳
	fmt.Println(timeStamp)
}
