package auth

import (
	"errors"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"github.com/xykong/tlog"
	"github.com/gin-gonic/gin"
	"github.com/xykong/loveauth/services/Login/quick_sdk"
)

type QuickMeizuAuth struct {
	tlogger *tlog.Tlogger
}

func (r *QuickMeizuAuth) Name() string {
	return "quick_meizu"
}

func (r *QuickMeizuAuth) GetVendor() model.Vendor {
	return model.VendorQuickMeiZu
}

func (r *QuickMeizuAuth) Start() error {

	var db = storage.AuthDatabase()
	if db == nil {
		return errors.New("quick require mysql database")
	}

	db.AutoMigrate(&model.AuthQuickMeizu{})

	r.tlogger = NewTLogger(model.VendorQuickMeiZu)

	return nil
}

func (r *QuickMeizuAuth) Verify(c *gin.Context, user *model.DoAuthRequest) error {
	// verify by msdk
	isVerify, err := quick_sdk.VerifyLoginQuick(user.VendorQuickMeiZu.Token, user.VendorQuickMeiZu.UserId)
	if err != nil {
		return err
	}

	if !isVerify {
		return errors.New("quick verify failed.")
	}

	return nil
}

func (r *QuickMeizuAuth) FetchUser(c *gin.Context, user *model.DoAuthRequest) (*model.Account, error) {

	// check db, if not exist create new account.

	var authQuick model.AuthQuickMeizu
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

func (r *QuickMeizuAuth) CreateUser(c *gin.Context, user *model.DoAuthRequest) (*model.Account, error) {

	account := createAccount(user)
	err := storage.Insert(storage.AuthDatabase(), account)
	if err != nil {

		return nil, err
	}

	AuthQuick := &model.AuthQuickMeizu{
		GlobalId:       account.GlobalId,
		OpenId:         user.OpenId,
		Platform:       user.Platform,
		UserId:         user.VendorQuickMeiZu.UserId,
		UserName:       user.VendorQuickMeiZu.UserName,
		Token:          user.VendorQuickMeiZu.Token,
		ChannelVersion: user.VendorQuickMeiZu.ChannelVersion,
		ChannelName:    user.VendorQuickMeiZu.ChannelName,
		ChannelType:    user.VendorQuickMeiZu.ChannelType,
	}

	err = storage.Insert(storage.AuthDatabase(), AuthQuick)
	if err != nil {
		return nil, err
	}

	LogPlayerRegister(user, account.GlobalId)

	return account, nil
}

func (r *QuickMeizuAuth) UpdateVendor(globalId int64, user *model.DoAuthRequest) {
	var authQuick model.AuthQuickMeizu
	data := map[string]interface{}{
		"user_id":         user.VendorQuickMeiZu.UserId,
		"user_name":       user.VendorQuickMeiZu.UserName,
		"token":           user.VendorQuickMeiZu.Token,
		"channel_version": user.VendorQuickMeiZu.ChannelVersion,
		"channel_name":    user.VendorQuickMeiZu.ChannelName,
		"channel_type":    user.VendorQuickMeiZu.ChannelType,
	}

	storage.UpdateSNSProfile(&authQuick, globalId, data)
}
