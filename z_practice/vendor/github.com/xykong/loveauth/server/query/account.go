package query

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/utils"
	"github.com/xykong/loveauth/storage/model"
	"github.com/xykong/loveauth/errors"
)

func init() {
	handlers["/account"] = account
}

// Binding from JSON
type DoQueryAccountReq struct {
	OpenId string       `form:"openId" json:"openId" binding:"required"`
	Vendor model.Vendor `form:"vendor" json:"vendor" binding:"required"`
}

//
// in: body
// swagger:parameters query_account
type DoQueryAccountReqBodyParams struct {
	//
	// swagger:allOf
	// in: body
	Request DoQueryAccountReq `form:"Request" json:"Request" binding:"required"`
}

// A DoQueryAccountRsp is an response message to client.
// swagger:response DoQueryAccountRsp
type DoQueryAccountRsp struct {
	// in: body
	Body struct {
		Code    int64  `json:"code"`
		Message string `json:"message"`

		// swagger:allOf
		// in: body
		Account model.Account `form:"account" json:"account" binding:"required"`
	}
}

// swagger:route POST /query/account query query_account
//
// Query account received from client.
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
//       200: DoQueryAccountRsp
func account(c *gin.Context) {

	var request DoQueryAccountReq

	// validation
	if err := c.BindJSON(&request); err != nil {

		utils.QuickReply(c, errors.Failed, "BindJSON failed: %v", err)
		return
	}

	if len(request.OpenId) == 0 {

		utils.QuickReply(c, errors.Failed, "OpenId is not valid.")
		return
	}

	account := storage.QueryAccountByOpenId(request.OpenId, request.Vendor)
	if account == nil {
		utils.QuickReply(c, errors.Failed, "account not found.")
		return
	}

	resp := DoQueryAccountRsp{}
	resp.Body.Code = int64(errors.Ok)
	resp.Body.Message = "Query account successfully!"
	resp.Body.Account = *account
	c.JSON(http.StatusOK, resp.Body)
}
