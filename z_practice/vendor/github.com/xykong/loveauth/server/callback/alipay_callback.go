package callback

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/server/watch_waring"
	"github.com/xykong/loveauth/services/payment/alipay"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"net/http"
)

func init() {
	handlers["/alipay/callback"] = alipayCallback
	getHandlers["/alipay/callback"] = alipayCallback
}

//
// in: body
// swagger:parameters alipay_callback
type DoAlipayCallbackReqBodyParams struct {
	//
	// swagger:allOf
	// in: body
	DoAlipayCallbackReq model.DoAlipayCallbackReq
}

//
// swagger:route POST /alipay/callback callback alipay_callback
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
//       200:
func alipayCallback(c *gin.Context) {

	var request model.DoAlipayCallbackReq
	// validation
	if err := c.Bind(&request); err != nil {

		logrus.WithFields(logrus.Fields{
			"request": request,
			"error":   err,
		}).Error("alipay callback failed.")

		c.String(http.StatusOK, "fail")

		return
	}

	if ok, err := alipay.VerifySign(c.Request.Form); !ok {

		logrus.WithFields(logrus.Fields{
			"error":   err,
			"ok":      ok,
			"request": request,
		}).Error("alipay verifySign error.")

		c.String(http.StatusOK, "fail")

		return
	}

	if !(request.TradeStatus == alipay.K_TRADE_STATUS_TRADE_SUCCESS ||
		request.TradeStatus == alipay.K_TRADE_STATUS_TRADE_FINISHED) {

		logrus.WithFields(logrus.Fields{
			"request": request,
		}).Error("alipay trade status error.")

		c.String(http.StatusOK, "fail")

		return
	}

	order := storage.QueryOrderPlacedWithSequence(request.PassbackParams)
	if order == nil {

		logrus.WithFields(logrus.Fields{
			"request": request,
			"order":   order,
		}).Error("alipay callback failed.")

		c.String(http.StatusOK, "fail")

		return
	}

	if order.State >= model.OrderStatePlace {

		c.String(http.StatusOK, "success")

		return
	}

	order.State = model.OrderStatePlace
	order.SNSOrderId = request.TradeNo

	storage.Insert(storage.PayDatabase(), &request)

	//订单再处理
	err := storage.Save(storage.PayDatabase(), order)
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"request": request,
			"error":   err,
		}).Error("alipay callback failed.")

		c.String(http.StatusOK, "fail")

		return
	}

	c.String(http.StatusOK, "success")
	watch_waring.PaymentWatch(order.GlobalId, order.Vendor, order.Amount)
	return
}
