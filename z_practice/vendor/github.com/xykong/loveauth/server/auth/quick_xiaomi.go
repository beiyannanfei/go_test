package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/xykong/loveauth/services/Login/quick_sdk"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"github.com/xykong/tlog"
)

type QuickXiaomiAuth struct {
	tlogger *tlog.Tlogger
}

func (r *QuickXiaomiAuth) Name() string {
	return "quick_xiaomi"
}

func (r *QuickXiaomiAuth) GetVendor() model.Vendor {
	return model.VendorQuickXiaomi
}

func (r *QuickXiaomiAuth) Start() error {

	var db = storage.AuthDatabase()
	if db == nil {
		return errors.New("quick require mysql database")
	}

	db.AutoMigrate(&model.AuthQuickXiaomi{})

	r.tlogger = NewTLogger(model.VendorQuickXiaomi)

	return nil
}

func (r *QuickXiaomiAuth) Verify(c *gin.Context, user *model.DoAuthRequest) error {
	// verify by msdk
	isVerify, err := quick_sdk.VerifyLoginQuick(user.VendorQuickXiaomi.Token, user.VendorQuickXiaomi.UserId)
	if err != nil {
		return err
	}

	if !isVerify {
		return errors.New("quick verify failed.")
	}

	return nil
}

func (r *QuickXiaomiAuth) FetchUser(c *gin.Context, user *model.DoAuthRequest) (*model.Account, error) {

	// check db, if not exist create new account.

	var authQuick model.AuthQuickXiaomi
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

func (r *QuickXiaomiAuth) CreateUser(c *gin.Context, user *model.DoAuthRequest) (*model.Account, error) {

	account := createAccount(user)
	err := storage.Insert(storage.AuthDatabase(), account)
	if err != nil {

		return nil, err
	}

	AuthQuick := &model.AuthQuickXiaomi{
		GlobalId:       account.GlobalId,
		OpenId:         user.OpenId,
		Platform:       user.Platform,
		UserId:         user.VendorQuickXiaomi.UserId,
		UserName:       user.VendorQuickXiaomi.UserName,
		Token:          user.VendorQuickXiaomi.Token,
		ChannelVersion: user.VendorQuickXiaomi.ChannelVersion,
		ChannelName:    user.VendorQuickXiaomi.ChannelName,
		ChannelType:    user.VendorQuickXiaomi.ChannelType,
	}

	err = storage.Insert(storage.AuthDatabase(), AuthQuick)
	if err != nil {
		return nil, err
	}

	LogPlayerRegister(user, account.GlobalId)

	return account, nil
}

func (r *QuickXiaomiAuth) UpdateVendor(globalId int64, user *model.DoAuthRequest) {
	var authQuick model.AuthQuickXiaomi
	data := map[string]interface{}{
		"user_id":         user.VendorQuickXiaomi.UserId,
		"user_name":       user.VendorQuickXiaomi.UserName,
		"token":           user.VendorQuickXiaomi.Token,
		"channel_version": user.VendorQuickXiaomi.ChannelVersion,
		"channel_name":    user.VendorQuickXiaomi.ChannelName,
		"channel_type":    user.VendorQuickXiaomi.ChannelType,
	}

	storage.UpdateSNSProfile(&authQuick, globalId, data)
}
