package bind

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/storage/model"
	"github.com/xykong/loveauth/utils"
	"github.com/xykong/loveauth/errors"
	"github.com/xykong/loveauth/storage"
	"net/http"
)

var handlers = make(map[string]gin.HandlerFunc)
var getHandlers = make(map[string]gin.HandlerFunc)

type Supplier interface {
	Name() string
	BindCheck(globalId int64, user *model.DoAuthRequest) error
	BindVerify(c *gin.Context, user *model.DoAuthRequest) error
	CreateBind(globalId int64, user *model.DoAuthRequest) error
	UndoBind(globalId int64, req *RequestUnBind) error
}

type Suppliers map[string]Supplier

var usedSuppliers = Suppliers{}

func Use(aps ...Supplier) {
	for _, Supplier := range aps {
		if usedSuppliers[Supplier.Name()] != nil {
			logrus.WithFields(logrus.Fields{
				"provider": Supplier.Name(),
			}).Warn("Supplier replaced.")
		}

		usedSuppliers[Supplier.Name()] = Supplier
	}
}

type Binding struct {
	Supplier Supplier
}

type UnBinding struct {
	Supplier Supplier
}

type RequestBind struct {
	Token         string              `form:"Token" json:"Token" binding:"required"`
	DoAuthRequest model.DoAuthRequest `form:"DoAuthRequest" json:"DoAuthRequest" binding:"required"`
}

type ResponseBind struct {
	Body struct {
		Code    int64  `json:"code"`
		Message string `json:"message"`
	}
}

type RequestUnBind struct {
	Token            string `form:"Token" json:"Token" binding:"required"`
	OpenId           string `form:"OpenId" json:"OpenId"`
	AccessTokenValue string `form:"AccessTokenValue" json:"AccessTokenValue"`
}

type ResponseUnBind struct {
	Body struct {
		Code    int64  `json:"code"`
		Message string `json:"message"`
	}
}

func (o *UnBinding) UnBind(c *gin.Context) {
	var request RequestUnBind
	if err := c.BindJSON(&request); err != nil {
		utils.QuickReply(c, errors.ServerFail, "连接服务器异常")
		return
	}

	//1 验证登录信息是否有效
	if request.Token == "" {
		utils.QuickReply(c, errors.ServerFail, "连接服务器异常")
		return
	}

	tokenRecord, err := storage.QueryAccessToken(request.Token)
	if err != nil {
		utils.QuickReply(c, errors.ServerFail, "连接服务器异常")
		return
	}

	if tokenRecord.GlobalId == 0 {
		utils.QuickReply(c, errors.ServerFail, "连接服务器异常")
		return
	}

	//2 解绑
	err = o.Supplier.UndoBind(tokenRecord.GlobalId, &request)
	if err != nil {
		ec, ok := err.(*errors.Type)
		if !ok { //系统错误
			utils.QuickReply(c, errors.BindCheckErr, "绑定校验失败")
			return
		}

		//逻辑类错误
		utils.QuickReply(c, ec.Code, ec.Message)
		return
	}

	info, err := storage.QueryProfile(tokenRecord.GlobalId)
	if nil == err {
		unBindVendor := o.Supplier.Name()
		LogBind(tokenRecord.GlobalId, info.Auth, unBindVendor, Unbind)
	}

	resp := ResponseUnBind{}
	resp.Body.Code = int64(errors.Ok)
	resp.Body.Message = "解绑成功"
	c.JSON(http.StatusOK, resp.Body)
	return
}

func (o *Binding) Bind(c *gin.Context) {
	var request RequestBind
	if err := c.BindJSON(&request); err != nil {
		utils.QuickReply(c, errors.ServerFail, "连接服务器异常")
		return
	}

	//1 验证登录信息是否有效
	if request.Token == "" {
		utils.QuickReply(c, errors.ServerFail, "连接服务器异常")
		return
	}

	tokenRecord, err := storage.QueryAccessToken(request.Token)
	if err != nil {
		utils.QuickReply(c, errors.QueryAccessTokenFailed, "登录信息已失效",)
		return
	}

	if tokenRecord.GlobalId == 0 {
		utils.QuickReply(c, errors.ServerFail, "连接服务器异常")
		return
	}

	//2 验证是否可以绑定
	err = o.Supplier.BindCheck(tokenRecord.GlobalId, &request.DoAuthRequest)
	if err != nil {
		ec, ok := err.(*errors.Type)
		if !ok { //系统错误
			utils.QuickReply(c, errors.BindCheckErr, "绑定校验失败")
			return
		}

		//逻辑类错误
		utils.QuickReply(c, ec.Code, ec.Message)
		return
	}

	//3 鉴权
	err = o.Supplier.BindVerify(c, &request.DoAuthRequest)
	if err != nil {
		ec, ok := err.(*errors.Type)
		if !ok { //系统错误
			utils.QuickReply(c, errors.BindCheckErr, "绑定校验失败")
			return
		}

		//逻辑类错误
		utils.QuickReply(c, ec.Code, ec.Message)
		return
	}

	//4 创建绑定信息
	err = o.Supplier.CreateBind(tokenRecord.GlobalId, &request.DoAuthRequest)
	if err != nil {
		utils.QuickReply(c, errors.ServerFail, "连接服务器异常")
		return
	}

	//5 删除设备绑定信息
	var authDevice model.AuthDevice
	err = storage.DeleteVendorByGlobalId(tokenRecord.GlobalId, &authDevice)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":      err,
			"GlobalId": tokenRecord.GlobalId,
		}).Error("delete device bind failed.")
		utils.QuickReply(c, errors.ServerFail, "连接服务器异常")
		return
	}

	bindVendor := o.Supplier.Name()
	LogBind(tokenRecord.GlobalId, request.DoAuthRequest, bindVendor, Bind)

	resp := ResponseBind{}
	resp.Body.Code = int64(errors.Ok)
	resp.Body.Message = "绑定成功"
	c.JSON(http.StatusOK, resp.Body)
	return
}

func Start(group *gin.RouterGroup) {
	for key, value := range handlers {
		group.POST(key, value)
	}

	for key, value := range getHandlers {
		group.GET(key, value)
	}

	for _, s := range usedSuppliers {
		var Supplier = s
		group.POST("/"+s.Name(), func(c *gin.Context) { //绑定
			var binding = Binding{Supplier}
			binding.Bind(c)
		})

		group.POST("/"+s.Name()+"/undo", func(c *gin.Context) { //解绑
			var unbind = UnBinding{Supplier}
			unbind.UnBind(c)
		})
	}
}
