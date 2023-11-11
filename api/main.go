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

	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.tmpl", nil)
	})

	r.GET("/graph", func(c *gin.Context) {
		c.HTML(200, "graph.tmpl", nil)
	})

	r.Run(":" + port)
}
