package quick_sdk

import (
	"github.com/xykong/loveauth/settings"
	"fmt"
	"net/http"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"github.com/xykong/loveauth/errors"
)

func VerifyLoginQuick(token string, uid string) (bool, error) {
	product_code := settings.GetString("tencent", "quickSdk.ProductCode")
	url := fmt.Sprintf("http://checkuser.sdk.quicksdk.net/v2/checkUserInfo?token=%s&uid=%s&product_code=%s", token, uid, product_code)

	logrus.WithFields(logrus.Fields{
		"url": url,
	}).Info("VerifyLoginQuick")

	resp, err := http.Get(url)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("VerifyLoginQuick http failed.")
		return false, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("VerifyLoginQuick ReadAll body failed.")
		return false, err
	}

	if "1" != string(body) {
		logrus.WithFields(logrus.Fields{
			"body": string(body),
		}).Error("VerifyLoginQuick response failed.")
		return false, errors.NewCodeString(errors.Failed, "quick response: %v failed.", string(body))
	}

	return true, nil
}
