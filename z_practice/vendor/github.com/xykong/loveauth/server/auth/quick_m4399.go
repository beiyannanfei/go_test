package auth

import (
	"errors"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"github.com/xykong/tlog"
	"github.com/gin-gonic/gin"
	"github.com/xykong/loveauth/services/Login/quick_sdk"
)

type QuickM4399Auth struct {
	tlogger *tlog.Tlogger
}

func (r *QuickM4399Auth) Name() string {
	return "quick_m4399"
}

func (r *QuickM4399Auth) GetVendor() model.Vendor {
	return model.VendorQuickM4399
}

func (r *QuickM4399Auth) Start() error {

	var db = storage.AuthDatabase()
	if db == nil {
		return errors.New("quick require mysql database")
	}

	db.AutoMigrate(&model.AuthQuickM4399{})

	r.tlogger = NewTLogger(model.VendorQuickM4399)

	return nil
}

func (r *QuickM4399Auth) Verify(c *gin.Context, user *model.DoAuthRequest) error {
	// verify by msdk
	isVerify, err := quick_sdk.VerifyLoginQuick(user.VendorQuickM4399.Token, user.VendorQuickM4399.UserId)
	if err != nil {
		return err
	}

	if !isVerify {
		return errors.New("quick verify failed.")
	}

	return nil
}

func (r *QuickM4399Auth) FetchUser(c *gin.Context, user *model.DoAuthRequest) (*model.Account, error) {

	// check db, if not exist create new account.

	var authQuick model.AuthQuickM4399
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

func (r *QuickM4399Auth) CreateUser(c *gin.Context, user *model.DoAuthRequest) (*model.Account, error) {

	account := createAccount(user)
	err := storage.Insert(storage.AuthDatabase(), account)
	if err != nil {

		return nil, err
	}

	AuthQuick := &model.AuthQuickM4399{
		GlobalId:       account.GlobalId,
		OpenId:         user.OpenId,
		Platform:       user.Platform,
		UserId:         user.VendorQuickM4399.UserId,
		UserName:       user.VendorQuickM4399.UserName,
		Token:          user.VendorQuickM4399.Token,
		ChannelVersion: user.VendorQuickM4399.ChannelVersion,
		ChannelName:    user.VendorQuickM4399.ChannelName,
		ChannelType:    user.VendorQuickM4399.ChannelType,
	}

	err = storage.Insert(storage.AuthDatabase(), AuthQuick)
	if err != nil {
		return nil, err
	}

	LogPlayerRegister(user, account.GlobalId)

	return account, nil
}

func (r *QuickM4399Auth) UpdateVendor(globalId int64, user *model.DoAuthRequest) {
	var authQuick model.AuthQuickM4399
	data := map[string]interface{}{
		"user_id":         user.VendorQuickM4399.UserId,
		"user_name":       user.VendorQuickM4399.UserName,
		"token":           user.VendorQuickM4399.Token,
		"channel_version": user.VendorQuickM4399.ChannelVersion,
		"channel_name":    user.VendorQuickM4399.ChannelName,
		"channel_type":    user.VendorQuickM4399.ChannelType,
	}

	storage.UpdateSNSProfile(&authQuick, globalId, data)
}
