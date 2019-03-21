package ysdk

import (
	"testing"
	"fmt"
)

func TestYsdkVerifyLogin(t *testing.T) {
	test := struct {
		name   string
		openId string
		token  string
		userIp string
	}{
		"TestYsdkVerifyLoginQQ",
		"30D0F2482CE9B2C90F0064485A261A29",
		"123FD4E8A0E448EACBBEC2596821043E",
		"127.0.0.1",
	}

	t.Run(test.name, func(t *testing.T) {
		resp, err := VerifyLoginQQ(test.openId, test.token, test.userIp)
		fmt.Println(resp, err)
	})

	/*test_Wechat := struct {
		name   string
		openId string
		token  string
		userIp string
	}{
		"TestYsdkVerifyLoginWechat",
		"oxSlRwdnw5bUh92hm9jqxrzSf7tI",
		"16_HGS99CJGpnU6puT28WDgb4VLLgAQW9c5R5Msr4kUGsYI6oW91YBvhTEjG369pm5OXZCJkJF_Ti-XMHPyTOtfOCK5eduF_nQF8GEziySDbiA",
		"127.0.0.1",
	}

	t.Run(test_Wechat.name, func(t *testing.T) {
		resp, err := VerifyLoginWechat(test_Wechat.openId, test_Wechat.token, test_Wechat.userIp)
		fmt.Println(resp, err)
	})*/
}
