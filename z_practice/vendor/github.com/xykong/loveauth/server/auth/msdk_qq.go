package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/services/msdk"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"github.com/xykong/tlog"
)

type QQAuth struct {
	tlogger *tlog.Tlogger
}

func (r *QQAuth) Name() string {
	return "qq"
}

func (r *QQAuth) GetVendor() model.Vendor {
	return model.VendorMsdkQQ
}

// swagger:route POST /auth/qq auth auth_qq
//
// Authenticate user by msdk qq
//
// This will accept message get from msdk client by qq choice.
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
func (r *QQAuth) Start() error {

	var db = storage.AuthDatabase()
	if db == nil {
		return errors.New("qq require mysql database")
	}

	if err := db.AutoMigrate(&model.AuthQQ{}).Error; err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Info("AutoMigrate AuthQQ failed.")

		return nil
	}

	r.tlogger = NewTLogger(model.VendorMsdkQQ)

	return nil
}

func (r *QQAuth) Verify(c *gin.Context, user *model.DoAuthRequest) error {

	userIp := c.Request.RemoteAddr

	// verify by msdk
	_, err := msdk.VerifyLogin(user.OpenId, user.VendorMsdkQQ.AccessTokenValue, userIp)

	return err
}

func (r *QQAuth) FetchUser(c *gin.Context, user *model.DoAuthRequest) (*model.Account, error) {

	// check db, if not exist create new account.

	var authQQ model.AuthQQ
	err := storage.QueryVendorByOpenId(user.OpenId, &authQQ)
	if err != nil {
		// no 3rdParty account found, need create new one.
		return nil, nil
	}

	account := storage.QueryAccount(authQQ.GlobalId)

	// some thing wrong like: account removed from table.
	if account == nil {
		return nil, errors.New("account not found in database")
	}

	return account, nil
}

func (r *QQAuth) CreateUser(c *gin.Context, user *model.DoAuthRequest) (*model.Account, error) {

	account := createAccount(user)
	err := storage.Insert(storage.AuthDatabase(), account)
	if err != nil {

		return nil, err
	}

	authQQ := &model.AuthQQ{

		GlobalId:         account.GlobalId,
		OpenId:           user.OpenId,
		TokenAccess:      user.VendorMsdkQQ.AccessTokenValue,
		ExpirationAccess: user.VendorMsdkQQ.AccessTokenExpiration,
		Pf:               user.VendorMsdkQQ.Pf,
		PfKey:            user.VendorMsdkQQ.PfKey,
	}

	err = storage.Insert(storage.AuthDatabase(), authQQ)
	if err != nil {

		return nil, err
	}

	LogPlayerRegister(user, account.GlobalId)

	return account, nil
}

func (r *QQAuth) UpdateVendor(globalId int64, user *model.DoAuthRequest) {
	var authQQ model.AuthQQ
	data := map[string]interface{}{
		"token_access":      user.VendorMsdkQQ.AccessTokenValue,
		"expiration_access": user.VendorMsdkQQ.AccessTokenExpiration,
		"pf":                user.VendorMsdkQQ.Pf,
		"pf_key":            user.VendorMsdkQQ.PfKey,
	}
	storage.UpdateSNSProfile(&authQQ, globalId, data)
}
