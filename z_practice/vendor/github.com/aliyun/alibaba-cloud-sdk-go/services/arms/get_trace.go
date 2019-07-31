package arms

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

// GetTrace invokes the arms.GetTrace API synchronously
// api document: https://help.aliyun.com/api/arms/gettrace.html
func (client *Client) GetTrace(request *GetTraceRequest) (response *GetTraceResponse, err error) {
	response = CreateGetTraceResponse()
	err = client.DoAction(request, response)
	return
}

// GetTraceWithChan invokes the arms.GetTrace API asynchronously
// api document: https://help.aliyun.com/api/arms/gettrace.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) GetTraceWithChan(request *GetTraceRequest) (<-chan *GetTraceResponse, <-chan error) {
	responseChan := make(chan *GetTraceResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.GetTrace(request)
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

// GetTraceWithCallback invokes the arms.GetTrace API asynchronously
// api document: https://help.aliyun.com/api/arms/gettrace.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) GetTraceWithCallback(request *GetTraceRequest, callback func(response *GetTraceResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *GetTraceResponse
		var err error
		defer close(result)
		response, err = client.GetTrace(request)
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

// GetTraceRequest is the request struct for api GetTrace
type GetTraceRequest struct {
	*requests.RpcRequest
	TraceID string `position:"Query" name:"TraceID"`
	AppType string `position:"Query" name:"AppType"`
}

// GetTraceResponse is the response struct for api GetTrace
type GetTraceResponse struct {
	*responses.BaseResponse
	RequestId string         `json:"RequestId" xml:"RequestId"`
	Data      DataInGetTrace `json:"Data" xml:"Data"`
}

// CreateGetTraceRequest creates a request to invoke GetTrace API
func CreateGetTraceRequest() (request *GetTraceRequest) {
	request = &GetTraceRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("ARMS", "2019-02-19", "GetTrace", "", "")
	return
}

// CreateGetTraceResponse creates a response to parse from GetTrace response
func CreateGetTraceResponse() (response *GetTraceResponse) {
	response = &GetTraceResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
