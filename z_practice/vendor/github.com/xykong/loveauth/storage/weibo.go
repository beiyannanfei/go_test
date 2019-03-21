package storage

import (
	"github.com/sirupsen/logrus"
	"fmt"
)

var RedisKeyWeiboCode = "auth:weibo_code:%v"

func SaveWeiboCode(code string) error {
	client, err := redisPool.Get()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("SaveWeiboCode get redis client failed.")
		return err
	}
	defer redisPool.Put(client)

	key := fmt.Sprintf(RedisKeyWeiboCode, code)
	err = client.Cmd("SETEX", key, 10*60, code).Err
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
			"key":   key,
		}).Error("SaveWeiboCode failed.")
		return err
	}

	return nil
}

func DeleteWeiboCode(code string) {
	client, err := redisPool.Get()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("DeleteWeiboCode get redis client failed.")
	}
	defer redisPool.Put(client)

	key := fmt.Sprintf(RedisKeyWeiboCode, code)
	err = client.Cmd("DEL", key).Err
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
			"key":   key,
		}).Error("DeleteWeiboCode failed.")
	}
}

func QueryWeiboCode(code string) (string, error) {
	client, err := redisPool.Get()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("QueryWeiboCode get redis client failed.")
		return "", err
	}
	defer redisPool.Put(client)

	key := fmt.Sprintf(RedisKeyWeiboCode, code)
	value, err := client.Cmd("GET", key).Str()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
			"key":   key,
		}).Error("QueryWeiboCode failed.")
		return "", err
	}

	return value, nil
}
