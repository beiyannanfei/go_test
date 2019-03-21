package idip

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {

	handlerFuncList[DO_DEL_ITEM_REQ] =
		handlerFuncInfo{"del_item", del_item}
}

const DO_DEL_ITEM_REQ = 0x100b
const DO_DEL_ITEM_RSP = 0x100c

//
// 删除道具请求
// Binding from JSON
type DoDelItemReq struct {
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
		// 道具ID
		//
		ItemId uint64 `json:"ItemId"`
		//
		// 道具等级：（备选）
		//
		ItemLv uint32 `json:"ItemLv"`
		//
		// 删除数量
		//
		ItemNum uint32 `json:"ItemNum"`
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
// swagger:parameters idip_del_item
type DoDelItemReqBodyParams struct {
	//
	// swagger:allOf
	// in: body
	DoDelItemReq DoDelItemReq
}

//
// 删除道具应答
// swagger:response DoDelItemRsp
type DoDelItemRsp struct {
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
			// 道具名称
			//
			ItemName string `json:"ItemName"`
			//
			// 删除前数量
			//
			BeforeNum int64 `json:"BeforeNum"`
			//
			// 当前道具数量
			//
			AfterNum int64 `json:"AfterNum"`
		} `json:"body"`
	}
}

/*
    <entry cmd="10356005" req="IDIP_DO_DEL_ITEM_REQ" rsp="IDIP_DO_DEL_ITEM_RSP" req_value="area=AreaId#platid=PlatId#openid=OpenId#item_id=ItemId#item_lv=ItemLv#item_num=ItemNum#source=Source#serial=Serial" rsp_value="result=Result#error_info=RetMsg#item_name=ItemName#before_num=BeforeNum#after_num=AfterNum" desc="删除道具"/>

    <macro name="IDIP_DO_DEL_ITEM_REQ" value="0x100b" desc="删除道具请求"/>
    <macro name="IDIP_DO_DEL_ITEM_RSP" value="0x100c" desc="删除道具应答"/>

  <!--删除道具请求-->
  <struct name="DoDelItemReq" id="IDIP_DO_DEL_ITEM_REQ" desc="删除道具请求">
    <entry name="AreaId" type="uint32" desc="服务器：微信（1），手Q（2），微信测试服（998），手Q测试服（999），预发布环境（99 ）      " test="1" isverify="true" isnull="true"/>
    <entry name="PlatId" type="uint8" desc="平台：IOS（0），安卓（1）" test="1" isverify="true" isnull="true"/>
    <entry name="OpenId" type="string" size="MAX_OPENID_LEN" desc="OpenId" test="732945400" isverify="false" isnull="true"/>
    <entry name="ItemId" type="uint64" desc="道具ID" test="1" isverify="true" isnull="true"/>
    <entry name="ItemLv" type="uint32" desc="道具等级：（备选）" test="1" isverify="true" isnull="true"/>
    <entry name="ItemNum" type="uint32" desc="删除数量" test="1" isverify="true" isnull="true"/>
    <entry name="Source" type="uint32" desc="渠道号，由前端生成，不需填写" test="11" isverify="true" isnull="true"/>
    <entry name="Serial" type="string" size="MAX_SERIAL_LEN" desc="流水号，由前端生成，不需要填写" test="M-PAYO-20140414124009-58382166" isverify="false" isnull="true"/>
  </struct>
  <!--删除道具应答-->
  <struct name="DoDelItemRsp" id="IDIP_DO_DEL_ITEM_RSP" desc="删除道具应答">
    <entry name="Result" type="int32" desc="结果：成功(0)，玩家不存在(1)，失败(其他)" test="0" isverify="false" isnull="true"/>
    <entry name="RetMsg" type="string" size="MAX_RETMSG_LEN" desc="返回消息" test="success" isverify="false" isnull="true"/>
    <entry name="ItemName" type="string" size="MAX_ITEMNAME_LEN" desc="道具名称" test="test" isverify="false" isnull="true"/>
    <entry name="BeforeNum" type="uint32" desc="删除前数量" test="1" isverify="false" isnull="true"/>
    <entry name="AfterNum" type="uint32" desc="当前道具数量" test="1" isverify="false" isnull="true"/>
  </struct>
*/
//
// swagger:route POST /idip/del_item idip idip_del_item
//
// del_item: 删除道具
//
// IDIP命令编码 <br>
//   IDIP_DO_DEL_ITEM_REQ = 0x100b <br>
//   IDIP_DO_DEL_ITEM_RSP = 0x100c <br>
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
//       200: DoDelItemRsp
func del_item(c *gin.Context) {

	resp := DoDelItemRsp{}
	resp.Body.Head.PacketLen = 0
	resp.Body.Head.Cmdid = DO_DEL_ITEM_RSP
	resp.Body.Head.ServiceName = "WorldOfLove"
	resp.Body.Head.SendTime = time.Now().Unix()
	resp.Body.Head.Version = 1
	resp.Body.Head.Authenticate = ""
	resp.Body.Head.Result = Ok
	resp.Body.Head.RetErrMsg = "ok"

	var request DoDelItemReq
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

		url := fmt.Sprintf("%s/%s", game, "del_item")

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

		SendIDIPTLog(int(AreaId), OpenId, int(request.Body.ItemId), int(request.Body.ItemNum), request.Body.Serial, int(request.Body.Source),
			DO_DEL_ITEM_REQ, request.Head.Seqid, int(PlatId))

		return
	}
}
