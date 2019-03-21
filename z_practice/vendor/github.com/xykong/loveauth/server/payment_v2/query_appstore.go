package payment_v2

import (
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/errors"
	"github.com/xykong/loveauth/server/watch_waring"
	"github.com/xykong/loveauth/services/payment/app_store"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"time"
)

type QueryAppStoreOrder struct {
}

func (q *QueryAppStoreOrder) Name() string {
	return "query_appstore"
}

func (q *QueryAppStoreOrder) QueryOrder(sequence string, globalId int64) (*model.Order, error) {
	var order *model.Order
	if len(sequence) > 0 {
		order = storage.QueryOrderPlacedWithSequence(sequence)
	} else {
		order = storage.QueryOrderPlaced(globalId)
	}

	if order == nil || order.GlobalId != globalId {
		return nil, errors.NewCodeString(errors.ShopOrderMiss, "query appStore QueryOrderPlacedWithSequence failed.")
	}

	if order.State == model.OrderStateComplete {
		return nil, errors.NewCodeString(errors.ShopQueryGeted, "query QueryOrderPlacedWithSequence failed.")
	}

	if order.State == model.OrderStateFailed {
		return nil, errors.NewCodeString(errors.QuickPayFailed, "query QueryOrder pay failed.")
	}

	return order, nil
}

func (q *QueryAppStoreOrder) GetShopItem(shopId int, globalId int64, orderType int, timestamp time.Time) *ShopItem {
	shopItem := getShopItem(shopId, globalId, orderType, timestamp)
	if shopItem == nil || shopItem.ItemCount == 0 {
		// something wrong. user have balance without place order. may be shopItem removed.
		logrus.WithFields(logrus.Fields{
			"shopId":    shopId,
			"globalId":  globalId,
			"orderType": orderType,
		}).Warn("query_appStore: shopItem without.")
	}

	return shopItem
}

func (q *QueryAppStoreOrder) DealOrder(request *RequestQueryOrder, order *model.Order, shopItem *ShopItem) (*AcquiredItem, error) {
	var acquiredItem *AcquiredItem
	resp, err := app_store.VerifyReceipt(app_store.DoVerifyReceiptReq{
		ReceiptData: request.AppStoreReceipt,
	})
	if err != nil {
		return nil, errors.NewCodeString(errors.ShopDataException, "verifyreceipt err %v", err)
	}

	isOk := false
	for _, app := range resp.Body.Receipt.InApp {
		if app.ProductId == shopItem.ProductId {
			//if order.Type == model.OrderTypeDiamond {

			if storage.QueryOrderPlacedWithAppStoreNumber(app.TransactionId) != nil {

				continue
			}

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

			acquiredItem.TransactionId = app.TransactionId
			order.SNSOrderId = app.TransactionId

			isOk = true
			break
		}
	}

	if !isOk {
		return nil, errors.NewCodeString(errors.ShopDataException, "shop_productid(%v)not int app_store resp productIds %v, orderNumber = %v", shopItem.ProductId, resp.Body.Receipt.InApp, request.OrderNumber)
	}
	watch_waring.PaymentWatch(order.GlobalId, order.Vendor, order.Amount)
	return acquiredItem, nil
}

func (q *QueryAppStoreOrder) SaveOrder(order *model.Order) {
	order.State = model.OrderStateComplete
	storage.Save(storage.PayDatabase(), order)
}
