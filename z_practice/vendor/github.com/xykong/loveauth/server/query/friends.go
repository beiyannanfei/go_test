package query

import (
	"github.com/gin-gonic/gin"
	"github.com/xykong/loveauth/errors"
	"github.com/xykong/loveauth/services/msdk"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"github.com/xykong/loveauth/utils"
	"net/http"
	"strings"
)

func init() {
	handlers["/friends"] = friends
}

// Binding from JSON
type RequestFriends struct {
	GlobalId int64 `form:"globalId" json:"globalId" binding:"required"`
}

//
// in: body
// swagger:parameters query_friends
type DoRequestFriendsBodyParams struct {
	//
	// swagger:allOf
	// in: body
	Request RequestFriends `form:"Request" json:"Request" binding:"required"`
}

type Friend struct {
	// GlobalId
	//
	// Required: true
	GlobalId int64 `json:"globalId"`
	// 用户的昵称
	//
	// Required: true
	NickName string `json:"nickName"`
	// 用户头像URL
	//
	// Required: true
	Picture string `json:"picture"`
}

// A ResponseFriends is an response message to client.
// swagger:response ResponseFriends
type ResponseFriends struct {
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
		// 同玩好友信息
		//
		// Required: true
		// swagger:allOf
		Friends []Friend `json:"friends"`
	}
}

// swagger:route POST /query/friends query query_friends
//
// Query friends retrieve from vendor.
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
//       200: ResponseFriends
func friends(c *gin.Context) {

	var request RequestFriends

	// validation
	if err := c.BindJSON(&request); err != nil {

		utils.QuickReply(c, errors.Failed, "friends BindJSON failed: %v", err)
		return
	}

	if request.GlobalId == 0 {

		utils.QuickReply(c, errors.Failed, "friends GlobalId is not valid.")
		return
	}

	profile, err := storage.QueryProfile(request.GlobalId)
	if err != nil {
		utils.QuickReply(c, errors.Failed, "friends QueryProfile failed: %v", err)
		return
	}

	resp := ResponseFriends{}
	resp.Body.Code = int64(errors.Ok)
	resp.Body.Message = "Query friends successfully!"

	switch profile.Vendor {
	case model.VendorDevice:
	case model.VendorMsdkGuest:
	case model.VendorWeibo:
	case model.VendorYsdkQQ:
	case model.VendorYsdkWechat:
	case model.VendorVivo:
	case model.VendorBilibili:
	case model.VendorHuawei:
	case model.VendorMobile:
	case model.VendorQuickAligames:
	case model.VendorQuickOppo:
	case model.VendorQuickM4399:
	case model.VendorQuickYsdk:
	case model.VendorQuickIqiyi:
	case model.VendorQuickMeiZu:
	case model.VendorQuickKuaiKan:

	case model.VendorMsdkQQ:

		info, err := msdk.QQFriendsDetail(profile.Auth.OpenId, profile.Auth.VendorMsdkQQ.AccessTokenValue)
		if err != nil {
			utils.QuickReply(c, errors.Failed, "friends QQFriendsDetail failed: %v", err)
			return
		}

		var openIds []string
		for _, item := range info.Body.Lists {
			openIds = append(openIds, item.OpenId)
		}

		var authQQs []model.AuthQQ
		err = storage.QueryVendorByOpenIds(openIds, &authQQs)
		if err != nil {
			utils.QuickReply(c, errors.Failed, "friends QueryVendorByOpenIds failed: %v", err)
			return
		}

		var openIdMap = make(map[string]model.AuthQQ)
		for _, item := range authQQs {
			openIdMap[item.OpenId] = item
		}

		if err != nil {
			utils.QuickReply(c, errors.Failed, "friends QueryVendorByOpenIds failed: %v", err)
			return
		}

		for _, item := range info.Body.Lists {

			if val, ok := openIdMap[item.OpenId]; ok {
				//noinspection ALL
				resp.Body.Friends = append(resp.Body.Friends, Friend{
					GlobalId: val.GlobalId,
					NickName: item.NickName,
					Picture:  strings.Replace(item.FigureUrlQQ, "http", "https", 1),
				})
			}
		}

	case model.VendorMsdkWechat:

		info, err := msdk.WXFriendsProfile(profile.Auth.OpenId, profile.Auth.VendorMsdkWechat.AccessTokenValue)
		if err != nil {
			utils.QuickReply(c, errors.Failed, "friends WXFriendsProfile failed: %v", err)
			return
		}

		var openIds []string
		for _, item := range info.Body.Lists {
			openIds = append(openIds, item.OpenId)
		}

		var authWechats []model.AuthWechat
		err = storage.QueryVendorByOpenIds(openIds, &authWechats)
		if err != nil {
			utils.QuickReply(c, errors.Failed, "friends QueryVendorByOpenIds failed: %v", err)
			return
		}

		var openIdMap = make(map[string]model.AuthWechat)
		for _, item := range authWechats {
			openIdMap[item.OpenId] = item
		}

		if err != nil {
			utils.QuickReply(c, errors.Failed, "friends QueryVendorByOpenIds failed: %v", err)
			return
		}

		for _, item := range info.Body.Lists {
			if val, ok := openIdMap[item.OpenId]; ok {
				//noinspection ALL
				resp.Body.Friends = append(resp.Body.Friends, Friend{
					GlobalId: val.GlobalId,
					NickName: item.NickName,
					Picture:  strings.Replace(item.Picture, "http", "https", 1),
				})
			}
		}

	default:
		utils.QuickReply(c, errors.Failed, "friends vendor not support.")
		return
	}

	c.JSON(http.StatusOK, resp.Body)
}
