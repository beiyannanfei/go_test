package payment_v2

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/errors"
	"github.com/xykong/loveauth/services/payment/bilibili"
	"github.com/xykong/loveauth/services/payment/qpay"
	"github.com/xykong/loveauth/services/payment/vivo"
	"github.com/xykong/loveauth/services/payment/wechat_pay"
	"github.com/xykong/loveauth/utils"
	"net/http"
	"time"
)

var handlers = make(map[string]gin.HandlerFunc)
var getHandlers = make(map[string]gin.HandlerFunc)

type BuyProvider interface {
	Name() string
	GetShopItem(request *RequestPay, resp *ResponsePay, testTime time.Time) *ShopItem
	CheckItemLegality(shopItem *ShopItem) bool
	GeneratePaymentSequence(request *RequestPay, resp *ResponsePay) string
	SaveOrder(request *RequestPay, resp *ResponsePay, testTime time.Time, shopItem *ShopItem) error
}

type BuyProviders map[string]BuyProvider

var usedBuyProviders = BuyProviders{}

type Buy struct {
	BuyProvider BuyProvider
}

func UseBuy(aps ...BuyProvider) {
	for _, BuyProvider := range aps {
		if usedBuyProviders[BuyProvider.Name()] != nil {
			logrus.WithFields(logrus.Fields{
				"BuyProvider": BuyProvider.Name(),
			}).Warn("BuyProvider replaced.")
		}

		usedBuyProviders[BuyProvider.Name()] = BuyProvider
	}
}

type RequestPay struct {
	GlobalId     int64        `form:"globalId" json:"globalId" binding:"required"`
	ShopId       int          `form:"shopId" json:"shopId" binding:"required"`
	Type         int          `form:"type" json:"type"`
	BuyItemMidas BuyItemMidas `form:"buyItemMidas" json:"buyItemMidas"`
	Num          int          `form:"num" json:"num"`
	MaxCostNum   int          `form:"maxCostNum" json:"maxCostNum"`
	ClientIp     string       `form:"client_ip" json:"client_ip"`
	ShopItemInfo ShopItemInfo `form:"shopItemInfo" json:"shopItemInfo"`
}

type ResponsePay struct {
	Body struct {
		Code        int64   `json:"code"`
		Message     string  `json:"message"`
		Sequence    string  `json:"sequence"`
		ProductId   string  `json:"productId"`
		Params      string  `json:"params"`
		CostNum     float64 `json:"costNum"`
		ItemId      int     `json:"itemId"`
		ItemCount   int     `json:"itemCount"`

		WechatPullPay *wechat_pay.PullPay `json:"wechat_pull_pay"`
		QqPullPay     *qpay.PullQPay      `json:"qq_pull_pay"`

		BiliOrder *bilibili.BilibiliOrder `json:"bili_order"`
		VivoOrder *vivo.VivoBuyGoodsResponse `json:"vivo_order"`
	}
}

func (b *Buy) MakeOrder(c *gin.Context) {
	var request RequestPay

	// validation
	if err := c.BindJSON(&request); err != nil {
		utils.QuickReply(c, errors.Failed, "buy BindJSON failed: %v", err)
		return
	}

	if request.ClientIp == "" {
		request.ClientIp = c.ClientIP()
	}

	resp := ResponsePay{}
	testTime := time.Now()
	shopItem := b.BuyProvider.GetShopItem(&request, &resp, testTime)
	if shopItem == nil {
		utils.QuickReply(c, errors.ShopIdMiss, "buy is not valid for this user.")
		return
	}

	if !b.BuyProvider.CheckItemLegality(shopItem) {
		utils.QuickReply(c, errors.Failed, "buy item type error.")
		return
	}

	b.BuyProvider.GeneratePaymentSequence(&request, &resp)

	err := b.BuyProvider.SaveOrder(&request, &resp, testTime, shopItem)
	if err != nil {
		utils.QuickReply(c, errors.Failed, "buy save order failed: %v", err)
		return
	}

	resp.Body.Code = int64(errors.Ok)
	c.JSON(http.StatusOK, resp.Body)
	return
}

func StartBuy(group *gin.RouterGroup) {
	UseBuy(&BuyDiamond{}, &BuyItem{}, &DiamondBuy{})

	for _, b := range usedBuyProviders {
		var BuyProvider = b
		group.POST("/"+b.Name(), func(c *gin.Context) {
			var buy = Buy{BuyProvider}
			buy.MakeOrder(c)
		})
	}
}

func Start(group *gin.RouterGroup) {
	loadPaymentTemplates()

	for key, value := range handlers {
		group.POST(key, value)
	}

	for key, value := range getHandlers {
		group.GET(key, value)
	}

	StartBuy(group)
	StartQuery(group)
}
