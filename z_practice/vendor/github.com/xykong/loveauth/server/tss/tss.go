package tss

import (
	"github.com/gin-gonic/gin"
	"github.com/xykong/loveauth/services/tss"
)

var handlers = make(map[string]gin.HandlerFunc)

func Start(group *gin.RouterGroup) {

	tss.Start()

	for key, value := range handlers {
		group.POST(key, value)
	}
}
