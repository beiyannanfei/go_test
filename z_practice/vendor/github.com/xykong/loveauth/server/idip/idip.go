package idip

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	waring "github.com/xykong/loveauth/server/send_mail"
	"github.com/xykong/loveauth/settings"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"github.com/xykong/tlog"
)

type Header struct {
	//
	// 包长
	//
	PacketLen int `json:"PacketLen"`
	//
	// 命令ID
	//
	Cmdid int `json:"Cmdid"`
	//
	// 流水号
	//
	Seqid int `json:"Seqid"`
	//
	// 服务名
	//
	// maxLength: 16
	ServiceName string `json:"ServiceName"`
	//
	// 发送时间YYYYMMDD对应的整数
	//
	SendTime int64 `json:"SendTime"`
	//
	// 版本号
	//
	Version int `json:"Version"`
	//
	// 加密串
	//
	// maxLength: 32
	Authenticate string `json:"Authenticate"`
	//
	// 错误码,返回码类型：
	// <br>  0：处理成功，需要解开包体获得详细信息,
	// <br>  1：处理成功，但包体返回为空，不需要处理包体（eg：查询用户角色，用户角色不存在等）,
	// <br> -1: 网络通信异常,
	// <br> -2：超时,
	// <br> -3：数据库操作异常,
	// <br> -4：API返回异常,
	// <br> -5：服务器忙,
	// <br> -6：其他错误,小于
	// <br> -100 ：用户自定义错误，需要填写RetErrMsg
	//
	Result int `json:"Result"`
	//
	// 错误信息
	//
	// maxLength: 100
	RetErrMsg string `json:"RetErrMsg"`
}

const (
	Ok = iota
	UserNotFoundErr
	NetWorkErr = -1
	//TimeOutErr      = -2
	StorageErr = -3
	PostApiErr = -4
	//ServerErr       = -5
	//OtherErr        = -6
	CmdNotFoundErr  = -7
	ZoneNotFoundErr = -8
	//DefineErr       = -100
)

const RandReqCount = 1

const (
	Item_Diamond = 111
	Item_Power   = 101
)

type ExecRequest struct {
	Head Header `json:"head"`
	Body struct {
		AreaId int    `json:"AreaId"`
		PlatId int    `json:"PlatId"`
		OpenId string `json:"OpenId"`
		Source int    `json:"Source"`
		Serial string `json:"Serial"`
	} `json:"body"`
}

//
// in: body
// swagger:parameters idip_exec
type ExecRequestBodyParams struct {
	//
	// swagger:allOf
	// in: body
	ExecRequest ExecRequest
}

// A idip.Response is an response message to client.
// swagger:response ExecResponse
type ExecResponse struct {
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

type handlerFuncInfo struct {
	relativePath string
	handlerFunc  gin.HandlerFunc
}

var handlerFuncList = make(map[int]handlerFuncInfo)
var idipV2 = make(map[string]gin.HandlerFunc)

//var games = make(map[uint32]map[uint8]string)
var games = make(map[uint32]map[uint8][]string)
var configs = make(map[uint32]model.Vendor)

func Start(group *gin.RouterGroup) {

	group.POST("/exec", exec)

	for _, info := range handlerFuncList {
		group.POST(info.relativePath, info.handlerFunc)
	}

	for p, f := range idipV2 {

		group.POST(p, func(c *gin.Context) {

			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Methods", "*")
			c.Header("Access-Control-Allow-Headers", "Content-Type,Access-Token")
			c.Header("Access-Control-Expose-Headers", "*")
		}, f)
	}

	g := settings.GetStringMap("tencent", "idip.games")

	for a, v := range g {
		value, err := strconv.Atoi(a)
		if err != nil {
			continue
		}

		AreaId := uint32(value)

		for p, s := range v.(map[string]interface{}) {

			value, err := strconv.Atoi(p)
			if err != nil {
				continue
			}

			PlatId := uint8(value)

			if _, ok := games[AreaId]; !ok {
				games[AreaId] = make(map[uint8][]string)
			}

			if s != nil {

				for _, u := range s.([]interface{}) {

					games[AreaId][PlatId] = append(games[AreaId][PlatId], u.(string))
				}
			}
			//games[AreaId][PlatId] = s.(string)
		}
	}

	cs := settings.GetStringMapStringSlice("tencent", "idip.configs")

	for v, as := range cs {

		vendor := model.Vendor(v)

		for _, a := range as {

			value, err := strconv.Atoi(a)
			if err != nil {

				continue
			}

			areaId := uint32(value)
			configs[areaId] = vendor
		}
	}
}

// swagger:route POST /idip/exec idip idip_exec
//
// Refresh access token by refresh token
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
//       200: ExecResponse
func exec(c *gin.Context) {

	var prefix = []byte("data_packet=")

	n, err := c.Request.Body.Read(prefix)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"n":     n,
			"error": err,
		}).Error("exec prefix")
	}

	// Read the Body content
	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)

		// Restore the io.ReadCloser to its original state
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	//logrus.WithFields(logrus.Fields{
	//	"prefix": prefix,
	//	"body":   string(utils.I(ioutil.ReadAll(c.Request.Body))[0].([]byte)),
	//	"n":      n,
	//	"error":  err,
	//}).Warn("exec debug")

	resp := ExecResponse{}
	resp.Body.Head.PacketLen = 0
	resp.Body.Head.Cmdid = DO_MODIFY_DIAMOND_RSP
	resp.Body.Head.ServiceName = "WorldOfLove"
	resp.Body.Head.SendTime = time.Now().Unix()
	resp.Body.Head.Version = 1
	resp.Body.Head.Authenticate = ""
	resp.Body.Head.Result = Ok
	resp.Body.Head.RetErrMsg = "ok"

	var request ExecRequest

	// validation
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusOK, resp.Body)
		return
	}

	if _, ok := handlerFuncList[request.Head.Cmdid]; !ok {

		resp.Body.Head.Seqid = request.Head.Seqid
		resp.Body.Head.Result = CmdNotFoundErr
		resp.Body.Head.RetErrMsg = "Command not support"
		resp.Body.Body.Result = CmdNotFoundErr
		resp.Body.Body.RetMsg = "Command not support"
		c.JSON(http.StatusOK, resp.Body)

		return
	}

	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	handlerFuncList[request.Head.Cmdid].handlerFunc(c)
}

func PostRequest(url string, body interface{}) ([]byte, error) {

	start := time.Now()
	jsonValue, _ := json.Marshal(body)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	elapsed := time.Since(start)

	logrus.WithFields(logrus.Fields{
		"elapsed":  elapsed,
		"request":  url,
		"body":     string(jsonValue),
		"response": string(respBody),
	}).Info("Idip post request.")

	return respBody, nil
}

func getGlobalIdByOpenIdAreaId(openId string, areaId uint32) int64 {

	vendor := areaIdToVendor(areaId)
	if vendor == "" {

		return 0
	}

	account := storage.QueryAccountByOpenId(openId, vendor)
	if account == nil {

		return 0
	}

	return account.GlobalId
}

func areaIdToVendor(areaId uint32) model.Vendor {

	if vendor, ok := configs[areaId]; ok {

		return vendor
	}

	switch areaId {
	case 1: //wechat

		return model.VendorMsdkWechat
	case 2: //qq

		return model.VendorMsdkQQ
	case 998: //wechat_test

		return model.VendorMsdkWechat
	case 999: //qq_test

		return model.VendorMsdkQQ
	case 98: //wechat_iosreview

		return model.VendorMsdkWechat
	case 99: //qq_iosreview

		return model.VendorMsdkQQ
	case 91: //test

		return model.VendorDevice
	case 804: //ios_review

		return model.VendorMsdkGuest
	default:

	}

	return ""
}

func randGames(areaId uint32, platId uint8) string {

	count := len(games[areaId][platId])
	if count <= 0 {

		return ""
	}

	return games[areaId][platId][rand.Intn(count)]
}

//<struct  name="IDIPFLOW" version="1" desc="IDIP接口日志">
//	<entry  name="dtEventTime" type="datetime" desc="(必填)游戏事件的时间, 格式 YYYY-MM-DD HH:MM:SS" />
//	<entry  name="Area_id" type="int" desc="大区id，没有大区可以去掉"/>
//	<entry  name="vopenid" type="string" size="64" desc="(必填)用户OPENID号，如果是端游改为iuin填QQ号" />
//	<entry  name="Item_id" type="int" defaultvalue="0" desc="(必填)道具id，非道具填0" />
//	<entry  name="Item_num" type="int" desc="(必填)涉及的物品数量" />
//	<entry  name="Serial" type="string" size="64" desc="(必填)流水号前端传入，默认0" />
//	<entry  name="Source" type="int" desc="(必填)渠道id前端传入，默认0" />
//	<entry  name="Cmd" type="int" desc="(必填)接口指令id，需要转为十进制" />
//	<entry  name="Seqid" type="int" desc="(必填)流水id" />
//	<entry  name="PlatId" type="int" desc="ios 0 /android 1"/>
//  <entry name="Format" type="string" size="64" defaultvalue="1.0.0" desc="Format version of this log"/>
//</struct>
func SendIDIPTLog(areaId int, openId string, itemId int, itemNum int, serial string, sourceId int, cmd int, seqId int, platId int) {

	cmdStr := strconv.Itoa(cmd)
	cmdD, err := strconv.Atoi(cmdStr)
	if nil != err {
		logrus.Println("SendIDIPTLog ", err)
	}

	dtEventTime := time.Now().Format("2006-01-02 15:04:05")

	format := settings.GetString("tencent", "tlog.format")
	address := settings.GetString("tencent", "tlog.address")
	if address != "" {

		tlog.LogRaw(fmt.Sprintf("IDIPFLOW|%s|%d|%s|%d|%d|%s|%d|%d|%d|%d|%s\n",
			dtEventTime, areaId, openId, itemId, itemNum, serial, sourceId, cmdD, seqId, platId, format))
	}

	logrus.WithFields(logrus.Fields{
		"dtEventTime": dtEventTime,
		"Area_id":     areaId,
		"vopenid":     openId,
		"Item_id":     itemId,
		"Item_num":    itemNum,
		"Serial":      serial,
		"Source":      sourceId,
		"Cmd":         cmdD,
		"Seqid":       seqId,
		"PlatId":      platId,
		"FlowName":    "IDIPFLOW",
		"Format":      format,
	}).Info("Tlog IDIPFLOW")
}

func sendWaringMail(uid int64, gm, channel, title, mailContent, awardList string) {

	defer func() {

		if err := recover(); err != nil {

			logrus.WithFields(logrus.Fields{
				"err": err,
			}).Error("idip sendwaringmail err")
		}
	}()
	var uidStr string
	if uid == 0 {

		uidStr = ""
	} else {

		uidStr = fmt.Sprintf("%d", uid)
	}

	var body = `
<html>
    <body>
		<h3>gm邮件</h3>
    	<p>操作员：</p>
		<p>&nbsp; &nbsp; %s</p>
    	<p>玩家id：</p>
		<p>&nbsp; &nbsp; %s</p>
    	<p>渠道：</p>
		<p>&nbsp; &nbsp; %s</p>
    	<p>邮件标题：</p>
		<p>&nbsp; &nbsp; %s</p>
    	<p>邮件内容：</p>
		<p>&nbsp; &nbsp; %s</p>
		<p>邮件道具：</p>
		<p>&nbsp; &nbsp; %s</p>
    </body>
</html>
`
	content := fmt.Sprintf(body, gm, uidStr, channel, title, mailContent, awardList)

	waring.SendMail("gm邮件", content, "gm", settings.GetString("waring", "gmSenderName"))
}
