package bind

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/utils"
	"github.com/xykong/loveauth/errors"
	"regexp"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/services/pic_code"
	"github.com/xykong/loveauth/services/sms"
	"github.com/xykong/loveauth/settings"
	"net/http"
)

func init() {
	handlers["/code"] = code
	handlers["/code/pic/refresh"] = refreshPic
}

type RequestCode struct {
	Mobile  string `form:"Mobile" json:"Mobile" binding:"required"`
	Type    int    `form:"Type" json:"Type"`
	PicCode string `form:"PicCode" json:"PicCode"`
}

type ResponseCode struct {
	Body struct {
		Code      int64  `json:"code"`
		Message   string `json:"message"`
		Base64Png string `json:"base64Png"`
	}
}

/*
ID  类型  申请时间  模板名称  内容  状态  操作
237072  普通短信  2018-11-28 17:04:53  绑定手机验证码短信  验证码：{1}，您正在绑定手机，验证码{2}分钟内有效，请勿泄露。祝体验愉快！
待审核
编辑删除
237071  普通短信  2018-11-28 17:04:09  实名注册验证码短信  验证码：{1}，您正在进行实名认证，验证码{2}分钟内有效，请勿泄露。祝体验愉快！
待审核
编辑删除
237066  普通短信  2018-11-28 17:03:25  登录验证码短信  验证码：{1}，您正在进入恋爱的世界，验证码{2}分钟内有效，请勿泄露。祝体验愉快！
已审核 预计5分钟后生效
复制到国际模板删除
104975  普通短信  2018-04-08 21:54:35  登录验证码  {1}为您的登录验证码，请于{2}分钟内填写。如非本人操作，请忽略本短信。
已通过
复制到国际模板删除
*/
func code(c *gin.Context) {
	var request RequestCode

	// validation
	if err := c.BindJSON(&request); err != nil {
		utils.QuickReply(c, errors.Failed, "code BindJSON failed: %v", err)
		return
	}

	if request.Mobile == "" {
		utils.QuickReply(c, errors.Failed, "请输入正确的手机号码")
		return
	}

	regexpReguler := `^1[3-9]\d{9}$`
	if m, _ := regexp.MatchString(regexpReguler, request.Mobile); !m {
		logrus.WithFields(logrus.Fields{
			"Mobile": request.Mobile,
		}).Error("code input mobile illegal")
		utils.QuickReply(c, errors.Failed, "请输入正确的手机号码")
		return
	}

	resp := ResponseCode{}
	var picReults, base64Png string
	smsToken, smsErr := storage.QuerySMSToken(request.Mobile)
	if err := storage.SMSSendCheck(request.Mobile); err != nil { //验证是否可以发送验证码
		ec, ok := err.(*errors.Type)
		if !ok { //系统累错误
			utils.QuickReply(c, errors.Failed, "code SMSSendCheck error: %v", err)
			return
		}

		if ec.Code == errors.SendSMSFrequently {
			utils.QuickReply(c, ec.Code, "服务器有点忙呢，请稍后几秒再试一次吧")
			return
		}

		if ec.Code == errors.SendSMSOneDayLimit {
			utils.QuickReply(c, ec.Code, "今日获取验证码次数已达上限")
			return
		}

		if ec.Code == errors.SendSMSNeedPicture { //需要图片验证码
			//1 第一次触发，给前端生成验证码
			//2 用户已经输入验证，但是验证错误
			if smsToken.PicToken == "" || (smsErr == nil && smsToken.PicToken != "" && smsToken.PicToken != request.PicCode) {
				picReults, base64Png = pic_code.Generate()
				storage.WriteSMSToken(request.Mobile, picReults, "")
				resp.Body.Code = int64(errors.PictureCodeError)
				if smsToken.PicToken == "" {
					resp.Body.Code = int64(errors.SendSMSNeedPicture)
				}
				resp.Body.Message = "请输入正确的图形验证码"
				resp.Body.Base64Png = base64Png
				c.JSON(http.StatusOK, resp.Body)
				return
			}
			//3 用户输入正确验证码,直接发送
		}
	}

	if request.PicCode != "" && request.PicCode != smsToken.PicToken { //当出现图形验证码后必须输入正确才通过
		picReults, base64Png = pic_code.Generate()
		storage.WriteSMSToken(request.Mobile, picReults, "")
		resp.Body.Code = int64(errors.PictureCodeError)
		resp.Body.Message = "请输入正确的图形验证码"
		resp.Body.Base64Png = base64Png
		c.JSON(http.StatusOK, resp.Body)
		return
	}

	var codeTemplateId int //短信模板
	switch request.Type {
	case sms.Bind:
		codeTemplateId = settings.GetInt("sms", "qcloud.TemplateBind")
	case sms.RealName:
		codeTemplateId = settings.GetInt("sms", "qcloud.TemplateRealName")
	case sms.LoginCode:
		codeTemplateId = settings.GetInt("sms", "qcloud.TemplateLogin")
	default:
		codeTemplateId = settings.GetInt("sms", "qcloud.TemplateAuth")
	}

	token := utils.GenerateAuthCode() //短信验证码

	whiteMobiles := settings.GetStringSlice("loveauth_white_list", "mobiles")
	for _, m := range whiteMobiles {	//如果是白名单手机号则验证码发送白名单验证码
		if request.Mobile == m[:11] {
			token = m[12:]
		}
	}

	sms.SendAuth(token, request.Mobile, codeTemplateId)
	storage.WriteSMSToken(request.Mobile, "", token)
	storage.PushSMSSend(request.Mobile) //记录发送内容，一天内发送总条数
	resp.Body.Code = int64(errors.Ok)
	resp.Body.Message = "验证码已发送"
	c.JSON(http.StatusOK, resp.Body)
	return
}

type RequestRefresh struct {
	Mobile string `form:"Mobile" json:"Mobile" binding:"required"`
}

type ResponseRefresh struct {
	Body struct {
		Code      int64  `json:"code"`
		Message   string `json:"message"`
		Base64Png string `json:"base64Png"`
	}
}

func refreshPic(c *gin.Context) {
	var request RequestRefresh

	// validation
	if err := c.BindJSON(&request); err != nil {
		utils.QuickReply(c, errors.ServerFail, "连接服务器异常")
		return
	}

	regexpReguler := `^1[3-9]\d{9}$`
	if m, _ := regexp.MatchString(regexpReguler, request.Mobile); !m {
		logrus.WithFields(logrus.Fields{
			"Mobile": request.Mobile,
		}).Error("code input mobile illegal")
		utils.QuickReply(c, errors.Failed, "请输入正确的手机号码")
		return
	}

	picReults, base64Png := pic_code.Generate()
	storage.WriteSMSToken(request.Mobile, picReults, "")

	resp := ResponseRefresh{}
	resp.Body.Code = int64(errors.Ok)
	resp.Body.Message = "刷新成功"
	resp.Body.Base64Png = base64Png
	c.JSON(http.StatusOK, resp.Body)
	return
}
