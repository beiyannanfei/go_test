package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/xykong/loveauth/services/msdk"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"github.com/xykong/tlog"
)

type WechatAuth struct {
	tlogger *tlog.Tlogger
}

func (r *WechatAuth) Name() string {
	return "wechat"
}

func (r *WechatAuth) GetVendor() model.Vendor {
	return model.VendorMsdkWechat
}

// swagger:route POST /auth/wechat auth auth_wechat
//
// Authenticate user by msdk wechat
//
// This will accept message get from msdk client by wechat choice.
// And verify it by msdk server side.
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       200: ResponseAuth
func (r *WechatAuth) Start() error {

	var db = storage.AuthDatabase()
	if db == nil {
		return errors.New("wechat require mysql database")
	}

	db.AutoMigrate(&model.AuthWechat{})

	r.tlogger = NewTLogger(model.VendorMsdkWechat)

	return nil
}

func (r *WechatAuth) Verify(c *gin.Context, user *model.DoAuthRequest) error {

	// verify by msdk
	_, err := msdk.CheckToken(user.OpenId, user.VendorMsdkWechat.AccessTokenValue)

	return err
}

func (r *WechatAuth) FetchUser(c *gin.Context, user *model.DoAuthRequest) (*model.Account, error) {

	// check db, if not exist create new account.

	var authWechat model.AuthWechat
	err := storage.QueryVendorByOpenId(user.OpenId, &authWechat)
	if err != nil {
		// no 3rdParty account found, need create new one.
		return nil, nil
	}

	account := storage.QueryAccount(authWechat.GlobalId)

	// some thing wrong like: account removed from table.
	if account == nil {
		return nil, errors.New("account not found in database")
	}

	return account, nil
}

func (r *WechatAuth) CreateUser(c *gin.Context, user *model.DoAuthRequest) (*model.Account, error) {

	account := createAccount(user)
	err := storage.Insert(storage.AuthDatabase(), account)
	if err != nil {

		return nil, err
	}

	authWechat := &model.AuthWechat{

		GlobalId:         account.GlobalId,
		OpenId:           user.OpenId,
		Platform:         user.Platform,
		TokenAccess:      user.VendorMsdkWechat.AccessTokenValue,
		ExpirationAccess: user.VendorMsdkWechat.AccessTokenExpiration,
		Pf:               user.VendorMsdkWechat.Pf,
		PfKey:            user.VendorMsdkWechat.PfKey,
	}

	err = storage.Insert(storage.AuthDatabase(), authWechat)
	if err != nil {

		return nil, err
	}

	LogPlayerRegister(user, account.GlobalId)

	return account, nil
}

func (r *WechatAuth) UpdateVendor(globalId int64, user *model.DoAuthRequest) {
	var authWechat model.AuthWechat
	data := map[string]interface{}{
		"token_access":      user.VendorMsdkWechat.AccessTokenValue,
		"expiration_access": user.VendorMsdkWechat.AccessTokenExpiration,
		"pf":                user.VendorMsdkWechat.Pf,
		"pf_key":            user.VendorMsdkWechat.PfKey,
	}
	storage.UpdateSNSProfile(&authWechat, globalId, data)
}
