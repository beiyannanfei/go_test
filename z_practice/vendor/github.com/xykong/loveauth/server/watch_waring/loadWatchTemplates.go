package watch_waring

import (
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/utils"
)

type WatchTemplate struct {
	Id   string
	Add  int
	Sub  int
	Desc string
	Type int
}


var itemWatchMap map[string]WatchTemplate
var failedOrderWaring = -1
var deviceWatch = -1
var accountWatch = -1
var paymentWatch = -1

func loadWatchTemplates() {
	var watchConfigList []WatchTemplate
	err := utils.LoadTemplates("waring_config.xlsx", &watchConfigList)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("loadPaymentTemplates load failed.")
	}

	for _, t := range watchConfigList {
		switch t.Type {
		case 1: {
				if nil == itemWatchMap {
					itemWatchMap = make(map[string]WatchTemplate)
				}
				itemWatchMap[t.Id] = t
			}
		case 2: {
			failedOrderWaring = t.Add
		}
		case 3 : {
			deviceWatch = t.Add
		}
		case 4 : {
			accountWatch = t.Add
		}
		case 5 : {
			paymentWatch = t.Add
		}
		default:
			 logrus.Info("not exist watch or waring type")
		}
	}
}

func GetItemWatchMap() map[string]WatchTemplate {
	return itemWatchMap
}

func GetFailedOrderWaring() int {
	return failedOrderWaring
}

func GetDeviceWatch() int {
	return deviceWatch
}

func GetAccountWatch() int {
	return accountWatch
}

func GetPaymentWatch() int {
	return paymentWatch
}