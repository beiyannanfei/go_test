package auth

import (
	"errors"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"github.com/xykong/tlog"
	"github.com/gin-gonic/gin"
	"github.com/xykong/loveauth/services/Login/quick_sdk"
)

type QuickYsdkAuth struct {
	tlogger *tlog.Tlogger
}

func (r *QuickYsdkAuth) Name() string {
	return "quick_ysdk"
}

func (r *QuickYsdkAuth) GetVendor() model.Vendor {
	return model.VendorQuickYsdk
}

func (r *QuickYsdkAuth) Start() error {

	var db = storage.AuthDatabase()
	if db == nil {
		return errors.New("quick require mysql database")
	}

	db.AutoMigrate(&model.AuthQuickYsdk{})

	r.tlogger = NewTLogger(model.VendorQuickYsdk)

	return nil
}

func (r *QuickYsdkAuth) Verify(c *gin.Context, user *model.DoAuthRequest) error {
	// verify by msdk
	isVerify, err := quick_sdk.VerifyLoginQuick(user.VendorQuickYsdk.Token, user.VendorQuickYsdk.UserId)
	if err != nil {
		return err
	}

	if !isVerify {
		return errors.New("quick verify failed.")
	}

	return nil
}

func (r *QuickYsdkAuth) FetchUser(c *gin.Context, user *model.DoAuthRequest) (*model.Account, error) {

	// check db, if not exist create new account.

	var authQuick model.AuthQuickYsdk
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

func (r *QuickYsdkAuth) CreateUser(c *gin.Context, user *model.DoAuthRequest) (*model.Account, error) {

	account := createAccount(user)
	err := storage.Insert(storage.AuthDatabase(), account)
	if err != nil {

		return nil, err
	}

	AuthQuick := &model.AuthQuickYsdk{
		GlobalId:       account.GlobalId,
		OpenId:         user.OpenId,
		Platform:       user.Platform,
		UserId:         user.VendorQuickYsdk.UserId,
		UserName:       user.VendorQuickYsdk.UserName,
		Token:          user.VendorQuickYsdk.Token,
		ChannelVersion: user.VendorQuickYsdk.ChannelVersion,
		ChannelName:    user.VendorQuickYsdk.ChannelName,
		ChannelType:    user.VendorQuickYsdk.ChannelType,
	}

	err = storage.Insert(storage.AuthDatabase(), AuthQuick)
	if err != nil {
		return nil, err
	}

	LogPlayerRegister(user, account.GlobalId)

	return account, nil
}

func (r *QuickYsdkAuth) UpdateVendor(globalId int64, user *model.DoAuthRequest) {
	var authQuick model.AuthQuickYsdk
	data := map[string]interface{}{
		"user_id":         user.VendorQuickYsdk.UserId,
		"user_name":       user.VendorQuickYsdk.UserName,
		"token":           user.VendorQuickYsdk.Token,
		"channel_version": user.VendorQuickYsdk.ChannelVersion,
		"channel_name":    user.VendorQuickYsdk.ChannelName,
		"channel_type":    user.VendorQuickYsdk.ChannelType,
	}

	storage.UpdateSNSProfile(&authQuick, globalId, data)
}
