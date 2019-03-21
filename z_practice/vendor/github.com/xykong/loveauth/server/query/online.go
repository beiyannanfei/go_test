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
	handlers["/online"] = online
}

// Binding from JSON
type DoQueryOnlineReq struct {
	Channel  string   `form:"channel" json:"channel" binding:"required"`
	Platform model.Platform `form:"platform" json:"platform" binding:""`
}

//
// in: body
// swagger:parameters query_online
type DoQueryOnlineReqBodyParams struct {
	//
	// swagger:allOf
	// in: body
	Request DoQueryOnlineReq `form:"Request" json:"Request" binding:"required"`
}

// A DoQueryOnlineRsp is an response message to client.
// swagger:response DoQueryOnlineRsp
type DoQueryOnlineRsp struct {
	// in: body
	Body struct {
		Code    int64                                   `json:"code"`
		Message string                                  `json:"message"`
		Summary map[string]map[model.Platform]int `form:"channels" json:"channels" binding:"required"`
	}
}

// swagger:route POST /query/online query query_online
//
// Query online count by specified vendor
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
//       200: DoQueryOnlineRsp
func online(c *gin.Context) {

	var request DoQueryOnlineReq

	// validation
	if err := c.BindJSON(&request); err != nil {

		utils.QuickReply(c, errors.Failed, "BindJSON failed: %v", err)
		return
	}

	platforms := []model.Platform{request.Platform}
	if request.Platform == model.PlatformAll {
		platforms = model.Platforms
	}

	channels := []string{request.Channel}
	//vendors := []model.Vendor{request.Channel}
	if request.Channel == model.ChannelAll {
		channels = storage.GetChannelList()
	}

	resp := DoQueryOnlineRsp{}
	resp.Body.Summary = make(map[string]map[model.Platform]int)
	for _, channel := range channels {

		resp.Body.Summary[channel] = make(map[model.Platform]int)

		for _, platform := range platforms {

			count := storage.CountOnlineToken(channel, platform)

			resp.Body.Summary[channel][platform] = count
		}
	}

	resp.Body.Code = int64(errors.Ok)
	resp.Body.Message = "Query online successfully!"
	c.JSON(http.StatusOK, resp.Body)
}
