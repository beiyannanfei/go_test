package idip

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	handlerFuncList[QUERY_READ_BOOK_INFO_REQ] =
		handlerFuncInfo{"query_read_book_info", query_read_book_info}
}

const QUERY_READ_BOOK_INFO_REQ = 0x101b
const QUERY_READ_BOOK_INFO_RSP = 0x101c

//
// 个人阅读信息查询请求
// Binding from JSON
type DoQueryReadBookInfoReq struct {
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
// swagger:parameters idip_query_read_book_info
type DoQueryReadBookInfoReqBodyParams struct {
	//
	// swagger:allOf
	// in: body
	DoQueryReadBookInfoReq DoQueryReadBookInfoReq
}

//
// 个人阅读信息查询应答
// swagger:response DoQueryReadBookInfoRsp
type DoQueryReadBookInfoRsp struct {
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
			// 书籍阅读信息
			//
			BookRead string `json:"BookRead"`
			//
			// 相册获取信息
			//
			Photo string `json:"Photo"`
			//
			// 声音获取信息
			//
			Voice string `json:"Voice"`
			//
			// 服装获取信息
			//
			Clothes string `json:"Clothes"`
			//
			// 物品获取信息
			//
			Goods string `json:"Goods"`
			//
			// 故事获取信息
			//
			Story string `json:"Story"`
		} `json:"body"`
	}
}

/*
    <entry cmd="10356013" req="IDIP_QUERY_READ_BOOK_INFO_REQ" rsp="IDIP_QUERY_READ_BOOK_INFO_RSP" req_value="area=AreaId#platid=PlatId#openid=OpenId" rsp_value="result=Result#error_info=RetMsg#book_read=Book_count|BookRead|#photo=Photo_count|Photo|#voice=Voice_count|Voice|#clothes=Clothes_count|Clothes|#goods=Goods_count|Goods|#story=Story_count|Story|" desc="个人阅读信息查询"/>

    <macro name="IDIP_QUERY_READ_BOOK_INFO_REQ" value="0x101b" desc="个人阅读信息查询请求"/>
    <macro name="IDIP_QUERY_READ_BOOK_INFO_RSP" value="0x101c" desc="个人阅读信息查询应答"/>

<!--书籍阅读信息表-->
  <struct name="BookReadList" desc="书籍阅读信息表">
    <entry name="BookId" type="uint32" desc="书籍ID" test="1" isverify="true" isnull="true"/>
    <entry name="BookName" type="string" desc="书籍名字" test="1" isverify="true" isnull="true"/>
    <entry name="MaxChapter" type="uint32" desc="最大章节" test="1" isverify="true" isnull="true"/>
    <entry name="CurrChapter" type="uint32" desc="当前章节" test="1" isverify="true" isnull="true"/>
    <entry name="SaveCount" type="uint32" desc="开启存档数量" test="1" isverify="true" isnull="true"/>
  </struct>
  <!--书籍角色收集获取表-->
  <struct name="CollectionList" desc="书籍角色收集获取表">
    <entry name="BookId" type="uint32" desc="书籍ID" test="1" isverify="true" isnull="true"/>
    <entry name="BookName" type="string" desc="书籍名字" test="1" isverify="true" isnull="true"/>
    <entry name="CharacterId" type="uint32" desc="角色ID" test="1" isverify="true" isnull="true"/>
    <entry name="CharacterName" type="string" desc="角色名字" test="1" isverify="true" isnull="true"/>
    <entry name="CollectionId" type="uint32" desc="收集品ID" test="1" isverify="true" isnull="true"/>
    <entry name="CollectionName" type="string" desc="收集品名字" test="1" isverify="true" isnull="true"/>
    <entry name="CollectionState" type="uint32" desc="持有状态0：未获取1：获取" test="1" isverify="true" isnull="true"/>
  </struct>

  <!--个人阅读信息查询请求-->
  <struct name="QueryReadBookInfoReq" id="IDIP_QUERY_READ_BOOK_INFO_REQ" desc="个人阅读信息查询请求">
    <entry name="AreaId" type="uint32" desc="服务器：微信（1），手Q（2），微信测试服（998），手Q测试服（999），预发布环境（99），ios审核服（98），体验服环境（97）     " test="1" isverify="true" isnull="true"/>
    <entry name="PlatId" type="uint8" desc="平台：IOS（0），安卓（1）" test="1" isverify="true" isnull="true"/>
    <entry name="OpenId" type="string" size="MAX_OPENID_LEN" desc="OpenId" test="732945400" isverify="false" isnull="true"/>
  </struct>
  <!--个人阅读信息查询应答-->
  <struct name="QueryReadBookInfoRsp" id="IDIP_QUERY_READ_BOOK_INFO_RSP" desc="个人阅读信息查询应答">
    <entry name="Result" type="int32" desc="结果：成功(0)，玩家不存在(1)，失败(其他)" test="0" isverify="false" isnull="true"/>
    <entry name="RetMsg" type="string" size="MAX_RETMSG_LEN" desc="返回消息" test="success" isverify="false" isnull="true"/>
	<entry name="Book_count" type="uint32" desc="书籍数量 " test="1" isverify="false" isnull="true"/>
	<entry name="BookRead" type="BookReadList" size="MAX_INFO_NUM" param="vector,struct,Book_count,1,|, " desc="书籍阅读信息" isverify="false" isnull="true"/>
	<entry name="Photo_count" type="uint32" desc="相册数量 " test="1" isverify="false" isnull="true"/>
    <entry name="Voice_count" type="uint32" desc="声音数量 " test="1" isverify="false" isnull="true"/>
    <entry name="Clothes_count" type="uint32" desc="服装数量 " test="1" isverify="false" isnull="true"/>
    <entry name="Goods_count" type="uint32" desc="物品数量 " test="1" isverify="false" isnull="true"/>
    <entry name="Story_count" type="uint32" desc="故事数量 " test="1" isverify="false" isnull="true"/>
	<entry name="Photo" type="CollectionList" size="MAX_INFO_NUM" param="vector,struct,Photo_count,1,|, " desc="相册获取信息" isverify="false" isnull="true"/>
    <entry name="Voice" type="CollectionList" size="MAX_INFO_NUM" param="vector,struct,Voice_count,1,|, " desc="声音获取信息" isverify="false" isnull="true"/>
    <entry name="Clothes" type="CollectionList" size="MAX_INFO_NUM" param="vector,struct,Clothes_count,1,|, " desc="服装获取信息" isverify="false" isnull="true"/>
    <entry name="Goods" type="CollectionList" size="MAX_INFO_NUM" param="vector,struct,Goods_count,1,|, " desc="物品获取信息" isverify="false" isnull="true"/>
    <entry name="Story" type="CollectionList" size="MAX_INFO_NUM" param="vector,struct,Story_count,1,|, " desc="故事获取信息" isverify="false" isnull="true"/>
  </struct>
*/

// swagger:route POST /idip/query_read_book_info idip idip_query_read_book_info
//
// query_read_book_info: 个人阅读信息查询
//
// IDIP命令编码 <br>
//   IDIP_QUERY_READ_BOOK_INFO_REQ = 0x101b <br>
//   IDIP_QUERY_READ_BOOK_INFO_RSP = 0x101c <br>
//
//    Parameters:
//		DoQueryReadBookInfoReq
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
//       200: DoQueryReadBookInfoRsp
func query_read_book_info(c *gin.Context) {

	resp := DoQueryReadBookInfoRsp{}
	resp.Body.Head.PacketLen = 0
	resp.Body.Head.Cmdid = QUERY_READ_BOOK_INFO_RSP
	resp.Body.Head.ServiceName = "WorldOfLove"
	resp.Body.Head.SendTime = time.Now().Unix()
	resp.Body.Head.Version = 1
	resp.Body.Head.Authenticate = ""
	resp.Body.Head.Result = Ok
	resp.Body.Head.RetErrMsg = "ok"

	var request DoQueryReadBookInfoReq
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

		url := fmt.Sprintf("%s/%s", game, "query_read_book_info")

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
