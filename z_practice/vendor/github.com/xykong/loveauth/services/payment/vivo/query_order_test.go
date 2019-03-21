package vivo

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestQueryOrder(t *testing.T) {

	test := struct {
		name   string
		params map[string]interface{}
		appKey string
	}{
		"QueryOrder",
		map[string]interface{}{
			"cpId":          "61f4119a42ee794e3c55",
			"appId":         "79676d36a3b15ac7e9401ead043259e9",
			"cpOrderNumber": "8430286606071-1547534399",
			"orderNumber":   "2019011514395944400012247300",
			"orderAmount":   "100",
		},
		"86abf10cc7c90347e1fc7299b870e891",
	}

	t.Run(test.name, func(t *testing.T) {

		//{"respCode":"200","respMsg":"查询成功","signMethod":"MD5","signature":"f69b328a7631484d43ddfb07b0986288","tradeType":"01","tradeStatus":"0001","cpId":"61f4119a42ee794e3c55","appId":"79676d36a3b15ac7e9401ead043259e9","uid":"","cpOrderNumber":"725584057435359099632877691","orderNumber":"2018122516231426500010584826","orderAmount":"1000","extInfo":"dddddddddd","payTime":"20181225162314"} <nil>
		resp, err := QueryOrder(test.params, test.appKey)
		j, _ := json.Marshal(resp)

		fmt.Println(string(j), err)
	})
}
