package idip

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/xykong/loveauth/storage"
	"io/ioutil"
	"net/http"
)

func init() {

	idipV2["delete_rank_message"] = delete_rank_message
}

type DeleteRankMessageRequest struct {
	Token     string `form:"Token" json:"Token"`
	Character int    `form:"Character" json:"Character"`
	Uid       uint64 `form:"Uid" json:"Uid"`
}

type DeleteRankMessageGsReq struct {
	Character int    `json:"character"`
	Uid       uint64 `json:"uid"`
}

type DeleteRankMessageGsResp struct {
	Result int    `json:"Result"`
	RetMsg string `json:"RetMsg"`
}

func delete_rank_message(c *gin.Context) {

	resp := &CommonResp{Msg: "success"}

	var req DeleteRankMessageRequest
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

	gsReq := &DeleteRankMessageGsReq{}
	gsReq.Uid = req.Uid
	gsReq.Character = req.Character

	reqBody, _ := json.Marshal(gsReq)
	deleteRankMessageResp, err := http.Post(url+"jdifoa/del_rank_message", "application/json", bytes.NewBuffer(reqBody))
	deleteRankMessageRespBody, err := ioutil.ReadAll(deleteRankMessageResp.Body)

	respData := &DeleteRankMessageGsResp{}
	err = json.Unmarshal(deleteRankMessageRespBody, respData)
	if err != nil {

		resp.Msg = "json err " + err.Error()
		c.JSON(http.StatusOK, resp)

		return
	}

	if respData.RetMsg != "success" || respData.Result != 0 {

		resp.Msg = "gs return err " + string(deleteRankMessageRespBody)
		c.JSON(http.StatusOK, resp)

		return
	}

	c.JSON(http.StatusOK, resp)
}
