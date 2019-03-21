package auth

import (
	"errors"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"github.com/xykong/tlog"
	"github.com/gin-gonic/gin"
	"github.com/xykong/loveauth/services/Login/quick_sdk"
)

type QuickIqiyiAuth struct {
	tlogger *tlog.Tlogger
}

func (r *QuickIqiyiAuth) Name() string {
	return "quick_iqiyi"
}

func (r *QuickIqiyiAuth) GetVendor() model.Vendor {
	return model.VendorQuickIqiyi
}

func (r *QuickIqiyiAuth) Start() error {

	var db = storage.AuthDatabase()
	if db == nil {
		return errors.New("quick require mysql database")
	}

	db.AutoMigrate(&model.AuthQuickIqiyi{})

	r.tlogger = NewTLogger(model.VendorQuickIqiyi)

	return nil
}

func (r *QuickIqiyiAuth) Verify(c *gin.Context, user *model.DoAuthRequest) error {
	// verify by msdk
	isVerify, err := quick_sdk.VerifyLoginQuick(user.VendorQuickIqiyi.Token, user.VendorQuickIqiyi.UserId)
	if err != nil {
		return err
	}

	if !isVerify {
		return errors.New("quick verify failed.")
	}

	return nil
}

func (r *QuickIqiyiAuth) FetchUser(c *gin.Context, user *model.DoAuthRequest) (*model.Account, error) {

	// check db, if not exist create new account.

	var authQuick model.AuthQuickIqiyi
	err := storage.QueryVendorByOpenId(user.OpenId, &authQuick)
	if err != nil {
		// no 3rdParty account found, need create new one.
		return nil, nil
	}

	account := storage.QueryAccount(authQuick.GlobalId)

	// some thing wrong like: account removed from table.
	if account == nil {
		return nil, errors.New("account not found in database")
	}

	return account, nil
}

func (r *QuickIqiyiAuth) CreateUser(c *gin.Context, user *model.DoAuthRequest) (*model.Account, error) {

	account := createAccount(user)
	err := storage.Insert(storage.AuthDatabase(), account)
	if err != nil {

		return nil, err
	}

	AuthQuick := &model.AuthQuickIqiyi{
		GlobalId:       account.GlobalId,
		OpenId:         user.OpenId,
		Platform:       user.Platform,
		UserId:         user.VendorQuickIqiyi.UserId,
		UserName:       user.VendorQuickIqiyi.UserName,
		Token:          user.VendorQuickIqiyi.Token,
		ChannelVersion: user.VendorQuickIqiyi.ChannelVersion,
		ChannelName:    user.VendorQuickIqiyi.ChannelName,
		ChannelType:    user.VendorQuickIqiyi.ChannelType,
	}

	err = storage.Insert(storage.AuthDatabase(), AuthQuick)
	if err != nil {
		return nil, err
	}

	LogPlayerRegister(user, account.GlobalId)

	return account, nil
}

func (r *QuickIqiyiAuth) UpdateVendor(globalId int64, user *model.DoAuthRequest) {
	var authQuick model.AuthQuickIqiyi
	data := map[string]interface{}{
		"user_id":         user.VendorQuickIqiyi.UserId,
		"user_name":       user.VendorQuickIqiyi.UserName,
		"token":           user.VendorQuickIqiyi.Token,
		"channel_version": user.VendorQuickIqiyi.ChannelVersion,
		"channel_name":    user.VendorQuickIqiyi.ChannelName,
		"channel_type":    user.VendorQuickIqiyi.ChannelType,
	}

	storage.UpdateSNSProfile(&authQuick, globalId, data)
}
