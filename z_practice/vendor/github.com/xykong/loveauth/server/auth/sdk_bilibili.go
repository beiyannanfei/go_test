package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/xykong/loveauth/errors"
	"github.com/xykong/loveauth/services/Login/bilibili"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"github.com/xykong/tlog"
	"github.com/sirupsen/logrus"
)

type BilibiliAuth struct {
	tlogger *tlog.Tlogger
}

func (r *BilibiliAuth) Name() string {
	return "bilibili"
}

func (r *BilibiliAuth) GetVendor() model.Vendor {
	return model.VendorBilibili
}

// swagger:route POST /auth/bilibili auth auth_bilibili
//
// Authenticate user by bilibili
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
func (r *BilibiliAuth) Start() error {
	var db = storage.AuthDatabase()
	if nil == db {
		return errors.New("bilibili require mysql database")
	}

	db.AutoMigrate(&model.AuthBilibili{})
	r.tlogger = NewTLogger(model.VendorBilibili)
	return nil
}

func (r *BilibiliAuth) Verify(c *gin.Context, user *model.DoAuthRequest) error {
	logrus.WithFields(logrus.Fields{
		"access_key": user.VendorBilibili.LoginResult.AccessToken,
		"userId": user.VendorBilibili.LoginResult.UserId,
	}).Info("bilibili")
	resp, err := bilibili.AuthToken(user.VendorBilibili.LoginResult.AccessToken, user.VendorBilibili.LoginResult.UserId)
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"user": user,
			"err":  err,
		}).Error("BilibiliAuth Verify failed.")

		return err
	}

	logrus.WithFields(logrus.Fields{
		"resp": resp,
	}).Info("BilibiliAuth Verify success")

	return nil
}

func (r *BilibiliAuth) FetchUser(c *gin.Context, user *model.DoAuthRequest) (*model.Account, error) {
	var authBilibili model.AuthBilibili
	err := storage.QueryVendorByOpenId(user.OpenId, &authBilibili)
	if nil != err {

		return nil, nil
	}

	account := storage.QueryAccount(authBilibili.GlobalId)

	if nil == account {
		return nil, errors.New("account not found in database for bilibili")
	}

	return account, nil
}

func (r *BilibiliAuth) CreateUser(c *gin.Context, user *model.DoAuthRequest) (*model.Account, error) {
	account := createAccount(user)
	err := storage.Insert(storage.AuthDatabase(), account)
	if nil != err {
		return nil, err
	}

	authBilibili := &model.AuthBilibili{
		GlobalId:     account.GlobalId,
		OpenId:       user.OpenId,
		Platform:     user.Platform,
		Code:         user.VendorBilibili.LoginResult.Code,
		Message:      user.VendorBilibili.LoginResult.Message,
		UserId:       user.VendorBilibili.LoginResult.UserId,
		UserName:     user.VendorBilibili.LoginResult.UserName,
		NickName:     user.VendorBilibili.LoginResult.NickName,
		AccessToken:  user.VendorBilibili.LoginResult.AccessToken,
		ExpireTimes:  user.VendorBilibili.LoginResult.ExpireTimes,
		RefreshToken: user.VendorBilibili.LoginResult.RefreshToken,
	}

	err = storage.Insert(storage.AuthDatabase(), authBilibili)
	if nil != err {
		return nil, err
	}

	LogPlayerRegister(user, account.GlobalId)

	return account, nil
}

func (r *BilibiliAuth) UpdateVendor(globalId int64, user *model.DoAuthRequest) {
	var authBilibli model.AuthBilibili
	data := map[string]interface{}{
		"code":          user.VendorBilibili.LoginResult.Code,
		"message":       user.VendorBilibili.LoginResult.Message,
		"user_id":       user.VendorBilibili.LoginResult.UserId,
		"user_name":     user.VendorBilibili.LoginResult.UserName,
		"nick_name":     user.VendorBilibili.LoginResult.NickName,
		"access_token":  user.VendorBilibili.LoginResult.AccessToken,
		"expire_times":  user.VendorBilibili.LoginResult.ExpireTimes,
		"refresh_token": user.VendorBilibili.LoginResult.RefreshToken,
	}

	storage.UpdateSNSProfile(&authBilibli, globalId, data)
}
