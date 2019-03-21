package callback

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/server/watch_waring"
	"github.com/xykong/loveauth/settings"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"github.com/xykong/loveauth/utils/encoding"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

func init() {
	handlers["/huawei/callback"] = dealHuaweiPay
}

type ResponseHuaweiPay struct {
	Body struct {
		//0: 表示成功 1: 验签失败 2: 超时 3: 业务信息错误，比如订单不存在 94: 系统错误
		//95: IO 错误 96: 错误的url 97: 错误的响应 98: 参数错误 99: 其他错误
		Result int `json:"result"` // 返回值
	}
}

func dealHuaweiPay(c *gin.Context) {
	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	requestList := strings.Split(string(bodyBytes), "&")
	var response ResponseHuaweiPay
	var rMap = make(map[string]string) //全量字段
	var qMap = make(map[string]string) //不包含sign signType的字段
	for _, v := range requestList {
		tList := strings.Split(v, "=")
		key := tList[0]
		value := tList[1]
		if key == "sign" || key == "extReserved" || key == "sysReserved" {
			value, _ = url.QueryUnescape(value)
		}
		rMap[key] = value

		if key == "signType" || key == "sign" {
			continue
		}

		qMap[key] = value
	}

	var keys []string
	for key, _ := range qMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var strArr []string
	for _, v := range keys {
		strArr = append(strArr, fmt.Sprintf("%s=%s", v, qMap[v]))
	}

	stringSignTemp := strings.Join(strArr, "&") //最终参与签名的字符串
	logrus.WithFields(logrus.Fields{
		"stringSignTemp": stringSignTemp,
	}).Info("dealHuaweiPay stringSignTemp")

	privateKeyStr := fmt.Sprintf("-----BEGIN PRIVATE KEY-----\r\n%s\r\n-----END PRIVATE KEY-----",
		settings.GetString("lovepay", "huawei.privateKey"))
	privateKey, err := encoding.GetRSAPrivateKey(privateKeyStr)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":        err,
			"requestStr": string(bodyBytes),
		}).Error("dealHuaweiPay GetRSAPrivateKey failed.")

		response.Body.Result = 94
		c.JSON(http.StatusOK, response.Body)
		return
	}

	signature, err := encoding.SignSHA256WithRSA(stringSignTemp, privateKey)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":        err,
			"requestStr": string(bodyBytes),
		}).Error("dealHuaweiPay SignSHA256WithRSA failed.")

		response.Body.Result = 94
		c.JSON(http.StatusOK, response.Body)
		return
	}

	if signature != rMap["sign"] { //验证签名失败
		logrus.WithFields(logrus.Fields{
			"requestStr": string(bodyBytes),
			"signature":  signature,
		}).Error("dealHuaweiPay sign error.")

		response.Body.Result = 1
		c.JSON(http.StatusOK, response.Body)
		return
	}

	// 获取服务器订单信息
	order := storage.QueryOrderPlacedWithSequence(rMap["extReserved"])
	if nil == order {
		logrus.WithFields(logrus.Fields{
			"requestStr": string(bodyBytes),
		}).Error("dealHuaweiPay QueryOrderPlacedWithSequence failed.")

		response.Body.Result = 3
		c.JSON(http.StatusOK, response.Body)
		return
	}

	// 获取订单状态
	if order.State > model.OrderStatePrepare {
		logrus.WithFields(logrus.Fields{
			"order":      order,
			"requestStr": string(bodyBytes),
		}).Error("dealHuaweiPay order State error")

		response.Body.Result = 0
		c.JSON(http.StatusOK, response.Body)
		return
	}

	// 校验订单金额
	requestAmount, _ := strconv.ParseFloat(rMap["amount"], 64)
	if order.Amount != int(requestAmount*100) {
		logrus.WithFields(logrus.Fields{
			"order":      order,
			"requestStr": string(bodyBytes),
		}).Error("dealHuaweiPay amount error.")

		response.Body.Result = 3
		c.JSON(http.StatusOK, response.Body)
		return
	}

	// 订单是否已存在
	huaweiOrder := storage.QueryHuaweiOrder(rMap["orderId"])
	if nil != huaweiOrder {
		logrus.WithFields(logrus.Fields{
			"order":      order,
			"requestStr": string(bodyBytes),
		}).Error("dealHuaweiPay order exists")

		response.Body.Result = 3
		c.JSON(http.StatusOK, response.Body)
		return
	}

	requestJson, err := json.Marshal(rMap)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":        err,
			"order":      order,
			"requestStr": string(bodyBytes),
		}).Error("dealHuaweiPay Marshal failed.")

		response.Body.Result = 3
		c.JSON(http.StatusOK, response.Body)
		return
	}

	var hwPayCbOrder model.HuaweiPayCallback
	err = json.Unmarshal(requestJson, &hwPayCbOrder)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":         err,
			"order":       order,
			"requestStr":  string(bodyBytes),
			"requestJson": string(requestJson),
		}).Error("dealHuaweiPay Unmarshal failed.")

		response.Body.Result = 3
		c.JSON(http.StatusOK, response.Body)
		return
	}

	err = storage.Save(storage.PayDatabase(), &hwPayCbOrder)
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err":          err,
			"order":        order,
			"requestStr":   string(bodyBytes),
			"hwPayCbOrder": hwPayCbOrder,
		}).Error("dealHuaweiPay huawei order save failed")

		response.Body.Result = 95
		c.JSON(http.StatusOK, response.Body)
		return
	}

	if rMap["result"] == "0" {
		order.State = model.OrderStatePlace
	} else {
		order.State = model.OrderStateFailed
	}
	order.SNSOrderId = hwPayCbOrder.OrderId
	//订单再处理
	err = storage.Save(storage.PayDatabase(), order)

	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err":          err,
			"order":        order,
			"requestStr":   string(bodyBytes),
			"hwPayCbOrder": hwPayCbOrder,
		}).Error("dealHuaweiPay order save failed")

		response.Body.Result = 95
		c.JSON(http.StatusOK, response.Body)
		return
	}

	response.Body.Result = 0
	c.JSON(http.StatusOK, response.Body)

	watch_waring.PaymentWatch(order.GlobalId, order.Vendor, order.Amount)
	return
}
