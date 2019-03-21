package msdk

import (
	"time"
	"github.com/xykong/loveauth/settings"
	"fmt"
	"crypto/md5"
	log "github.com/sirupsen/logrus"
	"net/http"
	"bytes"
	"io/ioutil"
	"encoding/json"
	"github.com/xykong/loveauth/utils"
	"github.com/xykong/loveauth/errors"
)

type Platform string

const (
	QQ     = Platform("QQ")
	Wechat = Platform("Wechat")
	Guest  = Platform("Guest")
)

func Timestamp() uint {
	return uint(time.Now().Unix())
}

//noinspection GoNameStartsWithPackageName
func MsdkKey() string {
	return settings.GetString("tencent", "msdk.MsdkKey")
}

func Host() string {
	return settings.GetString("tencent", "msdk.host")
}

func AppId(platform Platform) string {
	return settings.GetString("tencent", fmt.Sprintf("msdk.%v.AppId", platform))
}

func AppKey(platform Platform) string {
	return settings.GetString("tencent", fmt.Sprintf("msdk.%v.AppKey", platform))
}

/*
platform	string	平台标识(一般情况下：qq对应值为desktop_m_qq，wx对应值为desktop_m_wx，游客对应值为desktop_m_guest)
 */
func PlatformId(platform Platform) string {

	switch platform {
	case QQ:
		return "desktop_m_qq"
	case Wechat:
		return "desktop_m_wx"
	case Guest:
		return "desktop_m_guest"
	}

	return ""
}

/*

3）sig生成规则
当encode=1 时：sig = md5 ( appkey + timestamp )
当encode=2 时：sig = md5 ( msdkkey+ timestamp )
+ 表示两个字符串连接，非字符串"+";
 */
func Sig(encode int, platform Platform, timestamp uint) string {

	switch encode {
	case 1:
		return fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s%d", AppKey(platform), timestamp))))
	case 2:
		return fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s%d", MsdkKey(), timestamp))))
	}

	return ""
}

func UrlParams(openId string, platform Platform) string {

	var timestamp = Timestamp()

	var request = fmt.Sprintf("timestamp=%d&appid=%s&sig=%s&openid=%s&encode=2",
		timestamp, AppId(platform), Sig(2, platform, timestamp), openId)

	return request
}

func UrlRequest(openId string, platform Platform, module string) string {

	return fmt.Sprintf("%s/%s?%s", Host(), module, UrlParams(openId, platform))
}

func PostRequestRaw(openId string, platform Platform, module string, body map[string]interface{}) ([]byte, error) {

	newrelicApp := utils.GetNewRelic()
	if newrelicApp != nil {

		Transaction := fmt.Sprintf("%s/%s", Host(), module)

		//log.WithFields(log.Fields{
		//	"Transaction": Transaction,
		//}).Info("New Relic StartTransaction.")

		txn := newrelicApp.StartTransaction(Transaction, nil, nil)
		defer txn.End()
	}

	start := time.Now()

	request := UrlRequest(openId, platform, module)

	jsonValue, _ := json.Marshal(body)

	resp, err := http.Post(request, "application/json", bytes.NewBuffer(jsonValue))

	if err != nil {

		log.WithFields(log.Fields{
			"elapsed": time.Since(start),
			"request": request,
			"body":    string(jsonValue),
		}).Error("Msdk post request.")

		return nil, err
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	log.WithFields(log.Fields{
		"elapsed":  time.Since(start),
		"request":  request,
		"body":     string(jsonValue),
		"response": string(respBody),
	}).Info("Msdk post request.")

	return respBody, nil
}

func PostRequest(openId string, platform Platform, module string, body map[string]interface{}) (errors.Code, string) {

	resp, err := PostRequestRaw(openId, platform, module, body)

	if err != nil {
		log.Error(err)
		return errors.Failed, err.Error()
	}
	var data map[string]interface{}

	if err := json.Unmarshal(resp, &data); err != nil {
		return errors.Failed, err.Error()
	}

	if data["ret"].(float64) != 0 {
		return errors.Failed, string(resp)
	}

	return errors.Ok, string(resp)
}
