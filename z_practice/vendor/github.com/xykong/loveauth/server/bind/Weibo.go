package bind

import (
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"time"
	"github.com/xykong/loveauth/services/Login/webo_sdk"
)

type BindWeibo struct {
}

func (b *BindWeibo) Name() string {
	return "weibo"
}

func (b *BindWeibo) BindCheck(globalId int64, user *model.DoAuthRequest) error {
	//1 检测本登录方式是否已经绑定
	var weiboInfo model.AuthWeibo
	err := storage.QueryVendorInfoByGlobalId(globalId, &weiboInfo)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":      err,
			"globalId": globalId,
			"user":     user,
		}).Error("BindWeibo BindCheck QueryVendorInfoByGlobalId failed.")
		return err
	}

	if weiboInfo.GlobalId != 0 { //已绑定
		logrus.WithFields(logrus.Fields{
			"weiboInfo": weiboInfo,
			"user":      user,
		}).Warn("BindCheck rebind.")
		return errors.NewCodeString(errors.RepeatBind, "登录方式已绑定，请解绑后再次绑定.")
	}

	//2 检测要绑定的账号是不是只有这一种登录方式
	var destWeibo model.AuthWeibo
	err = storage.QueryVendorByOpenId(user.OpenId, &destWeibo)
	if err != nil && err != gorm.ErrRecordNotFound {
		logrus.WithFields(logrus.Fields{
			"err":      err,
			"globalId": globalId,
		}).Error("BindCheck QueryVendorByOpenId failed.")
		return err
	}

	if destWeibo.GlobalId == 0 { //微信号没有登录过
		return nil
	}

	return errors.NewCodeString(errors.ForbidBindOther, "已绑定其他游戏账号，绑定失败")
}

func (b *BindWeibo) BindVerify(c *gin.Context, user *model.DoAuthRequest) error {
	// verify by msdk
	_, err := webo_sdk.WeiBoGetTokenInfo(user.VendorWeibo.TokenAccess)
	if err != nil {
		return err
	}

	//注意：微博的过期时间存在差异(android: 剩余到期的秒数，ios: 到期时间的毫秒级时间戳)
	if user.Platform == model.PlatformAndroid { //android手动转换为ios的格式
		user.VendorWeibo.ExpirationAccess = (time.Now().Unix() + user.VendorWeibo.ExpirationAccess) * 1000
	}

	return nil
}

func (b *BindWeibo) CreateBind(globalId int64, user *model.DoAuthRequest) error {
	//1. 删除原有绑定信息
	var oriAuthWeibo model.AuthWeibo
	err := storage.DeleteVendorByOpenId4DB(user.OpenId, &oriAuthWeibo)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":      err,
			"user":     user,
			"globalId": globalId,
		}).Error("CreateBind DeleteVendorByOpenId4DB failed.")
		return err
	}

	//2. 创建新的绑定信息
	authWeibo := &model.AuthWeibo{
		GlobalId:         globalId,
		OpenId:           user.OpenId,
		TokenAccess:      user.VendorWeibo.TokenAccess,
		ExpirationAccess: user.VendorWeibo.ExpirationAccess,
		NickName:         user.VendorWeibo.NickName,
		Picture:          user.VendorWeibo.Picture,
	}

	err = storage.Insert(storage.AuthDatabase(), authWeibo)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":       err,
			"user":      user,
			"globalId":  globalId,
			"authWeibo": authWeibo,
		}).Error("CreateBind Insert failed.")
		return err
	}

	return nil
}

func (b *BindWeibo) UndoBind(globalId int64, req *RequestUnBind) error {
	//直接删除绑定关系
	var weiboInfo model.AuthWeibo
	err := storage.DeleteVendorByGlobalId(globalId, &weiboInfo)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":      err,
			"globalId": globalId,
		}).Error("BindWeibo UndoBind failed.")
		return err
	}

	return nil
}
