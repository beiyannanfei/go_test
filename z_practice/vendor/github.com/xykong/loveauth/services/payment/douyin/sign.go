package douyin

import (
	"crypto"
	"encoding/base64"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/services/payment/alipay/encoding"
	"net/url"
	"sort"
	"strings"
)

func CheckSign(value url.Values, payKey string) bool {

	var keys []string
	for k, _ := range value {

		if k == "tt_sign" || k == "tt_sign_type" {

			continue
		}

		keys = append(keys, k)
	}

	sort.Strings(keys)

	var varr []string
	for _, k := range keys {

		if value.Get(k) == "" {

			continue
		}

		varr = append(varr, k+"="+value.Get(k))
	}

	//varr = append(varr, "key="+payKey)

	originalData := strings.Join(varr, "&")

	key := encoding.ParsePublicKey(`MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDOZZ7iAkS3oN970+yDONe5TPhPrLHoNOZOjJjackEtgbptdy4PYGBGdeAUAz75TO7YUGESCM+JbyOz1YzkMfKl2HwYdoePEe8qzfk5CPq6VAhYJjDFA/M+BAZ6gppWTjKnwMcHVK4l2qiepKmsw6bwf/kkLTV9l13r6Iq5U+vrmwIDAQAB`)
	signBytes, err := base64.StdEncoding.DecodeString(value.Get("tt_sign"))
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"sign":         value.Get("tt_sign"),
			"originalData": originalData,
			"err":          err,
		}).Error("douyin base64.StdEncoding.DecodeString err")

		return false
	}

	err = encoding.VerifyPKCS1v15([]byte(originalData), signBytes, key, crypto.SHA1)
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"sign":         value.Get("tt_sign"),
			"originalData": originalData,
			"err":          err,
		}).Error("douyin RsaVerySignWithSha1Base64 err")

		return false
	}

	return true
}
