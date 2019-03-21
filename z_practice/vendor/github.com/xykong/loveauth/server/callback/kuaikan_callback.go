package callback

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/storage"
)

func init() {
	getHandlers["/kuaikan/callback"] = dealKuaikan
}

func dealKuaikan(c *gin.Context) {
	idfa := c.Query("idfa")
	if idfa == "" {
		c.Status(400)
		c.Writer.WriteString("idfa not exists.")
		return
	}

	err := storage.CreateFirstIdfa(idfa)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"idfa": idfa,
			"err":  err,
		}).Error("save idfa failed.")

		c.Status(400)
		c.Writer.WriteString("save idfa failed.")
		return
	}

	c.Writer.WriteString("success")
	return
}
