package router

import (
	"track/api/handler"

	"github.com/gin-gonic/gin"
)

func SetUpRouter(pubAddr string) *gin.Engine {
	r := gin.Default()

	handler.Init(pubAddr)

	r.POST("/track/:domain", handler.Track)
	r.GET("/hit/:url", handler.Hit)
	r.GET("/hits/:domain", handler.Hits)
	r.GET("/tracker", handler.Tracker)

	return r
}
