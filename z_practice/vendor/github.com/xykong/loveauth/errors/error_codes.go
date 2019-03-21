package errors

import "fmt"

const (
	Ok     Code = iota
	Failed

	PartialContent
	IllegalRequest = 5 // 非法请求

	AuthTokenInvalid         = 1001 // token无效
	AuthUserNotExist         = 1002 // 用户不存在
	AuthGoToLogin            = 1004 // 重新登录
	AuthRepeatLogin          = 1005 // 顶号
	AuthAccountAlreadyActive = 1006 // 账户已经处于激活状态
	AuthRegisterByCoupon     = 1007 // 需要注册邀请码
	AuthRegisterClosed       = 1008 // 关闭账户注册
	AuthRegisterByMobile     = 1009 // 需要手机验证码
	SendSMSNeedPicture       = 1010 //发送手机验证码需要图片验证码
	SendSMSFrequently        = 1011 //手机验证码发送频繁
	SendSMSOneDayLimit       = 1012 //一天内发送超限
	PictureCodeError         = 1013 //图片验证码错误
	SMSCodeExpire            = 1014 //短信验证码过期
	SMSCodeError             = 1015 //短信验证码错误
	RepeatBind               = 1016 //重复绑定
	ForbidBindOther          = 1017 //已绑定其他游戏账号，绑定失败
	CardIdError              = 1018 //身份证号码格式错误
	RealNameRepeat           = 1019 //重复实名认证
	RealNameError            = 1020 //用户名错误
	BindOtherAccount         = 1021 //已绑定其他游戏账号，绑定失败
	ThirdAccessTokenExpire   = 1022 //第三方授权已经过期
	ServerFail               = 1023 // 连接服务器异常
	InvalidPlatformSpecified = 1024 // 应用暂不支持该平台
	InvalidVendorSpecified   = 1025 // 应用暂不支持该第三方登录
	VerifyBy3rdPartyFailed   = 1026 // 认证失败
	CouponFailed             = 1027 // 请输入正确的激活码
	CouponUsed               = 1028 // 该激活码已被使用，请重试
	AccountWasBanned         = 1032 // 该账号由于多次违规，已被封停，如有疑问可以联系客服咨询
	BindCheckErr             = 1037 // 绑定校验失败
	SaveAccountRealNameFail  = 1044 // 实名认证失败

	DiamondShortage = 1202 //钻石不足
	GoldShortage    = 1206 //游戏币不足

	ShopInfoNotFound     = 2201 // 商城信息获取失败
	ShopCallBackAgain    = 2202 // 轮询请求失败
	ShopDataException    = 2203 // 数据异常
	ShopCommitRmb        = 2204 // rmb购买下单失败
	ShopDiamondBuy       = 2205 // 钻石购买道具失败
	ShopQueryGeted       = 2206 // 已领取
	ShopNotPay           = 2207 // 查询未成功,米大师消费失败
	ShopOrderMiss        = 2208 // 订单不存在
	ShopIdMiss           = 2209 // 商品已不存在
	ShopDiamondBuyNumOut = 2210 // 购买数量超限

	QuickPayFailed = 2215 //quick购买失败

	QueryAccessTokenFailed = 2501 // 登录信息已失效

	WxUnifiedOrderFailed   = 2216 //微信支付下单失败

	QPayUnifiedOrderFailed = 2218 //qq支付下单失败
)

type Code int

// errorString is a trivial implementation of error.
type Type struct {
	Code    Code
	Message string
}

// New returns an error that formats as the given text.
func New(text string) error {
	return &Type{Failed, text}
}

// New returns an error that formats as the given text.
func NewCode(code Code) error {
	return NewCodeString(code, "")
}

// New returns an error that formats as the given text.
func NewCodeString(code Code, format string, a ...interface{}) error {
	return &Type{code, fmt.Sprintf(format, a...)}
}

// New returns an error that formats as the given text.
func NewCodeError(code Code, err error) error {
	return &Type{code, err.Error()}
}

func (e *Type) Error() string {
	//return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
	return e.Message
}
