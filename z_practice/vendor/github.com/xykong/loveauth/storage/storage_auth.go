package storage

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/storage/model"
	"time"
)

var dbAuth *gorm.DB

func AuthDatabase() *gorm.DB {
	return dbAuth
}

func InitAuthDatabase() {

	//Migrate the schema
	dbAuth.AutoMigrate(
		&model.Account{},
		&model.AuthPassword{},
		&model.AccountTag{},
		&model.AccountRealName{},
		&model.KuaikanCb{},
		&model.AccountAdInfo{},
		&model.GmAccount{},
	)
}

func IsValidPlatforms(platform model.Platform) bool {

	for _, item := range model.Platforms {

		if platform == item {
			return true
		}
	}

	return false
}

func IsValidVendors(vendor model.Vendor) bool {

	for _, item := range model.Vendors {

		if vendor == item {
			return true
		}
	}

	return false
}

func QueryAccount(globalId int64) *model.Account {

	var data model.Account

	if err := dbAuth.Where("global_id = ?", globalId).First(&data).Error; err != nil {
		logrus.WithFields(logrus.Fields{
			"globalId": globalId,
			"error":    err.Error(),
		}).Warn("QueryAccount data not found.")

		return nil
	}

	return &data
}

func QueryAccountRealName(globalId int64) *model.AccountRealName {
	var data model.AccountRealName

	if err := dbAuth.Where("global_id = ?", globalId).First(&data).Error; err != nil {
		logrus.WithFields(logrus.Fields{
			"globalId": globalId,
			"error":    err.Error(),
		}).Warn("QueryAccountRealName data not found.")

		return nil
	}

	return &data
}

func QueryAccountTag(globalId int64) *model.AccountTag {

	var data model.AccountTag

	if err := dbAuth.Where("global_id = ?", globalId).First(&data).Error; err != nil {
		logrus.WithFields(logrus.Fields{
			"globalId": globalId,
			"error":    err.Error(),
		}).Warn("QueryAccountTag data not found.")

		return nil
	}

	return &data
}

func QueryAccountByOpenId(openId string, vendor model.Vendor) *model.Account {

	var globalId int64 = 0
	switch vendor {
	case model.VendorDevice:

		var data model.AuthDevice
		if err := QueryVendorByOpenId(openId, &data); err != nil {
			logrus.WithFields(logrus.Fields{
				"openId": openId,
				"vendor": vendor,
				"error":  err.Error(),
			}).Error("QueryAccountByOpenId failed.")
		}
		globalId = data.GlobalId
		break

	case model.VendorMsdkQQ:

		var data model.AuthQQ
		if err := QueryVendorByOpenId(openId, &data); err != nil {
			logrus.WithFields(logrus.Fields{
				"openId": openId,
				"vendor": vendor,
				"error":  err.Error(),
			}).Error("QueryAccountByOpenId failed.")
		}
		globalId = data.GlobalId
		break

	case model.VendorYsdkQQ:

		var data model.AuthYsdkQQ
		if err := QueryVendorByOpenId(openId, &data); err != nil {
			logrus.WithFields(logrus.Fields{
				"openId": openId,
				"vendor": vendor,
				"error":  err.Error(),
			}).Error("QueryAccountByOpenId failed.")
		}
		globalId = data.GlobalId
		break

	case model.VendorMsdkWechat:

		var data model.AuthWechat
		if err := QueryVendorByOpenId(openId, &data); err != nil {
			logrus.WithFields(logrus.Fields{
				"openId": openId,
				"vendor": vendor,
				"error":  err.Error(),
			}).Error("QueryAccountByOpenId failed.")
		}
		globalId = data.GlobalId
		break

	case model.VendorYsdkWechat:

		var data model.AuthYsdkWechat
		if err := QueryVendorByOpenId(openId, &data); err != nil {
			logrus.WithFields(logrus.Fields{
				"openId": openId,
				"vendor": vendor,
				"error":  err.Error(),
			}).Error("QueryAccountByOpenId failed.")
		}
		globalId = data.GlobalId
		break

	case model.VendorMsdkGuest:

		var data model.AuthGuest
		if err := QueryVendorByOpenId(openId, &data); err != nil {
			logrus.WithFields(logrus.Fields{
				"openId": openId,
				"vendor": vendor,
				"error":  err.Error(),
			}).Error("QueryAccountByOpenId failed.")
		}
		globalId = data.GlobalId
		break

	case model.VendorWeibo:

		var data model.AuthWeibo
		if err := QueryVendorByOpenId(openId, &data); err != nil {
			logrus.WithFields(logrus.Fields{
				"openId": openId,
				"vendor": vendor,
				"error":  err.Error(),
			}).Error("QueryAccountByOpenId failed.")
		}
		globalId = data.GlobalId
		break

	case model.VendorMobile:
		var data model.AuthMobile
		if err := QueryVendorByOpenId(openId, &data); err != nil {
			logrus.WithFields(logrus.Fields{
				"openId": openId,
				"vendor": vendor,
				"error":  err.Error(),
			}).Error("QueryAccountByOpenId failed.")
		}
		globalId = data.GlobalId
		break

	case model.VendorVivo:
		var data model.AuthVivo
		if err := QueryVendorByOpenId(openId, &data); err != nil {
			logrus.WithFields(logrus.Fields{
				"openId": openId,
				"vendor": vendor,
				"error":  err.Error(),
			}).Error("QueryAccountByOpenId failed.")
		}
		globalId = data.GlobalId
		break
	case model.VendorBilibili:
		{
			var data model.AuthBilibili
			if err := QueryVendorByOpenId(openId, &data); err != nil {
				logrus.WithFields(logrus.Fields{
					"openId": openId,
					"vendor": vendor,
					"error":  err.Error(),
				}).Error("QueryAccountByOpenId bilibili failed.")
			}
			globalId = data.GlobalId
			break
		}
	case model.VendorHuawei:
		{
			var data model.AuthHuawei
			if err := QueryVendorByOpenId(openId, &data); err != nil {
				logrus.WithFields(logrus.Fields{
					"openId": openId,
					"vendor": vendor,
					"error":  err.Error(),
				}).Error("QueryAccountByOpenId huawei failed.")
			}
			globalId = data.GlobalId
			break
		}
	}

	if globalId == 0 {

		logrus.WithFields(logrus.Fields{
			"openId":   openId,
			"vendor":   vendor,
			"globalId": globalId,
		}).Error("QueryAccountByOpenId account not found.")

		return nil
	}

	var data model.Account
	if err := dbAuth.Where("global_id = ?", globalId).First(&data).Error; err != nil {
		logrus.WithFields(logrus.Fields{
			"globalId": globalId,
			"error":    err.Error(),
		}).Info("QueryAccount data not found.")

		return nil
	}

	return &data
}

func QueryVendorByOpenId(openId string, out interface{}) error {

	if dbAuth == nil {
		return nil
	}

	return dbAuth.Where("open_id = ?", openId).First(out).Error
}

func QueryVendorByOpenIds(openIds []string, out interface{}) error {

	if dbAuth == nil {
		return nil
	}

	return dbAuth.Where("open_id in (?)", openIds).Find(out).Error
}

func UpdateAccountState(globalId int64, data map[string]interface{}) error {

	if err := dbAuth.Table("accounts").Where("global_id = ?", globalId).Updates(data).Error; err != nil {

		logrus.WithFields(logrus.Fields{
			"globalId": globalId,
			"data":     data,
			"error":    err.Error(),
		}).Error("UpdateAccountState data not success.")

		return err
	}

	return nil
}

func UpdateSNSProfile(m interface{}, globalId int64, data map[string]interface{}) error {
	if err := dbAuth.Model(m).Where("global_id = ?", globalId).Updates(data).Error; err != nil {
		logrus.WithFields(logrus.Fields{
			"globalId": globalId,
			"data":     data,
			"err":      err,
		}).Error("UpdateSNSProfile data failed.")

		return err
	}

	return nil
}

func QueryBindInfo(globalId int64) (qq *model.AuthYsdkQQ, wechat *model.AuthYsdkWechat, mobile *model.AuthMobile, weibo *model.AuthWeibo, err error) {
	if globalId == 0 {
		return nil, nil, nil, nil, errors.New("globalId is invalid")
	}

	var qqInfo model.AuthYsdkQQ
	err = QueryVendorInfoByGlobalId(globalId, &qqInfo)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	qq = &qqInfo
	if qqInfo.GlobalId == 0 {
		qq = nil
	}

	var weChatInfo model.AuthYsdkWechat
	err = QueryVendorInfoByGlobalId(globalId, &weChatInfo)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	wechat = &weChatInfo
	if weChatInfo.GlobalId == 0 {
		wechat = nil
	}

	var mobileInfo model.AuthMobile
	err = QueryVendorInfoByGlobalId(globalId, &mobileInfo)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	mobile = &mobileInfo
	if mobileInfo.GlobalId == 0 {
		mobile = nil
	}

	var weiboInfo model.AuthWeibo
	err = QueryVendorInfoByGlobalId(globalId, &weiboInfo)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	weibo = &weiboInfo
	if weiboInfo.GlobalId == 0 {
		weibo = nil
	}

	return
}

func QueryVendorInfoByGlobalId(globalId int64, out interface{}) error {
	if dbAuth == nil {
		return nil
	}

	err := dbAuth.Where("global_id = ?", globalId).First(out).Error
	if err == gorm.ErrRecordNotFound {
		logrus.WithFields(logrus.Fields{
			"globalId": globalId,
			"err":      err,
		}).Info("QueryVendorInfoByGlobalId data null.")

		return nil
	}

	logrus.WithFields(logrus.Fields{
		"globalId": globalId,
		"err":      err,
	}).Error("QueryVendorInfoByGlobalId failed.")

	return err
}

//数据库直接删除(不是标记删除，防止数据冗余)
func DeleteVendorByOpenId4DB(openId string, out interface{}) error {
	if dbAuth == nil {
		return nil
	}

	return dbAuth.Unscoped().Where("open_id = ?", openId).Delete(out).Error
}

//数据库直接删除(不是标记删除，防止数据冗余)
func DeleteVendorByGlobalId(globalId int64, out interface{}) error {
	if dbAuth == nil {
		return nil
	}

	return dbAuth.Unscoped().Where("global_id = ?", globalId).Delete(out).Error
}

func CreateFirstIdfa(idfa string) error {
	kuaiKanCb := model.KuaikanCb{}
	kuaiKanCb.Idfa = idfa
	kuaiKanCb.CallBackTime = time.Now().Unix()

	return dbAuth.Where(model.KuaikanCb{Idfa: idfa}).FirstOrCreate(&kuaiKanCb).Error
}

func QueryGmAccountByAccount(account string) *model.GmAccount {

	var gmAccount model.GmAccount
	if err := dbAuth.Where("account = ?", account).First(&gmAccount).Error; err != nil {

		logrus.WithFields(logrus.Fields{
			"account": account,
			"error":   err.Error(),
		}).Info("QueryGmAccountByAccount data not found.")

		return nil
	}

	return &gmAccount
}
