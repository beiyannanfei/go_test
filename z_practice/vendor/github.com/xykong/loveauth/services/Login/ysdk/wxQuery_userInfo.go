package ysdk

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/xykong/loveauth/errors"
)

//文档地址: https://open.weixin.qq.com/cgi-bin/showdocument?action=dir_list&t=resource/res_list&verify=1&id=open1419317853&lang=zh_CN
//获取用户个人信息（UnionID机制）
//接口说明
//此接口用于获取用户个人信息。开发者可通过OpenID来获取用户基本信息。特别需要注意的是，如果开发者拥有多个移动应用、网站应用和公众帐号，可通过获取用户基本信息中的unionid来区分用户的唯一性，因为只要是同一个微信开放平台帐号下的移动应用、网站应用和公众帐号，用户的unionid是唯一的。换句话说，同一用户，对同一个微信开放平台下的不同应用，unionid是相同的。请注意，在用户修改微信头像后，旧的微信头像URL将会失效，因此开发者应该自己在获取用户信息后，将头像图片保存下来，避免微信头像URL失效后的异常情况。
//请求说明
//http请求方式: GET
//https://api.weixin.qq.com/sns/userinfo?access_token=ACCESS_TOKEN&openid=OPENID

type QueryUserInfoRsp struct {
	Body struct {
		ErrCode int64  `json:"errcode"`
		ErrMsg  string `json:"errmsg"`

		OpenId     string   `json:"openid"`
		NickName   string   `json:"nickname"`
		Sex        int      `json:"sex"`
		Province   string   `json:"province"`
		City       string   `json:"city"`
		Country    string   `json:"country"`
		Headimgurl string   `json:"headimgurl"`
		Privilege  []string `json:"privilege"`
		UnionId    string   `json:"unionid"`
	}
}

func QueryUserInfo(access_token string, openid string) (*QueryUserInfoRsp, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s", access_token, openid)
	resp, err := http.Get(url)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"url": url,
			"err": err,
		}).Error("QueryUserInfo http failed.")

		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	data := QueryUserInfoRsp{}
	if err := json.Unmarshal(body, &data.Body); err != nil {
		logrus.WithFields(logrus.Fields{
			"err":  err,
			"body": string(body),
			"url":  url,
		}).Error("QueryUserInfo Unmarshal failed.")
		return nil, err
	}

	if data.Body.ErrCode != 0 {
		logrus.WithFields(logrus.Fields{
			"err":  err,
			"body": string(body),
			"url":  url,
			"data": data,
		}).Error("QueryUserInfo failed.")
		return nil, errors.New(fmt.Sprintf("errcode: %v, errmsg: %v", data.Body.ErrCode, data.Body.ErrMsg))
	}

	return &data, nil
}
