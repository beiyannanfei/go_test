package auth

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/storage"
	"math/rand"
	"net/http"
	"time"
)

type GmLoginRequest struct {
	Account  string `json:"account" form:"account"`
	Password string `json:"password" form:"password"`
	Channel  string `json:"channel" form:"channel"`
}

type GmLoginResponse struct {
	Msg   string `json:"msg"`
	Token string `json:"token"`
}

func gm_login(c *gin.Context) {

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "Content-Type,Access-Token")
	c.Header("Access-Control-Expose-Headers", "*")
	resp := &GmLoginResponse{Msg: "success"}

	var req GmLoginRequest
	err := c.Bind(&req)
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"request": req,
			"error":   err,
		}).Error("gm_login failed.")

		resp.Msg = "fail"
		c.JSON(http.StatusOK, resp)

		return
	}

	gmAccount := storage.QueryGmAccountByAccount(req.Account)
	if gmAccount == nil {

		logrus.WithFields(logrus.Fields{
			"request": req,
			"error":   err,
		}).Error("gm_login query gmaccount failed.")

		resp.Msg = "account nil"
		c.JSON(http.StatusOK, resp)

		return
	}

	if gmAccount.Password != req.Password {

		logrus.WithFields(logrus.Fields{
			"request":   req,
			"gmaccount": gmAccount,
		}).Error("gm_login password err.")

		resp.Msg = "password err"
		c.JSON(http.StatusOK, resp)

		return
	}

	accessToken := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf(fmt.Sprintf("%d%s%d%s%d", 0, "asdfghjk456", time.Now().Unix(), req.Account, rand.Int())))))
	resp.Token = accessToken

	url := "https://" + req.Channel + ".gs.worldoflove.cn/"
	if req.Channel == "test" {

		url = "http://127.0.0.1:80/"
	}
	err = storage.SetSession(accessToken, req.Account, url)
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"accessToken": accessToken,
			"gmaccount":   req.Account,
			"err":         err,
		}).Error("gm_login setsession err.")

		resp.Msg = "setsession err"
		c.JSON(http.StatusOK, resp)

		return
	}

	c.JSON(http.StatusOK, resp)
}
