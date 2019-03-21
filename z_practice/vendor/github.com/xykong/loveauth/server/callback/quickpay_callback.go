package callback

import (
	"crypto/md5"
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/server/watch_waring"
	"github.com/xykong/loveauth/services/Login/quick_sdk"
	"github.com/xykong/loveauth/settings"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"io/ioutil"
	"net/url"
	"strconv"
)

func init() {
	handlers["/quick/callback"] = dealQuickPay
	getHandlers["/quick/callback"] = dealQuickPay
}

type RequestQuickPay struct {
	NtData  string `form:"nt_data" json:"nt_data"`
	Sign    string `form:"sign" json:"sign"`
	Md5Sing string `form:"md5Sign" json:"md5Sign"`
}

func dealQuickPay(c *gin.Context) {

	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
	}

	value, _ := url.ParseQuery(string(bodyBytes))

	request := RequestQuickPay{
		NtData:  value.Get("nt_data"),
		Sign:    value.Get("sign"),
		Md5Sing: value.Get("md5Sign"),
	}

	//1. 验证签名
	md5Key := settings.GetString("lovepay", "quickSdk.Md5_Key")
	md5Sign := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s%s%s", request.NtData, request.Sign, md5Key))))
	if md5Sign != request.Md5Sing {
		logrus.WithFields(logrus.Fields{
			"request": request,
			"md5Sign": md5Sign,
		}).Error("dealQuickPay verify sign failed.")

		c.Writer.WriteString("SignError")
		return
	}

	//2. 解密
	xmlStr := quick_sdk.DecodeQuickCb(request.NtData)
	logrus.WithFields(logrus.Fields{
		"xmlStr":  xmlStr,
		"request": request,
	}).Debug("dealQuickPay xmlStr")

	//3. 解析xml
	type SkymoonsMessage struct {
		Message model.QuickPayCallback `xml:"message"`
	}

	cbMsg := SkymoonsMessage{}
	err := xml.Unmarshal([]byte(xmlStr), &cbMsg)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":     err,
			"xmlStr":  xmlStr,
			"request": request,
		}).Error("dealQuickPay parse xml failed.")

		c.Writer.WriteString("parseXmlError")
		return
	}

	//4. 获取服务器订单
	order := storage.QueryOrderPlacedWithSequence(cbMsg.Message.GameOrder)
	if order == nil {
		logrus.WithFields(logrus.Fields{
			"cbMsg": cbMsg,
		}).Error("dealQuickPay QueryOrderPlacedWithSequence failed.")
		c.Writer.WriteString("orderNotExists")
		return
	}

	//5. 检测订单状态(当订单状态不是未支付则返回成功，不需要再通知)
	if order.State != model.OrderStatePrepare {
		logrus.WithFields(logrus.Fields{
			"cbMsg": cbMsg,
			"order": order,
		}).Error("dealQuickPay order state error.")

		c.Writer.WriteString("SUCCESS")
		return
	}

	//6. 检测两侧订单金额是否一致
	quickAmount, _ := strconv.ParseFloat(cbMsg.Message.Amount, 64)
	if order.Amount != int(quickAmount*100) { //支付金额与订单金额不一致
		logrus.WithFields(logrus.Fields{
			"cbMsg": cbMsg,
			"order": order,
		}).Error("dealQuickPay amoung not equal.")
		c.Writer.WriteString("amountError")
		return
	}

	//7. 检查订单是否已存在
	quickOrder := storage.QueryQuickPayOrder(cbMsg.Message.GameOrder)
	if quickOrder != nil && quickOrder.GameOrder == cbMsg.Message.GameOrder {
		logrus.WithFields(logrus.Fields{
			"order":  order,
			"cbMsg":  cbMsg,
			"xmlStr": xmlStr,
		}).Warn("dealQuickPay order exists.")

		c.Writer.WriteString("SUCCESS")
		return
	}

	//8. 数据入库
	err = storage.Save(storage.PayDatabase(), &cbMsg.Message)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":   err,
			"cbMsg": cbMsg,
		}).Error("dealQuickPay save failed.")

		c.Writer.WriteString("saveOrderError")
		return
	}

	//9. 修改订单状态
	//status	string	必有		充值状态:0成功, 1失败(为1时 应返回FAILED失败)
	if cbMsg.Message.Status == "0" {
		order.State = model.OrderStatePlace
	} else {
		order.State = model.OrderStateFailed
	}

	order.SNSOrderId = cbMsg.Message.OrderNo

	//订单再处理
	err = storage.Save(storage.PayDatabase(), order)
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"err":   err,
			"order": order,
		}).Error("dealQuickPay save failed.")

		c.Writer.WriteString("saveOrderError")
		return
	}

	c.Writer.WriteString("SUCCESS")
	watch_waring.PaymentWatch(order.GlobalId, order.Vendor, order.Amount)
	return
}
