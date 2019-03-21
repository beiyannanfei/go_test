package callback

import (
	"github.com/xykong/loveauth/server/watch_waring"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
)

func init() {
	handlers["/midas/callback"] = midasCallback
}

//
// in: body
// swagger:parameters payment_midas_callback
type DoMidasCallbackReqBodyParams struct {
	//
	// swagger:allOf
	// in: body
	DoMidasCallbackReq model.DoMidasCallbackReq
}

//
// 应答: 协议返回包
// swagger:response DoMidasCallbackRsp
// noinspection ALL
type DoMidasCallbackRsp struct {
	// in: body
	Body struct {
		//
		// 应用的错误码应该从0开始，按照整数递增的方式进行定义，建议应用按照如下描述定义错误码：
		//
		//	0: 成功
		//
		//	1: 系统繁忙
		//
		//	2: token已过期
		//
		//	3: token不存在
		//
		//	4: 请求参数错误：（这里填写错误的具体参数）
		//
		Ret int `json:"ret"`
		//
		// 道具发放操作的结果，成功为“OK”，失败则表明错误原因（必须使用utf8编码）。
		//
		// 腾讯设置的调用开发者发货超时是2秒钟，请开发者注意超时时间设置不要超过2秒，
		//
		// 否则腾讯后台将返回“系统繁忙”的错误消息。
		//
		Msg string `json:"msg"`
	}
}

//
// swagger:route POST /payment/midas/callback payment payment_midas_callback
//
// 回调发货协议说明:
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
//       200: DoMidasCallbackRsp
func midasCallback(c *gin.Context) {

	resp := DoMidasCallbackRsp{}

	var request model.DoMidasCallbackReq
	// validation
	if err := c.BindJSON(&request); err != nil {
		resp.Body.Ret = 4
		resp.Body.Msg = err.Error()

		logrus.WithFields(logrus.Fields{
			"Request":  request,
			"Response": resp,
			"Error":    err,
		}).Error("Midas callback failed.")

		c.JSON(http.StatusOK, resp.Body)
		return
	}

	resp.Body.Ret = 0
	resp.Body.Msg = "ok"

	logrus.WithFields(logrus.Fields{
		"Request":  request,
		"Response": resp,
	}).Info("Midas callback success.")

	db := storage.PayDatabase()
	if db == nil {

		resp.Body.Ret = 1
		resp.Body.Msg = "database unavailable."

		logrus.WithFields(logrus.Fields{
			"Request":  request,
			"Response": resp,
		}).Error("Midas callback failed.")

		c.JSON(http.StatusOK, resp.Body)
	}

	storage.Insert(db, &request)

	order := storage.QueryOrderPlacedWithSequence(request.Appmeta)
	if order == nil {

		resp.Body.Ret = 4
		resp.Body.Msg = "appmeta err"

		logrus.WithFields(logrus.Fields{
			"Request":  request,
			"Response": resp,
			"Error":    "appmeta err",
		}).Error("Midas callback failed.")

		c.JSON(http.StatusOK, resp.Body)
		return
	}

	order.State = model.OrderStatePlace
	err := storage.Save(storage.PayDatabase(), &order)
	if err != nil {

		resp.Body.Ret = 1
		resp.Body.Msg = "order save err."

		logrus.WithFields(logrus.Fields{
			"Request":  request,
			"Response": resp,
			"Error":    err,
		}).Error("Midas callback failed.")

		c.JSON(http.StatusOK, resp.Body)
		return
	}

	c.JSON(http.StatusOK, resp.Body)
	watch_waring.PaymentWatch(order.GlobalId, order.Vendor, order.Amount)
}
