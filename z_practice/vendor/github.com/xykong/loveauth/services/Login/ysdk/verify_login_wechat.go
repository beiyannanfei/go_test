package ysdk

import (
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/errors"
	"encoding/json"
)

func VerifyLoginWechat(openId string, token string, userIp string) (*DoVerifyLoginRsp, error) {
	module := "/auth/wx_check_token"
	platform := YSDK_WECHAT

	body := map[string]interface{}{
		"appid":   AppId(platform),
		"openid":  openId,
		"openkey": token,
		"userip":  userIp,
	}

	resp, err := GetRequestRaw(platform, module, body)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":  err,
			"module": module,
			"body":   body,
		}).Error("ysdk VerifyLogin failed.")
		return nil, err
	}

	var data DoVerifyLoginRsp
	if err := json.Unmarshal(resp, &data.Body); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":  err,
			"module": module,
			"body":   body,
		}).Error("ysdk VerifyLogin Unmarshal failed.")
		return nil, err
	}

	if data.Body.Ret != 0 {
		logrus.WithFields(logrus.Fields{
			"data": data,
		}).Error("ysdk VerifyLogin Ret failed.")
		return nil, errors.NewCode(errors.Failed)
	}

	return &data, nil
}
