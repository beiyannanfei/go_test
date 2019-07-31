//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package iot

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// QueryProduct invokes the iot.QueryProduct API synchronously
// api document: https://help.aliyun.com/api/iot/queryproduct.html
func (client *Client) QueryProduct(request *QueryProductRequest) (response *QueryProductResponse, err error) {
	response = CreateQueryProductResponse()
	err = client.DoAction(request, response)
	return
}

// QueryProductWithChan invokes the iot.QueryProduct API asynchronously
// api document: https://help.aliyun.com/api/iot/queryproduct.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) QueryProductWithChan(request *QueryProductRequest) (<-chan *QueryProductResponse, <-chan error) {
	responseChan := make(chan *QueryProductResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.QueryProduct(request)
		if err != nil {
			errChan <- err
		} else {
			responseChan <- response
		}
	})
	if err != nil {
		errChan <- err
		close(responseChan)
		close(errChan)
	}
	return responseChan, errChan
}

// QueryProductWithCallback invokes the iot.QueryProduct API asynchronously
// api document: https://help.aliyun.com/api/iot/queryproduct.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) QueryProductWithCallback(request *QueryProductRequest, callback func(response *QueryProductResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *QueryProductResponse
		var err error
		defer close(result)
		response, err = client.QueryProduct(request)
		callback(response, err)
		result <- 1
	})
	if err != nil {
		defer close(result)
		callback(nil, err)
		result <- 0
	}
	return result
}

// QueryProductRequest is the request struct for api QueryProduct
type QueryProductRequest struct {
	*requests.RpcRequest
	AccessKeyId   string `position:"Query" name:"AccessKeyId"`
	IotInstanceId string `position:"Query" name:"IotInstanceId"`
	ProductKey    string `position:"Query" name:"ProductKey"`
}

// QueryProductResponse is the response struct for api QueryProduct
type QueryProductResponse struct {
	*responses.BaseResponse
	RequestId    string            `json:"RequestId" xml:"RequestId"`
	Success      bool              `json:"Success" xml:"Success"`
	Code         string            `json:"Code" xml:"Code"`
	ErrorMessage string            `json:"ErrorMessage" xml:"ErrorMessage"`
	Data         QueryProductData0 `json:"Data" xml:"Data"`
}

type QueryProductData0 struct {
	GmtCreate           int64  `json:"GmtCreate" xml:"GmtCreate"`
	DataFormat          int    `json:"DataFormat" xml:"DataFormat"`
	Description         string `json:"Description" xml:"Description"`
	DeviceCount         int    `json:"DeviceCount" xml:"DeviceCount"`
	NodeType            int    `json:"NodeType" xml:"NodeType"`
	ProductKey          string `json:"ProductKey" xml:"ProductKey"`
	ProductName         string `json:"ProductName" xml:"ProductName"`
	ProductSecret       string `json:"ProductSecret" xml:"ProductSecret"`
	CategoryName        string `json:"CategoryName" xml:"CategoryName"`
	CategoryKey         string `json:"CategoryKey" xml:"CategoryKey"`
	AliyunCommodityCode string `json:"AliyunCommodityCode" xml:"AliyunCommodityCode"`
	Id2                 bool   `json:"Id2" xml:"Id2"`
	ProtocolType        string `json:"ProtocolType" xml:"ProtocolType"`
	ProductStatus       string `json:"ProductStatus" xml:"ProductStatus"`
	Owner               bool   `json:"Owner" xml:"Owner"`
	NetType             int    `json:"NetType" xml:"NetType"`
}

// CreateQueryProductRequest creates a request to invoke QueryProduct API
func CreateQueryProductRequest() (request *QueryProductRequest) {
	request = &QueryProductRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Iot", "2018-01-20", "QueryProduct", "iot", "openAPI")
	return
}

// CreateQueryProductResponse creates a response to parse from QueryProduct response
func CreateQueryProductResponse() (response *QueryProductResponse) {
	response = &QueryProductResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
