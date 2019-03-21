package quick_sdk

import (
	"strings"
	"github.com/xykong/loveauth/settings"
	"strconv"
)

func DecodeQuickCb(ntData string) string {
	callback_Key := settings.GetString("lovepay", "quickSdk.Callback_Key")
	dataList := strings.Split(ntData, "@")[1:]
	keysByte := []byte(callback_Key)
	var dataByte []byte

	for i := 0; i < len(dataList); i++ {
		l, _ := strconv.Atoi(dataList[i])
		k := keysByte[i%len(keysByte)]
		cha := l - int(k)
		dataByte = append(dataByte, byte(cha))
	}

	return string(dataByte)
}
