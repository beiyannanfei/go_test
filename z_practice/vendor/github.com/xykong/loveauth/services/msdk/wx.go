/**
 * Generated by the loveworld utility tools.  DO NOT EDIT!
 * Source: wx
 * User: xy.kong@gmail.com
 * DateTime: 2018-03-27 12:17:25.098358 +0800 CST m=+1.458033251
 */
package msdk

import (
	"github.com/sirupsen/logrus"
	"encoding/json"
	"github.com/xykong/loveauth/errors"
)

//
// 应答: [微信]后台分享
// swagger:response DoWXRsp
// noinspection ALL
type DoWXRsp struct {
	// in: body
	Body struct {
		//
		// 返回码 0：正确，其它：失败
		//
		Ret int `json:"ret"`
		//
		// ret非0，则表示“错误码，错误提示”，详细注释参见错误码描述
		//
		Msg string `json:"msg"`
	}
}

func WX(openId string, token string) (*DoWXRsp, error) {

	//noinspection ALL
	module := "/share/wx"

	platform := QQ

	//noinspection SpellCheckingInspection
	body := map[string]interface{}{
		"appid":          AppId(platform),
		"openid":         openId,
		"fopenid":        "",
		"access_token":   "",
		"extinfo":        "",
		"title":          "",
		"description":    "",
		"media_tag_name": "",
		"thumb_media_id": "",
	}

	//noinspection SpellCheckingInspection
	resp, err := PostRequestRaw(openId, platform, module, body)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":  err,
			"module": module,
			"body":   body,
		}).Error("WX failed.")
		return nil, err
	}

	var data DoWXRsp
	if err := json.Unmarshal(resp, &data.Body); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":  err,
			"module": module,
			"body":   body,
		}).Error("WX failed.")
		return nil, err
	}

	if data.Body.Ret != 0 {
		return nil, errors.NewCode(errors.Failed)
	}

	return &data, nil
}