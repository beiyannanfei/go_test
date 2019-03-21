package midas

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"regexp"

	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/settings"
	"github.com/xykong/loveauth/storage/model"

	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sort"
	"strings"
	"time"
)

/**
//	session_id		《通用参数说明》
//	session_type	《通用参数说明》
//	org_loc
//	需要填写: /v3/r/mpay/get_balance_m （注意：如果经过接口机器，需要填写应用签名时使用的URI）
//
//	appip	（可选）来源的第三方应用的服务IP
*/
func makeCookies(urlPath string, appIp string, vendor model.Vendor) string {

	//noinspection ALL
	params := map[string]interface{}{
		"org_loc": urlPath,
		"appip":   appIp,
	}

	switch vendor {

	case model.VendorMsdkQQ:

		params["session_id"] = "openid"
		params["session_type"] = "kp_actoken"

	case model.VendorMsdkGuest:

		params["session_id"] = "hy_gameid"
		params["session_type"] = "st_dummy"

	case model.VendorMsdkWechat:

		params["session_id"] = "hy_gameid"
		params["session_type"] = "wc_actoken"
	}

	var cookies []string
	for k, v := range params {
		cookies = append(cookies, fmt.Sprintf("%v=%v", k, v))
	}

	return strings.Join(cookies, ";")
}

// http://openapi.sparta.html5.qq.com/v3/user/get_info?
// openid=11111111111111111&openkey=2222222222222222&
// appid=123456&pf=qzone&format=json&userip=112.90.139.30&sig=FdJkiDYwMj5Aj1UG2RUPc83iokk%3D
func makeUrl(method string, urlPath string, params map[string]interface{}, secret string) string {

	host := host()

	sig := makeSig(method, urlPath, params, secret)

	query := makeQuery(params, true)

	return fmt.Sprintf("%v%v?%v&sig=%v", host, urlPath, query, urlEncode(sig))
}

func host() string {

	return settings.GetString("tencent", "midas.host")
}

/**
 * 生成签名
 *
 * @param string 	method 请求方法 "get" or "post"
 * @param string 	url_path
 * @param array 	params 表单参数
 * @param string 	secret 密钥
 */
//noinspection ALL
func makeSig(method string, urlPath string, params map[string]interface{}, secret string) string {
	mk := makeSource(method, urlPath, params)

	//hmac ,use sha1
	key := []byte(secret)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(mk))

	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func makeSource(method string, urlPath string, params map[string]interface{}) string {

	encoded := strings.ToUpper(method) + "&" + urlEncode(urlPath) + "&"

	result := makeQuery(params, false)
	result = urlEncode(result)
	result = strings.Replace(result, "~", "%7E", -1)

	return encoded + result
}

func makeQuery(params map[string]interface{}, encode bool) string {

	keys := sortByKey(params)

	var queryString []string
	for _, k := range keys {
		v := params[k]
		if encode {
			v = urlEncode(fmt.Sprintf("%v", v))
		}

		queryString = append(queryString, fmt.Sprintf("%v=%v", k, v))
	}

	result := strings.Join(queryString, "&")

	return result
}

func sortByKey(params map[string]interface{}) []string {

	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return keys
}

func urlEncode(s string) (result string) {

	result = url.QueryEscape(s)

	result = strings.Replace(result, "+", "%20", -1)

	return result
}

func SendRequest(method string, urlPath string, params map[string]interface{}, platId int) ([]byte, error) {

	start := time.Now()
	settingKey := fmt.Sprintf("midas.%v.currency.", platId)
	appId := settings.GetString("tencent", settingKey+"AppId")
	appKey := settings.GetString("tencent", settingKey+"AppKey")

	vendor := getVendorByPf(params["openid"].(string))
	//noinspection SpellCheckingInspection
	params["appid"] = appId
	params["ts"] = time.Now().Unix()
	//params["ts"] = 1522292397
	//noinspection SpellCheckingInspection
	params["zoneid"] = getZoneIdByVendor(vendor)
	//noinspection SpellCheckingInspection
	params["accounttype"] = "common"
	params["format"] = "json"

	secret := appKey + "&"

	urlString := makeUrl(method, urlPath, params, secret)

	cookieJar, _ := cookiejar.New(nil)

	client := &http.Client{
		Jar: cookieJar,
	}

	req, err := http.NewRequest("GET", urlString, nil)

	cookie := makeCookies(urlPath, "", vendor)
	req.Header.Set("Cookie", cookie)

	resp, err := client.Do(req)
	//resp, err := client.Get(urlString)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":      err,
			"elapsed":    time.Since(start),
			"appId":      appId,
			"appKey":     appKey,
			"urlPath":    urlPath,
			"settingKey": settingKey,
			"urlString":  urlString,
			"cookie":     cookie,
			"params":     params,
		}).Error("midas send request")

		return nil, err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	logrus.WithFields(logrus.Fields{
		"elapsed":    time.Since(start),
		"appId":      appId,
		"appKey":     appKey,
		"urlPath":    urlPath,
		"settingKey": settingKey,
		"body":       string(body),
		"urlString":  urlString,
		"cookie":     cookie,
		"params":     params,
	}).Info("midas send request")

	return body, nil
}

func getVendorByPf(openId string) model.Vendor {

	match, _ := regexp.MatchString(`^oNB.*`, openId)
	if match {

		return model.VendorMsdkWechat
	}

	match, _ = regexp.MatchString(`^G_.*`, openId)
	if match {

		return model.VendorMsdkGuest
	}

	return model.VendorMsdkQQ
}

func getZoneIdByVendor(vendor model.Vendor) string {

	if vendor == model.VendorMsdkQQ {

		return settings.GetString("tencent", "midas.qq_zoneid")
	}

	if vendor == model.VendorMsdkWechat {

		return settings.GetString("tencent", "midas.wechat_zoneid")
	}

	return settings.GetString("tencent", "midas.default_zoneid")
}
