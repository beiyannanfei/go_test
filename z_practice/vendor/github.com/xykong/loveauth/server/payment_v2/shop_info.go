package payment_v2

import (
	"github.com/gin-gonic/gin"
	"github.com/xykong/loveauth/errors"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/utils"
	"net/http"
	"time"
)

func init() {
	handlers["/shop_info"] = shopInfo
}

//
// Request shop item
//
type DoShopInfoReq struct {
	// The sequence
	//
	//
	Sequence string `form:"sequence" json:"sequence"`
	//The shopId
	//
	//
	ShopId int `form:"shopId" json:"shopId"`
	// The globalId
	//
	//
	GlobalId int64 `form:"globalId" json:"globalId"`
}

//
// in: body
// swagger:parameters payment_shop_info
type DoShopInfoReqBodyParams struct {
	//
	// swagger:allOf
	// in: body
	DoShopInfoReq DoShopInfoReq
}

//
// 应答: 协议返回包
// swagger:response DoShopInfoRsp
// noinspection ALL
type DoShopInfoRsp struct {
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
		// The response state
		// 3:订单完成 其他:订单未完成
		// required: true
		State int `json:"state"`
		// The response shopId
		//
		// required: true
		ShopId int `json:"shopId"`
		// The response limitcount
		//
		//
		LimitCount int `json:"limitCount"`
		// The response availableCount
		//
		//
		AvailableCount int `json:"availableCount"`
	}
}

//
// swagger:route POST /payment/shop_info payment payment_shop_info
//
// Return shop info for the given user:
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
//       200: DoShopInfoRsp
func shopInfo(c *gin.Context) {

	var request DoShopInfoReq

	// validation
	if err := c.BindJSON(&request); err != nil {

		utils.QuickReply(c, errors.Failed, "shopInfo BindJSON failed: %v", err)
		return
	}

	resp := DoShopInfoRsp{}
	resp.Body.Code = int64(errors.Ok)
	resp.Body.Message = "shopInfo successfully!"

	if request.ShopId != 0 {

		shopItem := getShopItem(request.ShopId, request.GlobalId, 0, time.Now())
		resp.Body.ShopId = request.ShopId
		resp.Body.LimitCount = shopItem.LimitCount
		resp.Body.AvailableCount = shopItem.AvailableCount
	}

	if request.Sequence != "" {

		order := storage.QueryOrderPlacedWithSequence(request.Sequence)
		if order != nil {

			resp.Body.State = int(order.State)
			resp.Body.ShopId = order.ShopId
		}
	}

	c.JSON(http.StatusOK, resp.Body)
}
