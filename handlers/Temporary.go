package handlers

import (
	"ammunition/cache"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func TemporaryHandle(c *gin.Context) {
	switch c.Request.Method {
	case http.MethodGet:
		key := c.Query("key")
		if key != "" {
			switch key {
			case "iterator":
				res := cache.GetTemporaryIteratorValue()
				c.String(200, string(res))
			default:
				res, err := cache.GetTemporaryValue(key)
				if err != nil {
					c.JSON(500, gin.H{"Status": "error", "Message": err.Error()})
					return
				}
				c.String(200, string(res))
			}
		} else {
			c.JSON(400, gin.H{"Status": "error", "Message": "key is empty"})
		}

	case http.MethodPost:
		key := c.Query("key")
		if key != "" {
			body, err := ioutil.ReadAll(c.Request.Body)
			if err != nil {
				c.JSON(500, gin.H{"Status": "error", "Message": err.Error()})
				return
			}
			err = cache.SetTemporaryValue(key, body)
			if err != nil {
				c.JSON(500, gin.H{"Status": "error", "Message": err.Error()})
			} else {
				c.JSON(200, gin.H{"Status": "OK", "Message": "Value added"})
			}
		} else {
			c.JSON(400, gin.H{"Status": "error", "Message": "key is empty"})
		}
	case http.MethodDelete:
		key := c.Query("key")
		if key != "" {
			switch key {
			case "reset":
				err := cache.ResetTemporaryCache()
				if err != nil {
					c.JSON(500, gin.H{"Status": "error", "Message": err.Error()})
					return
				}
				c.JSON(200, gin.H{"Status": "OK", "Message": "Pool reseted"})
			default:
				err := cache.DeleteTemporaryValue(key)
				if err != nil {
					c.JSON(500, gin.H{"Status": "error", "Message": err.Error()})
					return
				}
				c.JSON(200, gin.H{"Status": "OK", "Message": "Value deleted"})
			}
		} else {
			c.JSON(400, gin.H{"Status": "error", "Message": "key is empty"})
		}

	}
}
