package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/xykong/loveauth/services/msdk"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"github.com/xykong/tlog"
)

type GuestAuth struct {
	tlogger *tlog.Tlogger
}

func (r *GuestAuth) Name() string {
	return "guest"
}

func (r *GuestAuth) GetVendor() model.Vendor {
	return model.VendorMsdkGuest
}

// swagger:route POST /auth/guest auth auth_guest
//
// Authenticate user by msdk guest
//
// This will accept message get from msdk client by guest choice.
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
func (r *GuestAuth) Start() error {

	var db = storage.AuthDatabase()
	if db == nil {
		return errors.New("guest require mysql database")
	}

	db.AutoMigrate(&model.AuthGuest{})

	r.tlogger = NewTLogger(model.VendorMsdkGuest)

	return nil
}

func (r *GuestAuth) Verify(c *gin.Context, user *model.DoAuthRequest) error {

	// verify by msdk
	_, err := msdk.GuestCheckToken(user.OpenId, user.VendorMsdkGuest.AccessTokenValue)

	return err
}

func (r *GuestAuth) FetchUser(c *gin.Context, user *model.DoAuthRequest) (*model.Account, error) {

	// check db, if not exist create new account.

	var authGuest model.AuthGuest
	err := storage.QueryVendorByOpenId(user.OpenId, &authGuest)
	if err != nil {
		// no 3rdParty account found, need create new one.
		return nil, nil
	}

	account := storage.QueryAccount(authGuest.GlobalId)

	// some thing wrong like: account removed from table.
	if account == nil {
		return nil, errors.New("account not found in database")
	}

	return account, nil
}

func (r *GuestAuth) CreateUser(c *gin.Context, user *model.DoAuthRequest) (*model.Account, error) {

	account := createAccount(user)
	err := storage.Insert(storage.AuthDatabase(), account)
	if err != nil {

		return nil, err
	}

	authGuest := &model.AuthGuest{

		GlobalId:         account.GlobalId,
		OpenId:           user.OpenId,
		Platform:         user.Platform,
		TokenAccess:      user.VendorMsdkGuest.AccessTokenValue,
		ExpirationAccess: user.VendorMsdkGuest.AccessTokenExpiration,
		Pf:               user.VendorMsdkGuest.Pf,
		PfKey:            user.VendorMsdkGuest.PfKey,
	}

	err = storage.Insert(storage.AuthDatabase(), authGuest)
	if err != nil {

		return nil, err
	}

	LogPlayerRegister(user, account.GlobalId)

	return account, nil
}

func (r *GuestAuth) UpdateVendor(globalId int64, user *model.DoAuthRequest) {
	var authGuest model.AuthGuest
	data := map[string]interface{}{
		"token_access":      user.VendorMsdkGuest.AccessTokenValue,
		"expiration_access": user.VendorMsdkGuest.AccessTokenExpiration,
		"pf":                user.VendorMsdkGuest.Pf,
		"pf_key":            user.VendorMsdkGuest.PfKey,
	}
	storage.UpdateSNSProfile(&authGuest, globalId, data)
}
