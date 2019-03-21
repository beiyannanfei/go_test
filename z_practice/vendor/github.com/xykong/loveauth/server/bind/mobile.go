package bind

import (
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/errors"
	"github.com/gin-gonic/gin"
	"github.com/xykong/loveauth/services/sms"
	"github.com/jinzhu/gorm"
)

type BindMobile struct {
}

func (b *BindMobile) Name() string {
	return "mobile"
}

func (b *BindMobile) BindCheck(globalId int64, user *model.DoAuthRequest) error {
	//1 检测本登录方式是否已经绑定
	var mobileInfo model.AuthMobile
	err := storage.QueryVendorInfoByGlobalId(globalId, &mobileInfo)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":      err,
			"globalId": globalId,
			"user":     user,
		}).Error("BindCheck QueryVendorInfoByGlobalId failed.")
		return err
	}

	if mobileInfo.GlobalId != 0 { //已绑定
		logrus.WithFields(logrus.Fields{
			"mobileInfo": mobileInfo,
			"user":       user,
		}).Warn("BindCheck rebind.")
		return errors.NewCodeString(errors.RepeatBind, "登录方式已绑定，请解绑后再次绑定.")
	}

	//2 检测要绑定的账号是不是只有这一种登录方式
	var destMobile model.AuthMobile
	err = storage.QueryVendorByOpenId(user.OpenId, &destMobile)
	if err != nil && err != gorm.ErrRecordNotFound {
		logrus.WithFields(logrus.Fields{
			"err":      err,
			"globalId": globalId,
		}).Error("BindCheck QueryVendorByOpenId failed.")
		return err
	}

	if destMobile.GlobalId == 0 { //手机号没有登录过
		return nil
	}

	return errors.NewCodeString(errors.ForbidBindOther, "已绑定其他游戏账号，绑定失败")
}

func (b *BindMobile) BindVerify(c *gin.Context, user *model.DoAuthRequest) error {
	return sms.VerifyMobileCode(user.OpenId, user.VendorMobile.VerifyCode, true)
}

func (b *BindMobile) CreateBind(globalId int64, user *model.DoAuthRequest) error {
	//1. 删除原有绑定信息
	var oriAuthMobile model.AuthMobile
	err := storage.DeleteVendorByOpenId4DB(user.OpenId, &oriAuthMobile)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":      err,
			"user":     user,
			"globalId": globalId,
		}).Error("CreateBind DeleteVendorByOpenId4DB failed.")
		return err
	}

	//2. 创建新的绑定信息
	authMobile := &model.AuthMobile{
		GlobalId: globalId,
		OpenId:   user.OpenId,
		Platform: user.Platform,
	}
	err = storage.Insert(storage.AuthDatabase(), authMobile)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":        err,
			"user":       user,
			"globalId":   globalId,
			"authMobile": authMobile,
		}).Error("CreateBind Insert failed.")
		return err
	}

	return nil
}

func (b *BindMobile) UndoBind(globalId int64, req *RequestUnBind) error {
	//更换步骤
	//1 检测新绑定手机号是否绑定其它账号
	var mobileInfo model.AuthMobile
	storage.QueryVendorByOpenId(req.OpenId, &mobileInfo)
	if mobileInfo.GlobalId != 0 { //存在绑定信息
		logrus.WithFields(logrus.Fields{
			"mobileInfo": mobileInfo,
			"globalId":   globalId,
			"req":        req,
		}).Error("BindMobile UndoBind QueryVendorByOpenId failed.")
		return errors.NewCodeString(errors.BindOtherAccount, "已绑定其他游戏账号，绑定失败")
	}

	//2 验证手机验证码
	err := sms.VerifyMobileCode(req.OpenId, req.AccessTokenValue, true)
	if err != nil {
		return err
	}

	//3 更换绑定关系
	var newMobileInfo model.AuthMobile
	err = storage.UpdateSNSProfile(&newMobileInfo, globalId, map[string]interface{}{"open_id": req.OpenId})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"globalId": globalId,
			"req":      req,
		}).Error("BindMobile UndoBind UpdateSNSProfile failed.")
		return err
	}

	return nil
}
