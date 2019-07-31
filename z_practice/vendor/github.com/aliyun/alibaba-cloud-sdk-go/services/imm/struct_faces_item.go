package imm

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

// FacesItem is a nested struct in imm response
type FacesItem struct {
	FaceConfidence       float64        `json:"FaceConfidence" xml:"FaceConfidence"`
	EmotionConfidence    float64        `json:"EmotionConfidence" xml:"EmotionConfidence"`
	ImageUri             string         `json:"ImageUri" xml:"ImageUri"`
	Similarity           float64        `json:"Similarity" xml:"Similarity"`
	FaceQuality          float64        `json:"FaceQuality" xml:"FaceQuality"`
	Attractive           float64        `json:"Attractive" xml:"Attractive"`
	AttractiveConfidence float64        `json:"AttractiveConfidence" xml:"AttractiveConfidence"`
	Age                  int            `json:"Age" xml:"Age"`
	AgeConfidence        float64        `json:"AgeConfidence" xml:"AgeConfidence"`
	Gender               string         `json:"Gender" xml:"Gender"`
	GenderConfidence     float64        `json:"GenderConfidence" xml:"GenderConfidence"`
	Emotion              string         `json:"Emotion" xml:"Emotion"`
	FaceId               string         `json:"FaceId" xml:"FaceId"`
	GroupId              string         `json:"GroupId" xml:"GroupId"`
	FaceAttributes       FaceAttributes `json:"FaceAttributes" xml:"FaceAttributes"`
	EmotionDetails       EmotionDetails `json:"EmotionDetails" xml:"EmotionDetails"`
}
