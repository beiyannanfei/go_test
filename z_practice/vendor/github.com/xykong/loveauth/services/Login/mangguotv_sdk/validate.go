package mangguotv_sdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

func ValidateUser(ticket, uuid string) error {

	value := url.Values{}
	value.Add("ticket", ticket)
	value.Add("uuid", uuid)

	resp, err := http.Post("http://cmop.mgtv.com/s/v2/account/cpValidateUser", "application/x-www-form-urlencoded", bytes.NewBuffer([]byte(value.Encode())))
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

	if response["result"] != "200" {

		return errors.New("mangguotv checkuser err " + string(respBody))
	}

	return nil
}
