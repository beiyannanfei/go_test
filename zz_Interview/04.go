package main

import (
	"net/http"
	"bytes"
	"fmt"
	"io/ioutil"
	"encoding/json"
)

type DoAuthTokenRsp1 struct {
	RtnCode int    `form:"rtnCode"  json:"rtnCode"  binding:"required"` // 0=成功, -1=失败  1=接口鉴权失败  3001=参数错误
	Ts      string `form:"ts"       json:"ts"       binding:"required"` // 时间戳，接口返回的当前时间戳
	RtnSign string `form:"rtnSign"  json:"rtnSign"  binding:"required"` // 返回参数签名值，请根据rtnCode和ts校验签名是否正确
	ErrMsg  string `form:"errMsg"  json:"errMsg"`
}

func main() {
	authUrl := ""
	resp, err := http.Post(authUrl, "application/x-www-form-urlencoded", bytes.NewBuffer([]byte("")))
	if err != nil {
		fmt.Println("Post err:", err)
		return
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ReadAll err: ", err)
		return
	}

	if http.StatusOK != resp.StatusCode {
		fmt.Println("StatusCode err code: ", resp.StatusCode)
		return
	}

	var data *DoAuthTokenRsp1
	if err := json.Unmarshal(respBody, &data); err != nil {
		fmt.Println("Unmarshal err: ", err)
		return
	}

	fmt.Printf("data: %#v\n, StatusCode: %v\n", data, resp.StatusCode)
}
