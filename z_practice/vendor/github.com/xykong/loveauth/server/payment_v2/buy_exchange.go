package payment_v2

import (
	"github.com/xykong/loveauth/errors"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"github.com/xykong/loveauth/utils"
	"time"
)

type DiamondBuy struct {
	//购买物品
}

func (b *DiamondBuy) Name() string {
	return "diamond_buy"
}

func (b *DiamondBuy) GetShopItem(request *RequestPay, resp *ResponsePay, testTime time.Time) *ShopItem {
	shopItem := getShopItem(request.ShopId, request.GlobalId, request.Type, testTime)
	return shopItem
}

func (b *DiamondBuy) CheckItemLegality(shopItem *ShopItem) bool {
	return !(shopItem.PriceId != model.DiamondId && shopItem.PriceId != model.GoldId)
}

func (b *DiamondBuy) GeneratePaymentSequence(request *RequestPay, resp *ResponsePay) string {
	sequence := utils.GeneratePaymentSequence(request.GlobalId, request.ShopId)
	resp.Body.Sequence = sequence
	return sequence
}

func (b *DiamondBuy) SaveOrder(request *RequestPay, resp *ResponsePay, testTime time.Time, shopItem *ShopItem) error {
	if int(resp.Body.CostNum) > request.MaxCostNum {
		if shopItem.PriceId == model.DiamondId {
			return errors.NewCodeString(errors.DiamondShortage, "diamondBuy diamond not enough.")
		}

		return errors.NewCodeString(errors.GoldShortage, "diamondBuy gold not enough.")
	}

	resp.Body.Message = "diamondBuy successfully!"

	if shopItem.LimitCount != 0 {
		if shopItem.AvailableCount < request.Num {
			return errors.NewCodeString(errors.ShopDiamondBuyNumOut, "diamondBuyNum outside activity.")
		}
	}

	amount := shopItem.PriceValue
	if shopItem.DiscountPriceValue != 0 { //打折价格
		amount = shopItem.DiscountPriceValue
	}

	profile, _ := storage.QueryProfile(request.GlobalId)
	order := model.Order{
		GlobalId:  request.GlobalId,
		ShopId:    request.ShopId,
		Sequence:  resp.Body.Sequence,
		Timestamp: testTime,
		State:     model.OrderStateComplete,
		Type:      model.OrderTypeCommon,
		Num:       request.Num,
		Amount:    int(amount) * request.Num * 100,
	}

	if profile != nil {
		order.Vendor = profile.Vendor
	}

	storage.Insert(storage.PayDatabase(), &order)

	if shopItem != nil {
		resp.Body.ProductId = shopItem.ProductId
		resp.Body.ItemId = shopItem.ItemId
		resp.Body.ItemCount = shopItem.ItemCount * request.Num
		resp.Body.CostNum = shopItem.PriceValue * float64(request.Num)
	}

	if shopItem.DiscountPriceValue != 0 {
		resp.Body.CostNum = shopItem.DiscountPriceValue * float64(request.Num)
	}

	resp.Body.Message = "diamondBuy successfully!"
	return nil
}
