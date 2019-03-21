package vivo

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

//var params = map[string]interface{}{
//"cpId":          productId,
//"appId":         appId,
//"cpOrderNumber": sequence,
//"notifyUrl":     callbackUrl,
//"orderAmount":   amount,
//"orderTitle":    "",
//"orderDesc":     "",
//"extInfo":       "",
//}
func BuyGoods(params map[string]interface{}, appKey string) (*VivoBuyGoodsResponse, error) {

	request := "https://pay.vivo.com.cn/vcoin/trade"

	params["version"] = "1.0.0"
	params["orderTime"] = time.Now().Format("20060102150405")

	params["signature"] = sign(params, appKey)
	params["signMethod"] = "MD5"

	request = request + "?" + buildPostBody(params)
	resp, err := http.Post(request, "application/json", nil)
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"error":  err,
			"params": params,
		}).Error("vivo buygoods http post failed.")

		return nil, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"error":  err,
			"params": params,
		}).Error("vivo buygoods readbody failed.")

		return nil, err
	}
	defer resp.Body.Close()

	buyGoodsResponse := &VivoBuyGoodsResponse{}
	err = json.Unmarshal(respBody, buyGoodsResponse)
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"error":    err,
			"params":   params,
			"respBody": string(respBody),
		}).Error("vivo buygoods json failed.")

		return nil, err
	}

	if buyGoodsResponse.Ret != "200" {

		logrus.WithFields(logrus.Fields{
			"params":   params,
			"respBody": string(respBody),
		}).Error("vivo buygoods ret failed.")

		return nil, errors.New("vivo buygoods ret:" + buyGoodsResponse.Ret)
	}

	var paramsRe map[string]interface{}
	err = json.Unmarshal(respBody, &paramsRe)
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"error":    err,
			"params":   params,
			"respBody": string(respBody),
		}).Error("vivo buygoods return params failed.")

		return nil, err
	}

	if !CheckSign(paramsRe, appKey) {

		logrus.WithFields(logrus.Fields{
			"params":   params,
			"respBody": string(respBody),
		}).Error("vivo buygoods check sign failed.")

		return nil, errors.New("vivo buygoods check sign failed.")
	}

	return buyGoodsResponse, nil
}

func buildPostBody(params map[string]interface{}) string {

	var varr []string
	for k, v := range params {

		varr = append(varr, k+"="+url.QueryEscape(fmt.Sprintf("%v", v)))
	}

	return strings.Join(varr, "&")
}

//{\"respCode\":\"400\",\"respMsg\":\"参数[version]错误\"}
type VivoBuyGoodsResponse struct {
	// 响应码
	//
	//成功返回：200，非200时，respMsg会提示错误信息
	Ret string `json:"respCode"`
	//响应信息
	//
	//对应响应码的响应信息
	Message string `json:"respMsg"`
	//签名方法
	//
	//对关键信息进行签名的算法名称：MD5
	SignMethod string `json:"signMethod"`
	//签名信息
	//
	//对关键信息签名后得到的字符串，用于商户验签，签名规则请参考签名计算说明
	Signature string `json:"signature"`
	//vivoSDK需要的参数
	//
	//vivoSDK使用
	AccessKey string `json:"accessKey"`
	//交易流水号
	//
	//vivo订单号
	OrderNumber string `json:"orderNumber"`
	//交易金额
	//
	//单位：分，币种：人民币，必须是整数
	OrderAmount string `json:"orderAmount"`
}
