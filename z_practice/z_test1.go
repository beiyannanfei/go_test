package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"encoding/base64"
	"net/url"
)

func main() {
	mac := hmac.New(sha1.New, []byte("228bf094169a40a3bd188ba37ebe8723&"))
	mac.Write([]byte("GET&%2Fv3%2Fuser%2Fget_info&appid%3D123456%26format%3Djson%26openid%3D11111111111111111%26openkey%3D2222222222222222%26pf%3Dqzone%26userip%3D112.90.139.30"))
	jmStr := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	fmt.Println(jmStr)

	fmt.Println("============================")

	uri := "/v3/user/get_info"
	encodeurl := url.QueryEscape(uri)
	fmt.Println(encodeurl)

}
