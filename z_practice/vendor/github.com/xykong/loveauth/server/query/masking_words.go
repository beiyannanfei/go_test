package query

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xykong/loveauth/errors"
	"github.com/xykong/loveauth/services/dict"
	"github.com/xykong/loveauth/utils"
)

func init() {
	handlers["/masking_words"] = masking_words
}

// Binding from JSON
type DoMaskingWordsReq struct {
	Text string `form:"text" json:"text" binding:"required"`
}

//
// in: body
// swagger:parameters query_masking_words
type DoMaskingWordsReqBodyParams struct {
	//
	// swagger:allOf
	// in: body
	DoMaskingWordsReq DoMaskingWordsReq
}

// A DoMaskingWordsRsp is an response message to client.
// swagger:response DoMaskingWordsRsp
type DoMaskingWordsRsp struct {
	// in: body
	Body struct {
		Code    int64  `json:"code"`
		Message string `json:"message"`
		// 过滤后字符串
		//
		// Required: true
		Text string `json:"text"`
		// 是否存在屏蔽词
		//
		// Required: true
		Exist bool `json:"exist"`
		// 涉及到的屏蔽词库
		//
		// Required: true
		DirtyWords []string `json:"dirty_words"`
	}
}

// swagger:route POST /query/masking_words query query_masking_words
//
// masking words.
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
//       200: DoMaskingWordsRsp
func masking_words(c *gin.Context) {

	var request DoMaskingWordsReq

	// validation
	if err := c.BindJSON(&request); err != nil {

		utils.QuickReply(c, errors.Failed, "BindJSON failed: %v", err)
		return
	}

	text, exist, dirtyWords := dict.ReplaceInvalidWords(request.Text)

	resp := DoMaskingWordsRsp{}
	resp.Body.Code = int64(errors.Ok)
	resp.Body.Message = "successfully!"
	resp.Body.Text = text
	resp.Body.Exist = exist
	resp.Body.DirtyWords = dirtyWords

	c.JSON(http.StatusOK, resp.Body)
}
