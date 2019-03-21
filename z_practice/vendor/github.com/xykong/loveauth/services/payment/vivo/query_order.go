package vivo

import (
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/storage/model"
	"io/ioutil"
	"net/http"
)

//var params = map[string]interface{}{
//"cpId":          productId,
//"appId":         appId,
//"cpOrderNumber": sequence,
//"orderAmount":   amount,
//"orderNumber":   orderNumber,
//}
func QueryOrder(params map[string]interface{}, appKey string) (*model.VivoQueryOrderResponse, error) {

	request := "https://pay.vivo.com.cn/vcoin/queryv2"

	params["version"] = "1.0.0"

	params["signature"] = sign(params, appKey)
	params["signMethod"] = "MD5"

	request = request + "?" + buildPostBody(params)

	resp, err := http.Post(request, "application/json", nil)
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"error":  err,
			"params": params,
		}).Error("vivo QueryOrder http post failed.")

		return nil, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"error":  err,
			"params": params,
		}).Error("vivo QueryOrder readbody failed.")

		return nil, err
	}
	defer resp.Body.Close()

	queryOrderResponse := &model.VivoQueryOrderResponse{}
	err = json.Unmarshal(respBody, queryOrderResponse)
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"error":    err,
			"params":   params,
			"respBody": string(respBody),
		}).Error("vivo QueryOrder json failed.")

		return nil, err
	}

	logrus.WithFields(logrus.Fields{
		"requestParams": params,
		"response":      string(respBody),
	}).Info("vivo queryorder request")

	if queryOrderResponse.Ret != "200" || queryOrderResponse.TradeStatus != "0000" {

		logrus.WithFields(logrus.Fields{
			"params":   params,
			"respBody": string(respBody),
		}).Error("vivo queryorder ret failed.")

		return nil, errors.New("vivo queryorder ret:" + queryOrderResponse.Ret + " tradeStatus:" + queryOrderResponse.TradeStatus)
	}

	var paramsRe map[string]interface{}
	err = json.Unmarshal(respBody, &paramsRe)
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"error":    err,
			"params":   params,
			"respBody": string(respBody),
		}).Error("vivo queryorder return params failed.")

		return nil, err
	}

	if !CheckSign(paramsRe, appKey) {

		logrus.WithFields(logrus.Fields{
			"params":   params,
			"respBody": string(respBody),
		}).Error("vivo queryorder check sign failed.")

		return nil, errors.New("vivo queryorder check sign failed.")
	}

	return queryOrderResponse, nil
}
