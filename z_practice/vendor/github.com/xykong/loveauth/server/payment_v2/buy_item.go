package payment_v2

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/errors"
	"github.com/xykong/loveauth/server/send_mail"
	"github.com/xykong/loveauth/services/payment/alipay"
	"github.com/xykong/loveauth/services/payment/bilibili"
	"github.com/xykong/loveauth/services/payment/midas"
	"github.com/xykong/loveauth/services/payment/qpay"
	"github.com/xykong/loveauth/services/payment/vivo"
	"github.com/xykong/loveauth/services/payment/wechat_pay"
	"github.com/xykong/loveauth/settings"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"github.com/xykong/loveauth/utils"
	"time"
)

type BuyItem struct {
	//购买物品
}

func (b *BuyItem) Name() string {
	return "buy_item"
}

func (b *BuyItem) GetShopItem(request *RequestPay, resp *ResponsePay, testTime time.Time) *ShopItem {
	shopItem := getShopItem(request.ShopId, request.GlobalId, request.Type, testTime)
	return shopItem
}

func (b *BuyItem) CheckItemLegality(shopItem *ShopItem) bool {
	// 111 is defined by 吉小星 in item_template.xlsx
	return shopItem.ItemId != 111
}

func (b *BuyItem) GeneratePaymentSequence(request *RequestPay, resp *ResponsePay) string {
	sequence := utils.GeneratePaymentSequence(request.GlobalId, request.ShopId)
	resp.Body.Sequence = sequence
	return sequence
}

func (b *BuyItem) SaveOrder(request *RequestPay, resp *ResponsePay, testTime time.Time, shopItem *ShopItem) error {
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
		Type:      model.OrderTypeItem,
		Amount:    int(amount * 100),
		PayMethod: request.Type,
	}

	if profile != nil {
		order.Vendor = profile.Vendor
	}

	switch request.Type {
	case model.Midas:
		midasRequest := midas.DoMidasBuyGoodsReq{
			OpenId:      request.BuyItemMidas.OpenId,
			OpenKey:     request.BuyItemMidas.OpenKey,
			Pf:          request.BuyItemMidas.Pf,
			Pfkey:       request.BuyItemMidas.Pfkey,
			UserIp:      request.BuyItemMidas.UserIp,
			PlatId:      request.BuyItemMidas.PlatId,
			PayItem:     fmt.Sprintf("%v*%v*%v", shopItem.ItemId, shopItem.PriceValue*10, 1),
			GoodsMeta:   request.BuyItemMidas.GoodsMeta,
			GoodsUrl:    "",
			Amount:      int(shopItem.PriceValue * 10),
			MaxNum:      0,
			AppMode:     1,
			AppMetadata: resp.Body.Sequence,
		}

		midasResp, err := midas.BuyGoods(midasRequest)
		if err != nil {
			return errors.New(fmt.Sprintf("buyItem failed: %v", err))
		}

		resp.Body.Params = midasResp.Body.UrlParams

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
			}).Error("BuyItem SaveOrder Wechat Unifiedorder failed.")
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
			"orderTitle":    request.ShopItemInfo.Title,
			"orderDesc":     request.ShopItemInfo.Desc,
			"extInfo":       order.Sequence,
		}
		vivoOrder, err := vivo.BuyGoods(params, key)
		if err != nil {
			return errors.New(fmt.Sprintf("vivo trade err: %v", err))
		}
		resp.Body.Params = vivoOrder.AccessKey
		resp.Body.VivoOrder = vivoOrder
		order.SNSOrderId = vivoOrder.OrderNumber

	case model.Alipay:

		alipayTradeAppPay := alipay.AliPayTradeAppPay{}
		alipayTradeAppPay.NotifyURL = settings.GetString("lovepay", "alipay.callback")
		alipayTradeAppPay.Body = request.ShopItemInfo.Desc
		alipayTradeAppPay.Subject = request.ShopItemInfo.Title
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
		}).Info("BuyItem Alipay")

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

	resp.Body.Message = "buyItem successfully!"
	err := storage.Save(storage.PayDatabase(), &order)

	send_mail.DoOrderWaring(storage.PreGenSequence, order.Sequence, order.Vendor)

	return err
}

//noinspection ALL
type BuyItemMidas struct {
	//
	// 腾讯登录态：
	//用户的id，QQ号码、openid等，String类型。例如：
	//
	//openid = “281348406”(QQ号)
	//
	//openid= “559B3E350A3AC6EB5CA98068AE5BA451”（openid）
	//
	//第三方登录态：
	//搜狗登录态： openid=搜狗帐号体系下的用户ID
	//
	//游客登录态：（外部登录，不校验登录态）openid=由应用自定义，保证其唯一性，最大长度32位
	//
	//不能包含：单引号 '  < > ( )  |  & = * ^等特殊字符
	//
	OpenId string `json:"openid" structs:"openid"`
	//
	// 腾讯登录态：
	//用户的登录态，skey、accessToken等,String类型。例如：
	//
	//openKey=“@8B8cFEpyi” (skey)
	//
	//openKey= “29d8443676b3be073ac56348417cbe65” （pay_token）
	//
	//特别注意如果使用的是手Q登录态，这里填的是支付时专用的pay_token
	//
	//openKey = “29d8443676b3be073ac56348417cbe65”  (accessToken)
	//
	//特别注意如果 使用的是微信登录态，这里填的是登录时获取到的accessToken
	//
	//第三方登录态：
	//搜狗登录态： openkey=搜狗帐号体系下的用户登录Token
	//
	//游客登录态：（外部登录，不校验登录态）
	//
	//openkey=由应用自定义后台不校验，但不能为空
	//
	OpenKey string `json:"openkey" structs:"openkey"`
	//
	// 平台标识信息：平台-注册渠道-版本-安装渠道-业务自定义(自定义)，最大150字节。（自定义部分不能包含单引号 '  < > ( )  |  & = * ^-等特殊字符，支持下划线_）
	//
	//例如：
	//
	//qq_m_qq-2001-android-2011-xxxx
	//
	//qq_m_wx-2001-android-2011-xxxx
	//
	//其中
	//
	//qq_m_qq 表示手Q平台启动，用qq登录态
	//
	//qq_m_wx 表示手Q平台启动，微信登录态
	//
	//渠道信息包括：业务侧自己定义
	//
	//版本信息包括：iap,android，html5
	//
	//平台：目前支持以下平台
	//
	//　　　　微信： wechat
	//　　　　手机QQ： qq_m
	//　　　　手机Qzone： qzone_m
	//　　　　手机QQ游戏大厅： mobile
	//　　　　应用宝： myapp_m
	//　　　　手机QQ浏览器： qqbrowser_m
	//　　　　3366： 3366_m
	//　　　　海外微信-wx帐号 wechat_abroad_wx
	//　　　　海外微信-qq帐号 wechat_abroad_qq
	//　　　　海外微信-pc老用户 wechat_abroad_pc
	//　　　　搜狗 sogou_m
	//　　　　第三方(非手Q非微信非搜狗) desktop_m_guest
	//
	//目前支持的pf有：
	//
	//应用宝：
	//
	//　　　　myapp_m_qq-2001-android-2011-xxxx
	//　　　　myapp_m_wx-2001-android-2011-xxxx
	//
	//手Q：
	//
	//　　　　qq_m_qq-2001-android-2011-xxxx
	//　　　　qq_m_wx-2001-android-2011-xxxx
	//
	//
	//手 Qzone:
	//
	//　　　　qzone_m_qq-2001-android-2011-xxx
	//　　　　qzone_m_wx-2001-android-2011-xxx
	//
	//微信：
	//
	//　　　　wechat_wx-2001-android-2011-xxxx
	//　　　　wechat_qq-2001-android-2011-xxxx
	//
	//手Q游戏大厅：
	//
	//　　　　moblie_wx-2001-android-2011-xxxx
	//　　　　mobile_qq-2001-android-2011-xxxx
	//
	//桌面启动：
	//
	//　　　　desktop_m_qq-2001-android-2011-xxxx
	//　　　　desktop_m_wx-2001-android-2011-xxxx
	//
	//手机QQ浏览器：
	//
	//　　　　qqbrowser_m_qq-2001-android-2011-xxxx
	//　　　　qqbrowser_m_wx-2001-android-2011-xxxx
	//
	//搜狗游戏：
	//
	//　　　　sougo_m-2001-android-2011-xxxx
	//
	//第三方：
	//
	//　　　　desktop_m_guest-2001-android-2011-xxxx
	//
	//注：使用了MSDK的游戏可以通过MSDK的接口：WGGetPf()获取pf的数值。
	//
	Pf string `json:"pf" structs:"pf"`
	//
	// String类型参数，由平台来源和openkey根据规则生成的一个密钥串，跳转到应用首页后，URL后会带该参数。平台直接传给应用，应用原样传给平台即可。
	//
	//自研和第三方登录的应用不校验，可以传递为pfKey="pfKey"
	//
	//非自研强校验,pfKey="58FCB2258B0BF818008382BD025E8022"（来自平台）
	//
	Pfkey string `json:"pfkey" structs:"pfkey"`
	//
	// （可选）用户的外网IP
	//
	UserIp string `json:"userip" structs:"userip"`
	//
	// （平台：IOS（0），安卓（1）
	//
	PlatId int `json:"platId" structs:"platId"`
	//
	// 物品信息，格式必须是“name*des”，批量购买套餐时也只能有1个道具名称和1个描述，即给出该套餐的名称和描述。
	// name表示物品的名称，des表示物品的描述信 息。用户购买物品的确认支付页面，将显示该物品信息。
	// 长度必须<=256字符，必须  使用utf8编码。目前goodsmeta超过76个字符后不能添加回车字符等特殊字符.
	//
	GoodsMeta string `json:"goodsmeta" structs:"goodsmeta"`
}

type ShopItemInfo struct {
	Title string `json:"title"`
	Desc  string `json:"desc"`
}
