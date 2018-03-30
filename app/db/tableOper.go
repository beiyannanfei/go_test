package db

import (
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

//表操作

//创建表
func createTb() {
	dbClient.AutoMigrate(&UserInfo{}, &UserInfo1{})            //根据表结构创建数据表
	dbClient.AutoMigrate(&UserInfo2{tableName: "myUserInfo2"}) //自定义表名，类似 myUserInfo2
}

func checkTableExists() { //检测表是否存在
	//通过调用db.HasTable来判断是否存在表，但是参数可以使用两种形式，一种是表名的字符串，一种是模型的地址类型。
	isHaveT1 := dbClient.HasTable(&UserInfo{})
	isHaveT2 := dbClient.HasTable("abcd")
	isHaveT3 := dbClient.HasTable(`user_info1`) //todo 如果是字符串类型必须是和数据库中的表名完全相同(不区分大小写)
	log.WithFields(log.Fields{
		`Table UserInfo`:   isHaveT1,
		`Table abcd`:       isHaveT2,
		`Table user_info1`: isHaveT3,
	}).Info(`checkTableExists`)
}

type (
	UserInfo struct {
		//数据库中数据表的名称为 user_infos
		gorm.Model                        //加上这个字段表中会添加 created_at updated_at deleted_at id 字段
		UserName string `json:"UserName"` //字段名 user_name
		Age      int    `json:"Age"`
		Addr     string `json:"Addr"`
	}

	UserInfo1 struct {
		Id       uint   `gorm:"primary_key"` //声明 id 字段为主键
		UserName string `json:"UserName"`    //字段名 user_name
		Age      int    `json:"Age"`
		Addr     string `json:"Addr"`
	}

	UserInfo2 struct {
		gorm.Model
		UserName  string `json:"UserName"`
		Age       int    `json:"Age"`
		Addr      string `json:"Addr"`
		tableName string `gorm:"-"` //创建表时会忽略这个字段
	}
)

func (u *UserInfo2) TableName() string { //自定义 UserInfo2 在数据表的表名
	return u.tableName
}
