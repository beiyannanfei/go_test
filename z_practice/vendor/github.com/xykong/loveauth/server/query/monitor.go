package query

import (
	"github.com/gin-gonic/gin"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"github.com/xykong/loveauth/utils"
	"github.com/xykong/loveauth/errors"
	"net/http"
	"github.com/jinzhu/gorm"
)

func init() {
	getHandlers["/monitor"] = monitor
}

type ResponseMonitor struct {
	Code       int64             `json:"code"`
	Message    string            `json:"message"`
	Account    model.Account     `json:"account"`
	Order      model.Order       `json:"order"`
	RedisValue map[string]string `json:"redis_value"`
}

func monitor(c *gin.Context) {
	dbAuth := storage.AuthDatabase()
	var account model.Account
	if err := dbAuth.Table("accounts").First(&account).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			utils.QuickReply(c, errors.Failed, "Query account failed: %v", err)
			return
		}
	}

	dbPay := storage.PayDatabase()
	var order model.Order
	if err := dbPay.Table("orders").First(&order).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			utils.QuickReply(c, errors.Failed, "Query order failed: %v", err)
			return
		}
	}

	redisValue, err := storage.MonitorRedis()
	if err != nil {
		utils.QuickReply(c, errors.Failed, "Query redis failed: %v", err)
		return
	}

	resp := ResponseMonitor{}
	resp.Code = int64(errors.Ok)
	resp.Message = "Query monitor successfully!"
	resp.Account = account
	resp.Order = order
	resp.RedisValue = redisValue
	c.JSON(http.StatusOK, resp)
	return
}
