package live

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

// DescribeLiveStreamDelayConfig invokes the live.DescribeLiveStreamDelayConfig API synchronously
// api document: https://help.aliyun.com/api/live/describelivestreamdelayconfig.html
func (client *Client) DescribeLiveStreamDelayConfig(request *DescribeLiveStreamDelayConfigRequest) (response *DescribeLiveStreamDelayConfigResponse, err error) {
	response = CreateDescribeLiveStreamDelayConfigResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeLiveStreamDelayConfigWithChan invokes the live.DescribeLiveStreamDelayConfig API asynchronously
// api document: https://help.aliyun.com/api/live/describelivestreamdelayconfig.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeLiveStreamDelayConfigWithChan(request *DescribeLiveStreamDelayConfigRequest) (<-chan *DescribeLiveStreamDelayConfigResponse, <-chan error) {
	responseChan := make(chan *DescribeLiveStreamDelayConfigResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeLiveStreamDelayConfig(request)
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

// DescribeLiveStreamDelayConfigWithCallback invokes the live.DescribeLiveStreamDelayConfig API asynchronously
// api document: https://help.aliyun.com/api/live/describelivestreamdelayconfig.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeLiveStreamDelayConfigWithCallback(request *DescribeLiveStreamDelayConfigRequest, callback func(response *DescribeLiveStreamDelayConfigResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeLiveStreamDelayConfigResponse
		var err error
		defer close(result)
		response, err = client.DescribeLiveStreamDelayConfig(request)
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

// DescribeLiveStreamDelayConfigRequest is the request struct for api DescribeLiveStreamDelayConfig
type DescribeLiveStreamDelayConfigRequest struct {
	*requests.RpcRequest
	DomainName string           `position:"Query" name:"DomainName"`
	OwnerId    requests.Integer `position:"Query" name:"OwnerId"`
}

// DescribeLiveStreamDelayConfigResponse is the response struct for api DescribeLiveStreamDelayConfig
type DescribeLiveStreamDelayConfigResponse struct {
	*responses.BaseResponse
	RequestId                 string                    `json:"RequestId" xml:"RequestId"`
	LiveStreamHlsDelayConfig  LiveStreamHlsDelayConfig  `json:"LiveStreamHlsDelayConfig" xml:"LiveStreamHlsDelayConfig"`
	LiveStreamFlvDelayConfig  LiveStreamFlvDelayConfig  `json:"LiveStreamFlvDelayConfig" xml:"LiveStreamFlvDelayConfig"`
	LiveStreamRtmpDelayConfig LiveStreamRtmpDelayConfig `json:"LiveStreamRtmpDelayConfig" xml:"LiveStreamRtmpDelayConfig"`
}

// CreateDescribeLiveStreamDelayConfigRequest creates a request to invoke DescribeLiveStreamDelayConfig API
func CreateDescribeLiveStreamDelayConfigRequest() (request *DescribeLiveStreamDelayConfigRequest) {
	request = &DescribeLiveStreamDelayConfigRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("live", "2016-11-01", "DescribeLiveStreamDelayConfig", "live", "openAPI")
	return
}

// CreateDescribeLiveStreamDelayConfigResponse creates a response to parse from DescribeLiveStreamDelayConfig response
func CreateDescribeLiveStreamDelayConfigResponse() (response *DescribeLiveStreamDelayConfigResponse) {
	response = &DescribeLiveStreamDelayConfigResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
