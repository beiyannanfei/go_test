package main

import (
	"io/ioutil"
	"fmt"
	"encoding/json"
)

//解析json文件

type GlobalConfigSetting struct {
	VersionRange   string        `json:"-"`
	ShowServerList bool          `json:"showServerList"`
	Tags           []string      `json:"tags"`
	GConfigList    []interface{} `json:"gConfigList"`
}

func main() {
	filePath := "./53_jsonfilesrc.json"
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("read file err:", err)
		return
	}

	var Settings map[string]GlobalConfigSetting

	err = json.Unmarshal(fileContent, &Settings)
	if err != nil {
		fmt.Println("Unmarshal json to map err:", err)
		return
	}

	fmt.Printf("Settings: %v\n", Settings)

	/*var SetInterface interface{}	//直接解析成interface会无法遍历
	err = json.Unmarshal(fileContent, &SetInterface)
	if err != nil {
		fmt.Println("Unmarshal json to interface err:", err)
		return
	}
	fmt.Printf("SetInterface: %#v\n", SetInterface)*/

	fmt.Println("----------------------------------------")

	ConfigSettings := make([]GlobalConfigSetting, 0)
	DefaultSetting := make(map[string]GlobalConfigSetting)

	for vr, v := range Settings {
		fmt.Printf("vr: %v, v: %v\n", vr, v)
		if vr == "default" {
			DefaultSetting[vr] = v
			continue
		}

		v.VersionRange = vr
		ConfigSettings = append(ConfigSettings, v)
	}

	fmt.Println("----------------------------------------")
	fmt.Printf("ConfigSettings: %v\n", ConfigSettings)
	fmt.Printf("DefaultSetting: %v\n", DefaultSetting)
}
