package idip

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func init() {

	idipV2["repire_order"] = repire_order
}

func repire_order(c *gin.Context) {

	resp := &CommonResp{Msg: "success"}

	var req model.RepireOrderRequest
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
	req.Operator = account

	order := storage.QueryOrderPlacedWithSequence(req.Sequence)
	if order == nil {

		resp.Msg = "sequence order nil"
		c.JSON(http.StatusOK, resp)

		return
	}

	if order.GlobalId != req.GlobalId {

		resp.Msg = "sequence globalId not match"
		c.JSON(http.StatusOK, resp)

		return
	}

	if order.State != model.OrderStatePlace {

		resp.Msg = fmt.Sprintf("order state is %d", order.State)
		c.JSON(http.StatusOK, resp)

		return
	}

	var sendMailReq DoSendItemMailReq
	sendMailReq.Body.Attatch, sendMailReq.Body.Attatch_count = awardListToAttatchList(req.AwardList)
	sendMailReq.Body.MailTitle = req.MailTitle
	sendMailReq.Body.MailContent = req.MailContent
	sendMailReq.Body.GlobalId = req.GlobalId
	sendMailReq.Body.Gm = account

	sendMailReqBody, _ := json.Marshal(sendMailReq.Body)
	sendMailResp, err := http.Post(url+"jdifoa/send_item_mail", "application/json", bytes.NewBuffer(sendMailReqBody))
	sendMailRespBody, err := ioutil.ReadAll(sendMailResp.Body)
	var sendMailGsResp GsResponse
	err = json.Unmarshal(sendMailRespBody, &sendMailGsResp)
	if err != nil {

		resp.Msg = "json err " + err.Error()
		c.JSON(http.StatusOK, resp)

		return
	}

	if sendMailGsResp.RetMsg != "success" || sendMailGsResp.Result != 0 {

		resp.Msg = "gs fail " + string(sendMailRespBody)
		c.JSON(http.StatusOK, resp)

		return
	}

	err = storage.UpdateOrderState(req.Sequence, model.OrderStateRepaired)
	if err != nil {

		resp.Msg = "update order state " + err.Error()
		c.JSON(http.StatusOK, resp)

		return
	}

	req.RepireTime = time.Now()
	err = storage.Insert(storage.PayDatabase(), &req)
	if err != nil {

		resp.Msg = "insert repire note " + err.Error()
		c.JSON(http.StatusOK, resp)

		return
	}

	c.JSON(http.StatusOK, resp)
	go sendWaringMail(req.GlobalId, account, req.Channel, req.MailTitle, req.MailContent, req.AwardList)
}

func awardListToAttatchList(awardList string) ([]Attatch, uint32) {

	var result []Attatch
	Items := strings.Split(awardList, ",")
	for _, item := range Items {

		info := strings.Split(item, "_")
		if len(info) == 2 {

			id, _ := strconv.Atoi(info[0])
			num, _ := strconv.Atoi(info[1])
			result = append(result, Attatch{id, num})
		}
	}

	return result, uint32(len(result))
}

type GsResponse struct {
	Result int    `json:"Result"`
	RetMsg string `json:"RetMsg"`
}
