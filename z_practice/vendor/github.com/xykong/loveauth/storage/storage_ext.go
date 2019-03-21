package storage

import (
	"github.com/jinzhu/gorm"
	"github.com/xykong/loveauth/settings"
	"github.com/xykong/loveauth/storage/model"
)

var dbExt *gorm.DB

func ExtDatabase() *gorm.DB {
	return dbExt
}

func InitExtDatabase() {
	//Migrate the schema
	dbExt.AutoMigrate(
		&PlayerRegister{},
		&PlayerPayment{},
	)
}

type PlayerRegister struct {
	TimeKey   int    `json:"timekey" gorm:"column:timekey; index"`
	GlobalId  int64  `json:"uid" gorm:"index"`
	DeviceId  string `json:"device_id" gorm:"type:varchar(191); index"`
	Channel   string `json:"channel" gorm:"type:varchar(100); index"`
	GameSvrId string `json:"gamesvrid" gorm:"type:varchar(20); index"`
}

func (r *PlayerRegister) TableName() string {
	return settings.GetString("tencent", "tlog.mysql.register_table")
}

type PlayerPayment struct {
	TimeKey   int            `json:"timekey" gorm:"column:timekey; index"`
	GlobalId  int64          `json:"uid" gorm:"index"`
	DeviceId  string         `json:"device_id" gorm:"type:varchar(191); index"`
	Amount    float64        `json:"amount"`
	Channel   string         `json:"channel" gorm:"type:varchar(100); index"`
	Platform  model.Platform `json:"platform" gorm:"type:varchar(100); index"`
	Sequence  string         `json:"sequence" gorm:"type:varchar(191)"`
	GameSvrId string         `json:"gamesvrid" gorm:"type:varchar(20); index"`
}

func (r *PlayerPayment) TableName() string {
	return settings.GetString("tencent", "tlog.mysql.payment_table")
}
