package wjx

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/errors"
	"github.com/xykong/loveauth/server/idip"
	"github.com/xykong/loveauth/utils"
	"net/http"
	"strconv"
	"strings"
)

func init() {

	postHandlers["/callback"] = callback
}

//
// Request Buy Item
//
//noinspection SpellCheckingInspection
type DoCallbackReq struct {
	Activity      string `json:"activity"`
	Name          string `json:"name"`
	PartnerUser   string `json:"parteruser"`
	PartnerJoiner string `json:"parterjoiner"`
	TimeTaken     string `json:"timetaken"`
	SubmitTime    string `json:"submittime"`
	TotalValue    string `json:"totalvalue"`
	SoJumpParm    string `json:"sojumpparm"`
	Index         string `json:"index"`
	JoinId        string `json:"joinid"`
}

//
// in: body
// swagger:parameters wjx_callback
type DoCallbackReqBodyParams struct {
	//
	// swagger:allOf
	// in: body
	DoCallbackReq DoCallbackReq
}

//
// 应答: 协议返回包
// swagger:response DoCallbackRsp
// noinspection ALL
type DoCallbackRsp struct {
	// in: body
	Body struct {
		// The response code
		//
		// Required: true
		Code int64 `json:"code"`
		// The response message
		//
		// Required: true
		Message string `json:"message"`
	}
}

//
// swagger:route POST /wjx/callback wjx wjx_callback
//
// 问卷星数据推送API
// https://www.wjx.cn/help/help.aspx?catid=35
// http://www.wjx.cn/help/help.aspx?helpid=377&h=1#%E6%95%B0%E6%8D%AE%E6%8E%A8%E9%80%81api
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
//       200: DoCallbackRsp
func callback(c *gin.Context) {
	var request DoCallbackReq

	// validation
	if err := c.BindJSON(&request); err != nil {

		utils.QuickReply(c, errors.Failed, "callback BindJSON failed: %v", err)
		return
	}
	logrus.WithFields(logrus.Fields{
		"request": request,
	}).Info("wjx callback.")

	resp := DoCallbackRsp{}
	resp.Body.Code = int64(errors.Ok)
	resp.Body.Message = "callback successfully!"

	c.JSON(http.StatusOK, resp.Body)

	temp := request.SoJumpParm
	infoArr := strings.Split(temp, "_")
	areaId, _ := strconv.Atoi(infoArr[1])
	platId, _ := strconv.Atoi(infoArr[2])

	for count := 0; count <= idip.RandReqCount; count++ {

		url := randGames(uint32(areaId), uint8(platId))
		logrus.WithFields(logrus.Fields{
			"url": url,
		}).Info("wjx callback.")

		_, err := idip.PostRequest(url, request)
		if err != nil {

			if count < idip.RandReqCount {

				continue
			}

			logrus.WithFields(logrus.Fields{
				"err":err,
			}).Error("wjx callback")
			return
		}
		logrus.WithFields(logrus.Fields{
			"url":url,
		}).Info("wjx callback success.")
		return
	}
}
