package ysdk

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"github.com/xykong/loveauth/errors"
)

//文档地址: https://open.weixin.qq.com/cgi-bin/showdocument?action=dir_list&t=resource/res_list&verify=1&id=open1419317853&lang=zh_CN
//接口说明
//检验授权凭证（access_token）是否有效
//请求说明
//http请求方式: GET
//https://api.weixin.qq.com/sns/auth?access_token=ACCESS_TOKEN&openid=OPENID

type VerifyAccessTokenRsp struct {
	Body struct {
		ErrCode int64  `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
}

func VerifyAccessToken(access_token string, openid string) error {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/auth?access_token=%s&openid=%s", access_token, openid)
	resp, err := http.Get(url)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"url": url,
			"err": err,
		}).Error("VerifyAccessToken http failed.")

		return err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	data := VerifyAccessTokenRsp{}
	if err := json.Unmarshal(body, &data.Body); err != nil {
		logrus.WithFields(logrus.Fields{
			"err":  err,
			"body": string(body),
			"url":  url,
		}).Error("VerifyAccessToken Unmarshal failed.")
		return err
	}

	if data.Body.ErrCode != 0 {
		logrus.WithFields(logrus.Fields{
			"err":  err,
			"body": string(body),
			"url":  url,
			"data": data,
		}).Error("VerifyAccessToken failed.")
		return errors.New(fmt.Sprintf("errcode: %v, errmsg: %v", data.Body.ErrCode, data.Body.ErrMsg))
	}

	return nil
}
