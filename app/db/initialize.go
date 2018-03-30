package db

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

const dbHost = "127.0.0.1"
const dbPort = 3306
const dbUser = "root"
const dbPwd = ""
const dbName = "gorm_test"

//创建数据库
func createDb() {
	connStr := fmt.Sprintf("%v：%v@tcp(%v:%v)/?charset=utf8", dbUser, dbPwd, dbHost, dbPort)
	log.Infof("connStr: %v", connStr)
}


