package auth

import (
	"errors"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"github.com/xykong/tlog"
	"github.com/gin-gonic/gin"
	"github.com/xykong/loveauth/services/Login/webo_sdk"
)

type WeiboAuth struct {
	tlogger *tlog.Tlogger
}

func (r *WeiboAuth) Name() string {
	return "weibo"
}

func (r *WeiboAuth) GetVendor() model.Vendor {
	return model.VendorWeibo
}

func (r *WeiboAuth) Start() error {

	var db = storage.AuthDatabase()
	if db == nil {
		return errors.New("weibo require mysql database")
	}

	db.AutoMigrate(&model.AuthWeibo{})

	r.tlogger = NewTLogger(model.VendorWeibo)

	return nil
}

func (r *WeiboAuth) Verify(c *gin.Context, user *model.DoAuthRequest) error {
	// verify by msdk
	_, err := webo_sdk.WeiBoGetTokenInfo(user.VendorWeibo.TokenAccess)
	if err != nil {
		return err
	}

	return nil
}

func (r *WeiboAuth) FetchUser(c *gin.Context, user *model.DoAuthRequest) (*model.Account, error) {

	// check db, if not exist create new account.

	var authWeibo model.AuthWeibo
	err := storage.QueryVendorByOpenId(user.OpenId, &authWeibo)
	if err != nil {
		// no 3rdParty account found, need create new one.
		return nil, nil
	}

	account := storage.QueryAccount(authWeibo.GlobalId)

	// some thing wrong like: account removed from table.
	if account == nil {
		return nil, errors.New("account not found in database")
	}

	return account, nil
}

func (r *WeiboAuth) CreateUser(c *gin.Context, user *model.DoAuthRequest) (*model.Account, error) {

	account := createAccount(user)
	err := storage.Insert(storage.AuthDatabase(), account)
	if err != nil {

		return nil, err
	}

	authWeibo := &model.AuthWeibo{
		GlobalId:         account.GlobalId,
		OpenId:           user.OpenId,
		Platform:         user.Platform,
		TokenAccess:      user.VendorWeibo.TokenAccess,
		ExpirationAccess: user.VendorWeibo.ExpirationAccess,
		NickName:         user.VendorWeibo.NickName,
		Picture:          user.VendorWeibo.Picture,
	}

	err = storage.Insert(storage.AuthDatabase(), authWeibo)
	if err != nil {

		return nil, err
	}

	LogPlayerRegister(user, account.GlobalId)

	return account, nil
}

func (r *WeiboAuth) UpdateVendor(globalId int64, user *model.DoAuthRequest) {
	var authWechat model.AuthWeibo
	data := map[string]interface{}{
		"token_access":      user.VendorWeibo.TokenAccess,
		"expiration_access": user.VendorWeibo.ExpirationAccess,
		"nick_name":         user.VendorWeibo.NickName,
		"picture":           user.VendorWeibo.Picture,
	}

	storage.UpdateSNSProfile(&authWechat, globalId, data)
}
