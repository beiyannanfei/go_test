package query

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xykong/loveauth/errors"
	"github.com/xykong/loveauth/services/msdk"
	"github.com/xykong/loveauth/settings"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"github.com/xykong/loveauth/utils"
	"net/http"
	"strings"
)

func init() {
	handlers["/profile"] = profile
}

// Binding from JSON
type RequestProfile struct {
	GlobalId int64 `form:"globalId" json:"globalId" binding:"required"`
}

//
// in: body
// swagger:parameters query_profile
type DoRequestProfileBodyParams struct {
	//
	// swagger:allOf
	// in: body
	Request RequestProfile `form:"Request" json:"Request" binding:"required"`
}

// A ResponseProfile is an response message to client.
// swagger:response ResponseProfile
type ResponseProfile struct {
	// in: body
	Body struct {
		// The response error code
		//
		// Required: true
		Code int64 `json:"code"`
		// The response message
		//
		// Required: true
		Message string `json:"message"`
		// 用户的昵称
		//
		// Required: true
		NickName string `json:"nickName"`
		// 用户头像URL
		//
		// Required: true
		Picture string `json:"picture"`
	}
}

// swagger:route POST /query/profile query query_profile
//
// Query profile retrieve from vendor.
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
//       200: ResponseProfile
func profile(c *gin.Context) {

	var request RequestProfile

	// validation
	if err := c.BindJSON(&request); err != nil {

		utils.QuickReply(c, errors.Failed, "profile BindJSON failed: %v", err)
		return
	}

	if request.GlobalId == 0 {

		utils.QuickReply(c, errors.Failed, "profile GlobalId is not valid.")
		return
	}

	profile, err := storage.QueryProfile(request.GlobalId)
	if err != nil {
		utils.QuickReply(c, errors.Failed, "profile QueryProfile failed: %v", err)
		return
	}

	var defaultPicture = settings.GetString("loveauth", "auth.DefaultProfilePicture")

	if len(profile.Picture) == 0 {
		profile.Picture = defaultPicture
	}

	if len(profile.Name) == 0 {
		profile.Name = utils.IdString(profile.GlobalId)
	}

	if len(profile.NickName) == 0 {
		profile.NickName = utils.IdString(profile.GlobalId)
	}

	resp := ResponseProfile{}
	resp.Body.Code = int64(errors.Ok)
	resp.Body.Message = "Query profile successfully!"

	switch profile.Vendor {
	case model.VendorDevice:
	case model.VendorMsdkGuest:
	case model.VendorMobile:
	case model.VendorQuickOppo: //quick-oppo无法获取渠道昵称
	case model.VendorDouyin: // 抖音无法获得渠道昵称
	case model.VendorMgtv: // 芒果tv
		authMgtv := model.AuthMgtv{}
		err := storage.QueryVendorInfoByGlobalId(request.GlobalId, &authMgtv)
		if err != nil {
			utils.QuickReply(c, errors.Failed, "profile VendorMgtv db failed: %v", err)
			return
		}

		nickName := strings.TrimSpace(authMgtv.NickName)
		if nickName != "" {
			profile.NickName = nickName
		}
	case model.VendorQuickAligames:
		authQuickAliGames := model.AuthQuickAliGames{}
		err := storage.QueryVendorInfoByGlobalId(request.GlobalId, &authQuickAliGames)
		if err != nil {
			utils.QuickReply(c, errors.Failed, "profile VendorQuickAligames db failed: %v", err)
			return
		}

		nickName := strings.TrimSpace(authQuickAliGames.UserName)
		if nickName != "" {
			profile.NickName = nickName
		}

	case model.VendorQuickM4399:
		authQuickM4399 := model.AuthQuickM4399{}
		err := storage.QueryVendorInfoByGlobalId(request.GlobalId, &authQuickM4399)
		if err != nil {
			utils.QuickReply(c, errors.Failed, "profile VendorQuickM4399 db failed: %v", err)
			return
		}

		nickName := strings.TrimSpace(authQuickM4399.UserName)
		if nickName != "" {
			profile.NickName = nickName
		}

	case model.VendorQuickYsdk:
		authQuickYsdk := model.AuthQuickYsdk{}
		err := storage.QueryVendorInfoByGlobalId(request.GlobalId, &authQuickYsdk)
		if err != nil {
			utils.QuickReply(c, errors.Failed, "profile VendorQuickYsdk db failed: %v", err)
			return
		}

		nickName := strings.TrimSpace(authQuickYsdk.UserName)
		if nickName != "" {
			profile.NickName = nickName
		}

	case model.VendorQuickIqiyi:
		authQuickIqiyi := model.AuthQuickIqiyi{}
		err := storage.QueryVendorInfoByGlobalId(request.GlobalId, &authQuickIqiyi)
		if err != nil {
			utils.QuickReply(c, errors.Failed, "profile VendorQuickIqiyi db failed: %v", err)
			return
		}

		nickName := strings.TrimSpace(authQuickIqiyi.UserName)
		if nickName != "" {
			profile.NickName = nickName
		}

	case model.VendorQuickMeiZu:
		authQuickMeiZu := model.AuthQuickMeizu{}
		err := storage.QueryVendorInfoByGlobalId(request.GlobalId, &authQuickMeiZu)
		if err != nil {
			utils.QuickReply(c, errors.Failed, "profile VendorQuickMeiZu db failed: %v", err)
			return
		}

		nickName := strings.TrimSpace(authQuickMeiZu.UserName)
		if nickName != "" {
			profile.NickName = nickName
		}

	case model.VendorQuickXiaomi:
		authQuickXiaomi := model.AuthQuickXiaomi{}
		err := storage.QueryVendorInfoByGlobalId(request.GlobalId, &authQuickXiaomi)
		if err != nil {
			utils.QuickReply(c, errors.Failed, "profile VendorQuickXiaomi db failed: %v", err)
			return
		}

		nickName := strings.TrimSpace(authQuickXiaomi.UserName)
		if nickName != "" {
			profile.NickName = nickName
		}

	case model.VendorQuickKuaiKan:
		authQuickKuaiKan := model.AuthQuickKuaikan{}
		err := storage.QueryVendorInfoByGlobalId(request.GlobalId, &authQuickKuaiKan)
		if err != nil {
			utils.QuickReply(c, errors.Failed, "profile VendorQuickKuaiKan db failed: %v", err)
			return
		}

		nickName := strings.TrimSpace(authQuickKuaiKan.UserName)
		if nickName != "" {
			profile.NickName = nickName
		}

	case model.VendorBilibili:
		authBilibili := model.AuthBilibili{}
		err := storage.QueryVendorInfoByGlobalId(request.GlobalId, &authBilibili)
		if err != nil {
			utils.QuickReply(c, errors.Failed, "profile VendorBilibili db failed: %v", err)
			return
		}

		nickName := strings.TrimSpace(authBilibili.UserName)
		if nickName != "" {
			profile.NickName = nickName
		}

	case model.VendorHuawei:
		authHuawei := model.AuthHuawei{}
		err := storage.QueryVendorInfoByGlobalId(request.GlobalId, &authHuawei)
		if err != nil {
			utils.QuickReply(c, errors.Failed, "profile VendorHuawei db failed: %v", err)
			return
		}

		nickName := strings.TrimSpace(authHuawei.DisplayName)
		if nickName != "" {
			profile.NickName = nickName
		}

	case model.VendorVivo:

		authVivo := model.AuthVivo{}
		err := storage.QueryVendorInfoByGlobalId(request.GlobalId, &authVivo)
		if err != nil {
			utils.QuickReply(c, errors.Failed, "profile VendorVivo db failed: %v", err)
			return
		}

		if authVivo.GlobalId == 0 {
			break
		}

		nickName := strings.TrimSpace(authVivo.NickName)
		if nickName != "" {
			profile.NickName = nickName
		}

	case model.VendorYsdkQQ: //ysdk形式在登陆时就返回昵称和头像，且已经直接入库
		authYsdkQQ := model.AuthYsdkQQ{}
		err := storage.QueryVendorInfoByGlobalId(request.GlobalId, &authYsdkQQ)
		if err != nil {
			utils.QuickReply(c, errors.Failed, "profile VendorYSDKQQ db failed: %v", err)
			return
		}

		if authYsdkQQ.GlobalId == 0 {
			break
		}

		nickName := strings.TrimSpace(authYsdkQQ.NickName)
		if nickName != "" {
			profile.NickName = nickName
		}

		picture := strings.TrimSpace(authYsdkQQ.Picture)
		if picture != "" {
			profile.Picture = picture
		}

	case model.VendorYsdkWechat:
		authYsdkWechat := model.AuthYsdkWechat{}
		err := storage.QueryVendorInfoByGlobalId(request.GlobalId, &authYsdkWechat)
		if err != nil {
			utils.QuickReply(c, errors.Failed, "profile authYsdkWechat db failed: %v", err)
			return
		}

		if authYsdkWechat.GlobalId == 0 {
			break
		}

		nickName := strings.TrimSpace(authYsdkWechat.NickName)
		if nickName != "" {
			profile.NickName = nickName
		}

		picture := strings.TrimSpace(authYsdkWechat.Picture)
		if picture != "" {
			profile.Picture = picture
		}

	case model.VendorMsdkQQ:

		info, err := msdk.QQProfile(profile.Auth.OpenId, profile.Auth.VendorMsdkQQ.AccessTokenValue)
		if err != nil {
			utils.QuickReply(c, errors.Failed, "profile QQProfile failed: %v", err)
			return
		}

		profileMap := make(map[string]interface{})
		nickName := strings.TrimSpace(info.Body.NickName)
		if nickName != "" {
			profile.NickName = nickName
			profileMap["nick_name"] = nickName
		}

		picture := strings.Replace(strings.TrimSpace(info.Body.Picture100), "http", "https", 1)
		if picture != "" {
			profile.Picture = picture
			profileMap["picture"] = picture
		}

		if len(profileMap) > 0 {
			var authWechat model.AuthWechat
			storage.UpdateSNSProfile(&authWechat, request.GlobalId, profileMap)
		}

	case model.VendorMsdkWechat:

		info, err := msdk.WXUserInfo(profile.Auth.OpenId, profile.Auth.VendorMsdkWechat.AccessTokenValue)
		if err != nil {
			utils.QuickReply(c, errors.Failed, "profile WXUserInfo failed: %v", err)
			return
		}

		profileMap := make(map[string]interface{})
		nickName := strings.TrimSpace(info.Body.Nickname)
		if nickName != "" {
			profile.NickName = nickName
			profileMap["nick_name"] = nickName
		}

		picture := strings.Replace(fmt.Sprintf("%v/0", strings.TrimSpace(info.Body.Picture)), "http", "https", 1)
		if picture != "" {
			profile.Picture = picture
			profileMap["picture"] = picture
		}

		if len(profileMap) > 0 {
			var authQQ model.AuthQQ
			storage.UpdateSNSProfile(&authQQ, request.GlobalId, profileMap)
		}

	case model.VendorWeibo:
		authWeibo := model.AuthWeibo{}
		err := storage.QueryVendorInfoByGlobalId(request.GlobalId, &authWeibo)
		if err != nil {
			utils.QuickReply(c, errors.Failed, "profile authWeibo db failed: %v", err)
			return
		}

		if authWeibo.GlobalId == 0 {
			break
		}

		nickName := strings.TrimSpace(authWeibo.NickName)
		if nickName != "" {
			profile.NickName = nickName
		}

		picture := strings.TrimSpace(authWeibo.Picture)
		if picture != "" {
			profile.Picture = picture
		}

	default:
		utils.QuickReply(c, errors.Failed, "profile vendor not support.")
		return
	}

	resp.Body.NickName = profile.NickName
	resp.Body.Picture = profile.Picture

	storage.WriteProfile(profile)

	c.JSON(http.StatusOK, resp.Body)
}
