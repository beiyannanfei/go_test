package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/xykong/loveauth/errors"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"github.com/xykong/tlog"
	"github.com/xykong/loveauth/services/sms"
	"fmt"
	"github.com/xykong/loveauth/settings"
)

type MobileAuth struct {
	tlogger *tlog.Tlogger
}

func (r *MobileAuth) Name() string {
	return "mobile"
}

func (r *MobileAuth) GetVendor() model.Vendor {
	return model.VendorMobile
}

// swagger:route POST /auth/mobile auth auth_mobile
//
// Authenticate user by mobile id
//
// This will accept any mobile id by default.
// You must make sure that different client has different openId
// There is no other guarantee for security
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
func (r *MobileAuth) Start() error {

	var db = storage.AuthDatabase()
	if db == nil {
		return errors.New("auth mobile require mysql database")
	}

	db.AutoMigrate(&model.AuthMobile{})

	r.tlogger = NewTLogger(model.VendorMobile)

	return nil
}

//接口只负责验证码验证，发送验证码改为专用接口
func (r *MobileAuth) Verify(c *gin.Context, user *model.DoAuthRequest) error {
	//手机号白名单验证
	whiteMobiles := settings.GetStringSlice("loveauth_white_list", "mobiles")
	for _, m := range whiteMobiles {
		if m == fmt.Sprintf("%s|%s", user.OpenId, user.VendorMobile.VerifyCode) {
			return nil
		}
	}

	return sms.VerifyMobileCode(user.OpenId, user.VendorMobile.VerifyCode, false)
}

func (r *MobileAuth) FetchUser(c *gin.Context, user *model.DoAuthRequest) (*model.Account, error) {

	// check db, if not exist create new account.

	var authMobile model.AuthMobile
	err := storage.QueryVendorByOpenId(user.OpenId, &authMobile)
	if err != nil {
		// no 3rdParty account found, need create new one.
		return nil, nil
	}

	account := storage.QueryAccount(authMobile.GlobalId)

	// some thing wrong like: account removed from table.
	if account == nil {
		return nil, errors.New("account not found in database")
	}

	return account, nil
}

func (r *MobileAuth) CreateUser(c *gin.Context, user *model.DoAuthRequest) (*model.Account, error) {

	account := createAccount(user)
	err := storage.Insert(storage.AuthDatabase(), account)
	if err != nil {

		return nil, err
	}

	authMobile := &model.AuthMobile{
		GlobalId: account.GlobalId,
		OpenId:   user.OpenId,
		Platform: user.Platform,
	}

	err = storage.Insert(storage.AuthDatabase(), authMobile)
	if err != nil {

		return nil, err
	}

	LogPlayerRegister(user, account.GlobalId)

	return account, nil
}

func (r *MobileAuth) UpdateVendor(globalId int64, user *model.DoAuthRequest) {}
