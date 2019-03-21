package idip

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	handlerFuncList[DO_SEND_ALL_MAIL_REQ] =
		handlerFuncInfo{"send_all_mail", send_all_mail}
}

const DO_SEND_ALL_MAIL_REQ = 0x1017
const DO_SEND_ALL_MAIL_RSP = 0x1018

//
// 发送带附件全服系统邮件请求
// Binding from JSON
type DoSendAllMailReq struct {
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
		// 邮件标题
		//
		MailTitle string `json:"MailTitle"`
		//
		// 邮件内容（3000字节）
		//
		MailContent string `json:"MailContent"`
		//
		// 附件：道具类型、道具ID、道具数量，安全起见，建议每次最多配置一种道具，如不填则视为纯文本邮件。的最大数量
		//
		Attatch_count uint32 `json:"Attatch_count"`
		//
		// 附件：道具类型、道具ID、道具数量，安全起见，建议每次最多配置一种道具，如不填则视为纯文本邮件。
		//
		Attatch []Attatch `json:"Attatch"`
		//
		// 生效时间
		//
		EffectTime uint32 `json:"EffectTime"`
		//
		// 失效时间
		//
		CloseTime uint32 `json:"CloseTime"`
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
// swagger:parameters idip_send_all_mail
type DoSendAllMailReqBodyParams struct {
	//
	// swagger:allOf
	// in: body
	DoSendAllMailReq DoSendAllMailReq
}

//
// 发送带附件全服系统邮件应答
// swagger:response DoSendAllMailRsp
type DoSendAllMailRsp struct {
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
   <entry cmd="10356011" req="IDIP_DO_SEND_ALL_MAIL_REQ" rsp="IDIP_DO_SEND_ALL_MAIL_RSP" req_value="area=AreaId#platid=PlatId#mail_title=MailTitle#mail_content=MailContent#attatch_list=Attatch_count|Attatch|#effect_time=EffectTime#source=Source#serial=Serial" rsp_value="result=Result#error_info=RetMsg" desc="发送带附件全服系统邮件"/>

   <macro name="IDIP_DO_SEND_ALL_MAIL_REQ" value="0x1017" desc="发送带附件全服系统邮件请求"/>
   <macro name="IDIP_DO_SEND_ALL_MAIL_RSP" value="0x1018" desc="发送带附件全服系统邮件应答"/>

  <!--发送带附件全服系统邮件请求-->
  <struct name="DoSendAllMailReq" id="IDIP_DO_SEND_ALL_MAIL_REQ" desc="发送带附件全服系统邮件请求">
    <entry name="AreaId" type="uint32" desc="服务器：微信（1），手Q（2），微信测试服（998），手Q测试服（999），预发布环境（99），ios审核服（98），体验服环境（97）      " test="1" isverify="true" isnull="true"/>
    <entry name="PlatId" type="uint8" desc="平台：IOS（0），安卓（1）" test="1" isverify="true" isnull="true"/>
    <entry name="MailTitle" type="string" size="MAX_MAILTITLE_LEN" desc="邮件标题" test="test" isverify="false" isnull="true"/>
    <entry name="MailContent" type="string" size="MAX_MAILCONTENT_LEN" desc="邮件内容（3000字节）" test="test" isverify="false" isnull="true"/>
    <entry name="Attatch_count" type="uint32" desc="附件：道具类型、道具ID、道具数量，安全起见，建议每次最多配置一种道具，如不填则视为纯文本邮件。的最大数量 " test="1" isverify="false" isnull="true"/>
    <entry name="Attatch" type="AttachList" size="MAX_ATTATCH_NUM" param="vector,struct,Attatch_count,1,|, " desc="附件：道具类型、道具ID、道具数量，安全起见，建议每次最多配置一种道具，如不填则视为纯文本邮件。" isverify="false" isnull="true"/>
    <entry name="EffectTime" type="uint32" desc="生效时间" test="1515640984" isverify="true" isnull="true"/>
    <entry name="Source" type="uint32" desc="渠道号，由前端生成，不需填写" test="11" isverify="true" isnull="true"/>
    <entry name="Serial" type="string" size="MAX_SERIAL_LEN" desc="流水号，由前端生成，不需要填写" test="M-PAYO-20140414124009-58382166" isverify="false" isnull="true"/>
  </struct>
  <!--发送带附件全服系统邮件应答-->
  <struct name="DoSendAllMailRsp" id="IDIP_DO_SEND_ALL_MAIL_RSP" desc="发送带附件全服系统邮件应答">
    <entry name="Result" type="int32" desc="结果：成功(0)，玩家不存在(1)，失败(其他)" test="0" isverify="false" isnull="true"/>
    <entry name="RetMsg" type="string" size="MAX_RETMSG_LEN" desc="返回消息" test="success" isverify="false" isnull="true"/>
  </struct>
*/

// swagger:route POST /idip/send_all_mail idip idip_send_all_mail
//
// send_all_mail: 发送带附件全服系统邮件
//
// IDIP命令编码 <br>
//   IDIP_DO_SEND_ALL_MAIL_REQ = 0x1017 <br>
//   IDIP_DO_SEND_ALL_MAIL_RSP = 0x1018 <br>
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
//       200: DoSendAllMailRsp
func send_all_mail(c *gin.Context) {
	resp := DoSendAllMailRsp{}
	resp.Body.Head.PacketLen = 0
	resp.Body.Head.Cmdid = DO_SEND_ALL_MAIL_RSP
	resp.Body.Head.ServiceName = "WorldOfLove"
	resp.Body.Head.SendTime = time.Now().Unix()
	resp.Body.Head.Version = 1
	resp.Body.Head.Authenticate = ""
	resp.Body.Head.Result = Ok
	resp.Body.Head.RetErrMsg = "ok"

	var request DoSendAllMailReq
	// validation
	if err := c.BindJSON(&request); err != nil {

		c.JSON(http.StatusOK, resp.Body)

		return
	}

	resp.Body.Head.Seqid = request.Head.Seqid

	AreaId := request.Body.AreaId
	PlatId := request.Body.PlatId

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

		url := fmt.Sprintf("%s/%s", game, "send_all_mail")

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

		for _, attatchInfo := range request.Body.Attatch {

			SendIDIPTLog(int(AreaId), "", int(attatchInfo.ItemID), int(attatchInfo.ItemNum), request.Body.Serial, int(request.Body.Source),
				DO_SEND_ITEM_MAIL_REQ, request.Head.Seqid, int(PlatId))
		}

		return
	}
}
