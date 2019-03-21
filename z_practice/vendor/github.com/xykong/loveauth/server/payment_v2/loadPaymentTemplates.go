package payment_v2

import (
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/utils"
)

var shopItems []ShopItemTemplate
var shopActivities []ShopActivityTemplate

// ShopId, ActivityType, Activities
var shopItemsActivities map[int]map[ActivityType][]ShopActivityTemplate

type ShopItemTemplate struct {
	ShopId           int
	Tag              int
	ItemId           int
	ItemCount        int
	PriceId          int
	PriceValue       float64
	HuaweiPriceValue float64
	Series           string
	Icon             string
	MidasId          string
	IosAppProductid  string
	HuaweiProductid  string
	XiaomiProductid  string
	KuaikanProductid string
	MangguoProductid string
	ShopTag          string
}

type ShopActivityTemplate struct {
	ActivityId       int
	ActivityGroupId  int
	ShopId           int
	Type             int
	PriceId          int
	PriceValue       float64
	HuaweiPriceValue float64
	Description      string
	Start            int64
	Finish           int64
	Limit            int
	SecondPriceId    int
	SecondPriceValue int
	Enable           int
}

func loadPaymentTemplates() {

	err := utils.LoadTemplates("shop.xlsx", &shopItems)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("loadPaymentTemplates load failed.")
	}

	err = utils.LoadTemplates("shop_activity.xlsx", &shopActivities)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("loadPaymentTemplates load failed.")
	}

	shopItemsActivities = make(map[int]map[ActivityType][]ShopActivityTemplate)
	for _, shopActivity := range shopActivities {

		if shopItemsActivities[shopActivity.ShopId] == nil {

			shopItemsActivities[shopActivity.ShopId] = make(map[ActivityType][]ShopActivityTemplate)
		}

		shopItemsActivities[shopActivity.ShopId][ActivityType(shopActivity.Type)] =
			append(shopItemsActivities[shopActivity.ShopId][ActivityType(shopActivity.Type)], shopActivity)
	}
}
