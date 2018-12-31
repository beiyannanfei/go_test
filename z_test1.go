package main

import (
	"encoding/json"
	"fmt"
)

type Ta struct {
	Aa string `json:"aa"`
	Bb int    `json:"bb"`
	Cc string `json:"cc,omitempty"`
	Dd int    `json:"dd,omitempty"`
}

func main() {
	t := Ta{}
	t.Aa = "aa"
	t.Bb = 20

	tStrBytes, _ := json.Marshal(t)

	var ti interface{}
	json.Unmarshal(tStrBytes, &ti)

	//fmt.Println(ti)
	mm, ok := ti.(map[string]interface{})
	fmt.Println(mm, ok)

	for k, v := range mm {
		fmt.Printf("%v=%v\n", k, v)
	}
}
