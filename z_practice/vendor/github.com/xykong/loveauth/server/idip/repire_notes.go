package idip

import (
	"github.com/gin-gonic/gin"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"net/http"
	"time"
)

func init() {

	idipV2["repire_notes"] = repire_notes
}

type RepireNotesRequest struct {
	GlobalId  int64  `form:"GlobalId" json:"GlobalId"`
	Operator  string `form:"Operator" json:"Operator"`
	StartTime string `form:"StartTime" json:"StartTime"`
	EndTime   string `form:"EndTime" json:"EndTime"`
	Token     string `form:"Token" json:"Token"`
}

type RepireNotesResponse struct {
	Msg  string                      `json:"msg"`
	Data []*model.RepireOrderRequest `json:"data"`
}

func repire_notes(c *gin.Context) {

	resp := &RepireNotesResponse{Msg: "success"}

	var req RepireNotesRequest
	err := c.Bind(&req)
	if err != nil {

		resp.Msg = "fail " + err.Error()
		c.JSON(http.StatusOK, resp)

		return
	}

	account, _ := storage.GetSession(req.Token)
	if account == "" {

		resp.Msg = "session err"
		c.JSON(http.StatusOK, resp)

		return
	}

	startTime, _ := time.Parse("20060102150405", req.StartTime)
	endTime, _ := time.Parse("20060102150405", req.EndTime)

	var queryResult []*model.RepireOrderRequest
	if req.GlobalId > 0 {

		queryResult = storage.QueryRepireOrderRequestWithGlobalId(req.GlobalId, startTime, endTime)
	} else if req.Operator != "" {

		queryResult = storage.QueryRepireOrderRequestWithOperator(req.Operator, startTime, endTime)
	}

	resp.Data = queryResult

	c.JSON(http.StatusOK, resp)
}
