package router

import (
	"time"
	"track/api/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetUpRouter(pubAddr string) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	handler.Init(pubAddr)

	r.POST("/track/:domain", handler.Track)
	r.GET("/hit/:url", handler.Hit)
	r.GET("/hits/:domain", handler.Hits)
	r.GET("/tracker", handler.Tracker)

	return r
}
