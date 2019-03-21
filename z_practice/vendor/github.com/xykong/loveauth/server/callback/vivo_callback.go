package callback

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/server/watch_waring"
	"github.com/xykong/loveauth/services/payment/vivo"
	"github.com/xykong/loveauth/settings"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"io/ioutil"
	"net/http"
	"net/url"
)

func init() {
	handlers["/vivo/callback"] = VivoCallBack
	getHandlers["/vivo/callback"] = VivoCallBack
}

func VivoCallBack(c *gin.Context) {

	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	urlValue, err := url.ParseQuery(string(bodyBytes))

	request := model.VivoQueryOrderResponse{
		Ret:           urlValue.Get("respCode"),
		Message:       urlValue.Get("respMsg"),
		SignMethod:    urlValue.Get("signMethod"),
		Signature:     urlValue.Get("signature"),
		TradeType:     urlValue.Get("tradeType"),
		TradeStatus:   urlValue.Get("tradeStatus"),
		CpId:          urlValue.Get("cpId"),
		AppId:         urlValue.Get("appId"),
		Uid:           urlValue.Get("uid"),
		CpOrderNumber: urlValue.Get("cpOrderNumber"),
		OrderNumber:   urlValue.Get("orderNumber"),
		OrderAmount:   urlValue.Get("orderAmount"),
		ExtInfo:       urlValue.Get("extInfo"),
		PayTime:       urlValue.Get("payTime"),
	}

	if request.Ret != "200" || request.TradeStatus != "0000" {

		logrus.WithFields(logrus.Fields{
			"request": request,
		}).Error("vivo callback failed.")

		c.String(http.StatusOK, "fail")

		return
	}

	requestBody, err := json.Marshal(request)
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"err":     err,
			"request": request,
		}).Error("vivo callback failed.")

		c.String(http.StatusOK, "fail")

		return
	}

	var params map[string]interface{}
	err = json.Unmarshal(requestBody, &params)
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"error":   err,
			"request": request,
		}).Error("vivo callback params failed.")

		c.String(http.StatusOK, "fail")

		return
	}

	key := settings.GetString("lovepay", "vivo.key")
	if !vivo.CheckSign(params, key) {

		logrus.WithFields(logrus.Fields{
			"params": params,
		}).Error("vivo callback check sign failed.")

		c.String(http.StatusOK, "fail")

		return
	}

	order := storage.QueryOrderPlacedWithSequence(request.ExtInfo)
	if order == nil {

		logrus.WithFields(logrus.Fields{
			"request": request,
			"order":   order,
		}).Error("vivo callback failed.")

		c.String(http.StatusOK, "fail")

		return
	}

	if order.State >= model.OrderStatePlace {

		c.String(http.StatusOK, "success")

		return
	}

	storage.Insert(storage.PayDatabase(), &request)

	order.State = model.OrderStatePlace
	order.SNSOrderId = request.OrderNumber

	//订单再处理
	err = storage.Save(storage.PayDatabase(), order)
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"request": request,
			"error":   err,
		}).Error("vivo callback failed.")

		c.String(http.StatusOK, "fail")

		return
	}

	c.String(http.StatusOK, "success")
	watch_waring.PaymentWatch(order.GlobalId, order.Vendor, order.Amount)
}
