package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"github.com/xykong/tlog"
	"github.com/xykong/loveauth/services/Login/ysdk"
)

type YSDKQQAuth struct {
	tlogger *tlog.Tlogger
}

func (r *YSDKQQAuth) Name() string {
	return "ysdk_qq"
}

func (r *YSDKQQAuth) GetVendor() model.Vendor {
	return model.VendorYsdkQQ
}

func (r *YSDKQQAuth) Start() error {
	var db = storage.AuthDatabase()
	if db == nil {
		return errors.New("ysdk qq require mysql database")
	}

	if err := db.AutoMigrate(&model.AuthYsdkQQ{}).Error; err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Info("AutoMigrate AuthYsdkQQ failed.")
		return nil
	}

	r.tlogger = NewTLogger(model.VendorYsdkQQ)

	return nil
}

func (r *YSDKQQAuth) Verify(c *gin.Context, user *model.DoAuthRequest) error {
	userIp := c.Request.RemoteAddr

	// verify by ysdk
	_, err := ysdk.VerifyLoginQQ(user.OpenId, user.VendorYsdkQQ.TokenAccess, userIp)

	return err
}

func (r *YSDKQQAuth) FetchUser(c *gin.Context, user *model.DoAuthRequest) (*model.Account, error) {
	// check db, if not exist create new account.
	var authYsdkQQ model.AuthYsdkQQ
	err := storage.QueryVendorByOpenId(user.OpenId, &authYsdkQQ)
	if err != nil {
		// no 3rdParty account found, need create new one.
		return nil, nil
	}

	account := storage.QueryAccount(authYsdkQQ.GlobalId)

	// some thing wrong like: account removed from table.
	if account == nil {
		return nil, errors.New("account not found in database")
	}

	return account, nil
}

func (r *YSDKQQAuth) CreateUser(c *gin.Context, user *model.DoAuthRequest) (*model.Account, error) {
	account := createAccount(user)
	err := storage.Insert(storage.AuthDatabase(), account)
	if err != nil {
		return nil, err
	}

	authYsdkQQ := &model.AuthYsdkQQ{
		GlobalId:         account.GlobalId,
		OpenId:           user.OpenId,
		Platform:         user.Platform,
		TokenAccess:      user.VendorYsdkQQ.TokenAccess,
		ExpirationAccess: user.VendorYsdkQQ.ExpirationAccess,
		TokenPay:         user.VendorYsdkQQ.TokenPay,
		Pf:               user.VendorYsdkQQ.Pf,
		PfKey:            user.VendorYsdkQQ.PfKey,
		NickName:         user.VendorYsdkQQ.NickName,
		Picture:          user.VendorYsdkQQ.Picture,
	}

	err = storage.Insert(storage.AuthDatabase(), authYsdkQQ)
	if err != nil {
		return nil, err
	}

	LogPlayerRegister(user, account.GlobalId)

	return account, nil
}

func (r *YSDKQQAuth) UpdateVendor(globalId int64, user *model.DoAuthRequest) {
	var authYsdkQQ model.AuthYsdkQQ
	data := map[string]interface{}{
		"token_access":      user.VendorYsdkQQ.TokenAccess,
		"expiration_access": user.VendorYsdkQQ.ExpirationAccess,
		"token_pay":         user.VendorYsdkQQ.TokenPay,
		"pf":                user.VendorYsdkQQ.Pf,
		"pf_key":            user.VendorYsdkQQ.PfKey,
		"nick_name":         user.VendorYsdkQQ.NickName,
		"picture":           user.VendorYsdkQQ.Picture,
	}
	storage.UpdateSNSProfile(&authYsdkQQ, globalId, data)
}
