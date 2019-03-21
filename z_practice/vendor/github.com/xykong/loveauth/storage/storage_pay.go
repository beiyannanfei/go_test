package storage

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/storage/model"
)

var dbPay *gorm.DB

func PayDatabase() *gorm.DB {
	return dbPay
}

func InitPayDatabase() {

	//Migrate the schema
	dbPay.AutoMigrate(
		&model.Order{},
		&model.Inventory{},
		&model.Activity{},
		&model.DoMidasCallbackReq{},
		&model.QuickPayCallback{},
		&model.BilibiliPayCallback{},
		&model.WechatPayCallback{},
		&model.HuaweiPayCallback{},
		&model.DoAlipayCallbackReq{},
		&model.VivoQueryOrderResponse{},
		&model.DouyinCallBackRequest{},
		&model.MgtvCallBackRequest{},
		&model.RepireOrderRequest{},
		&model.QPayCallback{},
	)
}

func QueryOrderPlacedWithSequence(sequence string) *model.Order {

	var data model.Order

	if err := dbPay.Where("sequence = ?", sequence).First(&data).Error; err != nil {

		logrus.WithFields(logrus.Fields{
			"sequence": sequence,
			"error":    err.Error(),
		}).Info("QueryOrderPlacedWithSequence data not found.")

		return nil
	}

	logrus.WithFields(logrus.Fields{
		"sequence": sequence,
		"data":     data,
	}).Info("QueryOrderPlacedWithSequence")

	return &data
}

func QueryOrderPlacedWithOrderNumber(orderNumber string, state model.OrderState) *model.Order {

	var data model.Order

	if err := dbPay.Where("sns_order_id = ? AND state = ?", orderNumber, state).First(&data).Error; err != nil {

		logrus.WithFields(logrus.Fields{
			"orderNumber": orderNumber,
			"error":       err.Error(),
		}).Info("QueryOrderPlacedWithOrderNumber data not found.")

		return nil
	}

	logrus.WithFields(logrus.Fields{
		"orderNumber": orderNumber,
		"data":        data,
	}).Info("QueryOrderPlacedWithOrderNumber")

	return &data
}

func QueryOrderPlacedWithAppStoreNumber(orderNumber string) *model.Order {

	var data model.Order

	if err := dbPay.Where("sns_order_id = ?", orderNumber).First(&data).Error; err != nil {

		logrus.WithFields(logrus.Fields{
			"sns_order_id": orderNumber,
			"error":        err.Error(),
		}).Info("QueryOrderPlacedWithAppStoreNumber data not found.")

		return nil
	}

	logrus.WithFields(logrus.Fields{
		"sns_order_id": orderNumber,
		"data":         data,
	}).Info("QueryOrderPlacedWithAppStoreNumber")

	return &data
}

func UpdateOrderState(sequence string, state model.OrderState) error {
	var order model.Order
	if err := dbPay.Model(&order).Where("sequence = ?", sequence).Update("state", state).Error; err != nil {
		logrus.WithFields(logrus.Fields{
			"sequence": sequence,
			"state":    state,
		}).Error("UpdateOrderState failed.")

		return err
	}

	return nil
}

func QueryOrderPlaced(globalId int64) *model.Order {

	var data model.Order

	if err := dbPay.Where("global_id = ? and state = ?", globalId, model.OrderStatePlace).
		Order("timestamp desc").First(&data).Error; err != nil {
		logrus.WithFields(logrus.Fields{
			"globalId": globalId,
			"error":    err.Error(),
		}).Info("QueryOrderPlaced data not found.")

		return nil
	}

	logrus.WithFields(logrus.Fields{
		"globalId": globalId,
		"data":     data,
	}).Info("QueryOrderPlaced")

	return &data
}

func QueryOrderByGlobalId(globalId int64, start, end time.Time) []*model.Order {

	var data []*model.Order

	if !start.IsZero() && !end.IsZero() {

		if err := dbPay.Where("global_id = ? AND type < 3 AND timestamp > ? AND timestamp < ?", globalId, start.String(), end.String()).Find(&data).Error; err != nil {

			logrus.WithFields(logrus.Fields{
				"globalId": globalId,
				"error":    err.Error(),
			}).Info("QueryOrderByGlobalId data not found.")

			return nil
		}
	} else {

		if err := dbPay.Where("global_id = ? AND type < 3", globalId).Find(&data).Error; err != nil {

			logrus.WithFields(logrus.Fields{
				"globalId": globalId,
				"error":    err.Error(),
			}).Info("QueryOrderByGlobalId data not found.")

			return nil
		}
	}

	return data
}

func QuerySameShopIdCount(globalId int64, shopId int, state model.OrderState, test time.Time) int {

	var count int

	dbPay.Table("orders").Where("global_id = ? and shop_id = ? and state = ? and timestamp < ?", globalId, shopId, state, test.String()).Count(&count)

	return count
}

func QuerySameShopIdCountPre(globalId int64, shopId int, test time.Time) int {

	var count int

	dbPay.Table("orders").Where("global_id = ? and shop_id = ? and state in (2,4) and timestamp < ?", globalId, shopId, test.String()).Count(&count)

	return count
}

func QueryShopIdAvailableCount(globalId int64, shopId int, state model.OrderState, start, end int64) int {

	var countstate struct {
		Count int
	}

	if start != 0 && end != 0 {

		startTime := time.Unix(start, 0)
		endTime := time.Unix(end, 0)
		dbPay.Raw("select sum(num) as count from orders where global_id = ? and shop_id = ? and state = ? and timestamp > ? and timestamp < ?", globalId, shopId, state, startTime.String(), endTime.String()).Scan(&countstate)
	} else {

		dbPay.Raw("select sum(num) as count from orders where global_id = ? and shop_id = ? and state = ?", globalId, shopId, state).Scan(&countstate)
	}

	return countstate.Count
}

func QueryQuickPayOrder(gameOrder string) *model.QuickPayCallback {
	data := model.QuickPayCallback{}

	if err := dbPay.Where("game_order = ?", gameOrder).First(&data).Error; err != nil {
		logrus.WithFields(logrus.Fields{
			"err":       err,
			"gameOrder": gameOrder,
		}).Info("QueryQuickPayOrder data not found.")
		return nil
	}

	return &data
}

func QueryBilibiliPayOrderWithSequence(outTradeNo string) *model.BilibiliPayCallback {
	var data *model.BilibiliPayCallback

	if err := dbPay.Where("out_trade_no = ?", outTradeNo).First(&data).Error; nil != err {
		logrus.WithFields(logrus.Fields{
			"err":        err,
			"outTradeNo": outTradeNo,
		}).Info("QueryBilibiliPayOrderWithSequence data not found")
		return nil
	}
	return data
}

func QueryBilibiliCbOrder(sequence string) *model.BilibiliPayCallback {
	var data *model.BilibiliPayCallback

	if err := dbPay.Where("extension_info = ?", sequence).First(&data).Error; nil != err {
		logrus.WithFields(logrus.Fields{
			"err":      err,
			"sequence": sequence,
		}).Info("QueryBilibiliPayOrder data not found")

		return nil
	}

	return data
}

func QueryWxPayOrder(transactionId string) *model.WechatPayCallback {
	data := model.WechatPayCallback{}

	if err := dbPay.Where("transaction_id = ?", transactionId).First(&data).Error; err != nil {
		logrus.WithFields(logrus.Fields{
			"err":           err,
			"transactionId": transactionId,
		}).Info("QueryWxPayOrder data not found.")
		return nil
	}

	return &data
}

func QueryWxPayOrderWithSequence(sequence string) *model.WechatPayCallback {
	data := model.WechatPayCallback{}

	if err := dbPay.Where("attach = ?", sequence).First(&data).Error; err != nil {
		logrus.WithFields(logrus.Fields{
			"err":      err,
			"sequence": sequence,
		}).Info("QueryWxPayOrderWithSequence data not found.")
		return nil
	}

	return &data
}

func QueryQPayOrder(transactionId string) *model.QPayCallback {
	data := model.QPayCallback{}

	if err := dbPay.Where("transaction_id = ?", transactionId).First(&data).Error; err != nil {
		logrus.WithFields(logrus.Fields{
			"err":           err,
			"transactionId": transactionId,
		}).Info("QueryQPayOrder data not found.")
		return nil
	}

	return &data
}

func QueryQPayOrderWithSequence(sequence string) *model.QPayCallback {
	data := model.QPayCallback{}

	if err := dbPay.Where("attach = ?", sequence).First(&data).Error; err != nil {
		logrus.WithFields(logrus.Fields{
			"err":      err,
			"sequence": sequence,
		}).Info("QueryQPayOrderWithSequence data not found.")
		return nil
	}

	return &data
}

func QueryHuaweiOrder(transactionId string) *model.HuaweiPayCallback {
	data := model.HuaweiPayCallback{}

	if err := dbPay.Where("order_id = ?", transactionId).First(&data).Error; err != nil {
		logrus.WithFields(logrus.Fields{
			"err":           err,
			"transactionId": transactionId,
		}).Info("QueryHuaweiPayOrder data not found.")
		return nil
	}

	return &data
}

func QueryDouyinOrder(orderNo string) *model.DouyinCallBackRequest {
	data := model.DouyinCallBackRequest{}

	if err := dbPay.Where("trade_no = ?", orderNo).First(&data).Error; err != nil {
		logrus.WithFields(logrus.Fields{
			"err":      err,
			"trade_no": orderNo,
		}).Info("QueryDouyinOrder data not found.")
		return nil
	}

	return &data
}

func QueryMgtvOrder(orderNo string) *model.MgtvCallBackRequest {
	data := model.MgtvCallBackRequest{}

	if err := dbPay.Where("business_order_id = ?", orderNo).First(&data).Error; err != nil {

		logrus.WithFields(logrus.Fields{
			"err":      err,
			"trade_no": orderNo,
		}).Info("QueryMangguotvOrder data not found.")

		return nil
	}

	return &data
}

func QueryHuaweiOrderWithSequence(sequence string) *model.HuaweiPayCallback {
	data := model.HuaweiPayCallback{}

	if err := dbPay.Where("ext_reserved = ?", sequence).First(&data).Error; err != nil {
		logrus.WithFields(logrus.Fields{
			"err":      err,
			"sequence": sequence,
		}).Info("QueryHuaweiPayOrderWithSequence data not found.")
		return nil
	}

	return &data
}

func QueryRepireOrderRequestWithGlobalId(globalId int64, start, end time.Time) []*model.RepireOrderRequest {

	var data []*model.RepireOrderRequest
	if !start.IsZero() && !end.IsZero() {

		if err := dbPay.Where("global_id = ? AND repire_time > ? AND repire_time < ?", globalId, start.String(), end.String()).Find(&data).Error; err != nil {

			logrus.WithFields(logrus.Fields{
				"err":      err,
				"globalId": globalId,
			}).Info("QueryRepireOrderRequestWithGlobalId data not found.")

			return nil
		}
	} else {

		if err := dbPay.Where("global_id = ?", globalId).Find(&data).Error; err != nil {

			logrus.WithFields(logrus.Fields{
				"err":      err,
				"globalId": globalId,
			}).Info("QueryRepireOrderRequestWithGlobalId data not found.")

			return nil
		}

	}

	return data
}

func QueryRepireOrderRequestWithOperator(operator string, start, end time.Time) []*model.RepireOrderRequest {

	var data []*model.RepireOrderRequest
	if !start.IsZero() && !end.IsZero() {

		if err := dbPay.Where("operator = ? AND repire_time > ? AND repire_time < ?", operator, start.String(), end.String()).Find(&data).Error; err != nil {
			logrus.WithFields(logrus.Fields{
				"err":      err,
				"operator": operator,
			}).Info("QueryRepireOrderRequestWithOperator data not found.")

			return nil
		}
	} else {

		if err := dbPay.Where("operator = ?", operator).Find(&data).Error; err != nil {

			logrus.WithFields(logrus.Fields{
				"err":      err,
				"operator": operator,
			}).Info("QueryRepireOrderRequestWithOperator data not found.")

			return nil
		}

	}

	return data
}
