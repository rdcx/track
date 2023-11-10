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

func randomKey() types.Key {
	runes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, 10)

	for i := range b {
		b[i] = runes[rand.Intn(len(runes))]
	}

	return types.Key(b)
}

func domainExists(hits types.HitMap, domain types.Domain) bool {
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
		hits[domain] = make(map[types.Key]map[types.Url][]time.Time)
	}

	key := randomKey()
	hits[domain][key] = make(map[types.Url][]time.Time)

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
		hits[domain][key][urlString] = make([]time.Time, 0)
	}

	hits[domain][key][urlString] = append(hits[domain][key][urlString], time.Now())

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
