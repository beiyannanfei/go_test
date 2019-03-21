package utils

import (
	"github.com/newrelic/go-agent"
	"github.com/xykong/loveauth/settings"
	"github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

var app newrelic.Application

func InitNewRelic() newrelic.Application {

	if app != nil {
		return app
	}

	enable := settings.GetBool("newrelic", "enable")

	if !enable {
		return nil
	}

	appName := settings.GetString("newrelic", "app_name")
	license := settings.GetString("newrelic", "license")

	if appName == "" || license == "" {
		return nil
	}

	logrus.WithFields(logrus.Fields{
		"app_name": appName,
		"license":  license,
	}).Info("New Relic start.")

	config := newrelic.NewConfig(appName, license)
	var err error
	app, err = newrelic.NewApplication(config)

	if err != nil {

		logrus.WithFields(logrus.Fields{
			"app_name": appName,
			"license":  license,
			"error":    err,
		}).Error("New Relic NewApplication failed.")

		return nil
	}

	return app
}

func NewRelic(app newrelic.Application) gin.HandlerFunc {
	return func(c *gin.Context) {

		//logrus.WithFields(logrus.Fields{
		//	"Transaction": c.Request.RequestURI,
		//}).Info("New Relic StartTransaction.")

		txn := app.StartTransaction(c.Request.RequestURI, c.Writer, c.Request)
		defer txn.End()

		c.Next()
	}
}

func GetNewRelic() newrelic.Application {
	return app
}
