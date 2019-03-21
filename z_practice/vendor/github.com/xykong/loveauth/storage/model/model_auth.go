package model

import (
	"github.com/jinzhu/gorm"
)

type AccountState int

const (
	Active AccountState = iota
	InActive
	Banned
	BanDied
)

type Account struct {
	gorm.Model
	GlobalId       int64        `json:"GlobalId" gorm:"index"`
	Name           string       `json:"Name"`
	State          AccountState `json:"State"`
	UnBanTime      int64        `json:"UnBanTime"`
	LoginTime      int64        `json:"LoginTime"`
	LogoutTime     int64        `json:"LogoutTime"`
	AccumLoginTime int64        `json:"AccumLoginTime"`
	LoginChannel   string       `json:"login_channel"`
}

type GmAccount struct {
	Account  string `json:"account" gorm:"type:varchar(191); index"`
	Password string `json:"password" gorm:"password"`
}

type AccountRealName struct {
	gorm.Model
	GlobalId       int64  `json:"GlobalId" gorm:"index"`
	RealNameMobile string `json:"RealNameMobile"`
	RealName       string `json:"RealName"`
	CardId         string `json:"CardId"`
}

type AccountTag struct {
	gorm.Model
	GlobalId int64  `json:"GlobalId" gorm:"index"`
	Tags     string `json:"Tags"`
}

type AccountAdInfo struct {
	gorm.Model
	GlobalId int64  `json:"GlobalId" gorm:"type:varchar(191); index"`
	Idfa     string `json:"idfa" gorm:"type:varchar(191); index"`
	Idfv     string `json:"idfv" gorm:"type:varchar(191)"`
}

type AuthDevice struct {
	gorm.Model
	GlobalId int64    `json:"GlobalId" gorm:"index"`
	OpenId   string   `json:"OpenId" gorm:"type:varchar(191); index"`
	Platform Platform `json:"Platform"`
}

type AuthGuest struct {
	gorm.Model
	GlobalId         int64    `json:"GlobalId" gorm:"index"`
	OpenId           string   `json:"OpenId" gorm:"type:varchar(191); index"`
	Platform         Platform `json:"Platform"`
	TokenAccess      string   `json:"token_access"`
	ExpirationAccess int64    `json:"expiration_access"`
	Pf               string   `json:"pf"`
	PfKey            string   `json:"pf_key"`
}

type AuthQQ struct {
	gorm.Model
	GlobalId         int64    `json:"GlobalId" gorm:"index"`
	OpenId           string   `json:"OpenId" gorm:"type:varchar(191); index"`
	Platform         Platform `json:"Platform"`
	TokenAccess      string   `json:"token_access"`
	ExpirationAccess int64    `json:"expiration_access"`
	Pf               string   `json:"pf"`
	PfKey            string   `json:"pf_key"`
	NickName         string   `json:"nick_name"`
	Picture          string   `json:"picture"`
}

type AuthYsdkQQ struct {
	gorm.Model
	GlobalId         int64    `json:"GlobalId" gorm:"index"`
	OpenId           string   `json:"OpenId" gorm:"type:varchar(191); index"`
	Platform         Platform `json:"Platform"`
	TokenAccess      string   `json:"token_access"`
	ExpirationAccess int64    `json:"expiration_access"`
	TokenPay         string   `json:"token_pay"`
	Pf               string   `json:"pf"`
	PfKey            string   `json:"pf_key"`
	NickName         string   `json:"nick_name"`
	Picture          string   `json:"picture"`
}

type AuthYsdkWechat struct {
	gorm.Model
	GlobalId         int64    `json:"GlobalId" gorm:"index"`
	OpenId           string   `json:"OpenId" gorm:"type:varchar(191); index"`
	Platform         Platform `json:"Platform"`
	TokenAccess      string   `json:"token_access"`
	ExpirationAccess int64    `json:"expiration_access"`
	TokenRefresh     string   `json:"token_refresh"`
	NickName         string   `json:"nick_name"`
	Picture          string   `json:"picture"`
	UnionId          string   `json:"union_id"`
}

type AuthWechat struct {
	gorm.Model
	GlobalId         int64    `json:"GlobalId" gorm:"index"`
	OpenId           string   `json:"OpenId" gorm:"type:varchar(191); index"`
	Platform         Platform `json:"Platform"`
	TokenAccess      string   `json:"token_access"`
	ExpirationAccess int64    `json:"expiration_access"`
	Pf               string   `json:"pf"`
	PfKey            string   `json:"pf_key"`
	NickName         string   `json:"nick_name"`
	Picture          string   `json:"picture"`
}

type AuthMobile struct {
	gorm.Model
	GlobalId   int64    `json:"GlobalId" gorm:"index"`
	OpenId     string   `json:"OpenId" gorm:"type:varchar(191); index"`
	Platform   Platform `json:"Platform"`
	Location   string   `json:"Location"`
	Token      string   `json:"Token"`
	Expiration int      `json:"Expiration"`
}

type AuthPassword struct {
	gorm.Model
	GlobalId int64  `json:"GlobalId"`
	UserName string `json:"UserName"`
	Password string `json:"Password"`
}

type AuthWeibo struct {
	gorm.Model
	GlobalId         int64    `json:"GlobalId" gorm:"index"`
	OpenId           string   `json:"OpenId" gorm:"type:varchar(191); index"`
	Platform         Platform `json:"Platform"`
	TokenAccess      string   `json:"token_access"`
	ExpirationAccess int64    `json:"expiration_access"`
	NickName         string   `json:"nick_name"`
	Picture          string   `json:"picture"`
}

type AuthQuickAliGames struct {
	gorm.Model
	GlobalId       int64    `json:"GlobalId" gorm:"index"`
	OpenId         string   `json:"OpenId" gorm:"type:varchar(191); index"`
	Platform       Platform `json:"Platform"`
	UserId         string   `json:"user_id"`
	UserName       string   `json:"user_name"`
	Token          string   `json:"token"`
	ChannelVersion string   `json:"channel_version"`
	ChannelName    string   `json:"channel_name"`
	ChannelType    int32    `json:"channel_type"`
}

type AuthQuickYsdk struct {
	gorm.Model
	GlobalId       int64    `json:"GlobalId" gorm:"index"`
	OpenId         string   `json:"OpenId" gorm:"type:varchar(191); index"`
	Platform       Platform `json:"Platform"`
	UserId         string   `json:"user_id"`
	UserName       string   `json:"user_name"`
	Token          string   `json:"token"`
	ChannelVersion string   `json:"channel_version"`
	ChannelName    string   `json:"channel_name"`
	ChannelType    int32    `json:"channel_type"`
}

type AuthQuickMeizu struct {
	gorm.Model
	GlobalId       int64    `json:"GlobalId" gorm:"index"`
	OpenId         string   `json:"OpenId" gorm:"type:varchar(191); index"`
	Platform       Platform `json:"Platform"`
	UserId         string   `json:"user_id"`
	UserName       string   `json:"user_name"`
	Token          string   `json:"token"`
	ChannelVersion string   `json:"channel_version"`
	ChannelName    string   `json:"channel_name"`
	ChannelType    int32    `json:"channel_type"`
}

type AuthQuickM4399 struct {
	gorm.Model
	GlobalId       int64    `json:"GlobalId" gorm:"index"`
	OpenId         string   `json:"OpenId" gorm:"type:varchar(191); index"`
	Platform       Platform `json:"Platform"`
	UserId         string   `json:"user_id"`
	UserName       string   `json:"user_name"`
	Token          string   `json:"token"`
	ChannelVersion string   `json:"channel_version"`
	ChannelName    string   `json:"channel_name"`
	ChannelType    int32    `json:"channel_type"`
}

type AuthQuickKuaikan struct {
	gorm.Model
	GlobalId       int64    `json:"GlobalId" gorm:"index"`
	OpenId         string   `json:"OpenId" gorm:"type:varchar(191); index"`
	Platform       Platform `json:"Platform"`
	UserId         string   `json:"user_id"`
	UserName       string   `json:"user_name"`
	Token          string   `json:"token"`
	ChannelVersion string   `json:"channel_version"`
	ChannelName    string   `json:"channel_name"`
	ChannelType    int32    `json:"channel_type"`
}

type AuthQuickOppo struct {
	gorm.Model
	GlobalId       int64    `json:"GlobalId" gorm:"index"`
	OpenId         string   `json:"OpenId" gorm:"type:varchar(191); index"`
	Platform       Platform `json:"Platform"`
	UserId         string   `json:"user_id"`
	UserName       string   `json:"user_name"`
	Token          string   `json:"token"`
	ChannelVersion string   `json:"channel_version"`
	ChannelName    string   `json:"channel_name"`
	ChannelType    int32    `json:"channel_type"`
}

type AuthQuickIqiyi struct {
	gorm.Model
	GlobalId       int64    `json:"GlobalId" gorm:"index"`
	OpenId         string   `json:"OpenId" gorm:"type:varchar(191); index"`
	Platform       Platform `json:"Platform"`
	UserId         string   `json:"user_id"`
	UserName       string   `json:"user_name"`
	Token          string   `json:"token"`
	ChannelVersion string   `json:"channel_version"`
	ChannelName    string   `json:"channel_name"`
	ChannelType    int32    `json:"channel_type"`
}

type AuthQuickXiaomi struct {
	gorm.Model
	GlobalId       int64    `json:"GlobalId" gorm:"index"`
	OpenId         string   `json:"OpenId" gorm:"type:varchar(191); index"`
	Platform       Platform `json:"Platform"`
	UserId         string   `json:"user_id"`
	UserName       string   `json:"user_name"`
	Token          string   `json:"token"`
	ChannelVersion string   `json:"channel_version"`
	ChannelName    string   `json:"channel_name"`
	ChannelType    int32    `json:"channel_type"`
}

type AuthVivo struct {
	gorm.Model
	GlobalId         int64    `json:"GlobalId" gorm:"index"`
	OpenId           string   `json:"OpenId" gorm:"type:varchar(191); index"`
	Platform         Platform `json:"Platform"`
	AuthToken        string   `json:"auth_token"`
	ExpirationAccess int64    `json:"expiration_access"`
	NickName         string   `json:"NickName"`
	Picture          string   `json:"Picture"`
}

type AuthBilibili struct {
	gorm.Model
	GlobalId     int64    `json:"GlobalId" gorm:"index"`
	OpenId       string   `json:"OpenId" gorm:"type:varchar(191); index"`
	Platform     Platform `json:"Platform"`
	Code         int32    `json:"code"`
	Message      string   `json:"message"`
	UserId       int32    `json:"user_id"`
	UserName     string   `json:"user_name"`
	NickName     string   `json:"nick_name"`
	AccessToken  string   `json:"access_token"`
	ExpireTimes  string   `json:"expire_times"`
	RefreshToken string   `json:"refresh_token"`
}

type AuthHuawei struct {
	gorm.Model
	GlobalId     int64    `json:"GlobalId" gorm:"index"`
	OpenId       string   `json:"OpenId" gorm:"type:varchar(191); index"`
	Platform     Platform `json:"Platform"`
	PlayerId     string   `json:"player_id"`
	DisplayName  string   `json:"display_name"`
	PlayerLevel  int      `json:"player_level"`
	IsAuth       int      `json:"is_auth"`
	Ts           string   `json:"ts"`
	GameAuthSign string   `json:"game_auth_sign" gorm:"type:varchar(1023);"`
}

type AuthDouyin struct {
	gorm.Model
	GlobalId         int64    `json:"GlobalId" gorm:"index"`
	OpenId           string   `json:"OpenId" gorm:"type:varchar(191); index"`
	Platform         Platform `json:"Platform"`
	Uid              uint64   `json:"uid"`
	UserType         int32    `json:"UserType"`
	AccessToken      string   `json:"AccessToken"`
	ExpirationAccess int64    `json:"ExpirationAccess"`
	NickName         string   `json:"NickName"`
	Picture          string   `json:"Picture"`
}

type AuthMgtv struct {
	gorm.Model
	GlobalId         int64    `json:"GlobalId" gorm:"index"`
	OpenId           string   `json:"OpenId" gorm:"type:varchar(191); index"`
	Platform         Platform `json:"Platform"`
	Uid              uint64   `json:"uid"`
	AccessToken      string   `json:"AccessToken"`
	Ticket           string   `json:"ticket"`
	ExpirationAccess int64    `json:"ExpirationAccess"`
	NickName         string   `json:"NickName"`
	Picture          string   `json:"Picture"`
}

//存储快看回调数据
type KuaikanCb struct {
	gorm.Model
	Idfa         string `json:"idfa" gorm:"type:varchar(191); index"`
	CallBackTime int64  `json:"callBackTime"`
}

// 平台项: ios, android, desktop, web
type Platform string

const (
	PlatformIOS       Platform = "ios"
	PlatformAndroid            = "android"
	PlatformDesktop            = "desktop"
	PlatformWeb                = "web"
	PlatformUndefined          = "undefined"
	PlatformAll                = "*"
)

var Platforms = []Platform{
	PlatformIOS,
	PlatformAndroid,
	PlatformDesktop,
	PlatformWeb,
	PlatformUndefined,
}

// 登录提供商: device, tencent_qq, tencent_wechat, tencent_guest
type Vendor string

const (
	VendorBilibili      = "bilibili"
	VendorDevice        = "device"
	VendorHuawei        = "huawei"
	VendorMobile        = "mobile"
	VendorMsdkGuest     = "madk_guest"
	VendorMsdkQQ        = "msdk_qq"
	VendorMsdkWechat    = "msdk_wechat"
	VendorVivo          = "vivo"
	VendorWeibo         = "weibo"
	VendorYsdkQQ        = "ysdk_qq"
	VendorYsdkWechat    = "ysdk_wechat"
	VendorQuickAligames = "quick_aligames" //uc九游
	VendorQuickOppo     = "quick_oppo"
	VendorQuickM4399    = "quick_m4399"
	VendorQuickYsdk     = "quick_ysdk"    //应用宝(YSDK)
	VendorQuickIqiyi    = "quick_iqiyi"   //爱奇艺
	VendorQuickMeiZu    = "quick_meizu"   //魅族
	VendorQuickKuaiKan  = "quick_kuaikan" //快看
	VendorQuickXiaomi   = "quick_xiaomi"  // 小米
	VendorDouyin        = "douyin"        // 抖音
	VendorMgtv          = "mgtv"          // 芒果tv
	VendorAll           = "*"
)

const ChannelAll = "*"

var Vendors = []Vendor{
	VendorBilibili,
	VendorDevice,
	VendorHuawei,
	VendorMobile,
	VendorMsdkGuest,
	VendorMsdkQQ,
	VendorMsdkWechat,
	VendorVivo,
	VendorWeibo,
	VendorYsdkQQ,
	VendorYsdkWechat,
	VendorQuickAligames,
	VendorQuickOppo,
	VendorQuickM4399,
	VendorQuickYsdk,
	VendorQuickIqiyi,
	VendorQuickMeiZu,
	VendorQuickKuaiKan,
	VendorQuickXiaomi,
	VendorDouyin,
	VendorMgtv,
}

// used for redis
type TokenRecord struct {
	GlobalId  int64    `json:"global_id"`
	Vendor    Vendor   `json:"vendor"`
	Platform  Platform `json:"Platform"`
	Timestamp int64    `json:"timestamp"`
}

// used for redis
type RefreshTokenRecord struct {
	GlobalId          int64    `json:"global_id"`
	Token             string   `json:"token"`
	RefreshToken      string   `json:"refresh_token"`
	Vendor            Vendor   `json:"vendor"`
	Platform          Platform `json:"platform"`
	Timestamp         int64    `json:"timestamp"`
	ExpirationSeconds int64    `json:"expiration_seconds"`
}

// used for redis
type Profile struct {
	GlobalId          int64         `json:"global_id"`
	Token             string        `json:"token"`
	RefreshToken      string        `json:"refresh_token"`
	ExpirationSeconds int64         `json:"expiration_seconds"`
	Timestamp         int64         `json:"timestamp"`
	Name              string        `json:"name"`
	Auth              DoAuthRequest `json:"state"`
	Vendor            Vendor        `json:"vendor"`
	Platform          Platform      `json:"platform"`
	NickName          string        `json:"nick_name"`
	Picture           string        `json:"picture"`
}

// used for redis
type SMSToken struct {
	//Auth              DoAuthRequest `json:"state"`
	Mobile            string `json:"mobile"`
	PicToken          string `json:"pic_token"`
	Token             string `json:"token"`
	Timestamp         int64  `json:"timestamp"`
	ExpirationSeconds int64  `json:"expiration_seconds"`
}

// Notable is a model in a transitive package.
// it's used for embedding in another model
//
// swagger:model Extra
//noinspection ALL
type Extra struct {
	GameSvrId      string  `form:"GameSvrId" json:"GameSvrId"`
	VGameAppid     string  `form:"vGameAppid" json:"vGameAppid"`
	PlatId         int     `form:"PlatID" json:"PlatID"`
	IZoneAreaID    int     `form:"iZoneAreaID" json:"iZoneAreaID"`
	Vopenid        string  `form:"vopenid" json:"vopenid"`
	ClientVersion  string  `form:"ClientVersion" json:"ClientVersion"`
	SystemSoftware string  `form:"SystemSoftware" json:"SystemSoftware"`
	SystemHardware string  `form:"SystemHardware" json:"SystemHardware"`
	TelecomOper    string  `form:"TelecomOper" json:"TelecomOper"`
	Network        string  `form:"Network" json:"Network"`
	ScreenWidth    int     `form:"ScreenWidth" json:"ScreenWidth"`
	ScreenHight    int     `form:"ScreenHight" json:"ScreenHight"`
	Density        float64 `form:"Density" json:"Density"`
	RegChannel     string  `form:"RegChannel" json:"RegChannel"`
	LoginChannel   string  `form:"LoginChannel" json:"LoginChannel"`
	CpuHardware    string  `form:"CpuHardware" json:"CpuHardware"`
	Memory         int     `form:"Memory" json:"Memory"`
	GLRender       string  `form:"GLRender" json:"GLRender"`
	GLVersion      string  `form:"GLVersion" json:"GLVersion"`
	DeviceId       string  `form:"DeviceId" json:"DeviceId"`
	VClientIP      string  `form:"vClientIP" json:"vClientIP"`
	MidasZoneId    string  `form:"midasZoneId" json:"midasZoneId"`
}

type VendorQuickReq struct {
	UserId         string `form:"UserId" json:"UserId"`
	UserName       string `form:"UserName" json:"UserName"`
	Token          string `form:"Token" json:"Token"`
	ChannelVersion string `form:"ChannelVersion" json:"ChannelVersion"`
	ChannelName    string `form:"ChannelName" json:"ChannelName"`
	ChannelType    int32  `form:"ChannelType" json:"ChannelType"`
}

type AuthExtra struct {
	IDFA string `form:"IDFA" json:"IDFA"`
	IDFV string `form:"IDFV" json:"IDFV"`
}

// DoAuthRequest represents the message get from msdk client
type DoAuthRequest struct {
	OpenId           string    `form:"OpenId" json:"OpenId" binding:"required,min=6,max=64"`
	Vendor           Vendor    `form:"Vendor" json:"Vendor" binding:"required"`
	Platform         Platform  `form:"Platform" json:"Platform" binding:"required"`
	Coupon           string    `form:"Coupon" json:"Coupon"`
	Channel          string    `form:"Channel" json:"Channel"`
	ClientVersion    string    `form:"ClientVersion" json:"ClientVersion"`
	ResourcesVersion string    `form:"ResourcesVersion" json:"ResourcesVersion"`
	ClientIp         string    `form:"ClientIp" json:"ClientIp"`
	DeviceId         string    `form:"DeviceId" json:"DeviceId"`
	Extra            AuthExtra `form:"Extra" json:"Extra"`

	VendorBilibili struct {
		LoginResult struct {
			Code         int32  `form:"Code" json:"Code"`
			Message      string `form:"Message" json:"Message"`
			UserId       int32  `form:"UserId" json:"UserId"`             //	用户ID	String	10001
			UserName     string `form:"UserName" json:"UserName"`         //	用户名昵称	String	gamenick
			NickName     string `form:"NickName" json:"NickName"`         //	用户名昵称	String	gamenick
			AccessToken  string `form:"AccessToken" json:"AccessToken"`   //	访问令牌	String	fdae8922a3b3d06a4e40882ac9f37a7e
			ExpireTimes  string `form:"ExpireTimes" json:"ExpireTimes"`   //	会话过期时间	String	1389262844（10位）
			RefreshToken string `form:"RefreshToken" json:"RefreshToken"` //	刷新令牌	String	0d5ddfa364d51359e6243892bf0a965c
		} `form:"LoginResult" json:"LoginResult"`
	} `form:"Bilibili" json:"Bilibili"`

	VendorMobile struct {
		VerifyCode string `form:"VerifyCode" json:"VerifyCode"`
	} `form:"Mobile" json:"Mobile"`

	VendorHuawei struct {
		GameUserData struct {
			PlayerId     string `form:"PlayerId" json:"PlayerId"`
			DisplayName  string `form:"DisplayName" json:"DisplayName"`
			PlayerLevel  int    `form:"PlayerLevel" json:"PlayerLevel"`
			IsAuth       int    `form:"IsAuth" json:"IsAuth"`
			Ts           string `form:"Ts" json:"Ts"`
			GameAuthSign string `form:"GameAuthSign" json:"GameAuthSign"`
		} `form:"GameUserData" json:"GameUserData"`
	} `form:"Huawei" json:"Huawei"`

	VendorMsdkGuest struct {
		AccessTokenValue      string `form:"AccessTokenValue" json:"AccessTokenValue"`
		AccessTokenExpiration int64  `form:"AccessTokenExpiration" json:"AccessTokenExpiration"`
		Pf                    string `form:"Pf" json:"Pf"`
		PfKey                 string `form:"PfKey" json:"PfKey"`
	} `form:"MsdkGuest" json:"MsdkGuest"`

	VendorMsdkQQ struct {
		AccessTokenValue      string `form:"AccessTokenValue" json:"AccessTokenValue"`
		AccessTokenExpiration int64  `form:"AccessTokenExpiration" json:"AccessTokenExpiration"`
		Pf                    string `form:"Pf" json:"Pf"`
		PfKey                 string `form:"PfKey" json:"PfKey"`
	} `form:"MsdkQQ" json:"MsdkQQ"`

	VendorMsdkWechat struct {
		AccessTokenValue      string `form:"AccessTokenValue" json:"AccessTokenValue"`
		AccessTokenExpiration int64  `form:"AccessTokenExpiration" json:"AccessTokenExpiration"`
		Pf                    string `form:"Pf" json:"Pf"`
		PfKey                 string `form:"PfKey" json:"PfKey"`
	} `form:"MsdkWechat" json:"MsdkWechat"`

	VendorVivo struct {
		AuthToken        string `form:"AuthToken" json:"AuthToken"`
		ExpirationAccess int64  `form:"ExpirationAccess" json:"ExpirationAccess"`
		NickName         string `form:"NickName" json:"NickName"`
		Picture          string `form:"Picture" json:"Picture"`
	} `form:"Vivo" json:"Vivo"`

	VendorWeibo struct {
		TokenAccess      string `form:"TokenAccess" json:"TokenAccess"`
		ExpirationAccess int64  `form:"ExpirationAccess" json:"ExpirationAccess"`
		NickName         string `form:"NickName" json:"NickName"`
		Picture          string `form:"Picture" json:"Picture"`
	} `form:"Weibo" json:"Weibo"`

	VendorYsdkQQ struct {
		TokenAccess      string `form:"TokenAccess" json:"TokenAccess"`
		ExpirationAccess int64  `form:"ExpirationAccess" json:"ExpirationAccess"`
		TokenPay         string `form:"TokenPay" json:"TokenPay"`
		Pf               string `form:"Pf" json:"Pf"`
		PfKey            string `form:"PfKey" json:"PfKey"`
		NickName         string `form:"NickName" json:"NickName"`
		Picture          string `form:"Picture" json:"Picture"`
	} `form:"YsdkQQ" json:"YsdkQQ"`

	VendorYsdkWechat struct {
		TokenAccess      string `form:"TokenAccess" json:"TokenAccess"`
		ExpirationAccess int64  `form:"ExpirationAccess" json:"ExpirationAccess"`
		TokenRefresh     string `form:"TokenRefresh" json:"TokenRefresh"`
		NickName         string `form:"NickName" json:"NickName"`
		Picture          string `form:"Picture" json:"Picture"`
		UnionId          string `form:"UnionId" json:"UnionId"`
	} `form:"YsdkWechat" json:"YsdkWechat"`

	VendorQuickAligames VendorQuickReq `form:"QuickAligames" json:"QuickAligames"`
	VendorQuickOppo     VendorQuickReq `form:"QuickOppo" json:"QuickOppo"`
	VendorQuickM4399    VendorQuickReq `form:"QuickM4399" json:"QuickM4399"`
	VendorQuickYsdk     VendorQuickReq `form:"QuickYsdk" json:"QuickYsdk"`
	VendorQuickIqiyi    VendorQuickReq `form:"QuickIqiyi" json:"QuickIqiyi"`
	VendorQuickMeiZu    VendorQuickReq `form:"QuickMeizu" json:"QuickMeizu"`
	VendorQuickKuaiKan  VendorQuickReq `form:"QuickKuaikan" json:"QuickKuaikan"`
	VendorQuickXiaomi   VendorQuickReq `form:"QuickXiaomi" json:"QuickXiaomi"`

	VendorDouyin struct {
		GameUserData struct {
			OpenId      string `form:"OpenId" json:"OpenId"`
			AccessToken string `form:"AccessToken" json:"AccessToken"`
			Uid         uint64 `form:"Uid" json:"Uid"`
			UserType    int32  `form:"UserType" json:"UserType"`
		} `form:"GameUserData" json:"GameUserData"`
	} `form:"Douyin" json:"Douyin"`

	VendorMgtv struct {
		GameUserData struct {
			OpenId       string `form:"OpenId" json:"OpenId"`
			NickName     string `form:"NickName" json:"NickName"`
			Ticket       string `form:"Ticket" json:"Ticket"`
			LoginAccount string `form:"LoginAccount" json:"LoginAccount"`
			Result       int32  `form:"Result" json:"Result"`
			ThirdId      string `form:"ThirdId" json:"ThirdId"`
		} `form:"GameUserData" json:"GameUserData"`
	} `form:"Mgtv" json:"Mgtv"`
}
