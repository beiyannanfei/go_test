package idip

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xykong/loveauth/storage"
)

func init() {

	handlerFuncList[DO_QUERY_USER_INFO_REQ] =
		handlerFuncInfo{"query_user_info", query_user_info}
}

const DO_QUERY_USER_INFO_REQ = 0x1015
const DO_QUERY_USER_INFO_RSP = 0x1016

//
// 当前个人信息查询请求
// Binding from JSON
type DoQueryUserInfoReq struct {
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
// swagger:parameters idip_query_user_info
type DoQueryUserInfoReqBodyParams struct {
	//
	// swagger:allOf
	// in: body
	DoQueryUserInfoReq DoQueryUserInfoReq
}

//
// 当前个人信息查询应答
// swagger:response DoQueryUserInfoRsp
type DoQueryUserInfoRsp struct {
	// in: body
	Body struct {
		Head Header `json:"head"`
		Body struct {
			//
			// 结果：成功(0)，玩家不存在(1)，失败(其他)
			//
			Result int32 `json:"Result"`
			//
			// 返回消息
			//
			RetMsg string `json:"RetMsg"`
			//
			// 用户ID
			//
			Uid uint64 `json:"Uid"`
			//
			// 用户名称
			//
			UserName string `json:"UserName"`
			//
			// 渠道号
			//
			Source uint32 `json:"Source"`
			//
			// 成就点数
			//
			Achievement uint32 `json:"Achievement"`
			//
			// 注册时间
			//
			RegisterTime uint32 `json:"RegisterTime"`
			//
			// 累计登录时长
			//
			AccumLoginTime uint32 `json:"AccumLoginTime"`
			//
			// 最后登出时间
			//
			LastLogoutTime uint32 `json:"LastLogoutTime"`
			//
			// 当前是否在线（0是 1 否）
			//
			IsOnline uint8 `json:"IsOnline"`
			//
			// 钻石数量
			//
			DiamondNum uint32 `json:"DiamondNum"`
			//
			// 书券值
			//
			Power uint32 `json:"Power"`
		} `json:"body"`
	}
}

/*
    <entry cmd="10356010" req="IDIP_QUERY_USER_INFO_REQ" rsp="IDIP_QUERY_USER_INFO_RSP" req_value="area=AreaId#platid=PlatId#openid=OpenId" rsp_value="result=Result#error_info=RetMsg#diamond_num=DiamondNum#gold_num=GoldNum#power=Power#exp=Exp#level=Level#week_top_score=WeekTopScore#rank=Rank#history_top_score=HistoryTopScore#register_time=RegisterTime#accum_login_time=AccumLoginTime#last_logout_time=LastLogoutTime#is_online=IsOnline" desc="当前个人信息查询"/>

    <macro name="IDIP_QUERY_USER_INFO_REQ" value="0x1015" desc="当前个人信息查询请求"/>
    <macro name="IDIP_QUERY_USER_INFO_RSP" value="0x1016" desc="当前个人信息查询应答"/>

	<!--当前个人信息查询请求-->
	<struct name="DoQueryUserInfoReq" id="IDIP_DO_QUERY_USER_INFO_REQ" desc="当前个人信息查询请求">
		<entry name="AreaId" type="uint32" desc="服务器：微信（1），手Q（2），微信测试服（998），手Q测试服（999），预发布环境（99 ）      " test="1" isverify="true" isnull="true"/>
		<entry name="PlatId" type="uint8" desc="平台：IOS（0），安卓（1）" test="1" isverify="true" isnull="true"/>
		<entry name="OpenId" type="string" size="MAX_OPENID_LEN" desc="OpenId" test="732945400" isverify="false" isnull="true"/>
		<entry name="ModifyNum" type="int64" desc="修改数量：-减，+加" test="1" isverify="true" isnull="true"/>
		<entry name="Source" type="uint32" desc="渠道号，由前端生成，不需填写" test="11" isverify="true" isnull="true"/>
		<entry name="Serial" type="string" size="MAX_SERIAL_LEN" desc="流水号，由前端生成，不需要填写" test="M-PAYO-20140414124009-58382166" isverify="false" isnull="true"/>
	</struct>
	<!--当前个人信息查询应答-->
	<struct name="DoQueryUserInfoRsp" id="IDIP_DO_QUERY_USER_INFO_RSP" desc="当前个人信息查询应答">
		<entry name="Result" type="int32" desc="结果：成功(0)，玩家不存在(1)，失败(其他)" test="0" isverify="false" isnull="true"/>
		<entry name="RetMsg" type="string" size="MAX_RETMSG_LEN" desc="返回消息" test="success" isverify="false" isnull="true"/>
		<entry name="BeforeNum" type="uint32" desc="修改前钻石量" test="1" isverify="false" isnull="true"/>
		<entry name="AfterNum" type="uint32" desc="修改后钻石数量" test="1" isverify="false" isnull="true"/>
	</struct>
*/
// swagger:route POST /idip/query_user_info idip idip_query_user_info
//
// query_user_info: 当前个人信息查询
//
// IDIP命令编码 <br>
//   IDIP_DO_QUERY_USER_INFO_REQ = 0x1015 <br>
//   IDIP_DO_QUERY_USER_INFO_RSP = 0x1016 <br>
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
//       200: DoQueryUserInfoRsp
func query_user_info(c *gin.Context) {

	resp := DoQueryUserInfoRsp{}
	resp.Body.Head.PacketLen = 0
	resp.Body.Head.Cmdid = DO_QUERY_USER_INFO_RSP
	resp.Body.Head.ServiceName = "WorldOfLove"
	resp.Body.Head.SendTime = time.Now().Unix()
	resp.Body.Head.Version = 1
	resp.Body.Head.Authenticate = ""
	resp.Body.Head.Result = Ok
	resp.Body.Head.RetErrMsg = "ok"

	var request DoQueryUserInfoReq
	// validation
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusOK, resp.Body)
		return
	}

	resp.Body.Head.Seqid = request.Head.Seqid

	AreaId := request.Body.AreaId
	PlatId := request.Body.PlatId
	OpenId := request.Body.OpenId

	vendor := areaIdToVendor(AreaId)
	if vendor == "" {

		resp.Body.Head.Result = ZoneNotFoundErr
		resp.Body.Head.RetErrMsg = fmt.Sprintf("[%v,%v]Zone not found.", AreaId, PlatId)
		resp.Body.Body.Result = ZoneNotFoundErr
		resp.Body.Body.RetMsg = fmt.Sprintf("[%v,%v]Zone not found.", AreaId, PlatId)
		c.JSON(http.StatusOK, resp.Body)

		return
	}

	account := storage.QueryAccountByOpenId(OpenId, vendor)
	if account == nil {

		resp.Body.Head.Result = UserNotFoundErr
		resp.Body.Head.RetErrMsg = fmt.Sprintf("[%v,%v]Account not found", AreaId, PlatId)
		resp.Body.Body.Result = UserNotFoundErr
		resp.Body.Body.RetMsg = fmt.Sprintf("[%v,%v]Account not found", AreaId, PlatId)
		c.JSON(http.StatusOK, resp.Body)

		return
	}

	request.Body.GlobalId = account.GlobalId

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

		url := fmt.Sprintf("%s/%s", game, "query_user_info")

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

		resp.Body.Body.RegisterTime = uint32(account.CreatedAt.Unix())
		resp.Body.Body.LastLogoutTime = uint32(account.LogoutTime)
		resp.Body.Body.AccumLoginTime = uint32(account.AccumLoginTime)
		resp.Body.Body.IsOnline = 1
		if account.LoginTime > account.LogoutTime {

			currTime := time.Now().Unix()
			resp.Body.Body.IsOnline = 0
			resp.Body.Body.AccumLoginTime = uint32(account.AccumLoginTime + currTime - account.LoginTime)
		}

		resp.Body.Body.RetMsg = "ok"
		c.JSON(http.StatusOK, resp.Body)

		return
	}
}
