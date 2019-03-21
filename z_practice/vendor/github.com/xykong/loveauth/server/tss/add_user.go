package tss

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/xykong/loveauth/utils"
	"github.com/xykong/loveauth/errors"
)

func init() {
	handlers["/add_user"] = addUser
}

// Binding from JSON
type DoTssAddUserReq struct {
	//
	// OpenId
	//
	// Required: true
	OpenId string `form:"OpenId" json:"OpenId" binding:"required"`
	//
	// 用户设备类型
	//
	// Required: true
	PlatId int `form:"PlatID" json:"PlatID" binding:"required"`
	//
	// 大区id
	//
	// Required: true
	WorldId int `form:"WorldId" json:"WorldId" binding:"required"`
	//
	// 登录的角色id
	//
	GlobalId int
	//
	// 游戏客户端的版本
	//
	ClientVersion string
	//
	// 游戏客户端ip，为网络字节序，必须填入
	//
	ClientIP string
	//
	// 用户当前的角色名，可选
	//
	Name string
}

//
// in: body
// swagger:parameters tss_add_user
type DoTssAddUserReqBodyParams struct {
	//
	// swagger:allOf
	// in: body
	Request DoTssAddUserReq `form:"Request" json:"Request" binding:"required"`
}

// A DoTssAddUserRsp is an response message to client.
// swagger:response DoTssAddUserRsp
type DoTssAddUserRsp struct {
	// in: body
	Body struct {
		Code    int64  `json:"code"`
		Message string `json:"message"`
	}
}

// swagger:route POST /tss/add_user tss tss_add_user
//
// TSS add_user received from client.
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
//       200: DoTssAddUserRsp
func addUser(c *gin.Context) {

	var request DoTssAddUserReq

	// validation
	if err := c.BindJSON(&request); err != nil {

		utils.QuickReply(c, errors.Failed, "BindJSON failed: %v", err)
		return
	}

	if len(request.OpenId) == 0 {
		utils.QuickReply(c, errors.Failed, "OpenId is not valid.")
		return
	}

	//add_user := storage.TssAddUserByOpenId(request.OpenId, request.Vendor)
	//if add_user == nil {
	//	utils.QuickReply(c, errors.Failed, "add_user not found.")
	//	return
	//}

	resp := DoTssAddUserRsp{}
	resp.Body.Code = int64(errors.Ok)
	resp.Body.Message = "Tss add_user successfully!"
	c.JSON(http.StatusOK, resp.Body)
}
