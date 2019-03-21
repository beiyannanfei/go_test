package storage

import (
	"encoding/json"
	"fmt"
	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/errors"
	"github.com/xykong/loveauth/settings"
	"github.com/xykong/loveauth/storage/model"
	"strings"
	"time"
)

//var client *redis.Client
var redisPool *pool.Pool

var RedisKeyTemplateGlobalIdAccessToken = "auth:id_token:%v"
var RedisKeyTemplateAccessToken = "auth:token_id:%v"
var RedisKeyTemplateRefreshToken = "auth:refresh_token:%v"
var RedisKeyTemplateActiveAccount = "auth:active_account:%v:%v"
var RedisKeyChannelList = "auth:channel_list"

func InitRedis() {

	var setting = settings.Get("loveauth")

	var host = setting.GetString("redis.host")
	var port = setting.GetInt64("redis.port")
	var size = setting.GetInt("redis.pool")
	var password = setting.GetString("redis.password")
	var database = setting.GetInt("redis.database")

	addr := fmt.Sprintf("%s:%d", host, port)

	logrus.WithFields(logrus.Fields{
		"addr": addr,
	}).Info("Connection to Redis")

	var err error
	//client, err = redis.Dial("tcp", addr)
	//if err != nil {
	//	logrus.WithFields(logrus.Fields{
	//		"error": err.Error(),
	//	}).Fatal("Connection redis failed.")
	//}

	df := func(network, addr string) (*redis.Client, error) {
		client, err := redis.Dial(network, addr)
		if err != nil {
			return nil, err
		}

		if password != "" {
			if err = client.Cmd("AUTH", password).Err; err != nil {
				client.Close()
				return nil, err
			}
		}

		if database > 0 {
			if err = client.Cmd("SELECT", database).Err; err != nil {
				client.Close()
				return nil, err
			}
		}
		return client, nil
	}

	redisPool, err = pool.NewCustom("tcp", addr, size, df)

	//redisPool, err = pool.New("tcp", addr, size)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("Connection redis failed.")
	}
}

type Command struct {
	Cmd  string
	Args []interface{}
}

type Pipe struct {
	commands []Command
}

func (pipe *Pipe) Append(cmd string, args ...interface{}) {
	pipe.commands = append(pipe.commands, Command{Cmd: cmd, Args: args})
}

func (pipe *Pipe) Run(name string) []error {

	client, err := redisPool.Get()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("WriteToken token failed.")
	}
	defer redisPool.Put(client)

	for _, value := range pipe.commands {
		client.PipeAppend(value.Cmd, value.Args...)
	}

	var errors []error
	for _, value := range pipe.commands {

		err = client.PipeResp().Err
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
				"cmd":   value.Cmd,
				"args":  value.Args,
			}).Errorf("Pipe Commands %v failed.", name)
		}

		errors = append(errors, err)
	}

	return errors
}

func execRedisCmd(tryTime int, client *redis.Client, cmd string, args ...interface{}) ([]byte, error) {
	savedValue, err := client.Cmd(cmd, args...).Bytes()
	if err != nil {
		if err != redis.ErrRespNil { //非空错误
			logrus.WithFields(logrus.Fields{
				"tryTime": tryTime,
				"err":     err,
				"cmd":     cmd,
				"args":    args,
			}).Error("execRedisCmd failed.")

			return nil, err
		}

		logrus.WithFields(logrus.Fields{
			"tryTime": tryTime,
			"err":     err,
			"cmd":     cmd,
			"args":    args,
		}).Info("execRedisCmd response is nil.")

		//结果为空的错误
		if tryTime > 2 { //重试次数超过3次
			return nil, err
		}

		tryTime++
		time.Sleep(time.Millisecond * 10 * time.Duration(tryTime)) //延时重试(根据重试次数不同，退避时间延长)
		return execRedisCmd(tryTime, client, cmd, args...)
	}

	return savedValue, nil
}

func cmdGetObject(name string, client *redis.Client, v interface{}, cmd string, args ...interface{}) error {
	savedValue, err := execRedisCmd(0, client, cmd, args...)
	if err != nil {
		if err != redis.ErrRespNil {
			logrus.WithFields(logrus.Fields{
				"error": err,
				"cmd":   cmd,
				"args":  args,
				"name":  name,
			}).Errorf("%v cmdGet struct failed.", name)
		} else {
			logrus.WithFields(logrus.Fields{
				"error": err,
				"cmd":   cmd,
				"args":  args,
				"name":  name,
			}).Infof("%v cmdGet struct failed.", name)
		}

		return err
	}

	err = json.Unmarshal(savedValue, v)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":      err.Error(),
			"savedValue": string(savedValue),
		}).Errorf("%v Unmarshal saved value failed.", name)
	}

	return err
}

func WriteToken(profile *model.Profile) (lastProfile *model.Profile) {

	client, err := redisPool.Get()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("WriteToken get redis client failed.")
		return
	}
	defer redisPool.Put(client)

	pipe := Pipe{}

	////////////////////////////////////////////////////////////////////////////////
	// process for multiple login.
	var savedProfile model.Profile
	key := fmt.Sprintf(RedisKeyTemplateGlobalIdAccessToken, profile.GlobalId)
	err = cmdGetObject("WriteToken", client, &savedProfile, "GET", key)
	if err == nil {

		last := QueryActiveAccount(&model.TokenRecord{
			GlobalId:  savedProfile.GlobalId,
			Vendor:    savedProfile.Vendor,
			Platform:  savedProfile.Platform,
			Timestamp: 0,
		})

		if last > 0 {
			key := fmt.Sprintf(RedisKeyTemplateActiveAccount, savedProfile.Auth.Channel, savedProfile.Platform)
			pipe.Append("ZREM", key, savedProfile.GlobalId)

			lastProfile = &savedProfile
		}

		key = fmt.Sprintf(RedisKeyTemplateAccessToken, savedProfile.Token)
		pipe.Append("DEL", key)

		key = fmt.Sprintf(RedisKeyTemplateRefreshToken, savedProfile.RefreshToken)
		pipe.Append("DEL", key)
	}
	////////////////////////////////////////////////////////////////////////////////

	// info for query token to id, auth verify.
	var tokenRecord = model.TokenRecord{
		GlobalId:  profile.GlobalId,
		Vendor:    profile.Vendor,
		Platform:  profile.Platform,
		Timestamp: profile.Timestamp,
	}
	jsonString, _ := json.Marshal(tokenRecord)
	key = fmt.Sprintf(RedisKeyTemplateAccessToken, profile.Token)
	pipe.Append("SETEX", key, profile.ExpirationSeconds, jsonString)

	// info for refresh token logic.
	var refreshTokenRecord = model.RefreshTokenRecord{
		GlobalId:          profile.GlobalId,
		Token:             profile.Token,
		RefreshToken:      profile.RefreshToken,
		Vendor:            profile.Vendor,
		Platform:          profile.Platform,
		Timestamp:         profile.Timestamp,
		ExpirationSeconds: 0,
	}
	jsonString, _ = json.Marshal(refreshTokenRecord)
	key = fmt.Sprintf(RedisKeyTemplateRefreshToken, profile.RefreshToken)
	pipe.Append("SETEX", key, profile.ExpirationSeconds, jsonString)

	// info for query user info by globalId.
	jsonString, _ = json.Marshal(profile)
	key = fmt.Sprintf(RedisKeyTemplateGlobalIdAccessToken, profile.GlobalId)
	pipe.Append("SETEX", key, profile.ExpirationSeconds, jsonString)

	// info for detect user logout by expiration
	key = fmt.Sprintf(RedisKeyTemplateActiveAccount, profile.Auth.Channel, profile.Platform)
	pipe.Append("ZADD", key, profile.Timestamp, profile.GlobalId)

	pipe.Run("WriteToken")

	return
}

func QueryProfile(globalId int64) (*model.Profile, error) {

	client, err := redisPool.Get()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("QueryProfile get redis client failed.")
		return nil, err
	}
	defer redisPool.Put(client)

	var record model.Profile
	key := fmt.Sprintf(RedisKeyTemplateGlobalIdAccessToken, globalId)
	err = cmdGetObject("QueryAccessToken", client, &record, "GET", key)

	return &record, err
}

func WriteProfile(profile *model.Profile) {

	client, err := redisPool.Get()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("WriteProfile get redis client failed.")
		return
	}
	defer redisPool.Put(client)

	jsonString, _ := json.Marshal(profile)
	key := fmt.Sprintf(RedisKeyTemplateGlobalIdAccessToken, profile.GlobalId)
	err = client.Cmd("SETEX", key, profile.ExpirationSeconds, jsonString).Err
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err,
			"profile": profile,
		}).Error("WriteProfile failed.")
	}
}

func QueryAccessToken(token string) (*model.TokenRecord, error) {

	client, err := redisPool.Get()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("QueryAccessToken failed.")
	}
	defer redisPool.Put(client)

	var tokenRecord model.TokenRecord
	key := fmt.Sprintf(RedisKeyTemplateAccessToken, token)
	err = cmdGetObject("QueryAccessToken", client, &tokenRecord, "GET", key)

	return &tokenRecord, err
}

func QueryRefreshToken(token string) (*model.RefreshTokenRecord, error) {

	client, err := redisPool.Get()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("QueryRefreshToken failed.")
	}
	defer redisPool.Put(client)

	var record model.RefreshTokenRecord
	key := fmt.Sprintf(RedisKeyTemplateRefreshToken, token)
	err = cmdGetObject("QueryRefreshToken", client, &record, "GET", key)

	return &record, err
}

func RefreshToken(record *model.RefreshTokenRecord) error {

	client, err := redisPool.Get()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("WriteToken token failed.")

		return err
	}
	defer redisPool.Put(client)

	key := fmt.Sprintf(RedisKeyTemplateAccessToken, record.Token)
	err = client.Cmd("EXPIRE", key, record.ExpirationSeconds).Err
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Set access token expire failed.")

		return err
	}

	key = fmt.Sprintf(RedisKeyTemplateRefreshToken, record.RefreshToken)
	err = client.Cmd("EXPIRE", key, record.ExpirationSeconds).Err
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Set refresh token expire failed.")

		return err
	}

	key = fmt.Sprintf(RedisKeyTemplateGlobalIdAccessToken, record.GlobalId)
	err = client.Cmd("EXPIRE", key, record.ExpirationSeconds).Err
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Set id to refresh token expire failed.")

		return err
	}

	return nil
}

func QueryActiveAccount(record *model.TokenRecord) int64 {

	client, err := redisPool.Get()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("QueryActiveAccount redis failed.")
	}
	defer redisPool.Put(client)

	profile, err := QueryProfile(record.GlobalId)
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err":      err.Error(),
			"globalId": record.GlobalId,
		}).Error("QueryActiveAccount QueryProfile")
	}

	key := fmt.Sprintf(RedisKeyTemplateActiveAccount, profile.Auth.Channel, record.Platform)
	ret, _ := client.Cmd("ZSCORE", key, record.GlobalId).Int64()

	return ret
}

func TouchAccount(tokenRecord *model.TokenRecord) {

	client, err := redisPool.Get()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("TouchAccount redis failed.")
	}
	defer redisPool.Put(client)

	pipe := Pipe{}

	profile, err := QueryProfile(tokenRecord.GlobalId)
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"error":    err.Error(),
			"globalId": tokenRecord.GlobalId,
		}).Error("TouchAccount get profile failed")
	}

	timestamp := uint(time.Now().Unix())
	key := fmt.Sprintf(RedisKeyTemplateActiveAccount, profile.Auth.Channel, tokenRecord.Platform)
	//key := fmt.Sprintf(RedisKeyTemplateActiveAccount, tokenRecord.Vendor, tokenRecord.Platform)
	pipe.Append("ZADD", key, timestamp, tokenRecord.GlobalId)
	pipe.Append("SADD", RedisKeyChannelList, profile.Auth.Channel)

	pipe.Run("TouchAccount")
}

func QueryExpiredLogout(channel string, platform model.Platform, start, stop int64) map[string]string {

	client, err := redisPool.Get()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("QueryExpiredLogout get redis client failed.")
	}
	defer redisPool.Put(client)

	//profile, err := QueryProfile(globalId)
	//if nil != err {
	//	logrus.WithFields(logrus.Fields{
	//		"err": err.Error(),
	//		"globalId": globalId,
	//	}).Error("PopExpiredLogout QueryProfile")
	//}

	key := fmt.Sprintf(RedisKeyTemplateActiveAccount, channel, platform)
	users, err := client.Cmd("ZRANGEBYSCORE", key, start, stop, "WITHSCORES").Map()
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"key":   key,
			"error": err.Error(),
		}).Error("query refresh token failed.")
	}

	//logrus.WithFields(logrus.Fields{
	//	"key":  key,
	//	"args": []interface{}{"ZRANGEBYSCORE", key, start, stop, "WITHSCORES"},
	//}).Warnf("users: %v", users)

	return users
}

func PopExpiredLogout(channel string, platform model.Platform, globalId int64, timestamp int64) (*model.Profile, error) {

	client, err := redisPool.Get()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("QueryExpiredLogout get redis client failed.")
	}
	defer redisPool.Put(client)

	key := fmt.Sprintf(RedisKeyTemplateActiveAccount, channel, platform)
	savedTime, err := client.Cmd("ZSCORE", key, globalId).Int64()
	if err != nil {
		if err != redis.ErrRespNil {
			logrus.WithFields(logrus.Fields{
				"key":       key,
				"globalId":  globalId,
				"timestamp": timestamp,
				"error":     err.Error(),
			}).Error("PopExpiredLogout ZSCORE failed.")
		} else {
			logrus.WithFields(logrus.Fields{
				"key":       key,
				"globalId":  globalId,
				"timestamp": timestamp,
				"error":     err.Error(),
			}).Info("PopExpiredLogout ZSCORE null.")
		}

		return nil, err
	}

	if timestamp != 0 && savedTime != 0 && timestamp != savedTime {
		logrus.WithFields(logrus.Fields{
			"key":       key,
			"globalId":  globalId,
			"timestamp": timestamp,
			"savedTime": savedTime,
		}).Error("PopExpiredLogout ZSCORE failed.")
	}

	count, err := client.Cmd("ZREM", key, globalId).Int64()
	if err != nil {
		if err != redis.ErrRespNil {
			logrus.WithFields(logrus.Fields{
				"key":       key,
				"globalId":  globalId,
				"timestamp": timestamp,
				"error":     err.Error(),
			}).Error("PopExpiredLogout ZREM failed.")
		} else {
			logrus.WithFields(logrus.Fields{
				"key":       key,
				"globalId":  globalId,
				"timestamp": timestamp,
				"error":     err.Error(),
			}).Info("PopExpiredLogout response null.")
		}

		return nil, err
	}

	if count != 1 {

		logrus.WithFields(logrus.Fields{
			"key":       key,
			"globalId":  globalId,
			"timestamp": timestamp,
			"count":     count,
		}).Warn("PopExpiredLogout count 0.")

		return nil, errors.New("PopExpiredLogout ZREM response 0")
	}

	logrus.WithFields(logrus.Fields{
		"channel":   channel,
		"platform":  platform,
		"globalId":  globalId,
		"timestamp": timestamp,
	}).Info("PopExpiredLogout success.")

	var record model.Profile
	key = fmt.Sprintf(RedisKeyTemplateGlobalIdAccessToken, globalId)
	err = cmdGetObject("PopExpiredLogout", client, &record, "GET", key)

	return &record, err
}

func CountOnlineToken(channel string, platform model.Platform) int {

	client, err := redisPool.Get()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("QueryExpiredLogout get redis client failed.")
	}
	defer redisPool.Put(client)

	var counts int
	var channels []string

	if channel != model.ChannelAll {

		channels = append(channels, channel)
	} else {

		channels = GetChannelList()
	}

	for _, queryChannel := range channels {

		key := fmt.Sprintf(RedisKeyTemplateActiveAccount, queryChannel, platform)
		count, err := client.Cmd("ZCOUNT", key, "-inf", "+inf").Int()

		if err != nil {

			logrus.WithFields(logrus.Fields{
				"error": err,
				"key":   key,
			}).Error("CountOnlineToken failed.")
		}

		counts = counts + count
	}

	return counts
}

func RemoveToken(globalId int64) {

	client, err := redisPool.Get()
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("BanOnLineUser get redis client failed.")

		return
	}
	defer redisPool.Put(client)

	var record model.Profile

	key := fmt.Sprintf(RedisKeyTemplateGlobalIdAccessToken, globalId)
	err = cmdGetObject("QueryAccessToken", client, &record, "GET", key)
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"key":   key,
			"error": err.Error(),
		}).Error("query key:id_token failed.")

		return
	}

	pipe := Pipe{}

	key = fmt.Sprintf(RedisKeyTemplateAccessToken, record.Token)
	pipe.Append("DEL", key)

	key = fmt.Sprintf(RedisKeyTemplateRefreshToken, record.RefreshToken)
	pipe.Append("DEL", key)

	pipe.Run("RemoveToken")
}

//是否打印过在线人数
func HasOnlineLog() bool {

	client, err := redisPool.Get()
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("HasOnlineLog get redis client failed.")

		return false
	}
	defer redisPool.Put(client)

	key := fmt.Sprintf("online:log:%s", time.Now().Format("200601021504"))
	resp := client.Cmd("INCR", key)
	if resp.Err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("incr online:log failed.")

		return false
	}

	data, _ := resp.Int()
	if data > 1 {

		return true
	}

	err = client.Cmd("EXPIRE", key, 180).Err
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Set online:log expire failed.")

		return false
	}

	return false
}

// 获取玩家使用的渠道列表
func GetChannelList() []string {
	client, err := redisPool.Get()
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("BanOnLineUser get redis client failed.")

		return []string{}
	}
	defer redisPool.Put(client)

	result, err := client.Cmd("SMEMBERS", RedisKeyChannelList).List()
	//logrus.WithFields(logrus.Fields{
	//	"result": result,
	//	"error": err,
	//}).Error("GetChannelList")
	return result
}

func SetSession(token, account, url string) error {

	client, err := redisPool.Get()
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("SetSession get redis client failed.")

		return nil
	}
	defer redisPool.Put(client)

	return client.Cmd("SETEX", "user:session:"+token, 600, account+"_-_"+url).Err
}

func GetSession(token string) (string, string) {

	if token == "" {

		return "", ""
	}

	client, err := redisPool.Get()
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("SetSession get redis client failed.")

		return "", ""
	}
	defer redisPool.Put(client)

	sessionInfo, err := client.Cmd("get", "user:session:"+token).Str()
	if sessionInfo == "" || err != nil {

		return "", ""
	}

	client.Cmd("SETEX", "user:session:"+token, 600, sessionInfo)

	infos := strings.Split(sessionInfo, "_-_")
	if len(infos) != 2 {

		return "", ""
	}
	return infos[0], infos[1]
}
