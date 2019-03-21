package idip

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/xykong/loveauth/storage"
	"io/ioutil"
	"net/http"
	"time"
)

func init() {

	idipV2["send_mail"] = send_mail
}

type DoSendMailReq struct {
	GlobalId    int64  `form:"GlobalId" json:"GlobalId"`
	MailTitle   string `form:"MailTitle" json:"MailTitle"`
	MailContent string `form:"MailContent" json:"MailContent"`
	AwardList   string `form:"AwardList" json:"AwardList"`
	Token       string `form:"Token" json:"Token"`
	EffectTime  string `form:"EffectTime" json:"EffectTime"`
	CloseTime   string `form:"CloseTime" json:"CloseTime"`
	Channel     string `form:"Channel" json:"Channel"`
}

type CommonResp struct {
	Msg string `json:"msg"`
}

func send_mail(c *gin.Context) {

	resp := &CommonResp{Msg: "success"}

	var req DoSendMailReq
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

	var sendMailReq DoSendItemMailReq
	sendMailReq.Body.Attatch, sendMailReq.Body.Attatch_count = awardListToAttatchList(req.AwardList)
	sendMailReq.Body.MailTitle = req.MailTitle
	sendMailReq.Body.MailContent = req.MailContent
	sendMailReq.Body.GlobalId = req.GlobalId
	sendMailReq.Body.Gm = account

	if req.EffectTime != "" && req.EffectTime != "20060102150405" {
		startTime, _ := time.Parse("20060102150405", req.EffectTime)
		sendMailReq.Body.EffectTime = uint32(startTime.Unix())
	}
	if req.CloseTime != "" && req.CloseTime != "20060102150405" {
		endTime, _ := time.Parse("20060102150405", req.CloseTime)
		sendMailReq.Body.EffectTime = uint32(endTime.Unix())
	}
	sendMailReq.Body.Channel = req.Channel

	sendMailReqBody, _ := json.Marshal(sendMailReq.Body)
	var sendMailResp *http.Response
	if req.GlobalId != 0 {

		sendMailResp, err = http.Post(url+"jdifoa/send_item_mail", "application/json", bytes.NewBuffer(sendMailReqBody))
	} else if req.Channel != "" {

		sendMailResp, err = http.Post(url+"jdifoa/send_all_mail", "application/json", bytes.NewBuffer(sendMailReqBody))
	} else {

		resp.Msg = "mail type err"
		c.JSON(http.StatusOK, resp)

		return
	}

	if err != nil {

		resp.Msg = "gs http req err " + err.Error()
		c.JSON(http.StatusOK, resp)

		return
	}

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

	c.JSON(http.StatusOK, resp)

	go sendWaringMail(req.GlobalId, account, req.Channel, req.MailTitle, req.MailContent, req.AwardList)
}
