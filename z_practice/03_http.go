package main

import (
	"encoding/json"
	"net/http"
	"bytes"
	"fmt"
	"io/ioutil"
	"time"
)

func main() {
	go goPost()
	go goGet()
	time.Sleep(5 * time.Second)
}

type GetData struct {
	General     int `json:"general"`
	Accurate    int `json:"accurate"`
	Idcard      int `json:"idcard"`
	Bankcard    int `json:"bankcard"`
	Drivecard   int `json:"drivecard"`
	Vehiclecard int `json:"vehiclecard"`
	License     int `json:"license"`
	Business    int `json:"business"`
	Receipt     int `json:"receipt"`
	Enhance     int `json:"enhance"`
}

type GetResult struct {
	Status string  `json:"status"`
	Code   int     `json:"code"`
	Data   GetData `json:"data"`
}

func goGet() {
	resultBuf, _ := HttpGet("https://www.iocr.vip/ai/users/times/default")

	var result GetResult
	json.Unmarshal(resultBuf, &result)
	fmt.Printf("result: %#v\n", result)
}

/*
返回结果格式类似
HTTP/1.1 200 OK
   {
     "status":"ok",    //成功标志
     "code":200,
     "data": {
       general: 2,       //文字识别(含位置信息)
       accurate: 2,      //高精度文字识别(含位置信息)
       idcard: 20,       //识别身份证(正反面)
       bankcard: 20,     //银行卡识别
       drivecard: 10,    //驾驶证识别
       vehiclecard: 10,  //行驶证识别
       license: 10,      //车牌识别
       business: 10,     //营业执照识别
       receipt: 20,      //通用票据识别
       enhance: 2,       //通用文字识别（含生僻字版）- 无位置信息
     }
   }
*/
func HttpGet(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("http.Get err: %v\n", err.Error())
		return nil, err
	}

	fmt.Printf("http.Get finish response: %v\n", response)

	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("ioutil.ReadAll err: %v\n", err.Error())
		return nil, err
	}

	fmt.Printf("ioutil.ReadAll finish responseBody: %v\n", responseBody)

	return responseBody, nil
}

type PostResult struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
	Data   string `json:"data"`
}

func goPost() {
	postBody := map[string]interface{}{
		"content": "测试测试测试",
		"email":   "a@b.c",
	}
	resultBuf, _ := HttpPost("https://www.iocr.vip/ai/users/feedback", "application/json", postBody)

	var result PostResult
	json.Unmarshal(resultBuf, &result)
	fmt.Printf("result: %#v\n", result)
}

/*
返回数据格式类似
HTTP/1.1 200 OK
   {
     "status":"ok",    //成功标志
     "code":200,
     "data": "您反馈的问题已知悉，请等待回复..."
   }
*/
func HttpPost(url string, contentType string, body map[string]interface{}) ([]byte, error) {
	bodyJson, _ := json.Marshal(body)
	response, err := http.Post(url, contentType, bytes.NewBuffer(bodyJson))
	if err != nil {
		fmt.Printf("http.Post err: %v, bodyJson: %v\n", err.Error(), bodyJson)
		return nil, err
	}

	fmt.Printf("http.Post finish response: %v\n", response)

	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("ioutil.ReadAll err: %v\n", err.Error())
		return nil, err
	}

	fmt.Printf("ioutil.ReadAll finish responseBody: %v\n", responseBody)

	return responseBody, nil
}
