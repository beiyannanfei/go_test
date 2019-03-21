package huawei

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/errors"
	"github.com/xykong/loveauth/settings"
	"github.com/xykong/loveauth/utils/encoding"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
)

type DoAuthTokenRsp struct {
	RtnCode int    `form:"rtnCode"  json:"rtnCode"  binding:"required"` // 0=成功, -1=失败  1=接口鉴权失败  3001=参数错误
	Ts      string `form:"ts"       json:"ts"       binding:"required"` // 时间戳，接口返回的当前时间戳
	RtnSign string `form:"rtnSign"  json:"rtnSign"  binding:"required"` // 返回参数签名值，请根据rtnCode和ts校验签名是否正确
	ErrMsg  string `form:"errMsg"   json:"errMsg"`                      // 报错信息
}

// ts 为毫秒数
func AuthToken(playerId, playerLevel, playerSSign, ts string) (*DoAuthTokenRsp, error) {
	var queryMap = make(map[string]string)
	queryMap["method"] = "external.hms.gs.checkPlayerSign"
	queryMap["appId"] = settings.GetString("tencent", "huawei.appId")
	queryMap["cpId"] = settings.GetString("tencent", "huawei.cpId")
	queryMap["ts"] = ts
	queryMap["playerId"] = playerId
	queryMap["playerLevel"] = playerLevel
	queryMap["playerSSign"] = playerSSign

	var keyList []string
	for key, _ := range queryMap {
		if "cpSign" == key || "rtnSign" == key {
			continue
		}
		keyList = append(keyList, key)
	}
	sort.Strings(keyList)

	p := url.Values{}
	for _, key := range keyList {
		p.Add(key, queryMap[key])
	}

	// 获取私钥
	privateKey, err := encoding.GetRSAPrivateKey("-----BEGIN PRIVATE KEY-----\n" + settings.GetString("tencent", "huawei.privateKey") + "\n-----END PRIVATE KEY-----\n")
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err":      err,
			"queryMap": queryMap,
		}).Error("huawei AuthToken GetRSAPrivateKey failed.")

		return nil, err
	}

	// 加签
	signature, err := encoding.SignSHA256WithRSA(p.Encode(), privateKey)
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"queryMap": queryMap,
			"err":      err,
		}).Error("huawei AuthToken SignSHA256WithRSA failed.")
		return nil, err
	}

	p.Add("cpSign", signature)

	// 发送请求
	authUrl := settings.GetString("tencent", "huawei.authUrl")
	resp, err := http.PostForm(authUrl, p)
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err": err,
			"p":   p.Encode(),
		}).Error("huawei AuthToken PostForm failed.")
		return nil, err
	}
	defer resp.Body.Close()

	// 解析返回信息
	respBody, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err":  err,
			"p":    p.Encode(),
			"resp": resp,
		}).Error("huawei AuthToken ReadAll failed")

		return nil, err
	}

	if http.StatusOK != resp.StatusCode {
		logrus.WithFields(logrus.Fields{
			"statusCode": resp.StatusCode,
			"p":          p.Encode(),
		}).Error("huawei AuthToken PostForm StatusCode failed.")

		return nil, errors.New("response status code = " + strconv.Itoa(resp.StatusCode))
	}

	var data *DoAuthTokenRsp
	if err := json.Unmarshal(respBody, &data); err != nil {
		logrus.WithFields(logrus.Fields{
			"err":      err,
			"respBody": string(respBody),
			"p":        p.Encode(),
		}).Error("huawei AuthToken Unmarshal failed")

		return nil, err
	}

	if 0 != data.RtnCode {
		logrus.WithFields(logrus.Fields{
			"p":    p.Encode(),
			"data": data,
		}).Error("huawei AuthToken RtnCode failed.")

		return nil, errors.New("rtnCode = " + strconv.Itoa(data.RtnCode) + ", reason = " + data.ErrMsg)
	}

	// 校验返回的签名
	respStr := "rtnCode=" + strconv.Itoa(data.RtnCode) + "&ts=" + data.Ts
	respSign, err := encoding.SignSHA256WithRSA(respStr, privateKey)
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err":     err,
			"data":    data,
			"respStr": respStr,
		}).Error("huawei AuthToken response SignSHA256WithRSA failed.")

		return nil, err
	}

	// todo sudo apt-get install pcscd libccid opensc opensc-pkcs11
	if respSign != data.RtnSign {
		logrus.WithFields(logrus.Fields{
			"responSign": respSign,
			"rtnSign":    data.RtnSign,
			"data":       data,
		}).Error("huawei AuthToken response check sign failed.")

		return nil, errors.New("response sign error.")
	}

	return data, nil
}
