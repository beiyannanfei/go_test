package payment_v2

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/errors"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"github.com/xykong/loveauth/utils"
)

func init() {

	handlers["/shop_table"] = shopTable
}

//
// Request Shop Table
//
type DoShopTableReq struct {
	//
	// The globalId
	//
	// Required: true
	GlobalId int64 `form:"globalId" json:"globalId" binding:"required"`
	//
	//type
	//
	Type int `form:"type" json:"type"`
}

//
// in: body
// swagger:parameters payment_shop_table
type DoShopTableReqBodyParams struct {
	//
	// swagger:allOf
	// in: body
	DoShopTableReq DoShopTableReq
}

//
// 应答: 协议返回包
// swagger:response DoShopTableRsp
// noinspection ALL
type DoShopTableRsp struct {
	// in: body
	Body struct {
		// The response code
		//
		// Required: true
		Code int64 `json:"code"`
		// The response message
		//
		// Required: true
		Message string `json:"message"`
		// The response ShopTable
		//
		// Required: true
		ShopTable ShopTable `json:"shopTable"`
	}
}

//
// swagger:route POST /payment/shop_table payment payment_shop_table
//
// Return shop table for the given user:
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       200: DoShopTableRsp
func shopTable(c *gin.Context) {

	var request DoShopTableReq

	// validation
	if err := c.BindJSON(&request); err != nil {

		utils.QuickReply(c, errors.Failed, "shopTable BindJSON failed: %v", err)
		return
	}

	shopTable := getShopTable(request.GlobalId, request.Type, time.Now())

	if shopTable == nil {

		utils.QuickReply(c, errors.Failed, "shopTable create failed")
		return
	}

	resp := DoShopTableRsp{}
	resp.Body.Code = int64(errors.Ok)
	resp.Body.Message = "shopTable successfully!"
	resp.Body.ShopTable = *shopTable

	c.JSON(http.StatusOK, resp.Body)
}

func getShopTable(globalId int64, pType int, test time.Time) *ShopTable {

	// build shop table for this user
	// 1. create shop table from shop xlsx
	var shopTable ShopTable

	for _, item := range shopItems {

		var shopItem = &ShopItem{
			ShopId:             item.ShopId,
			ItemId:             item.ItemId,
			ItemCount:          item.ItemCount,
			PriceId:            item.PriceId,
			PriceValue:         item.PriceValue,
			DiscountPriceId:    0,
			DiscountPriceValue: 0,
			GiftPriceId:        0,
			GiftPriceValue:     0,
			GroupId:            0,
			ActivityId:         0,
			Description:        "",
			FinishTime:         0,
			LimitCount:         0,
			AvailableCount:     0,
			Series:             item.Series,
			Icon:               item.Icon,
			ShopTag:            item.ShopTag,
			ProductId:          "NULL",
		}

		switch pType {
		case model.Midas:
			shopItem.ProductId = item.MidasId
		case model.AppStore:
			shopItem.ProductId = item.IosAppProductid
		case model.HUAWEI:
			shopItem.ProductId = item.HuaweiProductid
			shopItem.PriceValue = item.HuaweiPriceValue
		case model.QUICK_KUAIKAN:
			shopItem.ProductId = item.KuaikanProductid
		case model.Mgtv:
			shopItem.ProductId = item.MangguoProductid
		default:
		}

		if typeActivities, ok := shopItemsActivities[item.ShopId]; ok {

			for _, shopActivities := range typeActivities {

				shopActivities = filterByTime(test, shopActivities)

				var usedShopActivity = filterById(shopActivities)

				if shopItem == nil || usedShopActivity == nil {

					continue
				}

				shopItem.ActivityId = usedShopActivity.ActivityId

				switch ActivityType(usedShopActivity.Type) {

				case ActivityTypeEnable:

					if usedShopActivity.Enable == 0 {

						shopItem = nil
					}

				case ActivityTypeGift:

					isFirst := true

					if storage.QuerySameShopIdCount(globalId, item.ShopId, model.OrderStateComplete, test) > 0 {

						isFirst = false
					}

					if isFirst {

						shopItem.GiftPriceId = usedShopActivity.SecondPriceId
						shopItem.GiftPriceValue = usedShopActivity.SecondPriceValue
					} else {

						shopItem.GiftPriceId = usedShopActivity.PriceId
						shopItem.GiftPriceValue = int(usedShopActivity.PriceValue)
					}

					shopItem.GiftDescription = getItemDescription(usedShopActivity, isFirst)

				case ActivityTypeDiscount:

					shopItem.DiscountPriceId = usedShopActivity.PriceId
					shopItem.DiscountPriceValue = usedShopActivity.PriceValue
					shopItem.Description = usedShopActivity.Description

					if pType == model.HUAWEI {

						shopItem.DiscountPriceValue = usedShopActivity.HuaweiPriceValue
					}

					shopItem.LimitCount = usedShopActivity.Limit
					if usedShopActivity.Limit != 0 {

						count := storage.QueryShopIdAvailableCount(globalId, item.ShopId, model.OrderStateComplete, usedShopActivity.Start, usedShopActivity.Finish)
						if usedShopActivity.Limit > count {

							shopItem.AvailableCount = usedShopActivity.Limit - count
						}
					}

				default:

					logrus.WithFields(logrus.Fields{
						"ShopActivity": usedShopActivity,
					}).Warn("ActivityType not implement")

				}

				if shopItem != nil {

					shopItem.ActivityType = ActivityType(usedShopActivity.Type)

					if usedShopActivity.Finish != 0 {

						shopItem.FinishTime = usedShopActivity.Finish
					}
				}
			}
		}

		if shopItem != nil {

			index := shopTable.getTagIndex(item.Tag)
			shopTable.Tags[index].Items = append(shopTable.Tags[index].Items, *shopItem)
		}
	}

	return &shopTable
}

func filterByTime(test time.Time, activities []ShopActivityTemplate) []ShopActivityTemplate {

	if len(activities) == 0 {

		return activities
	}

	var result []ShopActivityTemplate

	for _, activity := range activities {

		if timeInRange(test, activity.Start, activity.Finish) {

			result = append(result, activity)
		}
	}

	// if no time match and ActivityType is ActivityTypeEnable
	//if len(result) == 0 && ActivityType(activities[0].Type) == ActivityTypeEnable {
	//
	//	return activities
	//}

	return result
}

// return the higher Id
func filterById(activities []ShopActivityTemplate) *ShopActivityTemplate {

	var result *ShopActivityTemplate

	for _, activity := range activities {

		if result == nil {

			result = &activity
			continue
		}

		if result.ActivityId < activity.ActivityId {

			result = &activity
			continue
		}
	}

	return result
}

func timeInRange(test time.Time, start, finish int64) bool {

	//if !start.IsZero() && test.Before(start) {
	if start != 0 && test.Unix() < start {

		return false
	}

	if start != 0 && finish != 0 && test.Unix() > finish {

		return false
	}

	return true
}

type ActivityType int

const (
	ActivityTypeInvalid  ActivityType = iota
	ActivityTypeEnable                // 强制上下架 0 下架 1 上架
	ActivityTypeGift                  // 充值多送
	ActivityTypeDiscount              // 现时打折
)

type ShopTable struct {
	//
	// 标签页列表
	//
	Tags []ShopTag `json:"tags"`
}

func (shopTable *ShopTable) getTagIndex(tagId int) int {

	for index, tag := range shopTable.Tags {

		if tag.TagId == tagId {

			return index
		}
	}

	shopTable.Tags = append(shopTable.Tags, ShopTag{TagId: tagId})

	return len(shopTable.Tags) - 1
}

type ShopTag struct {
	//
	// 标签页
	//
	TagId int `json:"tagId"`
	//
	// 标签页
	//
	TagName string `json:"tagName"`
	//
	// 物品列表
	//
	Items []ShopItem `json:"items"`
}

type ShopItem struct {
	//
	// 商品id
	//
	ShopId int `json:"shopId"`
	//
	// 道具id
	//
	ItemId int `json:"itemId"`
	//
	// 道具数量
	//
	ItemCount int `json:"itemCount"`
	//
	// 价格类型
	//
	PriceId int `json:"priceId"`
	//
	// 价格值
	//
	PriceValue float64 `json:"priceValue"`
	//
	// 打折值价格类型
	//
	DiscountPriceId int `json:"discountPriceId"`
	//
	// 打折值
	//
	DiscountPriceValue float64 `json:"discountPriceValue"`
	//
	// 活动奖励价格类型
	//
	GiftPriceId int `json:"giftPriceId"`
	//
	// 活动奖励值
	//
	GiftPriceValue int `json:"giftPriceValue"`
	//
	// 活动组的id
	//
	GroupId int `json:"groupId"`
	//
	// 活动id
	//
	ActivityId int `json:"activityId"`
	//
	// 活动描述
	//
	Description string `json:"description"`
	//
	//钻石充值描述信息
	//
	GiftDescription string `json:"giftDescription"`
	//
	// 活动结束时间
	//
	FinishTime int64 `json:"finishTime"`
	//
	// 限制购买次数
	//
	LimitCount int `json:"limitCount"`
	//
	// 可购买次数
	//
	AvailableCount int `json:"availableCount"`
	//
	// Icon 资源系列
	//
	Series string
	//
	// Icon 资源名
	//
	Icon string
	//
	//活动类型
	//
	ActivityType ActivityType `json:"activityType"`
	//
	//productId
	//
	ProductId string `json:"productId"`
	//
	//midasId
	//
	//MidasId string `json:"midasId"`
	//
	//iosAppProductid
	//
	//IosAppProductid string `json:"iosAppProductid"`
	//
	//shopTag
	//
	ShopTag string `json:"shopTag"`
}

func getItemDescription(usedShopActivity *ShopActivityTemplate, isFrist bool) string {

	if isFrist {

		return fmt.Sprintf(usedShopActivity.Description, "首充", usedShopActivity.SecondPriceValue)
	}

	return fmt.Sprintf(usedShopActivity.Description, "充值", int(usedShopActivity.PriceValue))
}

func GetShopItem(shopId int, globalId int64, payType int, test time.Time) *ShopItem {

	return getShopItem(shopId, globalId, payType, test)
}

func getShopItem(shopId int, globalId int64, payType int, test time.Time) *ShopItem {

	shopTable := getShopTable(globalId, payType, test)

	if shopTable == nil {
		return nil
	}

	for _, tag := range shopTable.Tags {

		for _, item := range tag.Items {

			if item.ShopId == shopId {
				return &item
			}
		}
	}

	return nil
}
