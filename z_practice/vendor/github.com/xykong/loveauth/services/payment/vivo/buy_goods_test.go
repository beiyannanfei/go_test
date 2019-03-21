package vivo

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestBuyGoods(t *testing.T) {

	test := struct {
		name   string
		params map[string]interface{}
		appKey string
	}{
		"BuyGoods",
		map[string]interface{}{
			"cpId":          "61f4119a42ee794e3c55",
			"appId":         "79676d36a3b15ac7e9401ead043259e9",
			"cpOrderNumber": fmt.Sprintf("%v-%v", 123131231231312312, time.Now().Unix()),
			"notifyUrl":     "http://127.0.0.1/callback",
			"orderAmount":   "1000",
			"orderTitle":    "ssss",
			"orderDesc":     "aaaaaaa",
			"extInfo":       "dddddddddd",
		},
		"86abf10cc7c90347e1fc7299b870e891",
	}

	t.Run(test.name, func(t *testing.T) {

		//{"respCode":"200","respMsg":"success","signMethod":"MD5","signature":"8dbd5166b43a34ad7e7f7a06fe29c71d","accessKey":"f1052ffb59fda92eaa254b8ca4074024","orderNumber":"2018122517563215300016326854","orderAmount":"1000"} <nil>
		resp, err := BuyGoods(test.params, test.appKey)
		j, _ := json.Marshal(resp)

		fmt.Println(string(j), err)
	})
}
