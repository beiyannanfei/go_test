// Package classification loveauth API.
//
// Auth server for World of Love.
//
// Terms Of Service:
//
// there are no TOS at this moment, use at your own risk we take no responsibility
//
//     Schemes: http, https
//     BasePath: /api/v1
//     Version: v1.0.64 -- develop(0efd925)
//     License: Apache License 2.0 http://www.apache.org/licenses/LICENSE-2.0.html
//     Contact: xy.kong<xy.kong@gmail.com>
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//
// swagger:meta
package server

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/server/auth"
	"github.com/xykong/loveauth/server/bind"
	"github.com/xykong/loveauth/server/callback"
	"github.com/xykong/loveauth/server/coupon"
	"github.com/xykong/loveauth/server/idip"
	"github.com/xykong/loveauth/server/payment_v2"
	"github.com/xykong/loveauth/server/query"
	"github.com/xykong/loveauth/server/send_mail"
	"github.com/xykong/loveauth/server/watch_waring"
	"github.com/xykong/loveauth/server/tss"
	"github.com/xykong/loveauth/server/verify"
	"github.com/xykong/loveauth/server/wjx"
	"github.com/xykong/loveauth/services/dict"
	"github.com/xykong/loveauth/services/logout"
	"github.com/xykong/loveauth/services/payment/alipay"
	"github.com/xykong/loveauth/services/tlog"
	"github.com/xykong/loveauth/settings"
	"github.com/xykong/loveauth/storage"
	"github.com/xykong/loveauth/utils"
	"github.com/xykong/loveauth/utils/graceful"
	"github.com/xykong/loveauth/utils/log"
	"net/http"
	"os"
	"time"
)

var logo = `
        __                               __  __  
       / /___ _   _____     ____ ___  __/ /_/ /_ 
      / / __ \ | / / _ \   / __ '/ / / / __/ __ \
     / / /_/ / |/ /  __/  / /_/ / /_/ / /_/ / / /
    /_/\____/|___/\___/   \__,_/\__,_/\__/_/ /_/

    BasePath: /api/v1
    Version: v1.0.64 -- develop(0efd925)
    Contact: xy.kong<xy.kong@gmail.com>

`

func Run() {

	logrus.Info(logo)

	if nil != log.InitBILogger() {
		return
	}

	log.InitLogger()

	storage.Initialize()

	var setting = settings.Get("loveauth")
	var listen = setting.GetString("gin.listen")
	var mode = setting.GetString("gin.mode")

	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	dict.Load(setting.GetString("dict_path"))

	//router := gin.Default()

	alipay.Init()

	router := gin.New()

	newrelicApp := utils.InitNewRelic()
	if newrelicApp != nil {
		router.Use(utils.NewRelic(newrelicApp))
	}

	// Add a ginrus middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.
	//   - RFC3339 with UTC time format.
	router.Use(utils.Ginrus(logrus.StandardLogger(), time.RFC3339, true))

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	var basePath = "/api/v1"

	callback.Start(router.Group(basePath))

	v1 := router.Group(basePath)
	{
		authGroup := v1.Group("/auth")
		{
			auth.Use(&auth.QQAuth{}, &auth.WechatAuth{}, &auth.GuestAuth{}, &auth.DeviceAuth{},
				&auth.MobileAuth{}, &auth.WeiboAuth{}, &auth.YSDKWechatAuth{}, &auth.YSDKQQAuth{},
				&auth.VivoAuth{}, &auth.BilibiliAuth{}, &auth.HuaweiAuth{}, &auth.QuickAliGamesAuth{},
				&auth.QuickIqiyiAuth{}, &auth.QuickKuaikanAuth{}, &auth.QuickM4399Auth{}, &auth.QuickMeizuAuth{},
				&auth.QuickOppoAuth{}, &auth.QuickYsdkAuth{}, &auth.QuickXiaomiAuth{}, &auth.DouyinAuth{}, &auth.MgtvAuth{})
			auth.Start(authGroup)
		}

		coupon.Start(v1.Group("/coupon"))
		idip.Start(v1.Group("/idip"))
		query.Start(v1.Group("/query"))
		payment_v2.Start(v1.Group("/payment"))
		tss.Start(v1.Group("/tss"))
		verify.Start(v1.Group("/verify"))
		wjx.Start(v1.Group("/wjx"))
		watch_waring.Start(v1.Group("/watch_waring"))

		bindGroup := v1.Group("/bind")
		{
			bind.Use(&bind.BindMobile{}, &bind.BindWechat{}, &bind.BindQQ{}, &bind.BindWeibo{})
			bind.Start(bindGroup)
		}
	}

	// support cors https://github.com/gin-gonic/gin/issues/29
	router.Use(func(c *gin.Context) {
		// Run this on all requests
		// Should be moved to a proper middleware
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, DELETE, PUT, PATCH, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "API-Key, accept, Content-Type, Token")
		c.Writer.Header().Set("Access-Control-Max-Age", "0")
		c.Writer.Header().Set("Content-Length", "0")
		c.Next()
	})

	router.OPTIONS("/*cors", func(c *gin.Context) {
		// Empty 200 response
	})

	if settings.GetBool("loveauth", "docs.swagger") {

		logrus.WithFields(logrus.Fields{
			"path": "docs/swagger.json",
		}).Info("Start swagger service.")

		router.StaticFS("/swagger", http.Dir("scripts/swagger-ui/node_modules/swagger-ui-dist"))
		router.StaticFile("/swagger.json", "docs/swagger.json")
	}

	//c := make(chan os.Signal, 1)
	//signal.Notify(c, os.Interrupt)
	//go func() {
	//	for sig := range c {
	//		// sig is a ^C, handle it
	//
	//		logrus.Warnf("sig is a ^C, handle it: %v", sig)
	//		os.Exit(0)
	//	}
	//}()

	var isRunLogout = setting.GetBool("run_logout")
	if isRunLogout {
		go logout.Start()
	}

	var mailSwitch = settings.GetBool("waring", "mailSwitch")
	if true == mailSwitch {
		go send_mail.Start()
	}

	go tlog.Start()

	//noinspection SpellCheckingInspection
	logrus.WithFields(logrus.Fields{
		"pid":     os.Getpid(),
		"listen":  listen,
		"version": "v1.0.64 -- develop(0efd925)",
		"path":    basePath,
	}).Info("Start loveauth service.")

	go utils.SavePid()

	graceful.ListenAndServe(listen, router)
}
