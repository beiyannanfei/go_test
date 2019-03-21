package payment_v2

import (
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/server/auth"
	"github.com/xykong/loveauth/settings"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/utils/log"
	"time"
)

func LogCurrencyPay(globalId int64, shopItem *ShopItem, payPlat int, thirdOrderId string, orderId string, result int) {
	now := time.Now()
	eventTime := now.Format("2006-01-02 15:04:05")
	format := settings.GetString("tencent", "tlog.format")

	logrus.Info("start log currencypay")
	userInfo, err := storage.QueryProfile(globalId)
	gameSvrId, gameAppID := auth.GetVendorConfigSetting(userInfo.Vendor)
	if nil != err {
		return
	}

	moneyCount := 0.00
	shopId := 0
	productId := ""
	if nil != shopItem {
		shopId = shopItem.ShopId
		productId = shopItem.ProductId
		moneyCount = shopItem.PriceValue
		if shopItem.DiscountPriceValue > 0 {
			moneyCount = shopItem.DiscountPriceValue
		}
	}

	logrus.WithFields(logrus.Fields{
		"FlowName":     "CurrencyPayFlow",
		"GameSvrId":    gameSvrId,
		"EventTime":    eventTime,
		"GameAppId":    gameAppID,
		"PlatId":       auth.GetPlatId(userInfo.Platform),
		"LoginChannel": userInfo.Auth.Channel,
		"OpenId":       userInfo.Auth.OpenId,
		"RoleId":       globalId,
		"Format":       format,
		"Sequence":     0,
		"ShopId":       shopId,
		"ProductId":    productId,
		"PayStep":      4, // Auth通知Gs验证结果
		"PayPlat":      payPlat,
		"ThirdOrderId": thirdOrderId,
		"OrderId":      orderId,
		"MoneyCount":   moneyCount,
		"Result":       result,
	}).Info("Tlog CurrencyPayFlow")

	if nil != log.BILogger {
		log.BILogger.Printf("%s,%s,%s,%s,%d,%s,%s,%d,%s,%d,%d,%s,%d,%d,%s,%s,%.2f,%d",
			"CurrencyPayFlow", gameSvrId, eventTime, gameAppID, auth.GetPlatId(userInfo.Platform),
			userInfo.Auth.Channel, userInfo.Auth.OpenId, globalId, format, 0, shopId, productId,
			4, payPlat, thirdOrderId, orderId, moneyCount, result)
	} else {
		logrus.WithFields(logrus.Fields{"reason": "BILogger is nil", "sequence": orderId}).Info("LogCurrencyPayFlow")
	}

	if true == settings.GetBool("tencent", "tlog.mysql.enable_ext") && 1 == result {
		// todo 只记录交易成功的订单
		payment := new(storage.PlayerPayment)
		payment.GlobalId = globalId
		payment.DeviceId = userInfo.Auth.DeviceId
		payment.TimeKey = int(now.Unix())
		payment.Channel = userInfo.Auth.Channel
		payment.Platform = userInfo.Platform
		payment.Amount = moneyCount
		payment.Sequence = orderId
		payment.GameSvrId = gameSvrId
		err = storage.Insert(storage.ExtDatabase(), payment)
		if nil != err {
			logrus.WithFields(logrus.Fields{
				"err": err.Error(),
			}).Error("LogCurrencyPay")
		}
	}
}
