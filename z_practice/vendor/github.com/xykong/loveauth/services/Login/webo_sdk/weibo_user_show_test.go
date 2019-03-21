package webo_sdk

import (
	"testing"
	"fmt"
	"encoding/json"
)

func TestWeiboGetUserInfo(t *testing.T) {
	test := struct {
		name   string
		openId string
	}{
		"TestWeiboGetUserInfo", "140437656",	//1404376560 可以正常返回
	}

	t.Run(test.name, func(t *testing.T) {
		resp, err := WeiBoGetUserInfo(test.openId)
		j, _ := json.Marshal(resp)
		fmt.Println(string(j), err)
	})
}
