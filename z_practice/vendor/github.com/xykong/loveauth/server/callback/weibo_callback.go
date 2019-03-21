package callback

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/errors"
	"github.com/xykong/loveauth/storage"
	"net/http"
)

func init() {
	getHandlers["/weibo/callback"] = codeCB
}

//
// in: body
// swagger:parameters weibo_callback
type RequestCodeCB struct {
	//
	// swagger:allOf
	// in: query
	Code string `json:"code"`
}

//
// 应答: 协议返回包
// swagger:response ResponseCodeCB
// noinspection ALL
type ResponseCodeCB struct {
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
// swagger:route GET /weibo/callback weibo weibo_callback
//
//     Consumes:
//     - application/x-www-form-urlencoded
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       200: ResponseCodeCB
func codeCB(c *gin.Context) {
	code := c.Query("code")
	resp := ResponseCodeCB{}
	if code == "" {
		resp.Body.Code = int64(errors.Failed)
		resp.Body.Message = "code empty."
		c.JSON(http.StatusBadRequest, resp.Body)
		return
	}

	err := storage.SaveWeiboCode(code)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":  err,
			"code": code,
		}).Error("codeCB SaveWeiboCode failed.")

		resp.Body.Code = int64(errors.Failed)
		resp.Body.Message = "save code failed."
		c.JSON(http.StatusOK, resp.Body)
	}

	resp.Body.Code = int64(errors.Ok)
	resp.Body.Message = "success."
	c.JSON(http.StatusOK, resp.Body)
	//watch_waring.PaymentWatch(order.GlobalId, order.Vendor, order.Amount)
}
