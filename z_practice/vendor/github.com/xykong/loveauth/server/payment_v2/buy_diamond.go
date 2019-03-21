package payment_v2

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/errors"
	"github.com/xykong/loveauth/server/send_mail"
	"github.com/xykong/loveauth/services/payment/alipay"
	"github.com/xykong/loveauth/services/payment/bilibili"
	"github.com/xykong/loveauth/services/payment/qpay"
	"github.com/xykong/loveauth/services/payment/vivo"
	"github.com/xykong/loveauth/services/payment/wechat_pay"
	"github.com/xykong/loveauth/settings"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"github.com/xykong/loveauth/utils"
	"time"
)

type BuyDiamond struct {
	//购买钻石
}

func (b *BuyDiamond) Name() string {
	return "buy_diamond"
}

func (b *BuyDiamond) GetShopItem(request *RequestPay, resp *ResponsePay, testTime time.Time) *ShopItem {
	shopItem := getShopItem(request.ShopId, request.GlobalId, request.Type, testTime)
	return shopItem
}

func (b *BuyDiamond) CheckItemLegality(shopItem *ShopItem) bool {
	// 111 is defined by 吉小星 in item_template.xlsx
	//return shopItem.ItemId == 111
	return true
}

func (b *BuyDiamond) GeneratePaymentSequence(request *RequestPay, resp *ResponsePay) string {
	sequence := utils.GeneratePaymentSequence(request.GlobalId, request.ShopId)
	resp.Body.Sequence = sequence
	return sequence
}

func (b *BuyDiamond) SaveOrder(request *RequestPay, resp *ResponsePay, testTime time.Time, shopItem *ShopItem) error {
	amount := shopItem.PriceValue
	if shopItem.DiscountPriceValue != 0 { //打折价格
		amount = shopItem.DiscountPriceValue
	}

	resp.Body.CostNum = amount

	profile, _ := storage.QueryProfile(request.GlobalId)
	order := model.Order{
		GlobalId:  request.GlobalId,
		ShopId:    request.ShopId,
		Sequence:  resp.Body.Sequence,
		Timestamp: testTime,
		State:     model.OrderStatePrepare,
		Type:      model.OrderTypeDiamond,
		Amount:    int(amount * 100),
		PayMethod: request.Type,
	}

	if profile != nil {
		order.Vendor = profile.Vendor
	}

	switch request.Type {
	case model.Midas:
	case model.QUICK, model.QUICK_XIAOMI, model.QUICK_MEIZU, model.QUICK_UC, model.QUICK_OPPO, model.QUICK_M4399, model.QUICK_YSDK, model.QUICK_IQIYI_PPS:
		resp.Body.ProductId = fmt.Sprintf("loveworld.quick.%d.%d", request.Type, shopItem.ShopId)

	case model.AppStore, model.HUAWEI, model.QUICK_KUAIKAN:
		resp.Body.ProductId = shopItem.ProductId

	case model.Wechat:
		wxOrder, err := wechat_pay.Unifiedorder(request.ClientIp, &order)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err":   err,
				"order": order,
			}).Error("BuyDiamond SaveOrder Wechat Unifiedorder failed.")
			return errors.NewCodeString(errors.WxUnifiedOrderFailed, "微信支付下单失败")
		}

		resp.Body.WechatPullPay = wxOrder

	case model.QPay:
		qpayOrder, err := qpay.Unifiedorder(request.ClientIp, &order)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err":   err,
				"order": order,
			}).Error("BuyItem SaveOrder QPay Unifiedorder failed.")
			return errors.NewCodeString(errors.QPayUnifiedOrderFailed, "qq支付下单失败")
		}

		resp.Body.QqPullPay = qpayOrder

	case model.Vivo:

		cpId := settings.GetString("lovepay", "vivo.cpId")
		appId := settings.GetString("lovepay", "vivo.appId")
		key := settings.GetString("lovepay", "vivo.key")
		callback := settings.GetString("lovepay", "vivo.callback")

		params := map[string]interface{}{
			"cpId":          cpId,
			"appId":         appId,
			"cpOrderNumber": fmt.Sprintf("%v-%v", order.GlobalId, order.Timestamp.Unix()),
			"notifyUrl":     callback,
			"orderAmount":   fmt.Sprintf("%d", order.Amount),
			"orderTitle":    fmt.Sprintf("%d钻石", shopItem.ItemCount),
			"orderDesc":     shopItem.GiftDescription,
			"extInfo":       order.Sequence,
		}

		vivoOrder, err := vivo.BuyGoods(params, key)
		if err != nil {
			return errors.New(fmt.Sprintf("vivo trade err: %v", err))
		}
		resp.Body.Params = vivoOrder.AccessKey
		order.SNSOrderId = vivoOrder.OrderNumber
		resp.Body.VivoOrder = vivoOrder

	case model.Alipay:

		alipayTradeAppPay := alipay.AliPayTradeAppPay{}
		alipayTradeAppPay.NotifyURL = settings.GetString("lovepay", "alipay.callback")
		alipayTradeAppPay.Body = shopItem.GiftDescription
		alipayTradeAppPay.Subject = fmt.Sprintf("%d钻石", shopItem.ItemCount)
		alipayTradeAppPay.OutTradeNo = fmt.Sprintf("%v-%v", order.GlobalId, order.Timestamp.Unix())
		alipayTradeAppPay.TotalAmount = fmt.Sprintf("%0.2f", amount)
		alipayTradeAppPay.ProductCode = shopItem.ProductId
		alipayTradeAppPay.PassbackParams = order.Sequence

		result, err := alipay.TradeAppPay(alipayTradeAppPay)
		if err != nil {

			return errors.New("buyItem failed:" + err.Error())
		}

		logrus.WithFields(logrus.Fields{
			"result": result,
		}).Info("BuyDiamond Alipay")

		resp.Body.Params = result

	case model.BILIBILI:

		resp.Body.BiliOrder = bilibili.MakeOrderSign(&order)

	case model.Douyin:

		resp.Body.CostNum = float64(order.Amount)

	case model.Mgtv:

		resp.Body.CostNum = float64(order.Amount)
		resp.Body.ProductId = shopItem.ProductId

	case model.OFFICIAL:

	default:
		return errors.NewCodeString(errors.ShopDataException, "paymethod not match.")
	}

	err := storage.Save(storage.PayDatabase(), &order)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":   err,
			"order": order,
		}).Error("BuyDiamond SaveOrder failed.")
	}

	// 人民币购买钻石
	send_mail.DoOrderWaring(storage.PreGenSequence, order.Sequence, order.Vendor)

	resp.Body.Message = "buyDiamond successfully!"
	resp.Body.ItemCount = shopItem.ItemCount

	return err
}
