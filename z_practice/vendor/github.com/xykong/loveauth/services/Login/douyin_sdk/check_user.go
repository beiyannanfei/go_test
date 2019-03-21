package douyin_sdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

func CheckUser(clientKey, clientSecret, openId string) error {

	value := url.Values{}
	value.Add("client_key", clientKey)
	value.Add("client_secret", clientSecret)
	value.Add("open_id", openId)

	resp, err := http.Post("https://i.snssdk.com/game_sdk/check_user/", "application/x-www-form-urlencoded", bytes.NewBuffer([]byte(value.Encode())))
	if err != nil {

		return err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {

		return err
	}

	response := make(map[string]interface{})

	err = json.Unmarshal(respBody, &response)
	if err != nil {

		return nil
	}

	if response["message"] != "success" {

		return errors.New("douyin checkuser err " + string(respBody))
	}

	return nil
}
