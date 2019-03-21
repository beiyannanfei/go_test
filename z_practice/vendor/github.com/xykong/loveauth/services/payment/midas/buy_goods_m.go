package midas

import (
	"github.com/sirupsen/logrus"
	"encoding/json"
	"github.com/fatih/structs"
	"github.com/xykong/loveauth/errors"
)

//
// 请求: 扣除游戏币接口
// Binding from JSON
// noinspection ALL
type DoMidasBuyGoodsReq struct {
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
	// （请使用x*p*num的格式，x表示物品ID，p表示单价（以Q点为单位，1Q币=10Q点，单价的制定需遵循腾讯定价规范），
	// num表示默认的购买数量。（格式：物品ID1*单价1*    建议数量1，批量购买物品时使用;分隔，
	// 如：id1*price1*num1;id2*price2*num2)长度必须<512
	// 注：道具直购物品信息在这里游戏自己定义自己管理
	//
	PayItem string `json:"payitem" structs:"payitem"`
	//
	// 物品信息，格式必须是“name*des”，批量购买套餐时也只能有1个道具名称和1个描述，即给出该套餐的名称和描述。
	// name表示物品的名称，des表示物品的描述信 息。用户购买物品的确认支付页面，将显示该物品信息。
	// 长度必须<=256字符，必须  使用utf8编码。目前goodsmeta超过76个字符后不能添加回车字符等特殊字符.
	//
	GoodsMeta string `json:"goodsmeta" structs:"goodsmeta"`
	//
	// 物品的图片url(长度<512字符) 注：参数已废弃直接传空
	//
	GoodsUrl string `json:"goodsurl" structs:"goodsurl"`
	//
	// (可选)道具总价格。（amt必须等于所有物品：单价*建议数量的总和 单位为1Q点）
	//
	Amount int `json:"amt" structs:"amt"`
	//
	// (可选) 用户可购买的道具数量的最大值。仅当appmode的值为2时，可以输入该参数。
	// 输入的值需大于参数“payitem”中的num，如果小于num，则自动调整为num的值。
	//
	MaxNum int `json:"max_num" structs:"max_num"`
	//
	// (可选)1表示用户不可以修改物品数量，2 表示用户可以选择购买物品的数量。默  认2（注：批量购买的时候，必须等于1）
	//
	AppMode int `json:"appmode" structs:"appmode"`
	//
	// （可选）发货时透传给应用。长度必须<=128字符
	//
	AppMetadata string `json:"app_metadata" structs:"app_metadata"`
}

//
// in: body
// swagger:parameters profile_query_vip
type DoMidasBuyGoodsReqBodyParams struct {
	//
	// swagger:allOf
	// in: body
	DoMidasBuyGoodsReq DoMidasBuyGoodsReq
}

//
// 应答: 扣除游戏币接口
// swagger:response DoMidasBuyGoodsRsp
// noinspection ALL
type DoMidasBuyGoodsRsp struct {
	// in: body
	Body struct {
		//
		// 返回码 0 ：成功， >=1000：失败
		//
		Ret int `json:"ret"`
		//
		// ret不为 0 的时候，错误信息（utf-8编码）
		//
		Msg string `json:"msg"`
		//
		// ret为0的时候，开发者需要保留。后续扣费成功后调用第三方发货时，会再传给开发者，作为本次交易的标识，有效期5分钟
		//
		Token string `json:"token"`
		//
		// ret为0的时候，返回真正购买物品的url的参数，开发者需要把该参数
		// 传给sdk跳转到相关页面使用户完成真正的购买动作。
		//
		UrlParams string `json:"url_params"`
	}
}

func BuyGoods(request DoMidasBuyGoodsReq) (*DoMidasBuyGoodsRsp, error) {

	params := structs.Map(&request)

	method := "GET"
	//noinspection SpellCheckingInspection
	urlPath := "/mpay/buy_goods_m"

	body, err := SendRequest(method, urlPath, params, request.PlatId)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
			"url":   urlPath,
		})
		return nil, errors.NewCodeError(errors.Failed, err)
	}

	var data DoMidasBuyGoodsRsp
	if err := json.Unmarshal(body, &data.Body); err != nil {
		return nil, errors.NewCodeError(errors.Failed, err)
	}

	// token校验失败(18)
	if data.Body.Ret == 1018 {
		return nil, errors.NewCodeString(errors.AuthTokenInvalid, data.Body.Msg)
	}

	if data.Body.Ret != 0 {
		return nil, errors.NewCodeString(errors.Failed, data.Body.Msg)
	}

	return &data, nil
}
