package callback

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/server/watch_waring"
	"github.com/xykong/loveauth/services/payment/bilibili"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"io/ioutil"
	"net/url"
	"strconv"
)

func init() {
	handlers["/bilibili/callback"] = dealBilibiliPay
	getHandlers["/bilibili/callback"] = dealBilibiliPay
}

type RequestBilibiliPay struct {
	Data model.BilibiliPayCallback `json:"data" form:"data"`
}

func dealBilibiliPay(c *gin.Context) {

	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	urlValue, err := url.ParseQuery(string(bodyBytes))
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"body": string(bodyBytes),
			"err":  err,
		}).Error("dealBilibiliPay parasequery failed.")

		c.Writer.WriteString("parsequeryErr")
		return
	}

	var request RequestBilibiliPay

	err = json.Unmarshal([]byte(urlValue.Get("data")), &request.Data)
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"err":  err,
			"data": urlValue.Get("data"),
		}).Error("dealBilibiliPay Unmarshal failed.")

		c.Writer.WriteString("unmarshalErr")
		return
	}

	sign := bilibili.MakeCallbackSign(urlValue.Get("data"))
	if sign != request.Data.Sign {
		logrus.WithFields(logrus.Fields{
			"sign":    sign,
			"request": request,
		}).Error("dealBilibiliPay check sign failed.")

		c.Writer.WriteString("signError")
		return
	}

	// 订单是否存在
	order := storage.QueryOrderPlacedWithSequence(request.Data.ExtensionInfo)
	if nil == order {
		logrus.WithFields(logrus.Fields{
			"request": request,
		}).Error("dealBilibiliPay QueryOrderPlacedWithSequence failed.")

		c.Writer.WriteString("Order data not exist")
		return
	}

	if order.State > model.OrderStatePrepare {
		// 订单已被处理过
		logrus.WithFields(logrus.Fields{
			"request": request,
			"order":   order,
		}).Info("dealBilibliPay already deal")

		c.Writer.WriteString("success")
		return
	}

	money, _ := strconv.Atoi(request.Data.Money)
	// 校验金额
	if order.Amount != money {
		logrus.WithFields(logrus.Fields{
			"order":   order,
			"request": request,
		}).Error("dealBilibliPay check amount failed.")

		c.Writer.WriteString("dealBilibliPay Amount not match")
		return
	}

	//检测第三方订单是否已经存在
	biliCbOrder := storage.QueryBilibiliCbOrder(request.Data.OrderNo)
	if biliCbOrder != nil {
		logrus.WithFields(logrus.Fields{
			"order":       order,
			"request":     request,
			"biliCbOrder": biliCbOrder,
		}).Error("dealBilibliPay order exists.")

		c.Writer.WriteString("success")
		return
	}

	// 第三方账单入库
	err = storage.Save(storage.PayDatabase(), &request.Data)
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err":     err,
			"request": request,
			"order":   order,
		}).Error("dealBilibiliPay save failed")

		c.Writer.WriteString("Save order failed")
		return
	}

	// 修改订单状态
	if 1 == request.Data.OrderStatus {
		order.State = model.OrderStatePlace
	} else {
		order.State = model.OrderStateFailed
	}
	order.SNSOrderId = request.Data.OrderNo

	err = storage.Save(storage.PayDatabase(), order)
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err":     err,
			"request": request,
			"order":   order,
		}).Error("dealBilibiliPay save failed")

		c.Writer.WriteString("update order failed")
		return
	}

	c.Writer.WriteString("success")

	watch_waring.PaymentWatch(order.GlobalId, order.Vendor, order.Amount)
	return
}
