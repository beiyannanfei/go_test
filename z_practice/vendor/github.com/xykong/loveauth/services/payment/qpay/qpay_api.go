package qpay

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/errors"
	"github.com/xykong/loveauth/settings"
	"github.com/xykong/loveauth/storage/model"
	"github.com/xykong/loveauth/utils"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

//统一下单
type RequestUnifiedorder struct {
	Appid          string `xml:"appid,omitempty" json:"appid,omitempty"`
	MchId          string `xml:"mch_id" json:"mch_id"`
	NonceStr       string `xml:"nonce_str" json:"nonce_str"`
	Sign           string `xml:"sign,omitempty" json:"sign,omitempty"`
	Body           string `xml:"body" json:"body"`
	Attach         string `xml:"attach" json:"attach"`
	OutTradeNo     string `xml:"out_trade_no" json:"out_trade_no"`
	FeeType        string `xml:"fee_type" json:"fee_type"`
	TotalFee       int    `xml:"total_fee" json:"total_fee,string"`
	SpbillCreateIp string `xml:"spbill_create_ip" json:"spbill_create_ip"`
	TradeType      string `xml:"trade_type" json:"trade_type"`
	NotifyUrl      string `xml:"notify_url" json:"notify_url"`
}

//下单返回
type ResponseUnifiedorder struct {
	ReturnCode string `xml:"return_code" json:"return_code"`
	ReturnMsg  string `xml:"return_msg" json:"return_msg"`
	Retcode    string `xml:"retcode" json:"retcode"`
	Retmsg     string `xml:"retmsg" json:"retmsg"`
	Appid      string `xml:"appid" json:"appid"`
	MchId      string `xml:"mch_id" json:"mch_id"`
	Sign       string `xml:"sign" json:"sign"`
	ResultCode string `xml:"result_code" json:"result_code"`
	ErrCode    string `xml:"err_code" json:"err_code"`
	ErrCodeDes string `xml:"err_code_des" json:"err_code_des"`
	NonceStr   string `xml:"nonce_str" json:"nonce_str"`
	TradeType  string `xml:"trade_type" json:"trade_type"`
	PrepayId   string `xml:"prepay_id" json:"prepay_id"`
	CodeUrl    string `xml:"code_url" json:"code_url"`
}

//调起支付
type PullQPay struct {
	TokenId   string `json:"tokenId"` // prepay_id
	NonceStr  string `json:"nonce"`
	Sign      string `json:"sign,omitempty"`
	PaySerial int32  `json:"paySerial"`
	TimeStamp int64  `json:"timeStamp"`
}

func GeneratePaySign2(signMap map[string]interface{}) string {

	var keys []string
	for key, _ := range signMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var strArr []string
	for _, v := range keys {
		value := fmt.Sprintf("%v", signMap[v])
		strArr = append(strArr, v+"="+value)
	}

	qPayKey := settings.GetString("lovepay", "qPay.appkey")
	qPayKey = qPayKey + "&"

	stringSignTemp := strings.Join(strArr, "&")

	logrus.WithFields(logrus.Fields{
		"stringSignTemp": stringSignTemp,
		"qPayKey":        qPayKey,
	}).Info("GeneratePaySign2")

	h := hmac.New(sha1.New, []byte(qPayKey))
	h.Write([]byte(stringSignTemp))

	result := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return result
}

func GeneratePaySign(jsonStr string) string {
	var tempInterface interface{}
	json.Unmarshal([]byte(jsonStr), &tempInterface)

	paramMap, _ := tempInterface.(map[string]interface{})

	var keys []string
	for key, _ := range paramMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var strArr []string
	for _, v := range keys {
		value := fmt.Sprintf("%v", paramMap[v])
		if value == "" {
			continue
		}
		strArr = append(strArr, fmt.Sprintf("%s=%s", v, value))
	}

	qPayKey := settings.GetString("lovepay", "qPay.key")
	strArr = append(strArr, fmt.Sprintf("key=%s", qPayKey))

	stringSignTemp := strings.Join(strArr, "&")
	logrus.WithFields(logrus.Fields{
		"stringSignTemp": stringSignTemp,
	}).Info("qpay GeneratePaySign")

	return strings.ToUpper(fmt.Sprintf("%x", md5.Sum([]byte(stringSignTemp))))
}

func Unifiedorder(ip string, order *model.Order) (*PullQPay, error) {
	mchId := settings.GetString("lovepay", "qPay.mch_id")
	appId := settings.GetString("lovepay", "qPay.appid")
	notifyUrl := settings.GetString("lovepay", "qPay.notify_url")

	request := RequestUnifiedorder{}
	//request.Appid = appId
	request.MchId = mchId
	request.NonceStr = utils.GenerateNonceStr()
	request.Body = fmt.Sprintf("恋世界购买%v-%v-%v", order.GlobalId, order.ShopId, order.Sequence)
	request.Attach = order.Sequence
	request.OutTradeNo = fmt.Sprintf("%v-%v", order.GlobalId, order.Timestamp.Unix())
	request.FeeType = "CNY"
	request.TotalFee = order.Amount
	request.SpbillCreateIp = ip
	request.TradeType = "APP"
	request.NotifyUrl = notifyUrl

	requestJson, _ := json.Marshal(request)
	request.Sign = GeneratePaySign(string(requestJson))

	var b bytes.Buffer
	enc := xml.NewEncoder(&b)
	start := xml.StartElement{}
	start.Name = xml.Name{Local: "xml"}
	enc.EncodeElement(request, start)

	resp, err := http.Post("https://qpay.qq.com/cgi-bin/pay/qpay_unified_order.cgi", "text/xml", bytes.NewBuffer(b.Bytes()))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"request": request,
			"err":     err,
		}).Error("qpay Unifiedorder http failed.")
		return nil, err
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"request": request,
			"err":     err,
		}).Error("qpay Unifiedorder ReadAll failed.")
		return nil, err
	}

	var unifiedorderResponse ResponseUnifiedorder
	if err := xml.Unmarshal(respBody, &unifiedorderResponse); err != nil {
		logrus.WithFields(logrus.Fields{
			"request":  request,
			"err":      err,
			"respBody": string(respBody),
		}).Error("qpay Unifiedorder Unmarshal failed.")
		return nil, err
	}

	logrus.WithFields(logrus.Fields{
		"unifiedorderResponse": unifiedorderResponse,
	}).Info("qpay Unifiedorder ResponseUnifiedorder")

	if unifiedorderResponse.ReturnCode == "FAIL" { //此字段是通信标识，非交易标识，交易是否成功需要查看result_code来判断
		logrus.WithFields(logrus.Fields{
			"request":              request,
			"unifiedorderResponse": unifiedorderResponse,
		}).Error("qpay Unifiedorder ReturnCode FAIL.")

		return nil, errors.NewCodeString(errors.Failed, "qpay unifiedorder failed return_code: %v, return_msg: %v, retcode: %v, retmsg: %v",
			unifiedorderResponse.ReturnCode, unifiedorderResponse.ReturnMsg, unifiedorderResponse.Retcode, unifiedorderResponse.Retmsg)
	}

	if unifiedorderResponse.ResultCode == "FAIL" { //业务结果
		logrus.WithFields(logrus.Fields{
			"request":              request,
			"unifiedorderResponse": unifiedorderResponse,
		}).Error("qpay Unifiedorder ResultCode FAIL.")

		return nil, errors.NewCodeString(errors.Failed, "qpay unifiedorder failed err_code: %v, err_code_des: %v", unifiedorderResponse.ErrCode, unifiedorderResponse.ErrCodeDes)
	}

	payParam := PullQPay{}
	payParam.TokenId = unifiedorderResponse.PrepayId
	payParam.NonceStr = unifiedorderResponse.NonceStr
	payParam.TimeStamp = order.Timestamp.Unix()
	payParam.PaySerial = int32(order.Timestamp.Unix())

	payParam.Sign = GeneratePaySign2(map[string]interface{}{
		"appId":       appId,
		"nonce":       payParam.NonceStr,
		"tokenId":     payParam.TokenId,
		"pubAcc":      "",
		"bargainorId": mchId,
	})

	logrus.WithFields(logrus.Fields{
		"payParam": payParam,
	}).Info("qpay Unifiedorder payParam")

	order.SNSOrderId = request.OutTradeNo //将第三方订单号存入服务器订单

	return &payParam, nil
}

type RequestOrderQuery struct {
	Appid      string `xml:"appid" json:"appid,omitempty"`
	MchId      string `xml:"mch_id" json:"mch_id"`
	OutTradeNo string `xml:"out_trade_no" json:"out_trade_no"`
	NonceStr   string `xml:"nonce_str" json:"nonce_str"`
	Sign       string `xml:"sign" json:"sign,omitempty"`
}

type ResponseOrderQuery struct {
	ReturnCode     string `xml:"return_code" json:"return_code"`
	ReturnMsg      string `xml:"return_msg" json:"return_msg,omitempty"`
	Retcode        string `xml:"retcode" json:"retcode"`
	Retmsg         string `xml:"retmsg" json:"retmsg,omitempty"`
	Appid          string `xml:"appid" json:"appid,omitempty"`
	MchId          string `xml:"mch_id" json:"mch_id"`
	Sign           string `xml:"sign" json:"-"`
	ResultCode     string `xml:"result_code" json:"result_code"`
	ErrCode        string `xml:"err_code" json:"err_code,omitempty"`
	ErrCodeDes     string `xml:"err_code_des" json:"err_code_des,omitempty"`
	NonceStr       string `xml:"nonce_str" json:"nonce_str"`
	DeviceInfo     string `xml:"device_info" json:"device_info,omitempty"`
	TradeType      string `xml:"trade_type" json:"trade_type"`
	TradeState     string `xml:"trade_state" json:"trade_state"`
	BankType       string `xml:"bank_type" json:"bank_type"`
	FeeType        string `xml:"fee_type" json:"fee_type,omitempty"`
	TotalFee       int    `xml:"total_fee" json:"total_fee,string"`
	CashFee        int    `xml:"cash_fee" json:"cash_fee,string"`
	CouponFee      int    `xml:"coupon_fee" json:"coupon_fee,string,omitempty"`
	TransactionId  string `xml:"transaction_id" json:"transaction_id"`
	OutTradeNo     string `xml:"out_trade_no" json:"out_trade_no"`
	Attach         string `xml:"attach" json:"attach,omitempty"`
	TimeEnd        string `xml:"time_end" json:"time_end"`
	TradeStateDesc string `xml:"trade_state_desc" json:"trade_state_desc"`
	Openid         string `xml:"openid" json:"openid,omitempty"`
}

func QPayOrderQuery(order *model.Order) (*ResponseOrderQuery, error) {
	mchId := settings.GetString("lovepay", "qPay.mch_id")
	//appId := settings.GetString("lovepay", "qPay.appid")

	request := RequestOrderQuery{}
	//request.Appid = appId
	request.MchId = mchId
	request.OutTradeNo = order.SNSOrderId
	if order.SNSOrderId == "" {
		request.OutTradeNo = fmt.Sprintf("%v-%v", order.GlobalId, order.Timestamp.Unix())
	}
	request.NonceStr = utils.GenerateNonceStr()

	requestJson, _ := json.Marshal(request)
	request.Sign = GeneratePaySign(string(requestJson))

	var b bytes.Buffer
	enc := xml.NewEncoder(&b)
	start := xml.StartElement{}
	start.Name = xml.Name{Local: "xml"}
	enc.EncodeElement(request, start)

	resp, err := http.Post("https://qpay.qq.com/cgi-bin/pay/qpay_order_query.cgi", "text/xml", bytes.NewBuffer(b.Bytes()))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"request": request,
			"err":     err,
		}).Error("QPayOrderQuery http failed.")
		return nil, err
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"request": request,
			"err":     err,
		}).Error("QPayOrderQuery ReadAll failed.")
		return nil, err
	}

	var orderQueryResponse ResponseOrderQuery
	if err := xml.Unmarshal(respBody, &orderQueryResponse); err != nil {
		logrus.WithFields(logrus.Fields{
			"err":      err,
			"respBody": string(respBody),
		}).Error("QPayOrderQuery Unmarshal failed.")
		return nil, err
	}

	logrus.WithFields(logrus.Fields{
		"orderQueryResponse": orderQueryResponse,
	}).Info("QPayOrderQuery ResponseOrderQuery")

	if orderQueryResponse.ReturnCode == "FAIL" { //此字段是通信标识，非交易标识，交易是否成功需要查看result_code来判断
		logrus.WithFields(logrus.Fields{
			"orderQueryResponse": orderQueryResponse,
		}).Error("WxPayOrderQuery ReturnCode FAIL.")

		return nil, errors.NewCodeString(errors.Failed, "qpay orderquery failed return_code: %v, return_msg: %v, retcode: %v, retmsg: %v",
			orderQueryResponse.ReturnCode, orderQueryResponse.ReturnMsg, orderQueryResponse.Retcode, orderQueryResponse.Retmsg)
	}

	if orderQueryResponse.ResultCode == "FAIL" { //业务结果
		logrus.WithFields(logrus.Fields{
			"orderQueryResponse": orderQueryResponse,
		}).Error("WxPayOrderQuery ResultCode FAIL.")

		return nil, errors.NewCodeString(errors.Failed, "qpayorderquery failed err_code: %v, err_code_des: %v", orderQueryResponse.ErrCode, orderQueryResponse.ErrCodeDes)
	}

	return &orderQueryResponse, nil
}
