package storage

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/storage/model"
	"time"
)

const PreGenSequence = 0 // 预下单
const QuerySequence = 1  // 查询发货

var RedisKeyWatchItem = "auth:watch_user_item:%v:%v_%v"        // date:itemId_opsType  subKey = globalId
var RedisKeyWatchOrders = "auth:watch_orders:%v:%v:%v"    // date:channel:orderType
var RedisKeyWatchOrdersFlag = "auth:waring_order_flag:%v" // date
var RedisKeyWatchDevice = "auth:watch_device:%v:%v"       // date:deviceId
var RedisKeyDeviceList = "auth:watch_device_count:%v"     // date
var RedisKeyWatchAccount = "auth:watch_account:%v:%v"     // date:globalId
var RedisKeyAccountList = "auth:watch_account_count:%v"   // date
var RedisKeyUserChannel = "auth:watch_channel_uid:%v"     // date
var RedisKeyWatchPayment = "auth:watch_payment:%v"        // date

// 更新玩家获得、消耗物品状态
func UpdateItemState(globalId int64, channel string, itemId string, opsType int, itemNum int) {
	client, err := redisPool.Get()
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("item watch get redis client failed")
		return
	}
	defer redisPool.Put(client)

	now := time.Now()
	nowStr := now.Format("2006-01-02")

	// 存储物品获得消耗
	redisKey := fmt.Sprintf(RedisKeyWatchItem, nowStr, itemId, opsType)
	err = client.Cmd("HINCRBY", redisKey, globalId, itemNum).Err
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("error while HINCRBY, key = ", redisKey, " ", globalId)
		return
	}

	// 存储玩家渠道信息
	channelKey := fmt.Sprintf(RedisKeyUserChannel, nowStr)
	err = client.Cmd("HSET", channelKey, globalId, channel).Err
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("error while store user channel, key = ", channelKey)
		return
	}

	// 设置过期时间
	expireAt := now.Add(time.Hour * 24).Unix()
	err = client.Cmd("EXPIREAT", redisKey, expireAt).Err
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("error while set redis key expire at ", expireAt, " key = ", redisKey)
	}
	err = client.Cmd("EXPIREAT", channelKey, expireAt).Err
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("error while set redis key expire at ", expireAt, " key = ", channelKey)
	}
}

// 拼装物品监控邮件内容
func GetItemWatchInfo(itemId string, opsType, threshold int) string {
	body := ""
	client, err := redisPool.Get()
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("item watch get redis client failed")
		return body
	}
	defer redisPool.Put(client)

	now := time.Now()
	nowStr := now.Format("2006-01-02")
	timeStr := now.Format("2006-01-02 03:04:05")

	channelKey := fmt.Sprintf(RedisKeyUserChannel, nowStr)

	redisKey := fmt.Sprintf(RedisKeyWatchItem, nowStr, itemId, opsType)
	globalIdList, err := client.Cmd("HKEYS", redisKey).List()
	for _, globalId := range globalIdList {
		logrus.Info(globalId)
		logrus.WithFields(logrus.Fields{
			"key":    redisKey,
			"subKey": globalId,
		}).Info("item watch")
		itemNum, err := client.Cmd("HGET", redisKey, globalId).Int()
		if nil != err {
			logrus.WithFields(logrus.Fields{
				"err": err.Error(),
			}).Error("error while get user item count, key = ", redisKey, " ", globalId)
			continue
		}
		if itemNum >= threshold {
			channel, err := client.Cmd("HGET", channelKey, globalId).Str()
			if nil != err {
				logrus.WithFields(logrus.Fields{
					"err":  err.Error(),
					"func": "GetItemWatchInfo",
				}).Error("error while get channel, key = ", channelKey, " ", globalId)
				channel = ""
			}
			row := fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%d</td><td>%d</td>", globalId, timeStr, channel, itemId, opsType, itemNum)
			body += row
		}
	}
	return body
}

// 更新账号、设备状态
func UpdateLoginState(globalId int64, channel string, deviceId string) {
	client, err := redisPool.Get()
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("order watch get redis client failed")
		return
	}
	defer redisPool.Put(client)

	now := time.Now()
	nowStr := now.Format("2006-01-02")
	// 更新设备监控
	deviceKey := fmt.Sprintf(RedisKeyWatchDevice, nowStr, deviceId)
	err = client.Cmd("SADD", deviceKey, globalId).Err
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("error while SADD globalId, key = ", deviceKey)
		return
	}
	deviceCount, err := client.Cmd("SCARD", deviceKey).Int()
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("error while scard, key = ", deviceKey)
		return
	}

	deviceCountKey := fmt.Sprintf(RedisKeyDeviceList, nowStr)
	err = client.Cmd("HSET", deviceCountKey, deviceId, deviceCount).Err
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("error while hset, key = ", deviceCountKey)
		return
	}

	// 更新账号监控
	accountKey := fmt.Sprintf(RedisKeyWatchAccount, nowStr, globalId )
	err = client.Cmd("SADD", accountKey, deviceId).Err
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err":  err.Error(),
			"func": "UpdateLoginState",
		}).Error("error while SADD globalId, key = ", accountKey)
		return
	}
	accountCount, err := client.Cmd("SCARD", accountKey).Int()
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("error while scard, key = ", accountKey)
		return
	}
	accountCountKey := fmt.Sprintf(RedisKeyAccountList, nowStr)
	err = client.Cmd("HSET", accountCountKey, globalId, accountCount).Err
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("error while hset, key = ", accountCountKey)
		return
	}

	// 更新玩家渠道信息
	channelKey := fmt.Sprintf(RedisKeyUserChannel, nowStr)
	err = client.Cmd("HSET", channelKey, globalId, channel).Err
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("error while set user channel, key = ", channelKey)
		return
	}

	// 设置过期时间
	expireAt := now.Add(time.Hour * 24).Unix()
	err = client.Cmd("EXPIREAT", deviceKey, expireAt).Err
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("error while set redis key expire at ", expireAt, " key = ", deviceKey)
	}
	err = client.Cmd("EXPIREAT", accountKey, expireAt).Err
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("error while set redis key expire at ", expireAt, " key = ", accountKey)
	}
	err = client.Cmd("EXPIREAT", accountKey, expireAt).Err
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("error while set redis key expire at ", expireAt, " key = ", accountKey)
	}
}

func GetAccountWatchInfo(threshold int) string {
	client, err := redisPool.Get()
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("order watch get redis client failed")
		return ""
	}
	defer redisPool.Put(client)

	now := time.Now()
	nowStr := now.Format("2006-01-02")
	timeStr := now.Format("2006-01-02 03:04:05")
	accountListKey := fmt.Sprintf(RedisKeyAccountList, nowStr)
	accountList, err := client.Cmd("HKEYS", accountListKey).List()
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("error while get account watch key list, key = ", accountList)
		return ""
	}

	body := ""
	for _, account := range accountList {
		count, err := client.Cmd("HGET", accountListKey, account).Int()
		if nil != err {
			logrus.WithFields(logrus.Fields{
				"err": err.Error(),
			}).Error("error while hget, key = ", accountList, " ", account)
			continue
		}
		if count >= threshold {
			channelKey := fmt.Sprintf(RedisKeyUserChannel, nowStr)
			channel, err := client.Cmd("HGET", channelKey, account).Str()
			if nil != err {
				logrus.WithFields(logrus.Fields{
					"err":  err.Error(),
					"func": "GetAccountWatchInfo",
				}).Error("error while get user channel, key = ", channelKey, " ", account)
				channel = ""
			}
			body = body + fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%d</td></tr>", account, timeStr, channel, count)
		}
	}

	return body
}

func GetDeviceWatchInfo(threshold int) string {
	client, err := redisPool.Get()
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("order watch get redis client failed")
		return ""
	}
	defer redisPool.Put(client)

	now := time.Now()
	nowStr := now.Format("2006-01-02")
	timeStr := now.Format("2006-01-02 03:04:05")
	deviceListKey := fmt.Sprintf(RedisKeyDeviceList, nowStr)
	deviceList, err := client.Cmd("HKEYS", deviceListKey).List()
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("error while get device watch key list, key = ", deviceListKey)
		return ""
	}

	body := ""
	for _, deviceId := range deviceList {
		count, err := client.Cmd("HGET", deviceListKey, deviceId).Int()
		if nil != err {
			logrus.WithFields(logrus.Fields{
				"err": err.Error(),
			}).Error("error while hget, key = ", deviceListKey, " ", deviceId)
			continue
		}
		if count >= threshold {
			body = body + fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%d</td></tr>", deviceId, timeStr, count)
		}
	}

	return body
}

// 获取对应状态订单数量
func GetOrderCount(orderType int, channel model.Vendor) int {
	client, err := redisPool.Get()
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("order watch get redis client failed")
		return 0
	}
	defer redisPool.Put(client)

	now := time.Now().Format("2006-01-02")
	redisKey := fmt.Sprintf(RedisKeyWatchOrders, now, channel, orderType)
	// 判断key是否存在
	exist, err := client.Cmd("EXISTS", redisKey).Int()
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("error while check redis key exist, key = ", redisKey)
		return 0
	}
	if 1 != exist {
		return 0
	}

	// 获取订单数量
	count, err := client.Cmd("SCARD", redisKey).Int()
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("error while get order count, key = ", redisKey, " orderType = ", orderType)
		return 0
	}
	return count
}

// 返回值标志今日该渠道是否已经发送过报警邮件
func UpdateOrdersState(orderType int, sequence string, channel model.Vendor) bool {
	client, err := redisPool.Get()
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("order watch get redis client failed")
		return false
	}
	defer redisPool.Put(client)

	now := time.Now()
	nowStr := now.Format("2006-01-02")
	// 更新订单
	redisKey := fmt.Sprintf(RedisKeyWatchOrders, nowStr, channel, orderType)
	err = client.Cmd("SADD", redisKey, sequence).Err
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("error while add order info, key = ", redisKey)
		return false
	}
	expireAt := now.Add(time.Hour * 24).Unix()
	client.Cmd("EXPIREAT", redisKey, expireAt)

	// 更新渠道订单报警标志
	flagKey := fmt.Sprintf(RedisKeyWatchOrdersFlag, nowStr)
	state := client.Cmd("HGET", flagKey, channel)
	if nil == state {
		client.Cmd("HSET", flagKey, channel, 0)
		client.Cmd("EXPIREAT", flagKey, expireAt)
		return true
	}
	f, _ := state.Int()
	if 1 != f {
		return true
	}
	return false
}

func SetOrderWaringDone(vendor model.Vendor) {
	client, err := redisPool.Get()
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("order watch get redis client failed")
		return
	}
	defer redisPool.Put(client)

	now := time.Now()
	// 更新渠道订单报警状态
	redisKey := fmt.Sprintf(RedisKeyWatchOrdersFlag, now.Format("2006-01-02"))
	err = client.Cmd("HSET", redisKey, vendor, 1).Err
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("error while set order alarm flag, key = ", redisKey)
		return
	}
}

// 获取订单渠道列表
func GetOrderChannelList() []string {
	client, err := redisPool.Get()
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("order watch get redis client failed")
		return []string{}
	}
	defer redisPool.Put(client)

	now := time.Now()
	redisKey := fmt.Sprintf(RedisKeyWatchOrdersFlag, now.Format("2006-01-02"))
	channelList, err := client.Cmd("HKEYS", redisKey).List()
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("error while get order key list, key = ", redisKey)
		return channelList
	}
	return channelList
}

func UpdatePaymentState(globalId int64, channel model.Vendor, moneyCount int) {
	client, err := redisPool.Get()
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("item watch get redis client failed")
		return
	}
	defer redisPool.Put(client)

	now := time.Now()
	nowStr := now.Format("2006-01-02")

	redisKey := fmt.Sprintf(RedisKeyWatchPayment, nowStr)
	err = client.Cmd("HINCRBY", redisKey, globalId, moneyCount).Err
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("error while HINCRBY, key = ", redisKey, " ", globalId)
		return
	}

	channelKey := fmt.Sprintf(RedisKeyUserChannel, nowStr)
	err = client.Cmd("HSET", channelKey, globalId, channel).Err
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("error while store user channel, key = ", channelKey)
		return
	}

	// 设置过期时间
	expireAt := now.Add(time.Hour * 24).Unix()
	err = client.Cmd("EXPIREAT", redisKey, expireAt).Err
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("error while set redis key expire at ", expireAt, " key = ", redisKey)
	}
	err = client.Cmd("EXPIREAT", channelKey, expireAt).Err
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("error while set redis key expire at ", expireAt, " key = ", channelKey)
	}
}

func GetPaymentInfo(threshold float64) string {
	client, err := redisPool.Get()
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("item watch get redis client failed")
		return ""
	}
	defer redisPool.Put(client)

	now := time.Now()
	nowStr := now.Format("2006-01-02")
	timeStr := now.Format("2006-01-02 03:04:05")
	redisKey := fmt.Sprintf(RedisKeyWatchPayment, nowStr)
	globalList, err := client.Cmd("HKEYS", redisKey).List()
	if nil != err {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error("error while hkeys, key = ", redisKey)
		return ""
	}

	body := ""
	for _, globalId := range globalList {
		moneyCount, err := client.Cmd("HGET", redisKey, globalId).Float64()
		if nil != err {
			logrus.WithFields(logrus.Fields{
				"err": err.Error(),
			}).Error("error while hget, key = ", redisKey, " ", globalId)
			continue
		}
		if float64(moneyCount) / 100.0 >= threshold {
			channelKey := fmt.Sprintf(RedisKeyUserChannel, nowStr)
			channel, err := client.Cmd("HGET", channelKey, globalId).Str()
			if nil != err {
				logrus.WithFields(logrus.Fields{
					"err": err.Error(),
				}).Error("error while hget, key = ", channelKey, " ", globalId)
				channel = ""
			}
			body = body + fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%.2f</td></tr>", globalId, timeStr, channel, moneyCount)
		}
	}
	return body
}
