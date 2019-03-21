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

type BindQQ struct {
}

func (b *BindQQ) Name() string {
	return "ysdk_qq"
}

func (b *BindQQ) BindCheck(globalId int64, user *model.DoAuthRequest) error {
	//1 检测本登录方式是否已经绑定
	var qqInfo model.AuthYsdkQQ
	err := storage.QueryVendorInfoByGlobalId(globalId, &qqInfo)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":      err,
			"globalId": globalId,
			"user":     user,
		}).Error("BindCheck QueryVendorInfoByGlobalId failed.")
		return err
	}

	if qqInfo.GlobalId != 0 { //已绑定
		logrus.WithFields(logrus.Fields{
			"qqInfo": qqInfo,
			"user":   user,
		}).Warn("BindCheck rebind.")
		return errors.NewCodeString(errors.RepeatBind, "登录方式已绑定，请解绑后再次绑定.")
	}

	//2 检测要绑定的账号是不是只有这一种登录方式
	var destQQ model.AuthYsdkQQ
	err = storage.QueryVendorByOpenId(user.OpenId, &destQQ)
	if err != nil && err != gorm.ErrRecordNotFound {
		logrus.WithFields(logrus.Fields{
			"err":      err,
			"globalId": globalId,
		}).Error("BindCheck QueryVendorByOpenId failed.")
		return err
	}

	if destQQ.GlobalId == 0 { //qq号没有登录过
		return nil
	}

	return errors.NewCodeString(errors.ForbidBindOther, "已绑定其他游戏账号，绑定失败")
}

func (b *BindQQ) BindVerify(c *gin.Context, user *model.DoAuthRequest) error {
	userIp := c.Request.RemoteAddr
	// verify by ysdk
	_, err := ysdk.VerifyLoginQQ(user.OpenId, user.VendorYsdkQQ.TokenAccess, userIp)
	return err
}

func (b *BindQQ) CreateBind(globalId int64, user *model.DoAuthRequest) error {
	//1. 删除原有绑定信息
	var oriAuthQQ model.AuthYsdkQQ
	err := storage.DeleteVendorByOpenId4DB(user.OpenId, &oriAuthQQ)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":      err,
			"user":     user,
			"globalId": globalId,
		}).Error("CreateBind DeleteVendorByOpenId4DB failed.")
		return err
	}

	//2. 创建新的绑定信息
	authYsdkQQ := &model.AuthYsdkQQ{
		GlobalId:         globalId,
		OpenId:           user.OpenId,
		Platform:         user.Platform,
		TokenAccess:      user.VendorYsdkQQ.TokenAccess,
		ExpirationAccess: user.VendorYsdkQQ.ExpirationAccess,
		TokenPay:         user.VendorYsdkQQ.TokenPay,
		Pf:               user.VendorYsdkQQ.Pf,
		PfKey:            user.VendorYsdkQQ.PfKey,
		NickName:         user.VendorYsdkQQ.NickName,
		Picture:          user.VendorYsdkQQ.Picture,
	}

	err = storage.Insert(storage.AuthDatabase(), authYsdkQQ)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":        err,
			"user":       user,
			"globalId":   globalId,
			"authYsdkQQ": authYsdkQQ,
		}).Error("CreateBind Insert failed.")
		return err
	}

	return nil
}

func (b *BindQQ) UndoBind(globalId int64, req *RequestUnBind) error {
	//直接删除绑定关系
	var ysdkQQInfo model.AuthYsdkQQ
	err := storage.DeleteVendorByGlobalId(globalId, &ysdkQQInfo)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":      err,
			"globalId": globalId,
		}).Error("BindQQ UndoBind failed.")
		return err
	}

	return nil
}
