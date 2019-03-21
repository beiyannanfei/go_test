package payment_v2

import (
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/errors"
	"github.com/xykong/loveauth/server/watch_waring"
	"github.com/xykong/loveauth/services/payment/bilibili"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"time"
)

type QueryBilibiliPay struct {
}

func (q *QueryBilibiliPay) Name() string {
	return "query_bilibili"
}

func (q *QueryBilibiliPay) QueryOrder(sequence string, globalId int64) (*model.Order, error) {
	var order *model.Order
	if len(sequence) > 0 {
		order = storage.QueryOrderPlacedWithSequence(sequence)
	} else {
		order = storage.QueryOrderPlaced(globalId)
	}

	if nil == order || globalId != order.GlobalId {
		return nil, errors.NewCodeString(errors.ShopOrderMiss, "query bilibili QueryOrderPlacedWithSequence failed")
	}

	if model.OrderStateComplete == order.State {
		return nil, errors.NewCodeString(errors.ShopQueryGeted, "query bilibili QueryOrderPlacedWithSequence already get")
	}

	if model.OrderStateFailed == order.State {
		return nil, errors.NewCodeString(errors.QuickPayFailed, "query bilibili QueryOrderPlacedWIthSequence pay failed")
	}
	return order, nil
}

func (q *QueryBilibiliPay) GetShopItem(shopId int, globalId int64, orderType int, timestamp time.Time) *ShopItem {
	shopItem := getShopItem(shopId, globalId, orderType, timestamp)
	if nil == shopItem || 0 == shopItem.ItemCount {
		logrus.WithFields(logrus.Fields{
			"shopId":    shopId,
			"globalId":  globalId,
			"orderType": orderType,
		}).Warn("bilibili: shopItem without")
	}
	return shopItem
}

func (q *QueryBilibiliPay) DealOrder(request *RequestQueryOrder, order *model.Order, shopItem *ShopItem) (*AcquiredItem, error) {
	var acquiredItem *AcquiredItem
	acquiredItem = &AcquiredItem{
		ShopId:         shopItem.ShopId,
		ItemId:         shopItem.ItemId,
		ItemCount:      shopItem.ItemCount,
		GiftCount:      0,
		TransactionId:  request.OrderNumber, // 这个是第三方的账单号
		ShopActivityId: shopItem.ActivityId,
		CostCount:      float64(order.Amount) / 100,
		ProductId:      shopItem.ProductId,
		Paymethod:      order.PayMethod,
		Sequence:       order.Sequence,
	}

	if nil != order && nil != shopItem && shopItem.GiftPriceValue > 0 {
		acquiredItem.GiftId = shopItem.GiftPriceId
		acquiredItem.GiftCount = shopItem.GiftPriceValue
	}

	order.SNSOrderId = request.OrderNumber

	if order.State == model.OrderStatePlace { //支付完成
		return acquiredItem, nil
	}

	profile, err := storage.QueryProfile(order.GlobalId)
	if err != nil {
		return nil, errors.NewCodeString(errors.ShopCallBackAgain, "query DealOrder bilibiliOrder null")
	}

	//支付未完成，主动查询
	respQuery, err := bilibili.QueryBiliOrder(request.OrderNumber, int(profile.Auth.VendorBilibili.LoginResult.UserId))
	if err != nil {
		return nil, err
	}

	if respQuery == nil { //查询返回nil
		return nil, errors.NewCodeString(errors.ShopCallBackAgain, "query DealOrder QueryBili respQuery nil")
	}

	if respQuery.OrderStatus == 1 { //支付成功
		watch_waring.PaymentWatch(order.GlobalId, order.Vendor, order.Amount)
		order.SNSOrderId = respQuery.OrderNo
		return acquiredItem, nil
	} else if respQuery.OrderStatus == 3 { //用户支付中
		return nil, errors.NewCodeString(errors.ShopCallBackAgain, "query DealOrder QueryBili USERPAYING state")
	} else { //支付失败
		order.State = model.OrderStateFailed
		storage.Save(storage.PayDatabase(), order)
		return nil, errors.NewCodeString(errors.QuickPayFailed, "query QueryOrder pay failed.")
	}
}

func (q *QueryBilibiliPay) SaveOrder(order *model.Order) {
	order.State = model.OrderStateComplete
	err := storage.Save(storage.PayDatabase(), order)
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("failed bilibili SaveOrder")
	}
}
