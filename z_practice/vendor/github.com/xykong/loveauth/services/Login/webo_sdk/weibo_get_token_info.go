package webo_sdk

import (
	"github.com/xykong/loveauth/settings"
	"fmt"
	"bytes"
	"net/http"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"github.com/xykong/loveauth/errors"
	"strconv"
)

type GetTokenInfoRsp struct {
	Body struct {
		//授权用户的uid。
		Uid int64 `json:"uid"`
		//access_token所属的应用appkey。
		Appkey string `json:"appkey"`
		//用户授权的scope权限。
		Scope string `json:"scope"`
		//access_token的创建时间，从1970年到创建时间的秒数。
		CreateAt int64 `json:"create_at"`
		//access_token的剩余时间，单位是秒数。
		ExpireIn int64 `json:"expire_in"`

		//错误信息
		Error     string `json:"error"`
		ErrorCode int64  `json:"error_code"`
		Request   string `json:"request"`
	}
}

func WeiBoGetTokenInfo(accessToken string) (*GetTokenInfoRsp, error) {
	host := settings.GetString("tencent", "weibo.Host")
	url := fmt.Sprintf("%s%s?access_token=%s", host, "/oauth2/get_token_info", accessToken)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte("")))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
			"url": url,
		}).Error("WeiBoGetTokenInfo http post failed.")
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":  err,
			"url":  url,
			"resp": resp,
		}).Error("WeiBoGetTokenInfo ReadAll failed.")
		return nil, err
	}

	var data GetTokenInfoRsp
	if err := json.Unmarshal(respBody, &data.Body); err != nil {
		logrus.WithFields(logrus.Fields{
			"err":      err,
			"url":      url,
			"resp":     resp,
			"respBody": string(respBody),
		}).Error("WeiBoGetTokenInfo Unmarshal failed.")
		return nil, err
	}

	if data.Body.ErrorCode != 0 {
		logrus.WithFields(logrus.Fields{
			"err":      err,
			"url":      url,
			"resp":     resp,
			"respBody": string(respBody),
			"data":     data,
		}).Error("WeiBoGetTokenInfo oauth2 get_token_info failed.")
		return nil, errors.New(data.Body.Error)
	}

	AppKey := settings.GetString("tencent", "weibo.AppKey")
	if AppKey != data.Body.Appkey {
		return nil, errors.New("AppKey not match.")
	}

	if "" == strconv.FormatInt(data.Body.Uid, 10) {
		return nil, errors.New("AppKey not exists.")
	}

	return &data, nil
}
