package storage

import (
	"github.com/xykong/loveauth/storage/model"
	"github.com/sirupsen/logrus"
	"fmt"
	"encoding/json"
	"time"
	"github.com/xykong/loveauth/settings"
	"github.com/xykong/loveauth/errors"
	"github.com/mediocregopher/radix.v2/redis"
)

var RedisKeyTemplateSMS = "auth:sms_token:%v"
var RedisKeySMSSendTimeList = "auth:sms_send_time:%v"
var RedisKeySMSSendCount = "auth:sms_send_count:%v"

func WriteSMSToken(mobile string, picToken string, token string) {

	client, err := redisPool.Get()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("WriteSMSToken get redis client failed.")
		return
	}
	defer redisPool.Put(client)

	profile := &model.SMSToken{
		Mobile:            mobile,
		PicToken:          picToken,
		Token:             token,
		Timestamp:         time.Now().Unix(),
		ExpirationSeconds: settings.GetInt64("loveauth", "mobile.ExpirationSeconds"),
	}

	jsonString, _ := json.Marshal(profile)
	key := fmt.Sprintf(RedisKeyTemplateSMS, profile.Mobile)
	err = client.Cmd("SETEX", key, profile.ExpirationSeconds, jsonString).Err
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err,
			"profile": profile,
		}).Error("WriteSMSToken failed.")
	}
}

func QuerySMSToken(openId string) (*model.SMSToken, error) {

	client, err := redisPool.Get()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("QuerySMSToken get redis client failed.")
		return nil, err
	}
	defer redisPool.Put(client)

	var record model.SMSToken
	key := fmt.Sprintf(RedisKeyTemplateSMS, openId)
	err = cmdGetObject("QuerySMSToken", client, &record, "GET", key)

	return &record, err
}

func DeleteSMSToken(mobile string) {
	client, err := redisPool.Get()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":  err.Error(),
			"mobile": mobile,
		}).Error("DeleteSMSToken get redis client failed.")
	}
	defer redisPool.Put(client)

	key := fmt.Sprintf(RedisKeyTemplateSMS, mobile)
	err = client.Cmd("DEL", key).Err
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
			"key":   key,
		}).Error("DeleteSMSToken failed.")
	}
}

func PushSMSSend(mobile string) error {
	client, err := redisPool.Get()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":  err.Error(),
			"mobile": mobile,
		}).Error("PushSMSSend get redis client failed.")
		return err
	}
	defer redisPool.Put(client)

	key := fmt.Sprintf(RedisKeySMSSendTimeList, mobile)
	value := time.Now().Unix()
	if err := client.Cmd("LPUSH", key, value).Err; err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
			"key":   key,
			"value": value,
		}).Error("PushSMSSend redis LPUSH failed.")
		return err
	}

	NeedPicCodeSecondes := settings.GetInt64("loveauth", "mobile.NeedPicCodeSecondes")
	client.Cmd("EXPIRE", key, NeedPicCodeSecondes*2)
	NeedPicCodeTimes := settings.GetInt("loveauth", "mobile.NeedPicCodeTimes")
	client.Cmd("LTRIM", key, 0, NeedPicCodeTimes-1)

	key = fmt.Sprintf(RedisKeySMSSendCount, time.Now().Format("2006-01-02"))
	if err := client.Cmd("HINCRBY", key, mobile, 1).Err; err != nil {
		logrus.WithFields(logrus.Fields{
			"error":  err.Error(),
			"key":    key,
			"mobile": mobile,
		}).Error("PushSMSSend redis HINCRBY failed.")
		return err
	}
	client.Cmd("EXPIRE", key, 24*3600)

	return nil
}

func SMSSendCheck(mobile string) error {
	client, err := redisPool.Get()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":  err.Error(),
			"mobile": mobile,
		}).Error("QuerySMSToken get redis client failed.")
		return err
	}
	defer redisPool.Put(client)

	key := fmt.Sprintf(RedisKeySMSSendCount, time.Now().Format("2006-01-02"))
	sendCount, err := client.Cmd("HGET", key, mobile).Int64()
	if err != nil && err != redis.ErrRespNil {
		logrus.WithFields(logrus.Fields{
			"error":  err.Error(),
			"key":    key,
			"mobile": mobile,
		}).Error("QuerySMSToken redis HGET failed.")
		return err
	}
	OneDaySendMaxCount := settings.GetInt64("loveauth", "mobile.OneDaySendMaxCount")
	if sendCount >= OneDaySendMaxCount { //一天内发送超限
		return errors.NewCodeString(errors.SendSMSOneDayLimit, "request sms one day limit")
	}

	key = fmt.Sprintf(RedisKeySMSSendTimeList, mobile)
	sendTimeList, err := client.Cmd("LRANGE", key, 0, 2).Array()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
			"key":   key,
		}).Error("QuerySMSToken redis LRANGE failed.")
		return err
	}

	if sendTimeList == nil || len(sendTimeList) == 0 { //最近10分钟没有发生短信记录
		return nil
	}

	now := time.Now().Unix()
	lastSendTime, _ := sendTimeList[0].Int64() //上次发送短信时间
	SendFrequencySeconds := settings.GetInt64("loveauth", "mobile.SendFrequencySeconds")
	if now-lastSendTime < SendFrequencySeconds { //1分钟内发送过
		return errors.NewCodeString(errors.SendSMSFrequently, "request sms request too fast")
	}

	NeedPicCodeTimes := settings.GetInt("loveauth", "mobile.NeedPicCodeTimes")
	if len(sendTimeList) < NeedPicCodeTimes {
		return nil
	}

	NeedPicCodeSecondes := settings.GetInt64("loveauth", "mobile.NeedPicCodeSecondes")
	last3rdSendTime, _ := sendTimeList[NeedPicCodeTimes-1].Int64() //前三条的发送时间
	if now-last3rdSendTime < NeedPicCodeSecondes { //规定时间内发送条数超限，需要图形验证码
		return errors.NewCodeString(errors.SendSMSNeedPicture, "request sms need picture")
	}

	return nil
}

func MonitorRedis() (map[string]string, error) {
	client, err := redisPool.Get()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("MonitorRedis get redis client failed.")
		return nil, err
	}
	defer redisPool.Put(client)

	key := fmt.Sprintf(RedisKeySMSSendCount, time.Now().Format("2006-01-02"))
	allCount, err := client.Cmd("HGETALL", key).Map()
	if err != nil && err != redis.ErrRespNil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
			"key":   key,
		}).Error("MonitorRedis redis HGETALL failed.")
		return nil, err
	}

	return allCount, nil
}
