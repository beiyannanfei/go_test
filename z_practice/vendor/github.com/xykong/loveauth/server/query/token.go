package query

import (
	"github.com/gin-gonic/gin"
	"github.com/xykong/loveauth/errors"
	"github.com/xykong/loveauth/server/auth"
	"github.com/xykong/loveauth/settings"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/utils"
	"net/http"
)

func init() {
	handlers["/token"] = token
}

// Binding from JSON
type RequestQueryTokenData struct {
	Token      string `form:"token" json:"token" binding:"required"`
	Expiration bool   `form:"expiration" json:"expiration"`
}

//
// in: body
// swagger:parameters query_token
type DoQueryTokenRequestBodyParams struct {
	//
	// swagger:allOf
	// in: body
	Request RequestQueryTokenData `form:"Request" json:"Request" binding:"required"`
}

// A ResponseQueryToken is an response message to client.
// swagger:response ResponseQueryToken
type ResponseQueryToken struct {
	// in: body
	Body struct {
		Code              int64  `json:"code"`
		Message           string `json:"message"`
		GlobalId          int64  `json:"globalId"`
		ExpirationSeconds int64  `json:"expiration_seconds"`
	}
}

// swagger:route POST /query/token query query_token
//
// Query token received from client.
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
//       200: ResponseQueryToken
func token(c *gin.Context) {

	var request RequestQueryTokenData

	// validation
	if err := c.BindJSON(&request); err != nil {

		utils.QuickReply(c, errors.Failed, "query_token BindJSON failed: %v", err)
		return
	}

	if request.Token == "" {

		utils.QuickReply(c, errors.Failed, "query_token Token is empty.")
		return
	}

	tokenRecord, err := storage.QueryAccessToken(request.Token)
	if err != nil {

		code := errors.Failed
		if ec, ok := err.(*errors.Type); ok {
			code = ec.Code
		}

		utils.QuickReply(c, code, "query_token QueryAccessToken failed: %v", err)
		return
	}

	if tokenRecord.GlobalId == 0 {

		utils.QuickReply(c, errors.Failed, "query_token globalId is invalid.")
		return
	}

	// process login and logout tlog message.
	err = auth.LogLongTimeSession(tokenRecord)
	if err != nil {

		utils.QuickReply(c, errors.Failed, "query_token LogLongTimeSession failed: %v", err)
		return
	}

	storage.TouchAccount(tokenRecord)

	resp := ResponseQueryToken{}
	resp.Body.Code = int64(errors.Ok)
	resp.Body.Message = "QueryToken token successfully!"
	resp.Body.GlobalId = tokenRecord.GlobalId

	if request.Expiration {
		resp.Body.ExpirationSeconds = settings.GetInt64("loveauth", "token.ExpirationSeconds")
	}

	c.JSON(http.StatusOK, resp.Body)
}
