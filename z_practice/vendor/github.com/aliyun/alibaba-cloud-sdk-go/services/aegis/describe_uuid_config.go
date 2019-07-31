package aegis

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

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// DescribeUuidConfig invokes the aegis.DescribeUuidConfig API synchronously
// api document: https://help.aliyun.com/api/aegis/describeuuidconfig.html
func (client *Client) DescribeUuidConfig(request *DescribeUuidConfigRequest) (response *DescribeUuidConfigResponse, err error) {
	response = CreateDescribeUuidConfigResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeUuidConfigWithChan invokes the aegis.DescribeUuidConfig API asynchronously
// api document: https://help.aliyun.com/api/aegis/describeuuidconfig.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeUuidConfigWithChan(request *DescribeUuidConfigRequest) (<-chan *DescribeUuidConfigResponse, <-chan error) {
	responseChan := make(chan *DescribeUuidConfigResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeUuidConfig(request)
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

// DescribeUuidConfigWithCallback invokes the aegis.DescribeUuidConfig API asynchronously
// api document: https://help.aliyun.com/api/aegis/describeuuidconfig.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeUuidConfigWithCallback(request *DescribeUuidConfigRequest, callback func(response *DescribeUuidConfigResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeUuidConfigResponse
		var err error
		defer close(result)
		response, err = client.DescribeUuidConfig(request)
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

// DescribeUuidConfigRequest is the request struct for api DescribeUuidConfig
type DescribeUuidConfigRequest struct {
	*requests.RpcRequest
	SourceIp string `position:"Query" name:"SourceIp"`
	Uuid     string `position:"Query" name:"Uuid"`
}

// DescribeUuidConfigResponse is the response struct for api DescribeUuidConfig
type DescribeUuidConfigResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	Data      Data   `json:"Data" xml:"Data"`
}

// CreateDescribeUuidConfigRequest creates a request to invoke DescribeUuidConfig API
func CreateDescribeUuidConfigRequest() (request *DescribeUuidConfigRequest) {
	request = &DescribeUuidConfigRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("aegis", "2016-11-11", "DescribeUuidConfig", "vipaegis", "openAPI")
	return
}

// CreateDescribeUuidConfigResponse creates a response to parse from DescribeUuidConfig response
func CreateDescribeUuidConfigResponse() (response *DescribeUuidConfigResponse) {
	response = &DescribeUuidConfigResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
