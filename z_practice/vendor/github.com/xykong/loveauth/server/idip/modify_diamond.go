package idip

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {

	handlerFuncList[DO_MODIFY_DIAMOND_REQ] =
		handlerFuncInfo{"modify_diamond", modify_diamond}
}

// 修改钻石（代币）请求
const DO_MODIFY_DIAMOND_REQ = 0x1001
const DO_MODIFY_DIAMOND_RSP = 0x1002

// Binding from JSON
type DoModifyDiamondReq struct {
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
		// 修改数量：-减，+加
		//
		ModifyNum int64 `json:"ModifyNum"`
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
// swagger:parameters idip_modify_diamond
type DoModifyDiamondReqBodyParams struct {
	//
	// swagger:allOf
	// in: body
	DoModifyDiamondReq DoModifyDiamondReq
}

//
// 修改钻石（代币）应答
// swagger:response DoModifyDiamondRsp
type DoModifyDiamondRsp struct {
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
			// 修改前钻石量
			//
			BeforeNum int64 `json:"BeforeNum"`
			//
			// 修改后钻石数量
			//
			AfterNum int64 `json:"AfterNum"`
		} `json:"body"`
	}
}

/*
	<entry
		cmd="10356000"
		req="IDIP_DO_MODIFY_DIAMOND_REQ"
		rsp="IDIP_DO_MODIFY_DIAMOND_RSP"
		req_value="area=AreaId#platid=PlatId#openid=OpenId#modify_num=ModifyNum#source=Source#serial=Serial"
		rsp_value="result=Result#error_info=RetMsg#before_num=BeforeNum#after_num=AfterNum"
		desc="修改钻石（代币）"
	/>

	<macro name="IDIP_DO_MODIFY_DIAMOND_REQ" value="0x1001" desc="修改钻石（代币）请求"/>
	<macro name="IDIP_DO_MODIFY_DIAMOND_RSP" value="0x1002" desc="修改钻石（代币）应答"/>

	<!--修改钻石（代币）请求-->
	<struct name="DoModifyDiamondReq" id="IDIP_DO_MODIFY_DIAMOND_REQ" desc="修改钻石（代币）请求">
		<entry name="AreaId" type="uint32" desc="服务器：微信（1），手Q（2），微信测试服（998），手Q测试服（999），预发布环境（99 ）      " test="1" isverify="true" isnull="true"/>
		<entry name="PlatId" type="uint8" desc="平台：IOS（0），安卓（1）" test="1" isverify="true" isnull="true"/>
		<entry name="OpenId" type="string" size="MAX_OPENID_LEN" desc="OpenId" test="732945400" isverify="false" isnull="true"/>
		<entry name="ModifyNum" type="int64" desc="修改数量：-减，+加" test="1" isverify="true" isnull="true"/>
		<entry name="Source" type="uint32" desc="渠道号，由前端生成，不需填写" test="11" isverify="true" isnull="true"/>
		<entry name="Serial" type="string" size="MAX_SERIAL_LEN" desc="流水号，由前端生成，不需要填写" test="M-PAYO-20140414124009-58382166" isverify="false" isnull="true"/>
	</struct>
	<!--修改钻石（代币）应答-->
	<struct name="DoModifyDiamondRsp" id="IDIP_DO_MODIFY_DIAMOND_RSP" desc="修改钻石（代币）应答">
		<entry name="Result" type="int32" desc="结果：成功(0)，玩家不存在(1)，失败(其他)" test="0" isverify="false" isnull="true"/>
		<entry name="RetMsg" type="string" size="MAX_RETMSG_LEN" desc="返回消息" test="success" isverify="false" isnull="true"/>
		<entry name="BeforeNum" type="uint32" desc="修改前钻石量" test="1" isverify="false" isnull="true"/>
		<entry name="AfterNum" type="uint32" desc="修改后钻石数量" test="1" isverify="false" isnull="true"/>
	</struct>
*/
// swagger:route POST /idip/modify_diamond idip idip_modify_diamond
//
// modify_diamond: 修改钻石（代币）
//
// IDIP命令编码 <br>
//   IDIP_DO_MODIFY_DIAMOND_REQ = 0x1001 <br>
//   IDIP_DO_MODIFY_DIAMOND_RSP = 0x1002 <br>
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
//       200: DoModifyDiamondRsp
func modify_diamond(c *gin.Context) {

	resp := DoModifyDiamondRsp{}
	resp.Body.Head.PacketLen = 0
	resp.Body.Head.Cmdid = DO_MODIFY_DIAMOND_RSP
	resp.Body.Head.ServiceName = "WorldOfLove"
	resp.Body.Head.SendTime = time.Now().Unix()
	resp.Body.Head.Version = 1
	resp.Body.Head.Authenticate = ""
	resp.Body.Head.Result = Ok
	resp.Body.Head.RetErrMsg = "ok"

	var request DoModifyDiamondReq
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

		url := fmt.Sprintf("%s/%s", game, "modify_diamond")

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

		SendIDIPTLog(int(AreaId), OpenId, Item_Diamond, int(request.Body.ModifyNum), request.Body.Serial, int(request.Body.Source),
			DO_MODIFY_DIAMOND_REQ, request.Head.Seqid, int(request.Body.PlatId))

		return
	}
}
