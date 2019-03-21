package bind

import (
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/server/auth"
	"github.com/xykong/loveauth/settings"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"github.com/xykong/loveauth/utils/log"
	"time"
)

/**
    <struct name="BindFlow" desc="账号绑定转化(账号服)">
        <entry name="GameSvrId" type="string" size="25" desc="(必填)登录的游戏服务器编号"/>
        <entry name="TimeKey" type="datetime" desc="(必填)游戏事件的时间, 格式 YYYY-MM-DD HH:MM:SS"/>
        <entry name="GameAppId" type="string" size="32" desc="(必填)游戏APPID"/>
        <entry name="PlatId" type="int" defaultvalue="0" desc="(必填)ios 0/android 1"/>
        <entry name="ZoneAreaId" type="int" index="1" defaultvalue="0"
               desc="(必填)针对分区分服的游戏填写分区id，用来唯一标示一个区；非分区分服游戏请填写0"/>
        <entry name="LoginChannel" type="int" defaultvalue="0" desc="登录渠道"/>
        <entry name="RoleId" type="string" size="64" desc="玩家游戏内Id"/>
        <entry name="Format" type="string" size="64" defaultvalue="1.0.0" desc="Format version of this log"/>
        <entry name="Vendor" type="string" size="256" desc="第三方平台"/>
        <entry name="VendorOpenId" type="string" size="256" desc="第三方平台提供的唯一Id"/>
        <entry name="BindType" type="int" desc="绑定操作类型"/>
        <entry name="BindCount" type="int" desc="当前绑定第三方账号数量"/>
    </struct>
 */

func LogBind(globalId int64, requestInfo model.DoAuthRequest, bindVendor string, bindType int) {
	eventTime := time.Now().Format("2006-01-02 15:04:05")
	format := settings.GetString("tencent", "tlog.format")
	bindCount := 0
	qq, wechat, mobile, weibo, err := storage.QueryBindInfo(globalId)
	if nil == err {
		if nil != qq && qq.GlobalId != 0 {
			bindCount += 1
		}
		if nil != wechat && wechat.GlobalId != 0 {
			bindCount += 1
		}
		if nil != mobile && mobile.GlobalId != 0 {
			bindCount += 1
		}
		if nil != weibo && weibo.GlobalId != 0 {
			bindCount += 1
		}
	}

	gameSvrId, gameAppId := auth.GetVendorConfigSetting(requestInfo.Vendor)

	logrus.WithFields(logrus.Fields{
		"FlowName":     "BindFlow",
		"GameSvrId":    gameSvrId,
		"EventTime":    eventTime,
		"GameAppId":    gameAppId,
		"PlatId":       auth.GetPlatId(requestInfo.Platform),
		"LoginChannel": requestInfo.Channel,
		"RoleId":       globalId,
		"Format":       format,
		"LoginVendor":  requestInfo.Vendor,
		"BindVendor":   bindVendor,
		"VendorOpenId": requestInfo.OpenId,
		"BindType":     bindType,
		"BindCount":    bindCount,
	}).Info("Tlog BindFlow")

	if nil != log.BILogger {
		log.BILogger.Printf("%s,%s,%s,%s,%d,%s,%d,%s,%s,%s,%s,%d,%d",
			"BindFlow", gameSvrId, eventTime, gameAppId, auth.GetPlatId(requestInfo.Platform),
			requestInfo.Channel, globalId, format, requestInfo.Vendor, bindVendor, requestInfo.OpenId,
			bindType, bindCount)
	} else {
		logrus.WithFields(logrus.Fields{"reason": "BILogger is nil"}).Info("LogBindFlow")
	}
}

const (
	Unbind = 0
	Bind   = 1
)
