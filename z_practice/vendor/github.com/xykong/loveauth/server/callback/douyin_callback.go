package callback

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/server/watch_waring"
	"github.com/xykong/loveauth/services/payment/douyin"
	"github.com/xykong/loveauth/settings"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

func init() {
	handlers["/douyin/callback"] = DouyinCallBack
	getHandlers["/douyin/callback"] = DouyinCallBack
}

func DouyinCallBack(c *gin.Context) {

	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {

		c.String(http.StatusOK, "douyin io read err "+err.Error())

		return
	}

	value, err := url.ParseQuery(string(reqBody))
	if err != nil {

		c.String(http.StatusOK, "douyin parseQuery err "+err.Error())

		return
	}

	payKey := settings.GetString("lovepay", "douyin.payKey")
	if !douyin.CheckSign(value, payKey) {

		c.String(http.StatusOK, "douyin checkoutsign err")

		return
	}

	request := model.DouyinCallBackRequest{
		NotifyId:    value.Get("notify_id"),
		NotifyType:  value.Get("notify_type"),
		NotifyTime:  value.Get("notify_time"),
		TradeStatus: value.Get("trade_status"),
		Way:         value.Get("way"),
		ClientId:    value.Get("client_id"),
		OutTradeNo:  value.Get("out_trade_no"),
		TradeNo:     value.Get("trade_no"),
		PayTime:     value.Get("pay_time"),
		TotalFee:    value.Get("total_fee"),
		BuyerId:     value.Get("buyer_id"),
		TtSign:      value.Get("tt_sign"),
		TtSignType:  value.Get("tt_sign_type"),
	}

	order := storage.QueryOrderPlacedWithSequence(request.OutTradeNo)
	if order == nil {

		logrus.WithFields(logrus.Fields{
			"request": request,
			"order":   order,
		}).Error("douyin callback failed.")

		c.String(http.StatusOK, "fail")

		return
	}

	amount, _ := strconv.Atoi(request.TotalFee)
	if order.Amount != amount {

		logrus.WithFields(logrus.Fields{
			"request": request,
			"order":   order,
		}).Error("douyin callback failed amount not equal.")

		c.String(http.StatusOK, "fail")

		return
	}

	douyinCallBack := storage.QueryDouyinOrder(request.TradeNo)
	if douyinCallBack != nil {

		c.String(http.StatusOK, "success")

		return
	}

	storage.Insert(storage.PayDatabase(), &request)

	if order.State >= model.OrderStatePlace {

		c.String(http.StatusOK, "success")

		return
	}

	order.SNSOrderId = request.TradeNo
	order.State = model.OrderStatePlace

	//订单再处理
	err = storage.Save(storage.PayDatabase(), order)
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"request": request,
			"error":   err,
		}).Error("douyin callback failed.")

		c.String(http.StatusOK, "fail")

		return
	}

	c.String(http.StatusOK, "success")

	watch_waring.PaymentWatch(order.GlobalId, order.Vendor, order.Amount)
}
