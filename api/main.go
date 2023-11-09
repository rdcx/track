package main

import (
	"flag"
	"track/api/router"

	"github.com/gin-gonic/gin"
)

func main() {
	var port string
	var pubAddr string

	flag.StringVar(&port, "port", "8080", "port to listen on")
	flag.StringVar(&pubAddr, "pubAddr", "http://localhost:8080", "public address of the server")

	flag.Parse()

	r := router.SetUpRouter(pubAddr)

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	r.Run(":" + port)
}
