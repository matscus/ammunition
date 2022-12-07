package handlers

import (
	"ammunition/cache"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func TemporaryInitHandle(c *gin.Context) {
	cache := cache.Temporary{}
	err := c.BindJSON(&cache)
	if err != nil {
		c.JSON(500, gin.H{"Status": "error", "Message": err.Error()})
		return
	}
	err = cache.New()
	if err != nil {
		c.JSON(500, gin.H{"Status": "error", "Message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"Status": "ok", "Message": "cache init"})
}

func TemporaryHandle(c *gin.Context) {
	switch c.Request.Method {
	case http.MethodGet:
		cacheName := c.Query("cache")
		if cacheName == "" {
			c.JSON(400, gin.H{"Status": "error", "Message": "search param \"cache\" is empty"})
			return
		}
		queue := c.Query("queue")
		if cacheName == "" {
			c.JSON(400, gin.H{"Status": "error", "Message": "search param \"queue\" is empty"})
			return
		}
		c.String(200, string(cache.GetTemporaryValue(cacheName, queue)))
	case http.MethodPost:
		cacheName := c.Query("cache")
		if cacheName == "" {
			c.JSON(400, gin.H{"Status": "error", "Message": "search param \"cache\" is empty"})
			return
		}
		key := c.Query("key")
		if key == "" {
			c.JSON(400, gin.H{"Status": "error", "Message": "search param \"key\" is empty"})
			return
		}
		queue := c.Query("queue")
		if cacheName == "" {
			c.JSON(400, gin.H{"Status": "error", "Message": "search param \"queue\" is empty"})
			return
		}
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(500, gin.H{"Status": "error", "Message": err.Error()})
			return
		}
		err = cache.SetTemporaryValue(cacheName, queue, key, body)
		if err != nil {
			c.JSON(500, gin.H{"Status": "error", "Message": err.Error()})
		} else {
			c.JSON(200, gin.H{"Status": "OK", "Message": "Value added"})
		}
	case http.MethodDelete:
		cacheName := c.Query("cache")
		if cacheName =="" {
			c.JSON(400, gin.H{"Status": "error", "Message": "search param \"cache\" is empty"})
		}
		err := cache.DeleteTemporaryCache(cacheName)
		if err != nil {
			c.JSON(500, gin.H{"Status": "error", "Message": err.Error()})
			return
		}
		c.JSON(200, gin.H{"Status": "OK", "Message": "Cache deleted"})
	}
}
