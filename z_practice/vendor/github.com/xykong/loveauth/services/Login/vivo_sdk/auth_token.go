package vivo_sdk

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/errors"
	"io/ioutil"
	"net/http"
)

//
// 应答: vivo authtoken
// swagger:response DoAuthTokenRsp
// noinspection ALL
type DoAuthTokenRsp struct {
	// in: body
	Body struct {
		//
		// 返回码 0：正确，其它：失败
		//
		Ret int `json:"retcode"`
		//
		// ret非0，则表示“错误码，错误提示”，详细注释参见错误码描述
		//
		Data struct {
			//
			// 状态
			//
			Success bool `json:"success"`
			//
			// 平台唯一id
			//
			Openid string `json:"openid"`
		} `json:"data"`
	}
}

func AuthToken(authToken string) (*DoAuthTokenRsp, error) {

	request := "https://usrsys.vivo.com.cn/sdk/user/auth.do?authtoken=" + authToken

	resp, err := http.Post(request, "application/json", nil)
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"error":  err,
			"module": "vivo/authtoken",
		}).Error("AuthToken failed.")

		return nil, err
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("vivo authtoken readbody.")

		return nil, err
	}

	logrus.WithFields(logrus.Fields{
		"request":  request,
		"response": string(respBody),
	}).Info("vivo authtoken post request.")

	var data DoAuthTokenRsp
	if err := json.Unmarshal(respBody, &data.Body); err != nil {

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("vivo authtoken failed.")

		return nil, err
	}

	if data.Body.Ret != 0 {

		return nil, errors.NewCode(errors.Failed)
	}

	return &data, nil
}
