package tlog

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/settings"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"github.com/xykong/tlog"
	"net"
	"time"
)

func Start() {

	go updateGameSvrState()
	var isRunLogout = settings.GetBool("loveauth", "run_logout")
	if isRunLogout {
		go updateOnlineTable()
	}

}

func updateOnlineTable() {

	for {

		doUpdateOnlineTable()

		time.Sleep(time.Second * time.Duration(settings.GetInt("tencent", "tlog.updateOnlineTableSeconds")))
	}
}

//noinspection ALL
func doUpdateOnlineTable() {

	if storage.HasOnlineLog() {

		return
	}

	user := settings.GetString("tencent", "tlog.mysql.user")
	password := settings.GetString("tencent", "tlog.mysql.password")
	host := settings.GetString("tencent", "tlog.mysql.host")
	port := settings.GetInt("tencent", "tlog.mysql.port")
	database := settings.GetString("tencent", "tlog.mysql.database")
	table := settings.GetString("tencent", "tlog.mysql.online_table")

	dataSourceName := fmt.Sprintf("%v:%v@tcp(%v:%v)/?charset=utf8", user, password, host, port)

	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"dataSourceName": dataSourceName,
			"error":          err,
		}).Error("Tlog mysql connection failed.")
	}
	defer db.Close()

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"dataSourceName": dataSourceName,
			"error":          err,
		}).Error("Tlog mysql ping failed.")
	}

	createDatabaseSql := fmt.Sprintf(`CREATE DATABASE IF NOT EXISTS %v DEFAULT CHARACTER SET utf8;`, database)

	// create database
	rows, err := db.Query(createDatabaseSql)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"dataSourceName":    dataSourceName,
			"createDatabaseSql": createDatabaseSql,
			"rows":              rows,
			"error":             err,
		}).Error("Tlog mysql create table failed.")
	}
	rows.Close()

	createTableSql := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %v.%v (
       channel varchar(100) not null DEFAULT '',
	   timekey int(11) NOT NULL DEFAULT 0,
	   gsid varchar(32) NOT NULL DEFAULT '',
	   onlinecntios int(11) NOT NULL DEFAULT '0',
	   onlinecntandroid int(11) NOT NULL DEFAULT '0',
	  KEY (timekey,gsid)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`, database, table)

	// create table
	rows, err = db.Query(createTableSql)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"dataSourceName": dataSourceName,
			"createTableSql": createTableSql,
			"rows":           rows,
			"error":          err,
		}).Error("Tlog mysql create table failed.")
	}
	rows.Close()

	insertSql := fmt.Sprintf("INSERT %v.%v SET channel=?, timekey=?, gsid=?, onlinecntios=?, onlinecntandroid=?",
		database, table)

	stmt, err := db.Prepare(insertSql)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"dataSourceName": dataSourceName,
			"insertSql":      insertSql,
			"rows":           rows,
			"error":          err,
		}).Error("Tlog mysql prepare failed.")
	}
	defer stmt.Close()

	timekey := time.Now().Unix()
	gsid := settings.GetString("tencent", "GameSvrId")

	channelList := storage.GetChannelList()
	for _, channel := range channelList {
	//for _, platform := range []string{"QQ", "Wechat", "Guest"} {

		//gameappid := settings.GetString("tencent", fmt.Sprintf("msdk.%v.AppId", platform))
		//zoneareaid := settings.GetInt("tencent", fmt.Sprintf("msdk.%v.iZoneAreaID", platform))

		//var vendor model.Vendor
		//switch platform {
		//case "QQ":
		//
		//	vendor = model.VendorTencentQQ
		//case "Wechat":
		//
		//	vendor = model.VendorTencentWechat
		//case "Guest":
		//
		//	vendor = model.VendorTencentGuest
		//}

		onlinecntios := storage.CountOnlineToken(channel, model.PlatformIOS)
		onlinecntandroid := storage.CountOnlineToken(channel, model.PlatformAndroid)

		//noinspection ALL
		logrus.WithFields(logrus.Fields{
			"dataSourceName":   dataSourceName,
			//"gameappid":        gameappid,
			"channel":          channel,
			"timekey":          timekey,
			"gsid":             gsid,
			//"zoneareaid":       zoneareaid,
			"onlinecntios":     onlinecntios,
			"onlinecntandroid": onlinecntandroid,
		}).Info("Tlog mysql exec.")

		res, err := stmt.Exec(channel, timekey, gsid, onlinecntios, onlinecntandroid)
		//res, err := stmt.Exec(gameappid, platform, timekey, gsid, zoneareaid, onlinecntios, onlinecntandroid)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"dataSourceName": dataSourceName,
				"insertSql":      insertSql,
				"rows":           rows,
				"error":          err,
			}).Error("Tlog mysql exec failed.")
		}

		affect, err := res.RowsAffected()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"dataSourceName": dataSourceName,
				"affect":         affect,
				"rows":           rows,
				"error":          err,
			}).Error("Tlog mysql Exec RowsAffected failed.")
		}
	}
}

//<!--//////////////////////////////////////////////
/////////服务器状态日志///////////////////////////////
///////////////////////////////////////////////////-->
//<struct  name="GameSvrState" filter="1"  version="1" desc="(必填)服务器状态流水，每5分钟一条日志">
//	<entry  name="dtEventTime"        type="datetime"                    desc="(必填) 格式 YYYY-MM-DD HH:MM:SS" />
//	<entry  name="vGameIP"            type="string"        size="32"                        desc="(必填)服务器IP" />
//	<entry  name="iZoneAreaID"        type="int"            index="1"            defaultvalue="0"    desc="(必填)针对分区分服的游戏填写分区id，用来唯一标示一个区；非分区分服游戏请填写0"/>
//  <entry name="Format" type="string" size="64" defaultvalue="1.0.0" desc="Format version of this log"/>
//</struct>
func updateGameSvrState() {

	vGameIP := getAddress()
	address := settings.GetString("tencent", "tlog.address")
	format := settings.GetString("tencent", "tlog.format")
	gameSvrId := settings.GetString("tencent", "GameSvrId")

	for {

		for _, platform := range []string{"QQ", "Wechat", "Guest"} {

			iZoneAreaID := settings.GetInt("tencent", fmt.Sprintf("msdk.%v.iZoneAreaID", platform))
			dtEventTime := time.Now().Format("2006-01-02 15:04:05")

			if address != "" {

				tlog.LogRaw(fmt.Sprintf("GameSvrState|%s|%s|%d|%s\n",
					dtEventTime, vGameIP, iZoneAreaID, format))
			}

			logrus.WithFields(logrus.Fields{
				"EventTime": dtEventTime,
				"GameIP":     vGameIP,
				"ZoneAreaId": iZoneAreaID,
				"FlowName":    "GameSvrState",
				"Format":      format,
				//为graylog过滤添加gamesvrid字段
				"GameSvrId": gameSvrId,
			}).Info("Tlog GameSvrState")
		}

		time.Sleep(time.Second * time.Duration(settings.GetInt("tencent", "tlog.updateGameSvrStateSeconds")))
	}
}

func getAddress() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("get ip address failed.")

		return "localhost"
	}

	for _, a := range addrs {
		if ipNet, ok := a.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}

	return "localhost"
}
