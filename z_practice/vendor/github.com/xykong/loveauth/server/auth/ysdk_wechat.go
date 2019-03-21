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

type YSDKWechatAuth struct {
	tlogger *tlog.Tlogger
}

func (r *YSDKWechatAuth) Name() string {
	return "ysdk_wechat"
}

func (r *YSDKWechatAuth) GetVendor() model.Vendor {
	return model.VendorYsdkWechat
}

func (r *YSDKWechatAuth) Start() error {
	var db = storage.AuthDatabase()
	if db == nil {
		return errors.New("ysdk wechat require mysql database")
	}

	if err := db.AutoMigrate(&model.AuthYsdkWechat{}).Error; err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Info("AutoMigrate AuthYsdkWechat failed.")
		return nil
	}

	r.tlogger = NewTLogger(model.VendorYsdkWechat)

	return nil
}

func (r *YSDKWechatAuth) Verify(c *gin.Context, user *model.DoAuthRequest) error {
	userIp := c.Request.RemoteAddr

	// verify by ysdk
	_, err := ysdk.VerifyLoginWechat(user.OpenId, user.VendorYsdkWechat.TokenAccess, userIp)

	return err
}

func (r *YSDKWechatAuth) FetchUser(c *gin.Context, user *model.DoAuthRequest) (*model.Account, error) {
	// check db, if not exist create new account.
	var authYsdkWechat model.AuthYsdkWechat
	err := storage.QueryVendorByOpenId(user.OpenId, &authYsdkWechat)
	if err != nil {
		// no 3rdParty account found, need create new one.
		return nil, nil
	}

	account := storage.QueryAccount(authYsdkWechat.GlobalId)

	// some thing wrong like: account removed from table.
	if account == nil {
		return nil, errors.New("account not found in database")
	}

	return account, nil
}

func (r *YSDKWechatAuth) CreateUser(c *gin.Context, user *model.DoAuthRequest) (*model.Account, error) {
	account := createAccount(user)
	err := storage.Insert(storage.AuthDatabase(), account)
	if err != nil {
		return nil, err
	}

	authYsdkWechat := &model.AuthYsdkWechat{
		GlobalId:         account.GlobalId,
		OpenId:           user.OpenId,
		Platform:         user.Platform,
		TokenAccess:      user.VendorYsdkWechat.TokenAccess,
		ExpirationAccess: user.VendorYsdkWechat.ExpirationAccess,
		TokenRefresh:     user.VendorYsdkWechat.TokenRefresh,
		NickName:         user.VendorYsdkWechat.NickName,
		Picture:          user.VendorYsdkWechat.Picture,
		UnionId:          user.VendorYsdkWechat.UnionId,
	}

	err = storage.Insert(storage.AuthDatabase(), authYsdkWechat)
	if err != nil {
		return nil, err
	}

	LogPlayerRegister(user, account.GlobalId)

	return account, nil
}

func (r *YSDKWechatAuth) UpdateVendor(globalId int64, user *model.DoAuthRequest) {
	var authYsdkWechat model.AuthYsdkWechat
	data := map[string]interface{}{
		"token_access":      user.VendorYsdkWechat.TokenAccess,
		"expiration_access": user.VendorYsdkWechat.ExpirationAccess,
		"token_refresh":     user.VendorYsdkWechat.TokenRefresh,
		"nick_name":         user.VendorYsdkWechat.NickName,
		"picture":           user.VendorYsdkWechat.Picture,
		"union_id":          user.VendorYsdkWechat.UnionId,
	}
	storage.UpdateSNSProfile(&authYsdkWechat, globalId, data)
}
