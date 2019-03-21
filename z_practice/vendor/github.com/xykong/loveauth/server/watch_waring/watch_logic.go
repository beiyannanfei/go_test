package watch_waring

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/errors"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"github.com/xykong/loveauth/utils"
	"net/http"
)

// 道具累计获得或消耗

func init() {
	handlers["/item"] = itemWatch
	handlers["/login"] = loginWatch
	//getHandlers["/item"] = itemWatch
}

type RequestItem struct {
	GlobalId int64  `form:"GlobalId" json:"GlobalId" binding:"exists"`
	Channel  string `form:"Channel" json:"Channel" binding:"exists"`
	ItemId   string    `form:"ItemId" json:"ItemId" binding:"exists"`
	OpsType  int    `form:"OpsType" json:"OpsType" binding:"exists"`
	ItemNum  int    `form:"ItemNum" json:"ItemNum" binding:"exists"`
}

type WatchResponse struct {
	Body struct {
		Code    int64  `json:"code"`
		Message string `json:"message"`
	}

}

func itemWatch(c *gin.Context) {
	var request RequestItem
	if err := c.BindJSON(&request); err != nil {
		utils.QuickReply(c, errors.ServerFail, "连接服务器异常")
		return
	}
	go doItemWatch(request)
	resp := WatchResponse{}
	resp.Body.Code = int64(errors.Ok)
	resp.Body.Message = "item waring received"
	c.JSON(http.StatusOK, resp.Body)

}

func doItemWatch(request RequestItem) {
	if request.OpsType != 0 && request.OpsType != 1 {
		logrus.Info("invalid ops type, opsType = ", request.OpsType)
		return
	}
	var itemMap = GetItemWatchMap()
	//logrus.Info(itemMap)
	//itemId := strconv.Itoa(request.ItemId)
	if _, exist := itemMap[request.ItemId]; false == exist {
		//logrus.Info("no need waring itemId, itemId = ", request.ItemId, " opsType = ", request.OpsType)
		return
	}
	// 累加操作数量
	storage.UpdateItemState(request.GlobalId, request.Channel, request.ItemId, request.OpsType, request.ItemNum)
}



type RequestLogin struct {
	GlobalId int64 `form:"GlobalId" json:"GlobalId" binding:"exists"`
	DeviceId string `form:"DeviceId" json:"DeviceId" binding:"exists"`
	Channel string `form:"Channel" json:"Channel" binding:"exists"`
}

func loginWatch(c *gin.Context) {
	var request RequestLogin
	if err := c.BindJSON(&request); err != nil {
		utils.QuickReply(c, errors.ServerFail, "连接服务器异常")
		return
	}
	go doLoginWatch(request)
	resp := WatchResponse{}
	resp.Body.Code = int64(errors.Ok)
	resp.Body.Message = "login waring received"
	c.JSON(http.StatusOK, resp.Body)
}

func doLoginWatch(request RequestLogin) {
	// 累加操作数量
	storage.UpdateLoginState(request.GlobalId, request.Channel, request.DeviceId)
}

func PaymentWatch(globalId int64, vendor model.Vendor, amount int) {
	storage.UpdatePaymentState(globalId, vendor, amount)
}