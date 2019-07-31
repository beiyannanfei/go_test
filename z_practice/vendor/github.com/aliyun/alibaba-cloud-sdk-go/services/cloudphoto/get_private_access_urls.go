package cloudphoto

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

// GetPrivateAccessUrls invokes the cloudphoto.GetPrivateAccessUrls API synchronously
// api document: https://help.aliyun.com/api/cloudphoto/getprivateaccessurls.html
func (client *Client) GetPrivateAccessUrls(request *GetPrivateAccessUrlsRequest) (response *GetPrivateAccessUrlsResponse, err error) {
	response = CreateGetPrivateAccessUrlsResponse()
	err = client.DoAction(request, response)
	return
}

// GetPrivateAccessUrlsWithChan invokes the cloudphoto.GetPrivateAccessUrls API asynchronously
// api document: https://help.aliyun.com/api/cloudphoto/getprivateaccessurls.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) GetPrivateAccessUrlsWithChan(request *GetPrivateAccessUrlsRequest) (<-chan *GetPrivateAccessUrlsResponse, <-chan error) {
	responseChan := make(chan *GetPrivateAccessUrlsResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.GetPrivateAccessUrls(request)
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

// GetPrivateAccessUrlsWithCallback invokes the cloudphoto.GetPrivateAccessUrls API asynchronously
// api document: https://help.aliyun.com/api/cloudphoto/getprivateaccessurls.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) GetPrivateAccessUrlsWithCallback(request *GetPrivateAccessUrlsRequest, callback func(response *GetPrivateAccessUrlsResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *GetPrivateAccessUrlsResponse
		var err error
		defer close(result)
		response, err = client.GetPrivateAccessUrls(request)
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

// GetPrivateAccessUrlsRequest is the request struct for api GetPrivateAccessUrls
type GetPrivateAccessUrlsRequest struct {
	*requests.RpcRequest
	LibraryId string    `position:"Query" name:"LibraryId"`
	PhotoId   *[]string `position:"Query" name:"PhotoId"  type:"Repeated"`
	StoreName string    `position:"Query" name:"StoreName"`
	ZoomType  string    `position:"Query" name:"ZoomType"`
}

// GetPrivateAccessUrlsResponse is the response struct for api GetPrivateAccessUrls
type GetPrivateAccessUrlsResponse struct {
	*responses.BaseResponse
	Code      string   `json:"Code" xml:"Code"`
	Message   string   `json:"Message" xml:"Message"`
	RequestId string   `json:"RequestId" xml:"RequestId"`
	Action    string   `json:"Action" xml:"Action"`
	Results   []Result `json:"Results" xml:"Results"`
}

// CreateGetPrivateAccessUrlsRequest creates a request to invoke GetPrivateAccessUrls API
func CreateGetPrivateAccessUrlsRequest() (request *GetPrivateAccessUrlsRequest) {
	request = &GetPrivateAccessUrlsRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("CloudPhoto", "2017-07-11", "GetPrivateAccessUrls", "cloudphoto", "openAPI")
	return
}

// CreateGetPrivateAccessUrlsResponse creates a response to parse from GetPrivateAccessUrls response
func CreateGetPrivateAccessUrlsResponse() (response *GetPrivateAccessUrlsResponse) {
	response = &GetPrivateAccessUrlsResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
