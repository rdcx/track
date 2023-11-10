package handler

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"time"
	"track/types"

	"github.com/gin-gonic/gin"
)

var hits types.HitMap
var serverAddr string = "http://localhost:8080"

func Init(addr string) {
	serverAddr = addr
	hits = make(types.HitMap)
}

func decodeUrl(u string) (*url.URL, error) {
	decoded, err := base64.URLEncoding.DecodeString(u)
	fmt.Println(err)
	if err != nil {
		return nil, err
	}
	return url.Parse(string(decoded))
}

func domainExists(hits types.HitMap, domain string) bool {
	if _, ok := hits[domain]; ok {
		return true
	}
	return false
}

func tracker() string {
	return `<script>
	(function() {
		var url = window.location.href;
		fetch("` + serverAddr + `/hit/" + btoa(url));
	})();
</script>
	`
}

func Track(c *gin.Context) {
	domain := c.Param("domain")
	if domainExists(hits, domain) {
		c.JSON(400, types.MessageResponse{
			Success: false,
			Message: "domain already exists",
		})
		return
	}

	hits[domain] = make(map[string][]time.Time)

	c.JSON(200, gin.H{
		"message": "ok",
	})
}

func Hit(c *gin.Context) {
	u := c.Param("url")

	url, err := decodeUrl(u)
	fmt.Println(u)
	if err != nil {
		c.JSON(400, types.MessageResponse{
			Success: false,
			Message: "invalid url",
		})
		return
	}

	domain := url.Hostname()
	if err != nil {
		c.JSON(400, gin.H{
			"message": "invalid domain",
		})
		return
	}

	if _, ok := hits[domain]; !ok {
		c.JSON(400, types.MessageResponse{
			Success: false,
			Message: "domain not found",
		})
		return
	}
	if _, ok := hits[domain][url.String()]; !ok {
		hits[domain][url.String()] = make([]time.Time, 0)
	}

	hits[domain][url.String()] = append(hits[domain][url.String()], time.Now())

	c.JSON(200, types.MessageResponse{
		Success: true,
		Message: "ok",
	})
}

func Hits(c *gin.Context) {
	domain := c.Param("domain")

	if !domainExists(hits, domain) {
		c.JSON(400, types.MessageResponse{
			Success: false,
			Message: "domain not found",
		})
		return
	}

	if _, ok := hits[domain]; !ok {
		c.JSON(400, types.MessageResponse{
			Success: false,
			Message: "domain not found",
		})
		return
	}

	c.JSON(200, gin.H{
		"hits": hits[domain],
	})
}

func Tracker(c *gin.Context) {
	c.String(200, tracker())
}
