package wechat_pay

import (
	"time"
	"fmt"
	"sort"
	"github.com/xykong/loveauth/settings"
	"strings"
	"crypto/md5"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/storage/model"
	"strconv"
	"net/http"
	"encoding/xml"
	"bytes"
	"io/ioutil"
	"github.com/xykong/loveauth/errors"
	"encoding/json"
	"github.com/xykong/loveauth/utils"
)

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
		strArr = append(strArr, fmt.Sprintf("%s=%s", v, paramMap[v]))
	}

	wxPayKey := settings.GetString("lovepay", "wechatPay.key")
	strArr = append(strArr, fmt.Sprintf("key=%s", wxPayKey))

	stringSignTemp := strings.Join(strArr, "&")
	logrus.WithFields(logrus.Fields{
		"stringSignTemp": stringSignTemp,
	}).Info("wxpay GeneratePaySign")

	return strings.ToUpper(fmt.Sprintf("%x", md5.Sum([]byte(stringSignTemp))))
}

//统一下单
type RequestUnifiedorder struct {
	Appid          string `xml:"appid" json:"appid"`
	MchId          string `xml:"mch_id" json:"mch_id"`
	NonceStr       string `xml:"nonce_str" json:"nonce_str"`
	Body           string `xml:"body" json:"body"`
	OutTradeNo     string `xml:"out_trade_no" json:"out_trade_no"`
	TotalFee       int    `xml:"total_fee" json:"total_fee,string"`
	SpbillCreateIp string `xml:"spbill_create_ip" json:"spbill_create_ip"`
	NotifyUrl      string `xml:"notify_url" json:"notify_url"`
	TradeType      string `xml:"trade_type" json:"trade_type"`
	Attach         string `xml:"attach" json:"attach"`
	Sign           string `xml:"sign" json:"sign,omitempty"`
}

//下单返回
type ResponseUnifiedorder struct {
	ReturnCode string `xml:"return_code" json:"return_code"`
	ReturnMsg  string `xml:"return_msg" json:"return_msg"`
	Appid      string `xml:"appid" json:"appid"`
	MchId      string `xml:"mch_id" json:"mch_id"`
	NonceStr   string `xml:"nonce_str" json:"nonce_str"`
	Sign       string `xml:"sign" json:"sign"`
	ResultCode string `xml:"result_code" json:"result_code"`
	ErrCode    string `xml:"err_code" json:"err_code"`
	ErrCodeDes string `xml:"err_code_des" json:"err_code_des"`
	TradeType  string `xml:"trade_type" json:"trade_type"`
	PrepayId   string `xml:"prepay_id" json:"prepay_id"`
}

//调起支付
type PullPay struct {
	Appid     string `json:"appid"`
	Partnerid string `json:"partnerid"`
	Prepayid  string `json:"prepayid"`
	Package   string `json:"package"`
	Noncestr  string `json:"noncestr"`
	Timestamp string `json:"timestamp"`
	Sign      string `json:"sign,omitempty"`
}

func Unifiedorder(ip string, order *model.Order) (*PullPay, error) {
	mchId := settings.GetString("lovepay", "wechatPay.mch_id")
	appId := settings.GetString("lovepay", "wechatPay.appid")
	notifyUrl := settings.GetString("lovepay", "wechatPay.notify_url")

	request := RequestUnifiedorder{}
	request.Appid = appId
	request.MchId = mchId
	request.NonceStr = utils.GenerateNonceStr()
	request.Body = fmt.Sprintf("恋世界购买%v-%v-%v", order.GlobalId, order.ShopId, order.Sequence)
	request.OutTradeNo = fmt.Sprintf("%v-%v", order.GlobalId, order.Timestamp.Unix())
	request.TotalFee = order.Amount
	request.SpbillCreateIp = ip
	request.NotifyUrl = notifyUrl
	request.TradeType = "APP"
	request.Attach = order.Sequence //附加字段存储订单号

	requestJson, _ := json.Marshal(request)
	request.Sign = GeneratePaySign(string(requestJson))

	xmlValue, _ := xml.Marshal(request)
	resp, err := http.Post("https://api.mch.weixin.qq.com/pay/unifiedorder", "text/xml", bytes.NewBuffer(xmlValue))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"request": request,
			"err":     err,
		}).Error("Unifiedorder http failed.")
		return nil, err
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"request": request,
			"err":     err,
		}).Error("Unifiedorder ReadAll failed.")
		return nil, err
	}

	var unifiedorderResponse ResponseUnifiedorder
	if err := xml.Unmarshal(respBody, &unifiedorderResponse); err != nil {
		logrus.WithFields(logrus.Fields{
			"request":  request,
			"err":      err,
			"respBody": string(respBody),
		}).Error("Unifiedorder Unmarshal failed.")
		return nil, err
	}

	logrus.WithFields(logrus.Fields{
		"unifiedorderResponse": unifiedorderResponse,
	}).Info("Unifiedorder ResponseUnifiedorder")

	if unifiedorderResponse.ReturnCode == "FAIL" { //此字段是通信标识，非交易标识，交易是否成功需要查看result_code来判断
		logrus.WithFields(logrus.Fields{
			"request":              request,
			"unifiedorderResponse": unifiedorderResponse,
		}).Error("Unifiedorder ReturnCode FAIL.")

		return nil, errors.NewCodeString(errors.Failed, "unifiedorder failed return_code: %v, return_msg: %v", unifiedorderResponse.ReturnCode, unifiedorderResponse.ReturnMsg)
	}

	if unifiedorderResponse.ResultCode == "FAIL" { //业务结果
		logrus.WithFields(logrus.Fields{
			"request":              request,
			"unifiedorderResponse": unifiedorderResponse,
		}).Error("Unifiedorder ResultCode FAIL.")

		return nil, errors.NewCodeString(errors.Failed, "unifiedorder failed err_code: %v, err_code_des: %v", unifiedorderResponse.ErrCode, unifiedorderResponse.ErrCodeDes)
	}

	payParam := PullPay{}
	payParam.Appid = appId
	payParam.Partnerid = mchId
	payParam.Prepayid = unifiedorderResponse.PrepayId
	payParam.Package = "Sign=WXPay"
	payParam.Noncestr = utils.GenerateNonceStr()
	payParam.Timestamp = strconv.FormatInt(time.Now().Unix(), 10)

	payParamJson, _ := json.Marshal(payParam)
	payParam.Sign = GeneratePaySign(string(payParamJson))

	logrus.WithFields(logrus.Fields{
		"payParam": payParam,
	}).Info("Unifiedorder payParam")

	//order.SNSOrderId = request.OutTradeNo //将第三方订单号存入服务器订单
	/*wxOrder := &model.WechatPayOrder{
		GlobalId:       order.GlobalId,
		AppId:          appId,
		MchId:          mchId,
		Body:           request.Body,
		OutTradeNo:     request.OutTradeNo,
		TotalFee:       request.TotalFee,
		SpbillCreateIp: request.SpbillCreateIp,
		Sequence:       request.Attach,
		Prepayid:       payParam.Prepayid,
	}

	err = storage.Insert(storage.PayDatabase(), wxOrder)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":     err,
			"wxOrder": wxOrder,
		}).Error("Unifiedorder save order failed.")

		return nil, err
	}*/

	return &payParam, nil
}

type RequestOrderQuery struct {
	Appid      string `xml:"appid" json:"appid"`
	MchId      string `xml:"mch_id" json:"mch_id"`
	OutTradeNo string `xml:"out_trade_no" json:"out_trade_no"`
	NonceStr   string `xml:"nonce_str" json:"nonce_str"`
	Sign       string `xml:"sign" json:"sign,omitempty"`
}

type ResponseOrderQuery struct {
	ReturnCode     string `xml:"return_code" json:"return_code"`
	ReturnMsg      string `xml:"return_msg" json:"return_msg,omitempty"`
	Appid          string `xml:"appid" json:"appid"`
	MchId          string `xml:"mch_id" json:"mch_id"`
	NonceStr       string `xml:"nonce_str" json:"nonce_str"`
	Sign           string `xml:"sign" json:"-"`
	ResultCode     string `xml:"result_code" json:"result_code"`
	ErrCode        string `xml:"err_code" json:"err_code,omitempty"`
	ErrCodeDes     string `xml:"err_code_des" json:"err_code_des,omitempty"`
	DeviceInfo     string `xml:"device_info" json:"device_info,omitempty"`
	Openid         string `xml:"openid" json:"openid"`
	IsSubscribe    string `xml:"is_subscribe" json:"is_subscribe"`
	TradeType      string `xml:"trade_type" json:"trade_type"`
	TradeState     string `xml:"trade_state" json:"trade_state"`
	BankType       string `xml:"bank_type" json:"bank_type"`
	TotalFee       int    `xml:"total_fee" json:"total_fee,string"`
	FeeType        string `xml:"fee_type" json:"fee_type,omitempty"`
	CashFee        int    `xml:"cash_fee" json:"cash_fee,string"`
	CashFeeType    string `xml:"cash_fee_type" json:"cash_fee_type,omitempty"`
	TransactionId  string `xml:"transaction_id" json:"transaction_id"`
	OutTradeNo     string `xml:"out_trade_no" json:"out_trade_no"`
	Attach         string `xml:"attach" json:"attach,omitempty"`
	TimeEnd        string `xml:"time_end" json:"time_end"`
	TradeStateDesc string `xml:"trade_state_desc" json:"trade_state_desc"`
}

func WxPayOrderQuery(order *model.Order) (*ResponseOrderQuery, error) {
	mchId := settings.GetString("lovepay", "wechatPay.mch_id")
	appId := settings.GetString("lovepay", "wechatPay.appid")

	request := RequestOrderQuery{}
	request.Appid = appId
	request.MchId = mchId
	request.OutTradeNo = fmt.Sprintf("%v-%v", order.GlobalId, order.Timestamp.Unix())
	request.NonceStr = utils.GenerateNonceStr()

	requestJson, _ := json.Marshal(request)
	request.Sign = GeneratePaySign(string(requestJson))

	xmlValue, _ := xml.Marshal(request)
	resp, err := http.Post("https://api.mch.weixin.qq.com/pay/orderquery", "text/xml", bytes.NewBuffer(xmlValue))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"request": request,
			"err":     err,
		}).Error("WxPayOrderQuery http failed.")
		return nil, err
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"request": request,
			"err":     err,
		}).Error("WxPayOrderQuery ReadAll failed.")
		return nil, err
	}

	var orderQueryResponse ResponseOrderQuery
	if err := xml.Unmarshal(respBody, &orderQueryResponse); err != nil {
		logrus.WithFields(logrus.Fields{
			"err":      err,
			"respBody": string(respBody),
		}).Error("WxPayOrderQuery Unmarshal failed.")
		return nil, err
	}

	logrus.WithFields(logrus.Fields{
		"orderQueryResponse": orderQueryResponse,
	}).Info("WxPayOrderQuery ResponseOrderQuery")

	if orderQueryResponse.ReturnCode == "FAIL" { //此字段是通信标识，非交易标识，交易是否成功需要查看result_code来判断
		logrus.WithFields(logrus.Fields{
			"orderQueryResponse": orderQueryResponse,
		}).Error("WxPayOrderQuery ReturnCode FAIL.")

		return nil, errors.NewCodeString(errors.Failed, "orderquery failed return_code: %v, return_msg: %v", orderQueryResponse.ReturnCode, orderQueryResponse.ReturnMsg)
	}

	if orderQueryResponse.ResultCode == "FAIL" { //业务结果
		logrus.WithFields(logrus.Fields{
			"orderQueryResponse": orderQueryResponse,
		}).Error("WxPayOrderQuery ResultCode FAIL.")

		return nil, errors.NewCodeString(errors.Failed, "orderquery failed err_code: %v, err_code_des: %v", orderQueryResponse.ErrCode, orderQueryResponse.ErrCodeDes)
	}

	return &orderQueryResponse, nil
}
