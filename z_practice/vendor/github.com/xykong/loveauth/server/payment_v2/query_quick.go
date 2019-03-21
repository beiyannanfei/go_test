package payment_v2

import (
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/errors"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"time"
)

type QueryQuickPay struct {
}

func (q *QueryQuickPay) Name() string {
	return "query_quick"
}

func (q *QueryQuickPay) QueryOrder(sequence string, globalId int64) (*model.Order, error) {
	var order *model.Order
	if len(sequence) > 0 {
		order = storage.QueryOrderPlacedWithSequence(sequence)
	} else {
		order = storage.QueryOrderPlaced(globalId)
	}

	if order == nil || order.GlobalId != globalId {
		return nil, errors.NewCodeString(errors.ShopOrderMiss, "query quick QueryOrderPlacedWithSequence failed.")
	}

	if order.State == model.OrderStateComplete {
		return nil, errors.NewCodeString(errors.ShopQueryGeted, "query QueryOrderPlacedWithSequence failed.")
	}

	if order.State == model.OrderStateFailed { //回调返回支付失败
		return nil, errors.NewCodeString(errors.QuickPayFailed, "query QueryOrder pay failed.")
	}

	return order, nil
}

func (q *QueryQuickPay) GetShopItem(shopId int, globalId int64, orderType int, timestamp time.Time) *ShopItem {
	shopItem := getShopItem(shopId, globalId, orderType, timestamp)
	if shopItem == nil || shopItem.ItemCount == 0 {
		// something wrong. user have balance without place order. may be shopItem removed.
		logrus.WithFields(logrus.Fields{
			"shopId":    shopId,
			"globalId":  globalId,
			"orderType": orderType,
		}).Warn("query_quick: shopItem without.")
	}

	return shopItem
}

func (q *QueryQuickPay) DealOrder(request *RequestQueryOrder, order *model.Order, shopItem *ShopItem) (*AcquiredItem, error) {
	var acquiredItem *AcquiredItem

	quickOrder := storage.QueryQuickPayOrder(order.Sequence)
	if quickOrder == nil || order.State == model.OrderStatePrepare { //第三方回调还没到
		return nil, errors.NewCodeString(errors.ShopCallBackAgain, "query DealOrder quickOrder null")
	}

	if quickOrder.Status != "0" || order.State != model.OrderStatePlace {
		return nil, errors.NewCodeString(errors.QuickPayFailed, "query DealOrder pay failed.")
	}

	acquiredItem = &AcquiredItem{
		ShopId:         shopItem.ShopId,
		ItemId:         shopItem.ItemId,
		ItemCount:      shopItem.ItemCount,
		GiftCount:      0,
		TransactionId:  quickOrder.OrderNo,
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

	return acquiredItem, nil
}

func (q *QueryQuickPay) SaveOrder(order *model.Order) {
	order.State = model.OrderStateComplete
	storage.Save(storage.PayDatabase(), order)
}
