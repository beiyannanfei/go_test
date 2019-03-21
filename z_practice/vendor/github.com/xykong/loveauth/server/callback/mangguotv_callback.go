package callback

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/server/watch_waring"
	"github.com/xykong/loveauth/services/payment/mangguotv"
	"github.com/xykong/loveauth/settings"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

func init() {
	handlers["/mgtv/callback"] = MangGuoTvCallBack
	getHandlers["/mgtv/callback"] = MangGuoTvCallBack
}

func MangGuoTvCallBack(c *gin.Context) {

	resp := &MgtvResponse{Result: "200", Msg: "success"}

	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {

		resp.Result = "201"
		resp.Msg = "read body err"

		c.JSON(http.StatusOK, resp)

		return
	}

	value, err := url.ParseQuery(string(reqBody))
	if err != nil {

		resp.Result = "201"
		resp.Msg = "parse query err"

		c.JSON(http.StatusOK, resp)

		return
	}

	appKey := settings.GetString("lovepay", "mgtv.appKey")
	if !mangguotv.CheckSign(value, appKey) {

		resp.Result = "201"
		resp.Msg = "check sign err"

		c.JSON(http.StatusOK, resp)

		return
	}

	request := model.MgtvCallBackRequest{
		Sign:            value.Get("sign"),
		Version:         value.Get("version"),
		Uuid:            value.Get("uuid"),
		BusinessOrderId: value.Get("business_order_id"),
		TradeStatus:     value.Get("trade_status"),
		TradeCreate:     value.Get("trade_create"),
		TotalFee:        value.Get("total_fee"),
		ExtData:         value.Get("ext_data"),
	}

	order := storage.QueryOrderPlacedWithSequence(request.ExtData)
	if order == nil {

		logrus.WithFields(logrus.Fields{
			"request": request,
			"order":   order,
		}).Error("mgtv callback failed.")

		resp.Result = "201"
		resp.Msg = "order nil"

		c.JSON(http.StatusOK, resp)

		return
	}

	amount, _ := strconv.Atoi(request.TotalFee)
	if order.Amount != amount {

		logrus.WithFields(logrus.Fields{
			"request": request,
			"order":   order,
		}).Error("mgtv callback failed amount not equal.")

		resp.Result = "201"
		resp.Msg = "amount not equal err"

		c.JSON(http.StatusOK, resp)

		return
	}

	mgtvCallBack := storage.QueryMgtvOrder(request.BusinessOrderId)
	if mgtvCallBack != nil {

		c.JSON(http.StatusOK, resp)

		return
	}

	storage.Insert(storage.PayDatabase(), &request)

	if order.State >= model.OrderStatePlace {

		c.JSON(http.StatusOK, resp)

		return
	}

	order.SNSOrderId = request.BusinessOrderId
	order.State = model.OrderStatePlace

	err = storage.Save(storage.PayDatabase(), order)
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"request": request,
			"error":   err,
		}).Error("mgtv callback failed.")

		resp.Result = "201"
		resp.Msg = "save order err"

		c.JSON(http.StatusOK, resp)

		return
	}

	c.JSON(http.StatusOK, resp)
	watch_waring.PaymentWatch(order.GlobalId, order.Vendor, order.Amount)
}

type MgtvResponse struct {
	Result string `json:"result"`
	Msg    string `json:"msg"`
}
