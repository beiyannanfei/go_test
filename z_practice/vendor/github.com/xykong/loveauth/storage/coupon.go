package storage

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type Coupon struct {
	gorm.Model

	CouponId  uint32    `json:"CouponId" gorm:"index"`
	Coupon    string    `json:"Coupon"`
	Used      time.Time `json:"Used"`
	Count     uint64    `json:"Count"`
	GlobalId  uint64    `json:"GlobalId"`
	tableName string    `gorm:"-"` // 忽略这个字段
}

// 设置User的表名为`profiles`
func (r *Coupon) TableName() string {
	return r.tableName
}

func (r *Coupon) SetName(name string) {
	r.tableName = "coupons_" + name
}

func (r *Coupon) BulkInsert(unsavedRows []*Coupon) error {

	valueStrings := make([]string, 0, len(unsavedRows))
	valueArgs := make([]interface{}, 0, len(unsavedRows)*3)
	for _, post := range unsavedRows {
		valueStrings = append(valueStrings, "(?, ?, ?)")
		valueArgs = append(valueArgs, post.CouponId)
		valueArgs = append(valueArgs, post.Coupon)
		valueArgs = append(valueArgs, post.Used)
	}

	stmt := fmt.Sprintf("INSERT INTO %s (coupon_id, coupon, used) VALUES %s",
		r.tableName, strings.Join(valueStrings, ","))
	_, err := dbAuth.DB().Exec(stmt, valueArgs...)

	return err
}

func (r *Coupon) QueryCoupon() *Coupon {

	if err := dbAuth.Where("coupon_id = ?", r.CouponId).First(r).Error; err != nil {
		logrus.WithFields(logrus.Fields{
			"couponId": r.CouponId,
			"error":    err.Error(),
		}).Info("QueryCoupon data not found.")

		return nil
	}

	logrus.WithFields(logrus.Fields{
		"couponId": r.CouponId,
		"data":     r,
	}).Info("QueryCoupon data found.")

	return r
}

func (r *Coupon) MarkCoupon() *Coupon {

	if err := dbAuth.Model(r).Where("coupon_id = ?", r.CouponId).Update("used", time.Now()).Error; err != nil {
		logrus.WithFields(logrus.Fields{
			"couponId": r.CouponId,
			"error":    err.Error(),
		}).Info("MarkCoupon failed.")

		return nil
	}

	logrus.WithFields(logrus.Fields{
		"couponId": r.CouponId,
		"data":     r,
	}).Info("MarkCoupon ok.")

	return r
}
