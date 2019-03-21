package verify

import "github.com/gin-gonic/gin"

var handlers = make(map[string]gin.HandlerFunc)

func Start(group *gin.RouterGroup) {

	for key, value := range handlers {
		group.POST(key, value)
	}
}
