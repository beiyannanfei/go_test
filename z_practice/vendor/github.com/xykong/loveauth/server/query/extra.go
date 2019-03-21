package query

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/utils"
	"github.com/xykong/loveauth/storage/model"
	"github.com/xykong/loveauth/errors"
	"github.com/xykong/loveauth/server/auth"
)

// Binding from JSON
type RequestExtraData struct {
	//Extra string `form:"extra" json:"extra" binding:"required"`
	GlobalId int64 `form:"globalId" json:"globalId" binding:"required"`
	//DeviceId              string `form:"deviceId" json:"deviceId" binding:"optional"`
}

func init() {
	handlers["/extra"] = extra
}

//
// in: body
// swagger:parameters query_extra
type DoQueryRequestBodyParams struct {
	//
	// swagger:allOf
	// in: body
	Request RequestExtraData `form:"Request" json:"Request" binding:"required"`
}

// A ResponseQuery is an response message to client.
// swagger:response ResponseQuery
type ResponseQuery struct {
	// in: body
	Body struct {
		Code    int64  `json:"code"`
		Message string `json:"message"`
		// OpenId
		//
		// Required: true
		OpenId string `json:"openId"`
		// Vendor
		//
		// Required: true
		Vendor model.Vendor `json:"vendor"`
		// swagger:allOf
		// in: body
		Extra model.Extra `form:"extra" json:"extra" binding:"required"`
	}
}

// swagger:route POST /query/extra query query_extra
//
// Query extra received from client.
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
//       200: ResponseQuery
func extra(c *gin.Context) {

	var request RequestExtraData

	// validation
	if err := c.BindJSON(&request); err != nil {

		utils.QuickReply(c, errors.Failed, "BindJSON failed: %v", err)
		return
	}

	if request.GlobalId == 0 {

		utils.QuickReply(c, errors.Failed, "GlobalId is not valid.")
		return
	}

	profile, err := storage.QueryProfile(request.GlobalId)
	if err != nil {
		utils.QuickReply(c, errors.Failed, "extra QueryProfile failed: %v", err)
		return
	}

	gameSveId, gameAppId := auth.GetVendorConfigSetting(profile.Vendor)
	extra := model.Extra{}
	extra.GameSvrId = gameSveId
	extra.VGameAppid = gameAppId
	extra.PlatId = 2
	if profile.Platform == model.PlatformIOS {
		extra.PlatId = 0
	}
	if profile.Platform == model.PlatformAndroid {
		extra.PlatId = 1
	}
	extra.Vopenid = profile.Auth.OpenId
	extra.ClientVersion = profile.Auth.ClientVersion
	extra.RegChannel = profile.Auth.Channel
	extra.LoginChannel = profile.Auth.Channel
	extra.DeviceId = profile.Auth.DeviceId
	extra.VClientIP = profile.Auth.ClientIp

	resp := ResponseQuery{}
	resp.Body.Code = int64(errors.Ok)
	resp.Body.Message = "Query extra successfully!"
	resp.Body.OpenId = profile.Auth.OpenId
	resp.Body.Vendor = profile.Vendor
	resp.Body.Extra = extra
	c.JSON(http.StatusOK, resp.Body)
}
