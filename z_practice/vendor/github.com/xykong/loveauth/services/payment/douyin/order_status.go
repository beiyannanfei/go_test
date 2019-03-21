package douyin

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
)

func OrderStatus(clientKey, clientSecret, outTradeNo string) error {

	value := url.Values{}
	value.Add("client_key", clientKey)
	value.Add("client_secret", clientSecret)
	value.Add("out_trade_no", outTradeNo)

	resp, err := http.Post("https://i.snssdk.com/game_sdk/order_status/", "application/x-www-form-urlencoded", bytes.NewBuffer([]byte(value.Encode())))
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"err":     err,
			"request": value.Encode(),
		}).Error("douyin orderstatus post err")

		return err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"err":     err,
			"request": value.Encode(),
		}).Error("douyin orderstatus ioread err")

		return err
	}

	var response DouyinOrderStatusResponse
	err = json.Unmarshal(respBody, &response)
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"err":      err,
			"request":  value.Encode(),
			"response": string(respBody),
		}).Error("douyin orderstatus json unmarshal err")

		return err
	}

	if response.Message != "success" || response.Data.OrderStatus != 0 {

		logrus.WithFields(logrus.Fields{
			"request":  value.Encode(),
			"response": string(respBody),
		}).Error("douyin orderstatus message err")

		return errors.New("douyin orderstatus err " + string(respBody))
	}

	return nil
}

type DouyinOrderStatusResponse struct {
	Message string `json:"message"`
	Data    struct {
		ErrorCode   int `json:"error_code"`
		OrderStatus int `json:"order_status"`
	} `json:"data"`
}
