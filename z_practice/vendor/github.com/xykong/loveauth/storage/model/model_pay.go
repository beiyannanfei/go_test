package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

// 111 is defined by 吉小星 in item_template.xlsx
const (
	DiamondId = 111
	GoldId    = 211
)

type Order struct {
	gorm.Model

	GlobalId   int64 `gorm:"index"`
	ShopId     int
	Sequence   string `gorm:"type:varchar(191); index"`
	Timestamp  time.Time
	State      OrderState
	Type       OrderType
	Num        int `gorm:"default:'1'"`
	Amount     int
	SNSOrderId string `gorm:"type:varchar(191); index"`
	PayMethod  int
	Vendor     Vendor
}

type OrderState int

const (
	OrderStatePrepare OrderState = iota
	OrderStatePlace
	OrderStateComplete
	OrderStateFailed
	OrderStateRepaired // 补单
)

type OrderType int

const (
	OrderTypeInvalid OrderType = iota
	OrderTypeDiamond
	OrderTypeItem
	OrderTypeCommon
)

type PayMethod int

const (
	UnDefine = iota
	Midas
	AppStore
	GooglePlay
	Wechat
	Alipay
	BILIBILI
	QUICK
	Vivo
	HUAWEI
	QPay
	QUICK_UC
	QUICK_OPPO
	QUICK_M4399
	QUICK_YSDK
	QUICK_IQIYI_PPS
	QUICK_MEIZU
	QUICK_KUAIKAN
	OFFICIAL
	QUICK_XIAOMI
	Douyin
	Mgtv
)

type Inventory struct {
	gorm.Model

	GlobalId int64
	OpenId   string
	ZoneId   int64

	ItemId        int64
	ItemCount     int64
	ItemHistory   int64
	ItemPurchased int64
}

type Activity struct {
	gorm.Model

	GlobalId int64
	OpenId   string
	ZoneId   int64

	ActivityId      int64
	ActivityGroupId int64
	ActivityCount   int64
	ActivityOccur   time.Time
}

type ActivityType int

const (
	ActivityTypeInvalid ActivityType = iota
	ActivityTypeAdditional
	ActivityTypeDiscount
	ActivityTypeFirstJoin
	ActivityTypeMaxCount
)

type ShopTag int

const (
	ShopTagInvalid ShopTag = iota
	ShopTagRecharge
	ShopTagHeadFrame
	ShopTagGift
)

type ActivityInfo struct {
	ActivityId      int64
	ActivityGroupId int64
	ShopId          int
	PriceId         int64
	PriceValue      int64
	Description     string
	Start           time.Time
	Finish          time.Time
	Type            ActivityType
	Limit           int64
	AccountCreate   time.Time
	Dependencies    []int64
}

type ShopInfo struct {
	ShopId     int
	ItemId     int64
	ItemCount  int64
	PriceId    int64
	PriceValue int64
	Tag        ShopTag
}

//
// 请求: [手Q]会员详情
// Binding from JSON
// noinspection ALL
type DoMidasCallbackReq struct {
	gorm.Model `json:"-"`

	//
	// 与APP通信的用户key，跳转到应用首页后，URL后会带该参数。由平台直接传给应用，应用原样传给平台即可。
	// 根据APPID以及QQ号码生成，即不同的appid下，同一个QQ号生成的OpenID是不一样的。
	//
	OpenId string `json:"openid"`
	//
	// 应用的唯一ID。可以通过appid查找APP基本信息。
	//
	AppId string `json:"appid"`
	//
	// linux时间戳。注意开发者的机器时间与计费服务器的时间相差不能超过15分钟。
	//
	TimeStamp          string `json:"ts"`
	PayItem            string `json:"payitem"`
	Token              string `json:"token"`
	BillNo             string `json:"billno"`
	Version            string `json:"version"`
	ZoneId             string `json:"zoneid"`
	ProvideType        string `json:"providetype"`
	Amount             string `json:"amt"`
	Appmeta            string `json:"appmeta"`
	CftId              string `json:"cftid"`
	ChannelId          string `json:"channel_id"`
	Clientver          string `json:"clientver"`
	PayamtCoins        string `json:"payamt_coins"`
	PubacctPayamtCoins string `json:"pubacct_payamt_coins"`
	Bazinga            string `json:"bazinga"`
	Sig                string `json:"sig"`
}

type QuickPayCallback struct {
	gorm.Model   `xml:"-" json:"-"`
	IsTest       string `xml:"is_test" json:"is_test"`
	Channel      string `xml:"channel" json:"channel"`
	ChannelUid   string `xml:"channel_uid" json:"channel_uid"`
	GameOrder    string `xml:"game_order" json:"game_order" gorm:"type:varchar(191); index"`
	OrderNo      string `xml:"order_no" json:"order_no"`
	PayTime      string `xml:"pay_time" json:"pay_time"`
	Amount       string `xml:"amount" json:"amount"`
	Status       string `xml:"status" json:"status"`
	ExtrasParams string `xml:"extras_params" json:"extras_params"`
}

type WechatPayCallback struct {
	gorm.Model    `xml:"-" json:"-"`
	ReturnCode    string `xml:"return_code" json:"return_code"`
	ReturnMsg     string `xml:"return_msg,omitempty" json:"return_msg,omitempty"`
	Appid         string `xml:"appid" json:"appid"`
	MchId         string `xml:"mch_id" json:"mch_id"`
	DeviceInfo    string `xml:"device_info" json:"device_info,omitempty"`
	NonceStr      string `xml:"nonce_str" json:"nonce_str"`
	Sign          string `xml:"sign" json:"-"`
	ResultCode    string `xml:"result_code" json:"result_code"`
	ErrCode       string `xml:"err_code" json:"err_code,omitempty"`
	ErrCodeDes    string `xml:"err_code_des" json:"err_code_des,omitempty"`
	Openid        string `xml:"openid" json:"openid"`
	IsSubscribe   string `xml:"is_subscribe" json:"is_subscribe"`
	TradeType     string `xml:"trade_type" json:"trade_type"`
	BankType      string `xml:"bank_type" json:"bank_type"`
	TotalFee      int    `xml:"total_fee" json:"total_fee,string"`
	FeeType       string `xml:"fee_type" json:"fee_type,omitempty"`
	CashFee       int    `xml:"cash_fee" json:"cash_fee,string"`
	CashFeeType   string `xml:"cash_fee_type" json:"cash_fee_type,omitempty"`
	TransactionId string `xml:"transaction_id" json:"transaction_id" gorm:"type:varchar(191); index"`
	OutTradeNo    string `xml:"out_trade_no" json:"out_trade_no"`
	Attach        string `xml:"attach" json:"attach,omitempty" gorm:"type:varchar(191); index"`
	TimeEnd       string `xml:"time_end" json:"time_end"`
}

type QPayCallback struct {
	gorm.Model    `xml:"-" json:"-"`
	Appid         string `xml:"appid" json:"appid,omitempty"`
	MchId         string `xml:"mch_id" json:"mch_id"`
	NonceStr      string `xml:"nonce_str" json:"nonce_str"`
	Sign          string `xml:"sign" json:"-"`
	DeviceInfo    string `xml:"device_info" json:"device_info,omitempty"`
	TradeType     string `xml:"trade_type" json:"trade_type"`
	TradeState    string `xml:"trade_state" json:"trade_state"`
	BankType      string `xml:"bank_type" json:"bank_type"`
	FeeType       string `xml:"fee_type" json:"fee_type"`
	TotalFee      int    `xml:"total_fee" json:"total_fee,string"`
	CashFee       int    `xml:"cash_fee" json:"cash_fee,string"`
	CouponFee     int    `xml:"coupon_fee" json:"coupon_fee,string,omitempty"`
	TransactionId string `xml:"transaction_id" json:"transaction_id" gorm:"type:varchar(191); index"`
	OutTradeNo    string `xml:"out_trade_no" json:"out_trade_no"`
	Attach        string `xml:"attach" json:"attach,omitempty" gorm:"type:varchar(191); index"`
	TimeEnd       string `xml:"time_end" json:"time_end"`
	Openid        string `xml:"openid" json:"openid,omitempty"`
}

type HuaweiPayCallback struct {
	gorm.Model  `json:"-"`
	Result      string `json:"result"`                                      // 支付结果 0=支付成功,1=退款成功
	UserName    string `json:"userName"`                                    // 商户名称，开发者注册的公司名
	ProductName string `json:"productName"`                                 // 商品名称
	PayType     int    `json:"payType,string"`                              // 支付类型
	Amount      string `json:"amount"`                                      // 商品支付金额，保留两位小数，退款通知下为退款金额，目前只支持全额退款
	Currency    string `json:"currency"`                                    // 国标货币，值为空或CNY时本参数不出现
	OrderId     string `json:"orderId" gorm:"type:varchar(191); index"`     // 华为支付平台订单号
	NotifyTime  string `json:"notifyTime"`                                  // 通知时间，毫秒
	RequestId   string `json:"requestId"`                                   // 商户生成的唯一订单号，有字母和数字组成
	BankId      string `json:"bankId"`                                      //银行编码
	OrderTime   string `json:"orderTime"`                                   // 下单时间 yyyy-MM-dd hh:mm:ss  仅在urlver为2时有效
	TradeTime   string `json:"tradeTime"`                                   // 交易/退款时间 yyyy-MM-dd hh:mm:ss 仅在urlver为2时有效
	AccessMode  string `json:"accessMode"`                                  // 接入方式 0=yidong  1=pc-web 2=mobile-web 3=TV 仅在urlvere为2时有效
	Spending    string `json:"spending"`                                    // 渠道开销，保留两位小数，单位元 仅在urlver为2时有效
	ExtReserved string `json:"extReserved" gorm:"type:varchar(191); index"` //商户侧保留信息
	SysReserved string `json:"sysReserved"`                                 //商户侧保留信息
	SignType    string `json:"signType"`                                    // 签名类型，默认值为RSA256 表示使用SHA256WithRSA算法
	Sign        string `json:"sign"`                                        // rsa签名，在"utf-8" urlencode后发送
}

// https://doc.open.alipay.com/docs/doc.htm?spm=a219a.7629140.0.0.8AmJwg&treeId=203&articleId=105286&docType=1
// 封号请求
// Binding from JSON
type DoAlipayCallbackReq struct {
	gorm.Model `json:"-"`
	AuthAppId  string `json:"auth_app_id"`
	//通知的发送时间。格式为yyyy-MM-dd HH:mm:ss
	//
	//2015-14-27 15:45:58
	NotifyTime string `json:"notify_time" form:"notify_time"`
	//通知的类型
	//
	//trade_status_sync
	NotifyType string `json:"notify_type" form:"notify_type"`
	//通知校验ID
	//
	//ac05099524730693a8b330c5ecf72da9786
	NotifyId string `json:"notify_id" form:"notify_id"`
	//支付宝分配给开发者的应用Id
	//
	//2014072300007148
	AppId string `json:"app_id" form:"app_id"`
	//编码格式，如utf-8、gbk、gb2312等
	//
	//utf-8
	Charset string `json:"charset" form:"charset"`
	//调用的接口版本，固定为：1.0
	//
	//1.0
	Version string `json:"version" form:"version"`
	//商户生成签名字符串所使用的签名算法类型，目前支持RSA2和RSA，推荐使用RSA2
	//
	//RSA2
	SignType string `json:"sign_type" form:"sign_type"`
	//请参考异步返回结果的验签
	//
	//601510b7970e52cc63db0f44997cf70e
	Sign string `json:"sign" form:"sign" gorm:"size:512"`
	//支付宝交易凭证号
	//
	//2013112011001004330000121536
	TradeNo string `json:"trade_no" form:"trade_no"`
	//原支付请求的商户订单号
	//
	//6823789339978248
	OutTradeNo string `json:"out_trade_no" from:"out_trade_no"`
	//商户业务ID，主要是退款通知中返回退款申请的流水号
	//
	//HZRF001
	OutBizNo string `json:"out_biz_no" form:"out_biz_no"`
	//买家支付宝账号对应的支付宝唯一用户号。以2088开头的纯16位数字
	//
	//2088102122524333
	BuyerId string `json:"buyer_id" form:"buyer_id"`
	//买家支付宝账号
	//
	//15901825620
	BuyerLogonId string `json:"buyer_logon_id" form:"buyer_logon_id"`
	//卖家支付宝用户号
	//
	//2088101106499364
	SellerId string `json:"seller_id" form:"seller_id"`
	//卖家支付宝账号
	//
	//zhuzhanghu@alitest.com
	SellerEmail string `json:"seller_email" form:"seller_email"`
	//交易目前所处的状态，见交易状态说明
	//
	//TRADE_CLOSED
	TradeStatus string `json:"trade_status" form:"trade_status"`
	//本次交易支付的订单金额，单位为人民币（元）
	//
	//20
	TotalAmount string `json:"total_amount" form:"total_amount"`
	//商家在交易中实际收到的款项，单位为元
	//
	//15
	ReceiptAmount string `json:"receipt_amount" form:"receipt_amount"`
	//用户在交易中支付的可开发票的金额
	//
	//10.00
	InvoiceAmount string `json:"invoice_amount" form:"invoice_amount"`
	//用户在交易中支付的金额
	//
	//13.88
	BuyerPayAmount string `json:"buyer_pay_amount" form:"buyer_pay_amount"`
	//使用集分宝支付的金额
	//
	//12.00
	PointAmount string `json:"point_amount" form:"point_amount"`
	//退款通知中，返回总退款金额，单位为元，支持两位小数
	//
	//2.58
	RefundFee string `json:"refund_fee" form:"refund_fee"`
	//商品的标题/交易标题/订单标题/订单关键字等，是请求时对应的参数，原样通知回来
	//
	//当面付交易
	Subject string `json:"subject" form:"subject"`
	//该订单的备注、描述、明细等。对应请求时的body参数，原样通知回来
	//
	//当面付交易内容
	Body string `json:"body" form:"body"`
	//该笔交易创建的时间。格式为yyyy-MM-dd HH:mm:ss
	//
	//2015-04-27 15:45:57
	GmtCreate string `json:"gmt_create" form:"gmt_create"`
	//该笔交易的买家付款时间。格式为yyyy-MM-dd HH:mm:ss
	//
	//2015-04-27 15:45:57
	GmtPayment string `json:"gmt_payment" form:"gmt_payment"`
	//该笔交易的退款时间。格式为yyyy-MM-dd HH:mm:ss.S
	//
	//2015-04-28 15:45:57.320
	GmtRefund string `json:"gmt_refund" form:"gmt_refund"`
	//该笔交易结束时间。格式为yyyy-MM-dd HH:mm:ss
	//
	//2015-04-29 15:45:57
	GmtClose string `json:"gmt_close" form:"gmt_close"`
	//支付成功的各个渠道金额信息，详见资金明细信息说明
	//
	//[{“amount”:“15.00”,“fundChannel”:“ALIPAYACCOUNT”}]
	FundBillList string `json:"fund_bill_list" form:"fund_bill_list"`
	//公共回传参数，如果请求时传递了该参数，则返回给商户时会在异步通知时将该参数原样返回。本参数必须进行UrlEncode之后才可以发送给支付宝
	//
	//merchantBizType%3d3C%26merchantBizNo%3d2016010101111
	PassbackParams string `json:"passback_params" form:"passback_params"`
	//本交易支付时所使用的所有优惠券信息，详见优惠券信息说明
	//
	//[{“amount”:“0.20”,“merchantContribute”:“0.00”,“name”:“一键创建券模板的券名称”,“otherContribute”:“0.20”,“type”:“ALIPAY_DISCOUNT_VOUCHER”,“memo”:“学生卡8折优惠”]
	VoucherDetailList string `json:"voucher_detail_list" form:"voucher_detail_list"`
}

type VivoQueryOrderResponse struct {
	gorm.Model `json:"-"`
	//响应码
	//
	//200
	Ret string `json:"respCode" form:"respCode"`
	//响应信息
	//
	//交易完成
	Message string `json:"respMsg" form:"respMsg"`
	//签名方法
	//
	//对关键信息进行签名的算法名称：MD5
	SignMethod string `json:"signMethod" form:"signMethod"`
	//签名信息
	//
	//对关键信息签名后得到的字符串1，用于商户验签签名规则请参考签名计算说明
	Signature string `json:"signature" form:"signature"`
	//交易种类
	//
	//目前固定01
	TradeType string `json:"tradeType" form:"tradeType"`
	//交易状态
	//
	//0000，代表支付成功
	TradeStatus string `json:"tradeStatus" form:"tradeStatus"`
	//Cp-id
	//
	//定长20位数字，由vivo分发的唯一识别码
	CpId string `json:"cpId" form:"cpId"`
	//appId
	//
	//应用ID
	AppId string `json:"appId" form:"appId"`
	//uid
	//
	//用户在vivo这边的唯一标识
	Uid string `json:"uid" form:"uid"`
	//商户自定义的订单号
	//
	//商户自定义，最长 64 位字母、数字和下划线组成
	CpOrderNumber string `json:"cpOrderNumber" form:"cpOrderNumber"`
	//交易流水号
	//
	//vivo订单号
	OrderNumber string `json:"orderNumber" form:"orderNumber"`
	//交易金额
	//
	//单位：分，币种：人民币，为长整型，如：101，10000
	OrderAmount string `json:"orderAmount" form:"orderAmount"`
	//商户透传参数
	//
	//64位
	ExtInfo string `json:"extInfo" form:"extInfo"`
	//交易时间
	//
	//yyyyMMddHHmmss
	PayTime string `json:"payTime" form:"payTime"`
}

type BilibiliPayCallback struct {
	gorm.Model    `json:"-"`
	BiliId        string `form:"id" json:"id"`                                             // 订单Id
	OrderNo       string `form:"order_no" json:"order_no" gorm:"type:varchar(191); index"` // 哔哩哔哩游戏SDK服务器方订单号
	OutTradeNo    string `form:"out_trade_no" json:"out_trade_no"`                         // 游戏CP厂商支付订单号
	Uid           string `form:"uid" json:"uid"`                                           // 用户ID
	UserName      string `form:"username" json:"username"`                                 // 用户名
	Role          string `form:"role" json:"role"`                                         // 角色名
	Money         string `form:"money" json:"money"`                                       // 支付金额（单位：分）
	PayMoney      string `form:"pay_money" json:"pay_money"`                               // 实际支付金额，单位：分
	GameMoney     string `form:"game_money" json:"game_money"`                             // 应用内货币
	MerchantId    string `form:"merchant_id" json:"merchant_id"`                           // 商户ID
	GameId        string `form:"game_id" json:"game_id"`                                   // 游戏ID
	ZoneId        string `form:"zone_id" json:"zone_id"`                                   // 区服ID
	ProductName   string `form:"product_name" json:"product_name"`                         // 商品名称
	ProductDesc   string `form:"product_desc" json:"product_desc"`                         // 商品描述
	PayTime       string `form:"pay_time" json:"pay_time"`                                 // 订单支付时间
	ClientIP      string `form:"client_ip" json:"client_ip"`                               // 客户端IP
	ExtensionInfo string `form:"extension_info" json:"extension_info"`                     // 额外信息，原样通知回来
	OrderStatus   int    `form:"order_status" json:"order_status"`                         // 订单状态：1为已完成
	Sign          string `form:"sign" json:"sign"`                                         // md5加密后的签名
}

type DouyinCallBackRequest struct {
	gorm.Model  `json:"-"`
	NotifyId    string `json:"notify_id" form:"notify_id"`
	NotifyType  string `json:"notify_type" form:"notify_type"`
	NotifyTime  string `json:"notify_time" form:"notify_time"`
	TradeStatus string `json:"trade_status" form:"trade_status"`
	Way         string `json:"way" form:"way"`
	ClientId    string `json:"client_id" form:"client_id"`
	OutTradeNo  string `json:"out_trade_no" form:"out_trade_no"`
	TradeNo     string `json:"trade_no" form:"trade_no" gorm:"type:varchar(191); index"`
	PayTime     string `json:"pay_time" form:"pay_time"`
	TotalFee    string `json:"total_fee" form:"total_fee"`
	BuyerId     string `json:"buyer_id" form:"buyer_id"`
	TtSign      string `json:"tt_sign" form:"tt_sign" gorm:"size:512"`
	TtSignType  string `json:"tt_sign_type" form:"tt_sign_type"`
}

type MgtvCallBackRequest struct {
	gorm.Model      `json:"-"`
	Sign            string `json:"sign" form:"sign"`
	Version         string `json:"version" form:"version"`
	Uuid            string `json:"uuid" form:"uuid"`
	BusinessOrderId string `json:"business_order_id" form:"business_order_id" gorm:"type:varchar(191); index"`
	TradeStatus     string `json:"trade_status" form:"trade_status"`
	TradeCreate     string `json:"trade_create" form:"trade_create"`
	TotalFee        string `json:"total_fee" form:"total_fee"`
	ExtData         string `json:"ext_data" form:"ext_data"`
}

type RepireOrderRequest struct {
	GlobalId    int64     `form:"GlobalId" gorm:"index" json:"GlobalId"`
	Sequence    string    `form:"Sequence" json:"Sequence"`
	MailTitle   string    `form:"MailTitle" json:"MailTitle"`
	MailContent string    `form:"MailContent" json:"MailContent"`
	AwardList   string    `form:"AwardList" json:"AwardList"`
	Operator    string    `form:"Operator" gorm:"type:varchar(191); index" json:"Operator"`
	Token       string    `form:"Token" json:"Token" gorm:"-"`
	Channel     string    `form:"Channel" json:"Channel"`
	RepireTime  time.Time `form:"RepireTime" gorm:"index" json:"RepireTime"`
}
