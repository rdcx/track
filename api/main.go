package main

import (
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

var serverAddr string = "http://localhost:8080/"
var domains []string
var hits map[string]map[string]time.Time

func uniqueAppend(slice []string, s string) []string {
	for _, v := range slice {
		if v == s {
			return slice
		}
	}
	return append(slice, s)
}

func domainExists(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

func parseDomain(u string) (string, error) {
	ul, err := url.Parse(u)

	if err != nil {
		return "", err
	}

	return ul.Hostname(), nil
}

func tracker() string {
	return `
		<script>
			(function() {
				var url = window.location.href;
				fetch("` + serverAddr + `/hit/" + encodeURIComponent(url));
			})();
		</script>
	`
}

func main() {
	r := gin.Default()

	r.POST("/track/{domain}", func(c *gin.Context) {
		domain := c.Param("domain")

		if domainExists(domains, domain) {
			c.JSON(400, gin.H{
				"message": "domain already exists",
			})
			return
		}

		domains = uniqueAppend(domains, domain)

		hits[domain] = make(map[string]time.Time)

		c.JSON(200, gin.H{
			"message": "ok",
		})
	})

	r.GET("/hit/{url}", func(c *gin.Context) {
		u := c.Param("url")

		unesc, err := url.QueryUnescape(u)

		if err != nil {
			c.JSON(400, gin.H{
				"message": "invalid url",
			})
			return
		}

		domain, err := parseDomain(unesc)

		if err != nil {
			c.JSON(400, gin.H{
				"message": "invalid domain",
			})
			return
		}

		if !domainExists(domains, domain) {
			c.JSON(400, gin.H{
				"message": "domain not found",
			})
			return
		}

		if _, ok := hits[domain]; !ok {
			c.JSON(400, gin.H{
				"message": "domain not found",
			})
			return
		}

		hits[domain][time.Now().String()] = time.Now()

		c.JSON(200, gin.H{
			"message": "ok",
		})
	})

	r.GET("/hits/{domain}", func(c *gin.Context) {
		domain := c.Param("domain")

		if !domainExists(domains, domain) {
			c.JSON(400, gin.H{
				"message": "domain not found",
			})
			return
		}

		if _, ok := hits[domain]; !ok {
			c.JSON(400, gin.H{
				"message": "domain not found",
			})
			return
		}

		c.JSON(200, gin.H{
			"hits": hits[domain],
		})
	})

	r.GET("/tracker", func(c *gin.Context) {
		c.String(200, tracker())
	})

	r.Run()
}
