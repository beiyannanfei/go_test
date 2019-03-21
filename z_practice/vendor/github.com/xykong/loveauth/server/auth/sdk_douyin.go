package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/xykong/loveauth/services/Login/douyin_sdk"
	"github.com/xykong/loveauth/settings"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"github.com/xykong/tlog"
)

type DouyinAuth struct {
	tlogger *tlog.Tlogger
}

func (r *DouyinAuth) Name() string {
	return "douyin"
}

func (r *DouyinAuth) GetVendor() model.Vendor {
	return model.VendorDouyin
}

// swagger:route POST /auth/douyin auth auth_douyin
//
// Authenticate user by douyin
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
func (r *DouyinAuth) Start() error {

	var db = storage.AuthDatabase()
	if db == nil {
		return errors.New("douyin require mysql database")
	}

	db.AutoMigrate(&model.AuthDouyin{})

	r.tlogger = NewTLogger(model.VendorDouyin)

	return nil
}

func (r *DouyinAuth) Verify(c *gin.Context, user *model.DoAuthRequest) error {
	// verify by douyin_sdk
	clientKey := settings.GetString("lovepay", "douyin.clientKey")
	clientSecret := settings.GetString("lovepay", "douyin.clientSecret")
	err := douyin_sdk.CheckUser(clientKey, clientSecret, user.OpenId)
	if err != nil {
		return err
	}

	return nil
}

func (r *DouyinAuth) FetchUser(c *gin.Context, user *model.DoAuthRequest) (*model.Account, error) {

	// check db, if not exist create new account.

	var authDouyin model.AuthDouyin
	err := storage.QueryVendorByOpenId(user.OpenId, &authDouyin)
	if err != nil {
		// no 3rdParty account found, need create new one.
		return nil, nil
	}

	account := storage.QueryAccount(authDouyin.GlobalId)

	// some thing wrong like: account removed from table.
	if account == nil {
		return nil, errors.New("account not found in database")
	}

	return account, nil
}

func (r *DouyinAuth) CreateUser(c *gin.Context, user *model.DoAuthRequest) (*model.Account, error) {

	account := createAccount(user)
	err := storage.Insert(storage.AuthDatabase(), account)
	if err != nil {

		return nil, err
	}

	authDouyin := &model.AuthDouyin{
		GlobalId:    account.GlobalId,
		OpenId:      user.OpenId,
		Platform:    user.Platform,
		AccessToken: user.VendorDouyin.GameUserData.AccessToken,
		Uid:         user.VendorDouyin.GameUserData.Uid,
		UserType:    user.VendorDouyin.GameUserData.UserType,
	}

	err = storage.Insert(storage.AuthDatabase(), authDouyin)
	if err != nil {

		return nil, err
	}

	LogPlayerRegister(user, account.GlobalId)

	return account, nil
}

func (r *DouyinAuth) UpdateVendor(globalId int64, user *model.DoAuthRequest) {
	var authDouyin model.AuthDouyin
	data := map[string]interface{}{
		"access_token": user.VendorDouyin.GameUserData.AccessToken,
	}

	storage.UpdateSNSProfile(&authDouyin, globalId, data)
}
