package callback

import (
	"encoding/json"
	"encoding/xml"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/server/watch_waring"
	"github.com/xykong/loveauth/services/payment/wechat_pay"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"io/ioutil"
)

func init() {
	handlers["/wxpay/callback"] = dealWechatPay
}

type ResponseWxPay struct {
	ReturnCode string `xml:"return_code" json:"return_code"`
	ReturnMsg  string `xml:"return_msg" json:"return_msg"`
}

func dealWechatPay(c *gin.Context) {
	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
	}

	var request model.WechatPayCallback
	err := xml.Unmarshal(bodyBytes, &request)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":       err,
			"bodyBytes": string(bodyBytes),
		}).Error("dealWechatPay Unmarshal failed.")
		respWxPay(c, "FAIL", err.Error())
		return
	}

	//0. 验证消息是否成功
	if request.ReturnCode != "SUCCESS" || request.ResultCode != "SUCCESS" {
		logrus.WithFields(logrus.Fields{
			"request": request,
		}).Error("dealWechatPay ReturnCode or ResultCode error.")

		respWxPay(c, "FAIL", "code-error")
		return
	}

	//1. 验证签名
	requestJson, _ := json.Marshal(request)
	sign := wechat_pay.GeneratePaySign(string(requestJson))
	if sign != request.Sign {
		logrus.WithFields(logrus.Fields{
			"request":     request,
			"sign":        sign,
			"requestJson": string(requestJson),
		}).Error("dealWechatPay sign error.")

		respWxPay(c, "FAIL", "signError")
		return
	}

	//2. 获取服务器订单
	order := storage.QueryOrderPlacedWithSequence(request.Attach)
	if order == nil {
		logrus.WithFields(logrus.Fields{
			"request": request,
		}).Error("dealWechatPay QueryOrderPlacedWithSequence failed.")

		respWxPay(c, "FAIL", "orderNull")
		return
	}

	//3. 检测订单状态(当订单状态不是未支付则返回成功，不需要再通知)
	if order.State > model.OrderStatePrepare {
		logrus.WithFields(logrus.Fields{
			"request": request,
			"order":   order,
		}).Error("dealWechatPay order state error.")

		respWxPay(c, "SUCCESS", "OK")
		return
	}

	//4. 检测两侧订单金额是否一致
	if order.Amount != request.TotalFee { //支付金额与订单金额不一致
		logrus.WithFields(logrus.Fields{
			"request": request,
			"order":   order,
		}).Error("dealWechatPay amoung not equal.")

		respWxPay(c, "FAIL", "total_fee error")
		return
	}

	//5. 检查订单是否已存在
	wxPayOrder := storage.QueryWxPayOrder(request.TransactionId)
	if wxPayOrder != nil {
		logrus.WithFields(logrus.Fields{
			"wxPayOrder": wxPayOrder,
			"request":    request,
		}).Warn("dealWechatPay order exists.")

		respWxPay(c, "SUCCESS", "OK")
		return
	}

	//6. 回调数据入库
	err = storage.Save(storage.PayDatabase(), &request)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":     err,
			"request": request,
		}).Error("dealWechatPay save failed.")

		respWxPay(c, "FAIL", "saveOrderError")
		return
	}

	//9. 修改订单状态
	order.State = model.OrderStatePlace
	order.SNSOrderId = request.TransactionId
	err = storage.Save(storage.PayDatabase(), order)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":   err,
			"order": order,
		}).Error("dealWechatPay save failed.")

		respWxPay(c, "FAIL", "saveOrderError")
		return
	}

	respWxPay(c, "SUCCESS", "OK")
	watch_waring.PaymentWatch(order.GlobalId, order.Vendor, order.Amount)
	return
}

func respWxPay(c *gin.Context, code string, msg string) {
	resp := ResponseWxPay{}
	resp.ReturnCode = code
	resp.ReturnMsg = msg

	xmlStr, _ := xml.Marshal(resp)
	c.Writer.Write(xmlStr)
	return
}
