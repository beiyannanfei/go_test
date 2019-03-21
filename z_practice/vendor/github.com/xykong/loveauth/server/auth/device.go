package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"github.com/xykong/tlog"
)

type DeviceAuth struct {
	tlogger *tlog.Tlogger
}

func (r *DeviceAuth) Name() string {
	return "device"
}

func (r *DeviceAuth) GetVendor() model.Vendor {
	return model.VendorDevice
}

// swagger:route POST /auth/device auth auth_device
//
// Authenticate user by device id
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
func (r *DeviceAuth) Start() error {

	var db = storage.AuthDatabase()
	if db == nil {
		return errors.New("auth device require mysql database")
	}

	db.AutoMigrate(&model.AuthDevice{})

	r.tlogger = NewTLogger(model.VendorDevice)

	return nil
}

func (r *DeviceAuth) Verify(c *gin.Context, user *model.DoAuthRequest) error {

	return nil
}

func (r *DeviceAuth) FetchUser(c *gin.Context, user *model.DoAuthRequest) (*model.Account, error) {

	// check db, if not exist create new account.

	var authDevice model.AuthDevice
	err := storage.QueryVendorByOpenId(user.OpenId, &authDevice)
	if err != nil {
		// no 3rdParty account found, need create new one.
		return nil, nil
	}

	account := storage.QueryAccount(authDevice.GlobalId)

	// some thing wrong like: account removed from table.
	if account == nil {
		return nil, errors.New("account not found in database")
	}

	return account, nil
}

func (r *DeviceAuth) CreateUser(c *gin.Context, user *model.DoAuthRequest) (*model.Account, error) {

	account := createAccount(user)
	err := storage.Insert(storage.AuthDatabase(), account)
	if err != nil {

		return nil, err
	}

	authDevice := &model.AuthDevice{
		GlobalId: account.GlobalId,
		OpenId:   user.OpenId,
		Platform: user.Platform,
	}

	err = storage.Insert(storage.AuthDatabase(), authDevice)
	if err != nil {

		return nil, err
	}

	LogPlayerRegister(user, account.GlobalId)

	return account, nil
}

func (r *DeviceAuth) UpdateVendor(globalId int64, user *model.DoAuthRequest) {}
