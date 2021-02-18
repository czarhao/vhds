package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 模拟了 cluster 下面有三个 endpoint
func main() {
	go start("8001")
	go start("8002")
	go start("8003")
	select {}
}

func start(port string) {
	router := gin.Default()
	port = ":" + port

	router.GET("/route1", func(c *gin.Context) {
		c.String(http.StatusOK, "endpoint: localhost"+port+", path: /route1")
	})
	router.GET("/route2", func(c *gin.Context) {
		c.String(http.StatusOK, "endpoint: localhost"+port+", path: /route2")
	})
	router.GET("/route3", func(c *gin.Context) {
		c.String(http.StatusOK, "endpoint: localhost"+port+", path: /route3")
	})
	panic(router.Run(port))
}
