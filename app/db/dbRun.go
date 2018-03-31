package db

import "github.com/sirupsen/logrus"

//用来执行db的各种操作

func Run() {
	dbClient= InitMySql()
	if dbClient == nil {
		logrus.Warn("InitMySql error db client nil")
		return
	}

	//createTb()

	//checkTableExists()

	do()
}
