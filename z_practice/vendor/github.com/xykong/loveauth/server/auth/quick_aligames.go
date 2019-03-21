package auth

import (
	"errors"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"github.com/xykong/tlog"
	"github.com/gin-gonic/gin"
	"github.com/xykong/loveauth/services/Login/quick_sdk"
)

type QuickAliGamesAuth struct {
	tlogger *tlog.Tlogger
}

func (r *QuickAliGamesAuth) Name() string {
	return "quick_aligames"
}

func (r *QuickAliGamesAuth) GetVendor() model.Vendor {
	return model.VendorQuickAligames
}

func (r *QuickAliGamesAuth) Start() error {

	var db = storage.AuthDatabase()
	if db == nil {
		return errors.New("quick require mysql database")
	}

	db.AutoMigrate(&model.AuthQuickAliGames{})

	r.tlogger = NewTLogger(model.VendorQuickAligames)

	return nil
}

func (r *QuickAliGamesAuth) Verify(c *gin.Context, user *model.DoAuthRequest) error {
	// verify by msdk
	isVerify, err := quick_sdk.VerifyLoginQuick(user.VendorQuickAligames.Token, user.VendorQuickAligames.UserId)
	if err != nil {
		return err
	}

	if !isVerify {
		return errors.New("quick verify failed.")
	}

	return nil
}

func (r *QuickAliGamesAuth) FetchUser(c *gin.Context, user *model.DoAuthRequest) (*model.Account, error) {

	// check db, if not exist create new account.

	var authQuick model.AuthQuickAliGames
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

func (r *QuickAliGamesAuth) CreateUser(c *gin.Context, user *model.DoAuthRequest) (*model.Account, error) {

	account := createAccount(user)
	err := storage.Insert(storage.AuthDatabase(), account)
	if err != nil {

		return nil, err
	}

	AuthQuick := &model.AuthQuickAliGames{
		GlobalId:       account.GlobalId,
		OpenId:         user.OpenId,
		Platform:       user.Platform,
		UserId:         user.VendorQuickAligames.UserId,
		UserName:       user.VendorQuickAligames.UserName,
		Token:          user.VendorQuickAligames.Token,
		ChannelVersion: user.VendorQuickAligames.ChannelVersion,
		ChannelName:    user.VendorQuickAligames.ChannelName,
		ChannelType:    user.VendorQuickAligames.ChannelType,
	}

	err = storage.Insert(storage.AuthDatabase(), AuthQuick)
	if err != nil {
		return nil, err
	}

	LogPlayerRegister(user, account.GlobalId)

	return account, nil
}

func (r *QuickAliGamesAuth) UpdateVendor(globalId int64, user *model.DoAuthRequest) {
	var authQuick model.AuthQuickAliGames
	data := map[string]interface{}{
		"user_id":         user.VendorQuickAligames.UserId,
		"user_name":       user.VendorQuickAligames.UserName,
		"token":           user.VendorQuickAligames.Token,
		"channel_version": user.VendorQuickAligames.ChannelVersion,
		"channel_name":    user.VendorQuickAligames.ChannelName,
		"channel_type":    user.VendorQuickAligames.ChannelType,
	}

	storage.UpdateSNSProfile(&authQuick, globalId, data)
}
