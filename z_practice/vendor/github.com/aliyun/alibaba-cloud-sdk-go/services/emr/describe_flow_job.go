package emr

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

// DescribeFlowJob invokes the emr.DescribeFlowJob API synchronously
// api document: https://help.aliyun.com/api/emr/describeflowjob.html
func (client *Client) DescribeFlowJob(request *DescribeFlowJobRequest) (response *DescribeFlowJobResponse, err error) {
	response = CreateDescribeFlowJobResponse()
	err = client.DoAction(request, response)
	return
}

// DescribeFlowJobWithChan invokes the emr.DescribeFlowJob API asynchronously
// api document: https://help.aliyun.com/api/emr/describeflowjob.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeFlowJobWithChan(request *DescribeFlowJobRequest) (<-chan *DescribeFlowJobResponse, <-chan error) {
	responseChan := make(chan *DescribeFlowJobResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.DescribeFlowJob(request)
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

// DescribeFlowJobWithCallback invokes the emr.DescribeFlowJob API asynchronously
// api document: https://help.aliyun.com/api/emr/describeflowjob.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) DescribeFlowJobWithCallback(request *DescribeFlowJobRequest, callback func(response *DescribeFlowJobResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *DescribeFlowJobResponse
		var err error
		defer close(result)
		response, err = client.DescribeFlowJob(request)
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

// DescribeFlowJobRequest is the request struct for api DescribeFlowJob
type DescribeFlowJobRequest struct {
	*requests.RpcRequest
	Id        string `position:"Query" name:"Id"`
	ProjectId string `position:"Query" name:"ProjectId"`
}

// DescribeFlowJobResponse is the response struct for api DescribeFlowJob
type DescribeFlowJobResponse struct {
	*responses.BaseResponse
	RequestId       string                        `json:"RequestId" xml:"RequestId"`
	Id              string                        `json:"Id" xml:"Id"`
	GmtCreate       int                           `json:"GmtCreate" xml:"GmtCreate"`
	GmtModified     int                           `json:"GmtModified" xml:"GmtModified"`
	Name            string                        `json:"Name" xml:"Name"`
	Type            string                        `json:"Type" xml:"Type"`
	Description     string                        `json:"Description" xml:"Description"`
	FailAct         string                        `json:"FailAct" xml:"FailAct"`
	MaxRetry        int                           `json:"MaxRetry" xml:"MaxRetry"`
	RetryInterval   int                           `json:"RetryInterval" xml:"RetryInterval"`
	Params          string                        `json:"Params" xml:"Params"`
	ParamConf       string                        `json:"ParamConf" xml:"ParamConf"`
	CustomVariables string                        `json:"CustomVariables" xml:"CustomVariables"`
	EnvConf         string                        `json:"EnvConf" xml:"EnvConf"`
	RunConf         string                        `json:"RunConf" xml:"RunConf"`
	MonitorConf     string                        `json:"MonitorConf" xml:"MonitorConf"`
	CategoryId      string                        `json:"CategoryId" xml:"CategoryId"`
	Mode            string                        `json:"mode" xml:"mode"`
	LastInstanceId  string                        `json:"LastInstanceId" xml:"LastInstanceId"`
	Adhoc           string                        `json:"Adhoc" xml:"Adhoc"`
	AlertConf       string                        `json:"AlertConf" xml:"AlertConf"`
	ResourceList    ResourceListInDescribeFlowJob `json:"ResourceList" xml:"ResourceList"`
}

// CreateDescribeFlowJobRequest creates a request to invoke DescribeFlowJob API
func CreateDescribeFlowJobRequest() (request *DescribeFlowJobRequest) {
	request = &DescribeFlowJobRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Emr", "2016-04-08", "DescribeFlowJob", "emr", "openAPI")
	return
}

// CreateDescribeFlowJobResponse creates a response to parse from DescribeFlowJob response
func CreateDescribeFlowJobResponse() (response *DescribeFlowJobResponse) {
	response = &DescribeFlowJobResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
