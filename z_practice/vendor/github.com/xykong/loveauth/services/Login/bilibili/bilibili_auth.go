package bilibili

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/errors"
	"github.com/xykong/loveauth/settings"
	"io/ioutil"
	"net/http"
)

type DoAuthTokenRsp struct {
	Body struct {
		OpenId     int    `json:"open_id"` // b站游戏用户的用户id
		UName      string `json:"uname"`
		Timestamp  int    `json:"timestamp"` // 对应request的时间戳，秒
		Code       int    `json:"code"`      // 状态码
		ErrMessage string `json:"message"`   // 错误信息
	}
}

func GoAuthBiliBili(domain string, accessKey string, uid int32) (*DoAuthTokenRsp, error) {
	p := GetCommonReqParams(uid)
	p.Add("access_key", accessKey)
	sign := GetSign(p)
	p.Add("sign", sign)
	//client := &http.Client{}
	//url := fmt.Sprintf("%s/api/server/session.verify?%s", domain, p.Encode())
	//req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte("")))
	//req.Header.Set("User-Agent", settings.GetString("lovepay", "bilibili.header.user-agent"))
	//req.Header.Set("Content-Type", settings.GetString("lovepay", "bilibili.header.content-type"))
	//resp, err := client.Do(req)
	url := fmt.Sprintf("%s/api/server/session.verify", domain)
	resp, err := http.PostForm(url, p)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
			"url": url,
			"p":   p.Encode(),
		}).Error("GoAuthBiliBili client.Do failed.")

		return nil, err
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
			"url": url,
			"p":   p.Encode(),
		}).Error("GoAuthBiliBili ReadAll failed.")

		return nil, err
	}

	var data DoAuthTokenRsp
	if err := json.Unmarshal(respBody, &data.Body); err != nil {
		logrus.WithFields(logrus.Fields{
			"err":      err,
			"url":      url,
			"p":        p.Encode(),
			"respBody": string(respBody),
		}).Error("GoAuthBiliBili Unmarshal failed.")

		return nil, err
	}

	if 0 != data.Body.Code {
		logrus.WithFields(logrus.Fields{
			"url":  url,
			"p":    p.Encode(),
			"data": data,
		}).Error("GoAuthBiliBili auth code failed.")

		return nil, errors.NewCodeString(errors.Failed, "bilibili AuthToken response code: %v err, errmsg: %v", data.Body.Code, data.Body.ErrMessage)
	}

	return &data, nil
}

func AuthToken(accessKey string, uid int32) (*DoAuthTokenRsp, error) {
	urls := settings.GetStringSlice("lovepay", "bilibili.url")
	for i := 0; i < len(urls); i += 1 {
		response, err := GoAuthBiliBili(urls[i], accessKey, uid)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"accessKey": accessKey,
				"uid":       uid,
				"url":       urls[i],
			}).Error("bilibili AuthToken failed.")

			if _, ok := err.(*errors.Type); ok { //认证失败错误
				return nil, err
			}

			//非认证失败，切线再次验证
			continue
		}

		return response, nil
	}

	return nil, errors.NewCodeString(errors.Failed, "bilibili AuthToken all domain failed.")
}
