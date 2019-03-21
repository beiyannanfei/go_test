package idip

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	handlerFuncList[DO_MODIFY_READ_PROGRESS_REQ] =
		handlerFuncInfo{"modify_read_progress", modify_read_progress}
}

const DO_MODIFY_READ_PROGRESS_REQ = 0x1019
const DO_MODIFY_READ_PROGRESS_RSP = 0x101a

//
// 修改当前阅读进度请求
// Binding from JSON
type DoModifyReadProgressReq struct {
	Head Header `json:"head"`
	Body struct {
		//
		// 服务器：微信（1），手Q（2），微信测试服（998），手Q测试服（999），预发布环境（99 ）
		//
		AreaId uint32 `json:"AreaId"`
		//
		// 平台：IOS（0），安卓（1）
		//
		PlatId uint8 `json:"PlatId"`
		//
		// OpenId
		//
		OpenId string `json:"OpenId"`
		//
		// 书籍ID
		//
		BookId uint32 `json:"BookId"`
		//
		// 跳转到的章节ID
		//
		ToChapterId uint32 `json:"ToChapterId"`
		//
		// 存档ID
		//
		SaveId uint32 `json:"SaveId"`
		//
		// 渠道号，由前端生成，不需填写
		//
		Source uint32 `json:"Source"`
		//
		// 流水号，由前端生成，不需要填写
		//
		Serial string `json:"Serial"`
		//
		// 用户全局id
		//
		GlobalId int64 `json:"GlobalId"`
	} `json:"body"`
}

//
// in: body
// swagger:parameters idip_modify_read_progress
type DoModifyReadProgressReqBodyParams struct {
	//
	// swagger:allOf
	// in: body
	DoModifyReadProgressReq DoModifyReadProgressReq
}

//
// 修改当前阅读进度应答
// swagger:response DoModifyReadProgressRsp
type DoModifyReadProgressRsp struct {
	// in: body
	Body struct {
		Head Header `json:"head"`
		Body struct {
			//
			// 结果：成功(0)，玩家不存在(1)，失败(其他)
			//
			Result int `json:"Result"`
			//
			// 返回消息
			//
			RetMsg string `json:"RetMsg"`
		} `json:"body"`
	}
}

/*
    <entry cmd="10356012" req="IDIP_DO_MODIFY_READ_PROGRESS_REQ" rsp="IDIP_DO_MODIFY_READ_PROGRESS_RSP" req_value="area=AreaId#platid=PlatId#openid=OpenId" rsp_value="result=Result#error_info=RetMsg" desc="修改当前阅读进度"/>

    <macro name="IDIP_DO_MODIFY_READ_PROGRESS_REQ" value="0x1019" desc="修改当前阅读进度请求"/>
    <macro name="IDIP_DO_MODIFY_READ_PROGRESS_RSP" value="0x101a" desc="修改当前阅读进度应答"/>

  <!--修改当前阅读进度请求-->
  <struct name="DoUpdateBookReadProgressReq" id="IDIP_DO_MODIFY_READ_PROGRESS_REQ" desc="修改当前阅读进度请求">
    <entry name="AreaId" type="uint32" desc="服务器：微信（1），手Q（2），微信测试服（998），手Q测试服（999），预发布环境（99），ios审核服（98），体验服环境（97）      " test="1" isverify="true" isnull="true"/>
    <entry name="PlatId" type="uint8" desc="平台：IOS（0），安卓（1）" test="1" isverify="true" isnull="true"/>
    <entry name="OpenId" type="string" size="MAX_OPENID_LEN" desc="OpenId" test="732945400" isverify="false" isnull="true"/>
    <entry name="BookId" type="uint32" desc="书籍ID" test="1" isverify="true" isnull="true"/>
    <entry name="ToChapterId" type="uint32" desc="跳转到的章节ID" test="1" isverify="true" isnull="true"/>
    <entry name="SaveId" type="uint32" desc="存档ID" test="1" isverify="true" isnull="true"/>
    <entry name="Source" type="uint32" desc="渠道号，由前端生成，不需填写" test="11" isverify="true" isnull="true"/>
    <entry name="Serial" type="string" size="MAX_SERIAL_LEN" desc="流水号，由前端生成，不需要填写" test="M-PAYO-20140414124009-58382166" isverify="false" isnull="true"/>
  </struct>
  <!--修改当前阅读进度应答-->
  <struct name="DoUpdateBookReadProgressRsp" id="IDIP_DO_MODIFY_READ_PROGRESS_RSP" desc="修改当前阅读进度应答">
    <entry name="Result" type="int32" desc="结果：成功(0)，玩家不存在(1)，失败(其他)" test="0" isverify="false" isnull="true"/>
    <entry name="RetMsg" type="string" size="MAX_RETMSG_LEN" desc="返回消息" test="success" isverify="false" isnull="true"/>
  </struct>
*/

// swagger:route POST /idip/modify_read_progress idip idip_modify_read_progress
//
// modify_read_progress: 修改当前阅读进度
//
// IDIP命令编码 <br>
//   IDIP_DO_MODIFY_READ_PROGRESS_REQ = 0x1019 <br>
//   IDIP_DO_MODIFY_READ_PROGRESS_RSP = 0x101a <br>
//
//    Parameters:
//		DoModifyReadProgressReq
//
//   Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Responses:
//       200: DoModifyReadProgressRsp
func modify_read_progress(c *gin.Context) {

	resp := DoModifyReadProgressRsp{}
	resp.Body.Head.PacketLen = 0
	resp.Body.Head.Cmdid = DO_MODIFY_READ_PROGRESS_RSP
	resp.Body.Head.ServiceName = "WorldOfLove"
	resp.Body.Head.SendTime = time.Now().Unix()
	resp.Body.Head.Version = 1
	resp.Body.Head.Authenticate = ""
	resp.Body.Head.Result = Ok
	resp.Body.Head.RetErrMsg = "ok"

	var request DoModifyReadProgressReq
	// validation
	if err := c.BindJSON(&request); err != nil {

		c.JSON(http.StatusOK, resp.Body)

		return
	}

	resp.Body.Head.Seqid = request.Head.Seqid

	AreaId := request.Body.AreaId
	PlatId := request.Body.PlatId
	OpenId := request.Body.OpenId

	globalId := getGlobalIdByOpenIdAreaId(OpenId, AreaId)
	if globalId == 0 {

		resp.Body.Head.Result = UserNotFoundErr
		resp.Body.Head.RetErrMsg = fmt.Sprintf("[%v,%v]Account not found", AreaId, PlatId)
		resp.Body.Body.Result = UserNotFoundErr
		resp.Body.Body.RetMsg = fmt.Sprintf("[%v,%v]Account not found", AreaId, PlatId)
		c.JSON(http.StatusOK, resp.Body)

		return
	}

	request.Body.GlobalId = globalId

	for count := 0; count <= RandReqCount; count++ {

		game := randGames(AreaId, PlatId)

		if game == "" {

			resp.Body.Head.Result = ZoneNotFoundErr
			resp.Body.Head.RetErrMsg = fmt.Sprintf("[%v,%v]Zone not found.", AreaId, PlatId)
			resp.Body.Body.Result = ZoneNotFoundErr
			resp.Body.Body.RetMsg = fmt.Sprintf("[%v,%v]Zone not found.", AreaId, PlatId)
			c.JSON(http.StatusOK, resp.Body)

			return
		}

		url := fmt.Sprintf("%s/%s", game, "modify_read_progress")

		data, err := PostRequest(url, request.Body)
		if err != nil {

			if count < RandReqCount {

				continue
			}

			resp.Body.Head.Result = NetWorkErr
			resp.Body.Head.RetErrMsg = fmt.Sprintf("[%v,%v]Game post failed: %v", AreaId, PlatId, err)
			resp.Body.Body.Result = NetWorkErr
			resp.Body.Body.RetMsg = fmt.Sprintf("[%v,%v]Game post failed: %v", AreaId, PlatId, err)
			c.JSON(http.StatusOK, resp.Body)

			return
		}

		if err := json.Unmarshal(data, &resp.Body.Body); err != nil || resp.Body.Body.Result != Ok {

			resp.Body.Head.Result = PostApiErr
			resp.Body.Head.RetErrMsg = fmt.Sprintf("[%v,%v]API failed: %v", AreaId, PlatId, err)
			resp.Body.Body.Result = PostApiErr
			resp.Body.Body.RetMsg = fmt.Sprintf("[%v,%v]API failed: %v", AreaId, PlatId, err)
			c.JSON(http.StatusOK, resp.Body)

			return
		}

		resp.Body.Body.RetMsg = "ok"
		c.JSON(http.StatusOK, resp.Body)

		return
	}
}
