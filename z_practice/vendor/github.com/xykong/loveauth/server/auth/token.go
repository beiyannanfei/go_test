package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/xykong/loveauth/storage"
	"net/http"
	"github.com/xykong/loveauth/settings"
	"github.com/sirupsen/logrus"
	"time"
	"github.com/xykong/loveauth/errors"
)

/*
curl -H "Content-Type: application/json" -X POST \
        -d '{
		"grantType" : "refresh_token",
        "refreshToken": "e00630ee6c8e9e8652996cce7268d1de",
        "GlobalId": 1515915429
        }' \
        http://localhost:8080/api/v1/auth/token
 */

// Binding from JSON
type DoTokenRequest struct {
	GrantType    string `form:"GrantType" json:"GrantType" binding:"required"`
	RefreshToken string `form:"RefreshToken" json:"RefreshToken" binding:"required"`
	GlobalId     int64  `form:"GlobalId" json:"GlobalId" binding:"required"`
}

//
// in: body
// swagger:parameters auth_token
type DoTokenRequestBodyParams struct {
	//
	// swagger:allOf
	// in: body
	Request DoTokenRequest `form:"Request" json:"Request" binding:"required"`
}

// swagger:route POST /auth/token auth auth_token
//
// Refresh access token by refresh token
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
func token(c *gin.Context) {

	var request DoTokenRequest

	// validation
	if err := c.BindJSON(&request); err != nil {
		resp := ResponseAuth{}
		resp.Body.Message = err.Error()
		resp.Body.Code = int64(errors.Failed)
		c.JSON(http.StatusOK, resp.Body)
		return
	}

	if request.GrantType != "refresh_token" {
		resp := ResponseAuth{}
		resp.Body.Message = "GrantType is invalid."
		resp.Body.Code = int64(errors.Failed)
		c.JSON(http.StatusOK, resp.Body)
		return
	}

	record, err := storage.QueryRefreshToken(request.RefreshToken)

	if err != nil || record == nil {
		resp := ResponseAuth{}
		resp.Body.Message = "Refresh token failed!"
		resp.Body.Code = int64(errors.Failed)
		c.JSON(http.StatusOK, resp.Body)
		return
	}

	var expirationSeconds = settings.GetInt64("loveauth", "token.ExpirationSeconds")

	record.Timestamp = time.Now().Unix()
	record.ExpirationSeconds = expirationSeconds

	err = storage.RefreshToken(record)
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"Request":           request,
			"GlobalId":          request.GlobalId,
			"AccessToken":       record.Token,
			"ExpirationSeconds": expirationSeconds,
			"RefreshToken":      request.RefreshToken,
			"err":               err,
		}).Info("refresh token failed.")

		resp := ResponseAuth{}
		resp.Body.Message = "Refresh token failed!"
		resp.Body.Code = int64(errors.Failed)
		c.JSON(http.StatusOK, resp.Body)

		return
	}

	logrus.WithFields(logrus.Fields{
		"Request":           request,
		"GlobalId":          request.GlobalId,
		"AccessToken":       record.Token,
		"ExpirationSeconds": expirationSeconds,
		"RefreshToken":      request.RefreshToken,
	}).Info("refresh token success.")

	resp := ResponseAuth{}
	resp.Body.Message = "Refresh token successfully!"
	resp.Body.Code = int64(errors.Ok)
	resp.Body.AccessToken = record.Token
	resp.Body.ExpirationSeconds = expirationSeconds
	resp.Body.RefreshToken = request.RefreshToken
	c.JSON(http.StatusOK, resp.Body)
}
