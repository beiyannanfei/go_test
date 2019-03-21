package idip

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
)

func init() {

	handlerFuncList[DO_BAN_ACCOUNT_REQ] =
		handlerFuncInfo{"ban_account", ban_account}
}

const DO_BAN_ACCOUNT_REQ = 0x1013
const DO_BAN_ACCOUNT_RSP = 0x1014

//
// 封号请求
// Binding from JSON
type DoBanAccountReq struct {
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
		// 封号时长：-1 永久，*秒
		//
		BanTime int64 `json:"BanTime"`
		//
		// 封号原因
		//
		BanReason string `json:"BanReason"`
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
// swagger:parameters idip_ban_account
type DoBanAccountReqBodyParams struct {
	//
	// swagger:allOf
	// in: body
	DoBanAccountReq DoBanAccountReq
}

//
// 封号应答
// swagger:response DoBanAccountRsp
type DoBanAccountRsp struct {
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
	<entry
		cmd="10356000"
		req="IDIP_DO_BAN_ACCOUNT_REQ"
		rsp="IDIP_DO_BAN_ACCOUNT_RSP"
		req_value="area=AreaId#platid=PlatId#openid=OpenId#modify_num=ModifyNum#source=Source#serial=Serial"
		rsp_value="result=Result#error_info=RetMsg#before_num=BeforeNum#after_num=AfterNum"
		desc="封号"
	/>

    <macro name="IDIP_DO_BAN_ACCOUNT_REQ" value="0x1013" desc="封号请求"/>
    <macro name="IDIP_DO_BAN_ACCOUNT_RSP" value="0x1014" desc="封号应答"/>

	<!--封号请求-->
	<struct name="DoBanAccountReq" id="IDIP_DO_BAN_ACCOUNT_REQ" desc="封号请求">
		<entry name="AreaId" type="uint32" desc="服务器：微信（1），手Q（2），微信测试服（998），手Q测试服（999），预发布环境（99 ）      " test="1" isverify="true" isnull="true"/>
		<entry name="PlatId" type="uint8" desc="平台：IOS（0），安卓（1）" test="1" isverify="true" isnull="true"/>
		<entry name="OpenId" type="string" size="MAX_OPENID_LEN" desc="OpenId" test="732945400" isverify="false" isnull="true"/>
		<entry name="ModifyNum" type="int64" desc="修改数量：-减，+加" test="1" isverify="true" isnull="true"/>
		<entry name="Source" type="uint32" desc="渠道号，由前端生成，不需填写" test="11" isverify="true" isnull="true"/>
		<entry name="Serial" type="string" size="MAX_SERIAL_LEN" desc="流水号，由前端生成，不需要填写" test="M-PAYO-20140414124009-58382166" isverify="false" isnull="true"/>
	</struct>
	<!--封号应答-->
	<struct name="DoBanAccountRsp" id="IDIP_DO_BAN_ACCOUNT_RSP" desc="封号应答">
		<entry name="Result" type="int32" desc="结果：成功(0)，玩家不存在(1)，失败(其他)" test="0" isverify="false" isnull="true"/>
		<entry name="RetMsg" type="string" size="MAX_RETMSG_LEN" desc="返回消息" test="success" isverify="false" isnull="true"/>
		<entry name="BeforeNum" type="uint32" desc="修改前钻石量" test="1" isverify="false" isnull="true"/>
		<entry name="AfterNum" type="uint32" desc="修改后钻石数量" test="1" isverify="false" isnull="true"/>
	</struct>
*/
// swagger:route POST /idip/ban_account idip idip_ban_account
//
// ban_account: 封号
//
// IDIP命令编码 <br>
//   IDIP_DO_BAN_ACCOUNT_REQ = 0x1013 <br>
//   IDIP_DO_BAN_ACCOUNT_RSP = 0x1014 <br>
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
//       200: DoBanAccountRsp
func ban_account(c *gin.Context) {

	resp := DoBanAccountRsp{}
	resp.Body.Head.PacketLen = 0
	resp.Body.Head.Cmdid = DO_BAN_ACCOUNT_RSP
	resp.Body.Head.ServiceName = "WorldOfLove"
	resp.Body.Head.SendTime = time.Now().Unix()
	resp.Body.Head.Version = 1
	resp.Body.Head.Authenticate = ""
	resp.Body.Head.Result = Ok
	resp.Body.Head.RetErrMsg = "ok"

	var request DoBanAccountReq
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

		url := fmt.Sprintf("%s/%s", game, "ban_account")

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

		err = storage.UpdateAccountState(globalId, getBanAccountMapInfo(request.Body.BanTime))
		if err != nil {

			resp.Body.Head.Result = StorageErr
			resp.Body.Head.RetErrMsg = fmt.Sprintf("[%v,%v]Ban account state failed: %v", AreaId, PlatId, err)
			resp.Body.Body.Result = StorageErr
			resp.Body.Body.RetMsg = fmt.Sprintf("[%v,%v]Ban account state failed: %v", AreaId, PlatId, err)
			c.JSON(http.StatusOK, resp.Body)

			return
		}

		storage.RemoveToken(globalId)

		resp.Body.Body.RetMsg = "ok"
		c.JSON(http.StatusOK, resp.Body)

		return
	}
}

func getBanAccountMapInfo(banTime int64) map[string]interface{} {

	data := map[string]interface{}{}
	data["state"] = model.BanDied

	if banTime != -1 {

		data["state"] = model.Banned
		data["un_ban_time"] = time.Now().Unix() + banTime
	}

	return data
}
