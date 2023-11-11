package handler

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/url"
	"time"
	"track/types"

	"github.com/gin-gonic/gin"
)

var hits types.DomainMap
var serverAddr string = "http://localhost:8080"

func Init(addr string) {
	serverAddr = addr
	hits = make(types.DomainMap)
}

func decodeUrl(u string) (*url.URL, error) {
	decoded, err := base64.URLEncoding.DecodeString(u)
	if err != nil {
		return nil, err
	}
	return url.Parse(string(decoded))
}

func randomKey() types.Key {
	runes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, 10)

	for i := range b {
		b[i] = runes[rand.Intn(len(runes))]
	}

	return types.Key(b)
}

func domainExists(hits types.DomainMap, domain types.Domain) bool {
	if _, ok := hits[domain]; ok {
		return true
	}
	return false
}

func tracker(key types.Key) string {
	return `<script>
	(function() {
		var url = window.location.href;
		fetch("` + serverAddr + `/hit/" + btoa(url) + "?k=` + string(key) + `");
	})();
</script>
	`
}

func Track(c *gin.Context) {
	domain := types.Domain(c.Param("domain"))

	if domain == "" {
		c.JSON(400, types.MessageResponse{
			Success: false,
			Message: "invalid domain",
		})
		return
	}

	if !domainExists(hits, domain) {
		hits[domain] = make(types.KeyMap)
	}

	key := randomKey()
	hits[domain][key] = make(types.UrlMap)

	c.JSON(200, types.TrackResponse{
		Message: "ok",
		Success: true,
		Key:     key,
	})
}

func Hit(c *gin.Context) {
	u := c.Param("url")

	key := types.Key(c.Query("k"))

	if key == "" {
		c.JSON(400, types.MessageResponse{
			Success: false,
			Message: "invalid key",
		})
		return
	}

	url, err := decodeUrl(u)

	fmt.Println(u)
	if err != nil {
		c.JSON(400, types.MessageResponse{
			Success: false,
			Message: "invalid url",
		})
		return
	}

	domain := types.Domain(url.Hostname())
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

	if _, ok := hits[domain][key]; !ok {
		c.JSON(400, types.MessageResponse{
			Success: false,
			Message: "key not found",
		})
		return
	}

	urlString := types.Url(url.String())

	if _, ok := hits[domain][key]; !ok {
		hits[domain][key][urlString] = make(types.HitSlice, 0)
	}

	loc := ""

	// if the HTTP_CF_IPCOUNTRY header is set, use that
	if c.GetHeader("HTTP_CF_IPCOUNTRY") != "" {
		loc = c.GetHeader("HTTP_CF_IPCOUNTRY")
	} else {
		loc = "XX"
	}

	hits[domain][key][urlString] = append(hits[domain][key][urlString], types.Hit{
		Loc:  loc,
		Time: time.Now(),
	})

	c.JSON(200, types.MessageResponse{
		Success: true,
		Message: "ok",
	})
}

func Hits(c *gin.Context) {
	domain := types.Domain(c.Param("domain"))
	key := types.Key(c.Query("k"))

	if !domainExists(hits, domain) {
		c.JSON(400, types.MessageResponse{
			Success: false,
			Message: "domain not found",
		})
		return
	}

	if _, ok := hits[domain][key]; !ok {
		c.JSON(400, types.MessageResponse{
			Success: false,
			Message: "key not found",
		})
		return
	}

	c.JSON(200, gin.H{
		"hits": hits[domain][key],
	})
}

func Tracker(c *gin.Context) {
	k := types.Key(c.Query("k"))
	if k == "" {
		c.JSON(400, types.MessageResponse{
			Success: false,
			Message: "invalid key",
		})

		return
	}
	c.String(200, tracker(k))
}
