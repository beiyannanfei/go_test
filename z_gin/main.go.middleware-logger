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
	//将请求日志输出到文件(不是认为输出的日志，类似：[GIN] 2018/07/09 - 14:43:53 | 200 [0m|       13.48µs |       127.0.0.1 |  GET [0m    /log/t1)
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

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
