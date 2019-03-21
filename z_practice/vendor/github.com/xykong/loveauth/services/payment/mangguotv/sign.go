package mangguotv

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/sirupsen/logrus"
	"io"
	"net/url"
	"sort"
	"strings"
)

func CheckSign(value url.Values, appKey string) bool {

	var keys []string
	for k, _ := range value {

		if k == "sign" {

			continue
		}

		keys = append(keys, k)
	}

	sort.Strings(keys)

	signStr := value.Get("sign")

	var varr []string
	for _, k := range keys {

		varr = append(varr, k+"="+value.Get(k))
	}

	varr = append(varr, "secret_key="+appKey)

	lowStr := strings.ToLower(strings.Join(varr, "&"))

	h := sha1.New()
	io.WriteString(h, lowStr)

	checkStr := strings.ToLower(hex.EncodeToString(h.Sum(nil)))
	if signStr == checkStr {

		return true
	}

	logrus.WithFields(logrus.Fields{
		"sign":   signStr,
		"check":  checkStr,
		"lowStr": lowStr,
	}).Error("mgtv check sign err")

	return false
}
