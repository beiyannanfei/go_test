package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/errors"
	"github.com/xykong/loveauth/services/Login/huawei"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"github.com/xykong/tlog"
	"strconv"
)

type HuaweiAuth struct {
	tlogger *tlog.Tlogger
}

func (r *HuaweiAuth) Name() string {
	return "huawei"
}

func (r *HuaweiAuth) GetVendor() model.Vendor {
	return model.VendorHuawei
}

func (r *HuaweiAuth) Start() error {
	var db = storage.AuthDatabase()
	if nil == db {
		return errors.New("huawei require mysql database")
	}

	db.AutoMigrate(&model.AuthHuawei{})
	r.tlogger = NewTLogger(model.VendorHuawei)
	return nil
}

// 需要playerLevel, playerSign
func (r *HuaweiAuth) Verify(c *gin.Context, user *model.DoAuthRequest) error {
	resp, err := huawei.AuthToken(user.OpenId, strconv.Itoa(user.VendorHuawei.GameUserData.PlayerLevel), user.VendorHuawei.GameUserData.GameAuthSign, user.VendorHuawei.GameUserData.Ts)
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err":  err,
			"user": user,
		}).Error("HuaweiAuth Verify err")

		return err
	}

	if 0 != resp.RtnCode {
		logrus.WithFields(logrus.Fields{
			"resp": resp,
			"user": user,
		}).Error("HuaweiAuth Verify failed")

		return errors.New("with error code " + strconv.Itoa(resp.RtnCode))
	}

	return nil
}

func (r *HuaweiAuth) FetchUser(c *gin.Context, user *model.DoAuthRequest) (*model.Account, error) {
	var authHuawei model.AuthHuawei
	err := storage.QueryVendorByOpenId(user.OpenId, &authHuawei)
	if nil != err {

		return nil, nil
	}

	account := storage.QueryAccount(authHuawei.GlobalId)

	if nil == account {
		return nil, errors.New("account not found in database for huawei")
	}
	return account, nil
}

func (r *HuaweiAuth) CreateUser(c *gin.Context, user *model.DoAuthRequest) (*model.Account, error) {
	account := createAccount(user)
	err := storage.Insert(storage.AuthDatabase(), account)
	if nil != err {
		return nil, err
	}

	authHuawei := &model.AuthHuawei{
		GlobalId:     account.GlobalId,
		OpenId:       user.OpenId,
		Platform:     user.Platform,
		PlayerId:     user.VendorHuawei.GameUserData.PlayerId,
		DisplayName:  user.VendorHuawei.GameUserData.DisplayName,
		PlayerLevel:  user.VendorHuawei.GameUserData.PlayerLevel,
		IsAuth:       user.VendorHuawei.GameUserData.IsAuth,
		Ts:           user.VendorHuawei.GameUserData.Ts,
		GameAuthSign: user.VendorHuawei.GameUserData.GameAuthSign,
	}

	err = storage.Insert(storage.AuthDatabase(), authHuawei)
	if nil != err {
		return nil, err
	}

	LogPlayerRegister(user, account.GlobalId)
	return account, nil
}

func (r *HuaweiAuth) UpdateVendor(globalId int64, user *model.DoAuthRequest) {
	var authHuawei model.AuthHuawei
	data := map[string]interface{}{
		"player_id":      user.VendorHuawei.GameUserData.PlayerId,
		"display_name":   user.VendorHuawei.GameUserData.DisplayName,
		"player_level":   user.VendorHuawei.GameUserData.PlayerLevel,
		"is_auth":        user.VendorHuawei.GameUserData.IsAuth,
		"ts":             user.VendorHuawei.GameUserData.Ts,
		"game_auth_sign": user.VendorHuawei.GameUserData.GameAuthSign,
	}

	storage.UpdateSNSProfile(&authHuawei, globalId, data)
}
