package bind

import (
	"github.com/gin-gonic/gin"
	"github.com/xykong/loveauth/utils"
	"github.com/xykong/loveauth/errors"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/services/card_id"
	"github.com/xykong/loveauth/services/sms"
	"github.com/xykong/loveauth/storage/model"
	"net/http"
	"strings"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"regexp"
)

func init() {
	handlers["/real/name"] = realName
	getHandlers["/real/info"] = realInfo
}

type RequestRealInfo struct {
	GlobalId int64 `form:"globalId" json:"globalId" binding:"required"`
}

type ResponseRealInfo struct {
	Body struct {
		Code       int64  `json:"code"`
		Message    string `json:"message"`
		IsRealName bool   `json:"isRealName"`
		IsAdult    bool   `json:"isAdult"`
	}
}

func realInfo(c *gin.Context) {
	var request RequestRealInfo

	// validation
	if err := c.Bind(&request); err != nil {
		utils.QuickReply(c, errors.ServerFail, "连接服务器异常")
		return
	}

	accountRealName := storage.QueryAccountRealName(request.GlobalId)

	resp := ResponseRealInfo{}
	resp.Body.Code = int64(errors.Ok)
	resp.Body.Message = "获取成功"
	resp.Body.IsRealName = false
	resp.Body.IsAdult = false
	if accountRealName != nil && accountRealName.CardId != "" {
		resp.Body.IsRealName = true
		resp.Body.IsAdult = card_id.IsAdult(accountRealName.CardId)
	}

	c.JSON(http.StatusOK, resp.Body)
	return
}

type RequestRealName struct {
	AccessTokenValue string `form:"accessTokenValue" json:"accessTokenValue" binding:"required"`
	Mobile           string `form:"mobile" json:"mobile" binding:"required"`
	Token            string `form:"token" json:"token" binding:"required"`
	Name             string `form:"name" json:"name" binding:"required"`
	CardId           string `form:"cardId" json:"cardId" binding:"required"`
}

type ResponseRealName struct {
	Body struct {
		Code    int64  `json:"code"`
		Message string `json:"message"`
		IsAdult bool   `json:"isAdult"`
	}
}

func realName(c *gin.Context) {
	var request RequestRealName

	// validation
	if err := c.BindJSON(&request); err != nil {
		utils.QuickReply(c, errors.ServerFail, "连接服务器异常")
		return
	}

	//0. 验证登录信息是否有效
	tokenRecord, err := storage.QueryAccessToken(request.Token)
	if err != nil {
		utils.QuickReply(c, errors.ServerFail, "连接服务器异常")
		return
	}

	if tokenRecord.GlobalId == 0 {
		utils.QuickReply(c, errors.ServerFail, "连接服务器异常")
		return
	}

	//1. 验证身份证号码格式
	ok := card_id.Validate(request.CardId)
	if !ok {
		utils.QuickReply(c, errors.CardIdError, "请填写正确的身份证号")
		return
	}

	realName := strings.TrimSpace(request.Name)
	if realName == "" {
		utils.QuickReply(c, errors.RealNameError, "请填写真实姓名")
		return
	}

	nameReguler := "^[^0-9a-zA-Z`~!@#$^&*()=|{}':;',\\[\\].<>/?~！@#￥……&*（）——|{}【】‘；：”“'。，、？]+$" //名字为：中文-_·
	m, _ := regexp.MatchString(nameReguler, realName)
	if !m {
		utils.QuickReply(c, errors.RealNameError, "请填写真实姓名")
		return
	}

	//2. 验证手机验证码
	err = sms.VerifyMobileCode(request.Mobile, request.AccessTokenValue, true)
	if err != nil {
		ec, ok := err.(*errors.Type)
		if !ok { //系统错误
			utils.QuickReply(c, errors.ServerFail, "连接服务器异常")
			return
		}

		//逻辑类错误
		utils.QuickReply(c, ec.Code, ec.Message)
		return
	}

	//3. 是否已经实名
	accountRealName := storage.QueryAccountRealName(tokenRecord.GlobalId)
	if accountRealName != nil && accountRealName.GlobalId != 0 { //已经认证
		utils.QuickReply(c, errors.RealNameRepeat, "实名认证已完成")
		return
	}

	//4. 记录实名验证时的手机号
	accountRealName = &model.AccountRealName{
		GlobalId:       tokenRecord.GlobalId,
		RealNameMobile: request.Mobile,
		RealName:       realName,
		CardId:         request.CardId,
	}
	err = storage.Insert(storage.AuthDatabase(), accountRealName)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":             err,
			"accountRealName": accountRealName,
		}).Error("Create accountRealName Insert failed.")
		utils.QuickReply(c, errors.ServerFail, "连接服务器异常")
		return
	}

	//5. 检测账号有没有绑定手机号
	resp := ResponseRealName{}
	resp.Body.Code = int64(errors.Ok)
	resp.Body.Message = "实名认证完成"
	resp.Body.IsAdult = card_id.IsAdult(accountRealName.CardId)
	var mobileInfo model.AuthMobile
	err = storage.QueryVendorInfoByGlobalId(tokenRecord.GlobalId, &mobileInfo)
	if err != nil || mobileInfo.GlobalId != 0 { //已绑定
		c.JSON(http.StatusOK, resp.Body)
		return
	}

	//6. 检测手机号有没有绑定账号，都没则将手机号绑定到当前账号
	var destMobile model.AuthMobile
	err = storage.QueryVendorByOpenId(request.Mobile, &destMobile)
	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusOK, resp.Body)
		return
	}

	if destMobile.GlobalId != 0 { //手机号已绑定其它账号
		c.JSON(http.StatusOK, resp.Body)
		return
	}

	authMobile := &model.AuthMobile{
		GlobalId: tokenRecord.GlobalId,
		OpenId:   request.Mobile,
		Platform: tokenRecord.Platform,
	}
	err = storage.Insert(storage.AuthDatabase(), authMobile)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"authMobile": authMobile,
		}).Warn("realName Insert failed.")
	}

	//实名时绑定手机成功则删除设备绑定
	var authDevice model.AuthDevice
	err = storage.DeleteVendorByGlobalId(tokenRecord.GlobalId, &authDevice)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":      err,
			"GlobalId": tokenRecord.GlobalId,
		}).Error("realName delete device bind failed.")
	}

	c.JSON(http.StatusOK, resp.Body)
	return
}
