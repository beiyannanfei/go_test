package db

//数据库连接操作

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"database/sql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"reflect"
	"github.com/jinzhu/gorm"
)

const dbHost = "127.0.0.1"
const dbPort = 3306
const dbUser = "root"
const dbPwd = ""
const dbName = "gorm_test"

var dbClient *gorm.DB

func InitMySql() (*gorm.DB) { //创建mysql连接
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPwd, dbHost, dbPort, dbName)
	log.Infof("connStr: %v", connStr)
	db, err := gorm.Open("mysql", connStr)
	if err != nil {
		log.Errorf("InitMySql Open mysql err: %v", err.Error())
		return nil
	}
	return db
}

//创建数据库
func init() {
	connStr := fmt.Sprintf("%v:%v@tcp(%v:%v)/?charset=utf8", dbUser, dbPwd, dbHost, dbPort)
	log.Infof("connStr: %v", connStr)

	dbClient, err := sql.Open("mysql", connStr)
	if err != nil {
		log.WithFields(log.Fields{
			"connStr": connStr,
			"err":     err.Error(),
		}).Error("createDb client err")
		return
	}
	defer dbClient.Close() //函数退出时关闭连接

	log.Infof(`type dbClient: %v`, reflect.TypeOf(dbClient)) //type dbClient: *sql.DB

	createDbSql := fmt.Sprintf(`CREATE DATABASE IF NOT EXISTS %v DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;`, dbName)
	log.Infof(`createDbSql: %v`, createDbSql)

	rows, err := dbClient.Query(createDbSql)
	if err != nil {
		log.WithFields(log.Fields{
			"createDbSql": createDbSql,
			"err":         err,
		}).Error("create database err")
	}
	rows.Close()

	resArr, _ := rows.Columns()
	log.WithFields(log.Fields{
		"dbName": dbName,
		"rows":   resArr,
	}).Warn("create db success!")
}
