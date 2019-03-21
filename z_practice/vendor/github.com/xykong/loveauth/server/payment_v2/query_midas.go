package payment_v2

import (
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/errors"
	"github.com/xykong/loveauth/services/payment/midas"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"github.com/xykong/loveauth/utils"
	"time"
)

type QueryMidasOrder struct {
}

func (q *QueryMidasOrder) Name() string {
	return "query_midas"
}

func (q *QueryMidasOrder) QueryOrder(sequence string, globalId int64) (*model.Order, error) {
	var order *model.Order
	if len(sequence) > 0 {
		order = storage.QueryOrderPlacedWithSequence(sequence)
	} else {
		order = storage.QueryOrderPlaced(globalId)
	}

	if order == nil || order.GlobalId != globalId {
		return nil, errors.NewCodeString(errors.ShopOrderMiss, "query midas QueryOrderPlacedWithSequence failed.")
	}

	if order.State == model.OrderStateComplete {
		return nil, errors.NewCodeString(errors.ShopQueryGeted, "query QueryOrderPlacedWithSequence failed.")
	}

	if order.State == model.OrderStateFailed {
		return nil, errors.NewCodeString(errors.QuickPayFailed, "query QueryOrder pay failed.")
	}

	return order, nil
}

func (q *QueryMidasOrder) GetShopItem(shopId int, globalId int64, orderType int, timestamp time.Time) *ShopItem {
	shopItem := getShopItem(shopId, globalId, orderType, timestamp)
	if shopItem == nil || shopItem.ItemCount == 0 {
		// something wrong. user have balance without place order. may be shopItem removed.
		logrus.WithFields(logrus.Fields{
			"shopId":    shopId,
			"globalId":  globalId,
			"orderType": orderType,
		}).Warn("query_midas: shopItem without.")
	}

	return shopItem
}

func (q *QueryMidasOrder) DealOrder(request *RequestQueryOrder, order *model.Order, shopItem *ShopItem) (*AcquiredItem, error) {
	var acquiredItem *AcquiredItem
	if order.Type == model.OrderTypeDiamond {
		// query balance
		midasRequest := midas.DoMidasGetBalanceReq{
			OpenId:  request.QueryMidas.OpenId,
			OpenKey: request.QueryMidas.OpenKey,
			Pf:      request.QueryMidas.Pf,
			Pfkey:   request.QueryMidas.Pfkey,
			UserIp:  request.QueryMidas.UserIp,
			PlatId:  request.QueryMidas.PlatId,
		}

		midasResp, err := midas.GetBalance(midasRequest)
		if err != nil {
			return nil, errors.NewCodeString(errors.ShopNotPay, "query GetBalance failed: %v", err)
		}

		if midasResp.Body.Balance == 0 {
			return nil, errors.NewCodeString(errors.ShopNotPay, "query GetBalance = 0")
		}

		acquiredItem, err = acquireBalance(order, shopItem, request, midasResp)
		if err != nil {
			return nil, errors.NewCodeString(errors.Failed, "query midas.Pay failed: %v", err)
		}
	} else if order.Type == model.OrderTypeItem {

	}

	return acquiredItem, nil
}

func (q *QueryMidasOrder) SaveOrder(order *model.Order) {
	order.State = model.OrderStateComplete
	storage.Save(storage.PayDatabase(), order)
}

func acquireBalance(order *model.Order,
	shopItem *ShopItem,
	request *RequestQueryOrder,
	doMidasGetBalanceRsp *midas.DoMidasGetBalanceRsp) (*AcquiredItem, error) {

	consume := 0
	sequence := ""
	if order != nil && shopItem != nil && doMidasGetBalanceRsp.Body.Balance >= shopItem.ItemCount {
		consume = shopItem.ItemCount
		sequence = order.Sequence
	} else {
		consume = doMidasGetBalanceRsp.Body.Balance
		sequence = utils.GeneratePaymentSequence(request.GlobalId, 0)
	}

	_, err := midas.Pay(midas.DoMidasPayReq{
		OpenId:    request.QueryMidas.OpenId,
		OpenKey:   request.QueryMidas.OpenKey,
		Pf:        request.QueryMidas.Pf,
		Pfkey:     request.QueryMidas.Pfkey,
		PlatId:    request.QueryMidas.PlatId,
		Amount:    consume,
		BillNo:    sequence,
		UserIp:    request.QueryMidas.UserIp,
		PayItem:   "",
		AppRemark: sequence,
	})
	if err != nil {
		return nil, err
	}

	acquiredItem := AcquiredItem{
		ShopId:         shopItem.ShopId,
		ItemId:         shopItem.ItemId,
		ItemCount:      consume,
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

	return &acquiredItem, nil
}

//noinspection ALL
type QueryMidas struct {
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
}
