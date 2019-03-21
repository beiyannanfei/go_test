/**
 * Created by wyq on 2019/3/7.
 */
package main

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"time"
)

type Sms_config_struct struct {
	accessKeyId     string
	accessKeySecret string
	mobile          string
	signName        string
	templateId      string
	templateParam   string
}

var confList []Sms_config_struct

func main() {
	confList = append(confList, Sms_config_struct{
		"accessKeyId",
		"accessKeySecret",
		"mobile",
		"signName",
		"templateId",
		"templateParam",
	})

	for _, item := range confList {
		goSendMsg(item.accessKeyId, item.accessKeySecret, item.mobile, item.signName, item.templateId, item.templateParam)
		time.Sleep(time.Second)
	}
}

func goSendMsg(accessKeyId, accessKeySecret, mobile, signName, templateId, templateParam string) {
	client, err := sdk.NewClientWithAccessKey("cn-hangzhou", accessKeyId, accessKeySecret)
	if err != nil {
		fmt.Println("SendAliAuth NewClientWithAccessKey failed.")
		return
	}

	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"
	request.QueryParams["RegionId"] = "cn-hangzhou"
	request.QueryParams["PhoneNumbers"] = mobile
	request.QueryParams["SignName"] = signName
	request.QueryParams["TemplateCode"] = templateId
	request.QueryParams["TemplateParam"] = templateParam

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		fmt.Printf("SendAliAuth ProcessCommonRequest failed err: %s. templateId: %s\n", err.Error(), templateId)
		return
	}

	fmt.Printf("SendAliAuth finish. response: %s, templateId: %s\n", response.GetHttpContentString(), templateId)
}
