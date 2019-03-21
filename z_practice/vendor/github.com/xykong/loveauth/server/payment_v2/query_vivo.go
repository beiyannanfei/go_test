package payment_v2

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/errors"
	"github.com/xykong/loveauth/services/payment/vivo"
	"github.com/xykong/loveauth/settings"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"time"
)

type QueryVivoOrder struct {
}

func (q *QueryVivoOrder) Name() string {
	return "query_vivo"
}

func (q *QueryVivoOrder) QueryOrder(sequence string, globalId int64) (*model.Order, error) {
	var order *model.Order
	if len(sequence) > 0 {
		order = storage.QueryOrderPlacedWithSequence(sequence)
	} else {
		order = storage.QueryOrderPlaced(globalId)
	}

	if order == nil || order.GlobalId != globalId {
		return nil, errors.NewCodeString(errors.ShopOrderMiss, "query vivo QueryOrderPlacedWithSequence failed.")
	}

	if order.State == model.OrderStateComplete {
		return nil, errors.NewCodeString(errors.ShopQueryGeted, "query QueryOrderPlacedWithSequence failed.")
	}

	if order.State == model.OrderStateFailed {
		return nil, errors.NewCodeString(errors.QuickPayFailed, "query QueryOrder pay failed.")
	}

	return order, nil
}

func (q *QueryVivoOrder) GetShopItem(shopId int, globalId int64, orderType int, timestamp time.Time) *ShopItem {
	shopItem := getShopItem(shopId, globalId, orderType, timestamp)
	if shopItem == nil || shopItem.ItemCount == 0 {
		// something wrong. user have balance without place order. may be shopItem removed.
		logrus.WithFields(logrus.Fields{
			"shopId":    shopId,
			"globalId":  globalId,
			"orderType": orderType,
		}).Warn("query_midas: shopItem without.")
	}

	return shopItem
}

func (q *QueryVivoOrder) DealOrder(request *RequestQueryOrder, order *model.Order, shopItem *ShopItem) (*AcquiredItem, error) {

	if order.State == model.OrderStatePrepare {

		cpId := settings.GetString("lovepay", "vivo.cpId")
		appId := settings.GetString("lovepay", "vivo.appId")
		key := settings.GetString("lovepay", "vivo.key")

		params := map[string]interface{}{
			"cpId":          cpId,
			"appId":         appId,
			"cpOrderNumber": fmt.Sprintf("%v-%v", order.GlobalId, order.Timestamp.Unix()),
			"orderNumber":   order.SNSOrderId,
			"orderAmount":   fmt.Sprintf("%d", order.Amount),
		}

		queryResp, err := vivo.QueryOrder(params, key)
		if err != nil {

			return nil, errors.NewCodeString(errors.ShopDataException, "queryorder err %v", err)
		}

		order.State = model.OrderStatePlace
		order.SNSOrderId = queryResp.OrderNumber
	}

	if order.State != model.OrderStatePlace {

		logrus.WithFields(logrus.Fields{
			"resp": order.State,
		}).Error("queryvivoorder state err.")

		return nil, errors.NewCodeString(errors.ShopDataException, "queryorder status err %v", order.State)
	}

	acquiredItem := &AcquiredItem{
		ShopId:         shopItem.ShopId,
		ItemId:         shopItem.ItemId,
		ItemCount:      shopItem.ItemCount,
		ShopActivityId: shopItem.ActivityId,
		CostCount:      float64(order.Amount) / 100,
		ProductId:      shopItem.ProductId,
		Paymethod:      order.PayMethod,
		Sequence:       order.Sequence,
	}

	if order != nil && shopItem != nil && shopItem.GiftPriceValue > 0 {

		acquiredItem.GiftId = shopItem.GiftPriceId
		acquiredItem.GiftCount = shopItem.GiftPriceValue
	}

	acquiredItem.TransactionId = order.SNSOrderId

	return acquiredItem, nil
}

func (q *QueryVivoOrder) SaveOrder(order *model.Order) {
	order.State = model.OrderStateComplete
	storage.Save(storage.PayDatabase(), order)
}
