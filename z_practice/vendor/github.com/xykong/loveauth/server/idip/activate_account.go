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

	handlerFuncList[DO_ACTIVATE_ACCOUNT_REQ] =
		handlerFuncInfo{"activate_account", activate_account}
}

const DO_ACTIVATE_ACCOUNT_REQ = 0x100f
const DO_ACTIVATE_ACCOUNT_RSP = 0x1010

//
// 激活帐号请求
// Binding from JSON
type DoActivateAccountReq struct {
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
// swagger:parameters idip_activate_account
type DoActivateAccountReqBodyParams struct {
	//
	// swagger:allOf
	// in: body
	DoActivateAccountReq DoActivateAccountReq
}

//
// 激活帐号应答
// swagger:response DoActivateAccountRsp
type DoActivateAccountRsp struct {
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
    <entry cmd="10356007" req="IDIP_DO_ACTIVATE_ACCOUNT_REQ" rsp="IDIP_DO_ACTIVATE_ACCOUNT_RSP" req_value="area=AreaId#platid=PlatId#openid=OpenId#source=Source#serial=Serial" rsp_value="result=Result#error_info=RetMsg" desc="激活帐号"/>

    <macro name="IDIP_DO_ACTIVATE_ACCOUNT_REQ" value="0x100f" desc="激活帐号请求"/>
    <macro name="IDIP_DO_ACTIVATE_ACCOUNT_RSP" value="0x1010" desc="激活帐号应答"/>

  <!--激活帐号请求-->
  <struct name="DoActivateAccountReq" id="IDIP_DO_ACTIVATE_ACCOUNT_REQ" desc="激活帐号请求">
    <entry name="AreaId" type="uint32" desc="服务器：微信（1），手Q（2），微信测试服（998），手Q测试服（999），预发布环境（99 ）      " test="1" isverify="true" isnull="true"/>
    <entry name="PlatId" type="uint8" desc="平台：IOS（0），安卓（1）" test="1" isverify="true" isnull="true"/>
    <entry name="OpenId" type="string" size="MAX_OPENID_LEN" desc="OpenId" test="732945400" isverify="false" isnull="true"/>
    <entry name="Source" type="uint32" desc="渠道号，由前端生成，不需填写" test="11" isverify="true" isnull="true"/>
    <entry name="Serial" type="string" size="MAX_SERIAL_LEN" desc="流水号，由前端生成，不需要填写" test="M-PAYO-20140414124009-58382166" isverify="false" isnull="true"/>
  </struct>
  <!--激活帐号应答-->
  <struct name="DoActivateAccountRsp" id="IDIP_DO_ACTIVATE_ACCOUNT_RSP" desc="激活帐号应答">
    <entry name="Result" type="int32" desc="结果：成功(0)，玩家不存在(1)，失败(其他)" test="0" isverify="false" isnull="true"/>
    <entry name="RetMsg" type="string" size="MAX_RETMSG_LEN" desc="返回消息" test="success" isverify="false" isnull="true"/>
  </struct>
*/
// swagger:route POST /idip/activate_account idip idip_activate_account
//
// activate_account: 激活帐号
//
// IDIP命令编码 <br>
//   IDIP_DO_ACTIVATE_ACCOUNT_REQ = 0x100f <br>
//   IDIP_DO_ACTIVATE_ACCOUNT_RSP = 0x1010 <br>
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
//       200: DoActivateAccountRsp
func activate_account(c *gin.Context) {

	resp := DoActivateAccountRsp{}
	resp.Body.Head.PacketLen = 0
	resp.Body.Head.Cmdid = DO_ACTIVATE_ACCOUNT_RSP
	resp.Body.Head.ServiceName = "WorldOfLove"
	resp.Body.Head.SendTime = time.Now().Unix()
	resp.Body.Head.Version = 1
	resp.Body.Head.Authenticate = ""
	resp.Body.Head.Result = Ok
	resp.Body.Head.RetErrMsg = "ok"

	var request DoActivateAccountReq
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

	err := storage.UpdateAccountState(globalId, map[string]interface{}{"state": model.Active, "un_ban_time": 0})
	if err != nil {

		resp.Body.Head.Result = StorageErr
		resp.Body.Head.RetErrMsg = fmt.Sprintf("[%v,%v]Activate account state failed: %v", AreaId, OpenId, err)
		resp.Body.Body.Result = StorageErr
		resp.Body.Body.RetMsg = fmt.Sprintf("[%v,%v]Activate account state failed: %v", AreaId, OpenId, err)
		c.JSON(http.StatusOK, resp.Body)

		return
	}

	resp.Body.Body.RetMsg = "ok"
	c.JSON(http.StatusOK, resp.Body)
}
