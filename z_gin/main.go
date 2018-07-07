package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
	"os"
	"io"
)

func main() {
	router := gin.Default()

	//最简单路由写法	curl "127.0.0.1:8090/get/t1"
	router.GET("/get/t1", func(c *gin.Context) {
		fmt.Println("Hello Gin")
		c.JSON(200, gin.H{"message": "Hello Gin"})
	})

	//获取路由参数 	curl "127.0.0.1:8090/get/param/bynf"
	router.GET("/get/param/:name", func(c *gin.Context) {
		name := c.Param("name")
		fmt.Printf("Hello %v\n", name)
		c.JSON(200, gin.H{"msg": "OK"})
	})

	//可选路由参数  本路由可以命中： /get/param/aaa/run  /get/param/aaa/
	//如果没有定义路由/get/param/:name 则也可以命中 /get/param/aaa
	//curl "127.0.0.1:8090/get/param/bynf/driver"
	//curl "127.0.0.1:8090/get/param/bynf/"
	router.GET("/get/param/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		fmt.Printf("%s is %s\n", name, action)
		c.JSON(http.StatusOK, gin.H{"msg": "ok"})
	})

	//获取查询参数	curl "127.0.0.1:8090/get/query?lname=yq"
	router.GET("/get/query", func(c *gin.Context) {
		fname := c.DefaultQuery("fname", "wang") //获取参数fname的值，如果没有则取默认值
		lname := c.Query("lname")
		fmt.Printf("Hello %v %v\n", fname, lname)
		c.JSON(200, gin.H{"msg": "OK"})
	})

	//获取post表单类型参数	curl "127.0.0.1:8090/post/form" -d "type=ABCD&msg=EFGH"
	router.POST("/post/form", func(c *gin.Context) {
		mytype := c.DefaultPostForm("type", "alert")
		msg := c.PostForm("msg")
		fmt.Printf("type: %s, msg: %s\n", mytype, msg)
		c.JSON(200, gin.H{"msg": "ok"})
		return
	})

	//query + post form
	//curl "127.0.0.1:8090/post/query/form?id=1234&page=1" -d "name=bynf&age=29"
	router.POST("/post/query/form", func(c *gin.Context) {
		id := c.Query("id")
		page := c.DefaultQuery("page", "0")
		name := c.PostForm("name")
		age := c.DefaultPostForm("age", "18")
		fmt.Printf("id: %s, page: %s, name: %s, age: %s\n", id, page, name, age)
		c.JSON(http.StatusOK, gin.H{"msg": "OK"})
	})

	//路由群组
	v1Group := router.Group("/v1")
	{
		//curl "127.0.0.1:8090/v1/get?q1=aaa&q2=123"
		v1Group.GET("/get", func(c *gin.Context) {
			q1 := c.Query("q1");
			q2 := c.Query("q2");
			fmt.Printf("v1Group get q1: %s, q2: %s", q1, q2)
			c.JSON(200, gin.H{"msg": "OK"})
		})

		//curl "127.0.0.1:8090/v1/post" -d "p1=bbb&p2=111"
		v1Group.POST("/post", func(c *gin.Context) {
			p1 := c.PostForm("p1")
			p2 := c.PostForm("p2")
			fmt.Printf("v1Group psot p1: %s, p2: %s", p1, p2)
			c.JSON(200, gin.H{"msg": "OK"})
		})
	}

	//控制器
	type Login struct {
		User     string `form:"user" json:"user" binding:"required"`
		Password string `form:"password" json:"password" binding:"required"`
	}

	//控制器 绑定JSON的例子
	// curl -X POST "127.0.0.1:8090/post/login/json" -H "Content-Type: application/json" -d "{\"user\":\"aaa\",\"password\":\"123\"}"
	router.POST("/post/login/json", func(c *gin.Context) {
		var json Login
		err := c.BindJSON(&json)
		if nil != err {
			fmt.Printf("bind json error: %v\n", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}

		user := json.User
		pwd := json.Password
		if user != "aaa" || pwd != "123" {
			fmt.Printf("user or pwd illegal user: %s, pwd: %s\n", user, pwd)
			c.JSON(http.StatusUnauthorized, gin.H{"err": "user or pwd illegal"})
			return
		}

		fmt.Printf("login success user: %s, pwd: %s\n", user, pwd)
		c.JSON(http.StatusOK, gin.H{"msg": "login success"})
		return
	})

	//控制器绑定form的例子
	// curl "127.0.0.1:8090/post/login/form" -d "user=aaa&password=123"
	router.POST("/post/login/form", func(c *gin.Context) {
		var form Login
		err := c.Bind(&form)
		if nil != err {
			fmt.Printf("bind form error: %v\n", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			return
		}

		user := form.User
		pwd := form.Password
		if user != "aaa" || pwd != "123" {
			fmt.Printf("user or pwd illegal user: %s, pwd: %s\n", user, pwd)
			c.JSON(http.StatusUnauthorized, gin.H{"err": "user or pwd illegal"})
			return
		}

		fmt.Printf("login success user: %s, pwd: %s\n", user, pwd)
		c.JSON(http.StatusOK, gin.H{"msg": "login success"})
		return
	})

	//上传文件
	//curl -X POST "127.0.0.1:8090/upload" -F "file=@/Users/wyq/workspace/lovejob/README.md" -H "Content-Type: multipart/form-data"
	router.POST("/upload", func(c *gin.Context) {
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, "Bad request")
			return
		}
		//文件的名称
		filename := header.Filename

		fmt.Printf("file: %v\n", file)
		fmt.Printf("header: %v\n", header)

		out, err := os.Create("./t_" + filename)
		if err != nil {
			fmt.Printf("os create file err: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"errmsg": err.Error()})
			return
		}

		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			fmt.Printf("io copy err: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"errmsg": err.Error()})
			return
		}

		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded", filename))
	})



	//http.ListenAndServe(":8090", router)		//两种方式均可以启动服务
	router.Run(":8090")
}
