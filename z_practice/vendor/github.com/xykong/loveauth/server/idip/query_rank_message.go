package idip

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/xykong/loveauth/storage"
	"io/ioutil"
	"net/http"
)

func init() {

	idipV2["query_rank_message"] = query_rank_message
}

type QueryRankMessageRequest struct {
	Token string `form:"Token" json:"Token"`
}

type QueryRankMessageResponse struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type QueryRankMessageGsResp struct {
	Result int         `json:"Result"`
	RetMsg string      `json:"RetMsg"`
	Data   interface{} `json:"Data"`
}

func query_rank_message(c *gin.Context) {

	resp := &QueryRankMessageResponse{Msg: "success"}

	var req RepireNotesRequest
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

	queryRankMessageResp, err := http.Post(url+"jdifoa/query_rank_message", "application/json", nil)
	queryRankMessageRespBody, err := ioutil.ReadAll(queryRankMessageResp.Body)

	respData := &QueryRankMessageGsResp{}
	err = json.Unmarshal(queryRankMessageRespBody, respData)
	if err != nil {

		resp.Msg = "json err " + err.Error()
		c.JSON(http.StatusOK, resp)

		return
	}

	if respData.RetMsg != "success" || respData.Result != 0 {

		resp.Msg = "gs return err " + string(queryRankMessageRespBody)
		c.JSON(http.StatusOK, resp)

		return
	}

	resp.Data = respData.Data

	c.JSON(http.StatusOK, resp)
}
