package vivo

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/errors"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"github.com/xykong/loveauth/utils"
	"net/http"
)

func VivoCallBack(c *gin.Context) {

	var request model.VivoQueryOrderResponse

	// validation
	if err := c.BindJSON(&request); err != nil {

		utils.QuickReply(c, errors.ServerFail, "连接服务器异常")

		return
	}

	//storage.Insert(storage.PayDatabase(), &request)

	if request.Ret != "200" || request.TradeStatus != "0000" {

		logrus.WithFields(logrus.Fields{
			"Request": request,
		}).Error("vivo callback failed.")

		c.String(http.StatusOK, "fail")

		return
	}

	order := storage.QueryOrderPlacedWithSequence(request.CpOrderNumber)
	if order == nil || order.State != model.OrderStatePrepare {

		logrus.WithFields(logrus.Fields{
			"Request": request,
			"order":   order,
		}).Error("vivo callback failed.")

		c.String(http.StatusOK, "fail")

		return
	}

	order.State = model.OrderStatePlace

	//订单再处理
	err := storage.Save(storage.PayDatabase(), order)
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"Request": request,
			"Error":   err,
		}).Error("vivo callback failed.")

		c.String(http.StatusOK, "fail")

		return
	}

	c.String(http.StatusOK, "success")
}
