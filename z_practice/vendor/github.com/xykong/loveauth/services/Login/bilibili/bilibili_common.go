package bilibili

import (
	"github.com/xykong/loveauth/settings"
	"net/url"
	"sort"
	"strconv"
	"time"
	"fmt"
	"crypto/md5"
)

const TimeFormat = "Mon, 02 Jan 2006 15:04:05 GMT"

func GetCommonReqParams(uid int32) url.Values {
	p := url.Values{}
	p.Add("game_id", settings.GetString("lovepay", "bilibili.gameId"))
	p.Add("merchant_id", settings.GetString("lovepay", "bilibili.merchantId"))
	p.Add("server_id", settings.GetString("lovepay", "bilibili.serverId"))
	p.Add("uid", strconv.Itoa(int(uid)))
	p.Add("version", settings.GetString("lovepay", "bilibili.version"))
	p.Add("timestamp", strconv.Itoa(int(time.Now().UnixNano()/1000000)))
	return p
}

func GetSign(p url.Values) string {
	var keyList []string
	for key, _ := range p {
		if "item_name" == key || "item_desc" == key || "sign" == key {
			continue
		}
		keyList = append(keyList, key)
	}
	sort.Strings(keyList)

	paramsStr := ""
	for _, k := range keyList {
		paramsStr += p.Get(k)
	}
	paramsStr = paramsStr + settings.GetString("lovepay", "bilibili.secretKey")

	return fmt.Sprintf("%x", md5.Sum([]byte(paramsStr)))
}
