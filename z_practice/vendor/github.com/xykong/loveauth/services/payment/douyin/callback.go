package douyin

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
	"net/url"
)

func callback(c *gin.Context) {

	reqBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {

		c.String(http.StatusOK, "douyin io read err "+err.Error())

		return
	}

	value, err := url.ParseQuery(string(reqBody))
	if err != nil {

		c.String(http.StatusOK, "douyin parseQuery err "+err.Error())

		return
	}

	payKey := ""
	if !CheckSign(value, payKey) {

		c.String(http.StatusOK, "douyin checkoutsign err")

		return
	}

}

type DouyinCallBackRequest struct {
	gorm.Model  `json:"-"`
	NotifyId    string `json:"notify_id" form:"notify_id"`
	NotifyType  string `json:"notify_type" form:"notify_type"`
	NotifyTime  string `json:"notify_time" form:"notify_time"`
	TradeStatus string `json:"trade_status" form:"trade_status"`
	Way         int    `json:"way" form:"way"`
	ClientId    string `json:"client_id" form:"client_id"`
	OutTradeNo  string `json:"out_trade_no" form:"out_trade_no"`
	TradeNo     string `json:"trade_no" form:"trade_no" gorm:"type:varchar(191); index"`
	PayTime     string `json:"pay_time" form:"pay_time"`
	TotalFee    int    `json:"total_fee" form:"total_fee"`
	BuyerId     string `json:"buyer_id" form:"buyer_id"`
	TtSign      string `json:"tt_sign" form:"tt_sign" gorm:"size:512"`
	TtSignType  string `json:"tt_sign_type" form:"tt_sign_type"`
}
