package logout

import (
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/server/auth"
	"github.com/xykong/loveauth/settings"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/storage/model"
	"strconv"
	"time"
)

func Start() {

	var logoutExpirationSeconds = settings.GetInt64("loveauth", "auth.LogoutExpirationSeconds")
	var logoutCheckingOffsetSeconds = settings.GetInt64("loveauth", "auth.LogoutCheckingOffsetSeconds")

	var offsetSeconds = logoutExpirationSeconds + logoutCheckingOffsetSeconds

	logrus.WithFields(logrus.Fields{
		"logoutExpirationSeconds":     logoutExpirationSeconds,
		"logoutCheckingOffsetSeconds": logoutCheckingOffsetSeconds,
	}).Info("logout services start.")

	channelList := storage.GetChannelList()
	for _, channel := range channelList {

		for _, platform := range model.Platforms {

			go processLogout(channel, platform, offsetSeconds)
		}
	}
}

func processLogout(channel string, platform model.Platform, offset int64) {
	logrus.WithFields(logrus.Fields{
		"channel":  channel,
		"platform": platform,
		"offset":   offset,
	}).Info("processLogout")

	for {

		logoutFromQueue(channel, platform, offset)

		time.Sleep(time.Second * 10)
	}
}

func logoutFromQueue(channel string, platform model.Platform, offset int64) {

	// fetch logout array from redis.

	now := time.Now().Unix()

	users := storage.QueryExpiredLogout(channel, platform, 0, now-offset)

	for k, v := range users {

		globalId, err := strconv.ParseInt(k, 10, 0)
		if err != nil {
			continue
		}
		timestamp, err := strconv.ParseInt(v, 10, 0)
		if err != nil {
			continue
		}

		profile, err := storage.PopExpiredLogout(channel, platform, globalId, timestamp)
		if profile != nil && err == nil {
			auth.LogPlayerLogout(profile, timestamp)

			account := storage.QueryAccount(globalId)
			if account != nil {

				account.LogoutTime = timestamp
				account.AccumLoginTime = account.AccumLoginTime + timestamp - profile.Timestamp
				storage.Save(storage.AuthDatabase(), account)
			}
		}
	}
}
