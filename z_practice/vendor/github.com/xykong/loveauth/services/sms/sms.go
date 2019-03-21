package sms

import (
	"github.com/xykong/qcloudsms_go"
	"github.com/xykong/loveauth/settings"
	"time"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/storage"
	"github.com/mediocregopher/radix.v2/redis"
	"github.com/xykong/loveauth/errors"
	"fmt"
)

const (
	Auth      = iota
	Bind
	RealName
	LoginCode
)

func SendAuth(code string, mobile string, templateId int) {

	opt := qcloudsms.NewOptions()
	opt.APPID = settings.GetString("sms", "qcloud.AppID")
	opt.APPKEY = settings.GetString("sms", "qcloud.AppKey")
	opt.SIGN = settings.GetString("sms", "qcloud.Sign")
	opt.HTTP.Timeout = 10 * time.Second
	opt.Debug = true

	var client = qcloudsms.NewClient(opt)

	var sm = qcloudsms.SMSSingleReq{
		TplID: templateId,
		Tel:   qcloudsms.SMSTel{Nationcode: "86", Mobile: mobile},
		Params: []string{
			code,
			"10",
		},
	}

	client.SendSMSSingle(sm)
}

func VerifyMobileCode(mobile string, code string, delCache bool) error {
	smsToken, err := storage.QuerySMSToken(mobile)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":    err,
			"mobile": mobile,
			"code":   code,
		}).Error("VerifyMobileCode QuerySMSToken failed.")

		if err == redis.ErrRespNil { //验证码已过期
			return errors.NewCodeString(errors.SMSCodeError, "请输入正确的验证码")
		}

		return err
	}

	if smsToken.Token != code { //短信验证码错误
		logrus.WithFields(logrus.Fields{
			"mobile":    mobile,
			"code":      code,
			"smsToken":  smsToken,
			"code-len":  len(code),
			"Token-len": len(smsToken.Token),
			"code-A":    fmt.Sprintf("AAAA%sAAAA", code),
			"Token-A":   fmt.Sprintf("AAAA%sAAAA", smsToken.Token),
		}).Warn("BindVerify smsToken.Token error")
		return errors.NewCodeString(errors.SMSCodeError, "请输入正确的验证码")
	}

	if delCache {
		//验证码正确，删除
		storage.DeleteSMSToken(mobile)
	}

	return nil
}
