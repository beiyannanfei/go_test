package query

import (
	"github.com/gin-gonic/gin"
	"github.com/xykong/loveauth/utils"
	"github.com/xykong/loveauth/errors"
	"github.com/xykong/loveauth/storage"
	"time"
	"net/http"
	"github.com/xykong/loveauth/storage/model"
	"github.com/xykong/loveauth/settings"
	"github.com/sirupsen/logrus"
	"strconv"
)

func init() {
	handlers["/grant"] = grant
}

type RequestGrant struct {
	Token string `form:"Token" json:"Token" binding:"required"`
}

type ResponseGrant struct {
	Body struct {
		Code     int64  `json:"code"`
		Message  string `json:"message"`
		GlobalId int64  `json:"globalId"`
	}
}

func grant(c *gin.Context) {
	var request RequestGrant

	// validation
	if err := c.BindJSON(&request); err != nil {
		utils.QuickReply(c, errors.ServerFail, "连接服务器异常")
		return
	}

	if request.Token == "" {
		utils.QuickReply(c, errors.ServerFail, "连接服务器异常")
		return
	}

	tokenRecord, err := storage.QueryAccessToken(request.Token)
	if err != nil {
		code := errors.Failed
		if ec, ok := err.(*errors.Type); ok {
			code = ec.Code
		}

		utils.QuickReply(c, code, "账号信息已过期，请重新登录")
		return
	}

	if tokenRecord.GlobalId == 0 {
		utils.QuickReply(c, errors.Failed, "用户并不存在")
		return
	}

	err = Check3rdToken(tokenRecord.GlobalId)
	if err != nil {
		code := errors.Failed
		if ec, ok := err.(*errors.Type); ok {
			code = ec.Code
		}

		utils.QuickReply(c, code, err.Error())
		return
	}

	resp := ResponseGrant{}
	resp.Body.Code = int64(errors.Ok)
	resp.Body.Message = "grant token successfully!"
	resp.Body.GlobalId = tokenRecord.GlobalId
	c.JSON(http.StatusOK, resp.Body)
	return
}

//检测第三方登陆用户token是否过期
func Check3rdToken(globalId int64) error {
	profile, err := storage.QueryProfile(globalId)
	if err != nil {
		return errors.NewCodeString(errors.Failed, "服务器连接错误，请稍后再试")
	}

	/*
		VendorDevice,
		VendorTencentQQ,
		VendorTencentWechat,
		VendorTencentGuest,
		VendorMobile,
		VendorWeibo,
		VendorYSDKQQ,
		VendorYSDKWechat,
	*/

	var dbGlobalId = globalId
	openId := profile.Auth.OpenId

	switch profile.Vendor {
	case model.VendorDevice:
	case model.VendorMsdkQQ:
	case model.VendorMsdkWechat:
	case model.VendorMsdkGuest:
	case model.VendorVivo:
	case model.VendorHuawei:

	case model.VendorMobile:
		var mobileInfo model.AuthMobile //验证是否解绑
		err = storage.QueryVendorByOpenId(openId, &mobileInfo)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err":     err,
				"profile": profile,
			}).Error("Check3rdToken QueryVendorByOpenId failed.")
			return err
		}

		dbGlobalId = mobileInfo.GlobalId

	case model.VendorWeibo:
		if time.Now().Unix()*1000 > profile.Auth.VendorWeibo.ExpirationAccess { //注意: 微博存储的为过期时间点的毫秒级时间戳
			return errors.NewCodeString(errors.ThirdAccessTokenExpire, "授权信息已过期，请重新登录")
		}

		var weiboInfo model.AuthWeibo //验证是否解绑
		err = storage.QueryVendorByOpenId(openId, &weiboInfo)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err":     err,
				"profile": profile,
			}).Error("Check3rdToken QueryVendorByOpenId failed.")
			return err
		}

		dbGlobalId = weiboInfo.GlobalId

	case model.VendorBilibili:
		bilibiliExpireTimes, _ := strconv.ParseInt(profile.Auth.VendorBilibili.LoginResult.ExpireTimes, 10, 64)
		if time.Now().Unix() > bilibiliExpireTimes {
			return errors.NewCodeString(errors.ThirdAccessTokenExpire, "授权信息已过期，请重新登录")
		}

	case model.VendorYsdkQQ: //ysdk_qq登陆形式
		if time.Now().Unix() > profile.Timestamp+profile.Auth.VendorYsdkQQ.ExpirationAccess {
			return errors.NewCodeString(errors.ThirdAccessTokenExpire, "授权信息已过期，请重新登录")
		}

		var ysdkQQInfo model.AuthYsdkQQ //验证是否解绑
		err = storage.QueryVendorByOpenId(openId, &ysdkQQInfo)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err":     err,
				"profile": profile,
			}).Error("Check3rdToken QueryVendorByOpenId failed.")
			return err
		}

		dbGlobalId = ysdkQQInfo.GlobalId

	case model.VendorYsdkWechat: //ysdk_wechat登陆形式,检测refresh_token是否过期
		RefreshTokenExpirationSeconds := settings.GetInt64("tencent", "ysdk.YSDK_Wechat.RefreshTokenExpirationSeconds")
		if time.Now().Unix() > profile.Timestamp+RefreshTokenExpirationSeconds {
			return errors.NewCodeString(errors.ThirdAccessTokenExpire, "授权信息已过期，请重新登录")
		}

		var ysdkWechatInfo model.AuthYsdkWechat //验证是否解绑
		err = storage.QueryVendorByOpenId(openId, &ysdkWechatInfo)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err":     err,
				"profile": profile,
			}).Error("Check3rdToken QueryVendorByOpenId failed.")
			return err
		}

		dbGlobalId = ysdkWechatInfo.GlobalId
	}

	if dbGlobalId != profile.GlobalId { //绑定关系不存在或已解除
		return errors.New("绑定关系已解除，请重新登录")
	}

	return nil
}
