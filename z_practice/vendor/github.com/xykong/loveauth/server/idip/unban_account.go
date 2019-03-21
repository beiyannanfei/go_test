package idip

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
)

func init() {

	handlerFuncList[DO_UNBAN_ACCOUNT_REQ] =
		handlerFuncInfo{"unban_account", unban_account}
}

const DO_UNBAN_ACCOUNT_REQ = 0x1011
const DO_UNBAN_ACCOUNT_RSP = 0x1012

//
// 解封请求
// Binding from JSON
type DoUnbanAccountReq struct {
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
// swagger:parameters idip_unban_account
type DoUnbanAccountReqBodyParams struct {
	//
	// swagger:allOf
	// in: body
	DoUnbanAccountReq DoUnbanAccountReq
}

//
// 解封应答
// swagger:response DoUnbanAccountRsp
type DoUnbanAccountRsp struct {
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
    <entry cmd="10356008" req="IDIP_DO_UNBAN_ACCOUNT_REQ" rsp="IDIP_DO_UNBAN_ACCOUNT_RSP" req_value="area=AreaId#platid=PlatId#openid=OpenId#source=Source#serial=Serial" rsp_value="result=Result#error_info=RetMsg" desc="解封"/>

    <macro name="IDIP_DO_UNBAN_ACCOUNT_REQ" value="0x1011" desc="解封请求"/>
    <macro name="IDIP_DO_UNBAN_ACCOUNT_RSP" value="0x1012" desc="解封应答"/>

	<!--解封请求-->
	<struct name="DoUnbanAccountReq" id="IDIP_DO_UNBAN_ACCOUNT_REQ" desc="解封请求">
		<entry name="AreaId" type="uint32" desc="服务器：微信（1），手Q（2），微信测试服（998），手Q测试服（999），预发布环境（99 ）      " test="1" isverify="true" isnull="true"/>
		<entry name="PlatId" type="uint8" desc="平台：IOS（0），安卓（1）" test="1" isverify="true" isnull="true"/>
		<entry name="OpenId" type="string" size="MAX_OPENID_LEN" desc="OpenId" test="732945400" isverify="false" isnull="true"/>
		<entry name="ModifyNum" type="int64" desc="修改数量：-减，+加" test="1" isverify="true" isnull="true"/>
		<entry name="Source" type="uint32" desc="渠道号，由前端生成，不需填写" test="11" isverify="true" isnull="true"/>
		<entry name="Serial" type="string" size="MAX_SERIAL_LEN" desc="流水号，由前端生成，不需要填写" test="M-PAYO-20140414124009-58382166" isverify="false" isnull="true"/>
	</struct>
	<!--解封应答-->
	<struct name="DoUnbanAccountRsp" id="IDIP_DO_UNBAN_ACCOUNT_RSP" desc="解封应答">
		<entry name="Result" type="int32" desc="结果：成功(0)，玩家不存在(1)，失败(其他)" test="0" isverify="false" isnull="true"/>
		<entry name="RetMsg" type="string" size="MAX_RETMSG_LEN" desc="返回消息" test="success" isverify="false" isnull="true"/>
		<entry name="BeforeNum" type="uint32" desc="修改前钻石量" test="1" isverify="false" isnull="true"/>
		<entry name="AfterNum" type="uint32" desc="修改后钻石数量" test="1" isverify="false" isnull="true"/>
	</struct>
*/
// swagger:route POST /idip/unban_account idip idip_unban_account
//
// unban_account: 解封
//
// IDIP命令编码 <br>
//   IDIP_DO_UNBAN_ACCOUNT_REQ = 0x1011 <br>
//   IDIP_DO_UNBAN_ACCOUNT_RSP = 0x1012 <br>
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
//       200: DoUnbanAccountRsp
func unban_account(c *gin.Context) {

	resp := DoUnbanAccountRsp{}
	resp.Body.Head.PacketLen = 0
	resp.Body.Head.Cmdid = DO_UNBAN_ACCOUNT_RSP
	resp.Body.Head.ServiceName = "WorldOfLove"
	resp.Body.Head.SendTime = time.Now().Unix()
	resp.Body.Head.Version = 1
	resp.Body.Head.Authenticate = ""
	resp.Body.Head.Result = Ok
	resp.Body.Head.RetErrMsg = "ok"

	var request DoUnbanAccountReq
	// validation
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusOK, resp.Body)
		return
	}

	AreaId := request.Body.AreaId
	PlatId := request.Body.PlatId
	OpenId := request.Body.OpenId

	resp.Body.Head.Seqid = request.Head.Seqid

	globalId := getGlobalIdByOpenIdAreaId(OpenId, AreaId)
	if globalId == 0 {

		resp.Body.Head.Result = UserNotFoundErr
		resp.Body.Head.RetErrMsg = fmt.Sprintf("[%v,%v]Account not found", AreaId, PlatId)
		resp.Body.Body.Result = UserNotFoundErr
		resp.Body.Body.RetMsg = fmt.Sprintf("[%v,%v]Account not found", AreaId, PlatId)
		c.JSON(http.StatusOK, resp.Body)

		return
	}
	//change mysql table accounts
	err := storage.UpdateAccountState(globalId, map[string]interface{}{"state": model.Active, "un_ban_time": 0})
	if err != nil {

		resp.Body.Head.Result = StorageErr
		resp.Body.Head.RetErrMsg = fmt.Sprintf("[%v,%v]Unban account state failed: %v", AreaId, PlatId, err)
		resp.Body.Body.Result = StorageErr
		resp.Body.Body.RetMsg = fmt.Sprintf("[%v,%v]Unban account state failed: %v", AreaId, PlatId, err)
		c.JSON(http.StatusOK, resp.Body)

		return
	}

	resp.Body.Body.RetMsg = "ok"
	c.JSON(http.StatusOK, resp.Body)
}
