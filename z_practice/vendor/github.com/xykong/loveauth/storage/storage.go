package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/xykong/loveauth/settings"
	"strconv"
	"strings"
)

func Initialize() {

	dbAuth = InitMySql(settings.Get("loveauth"))
	InitAuthDatabase()
	InitRedis()

	lovepay := settings.Get("lovepay")
	if lovepay != nil && lovepay.GetBool("enable") {
		dbPay = InitMySql(lovepay)
		InitPayDatabase()
	}

	if true == settings.GetBool("tencent", "tlog.mysql.enable_ext") {
		var extMap = make(map[string]string)
		extMap["host"] = settings.GetString("tencent", "tlog.mysql.host")
		extMap["port"] = settings.GetString("tencent", "tlog.mysql.port")
		extMap["user"] = settings.GetString("tencent", "tlog.mysql.user")
		extMap["password"] = settings.GetString("tencent", "tlog.mysql.password")
		extMap["database"] = settings.GetString("tencent", "tlog.mysql.database")
		dbExt = InitMySqlWithMap(extMap)
		InitExtDatabase()
	}
}

func InitMySql(setting *viper.Viper) *gorm.DB {

	var db *gorm.DB

	if setting == nil {
		logrus.Fatalf("GORM InitMySql failed. No setting specified.")
	}

	var host = setting.GetString("mysql.host")
	var port = setting.GetInt64("mysql.port")
	var user = setting.GetString("mysql.user")
	var password = setting.GetString("mysql.password")
	var database = setting.GetString("mysql.database")

	var createDatabaseIfNotExist = setting.GetBool("mysql.create_database_if_not_exist")

	var connString = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, database)

	logrus.Infof("Connection to MySql: %s", connString)

	//open a db connection
	var err error
	db, err = gorm.Open("mysql", connString)
	if err != nil {

		Error1049 := strings.Contains(err.Error(), "1049")

		if !Error1049 || !createDatabaseIfNotExist {
			logrus.WithFields(logrus.Fields{
				"connString": connString,
				"error":      err.Error(),
			}).Fatalf("Connection mysql failed. If database is not exist, create database with the command:"+
				"'CREATE DATABASE %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;'", database)
		}

		createDatabase(setting)

		db, err = gorm.Open("mysql", connString)
	}

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"connString": connString,
			"error":      err.Error(),
		}).Fatalf("Connection mysql failed. If database is not exist, create database with the command:"+
			"'CREATE DATABASE %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;'", database)
	}

	db.DB().SetMaxIdleConns(0)  //fix bug "unexpected EOF (Invalid Connection)"

	return db
}

func createDatabase(setting *viper.Viper) {

	var host = setting.GetString("mysql.host")
	var port = setting.GetInt64("mysql.port")
	var user = setting.GetString("mysql.user")
	var password = setting.GetString("mysql.password")
	var database = setting.GetString("mysql.database")

	dataSourceName := fmt.Sprintf("%v:%v@tcp(%v:%v)/?charset=utf8", user, password, host, port)

	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"dataSourceName": dataSourceName,
			"error":          err,
		}).Error("loveauth mysql connection failed.")
	}
	defer db.Close()

	createDatabaseSql := fmt.Sprintf(
		`CREATE DATABASE IF NOT EXISTS %v DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;`,
		database)

	// create database
	rows, err := db.Query(createDatabaseSql)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"dataSourceName":    dataSourceName,
			"createDatabaseSql": createDatabaseSql,
			"rows":              rows,
			"error":             err,
		}).Error("loveauth mysql create database failed.")
	}
	rows.Close()

	logrus.WithFields(logrus.Fields{
		"dataSourceName":    dataSourceName,
		"createDatabaseSql": createDatabaseSql,
		"rows":              rows,
	}).Warn("loveauth create database success.")
}

// 强制要求必须创建
func InitMySqlWithMap(setting map[string]string) *gorm.DB {
	var db *gorm.DB

	if setting == nil {
		logrus.Fatalf("GORM InitMySql failed. No setting specified.")
	}

	var host = setting["host"]
	var port, _ = strconv.Atoi(setting["port"])
	var user = setting["user"]
	var password = setting["password"]
	var database = setting["database"]

	//var createDatabaseIfNotExist = setting["create_database_if_not_exist"]

	var connString = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, database)

	logrus.Infof("Connection to MySql: %s", connString)

	//open a db connection
	var err error
	db, err = gorm.Open("mysql", connString)
	if err != nil {

		Error1049 := strings.Contains(err.Error(), "1049")

		if !Error1049 {
			logrus.WithFields(logrus.Fields{
				"connString": connString,
				"error":      err.Error(),
			}).Fatalf("Connection mysql failed. If database is not exist, create database with the command:"+
				"'CREATE DATABASE %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;'", database)
		}

		createDatabaseWithMap(setting)

		db, err = gorm.Open("mysql", connString)
	}

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"connString": connString,
			"error":      err.Error(),
		}).Fatalf("Connection mysql failed. If database is not exist, create database with the command:"+
			"'CREATE DATABASE %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;'", database)
	}

	db.DB().SetMaxIdleConns(0)  //fix bug "unexpected EOF (Invalid Connection)"

	return db
}

func createDatabaseWithMap(setting map[string]string) {

	var host = setting["host"]
	var port, _ = strconv.Atoi(setting["port"])
	var user = setting["user"]
	var password = setting["password"]
	var database = setting["database"]

	dataSourceName := fmt.Sprintf("%v:%v@tcp(%v:%v)/?charset=utf8", user, password, host, port)

	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"dataSourceName": dataSourceName,
			"error":          err,
		}).Error("loveauth mysql connection failed.")
	}
	defer db.Close()

	createDatabaseSql := fmt.Sprintf(
		`CREATE DATABASE IF NOT EXISTS %v DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;`,
		database)

	// create database
	rows, err := db.Query(createDatabaseSql)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"dataSourceName":    dataSourceName,
			"createDatabaseSql": createDatabaseSql,
			"rows":              rows,
			"error":             err,
		}).Error("loveauth mysql create database failed.")
	}
	rows.Close()

	logrus.WithFields(logrus.Fields{
		"dataSourceName":    dataSourceName,
		"createDatabaseSql": createDatabaseSql,
		"rows":              rows,
	}).Warn("loveauth create database success.")
}

func Insert(db *gorm.DB, model interface{}) error {

	if db == nil {

		logrus.WithFields(logrus.Fields{
			"db":    db,
			"model": model,
		}).Error("database unavailable")

		return errors.New("database unavailable")
	}

	if err := db.Create(model).Error; err != nil {

		logrus.WithFields(logrus.Fields{
			"db":    db,
			"model": model,
			"error": err,
		}).Error("database insert")

		return err
	}

	logrus.WithFields(logrus.Fields{
		"db":    db,
		"model": model,
	}).Info("database insert")

	return nil
}

func Save(db *gorm.DB, model interface{}) error {

	if db == nil {

		logrus.WithFields(logrus.Fields{
			"db":    db,
			"model": model,
		}).Error("database unavailable")

		return errors.New("database unavailable")
	}

	if err := db.Save(model).Error; err != nil {

		logrus.WithFields(logrus.Fields{
			"db":    db,
			"model": model,
			"error": err,
		}).Error("database save")

		return err
	}

	logrus.WithFields(logrus.Fields{
		"db":    db,
		"model": model,
	}).Info("database save")

	return db.Save(model).Error
}
