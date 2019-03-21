package payment_v2

import (
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/errors"
	"github.com/xykong/loveauth/server/watch_waring"
	"github.com/xykong/loveauth/services/payment/qpay"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"time"
)

type QueryQPay struct {
}

func (q *QueryQPay) Name() string {
	return "query_qpay"
}

func (q *QueryQPay) QueryOrder(sequence string, globalId int64) (*model.Order, error) {
	var order *model.Order
	if len(sequence) > 0 {
		order = storage.QueryOrderPlacedWithSequence(sequence)
	} else {
		order = storage.QueryOrderPlaced(globalId)
	}

	if order == nil || order.GlobalId != globalId {
		return nil, errors.NewCodeString(errors.ShopOrderMiss, "query qpay QueryOrderPlacedWithSequence failed.")
	}

	if order.State == model.OrderStateComplete {
		return nil, errors.NewCodeString(errors.ShopQueryGeted, "query qpay QueryOrderPlacedWithSequence failed.")
	}

	if order.State == model.OrderStateFailed { //回调返回支付失败
		return nil, errors.NewCodeString(errors.QuickPayFailed, "query qpay QueryOrder pay failed.")
	}

	return order, nil
}

func (q *QueryQPay) GetShopItem(shopId int, globalId int64, orderType int, timestamp time.Time) *ShopItem {
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

func (q *QueryQPay) DealOrder(request *RequestQueryOrder, order *model.Order, shopItem *ShopItem) (*AcquiredItem, error) {
	var acquiredItem *AcquiredItem

	acquiredItem = &AcquiredItem{
		ShopId:         shopItem.ShopId,
		ItemId:         shopItem.ItemId,
		ItemCount:      shopItem.ItemCount,
		GiftCount:      0,
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

	qPayOrder := storage.QueryQPayOrderWithSequence(order.Sequence)
	var respQuery *qpay.ResponseOrderQuery
	var err error
	if qPayOrder == nil {
		respQuery, err = qpay.QPayOrderQuery(order)
		if err != nil {
			return nil, err
		}

		acquiredItem.TransactionId = respQuery.TransactionId
	} else {
		acquiredItem.TransactionId = qPayOrder.TransactionId
	}

	order.SNSOrderId = acquiredItem.TransactionId

	if order.State == model.OrderStatePlace { //支付成功(已经回调)
		return acquiredItem, nil
	}

	if order.State == model.OrderStatePrepare { //未支付
		if respQuery == nil { //查询返回nil
			return nil, errors.NewCodeString(errors.ShopCallBackAgain, "query DealOrder QueryQPay respQuery nil")
		}

		if respQuery.TradeState == "SUCCESS" { //支付成功
			watch_waring.PaymentWatch(order.GlobalId, order.Vendor, order.Amount)
			return acquiredItem, nil
		} else if respQuery.TradeState == "USERPAYING" { //用户支付中
			return nil, errors.NewCodeString(errors.ShopCallBackAgain, "query DealOrder QueryQPay USERPAYING state")
		} else { //支付失败
			order.State = model.OrderStateFailed
			return nil, errors.NewCodeString(errors.QuickPayFailed, "query QueryOrder pay failed.")
		}

	} else {
		return nil, errors.NewCodeString(errors.QuickPayFailed, "query QueryOrder pay failed.")
	}
}

func (q *QueryQPay) SaveOrder(order *model.Order) {
	order.State = model.OrderStateComplete
	storage.Save(storage.PayDatabase(), order)
}
