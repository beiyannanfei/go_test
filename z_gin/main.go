package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"io"
)

func main() {
	r := gin.Default()

	// Logging to a file.
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)	//将请求日志输出到文件

	l := r.Group("/log")
	l.Use(gin.Logger())
	{
		//	curl "127.0.0.1:8092/log/t1"
		l.GET("/t1", func(context *gin.Context) {
			log.Printf("===============\n")
		})
	}

	r.Run(":8092")
}
