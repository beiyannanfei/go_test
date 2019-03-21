package payment_v2

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/errors"
	"github.com/xykong/loveauth/server/send_mail"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"github.com/xykong/loveauth/utils"
	"net/http"
	"time"
)

type QueryProvider interface {
	Name() string
	QueryOrder(sequence string, globalId int64) (*model.Order, error)
	GetShopItem(shopId int, globalId int64, orderType int, timestamp time.Time) *ShopItem
	DealOrder(request *RequestQueryOrder, order *model.Order, shopItem *ShopItem) (*AcquiredItem, error)
	SaveOrder(order *model.Order)
}

type QueryProviders map[string]QueryProvider

var usedQueryProviders = QueryProviders{}

type Query struct {
	QueryProvider QueryProvider
}

func UseQuery(aps ...QueryProvider) {
	for _, QueryProvider := range aps {
		if usedQueryProviders[QueryProvider.Name()] != nil {
			logrus.WithFields(logrus.Fields{
				"QueryProvider": QueryProvider.Name(),
			}).Warn("QueryProvider replaced.")
		}

		usedQueryProviders[QueryProvider.Name()] = QueryProvider
	}
}

type RequestQueryOrder struct {
	GlobalId        int64      `form:"globalId" json:"globalId" binding:"required"`
	Sequence        string     `form:"sequence" json:"sequence"`
	QueryMidas      QueryMidas `form:"queryMidas" json:"queryMidas"`
	Type            int        `form:"type" json:"type"`
	AppStoreReceipt string     `form:"appStoreReceipt" json:"appStoreReceipt"`
	OrderNumber     string     `form:"orderNumber" json:"orderNumber"`
}

type ResponseQueryOrder struct {
	Body struct {
		Code          int64        `json:"code"`
		Message       string       `json:"message"`
		AcquiredItems AcquiredItem `json:"acquiredItems"`
	}
}

type AcquiredItem struct {
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
	// 活动奖励道具id
	//
	GiftId int `json:"giftId"`
	//
	// 活动奖励数量
	//
	GiftCount int `json:"giftCount"`
	//
	//第三方订单号
	//
	TransactionId string `json:"transaction_id"`
	//
	//shop_activity Id
	//
	ShopActivityId int `json:"shopActivityId"`
	//
	// 费用，单位元
	//
	CostCount float64 `json:"cost_count"`
	//
	// productid
	//
	ProductId string `json:"productId"`
	//
	// paymethod
	//
	Paymethod int `json:"paymethod"`
	//
	// paymethod
	//
	Sequence string `json:"sequence"`
}

func (q *Query) QueryOrder(c *gin.Context) {
	var request RequestQueryOrder

	// validation
	if err := c.BindJSON(&request); err != nil {
		LogCurrencyPay(request.GlobalId, nil, 0, "NULL", "NULL", 0)
		utils.QuickReply(c, errors.Failed, "query BindJSON failed: %v", err)
		return
	}

	order, err := q.QueryProvider.QueryOrder(request.Sequence, request.GlobalId)
	if err != nil {
		LogCurrencyPay(request.GlobalId, nil, request.Type, "NULL", "NULL", 0)
		ec, ok := err.(*errors.Type)
		if !ok { //系统错误
			utils.QuickReply(c, errors.Failed, "QueryOrder err: %v", err)
			return
		}

		//逻辑类错误
		utils.QuickReply(c, ec.Code, ec.Message)
		return
	}

	if order == nil {
		LogCurrencyPay(request.GlobalId, nil, request.Type, "NULL", "NULL", 0)
		// something wrong. user have balance without place order. may be given by midas custom service.
		logrus.WithFields(logrus.Fields{
			"request": request,
		}).Warn("query: order without")

		utils.QuickReply(c, errors.ShopOrderMiss, "query order failed.")
		return
	}

	shopItem := q.QueryProvider.GetShopItem(order.ShopId, order.GlobalId, request.Type, order.Timestamp)
	if shopItem == nil {
		LogCurrencyPay(request.GlobalId, shopItem, request.Type, "NULL", order.Sequence, 0)
		utils.QuickReply(c, errors.ShopIdMiss, "ShopIdMiss")
		return
	}

	acquiredItem, err := q.QueryProvider.DealOrder(&request, order, shopItem)
	if err != nil {

		acquiredItem, order = QueryByOrderNumber(request.OrderNumber, request.Type, request.GlobalId)
		if acquiredItem == nil || order == nil {

			LogCurrencyPay(request.GlobalId, shopItem, request.Type, "NULL", order.Sequence, 0)
			ec, ok := err.(*errors.Type)
			if !ok { //系统错误
				utils.QuickReply(c, errors.Failed, "QueryOrder DealOrder err: %v", err)
				return
			}

			//逻辑类错误
			utils.QuickReply(c, ec.Code, ec.Message)
			return
		}
	}

	q.QueryProvider.SaveOrder(order)
	LogCurrencyPay(request.GlobalId, shopItem, request.Type, acquiredItem.TransactionId, order.Sequence, 1)
	send_mail.DoOrderWaring(storage.QuerySequence, order.Sequence, order.Vendor)
	// 发货
	resp := ResponseQueryOrder{}
	resp.Body.Code = int64(errors.Ok)
	resp.Body.Message = "query successfully"
	resp.Body.AcquiredItems = *acquiredItem
	c.JSON(http.StatusOK, resp.Body)
	return
}

func QueryByOrderNumber(orderNumber string, reqType int, globalId int64) (*AcquiredItem, *model.Order) {

	if orderNumber == "" {

		return nil, nil
	}

	order := storage.QueryOrderPlacedWithOrderNumber(orderNumber, model.OrderStatePlace)
	if order == nil || globalId != order.GlobalId {

		return nil, nil
	}

	shopItem := getShopItem(order.ShopId, order.GlobalId, reqType, order.Timestamp)
	if shopItem == nil {

		return nil, nil
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
	order.State = model.OrderStateComplete

	return acquiredItem, order
}

func StartQuery(group *gin.RouterGroup) {
	UseQuery(&QueryAppStoreOrder{}, &QueryBilibiliPay{}, &QueryMidasOrder{}, &QueryQuickPay{},
		&QueryVivoOrder{}, &QueryWechat{}, &QueryAlipayOrder{}, &QueryQPay{}, &QueryHuawei{}, &QueryOfficialOrder{}, &QueryDouyinOrder{},
		&QueryMgtvOrder{})

	for _, q := range usedQueryProviders {
		var QueryProvider = q
		group.POST("/"+q.Name(), func(c *gin.Context) {
			var query = Query{QueryProvider}
			query.QueryOrder(c)
		})
	}
}
