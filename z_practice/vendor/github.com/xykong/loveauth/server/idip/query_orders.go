package idip

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/xykong/loveauth/server/payment_v2"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"io/ioutil"
	"net/http"
	"time"
)

func init() {

	idipV2["query_orders"] = query_orders
}

type QueryOrdersRequest struct {
	GlobalId  int64  `form:"GlobalId" json:"GlobalId"`
	StartTime string `form:"StartTime" json:"StartTime"`
	EndTime   string `form:"EndTime" json:"EndTime"`
	Token     string `form:"Token" json:"Token"`
}

type QueryOrdersResponse struct {
	Msg  string       `json:"msg"`
	Data []*OrderInfo `json:"data"`
}

type OrderInfo struct {
	GlobalId  int64        `json:"global_id"`
	ShopId    int          `json:"shop_id"`
	ItemId    int          `json:"item_id"`
	ItemName  string       `json:"item_name"`
	ItemNum   int          `json:"item_num"`
	Rmb       int          `json:"rmb"`
	PayMethod int          `json:"pay_method"`
	IsFirst   bool         `json:"is_first"`
	IsPay     bool         `json:"is_pay"`
	IsSend    bool         `json:"is_send"`
	IsRepire  bool         `json:"is_repire"`
	Channel   model.Vendor `json:"channel"`
	Time      string       `json:"time"`
	Sequence  string       `json:"sequence"`
}

func query_orders(c *gin.Context) {

	resp := &QueryOrdersResponse{Msg: "success"}

	var req QueryOrdersRequest
	err := c.Bind(&req)
	if err != nil {

		resp.Msg = "fail " + err.Error()
		c.JSON(http.StatusOK, resp)

		return
	}

	account, url := storage.GetSession(req.Token)
	if account == "" {

		resp.Msg = "session err"
		c.JSON(http.StatusOK, resp)

		return
	}

	startTime, _ := time.Parse("20060102150405", req.StartTime)
	endTime, _ := time.Parse("20060102150405", req.EndTime)

	orders := storage.QueryOrderByGlobalId(req.GlobalId, startTime, endTime)

	var data []*OrderInfo
	var itemIds []int
	for _, order := range orders {

		var orderInfo OrderInfo
		orderInfo.GlobalId = order.GlobalId
		orderInfo.Sequence = order.Sequence
		orderInfo.PayMethod = order.PayMethod
		orderInfo.IsSend = order.State == model.OrderStateComplete
		orderInfo.IsRepire = order.State == model.OrderStateRepaired
		orderInfo.IsPay = order.State > model.OrderStatePrepare && order.State != model.OrderStateFailed

		orderInfo.IsFirst = storage.QuerySameShopIdCountPre(order.GlobalId, order.ShopId, order.Timestamp) == 0
		orderInfo.Rmb = order.Amount
		orderInfo.Channel = order.Vendor
		orderInfo.Time = order.Timestamp.Format("2006-01-02 15:04:05")
		orderInfo.ShopId = order.ShopId
		shopItem := payment_v2.GetShopItem(order.ShopId, order.GlobalId, order.PayMethod, order.Timestamp)
		orderInfo.ItemId = shopItem.ItemId
		if order.Type == model.OrderTypeDiamond {

			orderInfo.ItemNum = shopItem.ItemCount + shopItem.GiftPriceValue

			if order.ShopId == 1000 {

				orderInfo.ItemNum = shopItem.ItemCount
			}
		} else {

			orderInfo.ItemNum = shopItem.ItemCount
		}

		data = append(data, &orderInfo)
		itemIds = append(itemIds, orderInfo.ItemId)
	}

	if len(itemIds) != 0 {

		var queryItemsReq QueryItemsReq
		queryItemsReq.Data = itemIds
		queryItemsInfoBody, _ := json.Marshal(queryItemsReq)
		queryItemsInfoResp, err := http.Post(url+"jdifoa/query_items_info", "application/json", bytes.NewBuffer(queryItemsInfoBody))
		queryItemsInfoRespBody, err := ioutil.ReadAll(queryItemsInfoResp.Body)
		var queryItemsInfoGsResp QueryItemsInfoGsResp

		err = json.Unmarshal(queryItemsInfoRespBody, &queryItemsInfoGsResp)
		if err != nil {

			resp.Msg = "json err " + err.Error()
			c.JSON(http.StatusOK, resp)

			return
		}

		if queryItemsInfoGsResp.RetMsg != "success" || queryItemsInfoGsResp.Result != 0 {

			resp.Msg = "gs fail " + string(queryItemsInfoRespBody)
			c.JSON(http.StatusOK, resp)

			return
		}

		itemInfo := make(map[int]string, len(queryItemsInfoGsResp.Data))
		for _, item := range queryItemsInfoGsResp.Data {

			itemInfo[item.ItemId] = item.ItemName
		}

		for _, order := range data {

			if v, ok := itemInfo[order.ItemId]; ok {

				order.ItemName = v
			}
		}
	}

	resp.Data = data

	c.JSON(http.StatusOK, resp)
}

type QueryItemsReq struct {
	Data []int `json:"data"`
}

type QueryItemsInfoGsResp struct {
	Result int          `json:"Result"`
	RetMsg string       `json:"RetMsg"`
	Data   []*ItemsInfo `json:"Data"`
}

type ItemsInfo struct {
	ItemId   int    `json:"item_id"`
	ItemName string `json:"item_name"`
}
