package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
)

func main() {
	router := gin.Default()

	router.GET("/get/t1", func(c *gin.Context) {
		fmt.Println("Hello Gin")
		c.JSON(200, gin.H{
			"message": "Hello Gin",
		})
	})

	http.ListenAndServe(":8090", router)
}
