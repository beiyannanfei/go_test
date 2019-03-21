package ysdk

import (
	"github.com/xykong/loveauth/settings"
	"fmt"
	"github.com/xykong/loveauth/utils"
	"time"
	"crypto/md5"
	"net/http"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

type Platform string

const (
	YSDK_QQ = Platform("YSDK_QQ")
	YSDK_WECHAT = Platform("YSDK_Wechat")
)

func Timestamp() uint {
	return uint(time.Now().Unix())
}

func AppId(platform Platform) string {
	return settings.GetString("tencent", fmt.Sprintf("ysdk.%v.AppId", platform))
}

func Host() string {
	return settings.GetString("tencent", "ysdk.host")
}

func AppKey(platform Platform) string {
	return settings.GetString("tencent", fmt.Sprintf("ysdk.%v.AppKey", platform))
}

func Sig(platform Platform, timestamp uint) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s%d", AppKey(platform), timestamp))))
}

func UrlParams(platform Platform, body map[string]interface{}) string {
	var timestamp = Timestamp()
	var request = fmt.Sprintf("timestamp=%d&appid=%s&sig=%s&openid=%s&openkey=%s&userip=%s",
		timestamp, AppId(platform), Sig(platform, timestamp), body["openid"], body["openkey"], body["userip"])

	return request
}

func UrlRequest(platform Platform, module string, body map[string]interface{}) string {
	return fmt.Sprintf("%s%s?%s", Host(), module, UrlParams(platform, body))
}

func GetRequestRaw(platform Platform, module string, body map[string]interface{}) ([]byte, error) {
	newrelicApp := utils.GetNewRelic()
	if newrelicApp != nil {
		Transaction := fmt.Sprintf("%s/%s", Host(), module)
		txn := newrelicApp.StartTransaction(Transaction, nil, nil)
		defer txn.End()
	}

	start := time.Now()
	request := UrlRequest(platform, module, body)
	resp, err := http.Get(request)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"elapsed": time.Since(start),
			"request": request,
			"body":    body,
			"err":     err,
		}).Error("YSDK get request failed.")
		return nil, err
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"elapsed": time.Since(start),
			"request": request,
			"body":    body,
			"err":     err,
		}).Error("YSDK ReadAll failed.")
		return nil, err
	}

	logrus.WithFields(logrus.Fields{
		"elapsed": time.Since(start),
		"request": request,
		"body":    body,
	}).Info("YSDK get request.")

	return respBody, nil
}
