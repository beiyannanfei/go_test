package ysdk

import (
	"github.com/xykong/loveauth/settings"
	"fmt"
	"net/http"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"encoding/json"
	"github.com/xykong/loveauth/errors"
)

type RefreshAccessTokenRsp struct {
	Body struct {
		ErrCode      int64  `json:"errcode"`
		ErrMsg       string `json:"errmsg"`
		AccessToken  string `json:"access_token"`
		ExpiresIn    int64  `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		OpenId       string `json:"openid"`
		Scope        string `json:"scope"`
	}
}

//文档地址: https://open.weixin.qq.com/cgi-bin/showdocument?action=dir_list&t=resource/res_list&verify=1&id=open1419317853&lang=zh_CN
//刷新或续期access_token使用
//接口说明
//access_token是调用授权关系接口的调用凭证，由于access_token有效期（目前为2个小时）较短，当access_token超时后，可以使用refresh_token进行刷新，access_token刷新结果有两种：
//1.若access_token已超时，那么进行refresh_token会获取一个新的access_token，新的超时时间；
//2.若access_token未超时，那么进行refresh_token不会改变access_token，但超时时间会刷新，相当于续期access_token。
//refresh_token拥有较长的有效期（30天）且无法续期，当refresh_token失效的后，需要用户重新授权后才可以继续获取用户头像昵称。
//请求方法
//使用/sns/oauth2/access_token接口获取到的refresh_token进行以下接口调用：
//调用频率限制: 刷新access_token	10万/分钟
func RefreshAccessToken(refresh_token string) (*RefreshAccessTokenRsp, error) {
	AppId := settings.GetInt64("tencent", "ysdk.YSDK_Wechat.AppId")
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/refresh_token?appid=%s&grant_type=%s&refresh_token=%s", AppId, "refresh_token", refresh_token)
	resp, err := http.Get(url)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"url": url,
			"err": err,
		}).Error("RefreshAccessToken http failed.")

		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	data := RefreshAccessTokenRsp{}
	if err := json.Unmarshal(body, &data.Body); err != nil {
		logrus.WithFields(logrus.Fields{
			"err":  err,
			"body": string(body),
			"url":  url,
		}).Error("RefreshAccessToken Unmarshal failed.")
		return nil, err
	}

	if data.Body.ErrCode != 0 {
		logrus.WithFields(logrus.Fields{
			"err":  err,
			"body": string(body),
			"url":  url,
			"data": data,
		}).Error("RefreshAccessToken failed.")
		return nil, errors.New(fmt.Sprintf("errcode: %v, errmsg: %v", data.Body.ErrCode, data.Body.ErrMsg))
	}

	return &data, nil
}
