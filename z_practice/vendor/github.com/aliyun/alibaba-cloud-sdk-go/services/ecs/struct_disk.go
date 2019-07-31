package ecs

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

// Disk is a nested struct in ecs response
type Disk struct {
	DiskId                        string                        `json:"DiskId" xml:"DiskId"`
	RegionId                      string                        `json:"RegionId" xml:"RegionId"`
	ZoneId                        string                        `json:"ZoneId" xml:"ZoneId"`
	DiskName                      string                        `json:"DiskName" xml:"DiskName"`
	Description                   string                        `json:"Description" xml:"Description"`
	Type                          string                        `json:"Type" xml:"Type"`
	Category                      string                        `json:"Category" xml:"Category"`
	Size                          int                           `json:"Size" xml:"Size"`
	ImageId                       string                        `json:"ImageId" xml:"ImageId"`
	SourceSnapshotId              string                        `json:"SourceSnapshotId" xml:"SourceSnapshotId"`
	AutoSnapshotPolicyId          string                        `json:"AutoSnapshotPolicyId" xml:"AutoSnapshotPolicyId"`
	ProductCode                   string                        `json:"ProductCode" xml:"ProductCode"`
	Portable                      bool                          `json:"Portable" xml:"Portable"`
	Status                        string                        `json:"Status" xml:"Status"`
	InstanceId                    string                        `json:"InstanceId" xml:"InstanceId"`
	Device                        string                        `json:"Device" xml:"Device"`
	DeleteWithInstance            bool                          `json:"DeleteWithInstance" xml:"DeleteWithInstance"`
	DeleteAutoSnapshot            bool                          `json:"DeleteAutoSnapshot" xml:"DeleteAutoSnapshot"`
	EnableAutoSnapshot            bool                          `json:"EnableAutoSnapshot" xml:"EnableAutoSnapshot"`
	EnableAutomatedSnapshotPolicy bool                          `json:"EnableAutomatedSnapshotPolicy" xml:"EnableAutomatedSnapshotPolicy"`
	CreationTime                  string                        `json:"CreationTime" xml:"CreationTime"`
	AttachedTime                  string                        `json:"AttachedTime" xml:"AttachedTime"`
	DetachedTime                  string                        `json:"DetachedTime" xml:"DetachedTime"`
	DiskChargeType                string                        `json:"DiskChargeType" xml:"DiskChargeType"`
	ExpiredTime                   string                        `json:"ExpiredTime" xml:"ExpiredTime"`
	ResourceGroupId               string                        `json:"ResourceGroupId" xml:"ResourceGroupId"`
	Encrypted                     bool                          `json:"Encrypted" xml:"Encrypted"`
	MountInstanceNum              int                           `json:"MountInstanceNum" xml:"MountInstanceNum"`
	IOPS                          int                           `json:"IOPS" xml:"IOPS"`
	IOPSRead                      int                           `json:"IOPSRead" xml:"IOPSRead"`
	IOPSWrite                     int                           `json:"IOPSWrite" xml:"IOPSWrite"`
	KMSKeyId                      string                        `json:"KMSKeyId" xml:"KMSKeyId"`
	BdfId                         string                        `json:"BdfId" xml:"BdfId"`
	OperationLocks                OperationLocksInDescribeDisks `json:"OperationLocks" xml:"OperationLocks"`
	MountInstances                MountInstances                `json:"MountInstances" xml:"MountInstances"`
	Tags                          TagsInDescribeDisks           `json:"Tags" xml:"Tags"`
}
