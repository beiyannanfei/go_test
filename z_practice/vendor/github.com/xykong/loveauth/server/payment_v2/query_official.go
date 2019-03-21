package payment_v2

import (
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/errors"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"time"
)

type QueryOfficialOrder struct {
}

func (q *QueryOfficialOrder) Name() string {
	return "query_official"
}

func (q *QueryOfficialOrder) QueryOrder(sequence string, globalId int64) (*model.Order, error) {
	var order *model.Order
	if len(sequence) > 0 {
		order = storage.QueryOrderPlacedWithSequence(sequence)
	} else {
		order = storage.QueryOrderPlaced(globalId)
	}

	if order == nil || order.GlobalId != globalId {
		return nil, errors.NewCodeString(errors.ShopOrderMiss, "query official QueryOrderPlacedWithSequence failed.")
	}

	if order.State == model.OrderStateComplete {
		return nil, errors.NewCodeString(errors.ShopQueryGeted, "query QueryOrderPlacedWithSequence failed.")
	}

	if order.State == model.OrderStateFailed {
		return nil, errors.NewCodeString(errors.QuickPayFailed, "query QueryOrder pay failed.")
	}

	return order, nil
}

func (q *QueryOfficialOrder) GetShopItem(shopId int, globalId int64, orderType int, timestamp time.Time) *ShopItem {
	shopItem := getShopItem(shopId, globalId, orderType, timestamp)
	if shopItem == nil || shopItem.ItemCount == 0 {
		// something wrong. user have balance without place order. may be shopItem removed.
		logrus.WithFields(logrus.Fields{
			"shopId":    shopId,
			"globalId":  globalId,
			"orderType": orderType,
		}).Warn("query_official: shopItem without.")
	}

	return shopItem
}

func (q *QueryOfficialOrder) DealOrder(request *RequestQueryOrder, order *model.Order, shopItem *ShopItem) (*AcquiredItem, error) {

	switch order.PayMethod {

	case model.Alipay:

		var alipayOrder QueryAlipayOrder
		return alipayOrder.DealOrder(request, order, shopItem)
	case model.Wechat:

		var wechatOrder QueryWechat
		return wechatOrder.DealOrder(request, order, shopItem)
	case model.QPay:

		var qOrder QueryQPay
		return qOrder.DealOrder(request, order, shopItem)
	}

	return nil, errors.New("paymethod not match")
}

func (q *QueryOfficialOrder) SaveOrder(order *model.Order) {
	order.State = model.OrderStateComplete
	storage.Save(storage.PayDatabase(), order)
}
