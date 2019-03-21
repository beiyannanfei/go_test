package bind

import (
	"github.com/gin-gonic/gin"
	"github.com/xykong/loveauth/utils"
	"github.com/xykong/loveauth/errors"
	"github.com/xykong/loveauth/storage"
	"fmt"
	"net/http"
)

func init() {
	handlers["/info"] = bindInfo
}

type RequestInfo struct {
	GlobalId int64 `form:"globalId" json:"globalId" binding:"required"`
}

type ResponseInfo struct {
	Body struct {
		Code           int64  `json:"code"`
		Message        string `json:"message"`
		IsBindMobile   bool   `json:"isBindMobile"`
		BindMobileName string `json:"bindMobileName"`
		IsBindWechat   bool   `json:"isBindWechat"`
		BindWechatName string `json:"bindWechatName"`
		IsBindQQ       bool   `json:"isBindQq"`
		BindQQName     string `json:"bindQqName"`
		IsBindWeibo    bool   `json:"isBindWeibo"`
		BindWeiboName  string `json:"bindWeiboName"`
	}
}

func bindInfo(c *gin.Context) {
	var request RequestInfo
	if err := c.BindJSON(&request); err != nil {
		utils.QuickReply(c, errors.ServerFail, "连接服务器异常")
		return
	}

	qq, wechat, mobile, weibo, err := storage.QueryBindInfo(request.GlobalId)
	if err != nil {
		utils.QuickReply(c, errors.ServerFail, "连接服务器异常")
		return
	}

	resp := ResponseInfo{}
	if qq != nil && qq.GlobalId != 0 {
		resp.Body.IsBindQQ = true
		resp.Body.BindQQName = qq.NickName
	}

	if wechat != nil && wechat.GlobalId != 0 {
		resp.Body.IsBindWechat = true
		resp.Body.BindWechatName = wechat.NickName
	}

	if mobile != nil && mobile.GlobalId != 0 {
		resp.Body.IsBindMobile = true
		resp.Body.BindMobileName = fmt.Sprintf("%s****%s", mobile.OpenId[:3], mobile.OpenId[7:])
	}

	if weibo != nil && weibo.GlobalId != 0 {
		resp.Body.IsBindWeibo = true
		resp.Body.BindWeiboName = weibo.NickName
	}

	resp.Body.Code = int64(errors.Ok)
	resp.Body.Message = "获取成功"
	c.JSON(http.StatusOK, resp.Body)
	return
}
