package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/xykong/loveauth/services/Login/vivo_sdk"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"github.com/xykong/tlog"
)

type VivoAuth struct {
	tlogger *tlog.Tlogger
}

func (r *VivoAuth) Name() string {
	return "vivo"
}

func (r *VivoAuth) GetVendor() model.Vendor {
	return model.VendorVivo
}

// swagger:route POST /auth/vivo auth auth_vivo
//
// Authenticate user by vivo
//
// This will accept any device id by default.
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
func (r *VivoAuth) Start() error {

	var db = storage.AuthDatabase()
	if db == nil {
		return errors.New("vivo require mysql database")
	}

	db.AutoMigrate(&model.AuthVivo{})

	r.tlogger = NewTLogger(model.VendorVivo)

	return nil
}

func (r *VivoAuth) Verify(c *gin.Context, user *model.DoAuthRequest) error {
	// verify by vivo_sdk
	_, err := vivo_sdk.AuthToken(user.VendorVivo.AuthToken)
	if err != nil {
		return err
	}

	return nil
}

func (r *VivoAuth) FetchUser(c *gin.Context, user *model.DoAuthRequest) (*model.Account, error) {

	// check db, if not exist create new account.

	var authVivo model.AuthVivo
	err := storage.QueryVendorByOpenId(user.OpenId, &authVivo)
	if err != nil {
		// no 3rdParty account found, need create new one.
		return nil, nil
	}

	account := storage.QueryAccount(authVivo.GlobalId)

	// some thing wrong like: account removed from table.
	if account == nil {
		return nil, errors.New("account not found in database")
	}

	return account, nil
}

func (r *VivoAuth) CreateUser(c *gin.Context, user *model.DoAuthRequest) (*model.Account, error) {

	account := createAccount(user)
	err := storage.Insert(storage.AuthDatabase(), account)
	if err != nil {

		return nil, err
	}

	authVivo := &model.AuthVivo{
		GlobalId:         account.GlobalId,
		OpenId:           user.OpenId,
		Platform:         user.Platform,
		AuthToken:        user.VendorVivo.AuthToken,
		ExpirationAccess: user.VendorVivo.ExpirationAccess,
		NickName:         user.VendorVivo.NickName,
		Picture:          user.VendorVivo.Picture,
	}

	err = storage.Insert(storage.AuthDatabase(), authVivo)
	if err != nil {

		return nil, err
	}

	LogPlayerRegister(user, account.GlobalId)

	return account, nil
}

func (r *VivoAuth) UpdateVendor(globalId int64, user *model.DoAuthRequest) {
	var authVivo model.AuthVivo
	data := map[string]interface{}{
		"auth_token":        user.VendorVivo.AuthToken,
		"expiration_access": user.VendorVivo.ExpirationAccess,
	}

	storage.UpdateSNSProfile(&authVivo, globalId, data)
}
