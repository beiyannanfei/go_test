package auth

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/settings"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"github.com/xykong/loveauth/utils"
	"github.com/xykong/loveauth/utils/log"
	"github.com/xykong/tlog"
	"time"
)

func createAccount(user *model.DoAuthRequest) *model.Account {

	globalId := utils.GenerateId().Int64()

	state := model.Active

	workingMode := WorkingMode(settings.GetString("loveauth", "auth.WorkingMode"))

	switch workingMode {
	case Normal:

		state = model.Active
		if inRegisterByCouponChannels(user.Channel) {

			state = model.InActive
		}
	case RegisterByCoupon, RegisterClosed:

		state = model.InActive
	default:
		logrus.WithFields(logrus.Fields{
			"workingMode": settings.GetString("loveauth", "auth.WorkingMode"),
		}).Error("Invalid WorkingMode.")
	}

	//存储idfa信息
	if user.Extra.IDFA != "" || user.Extra.IDFV != "" {
		adInfo := &model.AccountAdInfo{
			GlobalId: globalId,
			Idfa:     user.Extra.IDFA,
			Idfv:     user.Extra.IDFV,
		}

		storage.Insert(storage.AuthDatabase(), adInfo)
	}

	account := &model.Account{

		GlobalId:     globalId,
		Name:         "",
		State:        state,
		LoginChannel: user.Channel,
	}

	return account
}

func inRegisterByCouponChannels(loginChannel string) bool {

	registerByCouponChannels := settings.GetStringSlice("loveauth", "auth.RegisterByCouponChannels")
	for _, channel := range registerByCouponChannels {

		if loginChannel == channel {

			return true
		}
	}

	return false
}

func NewTLogger(vendor model.Vendor) *tlog.Tlogger {

	address := settings.GetString("tencent", "tlog.address")
	if address == "" {

		logrus.Warning("NewTLogger setting get address is nil")
		return nil
	}

	tlogger, err := tlog.Dial(address)
	if tlogger == nil || err != nil {

		return nil
	}

	gameSvrId, gameAppID := GetVendorConfigSetting(vendor)

	tlogger.SetRequired(gameSvrId, gameAppID, 0)

	return tlogger
}

/*
  <!--//////////////////////////////////////////////
    ///////玩家注册表///////////////////////////////
   /////////////////////////////////////////////////-->
  <struct  name="PlayerRegister"  version="1" desc="(必填)玩家注册">
    <entry  name="GameSvrId"          type="string"        size="25"    desc="(必填)登录的游戏服务器编号" />
    <entry  name="dtEventTime"        type="datetime"                    desc="(必填)游戏事件的时间, 格式 YYYY-MM-DD HH:MM:SS" />
    <entry  name="vGameAppid"         type="string"        size="32"                        desc="(必填)游戏APPID" />
    <entry  name="PlatID"             type="int"                        defaultvalue="0"    desc="(必填)ios 0 /android 1"/>
    <entry  name="iZoneAreaID"        type="int"           index="1"    defaultvalue="0"    desc="(必填)针对分区分服的游戏填写分区id，用来唯一标示一个区；非分区分服游戏请填写0"/>
    <entry  name="vopenid"            type="string"        size="64"                        desc="(必填)用户OPENID号" />
    <entry  name="ClientVersion"      type="string"        size="64"    defaultvalue="NULL" desc="(可选)客户端版本"/>
    <entry  name="SystemSoftware"     type="string"        size="64"    defaultvalue="NULL" desc="(可选)移动终端操作系统版本"/>
    <entry  name="SystemHardware"     type="string"        size="64"    defaultvalue="NULL" desc="(可选)移动终端机型"/>
    <entry  name="TelecomOper"        type="string"        size="64"    defaultvalue="NULL" desc="(必填)运营商"/>
    <entry  name="Network"            type="string"        size="64"    defaultvalue="NULL" desc="(可选)3G/WIFI/2G"/>
    <entry  name="ScreenWidth"        type="int"                        defaultvalue="0"    desc="(可选)显示屏宽度"/>
    <entry  name="ScreenHight"        type="int"                        defaultvalue="0"    desc="(可选)显示屏高度"/>
    <entry  name="Density"            type="float"                      defaultvalue="0"    desc="(可选)像素密度"/>
    <entry  name="RegChannel"         type="int"                        defaultvalue="0"    desc="(必填)注册渠道"/>
    <entry  name="CpuHardware"        type="string"        size="64"    defaultvalue="NULL" desc="(可选)cpu类型|频率|核数"/>
    <entry  name="Memory"             type="int"                        defaultvalue="0"    desc="(可选)内存信息单位M"/>
    <entry  name="GLRender"           type="string"       size="64"     defaultvalue="NULL" desc="(可选)opengl render信息"/>
    <entry  name="GLVersion"          type="string"       size="64"     defaultvalue="NULL"  desc="(可选)opengl版本信息"/>
    <entry  name="DeviceId"           type="string"       size="64"     defaultvalue="NULL"  desc="(可选)设备ID"/>
    <entry  name="vClientIP"          type="string"       size="64"     defaultvalue="NULL"  desc="(必填)客户端IP(后台服务器记录与玩家通信时的IP地址)"/>
	<entry name="Format" type="string" size="64" defaultvalue="1.0.0" desc="Format version of this log"/>
   </struct>
*/
func LogPlayerRegister(user *model.DoAuthRequest, globalId int64) {

	now := time.Now()
	dtEventTime := now.Format("2006-01-02 15:04:05")
	format := settings.GetString("tencent", "tlog.format")
	gameSvrId, gameAppId := GetVendorConfigSetting(user.Vendor)

	/*if tlogger != nil {

		tlogger.Log(
			tlog.PublicConfig{"PlayerRegister", gameSvrId, dtEventTime, gameAppId, extra.PlatId, extra.IZoneAreaID},
			openId,
			extra.ClientVersion,
			extra.SystemSoftware,
			extra.SystemHardware,
			extra.TelecomOper,
			extra.Network,
			extra.ScreenWidth,
			extra.ScreenHight,
			extra.Density,
			extra.RegChannel,
			extra.CpuHardware,
			extra.Memory,
			extra.GLRender,
			extra.GLVersion,
			extra.DeviceId,
			extra.VClientIP,
			format)
	}*/

	logrus.WithFields(logrus.Fields{
		"FlowName":      "PlayerRegister",
		"GameSvrId":     gameSvrId,
		"EventTime":     dtEventTime,
		"GameAppId":     gameAppId,
		"PlatId":        GetPlatId(user.Platform),
		"OpenId":        user.OpenId,
		"RoleId":        globalId,
		"Format":        format,
		"ClientVersion": user.ClientVersion,
		"RegChannel":    user.Channel,
		"DeviceId":      user.DeviceId,
		"ClientIP":      user.ClientIp,
		"LoginChannel":  user.Channel,
	}).Info("Tlog PlayerRegister")

	if nil != log.BILogger {
		log.BILogger.Printf("%s,%s,%s,%s,%d,%s,%d,%s,%s,%s,%s,%s,%s",
			"PlayerRegister", gameSvrId, dtEventTime, gameAppId, GetPlatId(user.Platform),
			user.OpenId, globalId, format, user.ClientVersion, user.Channel, user.DeviceId, user.ClientIp, user.Channel)
	} else {
		logrus.WithFields(logrus.Fields{"reason": "BILogger is nil"}).Info("LogPlayerRegister")
	}

	if true == settings.GetBool("tencent", "tlog.mysql.enable_ext") {
		register := new(storage.PlayerRegister)
		register.Channel = user.Channel
		register.GlobalId = globalId
		register.TimeKey = int(now.Unix())
		register.DeviceId = user.DeviceId
		register.GameSvrId = gameSvrId
		err := storage.Insert(storage.ExtDatabase(), register)
		if nil != err {
			logrus.WithFields(logrus.Fields{
				"err": err.Error(),
			}).Error("LogPlayerRegister")
		}
	}
}

/*
     <!--//////////////////////////////////////////////
    ///////玩家登录表///////////////////////////////
   /////////////////////////////////////////////////-->
  <struct  name="PlayerLogin"  version="1" desc="(必填)玩家登陆">
    <entry  name="GameSvrId"          type="string"        size="25"                            desc="(必填)登录的游戏服务器编号" />
    <entry  name="dtEventTime"        type="datetime"                                           desc="(必填)游戏事件的时间, 格式 YYYY-MM-DD HH:MM:SS" />
    <entry  name="vGameAppid"         type="string"        size="32"                            desc="(必填)游戏APPID" />
    <entry  name="PlatID"             type="int"                        defaultvalue="0"        desc="(必填)ios 0/android 1"/>
    <entry  name="iZoneAreaID"        type="int"           index="1"    defaultvalue="0"        desc="(必填)针对分区分服的游戏填写分区id，用来唯一标示一个区；非分区分服游戏请填写0"/>
    <entry  name="vopenid"            type="string"        size="64"                            desc="(必填)用户OPENID号" />
    <entry  name="Level"              type="int"                                                desc="(必填)等级" />
    <entry  name="PlayerFriendsNum"   type="int"                                                desc="(必填)玩家好友数量"/>
    <entry  name="ClientVersion"      type="string"        size="64"    defaultvalue="NULL"     desc="(必填)客户端版本"/>
    <entry  name="SystemSoftware"     type="string"        size="64"    defaultvalue="NULL"     desc="(可选)移动终端操作系统版本"/>
    <entry  name="SystemHardware"     type="string"        size="64"    defaultvalue="NULL"     desc="(必填)移动终端机型"/>
    <entry  name="TelecomOper"        type="string"        size="64"    defaultvalue="NULL"     desc="(必填)运营商"/>
    <entry  name="Network"            type="string"        size="64"    defaultvalue="NULL"     desc="(必填)3G/WIFI/2G"/>
    <entry  name="ScreenWidth"        type="int"                        defaultvalue="0"        desc="(可选)显示屏宽度"/>
    <entry  name="ScreenHight"        type="int"                        defaultvalue="0"        desc="(可选)显示屏高度"/>
    <entry  name="Density"            type="float"                      defaultvalue="0"        desc="(可选)像素密度"/>
    <entry  name="LoginChannel"       type="int"                        defaultvalue="0"        desc="(必填)登录渠道"/>
    <entry  name="vRoleID"            type="string"        size="64"    defaultvalue="NULL"     desc="(必填)玩家角色ID"/>
    <entry  name="vRoleName"          type="string"        size="64"    defaultvalue="NULL"     desc="(必填)玩家角色名"/>
    <entry  name="CpuHardware"        type="string"        size="64"    defaultvalue="NULL"     desc="(可选)cpu类型-频率-核数"/>
    <entry  name="Memory"             type="int"                        defaultvalue="0"        desc="(可选)内存信息单位M"/>
    <entry  name="GLRender"           type="string"        size="64"    defaultvalue="NULL"     desc="(可选)opengl render信息"/>
    <entry  name="GLVersion"          type="string"        size="64"    defaultvalue="NULL"     desc="(可选)opengl版本信息"/>
    <entry  name="DeviceId"           type="string"        size="64"    defaultvalue="NULL"     desc="(可选)设备ID"/>
    <entry  name="vClientIP"          type="string"        size="64"    defaultvalue="NULL"     desc="(必填)客户端IP(后台服务器记录与玩家通信时的IP地址)"/>
	<entry name="Format" type="string" size="64" defaultvalue="1.0.0" desc="Format version of this log"/>
  </struct>
*/
//func LogPlayerLogin(profile *model.Profile) {
//
//	Level := 0
//	PlayerFriendsNum := 0
//
//	//address := settings.GetString("tencent", "tlog.address")
//	format := settings.GetString("tencent", "tlog.format")
//	dtEventTime := time.Now().Format("2006-01-02 15:04:05")
//	gameSveId, gameAppId := GetVendorConfigSetting(profile.Vendor)
//	/*if address != "" {
//
//		tlog.Log(
//			tlog.PublicConfig{"PlayerLogin", profile.Auth.Extra.GameSvrId, dtEventTime, profile.Auth.Extra.VGameAppid, profile.Auth.Extra.PlatId, profile.Auth.Extra.IZoneAreaID},
//			profile.Auth.OpenId,
//			Level,
//			PlayerFriendsNum,
//			profile.Auth.Extra.ClientVersion,
//			profile.Auth.Extra.SystemSoftware,
//			profile.Auth.Extra.SystemHardware,
//			profile.Auth.Extra.TelecomOper,
//			profile.Auth.Extra.Network,
//			profile.Auth.Extra.ScreenWidth,
//			profile.Auth.Extra.ScreenHight,
//			profile.Auth.Extra.Density,
//			profile.Auth.Extra.LoginChannel,
//			profile.GlobalId,
//			profile.Name,
//			profile.Auth.Extra.CpuHardware,
//			profile.Auth.Extra.Memory,
//			profile.Auth.Extra.GLRender,
//			profile.Auth.Extra.GLVersion,
//			profile.Auth.Extra.DeviceId,
//			profile.Auth.Extra.VClientIP,
//			format,
//		)
//	}*/
//
//	logrus.WithFields(logrus.Fields{
//		"GameSvrId":        gameSveId,
//		"dtEventTime":      dtEventTime,
//		"vGameAppid":       gameAppId,
//		"Platform":         profile.Platform,
//		"vopenid":          profile.Auth.OpenId,
//		"Level":            Level,
//		"PlayerFriendsNum": PlayerFriendsNum,
//		"ClientVersion":    profile.Auth.ClientVersion,
//		"LoginChannel":     profile.Auth.Channel,
//		"vRoleID":          profile.GlobalId,
//		"vRoleName":        profile.Name,
//		"DeviceId":         profile.Auth.DeviceId,
//		"vClientIP":        profile.Auth.ClientIp,
//		"FlowName":         "PlayerLogin",
//		"Format":           format,
//	}).Info("Tlog PlayerLogin")
//}

/*

  <!--//////////////////////////////////////////////
    ///////玩家登出表///////////////////////////////
   /////////////////////////////////////////////////-->
  <struct  name="PlayerLogout" version="1" desc="(必填)玩家登出">
    <entry  name="GameSvrId"          type="string"      size="25"                              desc="(必填)登录的游戏服务器编号" />
    <entry  name="dtEventTime"        type="datetime"                                           desc="(必填)游戏事件的时间, 格式 YYYY-MM-DD HH:MM:SS" />
    <entry  name="vGameAppid"         type="string"        size="32"                            desc="(必填)游戏APPID" />
    <entry  name="PlatID"             type="int"                        defaultvalue="0"        desc="(必填)ios 0/android 1"/>
    <entry  name="iZoneAreaID"        type="int"           index="1"    defaultvalue="0"        desc="(必填)针对分区分服的游戏填写分区id，用来唯一标示一个区；非分区分服游戏请填写0"/>
    <entry  name="vopenid"            type="string"        size="64"                            desc="(必填)用户OPENID号" />
    <entry  name="OnlineTime"         type="int"                                                desc="(必填)本次登录在线时间(秒)" />
    <entry  name="Level"              type="int"                                                desc="(必填)等级" />
    <entry  name="PlayerFriendsNum"   type="int"                                                desc="(必填)玩家好友数量"/>
    <entry  name="ClientVersion"      type="string"        size="64"    defaultvalue="NULL"     desc="(必填)客户端版本"/>
    <entry  name="SystemSoftware"     type="string"        size="64"    defaultvalue="NULL"     desc="(可选)移动终端操作系统版本"/>
    <entry  name="SystemHardware"     type="string"        size="64"    defaultvalue="NULL"     desc="(必填)移动终端机型"/>
    <entry  name="TelecomOper"        type="string"        size="64"    defaultvalue="NULL"     desc="(必填)运营商"/>
    <entry  name="Network"            type="string"        size="64"    defaultvalue="NULL"     desc="(必填)3G/WIFI/2G"/>
    <entry  name="ScreenWidth"        type="int"                        defaultvalue="0"        desc="(可选)显示屏宽度"/>
    <entry  name="ScreenHight"        type="int"                        defaultvalue="0"        desc="(可选)显示高度"/>
    <entry  name="Density"            type="float"                      defaultvalue="0"        desc="(可选)像素密度"/>
    <entry  name="LoginChannel"       type="int"                        defaultvalue="0"        desc="(可选)登录渠道"/>
    <entry  name="CpuHardware"        type="string"        size="64"    defaultvalue="NULL"     desc="(可选)cpu类型;频率;核数"/>
    <entry  name="Memory"             type="int"                        defaultvalue="0"        desc="(可选)内存信息单位M"/>
    <entry  name="GLRender"           type="string"        size="64"    defaultvalue="NULL"     desc="(可选)opengl render信息"/>
    <entry  name="GLVersion"          type="string"        size="64"    defaultvalue="NULL"     desc="(可选)opengl版本信息"/>
    <entry  name="DeviceId"           type="string"        size="64"    defaultvalue="NULL"     desc="(可选)设备ID"/>
    <entry  name="vClientIP"          type="string"        size="64"    defaultvalue="NULL"     desc="(必填)客户端IP(后台服务器记录与玩家通信时的IP地址)"/>
	<entry name="Format" type="string" size="64" defaultvalue="1.0.0" desc="Format version of this log"/>
  </struct>
*/
func LogPlayerLogout(profile *model.Profile, logoutTime int64) {

	//Level := 0
	PlayerFriendsNum := 0
	OnlineTime := logoutTime - profile.Timestamp
	if OnlineTime == 0 {

		OnlineTime = 1
	}

	//address := settings.GetString("tencent", "tlog.address")
	format := settings.GetString("tencent", "tlog.format")
	dtEventTime := time.Now().Format("2006-01-02 15:04:05")
	gameSvrId, gameAppId := GetVendorConfigSetting(profile.Vendor)

	/*if address != "" {

		tlog.Log(
			tlog.PublicConfig{"PlayerLogout", profile.Auth.Extra.GameSvrId, dtEventTime, profile.Auth.Extra.VGameAppid, profile.Auth.Extra.PlatId, profile.Auth.Extra.IZoneAreaID},
			profile.Auth.OpenId,
			OnlineTime,
			Level,
			PlayerFriendsNum,
			profile.Auth.Extra.ClientVersion,
			profile.Auth.Extra.SystemSoftware,
			profile.Auth.Extra.SystemHardware,
			profile.Auth.Extra.TelecomOper,
			profile.Auth.Extra.Network,
			profile.Auth.Extra.ScreenWidth,
			profile.Auth.Extra.ScreenHight,
			profile.Auth.Extra.Density,
			profile.Auth.Extra.LoginChannel,
			profile.Auth.Extra.CpuHardware,
			profile.Auth.Extra.Memory,
			profile.Auth.Extra.GLRender,
			profile.Auth.Extra.GLVersion,
			profile.Auth.Extra.DeviceId,
			profile.Auth.Extra.VClientIP,
			format,
		)
	}*/

	logrus.WithFields(logrus.Fields{
		"FlowName":         "PlayerLogout",
		"GameSvrId":        gameSvrId,
		"EventTime":        dtEventTime,
		"GameAppId":        gameAppId,
		"PlatId":           GetPlatId(profile.Platform),
		"OpenId":           profile.Auth.OpenId,
		"RoleId":           profile.GlobalId,
		"OnlineTime":       OnlineTime,
		"PlayerFriendsNum": PlayerFriendsNum,
		"ClientVersion":    profile.Auth.ClientVersion,
		"LoginChannel":     profile.Auth.Channel,
		"DeviceId":         profile.Auth.DeviceId,
		"ClientIP":         profile.Auth.ClientIp,
		"Format":           format,
	}).Info("Tlog PlayerLogout")

	if nil != log.BILogger {
		log.BILogger.Printf("%s,%s,%s,%s,%d,%s,%d,%d,%d,%s,%s,%s,%s,%s",
			"PlayerLogout", gameSvrId, dtEventTime, gameAppId, GetPlatId(profile.Platform),
			profile.Auth.OpenId, profile.GlobalId, OnlineTime, PlayerFriendsNum, profile.Auth.ClientVersion,
			profile.Auth.Channel, profile.Auth.DeviceId, profile.Auth.ClientIp, format)
	} else {
		logrus.WithFields(logrus.Fields{"reason": "BILogger is nil"}).Info("LogPlayerLogout")
	}
	//{
	//"FlowName":         "PlayerLogout",
	//"GameSvrId":        gameSvrId,
	//"EventTime":        dtEventTime,
	//"GameAppId":        gameAppId,
	//"PlatId":           GetPlatId(profile.Platform),
	//"OpenId":           profile.Auth.OpenId,
	//"RoleId":           profile.GlobalId,
	//"OnlineTime":       OnlineTime,
	//"PlayerFriendsNum": PlayerFriendsNum,
	//"ClientVersion":    profile.Auth.ClientVersion,
	//"LoginChannel":     profile.Auth.Channel,
	//"DeviceId":         profile.Auth.DeviceId,
	//"ClientIP":         profile.Auth.ClientIp,
	//"Format":           format,
	//}).Info("Tlog PlayerLogout")
}

type TlogPlayerLogoutInfo struct {
	Timestamp int64  `json:"Timestamp"`
	Login     string `json:"LoginFormat"`
	Logout    string `json:"LogoutFormat"`
}

func LogLongTimeSession(record *model.TokenRecord) error {

	// process login and logout tlog message.
	last := storage.QueryActiveAccount(record)

	if last > 0 {

		return nil
	}

	now := time.Now().Unix()

	profile, err := storage.QueryProfile(record.GlobalId)
	if err != nil || profile == nil {

		return errors.New("query profile failed")
	}

	account := storage.QueryAccount(record.GlobalId)
	if account == nil {

		return errors.New("query account failed")
	}

	account.LoginTime = now
	profile.Timestamp = now

	// fixme gs send this tlog now
	//LogPlayerLogin(profile)

	storage.WriteProfile(profile)
	storage.Save(storage.AuthDatabase(), account)

	return nil
}

func GetPlatId(platform model.Platform) int {
	switch platform {
	case model.PlatformIOS:
		return 0
	case model.PlatformAndroid:
		return 1
	default:
		return 2
	}
}

func Platform(platId int) model.Platform {

	if platId == 0 {
		return model.PlatformIOS
	}

	if platId == 1 {
		return model.PlatformAndroid
	}

	return model.PlatformUndefined
}

func GetVendorConfigSetting(vendor model.Vendor) (string, string) {

	var gameAppID string
	gameSvrId := settings.GetString("tencent", "GameSvrId")

	switch vendor {

	case model.VendorMsdkGuest:

		gameAppID = settings.GetString("tencent", "msdk.Guest.AppId")
	case model.VendorMsdkQQ:

		gameAppID = settings.GetString("tencent", "msdk.QQ.AppId")
	case model.VendorMsdkWechat:

		gameAppID = settings.GetString("tencent", "msdk.Wechat.AppId")
	case model.VendorMobile:

		gameAppID = settings.GetString("sms", "qcloud.AppID")
	case model.VendorDevice:

		gameAppID = settings.GetString("device", "local.AppId")
	case model.VendorYsdkQQ:

		gameAppID = settings.GetString("tencent", "ysdk.YSDK_QQ.AppId")
	case model.VendorYsdkWechat:

		gameAppID = settings.GetString("tencent", "ysdk.YSDK_Wechat.AppId")
	case model.VendorWeibo:

		gameAppID = settings.GetString("tencent", "weibo.AppKey")
	case model.VendorVivo:

		gameAppID = settings.GetString("lovepay", "vivo.appId")
	case model.VendorBilibili:

		gameAppID = settings.GetString("lovepay", "bilibili.gameId")
	case model.VendorHuawei:
		gameAppID = settings.GetString("tencent", "huawei.appId")
	case model.VendorMgtv:
		gameAppID = settings.GetString("lovepay", "mgtv.appId")
	case model.VendorDouyin:
		gameAppID = settings.GetString("lovepay", "douyin.appId")
	default:

		logrus.WithFields(logrus.Fields{
			"vendor": vendor,
		}).Warn("GetVendorConfigSetting vendor setting not match.")
	}

	return gameSvrId, gameAppID
}
