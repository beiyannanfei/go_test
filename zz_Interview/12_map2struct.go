package main

import (
	"fmt"
	"encoding/json"
)

type Ss struct {
	Aa string `json:"aa"`
	Bb int    `json:"bb,string"`
	Cc string `json:"cc"`
}

func main() {
	var m = make(map[string]string)
	m["aa"] = "10"
	m["bb"] = "20"
	m["cc"] = "30"

	fmt.Println(m)
	mJson, _ := json.Marshal(m)
	fmt.Println(string(mJson))

	var s Ss
	err := json.Unmarshal(mJson, &s)
	fmt.Println(err)
	fmt.Println(s)
}
