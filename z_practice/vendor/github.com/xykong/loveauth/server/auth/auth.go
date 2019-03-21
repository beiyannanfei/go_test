package auth

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	couponLib "github.com/xykong/loveauth/coupon"
	"github.com/xykong/loveauth/errors"
	"github.com/xykong/loveauth/settings"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"github.com/xykong/loveauth/utils"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"
)

type WorkingMode string

const (
	Normal           WorkingMode = "Normal"
	RegisterByCoupon WorkingMode = "RegisterByCoupon"
	RegisterClosed   WorkingMode = "RegisterClosed"
)

//
// in: body
// swagger:parameters auth_qq auth_guest auth_wechat auth_device auth_mobile auth_vivo auth_bilibili auth_douyin auth_mgtv
type DoAuthRequestBodyParams struct {
	//
	// swagger:allOf
	// in: body
	Request model.DoAuthRequest `form:"Request" json:"Request" binding:"required"`
}

// A auth.ResponseAuth is an response message to client.
// swagger:response ResponseAuth
type ResponseAuth struct {
	// in: body
	Body struct {
		// The response code
		//
		// Required: true
		Code int64 `json:"code"`
		// The response message
		//
		// Required: true
		Message           string `json:"message"`
		AccessToken       string `json:"accessToken"`
		ExpirationSeconds int64  `json:"expirationSeconds"`
		RefreshToken      string `json:"refreshToken"`
	}
}

func GenerateToken(globalId int64, request model.DoAuthRequest, expirationSeconds int64, salt string) string {

	var timeStamp = time.Now().Unix() + expirationSeconds*1000

	var base = fmt.Sprintf("%d%s%d%s%d", globalId, salt, timeStamp, request.DeviceId, rand.Int())

	return fmt.Sprintf("%x", md5.Sum([]byte(base)))
}

// Providers is list of known/available usedProviders.
type Providers map[string]Provider

var usedProviders = Providers{}

func Use(aps ...Provider) {
	for _, provider := range aps {

		if usedProviders[provider.Name()] != nil {
			logrus.WithFields(logrus.Fields{
				"provider": provider.Name(),
			}).Warn("provider replaced.")
		}

		usedProviders[provider.Name()] = provider
	}
}

func Start(group *gin.RouterGroup) {

	// checking working mode.
	workingMode := WorkingMode(settings.GetString("loveauth", "auth.WorkingMode"))

	if workingMode == "" {
		settings.Set("loveauth", "auth.WorkingMode", Normal)
		workingMode = Normal
	}

	switch workingMode {
	case Normal, RegisterByCoupon, RegisterClosed:
	default:
		logrus.WithFields(logrus.Fields{
			"workingMode": settings.GetString("loveauth", "auth.WorkingMode"),
		}).Fatal("Invalid WorkingMode.")
		return
	}

	logrus.WithFields(logrus.Fields{
		"workingMode": workingMode,
	}).Info("auth WorkingMode.")

	for _, p := range usedProviders {

		var provider = p
		provider.Start()
		group.POST("/"+p.Name(), func(context *gin.Context) {

			var authenticator = Authenticator{provider, workingMode}
			authenticator.Auth(context)
		})
	}

	group.POST("/token", token)
	group.POST("/gm", gm_login)
}

type Provider interface {
	Name() string
	Start() error
	Verify(c *gin.Context, user *model.DoAuthRequest) error
	FetchUser(c *gin.Context, user *model.DoAuthRequest) (*model.Account, error)
	CreateUser(c *gin.Context, user *model.DoAuthRequest) (*model.Account, error)
	GetVendor() model.Vendor
	UpdateVendor(globalId int64, user *model.DoAuthRequest)
}

type Authenticator struct {
	Provider    Provider
	workingMode WorkingMode
}

type RequestState struct {
	Context       *gin.Context
	Request       model.DoAuthRequest
	Account       *model.Account
	IsNewAccount  bool
	IsInWhiteList bool
	AuthClientIp  string
}

func (o *Authenticator) Auth(c *gin.Context) {

	var rs = RequestState{
		Context: c,
	}

	// validation
	if err := c.BindJSON(&rs.Request); err != nil {
		utils.QuickReply(c, errors.ServerFail, "连接服务器异常")
		return
	}

	request := rs.Request

	// ignore platform and vendor string for swagger
	if request.Platform == "string" {
		request.Platform = ""
	}
	if request.Vendor == "string" {
		request.Vendor = ""
	}

	if request.ClientIp == "" {
		request.ClientIp = c.ClientIP()
	}

	rs.AuthClientIp = c.ClientIP()

	if len(request.Platform) > 0 && !storage.IsValidPlatforms(request.Platform) {
		utils.QuickReply(c, errors.InvalidPlatformSpecified, "应用暂不支持该平台")
		return
	}

	if len(request.Vendor) > 0 && !storage.IsValidVendors(request.Vendor) {
		utils.QuickReply(c, errors.InvalidVendorSpecified, "应用暂不支持该第三方登录")
		return
	}

	//注意：微博的过期时间存在差异(android: 剩余到期的秒数，ios: 到期时间的毫秒级时间戳)
	if rs.Request.Vendor == model.VendorWeibo && rs.Request.Platform == model.PlatformAndroid { //android手动转换为ios的格式
		rs.Request.VendorWeibo.ExpirationAccess = (time.Now().Unix() + rs.Request.VendorWeibo.ExpirationAccess) * 1000
	}

	if err := o.Provider.Verify(c, &request); err != nil {
		if ec, ok := err.(*errors.Type); ok {
			utils.QuickReply(c, ec.Code, ec.Message)
			return
		}

		utils.QuickReply(c, errors.VerifyBy3rdPartyFailed, "认证失败")
		return
	}

	// check db, if not exist create new account.

	var err error
	if rs.Account, err = o.Provider.FetchUser(c, &request); err != nil {
		utils.QuickReply(c, errors.ServerFail, "连接服务器异常")
		return
	}

	rs.IsInWhiteList = isWhiteListContains(&rs)

	err = o.createAccount(&rs)
	if err != nil {
		utils.QuickReplyError(c, err)
		return
	}

	err = o.filterByCoupon(&rs)
	if err != nil {
		utils.QuickReplyError(c, err)
		return
	}

	err = o.filterByState(&rs)
	if err != nil {
		utils.QuickReplyError(c, err)
		return
	}

	var setting = settings.Get("loveauth")
	var expirationSeconds = setting.GetInt64("token.ExpirationSeconds")
	var salt = setting.GetString("token.Salt")

	// generate token
	accessToken := GenerateToken(rs.Account.GlobalId, request, expirationSeconds, salt)
	refreshToken := GenerateToken(rs.Account.GlobalId, request, expirationSeconds, salt+salt)

	platform := request.Platform

	currTime := time.Now().Unix()
	var profile = model.Profile{
		GlobalId:          rs.Account.GlobalId,
		Token:             accessToken,
		RefreshToken:      refreshToken,
		ExpirationSeconds: expirationSeconds,
		Timestamp:         currTime,
		Name:              rs.Account.Name,
		Auth:              rs.Request,
		Vendor:            o.Provider.GetVendor(),
		Platform:          platform,
		NickName:          "",
		Picture:           "",
	}

	//////////////////////////////////////////////////////////////////////
	// TLog record.
	savedProfile := storage.WriteToken(&profile)
	if savedProfile != nil {
		LogPlayerLogout(savedProfile, currTime)
	}

	// fixme gs send this tlog now
	//LogPlayerLogin(&profile)
	//////////////////////////////////////////////////////////////////////

	if savedProfile != nil {

		rs.Account.AccumLoginTime = rs.Account.AccumLoginTime + currTime - savedProfile.Timestamp
	}
	rs.Account.LoginTime = currTime

	if rs.Account.State != model.Active && rs.IsInWhiteList && !setting.GetBool("auth.StrictVerify") {

		rs.Account.State = model.Active
	}

	err = storage.Save(storage.AuthDatabase(), rs.Account)
	if err != nil {

		utils.QuickReply(c, errors.ServerFail, "连接服务器异常")
		return
	}

	if request.Vendor == model.VendorMobile {
		storage.DeleteSMSToken(request.OpenId)
	}

	logrus.WithFields(logrus.Fields{
		"Request":           request,
		"GlobalId":          rs.Account.GlobalId,
		"AccessToken":       accessToken,
		"ExpirationSeconds": expirationSeconds,
		"RefreshToken":      refreshToken,
	}).Info("auth success.")

	// return result.
	resp := ResponseAuth{}
	resp.Body.Message = "Verify login successfully!"
	resp.Body.Code = int64(errors.Ok)
	resp.Body.AccessToken = accessToken
	resp.Body.ExpirationSeconds = expirationSeconds
	resp.Body.RefreshToken = refreshToken
	c.JSON(http.StatusOK, resp.Body)
}

func (o *Authenticator) createAccount(rs *RequestState) error {

	if rs.Account != nil {
		o.Provider.UpdateVendor(rs.Account.GlobalId, &rs.Request)
		return nil
	}

	if o.workingMode == RegisterClosed && !rs.IsInWhiteList {

		return errors.NewCodeString(errors.AuthRegisterClosed, "当前暂时无法注册")
	}

	// create new account
	var err error
	rs.Account, err = o.Provider.CreateUser(rs.Context, &rs.Request)
	if err != nil {
		return errors.NewCodeString(errors.ServerFail, "连接服务器异常")
	}

	switch o.workingMode {
	case Normal:
		// nothing need to do.

		//inRegisterByCouponChannels need coupon
		if inRegisterByCouponChannels(rs.Request.Channel) && !rs.IsInWhiteList && len(rs.Request.Coupon) == 0 {

			return errors.NewCodeString(errors.AuthRegisterByCoupon, "需要使用激活码才能注册")
		}
	case RegisterByCoupon:

		if !rs.IsInWhiteList && len(rs.Request.Coupon) == 0 {
			return errors.NewCodeString(errors.AuthRegisterByCoupon, "需要使用激活码才能注册")
		}
	case RegisterClosed:

		if !rs.IsInWhiteList {
			return errors.NewCodeString(errors.AuthRegisterClosed, "当前暂时无法注册")
		}
	default:

		logrus.WithFields(logrus.Fields{
			"workingMode": settings.GetString("loveauth", "auth.WorkingMode"),
		}).Error("Invalid WorkingMode.")

		return errors.NewCodeString(errors.ServerFail, "连接服务器异常")
	}

	if rs.Account == nil {
		return errors.NewCodeString(errors.ServerFail, "连接服务器异常")
	}

	return nil
}

func (o *Authenticator) filterByCoupon(rs *RequestState) error {

	var account = rs.Account
	var request = &rs.Request

	if rs.IsInWhiteList {
		return nil
	}

	if account.State != model.InActive {
		return nil
	}

	var code = request.Coupon
	if len(code) == 0 {
		return errors.NewCodeString(errors.AuthRegisterByCoupon, "需要使用激活码才能注册")
	}

	var db = storage.AuthDatabase()
	if db == nil {
		return errors.NewCodeString(errors.ServerFail, "连接服务器异常")
	}

	couponName := settings.GetString("loveauth", "auth.CouponName")

	var aCoupon = storage.Coupon{}
	aCoupon.SetName(couponName)

	var couponSalt = settings.GetString("loveauth", "coupon.Salt")

	var key = fmt.Sprintf("%x", md5.Sum([]byte(couponName+couponSalt)))
	var generator = couponLib.NewEncoding(key, 4, 1, 4, 4, 4)

	index, err := generator.DecodeString(code)
	if err != nil {
		return errors.NewCodeString(errors.CouponFailed, "请输入正确的激活码")
	}

	aCoupon.CouponId = uint32(index)
	if aCoupon.QueryCoupon() == nil {
		return errors.NewCodeString(errors.CouponFailed, "请输入正确的激活码")
	}

	zero := time.Time{}
	if aCoupon.Used != zero {
		return errors.NewCodeString(errors.CouponUsed, "该激活码已被使用，请重试")
	}

	if aCoupon.MarkCoupon() == nil {
		return errors.NewCodeString(errors.ServerFail, "连接服务器异常")
	}

	account.State = model.Active

	err = storage.Save(storage.AuthDatabase(), account)
	if err != nil {

		return errors.NewCodeString(errors.ServerFail, "连接服务器异常")
	}

	return nil
}

func (o *Authenticator) filterByState(rs *RequestState) error {

	var account = rs.Account

	accountState, desc := checkAccountBanState(account)
	if accountState == model.Banned || accountState == model.BanDied {
		return errors.NewCodeString(errors.AccountWasBanned, desc)
	}

	return nil
}

func checkAccountBanState(account *model.Account) (model.AccountState, string) {

	if account.State == model.Banned {

		currTime := time.Now().Unix()
		if account.UnBanTime > currTime {

			return model.Banned, "该账户由于多次违规，已被封停，如有疑问可以联系客服咨询"
		}

		storage.UpdateAccountState(account.GlobalId, map[string]interface{}{"state": 0, "un_ban_time": 0})

		return model.Active, "Account state active."
	}

	if account.State == model.BanDied {

		return model.BanDied, "该账户由于多次违规，已被封停，如有疑问可以联系客服咨询"
	}

	return model.Active, "Account state active."
}

func isWhiteListContains(request *RequestState) bool {
	userIp := net.ParseIP(request.AuthClientIp)
	if userIp == nil {
		return false
	}

	ipWhiteList := settings.GetStringSlice("loveauth_white_list", "ip")

	logrus.WithFields(logrus.Fields{
		"userIp":      userIp,
		"ipWhiteList": ipWhiteList,
	}).Info("isWhiteListContains")

	for _, ip := range ipWhiteList {
		if strings.Index(ip, "/") == -1 {

			ipAddr := net.ParseIP(ip)
			if userIp.Equal(ipAddr) {
				return true
			}
		} else {
			_, ipNet, err := net.ParseCIDR(ip)
			if err == nil && ipNet.Contains(userIp) {
				return true
			}
		}
	}

	versoinWhiteList := settings.GetStringSlice("loveauth_white_list", "version")

	logrus.WithFields(logrus.Fields{
		"ClientVersion":    request.Request.ClientVersion,
		"versoinWhiteList": versoinWhiteList,
	}).Info("isWhiteListContains")

	for _, version := range versoinWhiteList {

		if version == request.Request.ClientVersion {

			return true
		}
	}

	return false
}
