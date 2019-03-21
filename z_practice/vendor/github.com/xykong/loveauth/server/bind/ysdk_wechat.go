package bind

import (
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/xykong/loveauth/services/Login/ysdk"
)

type BindWechat struct {
}

func (b *BindWechat) Name() string {
	return "ysdk_wechat"
}

func (b *BindWechat) BindCheck(globalId int64, user *model.DoAuthRequest) error {
	//1 检测本登录方式是否已经绑定
	var wechatInfo model.AuthYsdkWechat
	err := storage.QueryVendorInfoByGlobalId(globalId, &wechatInfo)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":      err,
			"globalId": globalId,
			"user":     user,
		}).Error("BindCheck QueryVendorInfoByGlobalId failed.")
		return err
	}

	if wechatInfo.GlobalId != 0 { //已绑定
		logrus.WithFields(logrus.Fields{
			"wechatInfo": wechatInfo,
			"user":       user,
		}).Warn("BindCheck rebind.")
		return errors.NewCodeString(errors.RepeatBind, "登录方式已绑定，请解绑后再次绑定.")
	}

	//2 检测要绑定的账号是不是只有这一种登录方式
	var destWechat model.AuthYsdkWechat
	err = storage.QueryVendorByOpenId(user.OpenId, &destWechat)
	if err != nil && err != gorm.ErrRecordNotFound {
		logrus.WithFields(logrus.Fields{
			"err":      err,
			"globalId": globalId,
		}).Error("BindCheck QueryVendorByOpenId failed.")
		return err
	}

	if destWechat.GlobalId == 0 { //微信号没有登录过
		return nil
	}

	return errors.NewCodeString(errors.ForbidBindOther, "已绑定其他游戏账号，绑定失败")
}

func (b *BindWechat) BindVerify(c *gin.Context, user *model.DoAuthRequest) error {
	userIp := c.Request.RemoteAddr
	// verify by msdk
	_, err := ysdk.VerifyLoginWechat(user.OpenId, user.VendorYsdkWechat.TokenAccess, userIp)
	return err
}

func (b *BindWechat) CreateBind(globalId int64, user *model.DoAuthRequest) error {
	//1. 删除原有绑定信息
	var oriAuthWechat model.AuthYsdkWechat
	err := storage.DeleteVendorByOpenId4DB(user.OpenId, &oriAuthWechat)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":      err,
			"user":     user,
			"globalId": globalId,
		}).Error("CreateBind DeleteVendorByOpenId4DB failed.")
		return err
	}

	//2. 创建新的绑定信息
	authWechat := &model.AuthYsdkWechat{
		GlobalId:         globalId,
		OpenId:           user.OpenId,
		Platform:         user.Platform,
		TokenAccess:      user.VendorYsdkWechat.TokenAccess,
		ExpirationAccess: user.VendorYsdkWechat.ExpirationAccess,
		TokenRefresh:     user.VendorYsdkWechat.TokenRefresh,
		NickName:         user.VendorYsdkWechat.NickName,
		Picture:          user.VendorYsdkWechat.Picture,
		UnionId:          user.VendorYsdkWechat.UnionId,
	}

	err = storage.Insert(storage.AuthDatabase(), authWechat)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":        err,
			"user":       user,
			"globalId":   globalId,
			"authWechat": authWechat,
		}).Error("CreateBind Insert failed.")
		return err
	}

	return nil
}

func (b *BindWechat) UndoBind(globalId int64, req *RequestUnBind) error {
	//直接删除绑定关系
	var wechatInfo model.AuthYsdkWechat
	err := storage.DeleteVendorByGlobalId(globalId, &wechatInfo)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":      err,
			"globalId": globalId,
		}).Error("BindWechat UndoBind failed.")
		return err
	}

	return nil
}
