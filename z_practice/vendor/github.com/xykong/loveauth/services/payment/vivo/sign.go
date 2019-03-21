package vivo

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

func sign(params map[string]interface{}, appKey string) string {

	var keys []string
	for k := range params {

		keys = append(keys, k)
	}

	sort.Strings(keys)

	var varr []string
	for _, k := range keys {

		if k == "signMethod" || k == "signature" {

			continue
		}

		v := fmt.Sprintf("%v", params[k])
		if v != "" {

			varr = append(varr, k+"="+v)
		}
	}

	stringA := strings.Join(varr, "&")

	fmt.Println(stringA)
	return md5Hex([]byte(stringA + "&" + md5Hex([]byte(appKey))))
}

func md5Hex(data []byte) string {

	h := md5.New()
	h.Write(data)

	return strings.ToLower(hex.EncodeToString(h.Sum(nil)))
}

func CheckSign(params map[string]interface{}, appKey string) bool {

	signStr := sign(params, appKey)

	return signStr == params["signature"]
}
