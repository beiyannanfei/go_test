/**
 * Generated by the loveworld utility tools.  DO NOT EDIT!
 * Source: qqscore_batch
 * User: xy.kong@gmail.com
 * DateTime: 2018-03-27 12:17:25.096109 +0800 CST m=+1.455784948
 */
package msdk

import (
	"github.com/sirupsen/logrus"
	"encoding/json"
	"github.com/xykong/loveauth/errors"
)

//
// 应答: [手Q]成就上报
// swagger:response DoQQScoreBatchRsp
// noinspection ALL
type DoQQScoreBatchRsp struct {
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

func QQScoreBatch(openId string, token string) (*DoQQScoreBatchRsp, error) {

	//noinspection ALL
	module := "/profile/qqscore_batch"

	platform := QQ

	//noinspection SpellCheckingInspection
	body := map[string]interface{}{
		"appid":       AppId(platform),
		"openid":      openId,
		"accessToken": token,
		"param":       nil,
	}

	//noinspection SpellCheckingInspection
	resp, err := PostRequestRaw(openId, platform, module, body)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":  err,
			"module": module,
			"body":   body,
		}).Error("QQScoreBatch failed.")
		return nil, err
	}

	var data DoQQScoreBatchRsp
	if err := json.Unmarshal(resp, &data.Body); err != nil {
		logrus.WithFields(logrus.Fields{
			"error":  err,
			"module": module,
			"body":   body,
		}).Error("QQScoreBatch failed.")
		return nil, err
	}

	if data.Body.Ret != 0 {
		return nil, errors.NewCode(errors.Failed)
	}

	return &data, nil
}