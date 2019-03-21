package auth

import (
	"errors"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"github.com/xykong/tlog"
	"github.com/gin-gonic/gin"
	"github.com/xykong/loveauth/services/Login/quick_sdk"
)

type QuickOppoAuth struct {
	tlogger *tlog.Tlogger
}

func (r *QuickOppoAuth) Name() string {
	return "quick_oppo"
}

func (r *QuickOppoAuth) GetVendor() model.Vendor {
	return model.VendorQuickOppo
}

func (r *QuickOppoAuth) Start() error {

	var db = storage.AuthDatabase()
	if db == nil {
		return errors.New("quick require mysql database")
	}

	db.AutoMigrate(&model.AuthQuickOppo{})

	r.tlogger = NewTLogger(model.VendorQuickOppo)

	return nil
}

func (r *QuickOppoAuth) Verify(c *gin.Context, user *model.DoAuthRequest) error {
	// verify by msdk
	isVerify, err := quick_sdk.VerifyLoginQuick(user.VendorQuickOppo.Token, user.VendorQuickOppo.UserId)
	if err != nil {
		return err
	}

	if !isVerify {
		return errors.New("quick verify failed.")
	}

	return nil
}

func (r *QuickOppoAuth) FetchUser(c *gin.Context, user *model.DoAuthRequest) (*model.Account, error) {

	// check db, if not exist create new account.

	var authQuick model.AuthQuickOppo
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

func (r *QuickOppoAuth) CreateUser(c *gin.Context, user *model.DoAuthRequest) (*model.Account, error) {

	account := createAccount(user)
	err := storage.Insert(storage.AuthDatabase(), account)
	if err != nil {

		return nil, err
	}

	AuthQuick := &model.AuthQuickOppo{
		GlobalId:       account.GlobalId,
		OpenId:         user.OpenId,
		Platform:       user.Platform,
		UserId:         user.VendorQuickOppo.UserId,
		UserName:       user.VendorQuickOppo.UserName,
		Token:          user.VendorQuickOppo.Token,
		ChannelVersion: user.VendorQuickOppo.ChannelVersion,
		ChannelName:    user.VendorQuickOppo.ChannelName,
		ChannelType:    user.VendorQuickOppo.ChannelType,
	}

	err = storage.Insert(storage.AuthDatabase(), AuthQuick)
	if err != nil {
		return nil, err
	}

	LogPlayerRegister(user, account.GlobalId)

	return account, nil
}

func (r *QuickOppoAuth) UpdateVendor(globalId int64, user *model.DoAuthRequest) {
	var authQuick model.AuthQuickOppo
	data := map[string]interface{}{
		"user_id":         user.VendorQuickOppo.UserId,
		"user_name":       user.VendorQuickOppo.UserName,
		"token":           user.VendorQuickOppo.Token,
		"channel_version": user.VendorQuickOppo.ChannelVersion,
		"channel_name":    user.VendorQuickOppo.ChannelName,
		"channel_type":    user.VendorQuickOppo.ChannelType,
	}

	storage.UpdateSNSProfile(&authQuick, globalId, data)
}
