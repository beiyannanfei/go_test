package db

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

//对数据库的增删改查操作

func do() {
	execSql()
}

func execSql() {
	log.Info("=================================")
	tbName := "myUserInfo2"
	field1 := "user_name"
	field2 := "age"
	field3 := "addr"
	sqlStr := fmt.Sprintf(`insert into %v (%v, %v, %v) values (?,?,?), (?,?,?)`, tbName, field1, field2, field3)
	log.Infof("execSql sqlStr: %v", sqlStr)

	args := make([]interface{}, 0, 6)
	args = append(args, "aaa", "25", "beijing", "bbb", "30", "shanghai")
	sqlResponse, err := dbClient.DB().Exec(sqlStr, args...)
	if err != nil {
		log.Errorf("sql exec err: %v", err.Error())
		return
	}
	lastInsertId, err := sqlResponse.LastInsertId()
	affectCount, err := sqlResponse.RowsAffected()
	if err != nil {
		log.Errorf("parse sqlResponse err: %v", err.Error())
		return
	}
	log.WithFields(log.Fields{
		"lastInsertId": lastInsertId,
		"affectCount":  affectCount,
	}).Info("sql cmd exec result")
}
