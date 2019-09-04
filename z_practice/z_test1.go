package main

import "fmt"

func main() {
	params := map[string]interface{}{
		"cpId":          "cpId",
		"appId":         "appId",
		"cpOrderNumber": "cpOrderNumber",
		"notifyUrl":     "notifyUrl",
		"orderAmount":   "orderAmount",
		"orderTitle":    "orderTitle",
		"orderDesc":     "orderDesc",
		"extInfo":       "extInfo",
	}

	var s string
	s = fmt.Sprintf("%v", params["orderTitle"])
	fmt.Println(s)
}
