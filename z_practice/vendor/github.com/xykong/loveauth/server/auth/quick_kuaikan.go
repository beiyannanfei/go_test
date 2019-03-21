package auth

import (
	"errors"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"github.com/xykong/tlog"
	"github.com/gin-gonic/gin"
	"github.com/xykong/loveauth/services/Login/quick_sdk"
)

type QuickKuaikanAuth struct {
	tlogger *tlog.Tlogger
}

func (r *QuickKuaikanAuth) Name() string {
	return "quick_kuaikan"
}

func (r *QuickKuaikanAuth) GetVendor() model.Vendor {
	return model.VendorQuickKuaiKan
}

func (r *QuickKuaikanAuth) Start() error {

	var db = storage.AuthDatabase()
	if db == nil {
		return errors.New("quick require mysql database")
	}

	db.AutoMigrate(&model.AuthQuickKuaikan{})

	r.tlogger = NewTLogger(model.VendorQuickKuaiKan)

	return nil
}

func (r *QuickKuaikanAuth) Verify(c *gin.Context, user *model.DoAuthRequest) error {
	// verify by msdk
	isVerify, err := quick_sdk.VerifyLoginQuick(user.VendorQuickKuaiKan.Token, user.VendorQuickKuaiKan.UserId)
	if err != nil {
		return err
	}

	if !isVerify {
		return errors.New("quick verify failed.")
	}

	return nil
}

func (r *QuickKuaikanAuth) FetchUser(c *gin.Context, user *model.DoAuthRequest) (*model.Account, error) {

	// check db, if not exist create new account.

	var authQuick model.AuthQuickKuaikan
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

func (r *QuickKuaikanAuth) CreateUser(c *gin.Context, user *model.DoAuthRequest) (*model.Account, error) {

	account := createAccount(user)
	err := storage.Insert(storage.AuthDatabase(), account)
	if err != nil {

		return nil, err
	}

	AuthQuick := &model.AuthQuickKuaikan{
		GlobalId:       account.GlobalId,
		OpenId:         user.OpenId,
		Platform:       user.Platform,
		UserId:         user.VendorQuickKuaiKan.UserId,
		UserName:       user.VendorQuickKuaiKan.UserName,
		Token:          user.VendorQuickKuaiKan.Token,
		ChannelVersion: user.VendorQuickKuaiKan.ChannelVersion,
		ChannelName:    user.VendorQuickKuaiKan.ChannelName,
		ChannelType:    user.VendorQuickKuaiKan.ChannelType,
	}

	err = storage.Insert(storage.AuthDatabase(), AuthQuick)
	if err != nil {
		return nil, err
	}

	LogPlayerRegister(user, account.GlobalId)

	return account, nil
}

func (r *QuickKuaikanAuth) UpdateVendor(globalId int64, user *model.DoAuthRequest) {
	var authQuick model.AuthQuickKuaikan
	data := map[string]interface{}{
		"user_id":         user.VendorQuickKuaiKan.UserId,
		"user_name":       user.VendorQuickKuaiKan.UserName,
		"token":           user.VendorQuickKuaiKan.Token,
		"channel_version": user.VendorQuickKuaiKan.ChannelVersion,
		"channel_name":    user.VendorQuickKuaiKan.ChannelName,
		"channel_type":    user.VendorQuickKuaiKan.ChannelType,
	}

	storage.UpdateSNSProfile(&authQuick, globalId, data)
}
