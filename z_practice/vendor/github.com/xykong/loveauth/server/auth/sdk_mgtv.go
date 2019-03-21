package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/xykong/loveauth/services/Login/mangguotv_sdk"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"github.com/xykong/tlog"
)

type MgtvAuth struct {
	tlogger *tlog.Tlogger
}

func (r *MgtvAuth) Name() string {
	return "mgtv"
}

func (r *MgtvAuth) GetVendor() model.Vendor {
	return model.VendorMgtv
}

// swagger:route POST /auth/mgtv auth auth_mgtv
//
// Authenticate user by mgtv
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
func (r *MgtvAuth) Start() error {

	var db = storage.AuthDatabase()
	if db == nil {
		return errors.New("mgtv require mysql database")
	}

	db.AutoMigrate(&model.AuthMgtv{})

	r.tlogger = NewTLogger(model.VendorMgtv)

	return nil
}

func (r *MgtvAuth) Verify(c *gin.Context, user *model.DoAuthRequest) error {
	// verify by mgtv_sdk

	err := mangguotv_sdk.ValidateUser(user.VendorMgtv.GameUserData.Ticket, user.VendorMgtv.GameUserData.ThirdId)
	if err != nil {
		return err
	}

	return nil
}

func (r *MgtvAuth) FetchUser(c *gin.Context, user *model.DoAuthRequest) (*model.Account, error) {

	// check db, if not exist create new account.

	var authMgtv model.AuthMgtv
	err := storage.QueryVendorByOpenId(user.OpenId, &authMgtv)
	if err != nil {
		// no 3rdParty account found, need create new one.
		return nil, nil
	}

	account := storage.QueryAccount(authMgtv.GlobalId)

	// some thing wrong like: account removed from table.
	if account == nil {
		return nil, errors.New("account not found in database")
	}

	return account, nil
}

func (r *MgtvAuth) CreateUser(c *gin.Context, user *model.DoAuthRequest) (*model.Account, error) {

	account := createAccount(user)
	err := storage.Insert(storage.AuthDatabase(), account)
	if err != nil {

		return nil, err
	}

	authMgtv := &model.AuthMgtv{
		GlobalId:    account.GlobalId,
		OpenId:      user.OpenId,
		Platform:    user.Platform,
		AccessToken: user.VendorMgtv.GameUserData.LoginAccount,
		NickName:    user.VendorMgtv.GameUserData.NickName,
		Ticket:      user.VendorMgtv.GameUserData.Ticket,
	}

	err = storage.Insert(storage.AuthDatabase(), authMgtv)
	if err != nil {

		return nil, err
	}

	LogPlayerRegister(user, account.GlobalId)

	return account, nil
}

func (r *MgtvAuth) UpdateVendor(globalId int64, user *model.DoAuthRequest) {
	var authMgtv model.AuthMgtv
	data := map[string]interface{}{
		"ticket":       user.VendorMgtv.GameUserData.Ticket,
		"access_token": user.VendorMgtv.GameUserData.LoginAccount,
	}

	storage.UpdateSNSProfile(&authMgtv, globalId, data)
}
