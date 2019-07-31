package afs

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

// CreateConfiguration invokes the afs.CreateConfiguration API synchronously
// api document: https://help.aliyun.com/api/afs/createconfiguration.html
func (client *Client) CreateConfiguration(request *CreateConfigurationRequest) (response *CreateConfigurationResponse, err error) {
	response = CreateCreateConfigurationResponse()
	err = client.DoAction(request, response)
	return
}

// CreateConfigurationWithChan invokes the afs.CreateConfiguration API asynchronously
// api document: https://help.aliyun.com/api/afs/createconfiguration.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) CreateConfigurationWithChan(request *CreateConfigurationRequest) (<-chan *CreateConfigurationResponse, <-chan error) {
	responseChan := make(chan *CreateConfigurationResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.CreateConfiguration(request)
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

// CreateConfigurationWithCallback invokes the afs.CreateConfiguration API asynchronously
// api document: https://help.aliyun.com/api/afs/createconfiguration.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) CreateConfigurationWithCallback(request *CreateConfigurationRequest, callback func(response *CreateConfigurationResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *CreateConfigurationResponse
		var err error
		defer close(result)
		response, err = client.CreateConfiguration(request)
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

// CreateConfigurationRequest is the request struct for api CreateConfiguration
type CreateConfigurationRequest struct {
	*requests.RpcRequest
	SourceIp            string `position:"Query" name:"SourceIp"`
	ConfigurationName   string `position:"Query" name:"ConfigurationName"`
	MaxPV               string `position:"Query" name:"MaxPV"`
	ConfigurationMethod string `position:"Query" name:"ConfigurationMethod"`
	ApplyType           string `position:"Query" name:"ApplyType"`
	Scene               string `position:"Query" name:"Scene"`
}

// CreateConfigurationResponse is the response struct for api CreateConfiguration
type CreateConfigurationResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
	BizCode   string `json:"BizCode" xml:"BizCode"`
	RefExtId  string `json:"RefExtId" xml:"RefExtId"`
}

// CreateCreateConfigurationRequest creates a request to invoke CreateConfiguration API
func CreateCreateConfigurationRequest() (request *CreateConfigurationRequest) {
	request = &CreateConfigurationRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("afs", "2018-01-12", "CreateConfiguration", "afs", "openAPI")
	return
}

// CreateCreateConfigurationResponse creates a response to parse from CreateConfiguration response
func CreateCreateConfigurationResponse() (response *CreateConfigurationResponse) {
	response = &CreateConfigurationResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
