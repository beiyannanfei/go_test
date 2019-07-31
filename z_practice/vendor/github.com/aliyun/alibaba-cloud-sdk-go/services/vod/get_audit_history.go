package vod

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

// GetAuditHistory invokes the vod.GetAuditHistory API synchronously
// api document: https://help.aliyun.com/api/vod/getaudithistory.html
func (client *Client) GetAuditHistory(request *GetAuditHistoryRequest) (response *GetAuditHistoryResponse, err error) {
	response = CreateGetAuditHistoryResponse()
	err = client.DoAction(request, response)
	return
}

// GetAuditHistoryWithChan invokes the vod.GetAuditHistory API asynchronously
// api document: https://help.aliyun.com/api/vod/getaudithistory.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) GetAuditHistoryWithChan(request *GetAuditHistoryRequest) (<-chan *GetAuditHistoryResponse, <-chan error) {
	responseChan := make(chan *GetAuditHistoryResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.GetAuditHistory(request)
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

// GetAuditHistoryWithCallback invokes the vod.GetAuditHistory API asynchronously
// api document: https://help.aliyun.com/api/vod/getaudithistory.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) GetAuditHistoryWithCallback(request *GetAuditHistoryRequest, callback func(response *GetAuditHistoryResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *GetAuditHistoryResponse
		var err error
		defer close(result)
		response, err = client.GetAuditHistory(request)
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

// GetAuditHistoryRequest is the request struct for api GetAuditHistory
type GetAuditHistoryRequest struct {
	*requests.RpcRequest
	PageNo   requests.Integer `position:"Query" name:"PageNo"`
	PageSize requests.Integer `position:"Query" name:"PageSize"`
	VideoId  string           `position:"Query" name:"VideoId"`
	SortBy   string           `position:"Query" name:"SortBy"`
}

// GetAuditHistoryResponse is the response struct for api GetAuditHistory
type GetAuditHistoryResponse struct {
	*responses.BaseResponse
	RequestId string    `json:"RequestId" xml:"RequestId"`
	Status    string    `json:"Status" xml:"Status"`
	Total     int64     `json:"Total" xml:"Total"`
	Histories []History `json:"Histories" xml:"Histories"`
}

// CreateGetAuditHistoryRequest creates a request to invoke GetAuditHistory API
func CreateGetAuditHistoryRequest() (request *GetAuditHistoryRequest) {
	request = &GetAuditHistoryRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("vod", "2017-03-21", "GetAuditHistory", "vod", "openAPI")
	return
}

// CreateGetAuditHistoryResponse creates a response to parse from GetAuditHistory response
func CreateGetAuditHistoryResponse() (response *GetAuditHistoryResponse) {
	response = &GetAuditHistoryResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
