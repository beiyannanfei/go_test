package log

import (
	"os"
	log "github.com/sirupsen/logrus"
)

func init() {
	//log.SetFormatter(&log.JSONFormatter{})		//输出为json结构
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	log.SetOutput(os.Stdout)

	log.SetLevel(log.DebugLevel) //日志输出基本
}

func LogTest() {
	argNum := 123
	argStr := "abc"
	log.Debugf("debug level argNum: %v, argStr: %v", argNum, argStr)

	log.Infof("debug level argNum: %v, argStr: %v", argNum, argStr)
	log.Warnf("debug level argNum: %v, argStr: %v", argNum, argStr)
	log.Errorf("debug level argNum: %v, argStr: %v", argNum, argStr)
	log.Fatalf("debug level argNum: %v, argStr: %v", argNum, argStr)	//exit status 1
	//log.Panicf("debug level argNum: %v, argStr: %v", argNum, argStr)


	log.WithFields(log.Fields{
		"animal": "walrus",
		"size":   10,
	}).Debugf("debug level argNum: %d, argStr: %s", argNum, argStr)
}
