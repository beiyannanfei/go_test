//+build windows

package graceful

import (
	"net"
	"github.com/sirupsen/logrus"
	"net/http"
	"github.com/gin-gonic/gin"
)

var (
	server   *http.Server
	listener net.Listener
)

func ListenAndServe(addr string, handler *gin.Engine) {

	err := handler.Run(addr)
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"listen": addr,
			"error":  err,
		}).Error("Start loveauth service failed.")
	}
}
