package webo_sdk

import (
	"github.com/xykong/loveauth/settings"
	"fmt"
	"net/http"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"encoding/json"
	"github.com/xykong/loveauth/errors"
)

type GetUserInfoRsp struct {
	Body struct {
		//用户昵称
		ScreenName string `json:"screen_name"`
		//用户头像地址（中图），50×50像素
		ProfileImageUrl string `json:"profile_image_url"`

		//错误信息
		Error     string `json:"error"`
		ErrorCode int64  `json:"error_code"`
		Request   string `json:"request"`
	}
}

func WeiBoGetUserInfo(openId string) (*GetUserInfoRsp, error) {
	host := settings.GetString("tencent", "weibo.Host")
	AppKey := settings.GetString("tencent", "weibo.AppKey")
	url := fmt.Sprintf("%s%s?source=%s&uid=%s", host, "/2/users/show.json", AppKey, openId)

	resp, err := http.Get(url)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"url": url,
			"err": err,
		}).Error("WeiBoGetUserInfo http failed.")
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	data := GetUserInfoRsp{}
	if err := json.Unmarshal(body, &data.Body); err != nil {
		logrus.WithFields(logrus.Fields{
			"err":  err,
			"body": string(body),
			"url":  url,
		}).Error("WeiBoGetUserInfo Unmarshal failed.")
		return nil, err
	}

	if data.Body.ErrorCode != 0 {
		logrus.WithFields(logrus.Fields{
			"err":  err,
			"body": string(body),
			"url":  url,
			"data": data,
		}).Error("WeiBoGetUserInfo get user info failed.")
		return nil, errors.New(data.Body.Error)
	}

	return &data, nil
}
