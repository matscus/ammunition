package handlers

import (
	"net/http"
	"strconv"

	"ammunition/cache"
	"ammunition/config"

	"github.com/gin-gonic/gin"
)

//Manage func from create(method post) or update(method put) or delete (method delete) datapool
func PersistHandle(c *gin.Context) {
	project := c.Query("project")
	if project == "" {
		c.JSON(400, gin.H{"Status": "error", "Message": "project is empty"})
		return
	}
	name := c.Query("name")
	if name == "" {
		c.JSON(400, gin.H{"Status": "error", "Message": "name is empty"})
		return
	}
	switch c.Request.Method {
	case http.MethodGet:
		res, err := cache.PersistedPool{Project: project, Name: name}.GetValue()
		if err != nil {
			c.JSON(500, gin.H{"Status": "error", "Message": err.Error()})
			return
		}
		if res == "" {
			c.JSON(200, gin.H{"Status": "OK", "Message": "chanel is empty"})
			return
		}
		c.String(200, res)
	case http.MethodPost:
		pool := cache.PersistedPool{Project: project, Name: name}
		bufferLen := c.Query("bufferlen")
		if bufferLen == "" {
			pool.BufferLen = config.Config.Persist.BufferLen
		} else {
			pool.BufferLen, _ = strconv.Atoi(bufferLen)
		}
		workers := c.Query("workers")
		if workers == "" {
			pool.Workers = config.Config.Persist.Worker
		} else {
			pool.Workers, _ = strconv.Atoi(workers)
		}
		fileHeader, err := c.FormFile("csvfile")
		if err != nil {
			c.JSON(400, gin.H{"Status": "error", "Message": "csvfile is empty"})
			return
		}
		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(500, gin.H{"Status": "error", "Message": err.Error()})
			return
		}
		err = pool.Create(&file)
		if err != nil {
			c.JSON(500, gin.H{"Status": "error", "Message": err.Error()})
			return
		}
		c.JSON(200, gin.H{"Status": "OK", "Message": "Pool created"})
	case http.MethodPut:
		pool := cache.PersistedPool{Project: project, Name: name}
		bufferLen := c.Query("bufferlen")
		if bufferLen == "" {
			pool.BufferLen = config.Config.Persist.BufferLen
		} else {
			pool.BufferLen, _ = strconv.Atoi(bufferLen)
		}
		workers := c.Query("workers")
		if workers == "" {
			pool.Workers = config.Config.Persist.Worker
		} else {
			pool.Workers, _ = strconv.Atoi(workers)
		}
		action := c.Query("action")
		if action == "" {
			c.JSON(400, gin.H{"error": "action is empty"})
			return
		}
		fileHeader, err := c.FormFile("csvfile")
		if err != nil {
			c.JSON(400, gin.H{"error": "csvfile is empty"})
			return
		}
		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(500, gin.H{"Status": "error", "Message": err.Error()})
			return
		}
		switch action {
		case "update":
			err = pool.Update(&file)
			if err != nil {
				c.JSON(500, gin.H{"Status": "error", "Message": err.Error()})
				return
			}
			c.JSON(200, gin.H{"Status": "OK", "Message": "Pool updated"})
		case "add":
			err = pool.AddValues(&file)
			if err != nil {
				c.JSON(500, gin.H{"Status": "error", "Message": err.Error()})
				return
			}
			c.JSON(200, gin.H{"Status": "OK", "Message": "value added"})
		default:
			c.JSON(400, gin.H{"Status": "error", "Message": "Invalid value action. Possible values Update or Add"})
			return
		}
	case http.MethodDelete:
		pool := cache.PersistedPool{Project: project, Name: name}
		err := pool.Delete()
		if err != nil {
			c.JSON(500, gin.H{"Status": "error", "Message": err.Error()})
			return
		}
		c.JSON(200, gin.H{"Status": "OK", "Message": "Pool deleted"})
	}
}
