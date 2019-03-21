package callback

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/server/watch_waring"
	"github.com/xykong/loveauth/services/payment/qpay"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
)

func init() {
	handlers["/qpay/callback"] = dealQPay
}

type ResponseQPay struct {
	ReturnCode string `xml:"return_code" json:"return_code"`
	ReturnMsg  string `xml:"return_msg" json:"return_msg"`
}

func dealQPay(c *gin.Context) {
	var request model.QPayCallback
	if err := c.MustBindWith(&request, binding.XML); err != nil {
		respQPay(c, "FAIL", err.Error())
		return
	}

	//1. 验证签名
	requestJson, _ := json.Marshal(request)
	sign := qpay.GeneratePaySign(string(requestJson))
	if sign != request.Sign {
		logrus.WithFields(logrus.Fields{
			"request":     request,
			"sign":        sign,
			"requestJson": string(requestJson),
		}).Error("dealQPay sign error.")

		respQPay(c, "FAIL", "signError")
		return
	}

	//2. 获取服务器订单
	order := storage.QueryOrderPlacedWithSequence(request.Attach)
	if order == nil {
		logrus.WithFields(logrus.Fields{
			"request": request,
		}).Error("dealQPay QueryOrderPlacedWithSequence failed.")

		respQPay(c, "FAIL", "orderNull")
		return
	}

	//3. 检测订单状态(当订单状态不是未支付则返回成功，不需要再通知)
	if order.State > model.OrderStatePrepare {
		logrus.WithFields(logrus.Fields{
			"request": request,
			"order":   order,
		}).Error("dealQPay order state error.")

		respQPay(c, "SUCCESS", "OK")
		return
	}

	//4. 检测两侧订单金额是否一致
	if order.Amount != request.TotalFee { //支付金额与订单金额不一致
		logrus.WithFields(logrus.Fields{
			"request": request,
			"order":   order,
		}).Error("dealQPay amoung not equal.")

		respQPay(c, "FAIL", "total_fee error")
		return
	}

	//5. 检查订单是否已存在
	qPayOrder := storage.QueryQPayOrder(request.TransactionId)
	if qPayOrder != nil {
		logrus.WithFields(logrus.Fields{
			"qPayOrder": qPayOrder,
			"request":   request,
		}).Warn("deal qPay order exists.")

		respQPay(c, "SUCCESS", "OK")
		return
	}

	//6. 回调数据入库
	err := storage.Save(storage.PayDatabase(), &request)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":     err,
			"request": request,
		}).Error("deal qPay save failed.")

		respQPay(c, "FAIL", "saveOrderError")
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
		}).Error("deal qPay save failed.")

		respQPay(c, "FAIL", "saveOrderError")
		return
	}

	respQPay(c, "SUCCESS", "OK")
	watch_waring.PaymentWatch(order.GlobalId, order.Vendor, order.Amount)
	return
}

func respQPay(c *gin.Context, code string, msg string) {
	resp := ResponseWxPay{}
	resp.ReturnCode = code
	resp.ReturnMsg = msg

	var b bytes.Buffer
	enc := xml.NewEncoder(&b)
	start := xml.StartElement{}
	start.Name = xml.Name{Local: "xml"}
	enc.EncodeElement(resp, start)

	c.Writer.Write(b.Bytes())
	return
}
