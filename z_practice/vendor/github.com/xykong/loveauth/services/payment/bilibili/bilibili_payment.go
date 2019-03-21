package bilibili

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/errors"
	"github.com/xykong/loveauth/services/Login/bilibili"
	"github.com/xykong/loveauth/settings"
	"github.com/xykong/loveauth/storage/model"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
)

type BilibiliOrder struct {
	GameMoney  int    `json:"game_money,string"`    //游戏内货币
	Money      int    `json:"money,string"`         //本次交易金额，单位：分（注意，total_fee的值必须为整数，并且在1~100000之间)
	NotifyUrl  string `json:"notify_url"`           //支付回调地址
	OutTradeNo string `json:"out_trade_no"`         //商户订单号，8-32位字符，用于对账用
	OrderSign  string `json:"order_sign,omitempty"` //订单参数签名，请在服务端完成订单参数签名
}

func GeneratePaySign(biliOrder BilibiliOrder) string {
	var data = fmt.Sprintf("%s%s%s%s%s", strconv.Itoa(biliOrder.GameMoney), strconv.Itoa(biliOrder.Money),
		biliOrder.NotifyUrl, biliOrder.OutTradeNo, settings.GetString("lovepay", "bilibili.secretKey"))

	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}

func MakeOrderSign(order *model.Order) *BilibiliOrder {
	request := BilibiliOrder{}
	request.GameMoney = order.Amount / 10
	request.Money = order.Amount
	//request.NotifyUrl = settings.GetString("lovepay", "bilibili.callbackUrl")
	request.NotifyUrl = ""
	request.OutTradeNo = fmt.Sprintf("%v-%v", order.GlobalId, order.Timestamp.Unix())
	request.OrderSign = GeneratePaySign(request)

	return &request
}

func MakeCallbackSign(jsonStr string) string {
	//商户需要对参数进行签名校验，去掉sign参数,将其他参数按照数组键值顺序升序排列,再把所有数组值连接起来，形成的字符串末尾拼接上
	// 开放平台所提供的密钥secret_key即服务端appkey，之后整体做MD5，然后全部转成小写，将产生的加密串与sign进行对比
	var tempInterface interface{}
	json.Unmarshal([]byte(jsonStr), &tempInterface)

	paramMap, _ := tempInterface.(map[string]interface{})

	var keys [] string
	for key, _ := range paramMap {
		if key != "sign" {
			keys = append(keys, key)
		}
	}

	sort.Strings(keys)
	data := ""
	for _, v := range keys {
		value := fmt.Sprintf("%v", paramMap[v])
		data = fmt.Sprintf("%s%s", data, value)
	}

	data = fmt.Sprintf("%s%s", data, settings.GetString("lovepay", "bilibili.secretKey"))

	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}

type DoQueryOrderRsp struct {
	Uid           int    `json:"uid"`
	UserName      string `json:"username"`       // 进行付款的用户账号
	PayTime       int    `json:"pay_time"`       // 支付订单时间，时间戳
	PayMoney      int    `json:"pay_money"`      // 支付金额，单位为“分”
	OrderNo       string `json:"order_no"`       // b站sdk服务器的订单号
	OutTradeNo    string `json:"out_trade_no"`   // 游戏cp厂商支付订单号
	ServerId      int    `json:"server_id"`      // 游戏区服Id
	Subject       string `json:"subject"`        // 订单名称
	Remark        string `json:"remark"`         // 订单备注信息
	OrderStatus   int    `json:"order_status"`   // 订单状态  1=完成 2=失败 3=处理中
	NotifyStatus  int    `json:"notify_status"`  // 异步通知状态 0=未完成 1=已经完成 2=异步通知失败
	ExtensionInfo string `json:"extension_info"` // 额外信息，联运商返回的扩展信息
	Timestamp     int    `json:"timestamp"`      // 对应request的时间戳，秒
	Code          int    `json:"code"`           // 状态码
	ErrMessage    string `json:"message"`        // 错误信息
}

func DoQueryBiliOrder(url string, orderNo string, uid int) (*DoQueryOrderRsp, error) {
	p := bilibili.GetCommonReqParams(int32(uid))
	p.Add("order_no", orderNo)
	sign := bilibili.GetSign(p)
	p.Add("sign", sign)

	resp, err := http.PostForm(url, p)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
			"p":   p.Encode(),
			"url": url,
		}).Error("DoQueryBiliOrder PostForm failed.")

		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
			"p":   p.Encode(),
			"url": url,
		}).Error("DoQueryBiliOrder ReadAll failed.")

		return nil, err
	}

	var data DoQueryOrderRsp
	if err := json.Unmarshal(respBody, &data); err != nil {
		logrus.WithFields(logrus.Fields{
			"err":      err,
			"p":        p.Encode(),
			"url":      url,
			"respBody": string(respBody),
		}).Error("DoQueryBiliOrder Unmarshal failed.")

		return nil, err
	}

	if 0 != data.Code {
		logrus.WithFields(logrus.Fields{
			"url":  url,
			"p":    p.Encode(),
			"data": data,
		}).Error("DoQueryBiliOrder code failed.")

		return nil, errors.NewCodeString(errors.Failed, "bilibili DoQueryBiliOrder response code: %v err, errmsg: %v", data.Code, data.ErrMessage)
	}

	return &data, nil
}

func QueryBiliOrder(orderNo string, uid int) (*DoQueryOrderRsp, error) {
	urls := settings.GetStringSlice("lovepay", "bilibili.url")
	for i := 0; i < len(urls); i += 1 {
		url := urls[i] + "/api/server/query.pay.order"
		response, err := DoQueryBiliOrder(url, orderNo, uid)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"orderNo": orderNo,
				"uid":     uid,
				"url":     urls[i],
			}).Error("bilibili QueryOrder failed.")

			if _, ok := err.(*errors.Type); ok { //认证失败错误
				return nil, err
			}

			//非认证失败，切线再次验证
			continue
		}

		return response, nil
	}

	return nil, errors.NewCodeString(errors.Failed, "bilibili query order all domain failed.")
}
