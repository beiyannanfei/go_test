package idip

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	handlerFuncList[QUERY_BOOK_CHARACTER_INFO_REQ] =
		handlerFuncInfo{"query_book_character_info", query_book_character_info}
}

const QUERY_BOOK_CHARACTER_INFO_REQ = 0x101d
const QUERY_BOOK_CHARACTER_INFO_RSP = 0x101e

//
// 个人书籍角色信息查询请求
// Binding from JSON
type DoQueryBookCharacterInfoReq struct {
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
		// 用户全局id
		//
		GlobalId int64 `json:"GlobalId"`
	} `json:"body"`
}

//
// in: body
// swagger:parameters idip_query_book_character_info
type DoQueryBookCharacterInfoReqBodyParams struct {
	//
	// swagger:allOf
	// in: body
	DoQueryBookCharacterInfoReq DoQueryBookCharacterInfoReq
}

//
// 个人书籍角色信息查询应答
// swagger:response DoQueryBookCharacterInfoRsp
type DoQueryBookCharacterInfoRsp struct {
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
			//
			// 书籍角色信息
			//
			Character string `json:"Character"`
			//
			// 守护道具信息
			//
			GuardItem string `json:"GuardItem"`
		} `json:"body"`
	}
}

/*
    <entry cmd="10356014" req="IDIP_QUERY_BOOK_CHARACTER_INFO_REQ" rsp="IDIP_QUERY_BOOK_CHARACTER_INFO_RSP" req_value="area=AreaId#platid=PlatId#openid=OpenId" rsp_value="result=Result#error_info=RetMsg#character=Character_count|Character|#guard_item=GuardItem_count|GuardItem|" desc="个人书籍角色信息查询"/>

	<macro name="IDIP_QUERY_BOOK_CHARACTER_INFO_REQ" value="0x101d" desc="个人书籍角色信息查询请求"/>
    <macro name="IDIP_QUERY_BOOK_CHARACTER_INFO_RSP" value="0x101e" desc="个人书籍角色信息查询应答"/>

	<!--个人书籍角色信息查询请求-->
  <struct name="QueryBookCharacterInfoReq" id="IDIP_QUERY_BOOK_CHARACTER_INFO_REQ" desc="个人书籍角色信息查询请求">
    <entry name="AreaId" type="uint32" desc="服务器：微信（1），手Q（2），微信测试服（998），手Q测试服（999），预发布环境（99），ios审核服（98），体验服环境（97）     " test="1" isverify="true" isnull="true"/>
    <entry name="PlatId" type="uint8" desc="平台：IOS（0），安卓（1）" test="1" isverify="true" isnull="true"/>
    <entry name="OpenId" type="string" size="MAX_OPENID_LEN" desc="OpenId" test="732945400" isverify="false" isnull="true"/>
  </struct>
  <!--个人书籍角色信息查询应答-->
  <struct name="QueryBookCharacterInfoRsp" id="IDIP_QUERY_BOOK_CHARACTER_INFO_RSP" desc="个人书籍角色信息查询应答">
    <entry name="Result" type="int32" desc="结果：成功(0)，玩家不存在(1)，失败(其他)" test="0" isverify="false" isnull="true"/>
    <entry name="RetMsg" type="string" size="MAX_RETMSG_LEN" desc="返回消息" test="success" isverify="false" isnull="true"/>
	<entry name="Character_count" type="uint32" desc="书籍角色数量 " test="1" isverify="false" isnull="true"/>
    <entry name="GuardItem_count" type="uint32" desc="守护道具数量 " test="1" isverify="false" isnull="true"/>
    <entry name="Character" type="CharacterList" size="MAX_INFO_NUM" param="vector,struct,Character_count,1,|, " desc="书籍角色信息" isverify="false" isnull="true"/>
    <entry name="GuardItem" type="AttachList" size="MAX_INFO_NUM" param="vector,struct,GuardItem_count,1,|, " desc="守护道具信息" isverify="false" isnull="true"/>
  </struct>
*/

// swagger:route POST /idip/query_book_character_info idip idip_query_book_character_info
//
// query_book_character_info: 个人书籍角色信息查询
//
// IDIP命令编码 <br>
//   IDIP_QUERY_BOOK_CHARACTER_INFO_REQ = 0x101d <br>
//   IDIP_QUERY_BOOK_CHARACTER_INFO_RSP = 0x101e <br>
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
//       200: DoQueryBookCharacterInfoRsp
func query_book_character_info(c *gin.Context) {

	resp := DoQueryBookCharacterInfoRsp{}
	resp.Body.Head.PacketLen = 0
	resp.Body.Head.Cmdid = QUERY_BOOK_CHARACTER_INFO_RSP
	resp.Body.Head.ServiceName = "WorldOfLove"
	resp.Body.Head.SendTime = time.Now().Unix()
	resp.Body.Head.Version = 1
	resp.Body.Head.Authenticate = ""
	resp.Body.Head.Result = Ok
	resp.Body.Head.RetErrMsg = "ok"

	var request DoQueryBookCharacterInfoReq
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

		url := fmt.Sprintf("%s/%s", game, "query_book_character_info")

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
